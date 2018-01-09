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
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskderrors"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"net/http"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
)

const (
	CONFLICT_CODE    = 153
	CONFLICT_MESSAGE = "Concurrent modification to resource detected"
	DEFAULT_ATTEMPTS = 3
	DEFAULT_INTERVAL = 1 * time.Second
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
	ProjectName     string
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
	ManagedAnnotation     whisk.KeyValue
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

	var manifestReader = NewManifestReader(deployer)
	manifestReader.IsUndeploy = false
	var err error
	manifest, manifestParser, err := manifestReader.ParseManifest()
	if err != nil {
		return err
	}

	deployer.RootPackageName = manifest.Package.Packagename
	deployer.ProjectName = manifest.GetProject().Name

	// Generate Managed Annotations if its marked as a Managed Deployment
	// Managed deployments are the ones when OpenWhisk entities are deployed with command line flag --managed.
	// Which results in a hidden annotation in every OpenWhisk entity in manifest file.
	if utils.Flags.Managed {
		// OpenWhisk entities are annotated with Project Name and therefore
		// Project Name in manifest/deployment file is mandatory for managed deployments
		if deployer.ProjectName == "" {
			errmsg := wski18n.T(wski18n.ID_ERR_MISSING_MANDATORY_KEY_X_key_X,
				map[string]interface{}{"key": "project name"})

			return wskderrors.NewYAMLFileFormatError(manifest.Filepath, errmsg)
		}
		// Every OpenWhisk entity in the manifest file will be annotated with:
		//managed: '{"__OW__PROJECT__NAME": <name>, "__OW__PROJECT_HASH": <hash>, "__OW__FILE": <path>}'
		deployer.ManagedAnnotation, err = utils.GenerateManagedAnnotation(deployer.ProjectName, manifest.Filepath)
		if err != nil {
			return wskderrors.NewYAMLFileFormatError(manifest.Filepath, err.Error())
		}
	}

	manifestReader.InitRootPackage(manifestParser, manifest, deployer.ManagedAnnotation)

	if deployer.IsDefault == true {
		fileReader := NewFileSystemReader(deployer)
		fileActions, err := fileReader.ReadProjectDirectory(manifest)
		if err != nil {
			return err
		}
		fileReader.SetFileActions(fileActions)
	}

	// process manifest file
	err = manifestReader.HandleYaml(deployer, manifestParser, manifest, deployer.ManagedAnnotation)
	if err != nil {
		return err
	}

	projectName := ""
	if len(manifest.GetProject().Packages) != 0 {
		projectName = manifest.GetProject().Name
	}

	// (TODO) delete this warning after deprecating application in manifest file
	if manifest.Application.Name != "" {
		wskprint.PrintOpenWhiskWarning(wski18n.T(wski18n.ID_WARN_DEPRECATED_KEY_REPLACED,
			map[string]interface{}{"oldkey": "application", "newkey": "project", "filetype": "manifest"}))
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
			wskprint.PrintOpenWhiskWarning(wski18n.T(wski18n.ID_WARN_DEPRECATED_KEY_REPLACED,
				map[string]interface{}{"oldkey": "application", "newkey": "project", "filetype": "deployment"}))
		}

		// compare the name of the project
		if len(deploymentReader.DeploymentDescriptor.GetProject().Packages) != 0 && len(projectName) != 0 {
			projectNameDeploy := deploymentReader.DeploymentDescriptor.GetProject().Name
			if projectNameDeploy != projectName {
				errorString := wski18n.T(wski18n.ID_ERR_MISMATCH_NAME_X_key_X_dname_X_dpath_X_mname_X_moath_X,
					map[string]interface{}{
						"key": "project",
						"dname": projectNameDeploy,
						"dpath": deployer.DeploymentPath,
						"mname": projectName,
						"mpath": deployer.ManifestPath})
				return wskderrors.NewYAMLFileFormatError(manifest.Filepath, errorString)
			}
		}
		if err := deploymentReader.BindAssets(); err != nil {
			return err
		}
	}

	return err
}

