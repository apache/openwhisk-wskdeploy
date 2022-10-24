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
	"encoding/json"
	"fmt"

	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apache/openwhisk-wskdeploy/webaction"
	reflect "github.com/goccy/go-reflect"

	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/apache/openwhisk-wskdeploy/dependencies"
	"github.com/apache/openwhisk-wskdeploy/parsers"
	"github.com/apache/openwhisk-wskdeploy/utils"
	"github.com/apache/openwhisk-wskdeploy/wskderrors"
	"github.com/apache/openwhisk-wskdeploy/wski18n"
	"github.com/apache/openwhisk-wskdeploy/wskprint"
)

const (
	CONFLICT_MESSAGE = "Concurrent modification to resource detected"
	CONFLICT_CODE    = 153
	DEFAULT_ATTEMPTS = 3
	DEFAULT_INTERVAL = 1 * time.Second
)

type DeploymentProject struct {
	Packages          map[string]*DeploymentPackage
	Triggers          map[string]*whisk.Trigger
	Rules             map[string]*whisk.Rule
	Apis              map[string]*whisk.ApiCreateRequest
	ApiOptions        map[string]*whisk.ApiCreateRequestOptions
	SwaggerApi        *whisk.ApiCreateRequest
	SwaggerApiOptions *whisk.ApiCreateRequestOptions
}

func NewDeploymentProject() *DeploymentProject {
	var dep DeploymentProject
	dep.Packages = make(map[string]*DeploymentPackage)
	dep.Triggers = make(map[string]*whisk.Trigger)
	dep.Rules = make(map[string]*whisk.Rule)
	dep.Apis = make(map[string]*whisk.ApiCreateRequest)
	dep.ApiOptions = make(map[string]*whisk.ApiCreateRequestOptions)
	return &dep
}

type DeploymentPackage struct {
	Package      *whisk.Package
	Dependencies map[string]dependencies.DependencyRecord
	Actions      map[string]utils.ActionRecord
	Sequences    map[string]utils.ActionRecord
	Inputs       parsers.PackageInputs
}

func NewDeploymentPackage() *DeploymentPackage {
	var dep DeploymentPackage
	dep.Dependencies = make(map[string]dependencies.DependencyRecord)
	dep.Actions = make(map[string]utils.ActionRecord)
	dep.Sequences = make(map[string]utils.ActionRecord)
	dep.Inputs = parsers.PackageInputs{}
	return &dep
}

// ServiceDeployer defines a prototype service deployer.
// It should do the following:
//   1. Collect information from the manifest file (if any)
//   2. Collect information from the deployment file (if any)
//   3. Collect information about the source code files in the working directory
//   4. Create a deployment plan to create OpenWhisk service
type ServiceDeployer struct {
	ProjectName       string
	ProjectInputs     map[string]parsers.Parameter
	Deployment        *DeploymentProject
	Client            *whisk.Client
	mt                sync.RWMutex
	Preview           bool
	Report            bool
	ManifestPath      string
	ProjectPath       string
	DeploymentPath    string
	ClientConfig      *whisk.Config
	DependencyMaster  map[string]dependencies.DependencyRecord
	ManagedAnnotation whisk.KeyValue
}

