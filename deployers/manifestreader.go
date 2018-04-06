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

func (reader *ManifestReader) InitPackages(manifestParser *parsers.YAMLParser, manifest *parsers.YAML, managedAnnotations whisk.KeyValue) error {
	packages, err := manifestParser.ComposeAllPackages(manifest, reader.serviceDeployer.ManifestPath, managedAnnotations)
	if err != nil {
		return err
	}
	reader.SetPackages(packages)
	return nil
}

// Wrapper parser to handle yaml dir
func (reader *ManifestReader) HandleYaml(sdeployer *ServiceDeployer, manifestParser *parsers.YAMLParser, manifest *parsers.YAML, managedAnnotations whisk.KeyValue) error {

	var err error
	var manifestName = manifest.Filepath

	deps, err := manifestParser.ComposeDependenciesFromAllPackages(manifest, reader.serviceDeployer.ProjectPath, reader.serviceDeployer.ManifestPath, managedAnnotations)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	actions, err := manifestParser.ComposeActionsFromAllPackages(manifest, reader.serviceDeployer.ManifestPath, managedAnnotations)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	sequences, err := manifestParser.ComposeSequencesFromAllPackages(reader.serviceDeployer.ClientConfig.Namespace, manifest, reader.serviceDeployer.ManifestPath, managedAnnotations)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	triggers, err := manifestParser.ComposeTriggersFromAllPackages(manifest, reader.serviceDeployer.ManifestPath, managedAnnotations)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	rules, err := manifestParser.ComposeRulesFromAllPackages(manifest, managedAnnotations)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	apis, err := manifestParser.ComposeApiRecordsFromAllPackages(reader.serviceDeployer.ClientConfig, manifest)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	err = reader.SetDependencies(deps)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	err = reader.SetActions(actions)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	err = reader.SetSequences(sequences)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	err = reader.SetTriggers(triggers)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	err = reader.SetRules(rules)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	err = reader.SetApis(apis)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	return nil
}

func (reader *ManifestReader) SetDependencies(deps map[string]utils.DependencyRecord) error {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for name, dependency := range deps {
		n := strings.Split(name, ":")
		depName := n[1]
		if depName == "" {
			return nil
		}
		if !dependency.IsBinding && !reader.IsUndeploy {
			if _, exists := dep.DependencyMaster[depName]; exists {
				if !utils.CompareDependencyRecords(dep.DependencyMaster[depName], dependency) {
					// return error
					err := errors.New("Dependecies have same name")
					return wskderrors.NewYAMLFileFormatError(depName, err)
				}
			}
			gitReader := utils.NewGitReader(depName, dependency)
			err := gitReader.CloneDependency()
			if err != nil {
				return wskderrors.NewYAMLFileFormatError(depName, err)
			}
		}
		// store in two places (one local to package to preserve relationship, one in master record to check for conflics
		dep.Deployment.Packages[dependency.Packagename].Dependencies[depName] = dependency
		dep.DependencyMaster[depName] = dependency
	}

	return nil
}

func (reader *ManifestReader) SetPackages(packages map[string]*whisk.Package) error {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, pkg := range packages {
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
		// not a new action so update the action in the package
		err := reader.checkAction(manifestAction)
		if err != nil {
			return wskderrors.NewFileReadError(manifestAction.Filepath, err)
		}
		dep.Deployment.Packages[manifestAction.Packagename].Actions[manifestAction.Action.Name] = manifestAction
	}
	return nil
}

// TODO create named errors
// Check action record before deploying it
// action record is created by reading and composing action elements from manifest file
// Action.kind is mandatory which is set to
// (1) action runtime for an action and (2) set to "sequence" for a sequence
// Also, action executable code should be specified for any action
func (reader *ManifestReader) checkAction(action utils.ActionRecord) error {
	if action.Action.Exec.Kind == "" {
		// TODO(): i18n of error message (or create a new named error)
		err := errors.New("Action [" + action.Action.Name + "] has no kind set.")
		return wskderrors.NewYAMLParserErr(reader.serviceDeployer.ManifestPath, err)
	}

	if action.Action.Exec.Code != nil {
		code := *action.Action.Exec.Code
		if code == "" && action.Action.Exec.Kind != parsers.YAML_KEY_SEQUENCE {
			// TODO(): i18n of error message (or create a new named error)
			err := errors.New("Action [" + action.Action.Name + "] has no source code.")
			return wskderrors.NewYAMLParserErr(reader.serviceDeployer.ManifestPath, err)
		}
	}

	return nil
}

func (reader *ManifestReader) SetSequences(sequences []utils.ActionRecord) error {
	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, sequence := range sequences {
		// If the sequence name matches with any of the actions defined, return error
		if _, exists := dep.Deployment.Packages[sequence.Packagename].Actions[sequence.Action.Name]; exists {
			// TODO(798): i18n of error message (or create a new named error)
			err := errors.New("Sequence action's name [" +
				sequence.Action.Name + "] is already used by an action.")
			return wskderrors.NewYAMLParserErr(reader.serviceDeployer.ManifestPath, err)
		}
		err := reader.checkAction(sequence)
		if err != nil {
			// TODO() Need a better error type here
			return wskderrors.NewFileReadError(sequence.Filepath, err)
		}
		dep.Deployment.Packages[sequence.Packagename].Sequences[sequence.Action.Name] = sequence
	}

	return nil

}

func (reader *ManifestReader) SetTriggers(triggers []*whisk.Trigger) error {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, trigger := range triggers {
		dep.Deployment.Triggers[trigger.Name] = trigger
	}
	return nil
}

func (reader *ManifestReader) SetRules(rules []*whisk.Rule) error {
	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, rule := range rules {
		dep.Deployment.Rules[rule.Name] = rule
	}
	return nil
}

func (reader *ManifestReader) SetApis(ar []*whisk.ApiCreateRequest) error {
	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, api := range ar {
		dep.Deployment.Apis[api.ApiDoc.Action.Name] = api
	}
	return nil
}
