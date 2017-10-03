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
	"errors"
	"fmt"
	"strings"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
)

var clientConfig *whisk.Config

type ManifestReader struct {
	serviceDeployer *ServiceDeployer
	IsUndeploy      bool
}

func NewManfiestReader(serviceDeployer *ServiceDeployer) *ManifestReader {
	var dep ManifestReader
	dep.serviceDeployer = serviceDeployer

	return &dep
}

func (deployer *ManifestReader) ParseManifest() (*parsers.ManifestYAML, *parsers.YAMLParser, error) {
	dep := deployer.serviceDeployer
	manifestParser := parsers.NewYAMLParser()
	manifest, err := manifestParser.ParseManifest(dep.ManifestPath)

	if err != nil {
		return manifest, manifestParser, utils.NewInputYamlFileError(err.Error())
	}
	return manifest, manifestParser, nil
}

func (reader *ManifestReader) InitRootPackage(manifestParser *parsers.YAMLParser, manifest *parsers.ManifestYAML) error {
	packages, err := manifestParser.ComposeAllPackages(manifest, reader.serviceDeployer.ManifestPath)
	if err != nil {
		return utils.NewInputYamlFormatError(err.Error())
	}
	reader.SetPackage(packages)

	return nil
}

// Wrapper parser to handle yaml dir
func (deployer *ManifestReader) HandleYaml(sdeployer *ServiceDeployer, manifestParser *parsers.YAMLParser, manifest *parsers.ManifestYAML) error {

	var err error
	deps, err := manifestParser.ComposeDependenciesFromAllPackages(manifest, deployer.serviceDeployer.ProjectPath, deployer.serviceDeployer.ManifestPath)
	if err != nil {
		return utils.NewInputYamlFormatError(err.Error())
	}

	actions, err := manifestParser.ComposeActionsFromAllPackages(manifest, deployer.serviceDeployer.ManifestPath)
	if err != nil {
		return utils.NewInputYamlFormatError(err.Error())
	}

	sequences, err := manifestParser.ComposeSequencesFromAllPackages(deployer.serviceDeployer.ClientConfig.Namespace, manifest)
	if err != nil {
		return utils.NewInputYamlFormatError(err.Error())
	}

	triggers, err := manifestParser.ComposeTriggersFromAllPackages(manifest, deployer.serviceDeployer.ManifestPath)
	if err != nil {
		return utils.NewInputYamlFormatError(err.Error())
	}

	rules, err := manifestParser.ComposeRulesFromAllPackages(manifest)
	if err != nil {
		return utils.NewInputYamlFormatError(err.Error())
	}

	apis, err := manifestParser.ComposeApiRecordsFromAllPackages(manifest)
	if err != nil {
		return utils.NewInputYamlFormatError(err.Error())
	}

	err = deployer.SetDependencies(deps)
	if err != nil {
		return utils.NewInputYamlFormatError(err.Error())
	}

	err = deployer.SetActions(actions)
	if err != nil {
		return utils.NewInputYamlFormatError(err.Error())
	}

	err = deployer.SetSequences(sequences)
	if err != nil {
		return utils.NewInputYamlFormatError(err.Error())
	}

	err = deployer.SetTriggers(triggers)
	if err != nil {
		return utils.NewInputYamlFormatError(err.Error())
	}

	err = deployer.SetRules(rules)
	if err != nil {
		return utils.NewInputYamlFormatError(err.Error())
	}

	err = deployer.SetApis(apis)
	if err != nil {
		return utils.NewInputYamlFormatError(err.Error())
	}

	return nil

}

func (reader *ManifestReader) SetDependencies(deps map[string]utils.DependencyRecord) error {

	d := reader.serviceDeployer

	for name, dep := range deps {
		d.mt.Lock()
		n := strings.Split(name, ":")
		depName := n[1]
		if (depName == "") {
			return nil
		}
		if !dep.IsBinding && !reader.IsUndeploy {
			if _, exists := reader.serviceDeployer.DependencyMaster[depName]; !exists {
				// dependency
				gitReader := utils.NewGitReader(depName, dep)
				err := gitReader.CloneDependency()
				if err != nil {
					return utils.NewInputYamlFormatError(err.Error())
				}
			} else {
				// TODO: we should do a check to make sure this dependency is compatible with an already installed one.
				// If not, we should throw dependency mismatch error.
			}
		}

		// store in two places (one local to package to preserve relationship, one in master record to check for conflics
		reader.serviceDeployer.Deployment.Packages[dep.Packagename].Dependencies[depName] = dep
		reader.serviceDeployer.DependencyMaster[depName] = dep
		d.mt.Unlock()
	}

	return nil
}

func (reader *ManifestReader) SetPackage(packages map[string]*whisk.Package) error {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, pkg := range packages {
		depPkg, exist := dep.Deployment.Packages[pkg.Name]
		if exist {
			if dep.IsDefault == true {
				existPkg := depPkg.Package
				existPkg.Annotations = pkg.Annotations
				existPkg.Namespace = pkg.Namespace
				existPkg.Parameters = pkg.Parameters
				existPkg.Publish = pkg.Publish
				existPkg.Version = pkg.Version

				dep.Deployment.Packages[pkg.Name].Package = existPkg
				return nil
			} else {
				return errors.New("Package " + pkg.Name + "exists twice")
			}
		}
		newPack := NewDeploymentPackage()
		newPack.Package = pkg
		dep.Deployment.Packages[pkg.Name] = newPack
	}
	return nil
}

