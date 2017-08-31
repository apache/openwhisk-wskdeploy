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

package deployers

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
    "github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
)

type DeploymentApplication struct {
	Packages map[string]*DeploymentPackage
	Triggers map[string]*whisk.Trigger
	Rules    map[string]*whisk.Rule
	Apis     map[string]*whisk.ApiCreateRequest
}

func NewDeploymentApplication() *DeploymentApplication {
	var dep DeploymentApplication
	dep.Packages = make(map[string]*DeploymentPackage)
	dep.Triggers = make(map[string]*whisk.Trigger)
	dep.Rules = make(map[string]*whisk.Rule)
	dep.Apis = make(map[string]*whisk.ApiCreateRequest)
	return &dep
}

type DeploymentPackage struct {
	Package      *whisk.Package
	Dependencies map[string]utils.DependencyRecord
	Actions      map[string]utils.ActionRecord
	Sequences    map[string]utils.ActionRecord
}

func NewDeploymentPackage() *DeploymentPackage {
	var dep DeploymentPackage
	dep.Dependencies = make(map[string]utils.DependencyRecord)
	dep.Actions = make(map[string]utils.ActionRecord)
	dep.Sequences = make(map[string]utils.ActionRecord)
	return &dep
}

// ServiceDeployer defines a prototype service deployer.
// It should do the following:
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
	// whether to deploy the action under the package
	DeployActionInPackage bool
	InteractiveChoice     bool
	ClientConfig          *whisk.Config
	DependencyMaster      map[string]utils.DependencyRecord
}

// NewServiceDeployer is a Factory to create a new ServiceDeployer
func NewServiceDeployer() *ServiceDeployer {
	var dep ServiceDeployer
	dep.Deployment = NewDeploymentApplication()
	dep.IsInteractive = true
	dep.DeployActionInPackage = true
	dep.DependencyMaster = make(map[string]utils.DependencyRecord)

	return &dep
}

// Check if the manifest yaml could be parsed by Manifest Parser.
// Check if the deployment yaml could be parsed by Manifest Parser.
func (deployer *ServiceDeployer) Check() {
	ps := parsers.NewYAMLParser()
	if utils.FileExists(deployer.DeploymentPath) {
		ps.ParseDeployment(deployer.DeploymentPath)
	}
	ps.ParseManifest(deployer.ManifestPath)
	// add more schema check or manifest/deployment consistency checks here if
	// necessary
}

func (deployer *ServiceDeployer) ConstructDeploymentPlan() error {

	var manifestReader = NewManfiestReader(deployer)
	manifestReader.IsUndeploy = false
    var err error
	manifest, manifestParser, err := manifestReader.ParseManifest()
	utils.Check(err)

	deployer.RootPackageName = manifest.Package.Packagename

	manifestReader.InitRootPackage(manifestParser, manifest)

	if deployer.IsDefault == true {
		fileReader := NewFileSystemReader(deployer)
		fileActions, err := fileReader.ReadProjectDirectory(manifest)
		utils.Check(err)

		fileReader.SetFileActions(fileActions)
	}

	// process manifest file
	err = manifestReader.HandleYaml(deployer, manifestParser, manifest)
	utils.Check(err)

	// process deploymet file
	if utils.FileExists(deployer.DeploymentPath) {
		var deploymentReader = NewDeploymentReader(deployer)
		deploymentReader.HandleYaml()

		deploymentReader.BindAssets()
	}

	return err
}

func (deployer *ServiceDeployer) ConstructUnDeploymentPlan() (*DeploymentApplication, error) {

	var manifestReader = NewManfiestReader(deployer)
	manifestReader.IsUndeploy = true
    var err error
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
	err = manifestReader.HandleYaml(deployer, manifestParser, manifest)
	utils.Check(err)

	// process deployment file
	if utils.FileExists(deployer.DeploymentPath) {
		var deploymentReader = NewDeploymentReader(deployer)
		deploymentReader.HandleYaml()

		deploymentReader.BindAssets()
	}

	verifiedPlan := deployer.Deployment

	return verifiedPlan, err
}

// Use reflect util to deploy everything in this service deployer
// TODO(TBD): according to some planning?
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
                errString := wski18n.T("Deployment did not complete sucessfully. Run `wskdeploy undeploy` to remove partially deployed assets.\n")
                whisk.Debug(whisk.DbgError, errString)
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
        errString := wski18n.T("Deployment did not complete sucessfully. Run `wskdeploy undeploy` to remove partially deployed assets.\n")
        whisk.Debug(whisk.DbgError, errString)
		return err
	}

	fmt.Println("\nDeployment completed successfully.")
	return nil

}

