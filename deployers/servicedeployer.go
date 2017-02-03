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

package deployers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
)

type DeploymentApplication struct {
	Packages map[string]*DeploymentPackage
}

func NewDeploymentApplication() *DeploymentApplication {
	var dep DeploymentApplication
	dep.Packages = make(map[string]*DeploymentPackage)
	return &dep
}

type DeploymentPackage struct {
	Package   *whisk.SentPackageNoPublish
	Actions   map[string]utils.ActionRecord
	Triggers  map[string]*whisk.Trigger
	Rules     map[string]*whisk.Rule
	Sequences map[string]utils.ActionRecord
}

func NewDeploymentPackage() *DeploymentPackage {
	var dep DeploymentPackage
	dep.Actions = make(map[string]utils.ActionRecord)
	dep.Triggers = make(map[string]*whisk.Trigger)
	dep.Rules = make(map[string]*whisk.Rule)
	dep.Sequences = make(map[string]utils.ActionRecord)
	return &dep
}

//ServiceDeployer defines a prototype service deployer.  It should do the following:
//   1. Collect information from the manifest file (if any)
//   2. Collect information from the deployment file (if any)
//   3. Collect information about the source code files in the working directory
//   4. Create a deployment plan to create OpenWhisk service
type ServiceDeployer struct {
	Deployment     *DeploymentApplication
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
	ClientConfig          *whisk.Config
}

// NewServiceDeployer is a Factory to create a new ServiceDeployer
func NewServiceDeployer() *ServiceDeployer {
	var dep ServiceDeployer
	dep.Deployment = NewDeploymentApplication()
	dep.IsInteractive = true
	dep.DeployActionInPackage = true

	return &dep
}

func (deployer *ServiceDeployer) ConstructDeploymentPlan() error {

	// process manifest file
	var manifestReader = NewManfiestReader(deployer)
	manifestReader.HandleYaml()

	// process deploymet file
	var deploymentReader = NewDeploymentReader(deployer)
	deploymentReader.HandleYaml()

	deploymentReader.BindAssets()

	// print report
	deployer.printDeploymentAssets()

	return nil
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
			return deployer.deployAssets()

		} else {
			deployer.InteractiveChoice = false
			fmt.Println("OK. Cancelling deployment")
			return nil
		}
	}

	// non-interactive
	return deployer.deployAssets()

}

func (deployer *ServiceDeployer) deployAssets() error {

	if err := deployer.DeployPackages(); err != nil {
		return err
	}

	if err := deployer.DeployActions(); err != nil {
		return err
	}

	if err := deployer.DeploySequences(); err != nil {
		return err
	}

	if err := deployer.DeployTriggers(); err != nil {
		return err
	}

	if err := deployer.DeployRules(); err != nil {
		return err
	}
	return nil
}

func (deployer *ServiceDeployer) DeployPackages() error {

	for _, pack := range deployer.Deployment.Packages {
		deployer.createPackage(pack.Package)
	}
	return nil
}

// DeployActions into OpenWhisk
func (deployer *ServiceDeployer) DeploySequences() error {

	for _, pack := range deployer.Deployment.Packages {
		for _, action := range pack.Sequences {
			deployer.createAction(pack.Package.Name, action.Action)
		}
	}
	return nil
}

// DeployActions into OpenWhisk
func (deployer *ServiceDeployer) DeployActions() error {

	for _, pack := range deployer.Deployment.Packages {
		for _, action := range pack.Actions {
			deployer.createAction(pack.Package.Name, action.Action)
		}
	}
	return nil
}

// Deploy Triggers into OpenWhisk
func (deployer *ServiceDeployer) DeployTriggers() error {
	for _, pack := range deployer.Deployment.Packages {
		for _, trigger := range pack.Triggers {
			deployer.createTrigger(pack.Package.Name, trigger)
		}
	}
	return nil

}

// Deploy Rules into OpenWhisk
func (deployer *ServiceDeployer) DeployRules() error {
	for _, pack := range deployer.Deployment.Packages {
		for _, rule := range pack.Rules {
			deployer.createRule(pack.Package.Name, rule)
		}
	}
	return nil
}

