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
	"strings"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskderrors"
)

var clientConfig *whisk.Config

type ManifestReader struct {
	serviceDeployer *ServiceDeployer
	IsUndeploy      bool
}

func NewManifestReader(serviceDeployer *ServiceDeployer) *ManifestReader {
	var dep ManifestReader
	dep.serviceDeployer = serviceDeployer

	return &dep
}

func (deployer *ManifestReader) ParseManifest() (*parsers.YAML, *parsers.YAMLParser, error) {
	dep := deployer.serviceDeployer
	manifestParser := parsers.NewYAMLParser()
	manifest, err := manifestParser.ParseManifest(dep.ManifestPath)

	if err != nil {
		return manifest, manifestParser, err
	}
	return manifest, manifestParser, nil
}

func (reader *ManifestReader) InitPackages(manifestParser *parsers.YAMLParser, manifest *parsers.YAML, ma whisk.KeyValue) error {
	packages, err := manifestParser.ComposeAllPackages(manifest, reader.serviceDeployer.ManifestPath, ma)
	if err != nil {
		return err
	}
	reader.SetPackages(packages)

	return nil
}

// Wrapper parser to handle yaml dir
func (deployer *ManifestReader) HandleYaml(sdeployer *ServiceDeployer, manifestParser *parsers.YAMLParser, manifest *parsers.YAML, ma whisk.KeyValue) error {

	var err error
	var manifestName = manifest.Filepath

	deps, err := manifestParser.ComposeDependenciesFromAllPackages(manifest, deployer.serviceDeployer.ProjectPath, deployer.serviceDeployer.ManifestPath)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	actions, err := manifestParser.ComposeActionsFromAllPackages(manifest, deployer.serviceDeployer.ManifestPath, ma)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	sequences, err := manifestParser.ComposeSequencesFromAllPackages(deployer.serviceDeployer.ClientConfig.Namespace, manifest, ma)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	triggers, err := manifestParser.ComposeTriggersFromAllPackages(manifest, deployer.serviceDeployer.ManifestPath, ma)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	rules, err := manifestParser.ComposeRulesFromAllPackages(manifest, ma)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	apis, err := manifestParser.ComposeApiRecordsFromAllPackages(deployer.serviceDeployer.ClientConfig, manifest)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	err = deployer.SetDependencies(deps)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	err = deployer.SetActions(actions)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	err = deployer.SetSequences(sequences)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	err = deployer.SetTriggers(triggers)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	err = deployer.SetRules(rules)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	err = deployer.SetApis(apis)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	return nil
}

func (reader *ManifestReader) SetDependencies(deps map[string]utils.DependencyRecord) error {
	for name, dep := range deps {
		n := strings.Split(name, ":")
		depName := n[1]
		if depName == "" {
			return nil
		}
		if !dep.IsBinding && !reader.IsUndeploy {
			if _, exists := reader.serviceDeployer.DependencyMaster[depName]; !exists {
				// dependency
				gitReader := utils.NewGitReader(depName, dep)
				err := gitReader.CloneDependency()
				if err != nil {
					return wskderrors.NewYAMLFileFormatError(depName, err)
				}
			} else {
				// TODO: we should do a check to make sure this dependency is compatible with an already installed one.
				// If not, we should throw dependency mismatch error.
			}
		}

		// store in two places (one local to package to preserve relationship, one in master record to check for conflics
		reader.serviceDeployer.Deployment.Packages[dep.Packagename].Dependencies[depName] = dep
		reader.serviceDeployer.DependencyMaster[depName] = dep

	}

	return nil
}

func (reader *ManifestReader) SetPackages(packages map[string]*whisk.Package) error {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, pkg := range packages {
		_, exist := dep.Deployment.Packages[pkg.Name]
		if exist {
			// TODO(): i18n of error message (or create a new named error)
			// TODO(): Is there a better way to handle an existing dependency of same name?
			err := errors.New("Package [" + pkg.Name + "] exists already.")
			return wskderrors.NewYAMLParserErr(reader.serviceDeployer.ManifestPath, err)
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
					return wskderrors.NewFileReadError(manifestAction.Filepath, err)
				}

			} else {
				// Action exists, but references two different sources
				// TODO(): i18n of error message (or create a new named error)
				err := errors.New("Conflict detected for action named [" +
					existAction.Action.Name + "].\nFound two locations for source file: [" +
					existAction.Filepath + "] and [" + manifestAction.Filepath + "]")
				return wskderrors.NewYAMLParserErr(reader.serviceDeployer.ManifestPath, err)
			}
		} else {
			// not a new action so update the action in the package
			err := reader.checkAction(manifestAction)
			if err != nil {
				return wskderrors.NewFileReadError(manifestAction.Filepath, err)
			}
			reader.serviceDeployer.Deployment.Packages[manifestAction.Packagename].Actions[manifestAction.Action.Name] = manifestAction
		}
	}

	return nil
}

// TODO create named errors
func (reader *ManifestReader) checkAction(action utils.ActionRecord) error {
	if action.Filepath == "" {
		// TODO(): i18n of error message (or create a new named error)
		err := errors.New("Action [" + action.Action.Name + "] has no source code location set.")
		return wskderrors.NewYAMLParserErr(reader.serviceDeployer.ManifestPath, err)
	}

	if action.Action.Exec.Kind == "" {
		// TODO(): i18n of error message (or create a new named error)
		err := errors.New("Action [" + action.Action.Name + "] has no kind set.")
		return wskderrors.NewYAMLParserErr(reader.serviceDeployer.ManifestPath, err)
	}

	if action.Action.Exec.Code != nil {
		code := *action.Action.Exec.Code
		if code == "" && action.Action.Exec.Kind != "sequence" {
			// TODO(): i18n of error message (or create a new named error)
			err := errors.New("Action [" + action.Action.Name + "] has no source code.")
			return wskderrors.NewYAMLParserErr(reader.serviceDeployer.ManifestPath, err)
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
			// TODO(): i18n of error message (or create a new named error)
			err := errors.New("Sequence action's name [" +
				seqAction.Action.Name + "] is already used by an action.")
			return wskderrors.NewYAMLParserErr(reader.serviceDeployer.ManifestPath, err)
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
				// TODO() Need a better error type here
				return wskderrors.NewFileReadError(seqAction.Filepath, err)
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

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, api := range ar {
		existApi, exist := dep.Deployment.Apis[api.ApiDoc.Action.Name]
		if exist {
			existApi.ApiDoc.ApiName = api.ApiDoc.ApiName
		} else {
			dep.Deployment.Apis[api.ApiDoc.Action.Name] = api
		}

	}
	return nil
}