// NewServiceDeployer is a Factory to create a new ServiceDeployer
func NewServiceDeployer() *ServiceDeployer {
	var dep ServiceDeployer
	dep.Deployment = NewDeploymentProject()
	dep.Preview = true
	dep.DependencyMaster = make(map[string]dependencies.DependencyRecord)
	dep.ProjectInputs = make(map[string]parsers.Parameter, 0)
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

func (deployer *ServiceDeployer) setProjectInputs(manifest *parsers.YAML) error {
	for parameterName, param := range manifest.Project.Inputs {
		p, err := parsers.ResolveParameter(parameterName, &param, manifest.Filepath)
		if err != nil {
			return err
		}
		param.Value = p
		deployer.ProjectInputs[parameterName] = param
	}
	return nil
}

func (deployer *ServiceDeployer) ConstructDeploymentPlan() error {

	var manifestReader = NewManifestReader(deployer)
	manifestReader.IsUndeploy = false
	var err error
	manifest, manifestParser, err := manifestReader.ParseManifest()
	if err != nil {
		return err
	}

	deployer.ProjectName = utils.Flags.ProjectName
	if deployer.ProjectName == "" {
		deployer.ProjectName = manifest.GetProject().Name
	} else {
		warningString := wski18n.T(
			wski18n.ID_WARN_PROJECT_NAME_OVERRIDDEN,
			map[string]interface{}{
				wski18n.KEY_PROJECT: deployer.ProjectName})
		wskprint.PrintOpenWhiskWarning(warningString)
	}

	err = deployer.setProjectInputs(manifest)
	if err != nil {
		return err
	}

	// Generate Managed Annotations if its marked as a Managed Deployment
	// Managed deployments are the ones when OpenWhisk entities are deployed with command line flag --managed.
	// Which results in a hidden annotation in every OpenWhisk entity in manifest file.
	if utils.Flags.Managed || utils.Flags.Sync {
		// OpenWhisk entities are annotated with Project Name and therefore
		// Project Name in manifest/deployment file is mandatory for managed deployments
		if deployer.ProjectName == "" {
			errmsg := wski18n.T(wski18n.ID_ERR_KEY_MISSING_X_key_X,
				map[string]interface{}{wski18n.KEY_KEY: wski18n.NAME_PROJECT})

			return wskderrors.NewYAMLFileFormatError(manifest.Filepath, errmsg)
		}
		// Every OpenWhisk entity in the manifest file will be annotated with:
		//managed: '{"__OW__PROJECT__NAME": <name>, "__OW__PROJECT_HASH": <hash>, "__OW__FILE": <path>}'
		deployer.ManagedAnnotation, err = utils.GenerateManagedAnnotation(deployer.ProjectName, manifest.Filepath)
		if err != nil {
			return wskderrors.NewYAMLFileFormatError(manifest.Filepath, err.Error())
		}
	}

	err = manifestReader.InitPackages(manifestParser, manifest, deployer.ManagedAnnotation)
	if err != nil {
		return err
	}

	projectName := ""
	if len(manifest.GetProject().Packages) != 0 {
		projectName = manifest.GetProject().Name
	}

	// process deployment file
	var deploymentReader = NewDeploymentReader(deployer)
	if utils.FileExists(deployer.DeploymentPath) {
		err = deploymentReader.HandleYaml()
		if err != nil {
			return err
		}

		// compare the name of the project
		if len(deploymentReader.DeploymentDescriptor.GetProject().Packages) != 0 && len(projectName) != 0 {
			projectNameDeploy := deploymentReader.DeploymentDescriptor.GetProject().Name
			if projectNameDeploy != projectName {
				errorString := wski18n.T(wski18n.ID_ERR_NAME_MISMATCH_X_key_X_dname_X_dpath_X_mname_X_moath_X,
					map[string]interface{}{
						wski18n.KEY_KEY:             parsers.YAML_KEY_PROJECT,
						wski18n.KEY_DEPLOYMENT_NAME: projectNameDeploy,
						wski18n.KEY_DEPLOYMENT_PATH: deployer.DeploymentPath,
						wski18n.KEY_MANIFEST_NAME:   projectName,
						wski18n.KEY_MANIFEST_PATH:   deployer.ManifestPath})
				return wskderrors.NewYAMLFileFormatError(manifest.Filepath, errorString)
			}
		}
	}

	// overwrite package inputs based on command line parameters
	// overwrite package inputs with the values from command line --param and/or --param-file in order
	err = deployer.UpdatePackageInputs()
	if err != nil {
		return err
	}

	// process manifest file
	err = manifestReader.HandleYaml(manifestParser, manifest, deployer.ManagedAnnotation)
	if err != nil {
		return err
	}

	// process deployment file
	if utils.FileExists(deployer.DeploymentPath) {
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

	manifestReader.InitPackages(manifestParser, manifest, whisk.KeyValue{})

	// process manifest file
	err = manifestReader.HandleYaml(manifestParser, manifest, whisk.KeyValue{})
	if err != nil {
		return deployer.Deployment, err
	}

	projectName := ""
	if len(manifest.GetProject().Packages) != 0 {
		projectName = manifest.GetProject().Name
	}

	// process deployment file
	if utils.FileExists(deployer.DeploymentPath) {
		var deploymentReader = NewDeploymentReader(deployer)
		err = deploymentReader.HandleYaml()
		if err != nil {
			return deployer.Deployment, err
		}

		// compare the name of the project
		if len(deploymentReader.DeploymentDescriptor.GetProject().Packages) != 0 && len(projectName) != 0 {
			projectNameDeploy := deploymentReader.DeploymentDescriptor.GetProject().Name
			if projectNameDeploy != projectName {
				errorString := wski18n.T(wski18n.ID_ERR_NAME_MISMATCH_X_key_X_dname_X_dpath_X_mname_X_moath_X,
					map[string]interface{}{
						wski18n.KEY_KEY:             parsers.YAML_KEY_PROJECT,
						wski18n.KEY_DEPLOYMENT_NAME: projectNameDeploy,
						wski18n.KEY_DEPLOYMENT_PATH: deployer.DeploymentPath,
						wski18n.KEY_MANIFEST_NAME:   projectName,
						wski18n.KEY_MANIFEST_PATH:   deployer.ManifestPath})
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

	if deployer.Preview {
		deployer.printDeploymentAssets(deployer.Deployment)
		return nil
	}

	if deployer.Report {
		deployer.reportInputs()
		return nil
	}

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
	if utils.Flags.Managed || utils.Flags.Sync {
		if err := deployer.RefreshManagedEntities(deployer.ManagedAnnotation); err != nil {
			errString := wski18n.T(wski18n.ID_MSG_MANAGED_UNDEPLOYMENT_FAILED)
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
				map[string]interface{}{wski18n.KEY_NAME: depName})
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
						map[string]interface{}{wski18n.KEY_NAME: depName})
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

				dependentPackages := []string{}
				for k := range depServiceDeployer.Deployment.Packages {
					dependentPackages = append(dependentPackages, k)
				}

				if len(dependentPackages) > 1 {
					// TODO(799) i18n
					errMessage := "GitHub dependency " + depName + " has multiple packages in manifest file: " +
						strings.Join(dependentPackages, ", ") + ". " +
						"One GitHub dependency can only be associated with single package in manifest file." +
						"There is no way to reference actions from multiple packages of any GitHub dependencies."
					return wskderrors.NewYAMLFileFormatError(deployer.ManifestPath, errMessage)
				}

				if err := depServiceDeployer.deployAssets(); err != nil {
					errString := wski18n.T(wski18n.ID_MSG_DEPENDENCY_DEPLOYMENT_FAILURE_X_name_X,
						map[string]interface{}{wski18n.KEY_NAME: depName})
					wskprint.PrintOpenWhiskError(errString)
					return err
				}

				// if the dependency name in the original package
				// is different from the package name in the manifest
				// file of dependent github repo, create a binding to the origin package
				if ok := depServiceDeployer.Deployment.Packages[depName]; ok == nil {
					bindingPackage := new(whisk.BindingPackage)
					bindingPackage.Namespace = pack.Package.Namespace
					bindingPackage.Name = depName
					pub := false
					bindingPackage.Publish = &pub

					qName, err := utils.ParseQualifiedName(dependentPackages[0], depServiceDeployer.Deployment.Packages[dependentPackages[0]].Package.Namespace)
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
							map[string]interface{}{wski18n.KEY_NAME: depName})
						whisk.Debug(whisk.DbgInfo, output)
					}
				}
			}
		}
	}

	return nil
}

// TODO() display "update" | "synced" messages pre/post
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

	if err := deployer.RefreshManagedPackagesWithDependencies(ma); err != nil {
		return err
	}
	return nil

}

// TODO() display "update" | "synced" messages pre/post
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

				output := wski18n.T(wski18n.ID_MSG_MANAGED_FOUND_DELETED_X_key_X_name_X_project_X,
					map[string]interface{}{
						wski18n.KEY_KEY:     parsers.YAML_KEY_ACTION,
						wski18n.KEY_NAME:    actionName,
						wski18n.KEY_PROJECT: aa[utils.OW_PROJECT_NAME]})
				wskprint.PrintOpenWhiskWarning(output)

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

// TODO() display "update" | "synced" messages pre/post
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
				output := wski18n.T(wski18n.ID_MSG_MANAGED_FOUND_DELETED_X_key_X_name_X_project_X,
					map[string]interface{}{
						wski18n.KEY_KEY:     parsers.YAML_KEY_TRIGGER,
						wski18n.KEY_NAME:    trigger.Name,
						wski18n.KEY_PROJECT: ma[utils.OW_PROJECT_NAME]})
				wskprint.PrintOpenWhiskWarning(output)

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