func (deployer *ServiceDeployer) createPackage(packa *whisk.SentPackageNoPublish) {
	_, _, err := deployer.Client.Packages.Insert(packa, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating package with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
}

func (deployer *ServiceDeployer) createTrigger(pkgname string, trigger *whisk.Trigger) {
	_, _, err := deployer.Client.Triggers.Insert(trigger, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating trigger with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
}

func (deployer *ServiceDeployer) createRule(pkgname string, rule *whisk.Rule) {
	// The rule's trigger should include the namespace with pattern /namespace/trigger
	rule.Trigger = deployer.getQualifiedName(rule.Trigger, deployer.ClientConfig.Namespace)
	// The rule's action should include the namespace and package with pattern /namespace/package/action
	// please refer https://github.com/openwhisk/openwhisk/issues/1577

	rule.Action = deployer.getQualifiedName(strings.Join([]string{pkgname, rule.Action}, "/"), deployer.ClientConfig.Namespace)
	_, _, err := deployer.Client.Rules.Insert(rule, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating rule with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
}

// Utility function to call go-whisk framework to make action
func (deployer *ServiceDeployer) createAction(pkgname string, action *whisk.Action) {
	// call ActionService Thru Client
	if deployer.DeployActionInPackage {

		// the action will be created under package with pattern 'packagename/actionname'
		action.Name = strings.Join([]string{pkgname, action.Name}, "/")
	}
	_, _, err := deployer.Client.Actions.Insert(action, false, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating action with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
}

// from whisk go client
func (deployer *ServiceDeployer) getQualifiedName(name string, namespace string) string {
	if strings.HasPrefix(name, "/") {
		return name
	} else if strings.HasPrefix(namespace, "/") {
		return fmt.Sprintf("%s/%s", namespace, name)
	} else {
		if len(namespace) == 0 {
			namespace = deployer.ClientConfig.Namespace
		}
		return fmt.Sprintf("/%s/%s", namespace, name)
	}
}

func (deployer *ServiceDeployer) printDeploymentAssets() {

	fmt.Println("         ____      ___                   _    _ _     _     _\n        /\\   \\    / _ \\ _ __   ___ _ __ | |  | | |__ (_)___| | __\n   /\\  /__\\   \\  | | | | '_ \\ / _ \\ '_ \\| |  | | '_ \\| / __| |/ /\n  /  \\____ \\  /  | |_| | |_) |  __/ | | | |/\\| | | | | \\__ \\   <\n  \\   \\  /  \\/    \\___/| .__/ \\___|_| |_|__/\\__|_| |_|_|___/_|\\_\\ \n   \\___\\/              |_|\n")

	fmt.Println("Packages:")
	for _, pack := range deployer.Deployment.Packages {
		fmt.Println("Name: " + pack.Package.Name)

		for _, action := range pack.Actions {
			fmt.Println("  * action: " + action.Action.Name)
			fmt.Println("    bindings: ")
			for _, p := range action.Action.Parameters {

				value := "?"
				if str, ok := p.Value.(string); ok {
					value = str
				}
				fmt.Println("        - name: " + p.Key + " value: " + value)
			}
		}

		fmt.Println("")
		for _, action := range pack.Sequences {
			fmt.Println("  * sequence: " + action.Action.Name)
		}

		fmt.Println("")
		for _, trigger := range pack.Triggers {
			fmt.Println("  * trigger: " + trigger.Name)
			fmt.Println("    bindings: ")

			for _, p := range trigger.Parameters {

				value := "?"
				if str, ok := p.Value.(string); ok {
					value = str
				}
				fmt.Println("        - name: " + p.Key + " value: " + value)
			}
		}

		fmt.Println("")
		for _, rule := range pack.Rules {
			fmt.Println("  rule: " + rule.Name)
			fmt.Println("    - trigger: " + rule.Trigger + "\n    - action: " + rule.Action)
		}

		fmt.Println("")

	}

}
