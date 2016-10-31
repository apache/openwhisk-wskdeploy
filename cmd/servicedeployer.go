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

package cmd

import (
  "fmt"
  "github.com/openwhisk/go-whisk/whisk"
  "net/http"
  "net/url"
  "io/ioutil"
  "bufio"
  "os"
  "path"
  "path/filepath"
  "strings"
)

var deployer = new(ServiceDeployer)


//ServiceDeployer defines a prototype service deployer.  It should do the following:
//   1. Collect information from the manifest file (if any)
//   2. Collect information from the deployment file (if any)
//   3. Collect information about the source code files in the working directory
//   4. Create a deployment plan to create OpenWhisk service
type ServiceDeployer struct {
  actions []*whisk.Action
  authtoken string
  namespace string
  apihost string
}

// NewServiceDeployer is a Factory to create a new ServiceDeployer
func NewServiceDeployer() *ServiceDeployer {
  return new(ServiceDeployer)
}

// Load configuration will load properties from a file
func (deployer *ServiceDeployer) LoadConfiguration(propPath string) error {
  fmt.Println("Loading configuration")

  props, err := readProps(propPath)
  Check(err)

  fmt.Println("Got props ", props)

  deployer.namespace = props["NAMEPSACE"]
  deployer.apihost = props["APIHOST"]
  deployer.authtoken = props["AUTH"]

  return nil
}

// ReadDirectory will collect information from the files on disk. These represent actions
func (deployer *ServiceDeployer) ReadDirectory(directoryPath string) error {
  err := filepath.Walk(directoryPath, processFilePath)
  Check(err)

  return nil
}

// DeployActions into OpenWhisk
func (deployer *ServiceDeployer) DeployActions() error {

  for _, action := range deployer.actions {
    fmt.Println("Got action ", action.Exec.Code)
    deployer.createAction(action)
  }
  return nil
}

// Utility function to call go-whisk framework to make action
func (deployer *ServiceDeployer) createAction(action *whisk.Action) {

	baseURL, err := getURLBase(deployer.apihost)
	if err != nil {
		fmt.Println("Got error making baseUrl ", err)
	}

	clientConfig := &whisk.Config{
		AuthToken: deployer.authtoken,
		Namespace: deployer.namespace,
		BaseURL:   baseURL,
		Version:   "v1",
		Insecure:  false, // true if you want to ignore certificate signing
	}

	// Setup network client
	client, err := whisk.NewClient(http.DefaultClient, clientConfig)
	if err != nil {
		fmt.Println("Got error making whisk client ", err)
	}

	action.Namespace = "MyNameSpace"
	action.Publish = false
	// action.Parameters =
	// action.Annotations =
	// action.Limits =

	// call ActionService Thru Client
	_, _, err = client.Actions.Insert(action, false, false)
	if err != nil {
		fmt.Println("Got error inserting action ", err)
	}
}

//===== Function code

// Utility to convert hostname to URL object
func getURLBase(host string) (*url.URL, error) {

	urlBase := fmt.Sprintf("%s/api/", host)
	url, err := url.Parse(urlBase)

	if len(url.Scheme) == 0 || len(url.Host) == 0 {
		urlBase = fmt.Sprintf("https://%s/api/", host)
		url, err = url.Parse(urlBase)
	}

	return url, err
}


func processFilePath(filePath string, f os.FileInfo, err error) error {
	fmt.Println("Visited ", filePath)

  ext := path.Ext(filePath)
  baseName:= path.Base(filePath)
  name := strings.TrimSuffix(baseName, filepath.Ext(baseName))

  // process source code files
  if ext == ".swift"  {

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

    action := new(whisk.Action)
    action.Exec = new(whisk.Exec)
    action.Exec.Code = string(dat)
    action.Exec.Kind = kind
    action.Name = name
    action.Publish = false

    deployer.actions = append(deployer.actions, action)
  }
	return nil
}

func readProps(path string) (map[string]string, error) {

    props := map[string]string{}

    file, err := os.Open(path)
    if err != nil {
        // If file does not exist, just return props
      fmt.Printf("Unable to read whisk properties file '%s' (file open error: %s); falling back to default properties\n" ,path, err)
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
