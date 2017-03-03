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
	"github.com/openwhisk/openwhisk-wskdeploy/parsers"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
)

type DeploymentApplication struct {
	Packages map[string]*DeploymentPackage
	Triggers map[string]*whisk.Trigger
	Rules    map[string]*whisk.Rule
}

func NewDeploymentApplication() *DeploymentApplication {
	var dep DeploymentApplication
	dep.Packages = make(map[string]*DeploymentPackage)
	dep.Triggers = make(map[string]*whisk.Trigger)
	dep.Rules = make(map[string]*whisk.Rule)
	return &dep
}

type DeploymentPackage struct {
	Package   *whisk.SentPackageNoPublish
	Actions   map[string]utils.ActionRecord
	Sequences map[string]utils.ActionRecord
}

func NewDeploymentPackage() *DeploymentPackage {
	var dep DeploymentPackage
	dep.Actions = make(map[string]utils.ActionRecord)
	dep.Sequences = make(map[string]utils.ActionRecord)
	return &dep
}

//ServiceDeployer defines a prototype service deployer.  It should do the following:
//   1. Collect information from the manifest file (if any)
//   2. Collect information from the deployment file (if any)
//   3. Collect information about the source code files in the working directory
//   4. Create a deployment plan to create OpenWhisk service
type ServiceDeployer struct {
	Deployment      *DeploymentApplication
	Client          *whisk.Client
	mt              sync.RWMutex
	RootPackageName string
	IsInteractive   bool
	IsDefault       bool
	ManifestPath    string
	ProjectPath     string
	DeploymentPath  string
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

// Check if the manifest yaml could be parsed by Manifest Parser.
// Check if the deployment yaml could be parsed by Manifest Parser.
func (deployer *ServiceDeployer) Check() {
	ps := parsers.NewYAMLParser()
	ps.ParseDeployment(deployer.DeploymentPath)
	ps.ParseManifest(deployer.ManifestPath)
	// add more schema check or manifest/deployment consistency checks here if
	// necessary
}

func (deployer *ServiceDeployer) ConstructDeploymentPlan() error {

	var manifestReader = NewManfiestReader(deployer)
	manifest, manifestParser, err := manifestReader.ParseManifest()
	deployer.RootPackageName = manifest.Package.Packagename

	utils.Check(err)

	manifestReader.InitRootPackage(manifestParser, manifest)

	if deployer.IsDefault == true {
		fileReader := NewFileSystemReader(deployer)
		fileActions, err := fileReader.ReadProjectDirectory(manifest)
		utils.Check(err)

		fileReader.SetFileActions(fileActions)
	}

	// process manifest file
	err = manifestReader.HandleYaml(manifestParser, manifest)
	utils.Check(err)

	// process deploymet file
	if utils.FileExists(deployer.DeploymentPath) {
		var deploymentReader = NewDeploymentReader(deployer)
		deploymentReader.HandleYaml()

		deploymentReader.BindAssets()
	}

	return nil
}

func (deployer *ServiceDeployer) ConstructUnDeploymentPlan() (*DeploymentApplication, error) {

	var manifestReader = NewManfiestReader(deployer)
	manifest, manifestParser, err := manifestReader.ParseManifest()
	utils.Check(err)

	manifestReader.InitRootPackage(manifestParser, manifest)

	// process file system
	if deployer.IsDefault == true {
		fileReader := NewFileSystemReader(deployer)
		fileActions, err := fileReader.ReadProjectDirectory(manifest)
		utils.Check(err)

		err = fileReader.SetFileActions(fileActions)
		utils.Check(err)

	}

	// process manifest file
	err = manifestReader.HandleYaml(manifestParser, manifest)
	utils.Check(err)

	// process deploymet file
	if utils.FileExists(deployer.DeploymentPath) {
		var deploymentReader = NewDeploymentReader(deployer)
		deploymentReader.HandleYaml()

		deploymentReader.BindAssets()
	}

	verifiedPlan := deployer.Deployment

	return verifiedPlan, nil
}

// Use relfect util to deploy everything in this service deployer
// according some planning?
func (deployer *ServiceDeployer) Deploy() error {

	if deployer.IsInteractive == true {
		deployer.printDeploymentAssets(deployer.Deployment)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Do you really want to deploy this? (y/N): ")

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "" {
			text = "n"
		}

		if strings.EqualFold(text, "y") || strings.EqualFold(text, "yes") {
			deployer.InteractiveChoice = true
			if err := deployer.deployAssets(); err != nil {
				fmt.Println("\nDeployment did not complete sucessfully. Run `wskdeploy undeploy` to remove partially deployed assets")
				return err
			}

			fmt.Println("\nDeployment completed successfully.")
			return nil

		} else {
			deployer.InteractiveChoice = false
			fmt.Println("OK. Cancelling deployment")
			return nil
		}
	}

	// non-interactive
	if err := deployer.deployAssets(); err != nil {
		fmt.Println("\nDeployment did not complete sucessfully. Run `wskdeploy undeploy` to remove partially deployed assets")
		return err
	}

	fmt.Println("\nDeployment completed successfully.")
	return nil

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
	for _, trigger := range deployer.Deployment.Triggers {

		if feedname, isFeed := utils.IsFeedAction(trigger); isFeed {
			deployer.createFeedAction(trigger, feedname)
		} else {
			deployer.createTrigger(trigger)
		}

	}
	return nil

}

// Deploy Rules into OpenWhisk
func (deployer *ServiceDeployer) DeployRules() error {
	for _, rule := range deployer.Deployment.Rules {
		deployer.createRule(rule)
	}
	return nil
}

func (deployer *ServiceDeployer) createPackage(packa *whisk.SentPackageNoPublish) {
	fmt.Print("Deploying pacakge " + packa.Name + " ... ")
	_, _, err := deployer.Client.Packages.Insert(packa, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating package with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
	fmt.Println("Done!")
}

func (deployer *ServiceDeployer) createTrigger(trigger *whisk.Trigger) {
	_, _, err := deployer.Client.Triggers.Insert(trigger, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating trigger with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
	fmt.Println("Done!")
}

func (deployer *ServiceDeployer) createFeedAction(trigger *whisk.Trigger, feedName string) {
	// to hold and modify trigger parameters, not passed by ref?
	params := make(map[string]interface{})

	// check for strings that are JSON
	for _, keyVal := range trigger.Parameters {
		fmt.Println("Checking trigger param " + keyVal.Key + " with value " + keyVal.Value.(string))
		if b, isJson := utils.IsJSON(keyVal.Value.(string)); isJson {
			fmt.Println("Setting JSON for trigger param " + keyVal.Key)
			params[keyVal.Key] = b
		} else {
			params[keyVal.Key] = keyVal.Value
		}
	}

	params["authKey"] = deployer.ClientConfig.AuthToken
	params["lifecycleEvent"] = "CREATE"
	params["triggerName"] = "/" + deployer.Client.Namespace + "/" + trigger.Name

	t := &whisk.Trigger{
		Name:        trigger.Name,
		Annotations: trigger.Annotations,
		Publish:     true,
	}

	_, _, err := deployer.Client.Triggers.Insert(t, false)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating trigger with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	} else {

		qName, err := utils.ParseQualifiedName(feedName, deployer.ClientConfig.Namespace)

		utils.Check(err)

		namespace := deployer.Client.Namespace
		deployer.Client.Namespace = qName.Namespace
		_, _, err = deployer.Client.Actions.Invoke(qName.EntityName, params, true)
		deployer.Client.Namespace = namespace

		if err != nil {
			wskErr := err.(*whisk.WskError)
			fmt.Printf("Got error creating trigger feed with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
		}
	}
	fmt.Println("Done!")
}

func (deployer *ServiceDeployer) createRule(rule *whisk.Rule) {
	// The rule's trigger should include the namespace with pattern /namespace/trigger
	rule.Trigger = deployer.getQualifiedName(rule.Trigger, deployer.ClientConfig.Namespace)
	// The rule's action should include the namespace and package with pattern /namespace/package/action
	// please refer https://github.com/openwhisk/openwhisk/issues/1577

	// if it contains a slash, then the action is qualified by a package name
	if strings.Contains(rule.Action, "/") {
		rule.Action = deployer.getQualifiedName(rule.Action, deployer.ClientConfig.Namespace)
	} else {
		// if not, we assume the action is inside the root package
		rule.Action = deployer.getQualifiedName(strings.Join([]string{deployer.RootPackageName, rule.Action}, "/"), deployer.ClientConfig.Namespace)
	}
	fmt.Print("Deploying rule " + rule.Name + " ... ")
	_, _, err := deployer.Client.Rules.Insert(rule, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating rule with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}

	_, _, err = deployer.Client.Rules.SetState(rule.Name, "active")
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error activating rule with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
	fmt.Println("Done!")
}

// Utility function to call go-whisk framework to make action
func (deployer *ServiceDeployer) createAction(pkgname string, action *whisk.Action) {
	// call ActionService Thru Client
	if deployer.DeployActionInPackage {
		// the action will be created under package with pattern 'packagename/actionname'
		action.Name = strings.Join([]string{pkgname, action.Name}, "/")
	}
	fmt.Print("Deploying action " + action.Name + " ... ")
	_, _, err := deployer.Client.Actions.Insert(action, false, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating action with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
	fmt.Println("Done!")
}

func (deployer *ServiceDeployer) UnDeploy(verifiedPlan *DeploymentApplication) error {
	if deployer.IsInteractive == true {
		deployer.printDeploymentAssets(verifiedPlan)
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Do you really want to undeploy this? (y/N): ")

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "" {
			text = "n"
		}

		if strings.EqualFold(text, "y") || strings.EqualFold(text, "yes") {
			deployer.InteractiveChoice = true

			if err := deployer.unDeployAssets(verifiedPlan); err != nil {
				fmt.Println("\nUndeployment did not complete sucessfully.")
				return err
			}

			fmt.Println("\nDeployment removed successfully.")
			return nil

		} else {
			deployer.InteractiveChoice = false
			fmt.Println("OK. Cancelling undeployment")
			return nil
		}
	}

	// non-interactive
	if err := deployer.unDeployAssets(verifiedPlan); err != nil {
		fmt.Println("\nUndeployment did not complete sucessfully.")
		return err
	}

	fmt.Println("\nDeployment removed successfully.")
	return nil

}

func (deployer *ServiceDeployer) unDeployAssets(verifiedPlan *DeploymentApplication) error {

	if err := deployer.UnDeployActions(verifiedPlan); err != nil {
		return err
	}

	if err := deployer.UnDeploySequences(verifiedPlan); err != nil {
		return err
	}

	if err := deployer.UnDeployTriggers(verifiedPlan); err != nil {
		return err
	}

	if err := deployer.UnDeployRules(verifiedPlan); err != nil {
		return err
	}

	if err := deployer.UnDeployPackages(verifiedPlan); err != nil {
		return err
	}

	return nil

}

func (deployer *ServiceDeployer) UnDeployPackages(deployment *DeploymentApplication) error {
	for _, pack := range deployment.Packages {
		deployer.deletePackage(pack.Package)
	}
	return nil
}

func (deployer *ServiceDeployer) UnDeploySequences(deployment *DeploymentApplication) error {

	for _, pack := range deployment.Packages {
		for _, action := range pack.Sequences {
			deployer.deleteAction(pack.Package.Name, action.Action)
		}
	}
	return nil
}

// DeployActions into OpenWhisk
func (deployer *ServiceDeployer) UnDeployActions(deployment *DeploymentApplication) error {

	for _, pack := range deployment.Packages {
		for _, action := range pack.Actions {
			deployer.deleteAction(pack.Package.Name, action.Action)
		}
	}
	return nil
}

// Deploy Triggers into OpenWhisk
func (deployer *ServiceDeployer) UnDeployTriggers(deployment *DeploymentApplication) error {

	for _, trigger := range deployment.Triggers {
		if feedname, isFeed := utils.IsFeedAction(trigger); isFeed {
			deployer.deleteFeedAction(trigger, feedname)
		} else {
			deployer.deleteTrigger(trigger)
		}
	}

	return nil

}

// Deploy Rules into OpenWhisk
func (deployer *ServiceDeployer) UnDeployRules(deployment *DeploymentApplication) error {

	for _, rule := range deployment.Rules {
		deployer.deleteRule(rule)

	}
	return nil
}

func (deployer *ServiceDeployer) deletePackage(packa *whisk.SentPackageNoPublish) {
	fmt.Print("Removing package " + packa.Name + " ... ")
	_, err := deployer.Client.Packages.Delete(packa.Name)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error deleteing package with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
	fmt.Println("Done!")
}

func (deployer *ServiceDeployer) deleteTrigger(trigger *whisk.Trigger) {
	fmt.Print("Removing trigger " + trigger.Name + " ... ")
	_, _, err := deployer.Client.Triggers.Delete(trigger.Name)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error deleting trigger with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
	fmt.Println("Done!")
}

func (deployer *ServiceDeployer) deleteFeedAction(trigger *whisk.Trigger, feedName string) {

	params := make(whisk.KeyValueArr, 0)
	params = append(params, whisk.KeyValue{Key: "authKey", Value: deployer.ClientConfig.AuthToken})
	params = append(params, whisk.KeyValue{Key: "lifecycleEvent", Value: "DELETE"})
	params = append(params, whisk.KeyValue{Key: "triggerName", Value: "/" + deployer.Client.Namespace + "/" + trigger.Name})

	trigger.Parameters = nil

	_, _, err := deployer.Client.Triggers.Delete(trigger.Name)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error deleting trigger with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	} else {
		parameters := make(map[string]interface{})
		for _, keyVal := range params {
			parameters[keyVal.Key] = keyVal.Value
		}

		qName, err := utils.ParseQualifiedName(feedName, deployer.ClientConfig.Namespace)

		utils.Check(err)

		namespace := deployer.Client.Namespace
		deployer.Client.Namespace = qName.Namespace
		_, _, err = deployer.Client.Actions.Invoke(qName.EntityName, parameters, true)
		deployer.Client.Namespace = namespace

		if err != nil {
			wskErr := err.(*whisk.WskError)
			fmt.Printf("Got error deleting trigger feed with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
		}
	}
	fmt.Println("Done!")
}

func (deployer *ServiceDeployer) deleteRule(rule *whisk.Rule) {
	fmt.Print("Removing rule " + rule.Name + " ... ")
	_, _, err := deployer.Client.Rules.SetState(rule.Name, "inactive")

	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error deleting rule with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	} else {

		_, err = deployer.Client.Rules.Delete(rule.Name)

		if err != nil {
			wskErr := err.(*whisk.WskError)
			fmt.Printf("Got error deleting rule with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
		}

		fmt.Println("Done!")
	}
}

// Utility function to call go-whisk framework to make action
func (deployer *ServiceDeployer) deleteAction(pkgname string, action *whisk.Action) {
	// call ActionService Thru Client
	if deployer.DeployActionInPackage {
		// the action will be deleted under package with pattern 'packagename/actionname'
		action.Name = strings.Join([]string{pkgname, action.Name}, "/")
	}

	fmt.Print("Removing action " + action.Name + " ... ")
	_, err := deployer.Client.Actions.Delete(action.Name)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error deleting action with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
	fmt.Println("Done!")
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

func (deployer *ServiceDeployer) printDeploymentAssets(assets *DeploymentApplication) {

	// pretty ASCII OpenWhisk graphic
	fmt.Println("         ____      ___                   _    _ _     _     _\n        /\\   \\    / _ \\ _ __   ___ _ __ | |  | | |__ (_)___| | __\n   /\\  /__\\   \\  | | | | '_ \\ / _ \\ '_ \\| |  | | '_ \\| / __| |/ /\n  /  \\____ \\  /  | |_| | |_) |  __/ | | | |/\\| | | | | \\__ \\   <\n  \\   \\  /  \\/    \\___/| .__/ \\___|_| |_|__/\\__|_| |_|_|___/_|\\_\\ \n   \\___\\/              |_|\n")

	fmt.Println("Packages:")
	for _, pack := range assets.Packages {
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
	}

	fmt.Println("Triggers:")
	for _, trigger := range assets.Triggers {
		fmt.Println("* trigger: " + trigger.Name)
		fmt.Println("    bindings: ")

		for _, p := range trigger.Parameters {

			value := "?"
			if str, ok := p.Value.(string); ok {
				value = str
			}
			fmt.Println("        - name: " + p.Key + " value: " + value)
		}

		fmt.Println("    annotations: ")
		for _, p := range trigger.Annotations {

			value := "?"
			if str, ok := p.Value.(string); ok {
				value = str
			}
			fmt.Println("        - name: " + p.Key + " value: " + value)
		}
	}

	fmt.Println("\n Rules")
	for _, rule := range assets.Rules {
		fmt.Println("* rule: " + rule.Name)
		fmt.Println("    - trigger: " + rule.Trigger + "\n    - action: " + rule.Action)
	}

	fmt.Println("")

}