func (deployer *ServiceDeployer) ConstructUnDeploymentPlan() (*DeploymentProject, error) {

	var manifestReader = NewManifestReader(deployer)
	manifestReader.IsUndeploy = true
	var err error
	manifest, manifestParser, err := manifestReader.ParseManifest()
	if err != nil {
		return deployer.Deployment, err
	}

	deployer.RootPackageName = manifest.Package.Packagename
	manifestReader.InitRootPackage(manifestParser, manifest, whisk.KeyValue{})

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
	err = manifestReader.HandleYaml(deployer, manifestParser, manifest, whisk.KeyValue{})
	if err != nil {
		return deployer.Deployment, err
	}

	projectName := ""
	if len(manifest.GetProject().Packages) != 0 {
		projectName = manifest.GetProject().Name
	}

	// (TODO) delete this warning after deprecating application in manifest file
	if manifest.Application.Name != "" {
		wskprint.PrintOpenWhiskWarning(wski18n.T(wski18n.ID_WARN_DEPRECATED_KEY_REPLACED,
			map[string]interface{}{"oldkey": "application", "newkey": "project", "filetype": "manifest"}))
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
			wskprint.PrintOpenWhiskWarning(wski18n.T(wski18n.ID_WARN_DEPRECATED_KEY_REPLACED,
				map[string]interface{}{"oldkey": "application", "newkey": "project", "filetype": "deployment"}))
		}

		// compare the name of the application
		if len(deploymentReader.DeploymentDescriptor.GetProject().Packages) != 0 && len(projectName) != 0 {
			projectNameDeploy := deploymentReader.DeploymentDescriptor.GetProject().Name
			if projectNameDeploy != projectName {
				errorString := wski18n.T(wski18n.ID_ERR_MISMATCH_NAME_X_key_X_dname_X_dpath_X_mname_X_moath_X,
					map[string]interface{}{
						"key": "project",
						"dname": projectNameDeploy,
						"dpath": deployer.DeploymentPath,
						"mname": projectName,
						"mpath": deployer.ManifestPath})
				return deployer.Deployment, wskderrors.NewYAMLFileFormatError(manifest.Filepath, errorString)
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

		// TODO() See if we can use the promptForValue() function
		fmt.Print(wski18n.T(wski18n.ID_MSG_PROMPT_DEPLOY))

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "" {
			text = "n"
		}

		// TODO() make possible responses constants (enum?) and create "No" corallary
		if strings.EqualFold(text, "y") || strings.EqualFold(text, "yes") {
			deployer.InteractiveChoice = true
			if err := deployer.deployAssets(); err != nil {
				wskprint.PrintOpenWhiskError(wski18n.T(wski18n.ID_MSG_DEPLOYMENT_FAILED))
				return err
			}

			wskprint.PrintOpenWhiskSuccess(wski18n.T(wski18n.ID_MSG_DEPLOYMENT_SUCCEEDED))
			return nil

		} else {
			// TODO() Should acknowledge if user typed (No/N/n) and if not still exit, but
			// indicate we took the response to mean "No", typically by displaying interpolated
			// response in parenthesis
			deployer.InteractiveChoice = false
			wskprint.PrintOpenWhiskSuccess(wski18n.T(wski18n.ID_MSG_DEPLOYMENT_CANCELLED))
			return nil
		}
	}

	// non-interactive
	if err := deployer.deployAssets(); err != nil {
		wskprint.PrintOpenWhiskError(wski18n.T(wski18n.ID_MSG_DEPLOYMENT_FAILED))
		return err
	}

	wskprint.PrintOpenWhiskSuccess(wski18n.T(wski18n.T(wski18n.ID_MSG_DEPLOYMENT_SUCCEEDED)))
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

	// During managed deployments, after deploying list of entities in a project
	// refresh previously deployed project entities, delete the assets which is no longer part of the project
	// i.e. in a subsequent managed deployment of the same project minus few OpenWhisk entities
	// from the manifest file must result in undeployment of those deleted entities
	if utils.Flags.Managed {
		if err := deployer.RefreshManagedEntities(deployer.ManagedAnnotation); err != nil {
			errString := wski18n.T(wski18n.ID_MSG_UNDEPLOYMENT_MANAGED_FAILED)
			whisk.Debug(whisk.DbgError, errString)
			return err
		}
	}

	return nil
}

func (deployer *ServiceDeployer) DeployDependencies() error {
	for _, pack := range deployer.Deployment.Packages {
		for depName, depRecord := range pack.Dependencies {
			output := wski18n.T(wski18n.ID_MSG_DEPENDENCY_DEPLOYING_X_name_X,
				map[string]interface{}{"name": depName})
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
					output := wski18n.T(wski18n.ID_MSG_DEPENDENCY_DEPLOYMENT_SUCCESS_X_name_X,
						map[string]interface{}{"name": depName})
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
					errString := wski18n.T(wski18n.ID_MSG_DEPENDENCY_DEPLOYMENT_FAILURE_X_name_X,
						map[string]interface{}{"name": depName})
					wskprint.PrintOpenWhiskError(errString)
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
						output := wski18n.T(wski18n.ID_MSG_DEPENDENCY_DEPLOYMENT_SUCCESS_X_name_X,
							map[string]interface{}{"name": depName})
						whisk.Debug(whisk.DbgInfo, output)
					}
				}
			}
		}
	}

	return nil
}

