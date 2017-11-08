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
	"strconv"
	"strings"
	"sync"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"reflect"
)

type DeploymentProject struct {
	Packages map[string]*DeploymentPackage
	Triggers map[string]*whisk.Trigger
	Rules    map[string]*whisk.Rule
	Apis     map[string]*whisk.ApiCreateRequest
}

func NewDeploymentProject() *DeploymentProject {
	var dep DeploymentProject
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
	Deployment      *DeploymentProject
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
	dep.Deployment = NewDeploymentProject()
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
	if err != nil {
		return err
	}

	deployer.RootPackageName = manifest.Package.Packagename

	manifestReader.InitRootPackage(manifestParser, manifest)

	if deployer.IsDefault == true {
		fileReader := NewFileSystemReader(deployer)
		fileActions, err := fileReader.ReadProjectDirectory(manifest)
		if err != nil {
			return err
		}
		fileReader.SetFileActions(fileActions)
	}

	// process manifest file
	err = manifestReader.HandleYaml(deployer, manifestParser, manifest)
	if err != nil {
		return err
	}

	projectName := ""
	if len(manifest.GetProject().Packages) != 0 {
		projectName = manifest.GetProject().Name
	}

	// (TODO) delete this warning after deprecating application in manifest file
	if manifest.Application.Name != "" {
		warningString := wski18n.T("WARNING: application in manifest file will soon be deprecated, please use project instead.\n")
		whisk.Debug(whisk.DbgWarn, warningString)
	}

	// process deployment file
	if utils.FileExists(deployer.DeploymentPath) {
		var deploymentReader = NewDeploymentReader(deployer)
		err = deploymentReader.HandleYaml()

		if err != nil {
			return err
		}

		// (TODO) delete this warning after deprecating application in deployment file
		if deploymentReader.DeploymentDescriptor.Application.Name != "" {
			warningString := wski18n.T("WARNING: application in deployment file will soon be deprecated, please use project instead.\n")
			whisk.Debug(whisk.DbgWarn, warningString)
		}

		// compare the name of the project/application
		if len(deploymentReader.DeploymentDescriptor.GetProject().Packages) != 0 && len(projectName) != 0 {
			projectNameDeploy := deploymentReader.DeploymentDescriptor.GetProject().Name
			if projectNameDeploy != projectName {
				errorString := wski18n.T("The name of the project/application {{.projectNameDeploy}} in deployment file at [{{.deploymentFile}}] does not match the name of the project/application {{.projectNameManifest}}} in manifest file at [{{.manifestFile}}].",
					map[string]interface{}{"projectNameDeploy": projectNameDeploy, "deploymentFile": deployer.DeploymentPath,
						"projectNameManifest": projectName, "manifestFile": deployer.ManifestPath})
				return utils.NewYAMLFormatError(errorString)
			}
		}
		if err := deploymentReader.BindAssets(); err != nil {
			return err
		}
	}

	fmt.Println(utils.GenerateManagedAnnotation(manifest.GetProject().Name, manifest.Filepath))

	return err
}

