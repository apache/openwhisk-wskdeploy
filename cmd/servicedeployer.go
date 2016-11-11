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
	"bufio"
	"bytes"
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
	Actions   map[string]utils.ActionRecord
	Triggers  map[string]*whisk.Trigger
	Packages  map[string]*whisk.SentPackageNoPublish
	Rules     map[string]*whisk.Rule
	Client    *whisk.Client
	mt        sync.RWMutex
	Authtoken string
	Namespace string
	Apihost   string
	IsInteractive bool
	IsDefault bool
	ManifestPath string
	ProjectPath string
	DeploymentPath string

}

// NewServiceDeployer is a Factory to create a new ServiceDeployer
func NewServiceDeployer() *ServiceDeployer {
	var dep ServiceDeployer
	dep.Actions = make(map[string]utils.ActionRecord)
	dep.Packages = make(map[string]*whisk.SentPackageNoPublish)
	dep.Triggers = make(map[string]*whisk.Trigger)
	dep.Rules = make(map[string]*whisk.Rule)
	dep.IsInteractive = true
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

// ConstructDeploymentPlan will collect information from the manifest, descriptors, and any
// defaults to determine what assets need to be installed.
func (deployer *ServiceDeployer) ConstructDeploymentPlan() error {

	if deployer.IsDefault == true {
		deployer.ReadDirectory()
	}

	deployer.HandleYamlDir()

	return nil
}

// ReadDirectory will collect information from the files on disk. These represent actions
func (deployer *ServiceDeployer) ReadDirectory() error {

	err := filepath.Walk(deployer.ProjectPath, func(filePath string, f os.FileInfo, err error) error {
		if filePath != deployer.ProjectPath {
			isDirectory := utils.IsDirectory(filePath)

			if isDirectory == true {
				baseName := path.Base(filePath)
				if strings.HasPrefix(baseName, ".") {
					return filepath.SkipDir
				}
				err = deployer.CreatePackageFromDirectory(baseName)

				} else {
					action, err := utils.CreateActionFromFile(deployer.ProjectPath, filePath)
					utils.Check(err)
					deployer.Actions[action.Name] = utils.ActionRecord{action, filePath}
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
			deployer.createAction(action.Action)
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
	func (deployer *ServiceDeployer) HandleYamlDir() error {
		mm := utils.NewManifestManager()
		packg, err := mm.ComposePackage(deployer.ManifestPath)
		utils.Check(err)
		actions, err := mm.ComposeActions(deployer.ManifestPath)
		utils.Check(err)
		if !deployer.SetActions(actions) {
			log.Panicln("duplication founded during deploy actions")
		}
		if !deployer.SetPackage(packg) {
			log.Panicln("duplication founded during deploy package")
		}

		deployer.createPackage(packg)

		return nil
	}

	func printDeploymentAssets(deployer *ServiceDeployer) {

		fmt.Println("----==== OpenWhisk Deployment Plan ====----")
		fmt.Println("Deploy Packages:")
		fmt.Println("----------------")
		for _, pkg := range deployer.Packages {
			var buffer bytes.Buffer
			buffer.WriteString(pkg.Namespace)
			buffer.WriteString("/")
			buffer.WriteString(pkg.Name)
			fmt.Printf("    %s (version: %s)\n", buffer.String(), pkg.Version)
		}

    fmt.Println("\nDeploy Actions:")
		fmt.Println("----------------")
		for _, action := range deployer.Actions {
			var buffer bytes.Buffer
			buffer.WriteString(action.Action.Namespace)
			buffer.WriteString("/")
			buffer.WriteString(action.Action.Name)
			fmt.Printf("    %s (version: %s)\n", buffer.String(), action.Action.Version)
		}


		fmt.Println("\nDeploy Triggers:")
		fmt.Println("----------------")
		for _, trigger := range deployer.Triggers {
			var buffer bytes.Buffer
			buffer.WriteString(trigger.Namespace)
			buffer.WriteString("/")
			buffer.WriteString(trigger.Name)
			fmt.Printf("    %s (version: %s)\n", buffer.String(), trigger.Version)
		}

		fmt.Println("\nDeploy Rules:")
		fmt.Println("----------------")
		for _, rule := range deployer.Rules {
			var buffer bytes.Buffer
			buffer.WriteString(rule.Namespace)
			buffer.WriteString("/")
			buffer.WriteString(rule.Name)
			fmt.Printf("    %s (version: %s)\n", buffer.String(), rule.Version)
		}

	}

	// Use relfect util to deploy everything in this service deployer
	// according some planning?
	func (deployer *ServiceDeployer) Deploy() error {
		if (deployer.IsInteractive == true) {
			printDeploymentAssets(deployer)
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Do you really want to deploy this? (y/n): ")

			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)

			if strings.EqualFold(text, "y") || strings.EqualFold(text, "yes") {
					err := deployer.DeployActions()
					if err != nil {
						return err
					}
			} else {
				fmt.Println("OK. Cancelling deployment")
			}
		}

		return nil
	}

	func (deployer *ServiceDeployer) SetPackage(pkg *whisk.SentPackageNoPublish) bool {
		deployer.mt.Lock()
		defer deployer.mt.Unlock()
		existPkg, exist := deployer.Packages[pkg.Name]
		if exist {
			if deployer.IsDefault == true {

				log.Printf("Updating package %s with values from manifest file ", pkg.Name)

				existPkg.Annotations = pkg.Annotations
				existPkg.Namespace = pkg.Namespace
				existPkg.Parameters = pkg.Parameters
				existPkg.Publish = pkg.Publish
				existPkg.Version = pkg.Version

				deployer.Packages[pkg.Name] = existPkg
				return true
			} else {
				return false
			}
		}

		deployer.Packages[pkg.Name] = pkg
		return true
	}

	func (deployer *ServiceDeployer) SetActions(actions []utils.ActionRecord) bool {
		deployer.mt.Lock()
		defer deployer.mt.Unlock()

		for _, action := range actions {
			fmt.Println(action.Action.Name)
			existAction, exist := deployer.Actions[action.Action.Name]

			if exist {
				if deployer.IsDefault == true {
					// look for actions declared in filesystem default as well as manifest
					// if one exists, merge if they are the same (either same Filepath or manifest doesn't specify a Filepath)
					// if they are not the same log error
					if action.Filepath != "" {
						if existAction.Filepath != action.Filepath {
							log.Printf("Action %s has location %s in manifest but already exists at %s", action.Action.Name, action.Filepath, existAction)
							return false
						} else {
							// merge the two, overwrite existing action with manifest values
							existAction.Action.Annotations = action.Action.Annotations
							existAction.Action.Exec.Kind = action.Action.Exec.Kind
							existAction.Action.Limits = action.Action.Limits
							existAction.Action.Namespace = action.Action.Namespace
							existAction.Action.Parameters = action.Action.Parameters
							existAction.Action.Publish = action.Action.Publish
							existAction.Action.Version = action.Action.Version

							deployer.Actions[action.Action.Name] = existAction

							return true
						}
					}
				} else {
					// no defaults, so assume everything is in the incoming ActionRecord
					// return false since it means the action is declared twice in the manifest
					log.Printf("Action %s is declared more than once", action.Action.Name)
					return false
				}
			}
			// doesn't exist so just add to deployer actions
				deployer.Actions[action.Action.Name] = action
			}
			return true
		}