func (deployer *ServiceDeployer) RefreshManagedEntities(maValue whisk.KeyValue) error {

	ma := maValue.Value.(map[string]interface{})
	if err := deployer.RefreshManagedTriggers(ma); err != nil {
		return err
	}

	if err := deployer.RefreshManagedRules(ma); err != nil {
		return err
	}

	if err := deployer.RefreshManagedPackages(ma); err != nil {
		return err
	}

	return nil

}
func (deployer *ServiceDeployer) RefreshManagedActions(packageName string, ma map[string]interface{}) error {
	options := whisk.ActionListOptions{}
	// get a list of actions in your namespace
	actions, _, err := deployer.Client.Actions.List(packageName, &options)
	if err != nil {
		return err
	}
	// iterate over list of actions to find an action with managed annotations
	// check if "managed" annotation is attached to an action
	for _, action := range actions {
		// an annotation with "managed" key indicates that an action was deployed as part of managed deployment
		// if such annotation exists, check if it belongs to the current managed deployment
		// this action has attached managed annotations
		if a := action.Annotations.GetValue(utils.MANAGED); a != nil {
			// decode the JSON blob and retrieve __OW_PROJECT_NAME and __OW_PROJECT_HASH
			aa := a.(map[string]interface{})
			// we have found an action which was earlier part of the current project
			// and this action was deployed as part of managed deployment and now
			// must be undeployed as its not part of the project anymore
			// The annotation with same project name but different project hash indicates
			// that this action is deleted from the project in manifest file
			if aa[utils.OW_PROJECT_NAME] == ma[utils.OW_PROJECT_NAME] && aa[utils.OW_PROJECT_HASH] != ma[utils.OW_PROJECT_HASH] {
				actionName := strings.Join([]string{packageName, action.Name}, "/")
				output := wski18n.T("Found the action {{.action}} which is deleted"+
					" from the current project {{.project}} in manifest file which is being undeployed.\n",
					map[string]interface{}{"action": actionName, "project": aa[utils.OW_PROJECT_NAME]})
				whisk.Debug(whisk.DbgInfo, output)
				var err error
				err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
					_, err := deployer.Client.Actions.Delete(actionName)
					return err
				})

				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (deployer *ServiceDeployer) RefreshManagedTriggers(ma map[string]interface{}) error {
	options := whisk.TriggerListOptions{}
	// Get list of triggers in your namespace
	triggers, _, err := deployer.Client.Triggers.List(&options)
	if err != nil {
		return err
	}
	// iterate over the list of triggers to determine whether any of them was part of managed project
	// and now deleted from manifest file we can determine that from the managed annotation
	// If a trigger has attached managed annotation with the project name equals to the current project name
	// but the project hash is different (project hash differs since the trigger is deleted from the manifest file)
	for _, trigger := range triggers {
		// trigger has attached managed annotation
		if a := trigger.Annotations.GetValue(utils.MANAGED); a != nil {
			// decode the JSON blob and retrieve __OW_PROJECT_NAME and __OW_PROJECT_HASH
			ta := a.(map[string]interface{})
			if ta[utils.OW_PROJECT_NAME] == ma[utils.OW_PROJECT_NAME] && ta[utils.OW_PROJECT_HASH] != ma[utils.OW_PROJECT_HASH] {
				// we have found a trigger which was earlier part of the current project
				output := wski18n.T("Found the trigger {{.trigger}} which is deleted"+
					" from the current project {{.project}} in manifest file which is being undeployed.\n",
					map[string]interface{}{"trigger": trigger.Name, "project": ma[utils.OW_PROJECT_NAME]})
				whisk.Debug(whisk.DbgInfo, output)
				var err error
				err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
					_, _, err := deployer.Client.Triggers.Delete(trigger.Name)
					return err
				})

				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (deployer *ServiceDeployer) RefreshManagedRules(ma map[string]interface{}) error {
	return nil
}

func (deployer *ServiceDeployer) RefreshManagedPackages(ma map[string]interface{}) error {
	options := whisk.PackageListOptions{}
	// Get the list of packages in your namespace
	packages, _, err := deployer.Client.Packages.List(&options)
	if err != nil {
		return err
	}
	// iterate over each package to find managed annotations
	// check if "managed" annotation is attached to a package
	// when managed project name matches with the current project name and project
	// hash differs, indicates that the package was part of the current project but
	// now is deleted from the manifest file and should be undeployed.
	for _, pkg := range packages {
		if a := pkg.Annotations.GetValue(utils.MANAGED); a != nil {
			// decode the JSON blob and retrieve __OW_PROJECT_NAME and __OW_PROJECT_HASH
			pa := a.(map[string]interface{})
			// perform the similar check on the list of actions from this package
			// since package can not be deleted if its not empty (has any action or sequence)
			if err := deployer.RefreshManagedActions(pkg.Name, ma); err != nil {
				return err
			}
			// we have found a package which was earlier part of the current project
			if pa[utils.OW_PROJECT_NAME] == ma[utils.OW_PROJECT_NAME] && pa[utils.OW_PROJECT_HASH] != ma[utils.OW_PROJECT_HASH] {
				output := wski18n.T("Found the package {{.package}} which is deleted"+
					" from the current project {{.project}} in manifest file which is being undeployed.\n",
					map[string]interface{}{"package": pkg.Name, "project": pa[utils.OW_PROJECT_NAME]})
				whisk.Debug(whisk.DbgInfo, output)
				var err error
				err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
					_, err := deployer.Client.Packages.Delete(pkg.Name)
					return err
				})

				if err != nil {
					return err
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
	var err error

	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Packages.Insert(packa, true)
		return err
	})

	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error creating package binding with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)
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
	var err error
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Packages.Insert(packa, true)
		return err
	})
	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error creating package with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)
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
	var err error
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Triggers.Insert(trigger, true)
		return err
	})
	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error creating trigger with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)
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

	var err error
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Triggers.Insert(t, true)
		return err
	})
	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error creating trigger with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)
	} else {

		qName, err := utils.ParseQualifiedName(feedName, deployer.ClientConfig.Namespace)
		if err != nil {
			return err
		}

		namespace := deployer.Client.Namespace
		deployer.Client.Namespace = qName.Namespace
		err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
			_, response, err = deployer.Client.Actions.Invoke(qName.EntityName, params, true, false)
			return err
		})
		deployer.Client.Namespace = namespace

		if err != nil {
			// Remove the created trigger
			deployer.Client.Triggers.Delete(trigger.Name)

			retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
				_, _, err := deployer.Client.Triggers.Delete(trigger.Name)
				return err
			})

			wskErr := err.(*whisk.WskError)
			errString := wski18n.T("Got error creating trigger feed with error message: {{.err}} and error code: {{.code}}.\n",
				map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
			whisk.Debug(whisk.DbgError, errString)
			return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)
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

	var err error
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Rules.Insert(rule, true)
		return err
	})

	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error creating rule with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)
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

	var err error
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Actions.Insert(action, true)
		return err
	})

	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error creating action with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)
	} else {
		output := wski18n.T("Action {{.output}} has been successfully deployed.\n",
			map[string]interface{}{"output": action.Name})
		whisk.Debug(whisk.DbgInfo, output)
	}
	return nil
}