func (deployer *ServiceDeployer) ConstructUnDeploymentPlan() (*DeploymentProject, error) {

	var manifestReader = NewManfiestReader(deployer)
	manifestReader.IsUndeploy = true
	var err error
	manifest, manifestParser, err := manifestReader.ParseManifest()
	if err != nil {
		return deployer.Deployment, err
	}

	deployer.RootPackageName = manifest.Package.Packagename
	manifestReader.InitRootPackage(manifestParser, manifest)

	// process file system
	if deployer.IsDefault == true {
		fileReader := NewFileSystemReader(deployer)
		fileActions, err := fileReader.ReadProjectDirectory(manifest)
		if err != nil {
			return deployer.Deployment, err
		}

		err = fileReader.SetFileActions(fileActions)
		if err != nil {
			return deployer.Deployment, err
		}

	}

	// process manifest file
	err = manifestReader.HandleYaml(deployer, manifestParser, manifest)
	if err != nil {
		return deployer.Deployment, err
	}

	projectName := ""
	if len(manifest.GetProject().Packages) != 0 {
		projectName = manifest.GetProject().Name
	}

	// (TODO) delete this warning after deprecating application in manifest file
	if manifest.Application.Name != "" {
		warningString := wski18n.T("WARNING: application in manifest file will soon be deprecated, please use project instead.\n")
		whisk.Debug(whisk.DbgWarn, warningString)
	}

	// process deployment file
	if utils.FileExists(deployer.DeploymentPath) {
		var deploymentReader = NewDeploymentReader(deployer)
		err = deploymentReader.HandleYaml()
		if err != nil {
			return deployer.Deployment, err
		}

		// (TODO) delete this warning after deprecating application in deployment file
		if deploymentReader.DeploymentDescriptor.Application.Name != "" {
			warningString := wski18n.T("WARNING: application in deployment file will soon be deprecated, please use project instead.\n")
			whisk.Debug(whisk.DbgWarn, warningString)
		}
		// compare the name of the application
		if len(deploymentReader.DeploymentDescriptor.GetProject().Packages) != 0 && len(projectName) != 0 {
			projectNameDeploy := deploymentReader.DeploymentDescriptor.GetProject().Name
			if projectNameDeploy != projectName {
				errorString := wski18n.T("The name of the project/application {{.projectNameDeploy}} in deployment file at [{{.deploymentFile}}] does not match the name of the application {{.projectNameManifest}}} in manifest file at [{{.manifestFile}}].",
					map[string]interface{}{"projectNameDeploy": projectNameDeploy, "deploymentFile": deployer.DeploymentPath,
						"projectNameManifest": projectName, "manifestFile": deployer.ManifestPath})
				return deployer.Deployment, utils.NewYAMLFormatError(errorString)
			}
		}

		if err := deploymentReader.BindAssets(); err != nil {
			return deployer.Deployment, err
		}
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

			utils.PrintOpenWhiskOutput(wski18n.T("Deployment completed successfully.\n"))
			return nil

		} else {
			deployer.InteractiveChoice = false
			utils.PrintOpenWhiskOutput(wski18n.T("OK. Cancelling deployment.\n"))
			return nil
		}
	}

	// non-interactive
	if err := deployer.deployAssets(); err != nil {
		errString := wski18n.T("Deployment did not complete sucessfully. Run `wskdeploy undeploy` to remove partially deployed assets.\n")
		whisk.Debug(whisk.DbgError, errString)
		return err
	}

	utils.PrintOpenWhiskOutput(wski18n.T("Deployment completed successfully.\n"))
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

	if err := deployer.DeployApis(); err != nil {
		return err
	}

	return nil
}

func (deployer *ServiceDeployer) DeployDependencies() error {
	for _, pack := range deployer.Deployment.Packages {
		for depName, depRecord := range pack.Dependencies {
			output := wski18n.T("Deploying dependency {{.output}} ...",
				map[string]interface{}{"output": depName})
			whisk.Debug(whisk.DbgInfo, output)

			if depRecord.IsBinding {
				bindingPackage := new(whisk.BindingPackage)
				bindingPackage.Namespace = pack.Package.Namespace
				bindingPackage.Name = depName
				pub := false
				bindingPackage.Publish = &pub

				qName, err := utils.ParseQualifiedName(depRecord.Location, pack.Package.Namespace)
				if err != nil {
					return err
				}
				bindingPackage.Binding = whisk.Binding{qName.Namespace, qName.EntityName}

				bindingPackage.Parameters = depRecord.Parameters
				bindingPackage.Annotations = depRecord.Annotations

				error := deployer.createBinding(bindingPackage)
				if error != nil {
					return error
				} else {
					output := wski18n.T("Dependency {{.output}} has been successfully deployed.\n",
						map[string]interface{}{"output": depName})
					whisk.Debug(whisk.DbgInfo, output)
				}

			} else {
				depServiceDeployer, err := deployer.getDependentDeployer(depName, depRecord)
				if err != nil {
					return err
				}

				err = depServiceDeployer.ConstructDeploymentPlan()
				if err != nil {
					return err
				}

				if err := depServiceDeployer.deployAssets(); err != nil {
					errString := wski18n.T("Deployment of dependency {{.depName}} did not complete sucessfully. Run `wskdeploy undeploy` to remove partially deployed assets.\n",
						map[string]interface{}{"depName": depName})
					utils.PrintOpenWhiskErrorMessage(errString)
					return err
				}

				// if the RootPackageName is different from depName
				// create a binding to the origin package
				if depServiceDeployer.RootPackageName != depName {
					bindingPackage := new(whisk.BindingPackage)
					bindingPackage.Namespace = pack.Package.Namespace
					bindingPackage.Name = depName
					pub := false
					bindingPackage.Publish = &pub

					qName, err := utils.ParseQualifiedName(depServiceDeployer.RootPackageName, depServiceDeployer.Deployment.Packages[depServiceDeployer.RootPackageName].Package.Namespace)
					if err != nil {
						return err
					}

					bindingPackage.Binding = whisk.Binding{qName.Namespace, qName.EntityName}

					bindingPackage.Parameters = depRecord.Parameters
					bindingPackage.Annotations = depRecord.Annotations

					err = deployer.createBinding(bindingPackage)
					if err != nil {
						return err
					} else {
						output := wski18n.T("Dependency {{.output}} has been successfully deployed.\n",
							map[string]interface{}{"output": depName})
						whisk.Debug(whisk.DbgInfo, output)
					}
				}
			}
		}
	}

	return nil
}