// TODO() display "update" | "synced" messages pre/post
func (deployer *ServiceDeployer) RefreshManagedRules(ma map[string]interface{}) error {
	options := whisk.RuleListOptions{}
	// Get list of rules in your namespace
	rules, _, err := deployer.Client.Rules.List(&options)
	if err != nil {
		return err
	}
	// iterate over the list of rules to determine whether any of them was part of managed project
	// and now deleted from manifest file we can determine that from the managed annotation
	// If a rule has attached managed annotation with the project name equals to the current project name
	// but the project hash is different (project hash differs since the rule is deleted from the manifest file)
	for _, rule := range rules {
		// rule has attached managed annotation
		if a := rule.Annotations.GetValue(utils.MANAGED); a != nil {
			// decode the JSON blob and retrieve __OW_PROJECT_NAME and __OW_PROJECT_HASH
			ta := a.(map[string]interface{})
			if ta[utils.OW_PROJECT_NAME] == ma[utils.OW_PROJECT_NAME] && ta[utils.OW_PROJECT_HASH] != ma[utils.OW_PROJECT_HASH] {
				// we have found a trigger which was earlier part of the current project
				output := wski18n.T(wski18n.ID_MSG_MANAGED_FOUND_DELETED_X_key_X_name_X_project_X,
					map[string]interface{}{
						wski18n.KEY_KEY:     parsers.YAML_KEY_RULE,
						wski18n.KEY_NAME:    rule.Name,
						wski18n.KEY_PROJECT: ma[utils.OW_PROJECT_NAME]})
				wskprint.PrintOpenWhiskWarning(output)

				var err error
				err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
					_, err := deployer.Client.Rules.Delete(rule.Name)
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

// TODO() display "update" | "synced" messages pre/post
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
				output := wski18n.T(wski18n.ID_MSG_MANAGED_FOUND_DELETED_X_key_X_name_X_project_X,
					map[string]interface{}{
						wski18n.KEY_KEY:     parsers.YAML_KEY_PACKAGE,
						wski18n.KEY_NAME:    pkg.Name,
						wski18n.KEY_PROJECT: pa[utils.OW_PROJECT_NAME]})
				wskprint.PrintOpenWhiskWarning(output)

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

func (deployer *ServiceDeployer) appendDepAnnotation(list whisk.KeyValueArr, pkg *whisk.Package) whisk.KeyValueArr {
	depExists := false
	if a := pkg.Annotations.GetValue(utils.MANAGED); a != nil {
		//append annotations from this package to deps
		pkgName := parsers.PATH_SEPARATOR + pkg.Namespace + parsers.PATH_SEPARATOR + pkg.Name
		for _, dep := range list {
			if dep.Key == pkgName {
				depExists = true
			}
		}
		if !depExists {
			list = append(list, whisk.KeyValue{Key: pkgName, Value: a.(map[string]interface{})})
		}
	}
	return list
}

func (deployer *ServiceDeployer) RefreshManagedPackagesWithDependencies(ma map[string]interface{}) error {
	// iterate over each package from the given project
	for _, p := range deployer.Deployment.Packages {
		dependencyAnnotations := make(whisk.KeyValueArr, 0)
		// iterate over the list of dependencies of the package
		// dependencies could be labeled same as dependent package name
		// for example, "helloworld" where the package it depends on is also called "helloworld"
		// dependencies could be labeled different from the dependent package name
		// for example, "custom-helloworld" where the package it depends on it called "helloworld"
		for label, n := range p.Dependencies {
			// we do not append dependencies in whisk-managed in case of a package binding
			// for example when a package has dependencies on /whisk.system/ or any other /<namespace>/
			// the dependent packages are pre-installed and should not be managed by current project
			if !n.IsBinding {
				// find the package using dependency label
				pkg, _, err := deployer.Client.Packages.Get(label)
				if err != nil {
					return err
				}
				// if dependency label (custom-helloworld) is different from the dependent package name,
				// it must have binding set to the original package ("helloworld")
				if len(pkg.Binding.Name) != 0 {
					// having dependency on packages under /whisk.system is treated in a different way
					// in which dependent package under /whisk.system are not modified to add managed annotation
					// and parent package does not show this dependency in its managed annotation
					// because whisk.system packages comes pre-packaged and deployed with OpenWhisk server and not
					// deployed along with application deployments.
					// get the original package to retrieve its managed annotations
					pkg, _, err := deployer.Client.Packages.Get(pkg.Binding.Name)
					if err != nil {
						return err
					}
					dependencyAnnotations = deployer.appendDepAnnotation(dependencyAnnotations, pkg)
				} else {
					dependencyAnnotations = deployer.appendDepAnnotation(dependencyAnnotations, pkg)
				}
			}
		}
		updatedAnnotation, err := utils.AddDependentAnnotation(ma, dependencyAnnotations)
		if err != nil {
			return err
		}
		p.Package.Annotations.AddOrReplace(&updatedAnnotation)
	}
	if err := deployer.DeployPackages(); err != nil {
		return err
	}
	return nil
}

func (deployer *ServiceDeployer) DeployPackages() error {
	for _, pack := range deployer.Deployment.Packages {
		// "default" package is a reserved package name
		// all openwhisk entities will be deployed under
		// /<namespace> instead of /<namespace>/<package> and
		// therefore skip creating a new package and set
		// deployer.DeployActionInPackage to false which is set to true by default
		if strings.ToLower(pack.Package.Name) != parsers.DEFAULT_PACKAGE {
			err := deployer.createPackage(pack.Package)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Deploy Sequences into OpenWhisk
func (deployer *ServiceDeployer) DeploySequences() error {

	for _, pack := range deployer.Deployment.Packages {
		for _, action := range pack.Sequences {
			err := deployer.createAction(pack.Package.Name, action.Action)
			if err != nil {
				return err
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
			err := deployer.createFeedAction(trigger, feedname)
			if err != nil {
				return err
			}
		} else {
			err := deployer.createTrigger(trigger)
			if err != nil {
				return err
			}
		}

	}
	return nil

}

// Deploy Rules into OpenWhisk
func (deployer *ServiceDeployer) DeployRules() error {
	for _, rule := range deployer.Deployment.Rules {
		err := deployer.createRule(rule)
		if err != nil {
			return err
		}
	}
	return nil
}

// Deploy Apis into OpenWhisk
func (deployer *ServiceDeployer) DeployApis() error {
	var err error
	// NOTE: Only deploy either swagger or manifest defined api, but not both
	// NOTE: Swagger API takes precedence
	if deployer.Deployment.SwaggerApi != nil && deployer.Deployment.SwaggerApiOptions != nil {
		err = deployer.createSwaggerApi(deployer.Deployment.SwaggerApi)
		if err != nil {
			return err
		}
	} else {
		for _, api := range deployer.Deployment.Apis {
			err = deployer.createApi(api)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (deployer *ServiceDeployer) createBinding(packa *whisk.BindingPackage) error {

	displayPreprocessingInfo(wski18n.PACKAGE_BINDING, packa.Name, true)

	var err error
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Packages.Insert(packa, true)
		return err
	})

	if err != nil {
		return createWhiskClientError(err.(*whisk.WskError), response, wski18n.PACKAGE_BINDING, true)
	}

	displayPostprocessingInfo(wski18n.PACKAGE_BINDING, packa.Name, true)
	return nil
}

func (deployer *ServiceDeployer) createPackage(packa *whisk.Package) error {

	displayPreprocessingInfo(parsers.YAML_KEY_PACKAGE, packa.Name, true)

	var err error
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Packages.Insert(packa, true)
		return err
	})
	if err != nil {
		return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_PACKAGE, true)
	}

	displayPostprocessingInfo(parsers.YAML_KEY_PACKAGE, packa.Name, true)
	return nil
}

func (deployer *ServiceDeployer) createTrigger(trigger *whisk.Trigger) error {

	displayPreprocessingInfo(parsers.YAML_KEY_TRIGGER, trigger.Name, true)

	var err error
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Triggers.Insert(trigger, true)
		return err
	})
	if err != nil {
		return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_TRIGGER, true)
	}

	displayPostprocessingInfo(parsers.YAML_KEY_TRIGGER, trigger.Name, true)
	return nil
}

func (deployer *ServiceDeployer) createFeedAction(trigger *whisk.Trigger, feedName string) error {

	displayPreprocessingInfo(wski18n.TRIGGER_FEED, trigger.Name, true)

	// to hold and modify trigger parameters, not passed by ref?
	params := make(map[string]interface{})

	// check for strings that are JSON
	for _, keyVal := range trigger.Parameters {
		params[keyVal.Key] = keyVal.Value
	}

	// TODO() define keys and lifecycle operation names as const
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
	if err = deployer.createTrigger(t); err != nil {
		return err
	}
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

		return createWhiskClientError(err.(*whisk.WskError), response, wski18n.TRIGGER_FEED, false)
	}

	displayPostprocessingInfo(wski18n.TRIGGER_FEED, trigger.Name, true)
	return nil
}

func (deployer *ServiceDeployer) createRule(rule *whisk.Rule) error {
	displayPreprocessingInfo(parsers.YAML_KEY_RULE, rule.Name, true)

	// The rule's trigger should include the namespace with pattern /namespace/trigger
	rule.Trigger = deployer.getQualifiedName(rule.Trigger.(string))
	// The rule's action should include the namespace and package with pattern
	// /namespace/package/action if that action was created under a package
	// otherwise action should include the namespace with pattern /namespace/action
	rule.Action = deployer.getQualifiedName(rule.Action.(string))

	var err error
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Rules.Insert(rule, true)
		return err
	})

	if err != nil {
		return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_RULE, true)
	}

	// Consecutive deployments of manifest containing trigger with feed action (and rule) result in inactive
	// rule. The rule seems to become inactive when its trigger get deleted (part of the wskdeploy feed action update)
	// Currently simply always setting rule status to active in case not specified implicitly
	_, _, err = deployer.Client.Rules.SetState(rule.Name, "active")
	if err != nil {
		return err
	}

	displayPostprocessingInfo(parsers.YAML_KEY_RULE, rule.Name, true)
	return nil
}