// create api (API Gateway functionality)
func (deployer *ServiceDeployer) createApi(api *whisk.ApiCreateRequest) error {
	var err error
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Apis.Insert(api, nil, true)
		return err
	})

	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error creating api with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)
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
				wskprint.PrintOpenWhiskError(wski18n.T("Undeployment did not complete sucessfully.\n"))
				return err
			}

			wskprint.PrintOpenWhiskSuccess(wski18n.T("Deployment removed successfully.\n"))
			return nil

		} else {
			deployer.InteractiveChoice = false
			wskprint.PrintOpenWhiskSuccess(wski18n.T("OK. Canceling undeployment.\n"))
			return nil
		}
	}

	// non-interactive
	if err := deployer.unDeployAssets(verifiedPlan); err != nil {
		errString := wski18n.T("Undeployment did not complete sucessfully.\n")
		whisk.Debug(whisk.DbgError, errString)
		return err
	}

	wskprint.PrintOpenWhiskSuccess(wski18n.T("Deployment removed successfully.\n"))
	return nil

}

func (deployer *ServiceDeployer) unDeployAssets(verifiedPlan *DeploymentProject) error {

	if err := deployer.UnDeployRules(verifiedPlan); err != nil {
		return err
	}

	if err := deployer.UnDeployTriggers(verifiedPlan); err != nil {
		return err
	}

	if err := deployer.UnDeploySequences(verifiedPlan); err != nil {
		return err
	}

	if err := deployer.UnDeployActions(verifiedPlan); err != nil {
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
			output := wski18n.T(wski18n.ID_MSG_DEPENDENCY_UNDEPLOYING_X_name_X,
				map[string]interface{}{"name": depName})
			whisk.Debug(whisk.DbgInfo, output)

			if depRecord.IsBinding {
				var err error
				err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
					_, err := deployer.Client.Packages.Delete(depName)
					return err
				})
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
						var err error
						var response *http.Response
						err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
							response, err = deployer.Client.Packages.Delete(depName)
							return err
						})
						if err != nil {
							wskErr := err.(*whisk.WskError)
							errString := wski18n.T("Got error deleting binding package with error message: {{.err}} and error code: {{.code}}.\n",
								map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
							whisk.Debug(whisk.DbgError, errString)
							return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)
						}
					}
				}

				if err := depServiceDeployer.unDeployAssets(plan); err != nil {
					errString := wski18n.T(wski18n.ID_MSG_DEPENDENCY_UNDEPLOYMENT_FAILURE_X_name_X,
						map[string]interface{}{"name": depName})
					whisk.Debug(whisk.DbgError, errString)
					return err
				}
			}
			output = wski18n.T(wski18n.ID_MSG_DEPENDENCY_UNDEPLOYMENT_SUCCESS_X_name_X,
				map[string]interface{}{"name": depName})
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
		var err error
		var response *http.Response
		err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
			response, err = deployer.Client.Packages.Delete(packa.Name)
			return err
		})

		if err != nil {
			wskErr := err.(*whisk.WskError)
			errString := wski18n.T("Got error deleting package with error message: {{.err}} and error code: {{.code}}.\n",
				map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
			whisk.Debug(whisk.DbgError, errString)
			return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)
		}
	}
	return nil
}