func (deployer *ServiceDeployer) deployAssets() error {

	if err := deployer.DeployPackages(); err != nil {
		return err
	}

	if err := deployer.DeployDependencies(); err != nil {
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

	if len(deployer.Deployment.Apis) != 0 {
		if err := deployer.DeployApis(); err != nil {
			return err
		}
	}

	return nil
}

func (deployer *ServiceDeployer) DeployDependencies() error {
	for _, pack := range deployer.Deployment.Packages {
		for depName, depRecord := range pack.Dependencies {
			fmt.Println("Deploying dependency " + depName + " ... ")

			if depRecord.IsBinding {
				bindingPackage := new(whisk.BindingPackage)
				bindingPackage.Namespace = pack.Package.Namespace
				bindingPackage.Name = depName
				pub := false
				bindingPackage.Publish = &pub

				qName, err := utils.ParseQualifiedName(depRecord.Location, pack.Package.Namespace)
				utils.Check(err)
				bindingPackage.Binding = whisk.Binding{qName.Namespace, qName.EntityName}

				bindingPackage.Parameters = depRecord.Parameters
				bindingPackage.Annotations = depRecord.Annotations

				deployer.createBinding(bindingPackage)

			} else {
				depServiceDeployer, err := deployer.getDependentDeployer(depName, depRecord)
				utils.Check(err)

				err = depServiceDeployer.ConstructDeploymentPlan()
				utils.Check(err)

				if err := depServiceDeployer.deployAssets(); err != nil {
					fmt.Println("\nDeployment of dependency " + depName + " did not complete sucessfully. Run `wskdeploy undeploy` to remove partially deployed assets")
					return err
				} else {
					fmt.Println("Done!")
				}
			}
		}
	}

	return nil
}

func (deployer *ServiceDeployer) DeployPackages() error {
	for _, pack := range deployer.Deployment.Packages {
		deployer.createPackage(pack.Package)
	}
	return nil
}

// Deploy Sequences into OpenWhisk
func (deployer *ServiceDeployer) DeploySequences() error {

	for _, pack := range deployer.Deployment.Packages {
		for _, action := range pack.Sequences {
			deployer.createAction(pack.Package.Name, action.Action)
		}
	}
	return nil
}

// Deploy Actions into OpenWhisk
func (deployer *ServiceDeployer) DeployActions() error {

	for _, pack := range deployer.Deployment.Packages {
		for _, action := range pack.Actions {
			err := deployer.createAction(pack.Package.Name, action.Action)
			if err != nil {
				return err
			}
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

// Deploy Apis into OpenWhisk
func (deployer *ServiceDeployer) DeployApis() error {
	for _, api := range deployer.Deployment.Apis {
		deployer.createApi(api)
	}
	return nil
}

func (deployer *ServiceDeployer) createBinding(packa *whisk.BindingPackage) {
	fmt.Print("Deploying package binding " + packa.Name + " ... ")
	_, _, err := deployer.Client.Packages.Insert(packa, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating package binding with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	}
	fmt.Println("Done!")
}

func (deployer *ServiceDeployer) createPackage(packa *whisk.Package) {
	fmt.Print("Deploying package " + packa.Name + " ... ")
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
	fmt.Println("Deploying trigger feed " + trigger.Name + " ... ")
	// to hold and modify trigger parameters, not passed by ref?
	params := make(map[string]interface{})

	// check for strings that are JSON
	for _, keyVal := range trigger.Parameters {
		params[keyVal.Key] = keyVal.Value
	}

	params["authKey"] = deployer.ClientConfig.AuthToken
	params["lifecycleEvent"] = "CREATE"
	params["triggerName"] = "/" + deployer.Client.Namespace + "/" + trigger.Name

	pub := true
	t := &whisk.Trigger{
		Name:        trigger.Name,
		Annotations: trigger.Annotations,
		Publish:     &pub,
	}

	_, _, err := deployer.Client.Triggers.Insert(t, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating trigger with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
	} else {

		qName, err := utils.ParseQualifiedName(feedName, deployer.ClientConfig.Namespace)

		utils.Check(err)

		namespace := deployer.Client.Namespace
		deployer.Client.Namespace = qName.Namespace
		_, _, err = deployer.Client.Actions.Invoke(qName.EntityName, params, true, false)
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
	rule.Trigger = deployer.getQualifiedName(rule.Trigger.(string), deployer.ClientConfig.Namespace)
	// The rule's action should include the namespace and package
	// with pattern /namespace/package/action
	// TODO(TBD): please refer https://github.com/openwhisk/openwhisk/issues/1577

	// if it contains a slash, then the action is qualified by a package name
	if strings.Contains(rule.Action.(string), "/") {
		rule.Action = deployer.getQualifiedName(rule.Action.(string), deployer.ClientConfig.Namespace)
	} else {
		// if not, we assume the action is inside the root package
		rule.Action = deployer.getQualifiedName(strings.Join([]string{deployer.RootPackageName, rule.Action.(string)}, "/"), deployer.ClientConfig.Namespace)
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
func (deployer *ServiceDeployer) createAction(pkgname string, action *whisk.Action) error {
	// call ActionService through the Client
	if deployer.DeployActionInPackage {
		// the action will be created under package with pattern 'packagename/actionname'
		action.Name = strings.Join([]string{pkgname, action.Name}, "/")
	}
	fmt.Print("Deploying action " + action.Name + " ... ")
	_, _, err := deployer.Client.Actions.Insert(action, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating action with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
		return err
	}
	fmt.Println("Done!")
	return nil
}

// create api (API Gateway functionality)
func (deployer *ServiceDeployer) createApi(api *whisk.ApiCreateRequest) {
	_, _, err := deployer.Client.Apis.Insert(api, nil, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		fmt.Printf("Got error creating api with error message: %v and error code: %v.\n", wskErr.Error(), wskErr.ExitCode)
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
			fmt.Println("OK. Canceling undeployment")
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

	if err := deployer.UnDeployDependencies(); err != nil {
		return err
	}

	return nil

}

func (deployer *ServiceDeployer) UnDeployDependencies() error {
	for _, pack := range deployer.Deployment.Packages {
		for depName, depRecord := range pack.Dependencies {
			fmt.Println("Undeploying dependency " + depName + " ... ")

			if depRecord.IsBinding {
				_, err := deployer.Client.Packages.Delete(depName)
				utils.Check(err)
			} else {

				depServiceDeployer, err := deployer.getDependentDeployer(depName, depRecord)
				utils.Check(err)

				plan, err := depServiceDeployer.ConstructUnDeploymentPlan()
				utils.Check(err)

				if err := depServiceDeployer.unDeployAssets(plan); err != nil {
					fmt.Println("\nUndeployment of dependency " + depName + " did not complete sucessfully. Run `wskdeploy undeploy` to remove partially deployed assets")
					return err
				} else {
					fmt.Println("Done!")
				}
			}
		}
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
			err := deployer.deleteAction(pack.Package.Name, action.Action)
			if err != nil {
				return err
			}
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

func (deployer *ServiceDeployer) deletePackage(packa *whisk.Package) {
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
		_, _, err = deployer.Client.Actions.Invoke(qName.EntityName, parameters, true, true)
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
func (deployer *ServiceDeployer) deleteAction(pkgname string, action *whisk.Action) error {
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
		return err
	}
	fmt.Println("Done!")
	return nil
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
		fmt.Println("    bindings: ")
		for _, p := range pack.Package.Parameters {
			fmt.Printf("        - %s : %v\n", p.Key, utils.PrettyJSON(p.Value))
		}

		for key, dep := range pack.Dependencies {
			fmt.Println("  * dependency: " + key)
			fmt.Println("    location: " + dep.Location)
			if !dep.IsBinding {
				fmt.Println("    local path: " + dep.ProjectPath)
			}
		}

		fmt.Println("")

		for _, action := range pack.Actions {
			fmt.Println("  * action: " + action.Action.Name)
			fmt.Println("    bindings: ")
			for _, p := range action.Action.Parameters {
				fmt.Printf("        - %s : %v\n", p.Key, utils.PrettyJSON(p.Value))
			}
			fmt.Println("    annotations: ")
			for _, p := range action.Action.Annotations {
				fmt.Printf("        - %s : %v\n", p.Key, p.Value)

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
			fmt.Printf("        - %s : %v\n", p.Key, utils.PrettyJSON(p.Value))
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
		fmt.Println("    - trigger: " + rule.Trigger.(string) + "\n    - action: " + rule.Action.(string))
	}

	fmt.Println("")

}

func (deployer *ServiceDeployer) getDependentDeployer(depName string, depRecord utils.DependencyRecord) (*ServiceDeployer, error) {
	depServiceDeployer := NewServiceDeployer()
	projectPath := path.Join(depRecord.ProjectPath, depName+"-"+depRecord.Version)
	manifestPath := path.Join(projectPath, utils.ManifestFileNameYml)
	deploymentPath := path.Join(projectPath, utils.DeploymentFileNameYaml)
	depServiceDeployer.ProjectPath = projectPath
	depServiceDeployer.ManifestPath = manifestPath
	depServiceDeployer.DeploymentPath = deploymentPath
	depServiceDeployer.IsInteractive = true

	depServiceDeployer.Client = deployer.Client
	depServiceDeployer.ClientConfig = deployer.ClientConfig

	depServiceDeployer.DependencyMaster = deployer.DependencyMaster

	// share the master dependency list
	depServiceDeployer.DependencyMaster = deployer.DependencyMaster

	return depServiceDeployer, nil
}