// Utility function to call go-whisk framework to make action
func (deployer *ServiceDeployer) createAction(pkgname string, action *whisk.Action) error {
	// call ActionService through the Client
	if strings.ToLower(pkgname) != parsers.DEFAULT_PACKAGE {
		// the action will be created under package with pattern 'packagename/actionname'
		action.Name = strings.Join([]string{pkgname, action.Name}, "/")
	}

	displayPreprocessingInfo(parsers.YAML_KEY_ACTION, action.Name, true)

	var err error
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Actions.Insert(action, true)
		return err
	})

	if err != nil {
		return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_ACTION, true)
	}

	displayPostprocessingInfo(parsers.YAML_KEY_ACTION, action.Name, true)
	return nil
}

func (deployer *ServiceDeployer) getAnnotationsFromPackageActionOrSequence(packageActionName string) *whisk.KeyValueArr {

	if len(packageActionName) != 0 {
		// Split the package name and action name being searched for
		aActionName := strings.Split(packageActionName, parsers.PATH_SEPARATOR)

		// Attempt to locate the named action (or sequence) to return its annotations
		if pkg, found := deployer.Deployment.Packages[aActionName[0]]; found {
			if atemp, found := pkg.Actions[aActionName[1]]; found {
				return &(atemp.Action.Annotations)
			} else if atemp, found := pkg.Sequences[aActionName[1]]; found {
				return &(atemp.Action.Annotations)
			}
		}
	}
	return nil
}

// create api (API Gateway functionality)
func (deployer *ServiceDeployer) createApi(api *whisk.ApiCreateRequest) error {

	apiPath := api.ApiDoc.ApiName + " " + api.ApiDoc.GatewayBasePath +
		api.ApiDoc.GatewayRelPath + " " + api.ApiDoc.GatewayMethod
	displayPreprocessingInfo(parsers.YAML_KEY_API, apiPath, true)

	var err error
	var response *http.Response

	apiCreateReqOptions := deployer.Deployment.ApiOptions[apiPath]

	// Retrieve annotations on the action we are attempting to create an API for
	var actionAnnotations *whisk.KeyValueArr
	actionAnnotations = deployer.getAnnotationsFromPackageActionOrSequence(api.ApiDoc.Action.Name)

	// Process any special annotations (e.g., "require-whisk-auth") on the associated Action
	// NOTE: we do not throw an error if annotations are NOT found (nil) since this is already done in
	// the parsing phase and would be redundant.
	if actionAnnotations != nil {
		wskprint.PrintlnOpenWhiskVerbose(utils.Flags.Verbose, fmt.Sprintf("Processing action annotations: %v", actionAnnotations))

		// If the "require-whisk-auth" annotation is present on the referenced action,
		// apply its user provided security key (i.e., the annotation's value) to the API
		if webaction.HasAnnotation(actionAnnotations, webaction.REQUIRE_WHISK_AUTH) {
			api.ApiDoc.Action.SecureKey = actionAnnotations.GetValue(webaction.REQUIRE_WHISK_AUTH)
		}
	}

	if len(deployer.Client.Config.ApigwTenantId) > 0 {
		// Use it to identify the IAM namespace
		apiCreateReqOptions.SpaceGuid = deployer.Client.Config.ApigwTenantId
	} else {
		//  assume a CF namespace (SpaceGuid) which is part of the authtoken
		apiCreateReqOptions.SpaceGuid = strings.Split(deployer.Client.Config.AuthToken, ":")[0]
	}

	apiCreateReqOptions.AccessToken = deployer.Client.Config.ApigwAccessToken

	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Apis.Insert(api, apiCreateReqOptions, true)
		return err
	})

	if err != nil {
		return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_API, true)
	}

	displayPostprocessingInfo(parsers.YAML_KEY_API, apiPath, true)
	return nil
}