func (deployer *ServiceDeployer) deleteTrigger(trigger *whisk.Trigger) error {
	output := wski18n.T("Removing trigger {{.trigger}} ...",
		map[string]interface{}{"trigger": trigger.Name})
	whisk.Debug(whisk.DbgInfo, output)

	var err error
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Triggers.Delete(trigger.Name)
		return err
	})

	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error deleting trigger with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)
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
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Actions.Invoke(qName.EntityName, parameters, true, false)
		return err
	})

	deployer.Client.Namespace = namespace

	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Failed to invoke the feed when deleting trigger feed with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)

	} else {
		trigger.Parameters = nil
		var err error
		err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
			_, response, err = deployer.Client.Triggers.Delete(trigger.Name)
			return err
		})

		if err != nil {
			wskErr := err.(*whisk.WskError)
			errString := wski18n.T("Got error deleting trigger with error message: {{.err}} and error code: {{.code}}.\n",
				map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
			whisk.Debug(whisk.DbgError, errString)
			return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)
		}
	}

	return nil
}

func (deployer *ServiceDeployer) deleteRule(rule *whisk.Rule) error {
	output := wski18n.T("Removing rule {{.rule}} ...",
		map[string]interface{}{"rule": rule.Name})
	whisk.Debug(whisk.DbgInfo, output)

	var err error
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		response, err = deployer.Client.Rules.Delete(rule.Name)
		return err
	})

	if err != nil {
		wskErr := err.(*whisk.WskError)
		errString := wski18n.T("Got error deleting rule with error message: {{.err}} and error code: {{.code}}.\n",
			map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)
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
		var err error
		var response *http.Response
		err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
			response, err = deployer.Client.Actions.Delete(action.Name)
			return err
		})

		if err != nil {
			wskErr := err.(*whisk.WskError)
			errString := wski18n.T("Got error deleting action with error message: {{.err}} and error code: {{.code}}.\n",
				map[string]interface{}{"err": wskErr.Error(), "code": strconv.Itoa(wskErr.ExitCode)})
			whisk.Debug(whisk.DbgError, errString)
			return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)

		}
		output = wski18n.T("Action {{.action}} has been removed.\n",
			map[string]interface{}{"action": action.Name})
		whisk.Debug(whisk.DbgInfo, output)
	}
	return nil
}

