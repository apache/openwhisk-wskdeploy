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

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/deployers"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskderrors"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:        "export",
	SuggestFor: []string{"capture"},
	Short:      "Export project assets from OpenWhisk",
	Long:       `Exports managed project assets from OpenWhisk to manifest and function files`,
	RunE:       ExportCmdImp,
}

var config *whisk.Config

func ExportRule(wskRule whisk.Rule, pkgName string, maniyaml *parsers.YAML) {
	if maniyaml.Packages[pkgName].Rules == nil {
		pkg := maniyaml.Packages[pkgName]
		pkg.Rules = make(map[string]parsers.Rule)
		maniyaml.Packages[pkgName] = pkg
	}

	// export rule to manifest
	maniyaml.Packages[pkgName].Rules[wskRule.Name] = *maniyaml.ComposeParsersRule(wskRule)
}

func ExportAction(actionName string, packageName string, maniyaml *parsers.YAML) error {

	pkg := maniyaml.Packages[packageName]
	if pkg.Actions == nil {
		pkg.Actions = make(map[string]parsers.Action)
		maniyaml.Packages[packageName] = pkg
	}

	wskAction, _, err := client.Actions.Get(actionName, true)
	if err != nil {
		return err
	}
	if wskAction.Exec.Kind == "sequence" {
		seq := new(parsers.Sequence)
		for _, component := range wskAction.Exec.Components {
			// must ommit namespace from seq component name
			ExportAction(strings.SplitN(component, "/", 3)[2], packageName, maniyaml)
			slices := strings.Split(component, "/")

			// save in the seq list only action names
			if len(seq.Actions) > 0 {
				seq.Actions += ","
			}

			seq.Actions += slices[len(slices)-1]
		}

		pkg = maniyaml.Packages[packageName]
		if pkg.Sequences == nil {
			pkg.Sequences = make(map[string]parsers.Sequence)
		}

		pkg.Sequences[wskAction.Name] = *seq
	} else {
		parsedAction := *maniyaml.ComposeParsersAction(*wskAction)

		// get the action file extension according to action kind (e.g. js for nodejs)
		ext := utils.FileRuntimeExtensionsMap[wskAction.Exec.Kind]

		manifestDir := filepath.Dir(utils.Flags.ManifestPath)

		// store function file under action package name subdirectory in the specified manifest folder
		functionDir := filepath.Join(manifestDir, packageName)
		os.MkdirAll(functionDir, os.ModePerm)

		// store function in manifest under path relative to manifest root
		functionFile := filepath.Join(packageName, wskAction.Name) + "." + ext
		parsedAction.Function = functionFile

		// create function file at the full path
		functionFile = filepath.Join(manifestDir, functionFile)
		f, err := os.Create(functionFile)
		if err != nil {
			return wskderrors.NewFileReadError(functionFile, err.Error())
		}

		defer f.Close()

		// store action function in the filesystem next to the manifest.yml
		// TODO: consider to name files by namespace + action to make function file names unique
		if wskAction.Exec.Code != nil {
			f.Write([]byte(*wskAction.Exec.Code))
		}
		pkg.Actions[wskAction.Name] = parsedAction
	}

	maniyaml.Packages[packageName] = pkg
	return nil
}

func ContainsManageProjectAnnotation(annotations whisk.KeyValueArr, projectName string) (res bool) {
	// get the list of managed projects
	if annotation := annotations.GetValue(utils.MANAGED_LIST); annotation != nil {
		// decode the JSON blob and retrieve __OW_PROJECT_NAME_LIST
		managedList := annotation.([]interface{})

		for _, a := range managedList {
			managed := a.(map[string]interface{})
			if managed[utils.OW_PROJECT_NAME] == projectName {
				return true
			}
		}
	}

	// for backward compatibility check 'managed' annotation as well
	if a := annotations.GetValue(utils.MANAGED); a != nil {
		// decode the JSON blob and retrieve __OW_PROJECT_NAME
		pa := a.(map[string]interface{})

		// we have found a package which is part of the current project
		if pa[utils.OW_PROJECT_NAME] == projectName {
			return true
		}
	}

	return false
}