// create api (API Gateway functionality) from swagger file
func (deployer *ServiceDeployer) createSwaggerApi(api *whisk.ApiCreateRequest) error {
	var err error
	var response *http.Response

	apiCreateReqOptions := deployer.Deployment.SwaggerApiOptions
	apiCreateReqOptions.AccessToken = deployer.Client.Config.ApigwAccessToken
	// In the case of IAM namespaces, we must use the ApigwTenantId as the SpaceGuid
	// IAM namespaces can be detected by seeing if the ApigwTenantId is populated
	if len(deployer.Client.Config.ApigwTenantId) > 0 {
		apiCreateReqOptions.SpaceGuid = deployer.Client.Config.ApigwTenantId
	} else {
		apiCreateReqOptions.SpaceGuid = strings.Split(deployer.Client.Config.AuthToken, ":")[0]
	}

	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		_, response, err = deployer.Client.Apis.Insert(api, apiCreateReqOptions, true)
		return err
	})

	if err != nil {
		return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_API, true)
	}

	return nil
}

func (deployer *ServiceDeployer) UnDeploy(verifiedPlan *DeploymentProject) error {
	if deployer.Preview == true {
		deployer.printDeploymentAssets(verifiedPlan)
		return nil
	}

	if err := deployer.unDeployAssets(verifiedPlan); err != nil {
		wskprint.PrintOpenWhiskError(wski18n.T(wski18n.T(wski18n.ID_MSG_UNDEPLOYMENT_FAILED)))
		return err
	}

	wskprint.PrintOpenWhiskSuccess(wski18n.T(wski18n.T(wski18n.ID_MSG_UNDEPLOYMENT_SUCCEEDED)))
	return nil
}

func (deployer *ServiceDeployer) UnDeployProject() error {
	if err := deployer.UnDeployProjectAssets(); err != nil {
		return err
	}
	return nil
}

func (deployer *ServiceDeployer) unDeployAssets(verifiedPlan *DeploymentProject) error {
	if err := deployer.UndeploySwaggerApis(verifiedPlan); err != nil {
		return err
	}
	if err := deployer.UnDeployApis(verifiedPlan); err != nil {
		return err
	}

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
				map[string]interface{}{wski18n.KEY_NAME: depName})
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

				dependentPackages := []string{}
				for k := range depServiceDeployer.Deployment.Packages {
					dependentPackages = append(dependentPackages, k)
				}

				// delete binding pkg if the origin package name is different
				if ok := depServiceDeployer.Deployment.Packages[depName]; ok == nil {
					if _, _, ok := deployer.Client.Packages.Get(depName); ok == nil {
						var err error
						var response *http.Response
						err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
							response, err = deployer.Client.Packages.Delete(depName)
							return err
						})
						if err != nil {
							return createWhiskClientError(err.(*whisk.WskError), response, wski18n.PACKAGE_BINDING, false)
						}
					}
				}

				if err := depServiceDeployer.unDeployAssets(plan); err != nil {
					errString := wski18n.T(wski18n.ID_MSG_DEPENDENCY_UNDEPLOYMENT_FAILURE_X_name_X,
						map[string]interface{}{wski18n.KEY_NAME: depName})
					whisk.Debug(whisk.DbgError, errString)
					return err
				}
			}
			output = wski18n.T(wski18n.ID_MSG_DEPENDENCY_UNDEPLOYMENT_SUCCESS_X_name_X,
				map[string]interface{}{wski18n.KEY_NAME: depName})
			whisk.Debug(whisk.DbgInfo, output)
		}
	}

	return nil
}

func (deployer *ServiceDeployer) UnDeployPackages(deployment *DeploymentProject) error {
	for _, pack := range deployment.Packages {
		// "default" package is a reserved package name
		// all openwhisk entities were deployed under
		// /<namespace> instead of /<namespace>/<package> and
		// therefore skip deleting default package during undeployment
		if strings.ToLower(pack.Package.Name) != parsers.DEFAULT_PACKAGE {
			err := deployer.deletePackage(pack.Package)
			if err != nil {
				return err
			}
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
		}
		err := deployer.deleteTrigger(trigger)
		if err != nil {
			return err
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

func (deployer *ServiceDeployer) UnDeployApis(deployment *DeploymentProject) error {

	for _, api := range deployment.Apis {
		err := deployer.deleteApi(api)
		if err != nil {
			return err
		}
	}
	return nil
}

func (deployer *ServiceDeployer) UndeploySwaggerApis(deployment *DeploymentProject) error {
	api := deployment.SwaggerApi
	err := deployer.deleteSwaggerApi(api)
	if err != nil {
		return err
	}
	return nil
}

func (deployer *ServiceDeployer) deletePackage(packa *whisk.Package) error {

	displayPreprocessingInfo(parsers.YAML_KEY_PACKAGE, packa.Name, false)

	if _, _, ok := deployer.Client.Packages.Get(packa.Name); ok == nil {
		var err error
		var response *http.Response
		err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
			response, err = deployer.Client.Packages.Delete(packa.Name)
			return err
		})

		if err != nil {
			return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_PACKAGE, false)
		}
	}
	displayPostprocessingInfo(parsers.YAML_KEY_PACKAGE, packa.Name, false)
	return nil
}

