/*
 * Copyright 2015-2016 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/openwhisk/openwhisk-client-go/whisk"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

// ServerlessBinaryCommand is the CLI name to run serverless
const ServerlessBinaryCommand = "serverless"

// ManifestProvider is a provider description in the manifest
type ManifestProvider struct {
	Name    string
	Runtime string
}

// Manifest is the main manifest file
type Manifest struct {
	Service  string
	Provider ManifestProvider
}

// ServerlessErr records errors from the Serverless binary
type ServerlessErr struct {
	Msg string
}

func (e *ServerlessErr) Error() string {
	return e.Msg
}

// Check is a util function to panic when there is an error.
func Check(e error) {

	if err := recover(); err != nil {
		log.Printf("runtime panic : %v", e)
	}

	if e != nil {
		erro := errors.New("Error happened during execution, please type wskdeploy -help for help messages, now exit.")
		log.Printf("error: %v", erro)
		panic(e)
	}

}

type URLReader struct {
}

func (urlReader *URLReader) ReadUrl(url string) (content []byte, err error) {
	resp, err := http.Get(url)
	Check(err)
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	Check(err)
	return b, nil
}

type LocalReader struct {
}

func (localReader *LocalReader) ReadLocal(path string) (content []byte, err error) {
	cont, err := ioutil.ReadFile(path)
	Check(err)
	return cont, nil
}

// agnostic util reader to fetch content from web or local path or potentially other places.
type ContentReader struct {
	URLReader
	LocalReader
}

func GetHomeDirectory() string {
	usr, err := user.Current()
	Check(err)

	return usr.HomeDir
}

// Utility to convert hostname to URL object
func GetURLBase(host string) (*url.URL, error) {

	urlBase := fmt.Sprintf("%s/api/", host)
	url, err := url.Parse(urlBase)

	if len(url.Scheme) == 0 || len(url.Host) == 0 {
		urlBase = fmt.Sprintf("https://%s/api/", host)
		url, err = url.Parse(urlBase)
	}

	return url, err
}

func ReadProps(path string) (map[string]string, error) {

	props := map[string]string{}

	file, err := os.Open(path)
	if err != nil {
		// If file does not exist, just return props
		fmt.Printf("Unable to read whisk properties file '%s' (file open error: %s); falling back to default properties\n", path, err)
		return props, nil
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	props = map[string]string{}
	for _, line := range lines {
		kv := strings.Split(line, "=")
		if len(kv) != 2 {
			// Invalid format; skip
			continue
		}
		props[kv[0]] = kv[1]
	}

	return props, nil

}

func IsDirectory(filePath string) bool {
	f, err := os.Open(filePath)
	Check(err)

	defer f.Close()

	fi, err := f.Stat()
	Check(err)

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	default:
		return false
	}
}

func CreateActionFromFile(manipath, filePath string) (*whisk.Action, error) {
	ext := path.Ext(filePath)
	baseName := path.Base(filePath)
	name := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	action := new(whisk.Action)
	//better refactor this
	splitmanipath := strings.Split(manipath, string(os.PathSeparator))
	filePath = strings.TrimRight(manipath, splitmanipath[len(splitmanipath)-1]) + filePath
	fmt.Println(filePath)
	// process source code files
	if ext == ".swift" || ext == ".js" || ext == ".py" {

		kind := "nodejs:default"

		switch ext {
		case ".swift":
			kind = "swift:default"
		case ".js":
			kind = "nodejs:default"
		case ".py":
			kind = "python"
		}

		dat, err := ioutil.ReadFile(filePath)
		Check(err)
		action.Exec = new(whisk.Exec)
		action.Exec.Code = string(dat)
		action.Exec.Kind = kind
		action.Name = name
		action.Publish = false
		return action, nil
	}
	return nil, nil
}