func (deployer *ServiceDeployer) DeployPackages() error {
	for _, pack := range deployer.Deployment.Packages {
		err := deployer.createPackage(pack.Package)
		if err != nil {
			return err
		}
	}
	return nil
}

// Deploy Sequences into OpenWhisk
func (deployer *ServiceDeployer) DeploySequences() error {

	for _, pack := range deployer.Deployment.Packages {
		for _, action := range pack.Sequences {
			error := deployer.createAction(pack.Package.Name, action.Action)
			if error != nil {
				return error
			}
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
			error := deployer.createFeedAction(trigger, feedname)
			if error != nil {
				return error
			}
		} else {
			error := deployer.createTrigger(trigger)
			if error != nil {
				return error
			}
		}

	}
	return nil

}

// Deploy Rules into OpenWhisk
func (deployer *ServiceDeployer) DeployRules() error {
	for _, rule := range deployer.Deployment.Rules {
		error := deployer.createRule(rule)
		if error != nil {
			return error
		}
	}
	return nil
}

// Deploy Apis into OpenWhisk
func (deployer *ServiceDeployer) DeployApis() error {
	for _, api := range deployer.Deployment.Apis {
		error := deployer.createApi(api)
		if error != nil {
			return error
		}
	}
	return nil
}

func (deployer *ServiceDeployer) createBinding(packa *whisk.BindingPackage) error {
	output := wski18n.T("Deploying package binding {{.output}} ...",
		map[string]interface{}{"output": packa.Name})
	whisk.Debug(whisk.DbgInfo, output)
	_, _, err := deployer.Client.Packages.Insert(packa, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error creating package binding with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
	} else {
		output := wski18n.T("Package binding {{.output}} has been successfully deployed.\n",
			map[string]interface{}{"output": packa.Name})
		whisk.Debug(whisk.DbgInfo, output)
	}
	return nil
}

func (deployer *ServiceDeployer) createPackage(packa *whisk.Package) error {
	output := wski18n.T("Deploying package {{.output}} ...",
		map[string]interface{}{"output": packa.Name})
	whisk.Debug(whisk.DbgInfo, output)
	_, _, err := deployer.Client.Packages.Insert(packa, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error creating package with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
	} else {
		output := wski18n.T("Package {{.output}} has been successfully deployed.\n",
			map[string]interface{}{"output": packa.Name})
		whisk.Debug(whisk.DbgInfo, output)
	}
	return nil
}

func (deployer *ServiceDeployer) createTrigger(trigger *whisk.Trigger) error {
	output := wski18n.T("Deploying trigger {{.output}} ...",
		map[string]interface{}{"output": trigger.Name})
	whisk.Debug(whisk.DbgInfo, output)
	_, _, err := deployer.Client.Triggers.Insert(trigger, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error creating trigger with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
	} else {
		output := wski18n.T("Trigger {{.output}} has been successfully deployed.\n",
			map[string]interface{}{"output": trigger.Name})
		whisk.Debug(whisk.DbgInfo, output)
	}
	return nil
}

