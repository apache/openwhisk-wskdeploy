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
	"strings"

	"fmt"
	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/dependencies"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskderrors"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
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
	packages, inputs, err := manifestParser.ComposeAllPackages(reader.serviceDeployer.ProjectInputs, manifest, reader.serviceDeployer.ManifestPath, managedAnnotations)
	if err != nil {
		return err
	}
	reader.SetPackages(packages, inputs)
	return nil
}

// Wrapper parser to handle yaml dir
func (reader *ManifestReader) HandleYaml(manifestParser *parsers.YAMLParser, manifest *parsers.YAML, managedAnnotations whisk.KeyValue) error {

	var err error
	var manifestName = manifest.Filepath

	// pull the package inputs out of Deployment.Packages
	// so that they can be sent to all OW entities while they are being parsed and constructed
	// the main reason for sending Package inputs (normalized) to parser is
	// manifest file can have $TRIGGER_NAME as a trigger name, $RULE_NAME as a rule name, etc
	// which are being read by the parser and the way it works is, these variables are
	// treated as environment variables and tried to interpolate them but when these
	// variables are not defined in env., they are returned as an empty string and we
	// loose the env. variable name(TRIGGER_NAME), so before we loose that label,
	// we need to substitute it with its corresponding value from the inputs section
	inputs := make(map[string]parsers.PackageInputs, 0)
	for _, pkg := range reader.serviceDeployer.Deployment.Packages {
		inputs[pkg.Package.Name] = pkg.Inputs
	}

	deps, err := manifestParser.ComposeDependenciesFromAllPackages(manifest, reader.serviceDeployer.ProjectPath, reader.serviceDeployer.ManifestPath, managedAnnotations, inputs)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	actions, err := manifestParser.ComposeActionsFromAllPackages(manifest, reader.serviceDeployer.ManifestPath, managedAnnotations, inputs)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	sequences, err := manifestParser.ComposeSequencesFromAllPackages(reader.serviceDeployer.ClientConfig.Namespace, manifest, reader.serviceDeployer.ManifestPath, managedAnnotations, inputs)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	triggers, err := manifestParser.ComposeTriggersFromAllPackages(manifest, reader.serviceDeployer.ManifestPath, managedAnnotations, inputs)
	if err != nil {
		return wskderrors.NewYAMLFileFormatError(manifestName, err)
	}

	rules, err := manifestParser.ComposeRulesFromAllPackages(manifest, managedAnnotations, inputs)
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

func (reader *ManifestReader) SetPackages(packages map[string]*whisk.Package, inputs map[string]parsers.PackageInputs) error {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, pkg := range packages {
		newPack := NewDeploymentPackage()
		newPack.Package = pkg
		newPack.Inputs = inputs[pkg.Name]
		dep.Deployment.Packages[pkg.Name] = newPack
	}
	return nil
}

func (reader *ManifestReader) SetDependencies(deps map[string]dependencies.DependencyRecord) error {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for name, dependency := range deps {
		// name is <packagename>:<dependencylabel>
		depName := strings.Split(name, ":")[1]
		if len(depName) == 0 {
			return nil
		}
		if !dependency.IsBinding && !reader.IsUndeploy {
			if _, exists := dep.DependencyMaster[depName]; exists {
				if !dependencies.CompareDependencyRecords(dep.DependencyMaster[depName], dependency) {
					location := strings.Join([]string{dep.DependencyMaster[depName].Location, dependency.Location}, ",")
					errmsg := wski18n.T(wski18n.ID_ERR_DEPENDENCIES_WITH_SAME_LABEL_X_dependency_X_location_X,
						map[string]interface{}{wski18n.KEY_DEPENDENCY: depName,
							wski18n.KEY_LOCATION: location})
					return wskderrors.NewYAMLParserErr(dep.ManifestPath, errmsg)
				}
			}
			gitReader := dependencies.NewGitReader(depName, dependency)
			err := gitReader.CloneDependency()
			if err != nil {
				return err
			}
		}
		// store in two places (one local to package to preserve relationship, one in master record to check for conflics
		dep.Deployment.Packages[dependency.Packagename].Dependencies[depName] = dependency
		dep.DependencyMaster[depName] = dependency
	}

	return nil
}

func (reader *ManifestReader) SetActions(actions []utils.ActionRecord) error {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, manifestAction := range actions {
		err := reader.checkAction(manifestAction)
		if err != nil {
			return err
		}
		dep.Deployment.Packages[manifestAction.Packagename].Actions[manifestAction.Action.Name] = manifestAction
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
			err := wski18n.T(wski18n.ID_ERR_SEQUENCE_HAVING_SAME_NAME_AS_ACTION_X_action_X,
				map[string]interface{}{wski18n.KEY_SEQUENCE: sequence.Action.Name})
			return wskderrors.NewYAMLParserErr(reader.serviceDeployer.ManifestPath, err)
		}
		err := reader.checkAction(sequence)
		if err != nil {
			return err
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
		if _, exists := dep.Deployment.Triggers[trigger.Name]; exists {
			var feed string
			var existingFeed string
			for _, a := range dep.Deployment.Triggers[trigger.Name].Annotations {
				if a.Key == parsers.YAML_KEY_FEED {
					existingFeed = a.Value.(string)
				}
			}
			for _, a := range trigger.Annotations {
				if a.Key == parsers.YAML_KEY_FEED {
					feed = a.Value.(string)
				}
			}
			if feed != existingFeed {
				feed = fmt.Sprintf("%q", feed)
				existingFeed = fmt.Sprintf("%q", existingFeed)
				err := wski18n.T(wski18n.ID_ERR_CONFLICTING_TRIGGERS_ACROSS_PACKAGES_X_trigger_X_feed_X,
					map[string]interface{}{wski18n.KEY_TRIGGER: trigger.Name,
						wski18n.KEY_TRIGGER_FEED: strings.Join([]string{feed, existingFeed}, ", ")})
				return wskderrors.NewYAMLParserErr(reader.serviceDeployer.ManifestPath, err)
			}
		}
		dep.Deployment.Triggers[trigger.Name] = trigger
	}
	return nil
}

func (reader *ManifestReader) SetRules(rules []*whisk.Rule) error {
	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, rule := range rules {
		if _, exists := dep.Deployment.Rules[rule.Name]; exists {
			action := rule.Action.(string)
			existingAction := dep.Deployment.Rules[rule.Name].Action.(string)
			trigger := rule.Trigger.(string)
			existingTrigger := dep.Deployment.Rules[rule.Name].Trigger.(string)
			if action != existingAction || trigger != existingTrigger {
				action = fmt.Sprintf("%q", action)
				existingAction = fmt.Sprintf("%q", existingAction)
				trigger = fmt.Sprintf("%q", trigger)
				existingTrigger = fmt.Sprintf("%q", existingTrigger)
				err := wski18n.T(wski18n.ID_ERR_CONFLICTING_RULES_ACROSS_PACKAGES_X_rule_X_action_X_trigger_X,
					map[string]interface{}{wski18n.KEY_RULE: rule.Name,
						wski18n.KEY_TRIGGER: strings.Join([]string{trigger, existingTrigger}, ", "),
						wski18n.KEY_ACTION:  strings.Join([]string{action, existingAction}, ", ")})
				return wskderrors.NewYAMLParserErr(reader.serviceDeployer.ManifestPath, err)
			}
		}
		dep.Deployment.Rules[rule.Name] = rule
	}
	return nil
}

func (reader *ManifestReader) SetApis(ar []*whisk.ApiCreateRequest) error {
	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, api := range ar {
		apiPath := api.ApiDoc.ApiName + api.ApiDoc.GatewayBasePath + api.ApiDoc.GatewayRelPath + api.ApiDoc.GatewayMethod

		// uniqueness issue when using action name as key as there can be multiple APIs pointing to same action.
		// using apiPath instead as it is uniqueue
		dep.Deployment.Apis[apiPath] = api
	}

	return nil
}

// Check action record before deploying it
// action record is created by reading and composing action elements from manifest file
// Action.kind is mandatory which is set to
// (1) action runtime for an action and (2) set to "sequence" for a sequence
// Also, action executable code should be specified for any action
func (reader *ManifestReader) checkAction(action utils.ActionRecord) error {
	if action.Action.Exec.Kind == "" {
		err := wski18n.T(wski18n.ID_ERR_ACTION_WITHOUT_KIND_X_action_X,
			map[string]interface{}{wski18n.KEY_ACTION: action.Action.Name})
		return wskderrors.NewYAMLParserErr(reader.serviceDeployer.ManifestPath, err)
	}

	if action.Action.Exec.Code != nil {
		code := *action.Action.Exec.Code
		if code == "" && action.Action.Exec.Kind != parsers.YAML_KEY_SEQUENCE {
			err := wski18n.T(wski18n.ID_ERR_ACTION_WITHOUT_SOURCE_X_action_X,
				map[string]interface{}{wski18n.KEY_ACTION: action.Action.Name})
			return wskderrors.NewYAMLParserErr(reader.serviceDeployer.ManifestPath, err)
		}
	}

	return nil
}
