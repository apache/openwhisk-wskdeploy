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
	"github.com/openwhisk/wsktool/utils"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var deployer = NewServiceDeployer()

//ServiceDeployer defines a prototype service deployer.  It should do the following:
//   1. Collect information from the manifest file (if any)
//   2. Collect information from the deployment file (if any)
//   3. Collect information about the source code files in the working directory
//   4. Create a deployment plan to create OpenWhisk service
type ServiceDeployer struct {
	actions   map[string]*whisk.Action
	triggers  map[string]string
	packages  map[string]string
	authtoken string
	namespace string
	apihost   string
}

// NewServiceDeployer is a Factory to create a new ServiceDeployer
func NewServiceDeployer() *ServiceDeployer {
	var dep ServiceDeployer
	dep.actions = make(map[string]*whisk.Action)
	return &dep
}

// Load configuration will load properties from a file
func (deployer *ServiceDeployer) LoadConfiguration(propPath string) error {
	fmt.Println("Loading configuration")

	props, err := utils.ReadProps(propPath)
	utils.Check(err)

	fmt.Println("Got props ", props)

	deployer.namespace = props["NAMEPSACE"]
	deployer.apihost = props["APIHOST"]
	deployer.authtoken = props["AUTH"]

	return nil
}

// ReadDirectory will collect information from the files on disk. These represent actions
func (deployer *ServiceDeployer) ReadDirectory(directoryPath string) error {

	err := filepath.Walk(directoryPath, func(filePath string, f os.FileInfo, err error) error {
		if filePath != directoryPath {
			isDirectory := utils.IsDirectory(filePath)

			if isDirectory == true {
				baseName := path.Base(filePath)
				if strings.HasPrefix(baseName, ".") {
					return filepath.SkipDir
				}
				err = deployer.CreatePackageFromDirectory(baseName)

			} else {
				err = deployer.CreateActionFromFile(filePath)
			}
		}
		return err
	})

	utils.Check(err)
	return nil
}

func (deployer *ServiceDeployer) CreatePackageFromDirectory(directoryName string) error {
	fmt.Println("Making a package ", directoryName)
	return nil
}

func (deployer *ServiceDeployer) CreateActionFromFile(filePath string) error {
	ext := path.Ext(filePath)
	baseName := path.Base(filePath)
	name := strings.TrimSuffix(baseName, filepath.Ext(baseName))

	// process source code files
	if ext == ".swift" || ext == ".js" || ext == ".py" {

		if _, ok := deployer.actions[name]; ok {
			return fmt.Errorf("Found a duplicate name %s when scanning file directory", name)

		} else {

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
			utils.Check(err)

			action := new(whisk.Action)
			action.Exec = new(whisk.Exec)
			action.Exec.Code = string(dat)
			action.Exec.Kind = kind
			action.Name = name
			action.Publish = false

			deployer.actions[name] = action
		}
	}
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

	baseURL, err := utils.GetURLBase(deployer.apihost)
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

	action.Namespace = deployer.namespace
	action.Publish = false
	// action.Parameters =
	// action.Annotations =
	// action.Limits =

	// call ActionService Thru Client
	_, _, err = client.Actions.Insert(action, false, true)
	if err != nil {
		fmt.Println("Got error inserting action ", err)
	}
}