func (deployer *ServiceDeployer) createFeedAction(trigger *whisk.Trigger, feedName string) error {
	output := wski18n.T("Deploying trigger feed {{.output}} ...",
		map[string]interface{}{"output": trigger.Name})
	whisk.Debug(whisk.DbgInfo, output)
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

	// triggers created using any of the feeds including cloudant, alarm, message hub etc
	// does not honor UPDATE or overwrite=true with CREATE
	// wskdeploy is designed such that, it updates trigger feeds if they exists
	// or creates new in case they are missing
	// To address trigger feed UPDATE issue, we are checking here if trigger feed
	// exists, if so, delete it and recreate it
	_, r, _ := deployer.Client.Triggers.Get(trigger.Name)
	if r.StatusCode == 200 {
		// trigger feed already exists so first lets delete it and then recreate it
		deployer.deleteFeedAction(trigger, feedName)
	}

	_, _, err := deployer.Client.Triggers.Insert(t, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error creating trigger with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
	} else {

		qName, err := utils.ParseQualifiedName(feedName, deployer.ClientConfig.Namespace)
		if err != nil {
			return err
		}

		namespace := deployer.Client.Namespace
		deployer.Client.Namespace = qName.Namespace
		_, _, err = deployer.Client.Actions.Invoke(qName.EntityName, params, true, false)
		deployer.Client.Namespace = namespace

		if err != nil {
			// Remove the created trigger
			deployer.Client.Triggers.Delete(trigger.Name)
			wskErr := err.(*whisk.WskError)
			errString := wski18n.T("Got error creating trigger feed with error message: {{.err}} and error code: {{.code}}.\n",
				map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
			whisk.Debug(whisk.DbgError, errString)
			return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
		}
	}
	output = wski18n.T("Trigger feed {{.output}} has been successfully deployed.\n",
		map[string]interface{}{"output": trigger.Name})
	whisk.Debug(whisk.DbgInfo, output)
	return nil
}

func (deployer *ServiceDeployer) createRule(rule *whisk.Rule) error {
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
	output := wski18n.T("Deploying rule {{.output}} ...",
		map[string]interface{}{"output": rule.Name})
	whisk.Debug(whisk.DbgInfo, output)

	_, _, err := deployer.Client.Rules.Insert(rule, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error creating rule with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
	}

	_, _, err = deployer.Client.Rules.SetState(rule.Name, "active")
	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error setting the status of rule with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
	}
	output = wski18n.T("Rule {{.output}} has been successfully deployed.\n",
		map[string]interface{}{"output": rule.Name})
	whisk.Debug(whisk.DbgInfo, output)
	return nil
}

// Utility function to call go-whisk framework to make action
func (deployer *ServiceDeployer) createAction(pkgname string, action *whisk.Action) error {
	// call ActionService through the Client
	if deployer.DeployActionInPackage {
		// the action will be created under package with pattern 'packagename/actionname'
		action.Name = strings.Join([]string{pkgname, action.Name}, "/")
	}
	output := wski18n.T("Deploying action {{.output}} ...",
		map[string]interface{}{"output": action.Name})
	whisk.Debug(whisk.DbgInfo, output)

	_, _, err := deployer.Client.Actions.Insert(action, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error creating action with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
	} else {
		output := wski18n.T("Action {{.output}} has been successfully deployed.\n",
			map[string]interface{}{"output": action.Name})
		whisk.Debug(whisk.DbgInfo, output)
	}
	return nil
}

// create api (API Gateway functionality)
func (deployer *ServiceDeployer) createApi(api *whisk.ApiCreateRequest) error {
	_, _, err := deployer.Client.Apis.Insert(api, nil, true)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error creating api with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
	}
	return nil
}

func (deployer *ServiceDeployer) UnDeploy(verifiedPlan *DeploymentProject) error {
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
				utils.PrintOpenWhiskErrorMessage(wski18n.T("Undeployment did not complete sucessfully.\n"))
				return err
			}

			utils.PrintOpenWhiskOutput(wski18n.T("Deployment removed successfully.\n"))
			return nil

		} else {
			deployer.InteractiveChoice = false
			utils.PrintOpenWhiskOutput(wski18n.T("OK. Canceling undeployment.\n"))
			return nil
		}
	}

	// non-interactive
	if err := deployer.unDeployAssets(verifiedPlan); err != nil {
		errString := wski18n.T("Undeployment did not complete sucessfully.\n")
		whisk.Debug(whisk.DbgError, errString)
		return err
	}

	utils.PrintOpenWhiskOutput(wski18n.T("Deployment removed successfully.\n"))
	return nil

}