func (reader *ManifestReader) SetActions(actions []utils.ActionRecord) error {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, manifestAction := range actions {
		existAction, exists := reader.serviceDeployer.Deployment.Packages[manifestAction.Packagename].Actions[manifestAction.Action.Name]

		if exists == true {
			if existAction.Filepath == manifestAction.Filepath || manifestAction.Filepath == "" {
				// we're adding a filesystem detected action so just updated code and filepath if needed
				if manifestAction.Action.Exec.Kind != "" {
					existAction.Action.Exec.Kind = manifestAction.Action.Exec.Kind
				}

				if manifestAction.Action.Exec.Code != nil {
					code := *manifestAction.Action.Exec.Code
					if code != "" {
						existAction.Action.Exec.Code = manifestAction.Action.Exec.Code
					}
				}

				existAction.Action.Annotations = manifestAction.Action.Annotations
				existAction.Action.Limits = manifestAction.Action.Limits
				existAction.Action.Parameters = manifestAction.Action.Parameters
				existAction.Action.Version = manifestAction.Action.Version

				if manifestAction.Filepath != "" {
					existAction.Filepath = manifestAction.Filepath
				}

				err := reader.checkAction(existAction)
				if err != nil {
					return utils.NewInputYamlFormatError(err.Error())
				}

			} else {
				// Action exists, but references two different sources
				return errors.New("manifestReader. Error: Conflict detected for action named " + existAction.Action.Name + ". Found two locations for source file: " + existAction.Filepath + " and " + manifestAction.Filepath)
			}
		} else {
			// not a new action so update the action in the package
			err := reader.checkAction(manifestAction)
			if err != nil {
				return utils.NewInputYamlFormatError(err.Error())
			}
			reader.serviceDeployer.Deployment.Packages[manifestAction.Packagename].Actions[manifestAction.Action.Name] = manifestAction
		}
	}

	return nil
}

func (reader *ManifestReader) checkAction(action utils.ActionRecord) error {
	if action.Filepath == "" {
		return errors.New("Error: Action " + action.Action.Name + " has no source code location set.")
	}

	if action.Action.Exec.Kind == "" {
		return errors.New("Error: Action " + action.Action.Name + " has no kind set")
	}

	if action.Action.Exec.Code != nil {
		code := *action.Action.Exec.Code
		if code == "" && action.Action.Exec.Kind != "sequence" {
			return errors.New("Error: Action " + action.Action.Name + " has no source code")
		}
	}

	return nil
}

func (reader *ManifestReader) SetSequences(actions []utils.ActionRecord) error {
	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, seqAction := range actions {
		// check if the sequence action is exist in actions
		// If the sequence action exists in actions, return error
		_, exists := reader.serviceDeployer.Deployment.Packages[seqAction.Packagename].Actions[seqAction.Action.Name]
		if exists == true {
			return errors.New("manifestReader. Error: Conflict sequence action with an action. " +
				"Found a sequence action with the same name of an action:" +
				seqAction.Action.Name)
		}
		existAction, exists := reader.serviceDeployer.Deployment.Packages[seqAction.Packagename].Sequences[seqAction.Action.Name]

		if exists == true {
			existAction.Action.Annotations = seqAction.Action.Annotations
			existAction.Action.Exec.Kind = "sequence"
			existAction.Action.Exec.Components = seqAction.Action.Exec.Components
			existAction.Action.Publish = seqAction.Action.Publish
			existAction.Action.Namespace = seqAction.Action.Namespace
			existAction.Action.Limits = seqAction.Action.Limits
			existAction.Action.Parameters = seqAction.Action.Parameters
			existAction.Action.Version = seqAction.Action.Version
		} else {
			// not a new action so update the action in the package
			err := reader.checkAction(seqAction)
			if err != nil {
				return utils.NewInputYamlFormatError(err.Error())
			}
			reader.serviceDeployer.Deployment.Packages[seqAction.Packagename].Sequences[seqAction.Action.Name] = seqAction
		}
	}

	return nil

}

func (reader *ManifestReader) SetTriggers(triggers []*whisk.Trigger) error {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, trigger := range triggers {
		existTrigger, exist := dep.Deployment.Triggers[trigger.Name]
		if exist {
			existTrigger.Name = trigger.Name
			existTrigger.ActivationId = trigger.ActivationId
			existTrigger.Namespace = trigger.Namespace
			existTrigger.Annotations = trigger.Annotations
			existTrigger.Version = trigger.Version
			existTrigger.Parameters = trigger.Parameters
			existTrigger.Publish = trigger.Publish
		} else {
			dep.Deployment.Triggers[trigger.Name] = trigger
		}

	}
	return nil
}

func (reader *ManifestReader) SetRules(rules []*whisk.Rule) error {
	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, rule := range rules {
		existRule, exist := dep.Deployment.Rules[rule.Name]
		if exist {
			existRule.Name = rule.Name
			existRule.Publish = rule.Publish
			existRule.Version = rule.Version
			existRule.Namespace = rule.Namespace
			existRule.Action = rule.Action
			existRule.Trigger = rule.Trigger
			existRule.Status = rule.Status
		} else {
			dep.Deployment.Rules[rule.Name] = rule
		}

	}
	return nil
}

func (reader *ManifestReader) SetApis(ar []*whisk.ApiCreateRequest) error {
	dep := reader.serviceDeployer
	var apis []*whisk.ApiCreateRequest = make([]*whisk.ApiCreateRequest, 0)

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, api := range apis {
		existApi, exist := dep.Deployment.Apis[api.ApiDoc.ApiName]
		if exist {
			existApi.ApiDoc.ApiName = api.ApiDoc.ApiName
		} else {
			dep.Deployment.Apis[api.ApiDoc.ApiName] = api
		}

	}
	return nil
}

// from whisk go client
func (deployer *ManifestReader) getQualifiedName(name string, namespace string) string {
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