func (deployer *ServiceDeployer) deleteTrigger(trigger *whisk.Trigger) error {

	displayPreprocessingInfo(parsers.YAML_KEY_TRIGGER, trigger.Name, false)

	if _, _, ok := deployer.Client.Triggers.Get(trigger.Name); ok == nil {
		var err error
		var response *http.Response
		err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
			_, response, err = deployer.Client.Triggers.Delete(trigger.Name)
			return err
		})

		if err != nil {
			return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_TRIGGER, false)
		}
	}

	displayPostprocessingInfo(parsers.YAML_KEY_TRIGGER, trigger.Name, false)
	return nil
}

func (deployer *ServiceDeployer) deleteFeedAction(trigger *whisk.Trigger, feedName string) error {

	displayPreprocessingInfo(parsers.YAML_KEY_FEED, trigger.Name, false)

	params := make(whisk.KeyValueArr, 0)
	// TODO() define keys and operations as const
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

	if _, _, ok := deployer.Client.Triggers.Get(trigger.Name); ok != nil {
		displayPostprocessingInfo(parsers.YAML_KEY_FEED, trigger.Name, false)
		return nil
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
		errString := wski18n.T(wski18n.ID_ERR_FEED_INVOKE_X_err_X_code_X,
			map[string]interface{}{wski18n.KEY_ERR: wskErr.Error(), wski18n.KEY_CODE: strconv.Itoa(wskErr.ExitCode)})
		whisk.Debug(whisk.DbgError, errString)
		return wskderrors.NewWhiskClientError(wskErr.Error(), wskErr.ExitCode, response)

	}
	displayPostprocessingInfo(parsers.YAML_KEY_FEED, trigger.Name, false)
	return nil
}

func (deployer *ServiceDeployer) deleteRule(rule *whisk.Rule) error {

	displayPreprocessingInfo(parsers.YAML_KEY_RULE, rule.Name, false)

	if _, _, ok := deployer.Client.Rules.Get(rule.Name); ok == nil {
		var err error
		var response *http.Response
		err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
			response, err = deployer.Client.Rules.Delete(rule.Name)
			return err
		})

		if err != nil {
			return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_RULE, false)
		}
	}
	displayPostprocessingInfo(parsers.YAML_KEY_RULE, rule.Name, false)
	return nil
}

func (deployer *ServiceDeployer) isApi(api *whisk.ApiCreateRequest) bool {
	apiReqOptions := new(whisk.ApiGetRequestOptions)
	apiReqOptions.AccessToken = deployer.Client.Config.ApigwAccessToken
	apiReqOptions.ApiBasePath = api.ApiDoc.GatewayBasePath
	if len(deployer.Client.Config.ApigwTenantId) > 0 {
		// Use it to identify the IAM namespace
		apiReqOptions.SpaceGuid = deployer.Client.Config.ApigwTenantId
	} else {
		//  assume a CF namespaces (SpaceGuid) which is part of the authtoken
		apiReqOptions.SpaceGuid = strings.Split(deployer.Client.Config.AuthToken, ":")[0]
	}

	a := new(whisk.ApiGetRequest)

	retApi, _, err := deployer.Client.Apis.Get(a, apiReqOptions)
	if err == nil {
		if retApi.Apis != nil && len(retApi.Apis) > 0 &&
			retApi.Apis[0].ApiValue != nil {
			return true
		}

	}
	return false
}

// delete api (API Gateway functionality)
func (deployer *ServiceDeployer) deleteApi(api *whisk.ApiCreateRequest) error {

	apiPath := api.ApiDoc.ApiName + " " + api.ApiDoc.GatewayBasePath +
		api.ApiDoc.GatewayRelPath + " " + api.ApiDoc.GatewayMethod
	displayPreprocessingInfo(parsers.YAML_KEY_API, apiPath, false)

	if deployer.isApi(api) {
		var err error
		var response *http.Response

		apiDeleteReqOptions := new(whisk.ApiDeleteRequestOptions)
		apiDeleteReqOptions.AccessToken = deployer.Client.Config.ApigwAccessToken
		if len(deployer.Client.Config.ApigwTenantId) > 0 {
			// Use it to identify the IAM namespace
			apiDeleteReqOptions.SpaceGuid = deployer.Client.Config.ApigwTenantId
		} else {
			//  assume a CF namespaces (SpaceGuid) which is part of the authtoken
			apiDeleteReqOptions.SpaceGuid = strings.Split(deployer.Client.Config.AuthToken, ":")[0]
		}
		apiDeleteReqOptions.ApiBasePath = api.ApiDoc.GatewayBasePath
		apiDeleteReqOptions.ApiRelPath = api.ApiDoc.GatewayRelPath
		apiDeleteReqOptions.ApiVerb = api.ApiDoc.GatewayMethod

		a := new(whisk.ApiDeleteRequest)

		err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
			response, err = deployer.Client.Apis.Delete(a, apiDeleteReqOptions)
			return err
		})

		if err != nil {
			return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_API, false)
		}
	}
	displayPostprocessingInfo(parsers.YAML_KEY_API, apiPath, false)
	return nil
}

// delete api (API Gateway functionality) from swagger file
func (deployer *ServiceDeployer) deleteSwaggerApi(api *whisk.ApiCreateRequest) error {
	var err error
	var response *http.Response

	// If there is no swagger file do nothing
	if api == nil {
		return nil
	}
	swaggerString := api.ApiDoc.Swagger
	swaggerObj := new(whisk.ApiSwagger)
	err = json.Unmarshal([]byte(swaggerString), swaggerObj)
	if err != nil {
		return err
	}

	apiDeleteReqOptions := new(whisk.ApiDeleteRequestOptions)
	apiDeleteReqOptions.AccessToken = deployer.Client.Config.ApigwAccessToken
	apiDeleteReqOptions.ApiBasePath = swaggerObj.BasePath
	// In the case of IAM namespaces, we must use the ApigwTenantId as the SpaceGuid
	// IAM namespaces can be detected by seeing if the ApigwTenantId is populated
	if len(deployer.Client.Config.ApigwTenantId) > 0 {
		apiDeleteReqOptions.SpaceGuid = deployer.Client.Config.ApigwTenantId
	} else {
		apiDeleteReqOptions.SpaceGuid = strings.Split(deployer.Client.Config.AuthToken, ":")[0]
	}

	a := new(whisk.ApiDeleteRequest)
	a.Swagger = swaggerString

	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		response, err = deployer.Client.Apis.Delete(a, apiDeleteReqOptions)
		return err
	})

	if err != nil {
		return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_API, true)
	}

	return nil
}