func retry(attempts int, sleep time.Duration, callback func() error) error {
	var err error
	for i := 0; ; i++ {
		err = callback()
		if i >= (attempts - 1) {
			break
		}
		if err != nil {
			wskErr := err.(*whisk.WskError)
			if wskErr.ExitCode == CONFLICT_CODE && strings.Contains(wskErr.Error(), CONFLICT_MESSAGE) {
				time.Sleep(sleep)
				whisk.Debug(whisk.DbgError, "Retrying [%s] after error: %s\n", strconv.Itoa(i+1), err)
			} else {
				return err
			}
		} else {
			return err
		}
	}
	return err
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
	// TODO() move to separate function and suppress using some flag
	wskprint.PrintlnOpenWhiskOutput("         ____      ___                   _    _ _     _     _\n        /\\   \\    / _ \\ _ __   ___ _ __ | |  | | |__ (_)___| | __\n   /\\  /__\\   \\  | | | | '_ \\ / _ \\ '_ \\| |  | | '_ \\| / __| |/ /\n  /  \\____ \\  /  | |_| | |_) |  __/ | | | |/\\| | | | | \\__ \\   <\n  \\   \\  /  \\/    \\___/| .__/ \\___|_| |_|__/\\__|_| |_|_|___/_|\\_\\ \n   \\___\\/              |_|\n")

	wskprint.PrintlnOpenWhiskOutput("Packages:")
	for _, pack := range assets.Packages {
		wskprint.PrintlnOpenWhiskOutput("Name: " + pack.Package.Name)
		wskprint.PrintlnOpenWhiskOutput("    bindings: ")
		for _, p := range pack.Package.Parameters {
			jsonValue, err := utils.PrettyJSON(p.Value)
			if err != nil {
				fmt.Printf("        - %s : %s\n", p.Key, wskderrors.STR_UNKNOWN_VALUE)
			} else {
				fmt.Printf("        - %s : %v\n", p.Key, jsonValue)
			}
		}

		for key, dep := range pack.Dependencies {
			wskprint.PrintlnOpenWhiskOutput("  * dependency: " + key)
			wskprint.PrintlnOpenWhiskOutput("    location: " + dep.Location)
			if !dep.IsBinding {
				wskprint.PrintlnOpenWhiskOutput("    local path: " + dep.ProjectPath)
			}
		}

		wskprint.PrintlnOpenWhiskOutput("")

		for _, action := range pack.Actions {
			wskprint.PrintlnOpenWhiskOutput("  * action: " + action.Action.Name)
			wskprint.PrintlnOpenWhiskOutput("    bindings: ")
			for _, p := range action.Action.Parameters {

				if reflect.TypeOf(p.Value).Kind() == reflect.Map {
					if _, ok := p.Value.(map[interface{}]interface{}); ok {
						var temp map[string]interface{} = utils.ConvertInterfaceMap(p.Value.(map[interface{}]interface{}))
						fmt.Printf("        - %s : %v\n", p.Key, temp)
					} else {
						jsonValue, err := utils.PrettyJSON(p.Value)
						if err != nil {
							fmt.Printf("        - %s : %s\n", p.Key, wskderrors.STR_UNKNOWN_VALUE)
						} else {
							fmt.Printf("        - %s : %v\n", p.Key, jsonValue)
						}
					}
				} else {
					jsonValue, err := utils.PrettyJSON(p.Value)
					if err != nil {
						fmt.Printf("        - %s : %s\n", p.Key, wskderrors.STR_UNKNOWN_VALUE)
					} else {
						fmt.Printf("        - %s : %v\n", p.Key, jsonValue)
					}
				}

			}
			wskprint.PrintlnOpenWhiskOutput("    annotations: ")
			for _, p := range action.Action.Annotations {
				fmt.Printf("        - %s : %v\n", p.Key, p.Value)

			}
		}

		wskprint.PrintlnOpenWhiskOutput("")
		for _, action := range pack.Sequences {
			wskprint.PrintlnOpenWhiskOutput("  * sequence: " + action.Action.Name)
		}

		wskprint.PrintlnOpenWhiskOutput("")
	}

	wskprint.PrintlnOpenWhiskOutput("Triggers:")
	for _, trigger := range assets.Triggers {
		wskprint.PrintlnOpenWhiskOutput("* trigger: " + trigger.Name)
		wskprint.PrintlnOpenWhiskOutput("    bindings: ")

		for _, p := range trigger.Parameters {
			jsonValue, err := utils.PrettyJSON(p.Value)
			if err != nil {
				fmt.Printf("        - %s : %s\n", p.Key, wskderrors.STR_UNKNOWN_VALUE)
			} else {
				fmt.Printf("        - %s : %v\n", p.Key, jsonValue)
			}
		}

		wskprint.PrintlnOpenWhiskOutput("    annotations: ")
		for _, p := range trigger.Annotations {

			value := "?"
			if str, ok := p.Value.(string); ok {
				value = str
			}
			wskprint.PrintlnOpenWhiskOutput("        - name: " + p.Key + " value: " + value)
		}
	}

	wskprint.PrintlnOpenWhiskOutput("\n Rules")
	for _, rule := range assets.Rules {
		wskprint.PrintlnOpenWhiskOutput("* rule: " + rule.Name)
		wskprint.PrintlnOpenWhiskOutput("    - trigger: " + rule.Trigger.(string) + "\n    - action: " + rule.Action.(string))
	}

	wskprint.PrintlnOpenWhiskOutput("")

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