func (deployer *ServiceDeployer) unDeployAssets(verifiedPlan *DeploymentProject) error {

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
			output := wski18n.T("Undeploying dependency {{.depName}} ...",
				map[string]interface{}{"depName": depName})
			whisk.Debug(whisk.DbgInfo, output)

			if depRecord.IsBinding {
				_, err := deployer.Client.Packages.Delete(depName)
				if err != nil {
					return err
				}
			} else {

				depServiceDeployer, err := deployer.getDependentDeployer(depName, depRecord)
				if err != nil {
					return err
				}

				plan, err := depServiceDeployer.ConstructUnDeploymentPlan()
				if err != nil {
					return err
				}

				// delete binding pkg if the origin package name is different
				if depServiceDeployer.RootPackageName != depName {
					if _, _, ok := deployer.Client.Packages.Get(depName); ok == nil {
						_, err := deployer.Client.Packages.Delete(depName)
						if err != nil {
							wskErr := err.(*whisk.WskError)
							errString := wski18n.T("Got error deleting binding package with error message: {{.err}} and error code: {{.code}}.\n",
								map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
							whisk.Debug(whisk.DbgError, errString)
							return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
						}
					}
				}

				if err := depServiceDeployer.unDeployAssets(plan); err != nil {
					errString := wski18n.T("Undeployment of dependency {{.depName}} did not complete sucessfully.\n",
						map[string]interface{}{"depName": depName})
					whisk.Debug(whisk.DbgError, errString)
					return err
				}
			}
			output = wski18n.T("Dependency {{.depName}} has been successfully undeployed.\n",
				map[string]interface{}{"depName": depName})
			whisk.Debug(whisk.DbgInfo, output)
		}
	}

	return nil
}

func (deployer *ServiceDeployer) UnDeployPackages(deployment *DeploymentProject) error {
	for _, pack := range deployment.Packages {
		err := deployer.deletePackage(pack.Package)
		if err != nil {
			return err
		}
	}
	return nil
}