// Utility function to call go-whisk framework to delete action
func (deployer *ServiceDeployer) deleteAction(pkgname string, action *whisk.Action) error {
	// call ActionService through Client
	if pkgname != parsers.DEFAULT_PACKAGE {
		// the action will be deleted under package with pattern 'packagename/actionname'
		action.Name = strings.Join([]string{pkgname, action.Name}, "/")
	}

	displayPreprocessingInfo(parsers.YAML_KEY_ACTION, action.Name, false)

	if _, _, ok := deployer.Client.Actions.Get(action.Name, false); ok == nil {
		var err error
		var response *http.Response
		err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
			response, err = deployer.Client.Actions.Delete(action.Name)
			return err
		})

		if err != nil {
			return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_ACTION, false)

		}
	}
	displayPostprocessingInfo(parsers.YAML_KEY_ACTION, action.Name, false)
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
				warningMsg := wski18n.T(wski18n.ID_WARN_COMMAND_RETRY,
					map[string]interface{}{
						wski18n.KEY_CMD: strconv.Itoa(i + 1),
						wski18n.KEY_ERR: err.Error()})
				wskprint.PrintlnOpenWhiskWarning(warningMsg)
			} else {
				return err
			}
		} else {
			return err
		}
	}
	return err
}

//  getQualifiedName(name) returns a fully qualified name given a
//      (possibly fully qualified) resource name.
//
//  Examples:
//      (foo) => /ns/foo
//      (pkg/foo) => /ns/pkg/foo
//      (/ns/pkg/foo) => /ns/pkg/foo
func (deployer *ServiceDeployer) getQualifiedName(name string) string {
	namespace := deployer.ClientConfig.Namespace
	if strings.HasPrefix(name, "/") {
		return name
	} else if strings.HasPrefix(namespace, "/") {
		return fmt.Sprintf("%s/%s", namespace, name)
	}
	return fmt.Sprintf("/%s/%s", namespace, name)
}

func (deployer *ServiceDeployer) printDeploymentAssets(assets *DeploymentProject) {

	// TODO() review format
	wskprint.PrintlnOpenWhiskOutput(strings.Title(parsers.YAML_KEY_PACKAGES) + ":")
	for _, pack := range assets.Packages {
		wskprint.PrintlnOpenWhiskOutput(strings.Title(wski18n.KEY_NAME) + ": " + pack.Package.Name)
		wskprint.PrintlnOpenWhiskOutput("    " + wski18n.KEY_BINDINGS + ": ")
		for _, p := range pack.Package.Parameters {
			jsonValue, err := utils.PrettyJSON(p.Value)
			if err != nil {
				fmt.Printf("        - %s : %s\n", p.Key, wskderrors.STR_UNKNOWN_VALUE)
			} else {
				fmt.Printf("        - %s : %v\n", p.Key, jsonValue)
			}
		}

		wskprint.PrintlnOpenWhiskOutput("    " + parsers.YAML_KEY_ANNOTATION + ": ")
		for _, p := range pack.Package.Annotations {
			fmt.Printf("        - %s : %v\n", p.Key, p.Value)

		}

		for key, dep := range pack.Dependencies {
			wskprint.PrintlnOpenWhiskOutput("  * " + wski18n.KEY_DEPENDENCY + ": " + key)
			wskprint.PrintlnOpenWhiskOutput("    " + wski18n.KEY_LOCATION + ": " + dep.Location)
			if !dep.IsBinding {
				wskprint.PrintlnOpenWhiskOutput("    " + wski18n.KEY_PATH + ": " + dep.ProjectPath)
			}
		}

		wskprint.PrintlnOpenWhiskOutput("")

		for _, action := range pack.Actions {
			wskprint.PrintlnOpenWhiskOutput("  * " + parsers.YAML_KEY_ACTION + ": " + action.Action.Name)
			wskprint.PrintlnOpenWhiskOutput("    " + wski18n.KEY_BINDINGS + ": ")
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
			wskprint.PrintlnOpenWhiskOutput("    " + parsers.YAML_KEY_ANNOTATION + ": ")
			for _, p := range action.Action.Annotations {
				fmt.Printf("        - %s : %v\n", p.Key, p.Value)

			}
		}

		wskprint.PrintlnOpenWhiskOutput("")
		for _, action := range pack.Sequences {
			wskprint.PrintlnOpenWhiskOutput("  * " + parsers.YAML_KEY_SEQUENCE + ": " + action.Action.Name)
			wskprint.PrintlnOpenWhiskOutput("    " + parsers.YAML_KEY_ANNOTATION + ": ")
			for _, p := range action.Action.Annotations {
				fmt.Printf("        - %s : %v\n", p.Key, p.Value)

			}
		}

		wskprint.PrintlnOpenWhiskOutput("")
	}

	wskprint.PrintlnOpenWhiskOutput(wski18n.TRIGGERS + ":")
	for _, trigger := range assets.Triggers {
		wskprint.PrintlnOpenWhiskOutput("* " + parsers.YAML_KEY_TRIGGER + ": " + trigger.Name)
		wskprint.PrintlnOpenWhiskOutput("    " + wski18n.KEY_BINDINGS + ": ")

		for _, p := range trigger.Parameters {
			jsonValue, err := utils.PrettyJSON(p.Value)
			if err != nil {
				fmt.Printf("        - %s : %s\n", p.Key, wskderrors.STR_UNKNOWN_VALUE)
			} else {
				fmt.Printf("        - %s : %v\n", p.Key, jsonValue)
			}
		}

		wskprint.PrintlnOpenWhiskOutput("    " + parsers.YAML_KEY_ANNOTATION + ": ")
		for _, p := range trigger.Annotations {
			fmt.Printf("        - %s : %v\n", p.Key, p.Value)

		}
	}

	wskprint.PrintlnOpenWhiskOutput("\n" + wski18n.RULES)
	for _, rule := range assets.Rules {
		wskprint.PrintlnOpenWhiskOutput("* " + parsers.YAML_KEY_RULE + ": " + rule.Name)
		wskprint.PrintlnOpenWhiskOutput("    " + parsers.YAML_KEY_ANNOTATION + ": ")
		for _, p := range rule.Annotations {
			fmt.Printf("        - %s : %v\n", p.Key, p.Value)

		}
		if reflect.TypeOf(rule.Trigger).Kind() == reflect.String {
			wskprint.PrintlnOpenWhiskOutput("    - " + parsers.YAML_KEY_TRIGGER + ": " + rule.Trigger.(string) + "\n    - " + parsers.YAML_KEY_ACTION + ": " + rule.Action.(string))
		} else if reflect.TypeOf(rule.Trigger).Kind() == reflect.Map {
			trigger := rule.Trigger.(map[string]interface{})
			triggerName := trigger["path"].(string) + parsers.PATH_SEPARATOR + trigger["name"].(string)
			action := rule.Action.(map[string]interface{})
			actionName := action["path"].(string) + parsers.PATH_SEPARATOR + action["name"].(string)
			wskprint.PrintlnOpenWhiskOutput("    - " + parsers.YAML_KEY_TRIGGER + ": " + triggerName + "\n    - " + parsers.YAML_KEY_ACTION + ": " + actionName)
		}

	}

	wskprint.PrintlnOpenWhiskOutput("")

}

