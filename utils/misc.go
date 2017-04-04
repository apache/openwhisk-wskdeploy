/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"archive/zip"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hokaccha/go-prettyjson"
	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/openwhisk-wskdeploy/wski18n"
	"io"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
)

// ActionRecord is a container to keep track of
// a whisk action struct and a location filepath we use to
// map files and manifest declared actions
type ActionRecord struct {
	Action      *whisk.Action
	Packagename string
	Filepath    string
}

type TriggerRecord struct {
	Trigger     *whisk.Trigger
	Packagename string
}

type RuleRecord struct {
	Rule        *whisk.Rule
	Packagename string
}

// Utility to convert hostname to URL object
func GetURLBase(host string) (*url.URL, error) {

	urlBase := fmt.Sprintf("%s/api", host)
	url, err := url.Parse(urlBase)

	if len(url.Scheme) == 0 || len(url.Host) == 0 {
		urlBase = fmt.Sprintf("https://%s/api", host)
		url, err = url.Parse(urlBase)
	}

	return url, err
}

func GetHomeDirectory() string {
	usr, err := user.Current()
	Check(err)

	return usr.HomeDir
}

// Potentially complex structures(such as DeploymentApplication, DeploymentPackage)
// could implement those interface which is convenient for put, get subtract in
// containers etc.
type Comparable interface {
	HashCode() uint32
	Equals() bool
}

func IsFeedAction(trigger *whisk.Trigger) (string, bool) {
	for _, annotation := range trigger.Annotations {
		if annotation.Key == "feed" {
			return annotation.Value.(string), true
		}
	}

	return "", false
}

func IsJSON(s string) (interface{}, bool) {
	var js interface{}
	if json.Unmarshal([]byte(s), &js) == nil {
		return js, true
	}
	return nil, false

}

func PrettyJSON(j interface{}) string {
	formatter := prettyjson.NewFormatter()
	bytes, err := formatter.Marshal(j)
	Check(err)
	return string(bytes)
}

// Common utilities

// Prompt for user input
func Ask(reader *bufio.Reader, question string, def string) string {
	fmt.Print(question + " (" + def + "): ")
	answer, _ := reader.ReadString('\n')
	len := len(answer)
	if len == 1 {
		return def
	}
	return answer[:len-1]
}

// Get the env variable value by key.
// Get the env variable if the key is start by $
func GetEnvVar(key interface{}) interface{} {
	if reflect.TypeOf(key).String() == "string" {
		if strings.HasPrefix(key.(string), "$") {
			envkey := strings.Split(key.(string), "$")[1]
			value := os.Getenv(envkey)
			if value != "" {
				return value
			}
			return envkey
		}
		return key.(string)
	}
	return key
}

var kindToJSON []string = []string{"", "boolean", "integer", "integer", "integer", "integer", "integer", "integer", "integer", "integer",
	"integer", "integer", "integer", "number", "number", "number", "number", "array", "", "", "", "object", "", "", "string", "", ""}

// Gets JSON type name
func GetJSONType(j interface{}) string {
	fmt.Print(reflect.TypeOf(j).Kind())
	return kindToJSON[reflect.TypeOf(j).Kind()]
}

func CreateZip(filename string, files []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	zipwriter := zip.NewWriter(file)
	defer zipwriter.Close()
	for _, name := range files {
		if err := writeFileToZip(zipwriter, name); err != nil {
			return err
		}
	}
	return nil
}

func writeFileToZip(zipwriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	finfo, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(finfo)
	if err != nil {
		return err
	}
	//add some filter logic if necessary
	//filter(file)
	writer, err := zipwriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}

func filter(filename string) interface{} {
	//To do
	return nil
}

// below codes is from wsk cli with tiny adjusts.
func GetExec(artifact string, kind string, isDocker bool, mainEntry string) (*whisk.Exec, error) {
	var err error
	var code string
	var content []byte
	var exec *whisk.Exec

	ext := filepath.Ext(artifact)
	exec = new(whisk.Exec)

	if !isDocker || ext == ".zip" {
		content, err = new(ContentReader).ReadLocal(artifact)
		Check(err)
		code = string(content)
		exec.Code = &code
	}

	if len(kind) > 0 {
		exec.Kind = kind
	} else if isDocker {
		exec.Kind = "blackbox"
		if ext != ".zip" {
			exec.Image = artifact
		} else {
			exec.Image = "openwhisk/dockerskeleton"
		}
	} else if ext == ".swift" {
		exec.Kind = "swift:default"
	} else if ext == ".js" {
		exec.Kind = "nodejs:default"
	} else if ext == ".py" {
		exec.Kind = "python:default"
	} else if ext == ".jar" {
		exec.Kind = "java:default"
		exec.Code = nil
	} else {
		if ext == ".zip" {
			return nil, zipKindError()
		} else {
			return nil, extensionError(ext)
		}
	}

	// Error if entry point is not specified for Java
	if len(mainEntry) != 0 {
		exec.Main = mainEntry
	} else {
		if exec.Kind == "java" {
			return nil, javaEntryError()
		}
	}

	// Base64 encode the zip file content
	if ext == ".zip" {
		code = base64.StdEncoding.EncodeToString([]byte(code))
		exec.Code = &code
	}

	return exec, nil
}

func zipKindError() error {
	errMsg := wski18n.T("creating an action from a .zip artifact requires specifying the action kind explicitly")

	return errors.New(errMsg)
}

func extensionError(extension string) error {
	errMsg := wski18n.T(
		"'{{.name}}' is not a supported action runtime",
		map[string]interface{}{
			"name": extension,
		})

	return errors.New(errMsg)
}

func javaEntryError() error {
	errMsg := wski18n.T("Java actions require --main to specify the fully-qualified name of the main class")

	return errors.New(errMsg)
}
