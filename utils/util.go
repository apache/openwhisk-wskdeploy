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

// ActionRecord is a container to keep track of
// a whisk action struct and a location filepath we use to
// map files and manifest declared actions
type ActionRecord struct {
	Action   *whisk.Action
	Filepath string
}

func NewClient(proppath string, deploymentPath string) (*whisk.Client, *whisk.Config) {
	var clientConfig *whisk.Config
	configs, err := LoadConfiguration(proppath)
	Check(err)
	//we need to get Apihost from property file which currently not defined in sample deployment file.
	baseURL, err := GetURLBase(configs[1])
	Check(err)
	if deploymentPath != "" {
		mm := NewYAMLParser()
		deployment := mm.ParseDeployment(deploymentPath)
		// We get the first package from the sample deployment file.
		pkg := deployment.Application.GetPackageList()[0]
		clientConfig = &whisk.Config{
			AuthToken: pkg.Credential, //Authtoken
			Namespace: pkg.Namespace,  //Namespace
			BaseURL:   baseURL,
			Version:   "v1",
			Insecure:  true,
		}

	} else {
		clientConfig = &whisk.Config{
			AuthToken: configs[2], //Authtoken
			Namespace: configs[0], //Namespace
			BaseURL:   baseURL,
			Version:   "v1",
			Insecure:  true, // true if you want to ignore certificate signing

		}

	}

	// Setup network client
	client, err := whisk.NewClient(http.DefaultClient, clientConfig)
	Check(err)
	return client, clientConfig
}

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
	defer func() {
		if err := recover(); err != nil {
			log.Printf("runtime panic : %v", err)
		}
	}()

	if e != nil {
		log.Printf("%v", e)
		erro := errors.New("Error happened during execution, please type 'wskdeploy -h' for help messages.")
		log.Printf("%v", erro)
		os.Exit(1)

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
	// If the action is not supported, we better to return an error.
	return nil, errors.New("Unsupported action type.")
}

func GetBoolFromString(value string) (bool, error) {
	if strings.EqualFold(value, "true") || strings.EqualFold(value, "t") || strings.EqualFold(value, "yes") || strings.EqualFold(value, "y") {
		return true, nil
	} else if strings.EqualFold(value, "false") || strings.EqualFold(value, "f") || strings.EqualFold(value, "no") || strings.EqualFold(value, "n") {
		return false, nil
	}

	return false, fmt.Errorf("Value %s not a valid comparison for boolean", value)
}

// Load configuration will load properties from a file
func LoadConfiguration(propPath string) ([]string, error) {
	props, err := ReadProps(propPath)
	Check(err)
	Namespace := props["NAMESPACE"]
	Apihost := props["APIHOST"]
	Authtoken := props["AUTH"]
	return []string{Namespace, Apihost, Authtoken}, nil
}
