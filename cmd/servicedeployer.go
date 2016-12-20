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
	Actions        map[string]utils.ActionRecord
	Triggers       map[string]*whisk.Trigger
	Packages       map[string]*whisk.SentPackageNoPublish
	Rules          map[string]*whisk.Rule
	Client         *whisk.Client
	mt             sync.RWMutex
	IsInteractive  bool
	IsDefault      bool
	ManifestPath   string
	ProjectPath    string
	DeploymentPath string
	//whether to deploy the action under the package
	DeployActionInPackage bool
	InteractiveChoice     bool
}

// NewServiceDeployer is a Factory to create a new ServiceDeployer
func NewServiceDeployer() *ServiceDeployer {
	var dep ServiceDeployer
	dep.Actions = make(map[string]utils.ActionRecord)
	dep.Packages = make(map[string]*whisk.SentPackageNoPublish)
	dep.Triggers = make(map[string]*whisk.Trigger)
	dep.Rules = make(map[string]*whisk.Rule)
	dep.IsInteractive = true
	dep.DeployActionInPackage = true
	return &dep
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

// DeployActions into OpenWhisk
func (deployer *ServiceDeployer) DeployActions() error {

	for _, action := range deployer.Actions {
		//fmt.Println("Got action ", action.Exec.Code)
		deployer.createAction(action.Action)
	}
	return nil
}

// Deploy Triggers into OpenWhisk
func (deployer *ServiceDeployer) DeployTriggers() error {
	for _, trigger := range deployer.Triggers {
		deployer.createTrigger(trigger)
	}
	return nil

}

// Deploy Rules into OpenWhisk
func (deployer *ServiceDeployer) DeployRules() error {
	for _, rule := range deployer.Rules {
		deployer.createRule(rule)
	}

	return nil
}

func (deployer *ServiceDeployer) createTrigger(trigger *whisk.Trigger) {
	_, _, err := deployer.Client.Triggers.Insert(trigger, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating trigger with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
}

func (deployer *ServiceDeployer) createRule(rule *whisk.Rule) {
	// The rule's trigger should include the namespace with pattern /namespace/trigger
	rule.Trigger = deployer.getQualifiedName(rule.Trigger, clientConfig.Namespace)
	// The rule's action should include the namespace and package with pattern /namespace/package/action
	// please refer https://github.com/openwhisk/openwhisk/issues/1577
	pkgName := deployer.getPackageName()
	rule.Action = deployer.getQualifiedName(strings.Join([]string{pkgName, rule.Action}, "/"), clientConfig.Namespace)
	_, _, err := deployer.Client.Rules.Insert(rule, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating rule with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
}

// Utility function to call go-whisk framework to make action
func (deployer *ServiceDeployer) createAction(action *whisk.Action) {
	// call ActionService Thru Client
	if deployer.DeployActionInPackage {
		pkgname := deployer.getPackageName()
		// the action will be created under package with pattern 'packagename/actionname'
		action.Name = strings.Join([]string{pkgname, action.Name}, "/")
	}
	_, _, err := deployer.Client.Actions.Insert(action, false, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating action with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
}

func (deployer *ServiceDeployer) getPackageName() string {
	if len(deployer.Packages) > 0 {
		//get the first available package name by default
		for _, pkg := range deployer.Packages {
			if pkg.Name != "" {
				return pkg.Name
			}
		}
	}
	return ""
}

func (deployer *ServiceDeployer) createPackage(packa *whisk.SentPackageNoPublish) {
	_, _, err := deployer.Client.Packages.Insert(packa, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating package with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
}

// Wrapper parser to handle yaml dir
func (deployer *ServiceDeployer) HandleYamlDir() error {
	mm := utils.NewYAMLParser()
	packg, err := mm.ComposePackage(deployer.ManifestPath)
	utils.Check(err)
	actions, err := mm.ComposeActions(deployer.ManifestPath)
	utils.Check(err)
	//for implementation of feed support(issue#47). The trigger configs are in manifest.yaml,
	//If there is feed nodes, then parse the deployment yaml to get the feed, and create
	//the triggers with feed configs as annonations of the trigger, so change the compose trigger
	//interface to include both manifest and deployment path.
	triggers, err := mm.ComposeTriggers(deployer.ManifestPath, deployer.DeploymentPath)
	utils.Check(err)
	rules, err := mm.ComposeRules(deployer.ManifestPath)
	utils.Check(err)
	if !deployer.SetActions(actions) {
		log.Panicln("duplication founded during deploy actions")
	}
	if !deployer.SetPackage(packg) {
		log.Panicln("duplication founded during deploy package")
	}
	if !deployer.SetTriggers(triggers) {
		log.Panicln("duplication founded during deploy triggers")
	}
	if !deployer.SetRules(rules) {
		log.Panicln("duplication founded during deploy rules")
	}

	deployer.createPackage(packg)

	return nil
}

func (deployer *ServiceDeployer) printDeploymentAssets() {

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
	if deployer.IsInteractive == true {
		deployer.printDeploymentAssets()
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Do you really want to deploy this? (y/n): ")

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if strings.EqualFold(text, "y") || strings.EqualFold(text, "yes") {
			deployer.InteractiveChoice = true
			if err := deployer.DeployActions(); err != nil {
				return err
			}
			if err := deployer.DeployTriggers(); err != nil {
				return err
			}
			if err := deployer.DeployRules(); err != nil {
				return err
			}
		} else {
			deployer.InteractiveChoice = false
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
		//fmt.Println(action.Action.Name)
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

func (deployer *ServiceDeployer) SetTriggers(triggers []*whisk.Trigger) bool {
	deployer.mt.Lock()
	defer deployer.mt.Unlock()

	for _, trigger := range triggers {
		existTrigger, exist := deployer.Triggers[trigger.Name]
		if exist {
			existTrigger.Name = trigger.Name
			existTrigger.ActivationId = trigger.ActivationId
			existTrigger.Namespace = trigger.Namespace
			existTrigger.Annotations = trigger.Annotations
			existTrigger.Version = trigger.Version
			existTrigger.Parameters = trigger.Parameters
			existTrigger.Publish = trigger.Publish
		} else {
			deployer.Triggers[trigger.Name] = trigger
		}

	}
	return true
}

func (deployer *ServiceDeployer) SetRules(rules []*whisk.Rule) bool {
	deployer.mt.Lock()
	defer deployer.mt.Unlock()

	for _, rule := range rules {
		existRule, exist := deployer.Rules[rule.Name]
		if exist {
			existRule.Name = rule.Name
			existRule.Publish = rule.Publish
			existRule.Version = rule.Version
			existRule.Namespace = rule.Namespace
			existRule.Action = rule.Action
			existRule.Trigger = rule.Trigger
			existRule.Status = rule.Status
		} else {
			deployer.Rules[rule.Name] = rule
		}

	}
	return true
}

// from whisk go client
func (deployer *ServiceDeployer) getQualifiedName(name string, namespace string) string {
	if strings.HasPrefix(name, "/") {
		return name
	} else if strings.HasPrefix(namespace, "/") {
		return fmt.Sprintf("%s/%s", namespace, name)
	} else {
		if len(namespace) == 0 {
			namespace = clientConfig.Namespace
		}
		return fmt.Sprintf("/%s/%s", namespace, name)
	}
}