func ExportCmdImp(cmd *cobra.Command, args []string) error {

	projectName := utils.Flags.ProjectPath
	maniyaml := &parsers.YAML{}
	maniyaml.Project.Name = projectName

	config, _ = deployers.NewWhiskConfig(wskpropsPath, utils.Flags.DeploymentPath, utils.Flags.ManifestPath, false)
	client, _ = deployers.CreateNewClient(config)

	// Init supported runtimes and action files extensions maps
	setSupportedRuntimes(config.Host)

	// Get the list of packages in your namespace
	packages, _, err := client.Packages.List(&whisk.PackageListOptions{})
	if err != nil {
		return err
	}

	// iterate over each package to find managed annotations
	// check if "managed" annotation is attached to a package
	// add to export when managed project name matches with the
	// specified project name
	for _, pkg := range packages {
		if ContainsManageProjectAnnotation(pkg.Annotations, projectName) {
			if maniyaml.Packages == nil {
				maniyaml.Packages = make(map[string]parsers.Package)
			}

			maniyaml.Packages[pkg.Name] = *maniyaml.ComposeParsersPackage(pkg)

			// perform the similar check on the list of actions from this package
			// get a list of actions in your namespace
			actions, _, err := client.Actions.List(pkg.Name, &whisk.ActionListOptions{})
			if err != nil {
				return err
			}

			// iterate over list of managed package actions to find an action with managed annotations
			// check if "managed" annotation is attached to an action
			for _, action := range actions {
				// TODO: consider to throw error when there unmanaged or "foreign" assets under managed package
				// an annotation with "managed" key indicates that an action was deployed as part of managed deployment
				// if such annotation exists, check if it belongs to the current managed deployment
				// this action has attached managed annotations
				if ContainsManageProjectAnnotation(action.Annotations, projectName) {
					actionName := strings.Join([]string{pkg.Name, action.Name}, "/")
					// export action to file system
					err = ExportAction(actionName, pkg.Name, maniyaml)
					if err != nil {
						return err
					}

				}
			}

		}
	}

	// Get list of triggers in your namespace
	triggers, _, err := client.Triggers.List(&whisk.TriggerListOptions{})
	if err != nil {
		return err
	}

	// iterate over the list of triggers to determine whether any of them part of specified managed project
	for _, trg := range triggers {
		// trigger has attached managed annotation
		if ContainsManageProjectAnnotation(trg.Annotations, projectName) {
			for pkgName := range maniyaml.Packages {
				if maniyaml.Packages[pkgName].Namespace == trg.Namespace {
					if maniyaml.Packages[pkgName].Triggers == nil {
						pkg := maniyaml.Packages[pkgName]
						pkg.Triggers = make(map[string]parsers.Trigger)
						maniyaml.Packages[pkgName] = pkg
					}

					// export trigger to manifest
					maniyaml.Packages[pkgName].Triggers[trg.Name] = *maniyaml.ComposeParsersTrigger(trg)
				}
			}
		}
	}

	// Get list of rules from OW
	rules, _, err := client.Rules.List(&whisk.RuleListOptions{})
	if err != nil {
		return err
	}

	// TODO: can be simplifyed once OW permits to add annotation to rules
	// iterate over the list of rules to determine whether any of them is part of
	// managed trigger -> action set for specified managed project. if yes, add to manifest
	for _, rule := range rules {
		// get rule from OW
		wskRule, _, _ := client.Rules.Get(rule.Name)
		ruleAction := wskRule.Action.(map[string]interface{})["name"].(string)
		ruleTrigger := wskRule.Trigger.(map[string]interface{})["name"].(string)

		// can be simplified once rules moved to top level next to triggers
		for pkgName := range maniyaml.Packages {
			if maniyaml.Packages[pkgName].Namespace == wskRule.Namespace {
				// iterate over all managed triggers in manifest
				for _, trigger := range maniyaml.Packages[pkgName].Triggers {
					// check that managed trigger equals to rule trigger
					if ruleTrigger == trigger.Name && trigger.Namespace == wskRule.Namespace {
						// check that there a managed action (or sequence action) equals to rule action
						for _, action := range maniyaml.Packages[pkgName].Actions {
							if action.Name == ruleAction {
								// export rule to manifest
								ExportRule(*wskRule, pkgName, maniyaml)
							}
						}

						// check that there a managed sequence action with name matching the rule action
						for name := range maniyaml.Packages[pkgName].Sequences {
							if name == ruleAction {
								// export rule to manifest
								ExportRule(*wskRule, pkgName, maniyaml)
							}
						}
					}
				}
			}
		}
	}

	parsers.Write(maniyaml, utils.Flags.ManifestPath)
	fmt.Println("Manifest exported to: " + utils.Flags.ManifestPath)

	return nil
}

func init() {
	RootCmd.AddCommand(exportCmd)
}