func (deployer *ServiceDeployer) UnDeploySequences(deployment *DeploymentProject) error {

	for _, pack := range deployment.Packages {
		for _, action := range pack.Sequences {
			err := deployer.deleteAction(pack.Package.Name, action.Action)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// DeployActions into OpenWhisk
func (deployer *ServiceDeployer) UnDeployActions(deployment *DeploymentProject) error {

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
func (deployer *ServiceDeployer) UnDeployTriggers(deployment *DeploymentProject) error {

	for _, trigger := range deployment.Triggers {
		if feedname, isFeed := utils.IsFeedAction(trigger); isFeed {
			err := deployer.deleteFeedAction(trigger, feedname)
			if err != nil {
				return err
			}
		} else {
			err := deployer.deleteTrigger(trigger)
			if err != nil {
				return err
			}
		}
	}

	return nil

}

// Deploy Rules into OpenWhisk
func (deployer *ServiceDeployer) UnDeployRules(deployment *DeploymentProject) error {

	for _, rule := range deployment.Rules {
		err := deployer.deleteRule(rule)
		if err != nil {
			return err
		}
	}
	return nil
}

func (deployer *ServiceDeployer) deletePackage(packa *whisk.Package) error {
	output := wski18n.T("Removing package {{.package}} ...",
		map[string]interface{}{"package": packa.Name})
	whisk.Debug(whisk.DbgInfo, output)
	if _, _, ok := deployer.Client.Packages.Get(packa.Name); ok == nil {
		_, err := deployer.Client.Packages.Delete(packa.Name)
		if err != nil {
			wskErr := err.(*whisk.WskError)
			errString := wski18n.T("Got error deleting package with error message: {{.err}} and error code: {{.code}}.\n",
				map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
			whisk.Debug(whisk.DbgError, errString)
			return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
		}
	}
	return nil
}

func (deployer *ServiceDeployer) deleteTrigger(trigger *whisk.Trigger) error {
	output := wski18n.T("Removing trigger {{.trigger}} ...",
		map[string]interface{}{"trigger": trigger.Name})
	whisk.Debug(whisk.DbgInfo, output)
	_, _, err := deployer.Client.Triggers.Delete(trigger.Name)
	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error deleting trigger with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
	} else {
		output := wski18n.T("Trigger {{.trigger}} has been removed.\n",
			map[string]interface{}{"trigger": trigger.Name})
		whisk.Debug(whisk.DbgInfo, output)
	}
	return nil
}

func (deployer *ServiceDeployer) deleteFeedAction(trigger *whisk.Trigger, feedName string) error {

	params := make(whisk.KeyValueArr, 0)
	params = append(params, whisk.KeyValue{Key: "authKey", Value: deployer.ClientConfig.AuthToken})
	params = append(params, whisk.KeyValue{Key: "lifecycleEvent", Value: "DELETE"})
	params = append(params, whisk.KeyValue{Key: "triggerName", Value: "/" + deployer.Client.Namespace + "/" + trigger.Name})

	parameters := make(map[string]interface{})
	for _, keyVal := range params {
		parameters[keyVal.Key] = keyVal.Value
	}

	qName, err := utils.ParseQualifiedName(feedName, deployer.ClientConfig.Namespace)
	if err != nil {
		return err
	}

	namespace := deployer.Client.Namespace
	deployer.Client.Namespace = qName.Namespace
	_, _, err = deployer.Client.Actions.Invoke(qName.EntityName, parameters, true, true)
	deployer.Client.Namespace = namespace

	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Failed to invoke the feed when deleting trigger feed with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)

	} else {
		trigger.Parameters = nil

		_, _, err := deployer.Client.Triggers.Delete(trigger.Name)
		if err != nil {
			wskErr := err.(*whisk.WskError)
			errString := wski18n.T("Got error deleting trigger with error message: {{.err}} and error code: {{.code}}.\n",
				map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
			whisk.Debug(whisk.DbgError, errString)
			return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
		}
	}

	return nil
}

func (deployer *ServiceDeployer) deleteRule(rule *whisk.Rule) error {
	output := wski18n.T("Removing rule {{.rule}} ...",
		map[string]interface{}{"rule": rule.Name})
	whisk.Debug(whisk.DbgInfo, output)
	_, _, err := deployer.Client.Rules.SetState(rule.Name, "inactive")

	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error setting the status of rule with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
	} else {

		_, err = deployer.Client.Rules.Delete(rule.Name)

		if err != nil {
			wskErr := err.(*whisk.WskError)
			errString := wski18n.T("Got error deleting rule with error message: {{.err}} and error code: {{.code}}.\n",
				map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
			whisk.Debug(whisk.DbgError, errString)
			return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)
		}
	}
	output = wski18n.T("Rule {{.rule}} has been removed.\n",
		map[string]interface{}{"rule": rule.Name})
	whisk.Debug(whisk.DbgInfo, output)
	return nil
}

// Utility function to call go-whisk framework to make action
func (deployer *ServiceDeployer) deleteAction(pkgname string, action *whisk.Action) error {
	// call ActionService Thru Client
	if deployer.DeployActionInPackage {
		// the action will be deleted under package with pattern 'packagename/actionname'
		action.Name = strings.Join([]string{pkgname, action.Name}, "/")
	}

	output := wski18n.T("Removing action {{.action}} ...",
		map[string]interface{}{"action": action.Name})
	whisk.Debug(whisk.DbgInfo, output)

	if _, _, ok := deployer.Client.Actions.Get(action.Name); ok == nil {
		_, err := deployer.Client.Actions.Delete(action.Name)
		if err != nil {
			wskErr := err.(*whisk.WskError)
			errString := wski18n.T("Got error deleting action with error message: {{.err}} and error code: {{.code}}.\n",
				map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
			whisk.Debug(whisk.DbgError, errString)
			return utils.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode)

		}
		output = wski18n.T("Action {{.action}} has been removed.\n",
			map[string]interface{}{"action": action.Name})
		whisk.Debug(whisk.DbgInfo, output)
	}
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

func (deployer *ServiceDeployer) printDeploymentAssets(assets *DeploymentProject) {

	// pretty ASCII OpenWhisk graphic
	utils.PrintOpenWhiskOutputln("         ____      ___                   _    _ _     _     _\n        /\\   \\    / _ \\ _ __   ___ _ __ | |  | | |__ (_)___| | __\n   /\\  /__\\   \\  | | | | '_ \\ / _ \\ '_ \\| |  | | '_ \\| / __| |/ /\n  /  \\____ \\  /  | |_| | |_) |  __/ | | | |/\\| | | | | \\__ \\   <\n  \\   \\  /  \\/    \\___/| .__/ \\___|_| |_|__/\\__|_| |_|_|___/_|\\_\\ \n   \\___\\/              |_|\n")

	utils.PrintOpenWhiskOutputln("Packages:")
	for _, pack := range assets.Packages {
		utils.PrintOpenWhiskOutputln("Name: " + pack.Package.Name)
		utils.PrintOpenWhiskOutputln("    bindings: ")
		for _, p := range pack.Package.Parameters {
			jsonValue, err := utils.PrettyJSON(p.Value)
			if err != nil {
				fmt.Printf("        - %s : %s\n", p.Key, utils.UNKNOWN_VALUE)
			} else {
				fmt.Printf("        - %s : %v\n", p.Key, jsonValue)
			}
		}

		for key, dep := range pack.Dependencies {
			utils.PrintOpenWhiskOutputln("  * dependency: " + key)
			utils.PrintOpenWhiskOutputln("    location: " + dep.Location)
			if !dep.IsBinding {
				utils.PrintOpenWhiskOutputln("    local path: " + dep.ProjectPath)
			}
		}

		utils.PrintOpenWhiskOutputln("")

		for _, action := range pack.Actions {
			utils.PrintOpenWhiskOutputln("  * action: " + action.Action.Name)
			utils.PrintOpenWhiskOutputln("    bindings: ")
			for _, p := range action.Action.Parameters {

				if reflect.TypeOf(p.Value).Kind() == reflect.Map {
					if _, ok := p.Value.(map[interface{}]interface{}); ok {
						var temp map[string]interface{} = utils.ConvertInterfaceMap(p.Value.(map[interface{}]interface{}))
						fmt.Printf("        - %s : %v\n", p.Key, temp)
					} else {
						jsonValue, err := utils.PrettyJSON(p.Value)
						if err != nil {
							fmt.Printf("        - %s : %s\n", p.Key, utils.UNKNOWN_VALUE)
						} else {
							fmt.Printf("        - %s : %v\n", p.Key, jsonValue)
						}
					}
				} else {
					jsonValue, err := utils.PrettyJSON(p.Value)
					if err != nil {
						fmt.Printf("        - %s : %s\n", p.Key, utils.UNKNOWN_VALUE)
					} else {
						fmt.Printf("        - %s : %v\n", p.Key, jsonValue)
					}
				}

			}
			utils.PrintOpenWhiskOutputln("    annotations: ")
			for _, p := range action.Action.Annotations {
				fmt.Printf("        - %s : %v\n", p.Key, p.Value)

			}
		}

		utils.PrintOpenWhiskOutputln("")
		for _, action := range pack.Sequences {
			utils.PrintOpenWhiskOutputln("  * sequence: " + action.Action.Name)
		}

		utils.PrintOpenWhiskOutputln("")
	}

	utils.PrintOpenWhiskOutputln("Triggers:")
	for _, trigger := range assets.Triggers {
		utils.PrintOpenWhiskOutputln("* trigger: " + trigger.Name)
		utils.PrintOpenWhiskOutputln("    bindings: ")

		for _, p := range trigger.Parameters {
			jsonValue, err := utils.PrettyJSON(p.Value)
			if err != nil {
				fmt.Printf("        - %s : %s\n", p.Key, utils.UNKNOWN_VALUE)
			} else {
				fmt.Printf("        - %s : %v\n", p.Key, jsonValue)
			}
		}

		utils.PrintOpenWhiskOutputln("    annotations: ")
		for _, p := range trigger.Annotations {

			value := "?"
			if str, ok := p.Value.(string); ok {
				value = str
			}
			utils.PrintOpenWhiskOutputln("        - name: " + p.Key + " value: " + value)
		}
	}

	utils.PrintOpenWhiskOutputln("\n Rules")
	for _, rule := range assets.Rules {
		utils.PrintOpenWhiskOutputln("* rule: " + rule.Name)
		utils.PrintOpenWhiskOutputln("    - trigger: " + rule.Trigger.(string) + "\n    - action: " + rule.Action.(string))
	}

	utils.PrintOpenWhiskOutputln("")

}

func (deployer *ServiceDeployer) getDependentDeployer(depName string, depRecord utils.DependencyRecord) (*ServiceDeployer, error) {
	depServiceDeployer := NewServiceDeployer()
	projectPath := path.Join(depRecord.ProjectPath, depName+"-"+depRecord.Version)
	if len(depRecord.SubFolder) > 0 {
		projectPath = path.Join(projectPath, depRecord.SubFolder)
	}
	manifestPath := utils.GetManifestFilePath(projectPath)
	deploymentPath := utils.GetDeploymentFilePath(projectPath)
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
