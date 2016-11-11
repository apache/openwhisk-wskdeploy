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
	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/wskdeploy/utils"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

//ServiceDeployer defines a prototype service deployer.  It should do the following:
//   1. Collect information from the manifest file (if any)
//   2. Collect information from the deployment file (if any)
//   3. Collect information about the source code files in the working directory
//   4. Create a deployment plan to create OpenWhisk service
type ServiceDeployer struct {
	Actions   map[string]*whisk.Action
	Triggers  map[string]*whisk.Trigger
	Packages  map[string]*whisk.SentPackageNoPublish
	Rules     map[string]*whisk.Rule
	Client    *whisk.Client
	mt        sync.RWMutex
	Authtoken string
	Namespace string
	Apihost   string
}

// NewServiceDeployer is a Factory to create a new ServiceDeployer
func NewServiceDeployer() *ServiceDeployer {
	var dep ServiceDeployer
	dep.Actions = make(map[string]*whisk.Action)
	dep.Packages = make(map[string]*whisk.SentPackageNoPublish)
	dep.Triggers = make(map[string]*whisk.Trigger)
	dep.Rules = make(map[string]*whisk.Rule)
	return &dep
}

// Load configuration will load properties from a file
func (deployer *ServiceDeployer) LoadConfiguration(propPath string) error {
	fmt.Println("Loading configuration")
	props, err := utils.ReadProps(propPath)
	utils.Check(err)
	fmt.Println("Got props ", props)
	deployer.Namespace = props["NAMESPACE"]
	deployer.Apihost = props["APIHOST"]
	deployer.Authtoken = props["AUTH"]
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
				//err = deployer.CreateActionFromFile(filePath)
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

func (deployer *ServiceDeployer) CreateClient() {
	baseURL, err := utils.GetURLBase(deployer.Apihost)
	utils.Check(err)
	clientConfig := &whisk.Config{
		AuthToken: deployer.Authtoken,
		Namespace: deployer.Namespace,
		BaseURL:   baseURL,
		Version:   "v1",
		Insecure:  true, // true if you want to ignore certificate signing
	}
	// Setup network client
	client, err := whisk.NewClient(http.DefaultClient, clientConfig)
	utils.Check(err)
	deployer.Client = client

}

// DeployActions into OpenWhisk
func (deployer *ServiceDeployer) DeployActions() error {

	for _, action := range deployer.Actions {
		//fmt.Println("Got action ", action.Exec.Code)
		deployer.createAction(action)
	}
	return nil
}

// Utility function to call go-whisk framework to make action
func (deployer *ServiceDeployer) createAction(action *whisk.Action) {
	// call ActionService Thru Client
	_, _, err := deployer.Client.Actions.Insert(action, false, true)
	if err != nil {
		fmt.Println("Got error inserting action ", err)
	}
}

func (deployer *ServiceDeployer) createPackage(packa *whisk.SentPackageNoPublish) {
	_, _, err := deployer.Client.Packages.Insert(packa, true)
	if err != nil {
		fmt.Errorf("Got error creating package %s.", err)
	}
}

// Wrapper parser to handle yaml dir
func (deployer *ServiceDeployer) HandleYamlDir(manifestpath string) error {
	mm := utils.NewManifestManager()
	packg, err := mm.ComposePackage(manifestpath)
	utils.Check(err)
	actions, err := mm.ComposeActions(manifestpath)
	utils.Check(err)
	if !deployer.SetActions(actions) {
		log.Panicln("duplication founded during deploy actions")
	}
	if !deployer.SetPackage(packg) {
		log.Panicln("duplication founded during deploy package")
	}

	deployer.createPackage(packg)
	deployer.DeployActions()
	return nil
}

// Use relfect util to deploy everything in this service deployer
// according some planning?
func (deployer *ServiceDeployer) Deploy() {
}

func (deployer *ServiceDeployer) SetPackage(pkg *whisk.SentPackageNoPublish) bool {
	deployer.mt.Lock()
	defer deployer.mt.Unlock()
	_, exist := deployer.Packages[pkg.Name]
	if exist {
		return false
	}
	deployer.Packages[pkg.Name] = pkg
	return true
}

func (deployer *ServiceDeployer) SetActions(actions []*whisk.Action) bool {
	deployer.mt.Lock()
	defer deployer.mt.Unlock()

	for _, action := range actions {
		fmt.Println(action.Name)
		_, exist := deployer.Actions[action.Name]
		if exist {
			return false
		}
		deployer.Actions[action.Name] = action
	}
	return true
}