func (deployer *ServiceDeployer) getDependentDeployer(depName string, depRecord dependencies.DependencyRecord) (*ServiceDeployer, error) {
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
	depServiceDeployer.Preview = true

	depServiceDeployer.Client = deployer.Client
	depServiceDeployer.ClientConfig = deployer.ClientConfig

	// share the master dependency list
	depServiceDeployer.DependencyMaster = deployer.DependencyMaster

	return depServiceDeployer, nil
}

func displayPreprocessingInfo(entity string, name string, onDeploy bool) {

	var msgKey string
	if onDeploy {
		msgKey = wski18n.ID_MSG_ENTITY_DEPLOYING_X_key_X_name_X
	} else {
		msgKey = wski18n.ID_MSG_ENTITY_UNDEPLOYING_X_key_X_name_X
	}
	msg := wski18n.T(msgKey,
		map[string]interface{}{
			wski18n.KEY_KEY:  entity,
			wski18n.KEY_NAME: name})
	wskprint.PrintlnOpenWhiskVerbose(utils.Flags.Verbose, msg)
}

func displayPostprocessingInfo(entity string, name string, onDeploy bool) {

	var msgKey string
	if onDeploy {
		msgKey = wski18n.ID_MSG_ENTITY_DEPLOYED_SUCCESS_X_key_X_name_X
	} else {
		msgKey = wski18n.ID_MSG_ENTITY_UNDEPLOYED_SUCCESS_X_key_X_name_X
	}
	msg := wski18n.T(msgKey,
		map[string]interface{}{
			wski18n.KEY_KEY:  entity,
			wski18n.KEY_NAME: name})
	wskprint.PrintlnOpenWhiskVerbose(utils.Flags.Verbose, msg)
}

func createWhiskClientError(err *whisk.WskError, response *http.Response, entity string, onCreate bool) *wskderrors.WhiskClientError {

	var msgKey string
	if onCreate {
		msgKey = wski18n.ID_ERR_ENTITY_CREATE_X_key_X_err_X_code_X
	} else {
		msgKey = wski18n.ID_ERR_ENTITY_DELETE_X_key_X_err_X_code_X
	}
	errString := wski18n.T(msgKey,
		map[string]interface{}{
			wski18n.KEY_KEY:  entity,
			wski18n.KEY_ERR:  err.Error(),
			wski18n.KEY_CODE: strconv.Itoa(err.ExitCode)})
	wskprint.PrintOpenWhiskVerbose(utils.Flags.Verbose, errString)

	// TODO() add errString as an AppendDetail() to WhiskClientError
	return wskderrors.NewWhiskClientError(err.Error(), err.ExitCode, response)
}

func (deployer *ServiceDeployer) reportInputs() error {
	// display project level inputs
	i := make(map[string]interface{}, 0)
	for name, param := range deployer.ProjectInputs {
		i[name] = param.Value
	}
	projectInputs := parsers.DisplayInputs{Name: deployer.ProjectName, Inputs: i}
	j, err := json.MarshalIndent(projectInputs, "", " ")
	if err != nil {
		return err
	}
	wskprint.PrintlnOpenWhiskOutput(string(j))

	// display package level inputs
	// iterate over each package and print inputs section of each package
	for _, pkg := range deployer.Deployment.Packages {
		i := make(map[string]interface{}, 0)
		for name, param := range pkg.Inputs.Inputs {
			if _, ok := deployer.ProjectInputs[name]; !ok {
				i[name] = param.Value
			}
		}
		packageInputs := parsers.DisplayInputs{Name: pkg.Package.Name, Inputs: i}
		j, err := json.MarshalIndent(packageInputs, "", "  ")
		if err != nil {
			return err
		}
		wskprint.PrintlnOpenWhiskOutput(string(j))

		for _, d := range pkg.Dependencies {
			i := make(map[string]interface{}, 0)
			for _, param := range d.Parameters {
				i[param.Key] = param.Value
			}
			depInputs := parsers.DisplayInputs{Name: d.Location, Inputs: i}
			j, err := json.MarshalIndent(depInputs, "", " ")
			if err != nil {
				return err
			}
			wskprint.PrintlnOpenWhiskOutput(string(j))
		}

		for _, a := range pkg.Actions {
			i := make(map[string]interface{}, 0)
			for _, param := range a.Action.Parameters {
				i[param.Key] = param.Value
			}

			actionInputs := parsers.DisplayInputs{Name: a.Action.Name, Inputs: i}
			j, err := json.MarshalIndent(actionInputs, "", " ")
			if err != nil {
				return err
			}
			wskprint.PrintlnOpenWhiskOutput(string(j))
		}

		for _, s := range pkg.Sequences {
			i := make(map[string]interface{}, 0)
			for _, param := range s.Action.Parameters {
				i[param.Key] = param.Value
			}
			seqInputs := parsers.DisplayInputs{Name: s.Action.Name, Inputs: i}
			j, err := json.MarshalIndent(seqInputs, "", " ")
			if err != nil {
				return err
			}
			wskprint.PrintlnOpenWhiskOutput(string(j))
		}
	}

	for _, trigger := range deployer.Deployment.Triggers {
		i := make(map[string]interface{}, 0)
		for _, param := range trigger.Parameters {
			i[param.Key] = param.Value
		}
		triggerInputs := parsers.DisplayInputs{Name: trigger.Name, Inputs: i}
		j, err := json.MarshalIndent(triggerInputs, "", " ")
		if err != nil {
			return err
		}
		wskprint.PrintlnOpenWhiskOutput(string(j))
	}
	return nil
}
