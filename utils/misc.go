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
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/hokaccha/go-prettyjson"
	"io/ioutil"
	"net/http"
	"path"
)

const (
	DEFAULT_HTTP_TIMEOUT = 30
	DEFAULT_PROJECT_PATH = "."
	// name of manifest and deployment files
	ManifestFileNameYaml   = "manifest.yaml"
	ManifestFileNameYml    = "manifest.yml"
	DeploymentFileNameYaml = "deployment.yaml"
	DeploymentFileNameYml  = "deployment.yml"
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

func GetHomeDirectory() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}

	return usr.HomeDir
}

// Potentially complex structures(such as DeploymentProject, DeploymentPackage)
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

func PrettyJSON(j interface{}) (string, error) {
	formatter := prettyjson.NewFormatter()
	bytes, err := formatter.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
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

// Test if a string
func isValidEnvironmentVar(value string) bool {

	// A valid Env. variable should start with or contain '$' (dollar) char.
	//
	// If the value is a single Env. variable, it should start with a '$' (dollar) char
	// and have at least 1 additional character after it, e.g. $ENV_VAR
	// If the value is a concatenation of a string and a Env. variable, it should contain '$' (dollar)
	// and have a string following which is surrounded with '{' and '}', e.g. xxx${ENV_VAR}xxx.
	if value != "" && strings.HasPrefix(value, "$") && len(value) > 1 {
		return true
	}
	if value != "" && strings.Contains(value, "${") && strings.Count(value, "{") == strings.Count(value, "}") {
		return true
	}
	return false
}

// Get the env variable value by key.
// Get the env variable if the key is start by $
func GetEnvVar(key interface{}) interface{} {
	// Assure the key itself is not nil
	if key == nil {
		return nil
	}

	if reflect.TypeOf(key).String() == "string" {
		keystr := key.(string)
		if isValidEnvironmentVar(keystr) {
			// retrieve the value of the env. var. from the host system.
			var thisValue string
			//split the string with ${}
			//test if the substr is a environment var
			//if it is, replace it with the value
			f := func(c rune) bool {
				return c == '$' || c == '{' || c == '}'
			}
			for _, substr := range strings.FieldsFunc(keystr, f) {
				//if the substr is a $ENV_VAR
				if strings.Contains(keystr, "$"+substr) {
					thisValue = os.Getenv(substr)
					if thisValue == "" {
						PrintOpenWhiskOutputln("WARNING: Missing Environment Variable " + substr + ".")
					}
					keystr = strings.Replace(keystr, "$"+substr, thisValue, -1)
					//if the substr is a ${ENV_VAR}
				} else if strings.Contains(keystr, "${"+substr+"}") {
					thisValue = os.Getenv(substr)
					if thisValue == "" {
						PrintOpenWhiskOutputln("WARNING: Missing Environment Variable " + substr + ".")
					}
					keystr = strings.Replace(keystr, "${"+substr+"}", thisValue, -1)
				}
			}
			return keystr
		}

		// The key was not a valid env. variable, simply return it as the value itself (of type string)
		return keystr
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

func NewZipWritter(src, des string) *ZipWritter {
	zw := &ZipWritter{src: src, des: des}
	return zw
}

type ZipWritter struct {
	src        string
	des        string
	zipWritter *zip.Writer
}

func (zw *ZipWritter) zipFile(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !f.Mode().IsRegular() || f.Size() == 0 {
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fileName := strings.TrimPrefix(path, zw.src+"/")
	wr, err := zw.zipWritter.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(wr, file)
	if err != nil {
		return err
	}
	return nil
}

func (zw *ZipWritter) Zip() error {
	// create zip file
	zipFile, err := os.Create(zw.des)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	zw.zipWritter = zip.NewWriter(zipFile)
	err = filepath.Walk(zw.src, zw.zipFile)
	if err != nil {
		return nil
	}
	err = zw.zipWritter.Close()
	if err != nil {
		return err
	}
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
		if err != nil {
			return nil, err
		}
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

func deleteKey(key string, keyValueArr whisk.KeyValueArr) whisk.KeyValueArr {
	for i := 0; i < len(keyValueArr); i++ {
		if keyValueArr[i].Key == key {
			keyValueArr = append(keyValueArr[:i], keyValueArr[i+1:]...)
			break
		}
	}

	return keyValueArr
}

func addKeyValue(key string, value interface{}, keyValueArr whisk.KeyValueArr) whisk.KeyValueArr {
	keyValue := whisk.KeyValue{
		Key:   key,
		Value: value,
	}

	return append(keyValueArr, keyValue)
}

func GetManifestFilePath(projectPath string) string {
	if _, err := os.Stat(path.Join(projectPath, ManifestFileNameYaml)); err == nil {
		return path.Join(projectPath, ManifestFileNameYaml)
	} else if _, err := os.Stat(path.Join(projectPath, ManifestFileNameYml)); err == nil {
		return path.Join(projectPath, ManifestFileNameYml)
	} else {
		return ""
	}
}

func GetDeploymentFilePath(projectPath string) string {
	if _, err := os.Stat(path.Join(projectPath, DeploymentFileNameYaml)); err == nil {
		return path.Join(projectPath, DeploymentFileNameYaml)
	} else if _, err := os.Stat(path.Join(projectPath, DeploymentFileNameYml)); err == nil {
		return path.Join(projectPath, DeploymentFileNameYml)
	} else {
		return ""
	}
}

// agnostic util reader to fetch content from web or local path or potentially other places.
type ContentReader struct {
	URLReader
	LocalReader
}

type URLReader struct {
}

func (urlReader *URLReader) ReadUrl(url string) (content []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return content, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return content, err
	} else {
		defer resp.Body.Close()
	}
	return b, nil
}

type LocalReader struct {
}

func (localReader *LocalReader) ReadLocal(path string) ([]byte, error) {
	cont, err := ioutil.ReadFile(path)
	return cont, err
}

func Read(url string) ([]byte, error) {
	if strings.HasPrefix(url, "http") {
		return new(ContentReader).URLReader.ReadUrl(url)
	} else {
		return new(ContentReader).LocalReader.ReadLocal(url)
	}
}
