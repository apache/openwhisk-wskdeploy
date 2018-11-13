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
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/deployers"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/runtimes"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:        "export",
	SuggestFor: []string{"capture"},
	Short:      wski18n.T(wski18n.ID_CMD_DESC_SHORT_EXPORT),
	Long:       wski18n.T(wski18n.ID_CMD_DESC_LONG_EXPORT),
	RunE:       ExportCmdImp,
}

var config *whisk.Config
var wskpropsPath string
var client *whisk.Client

func ExportAction(actionName string, packageName string, maniyaml *parsers.YAML, targetManifest string, projectName string) error {

	pkg := maniyaml.Packages[packageName]
	if pkg.Actions == nil {
		pkg.Actions = make(map[string]parsers.Action)
		maniyaml.Packages[packageName] = pkg
	}

	wskAction, _, err := client.Actions.Get(actionName, true)
	if err != nil {
		return err
	}

	if a := wskAction.Annotations.GetValue(utils.MANAGED); a != nil {
		// decode the JSON blob and retrieve __OW_PROJECT_NAME
		pa := a.(map[string]interface{})

		// we have found a package which is part of the current project
		if pa[utils.OW_PROJECT_NAME] != projectName {
			return nil
		}
	} else {
		return nil
	}

	if wskAction.Exec.Kind == "sequence" {
		seq := new(parsers.Sequence)
		for _, component := range wskAction.Exec.Components {
			// must ommit namespace from seq component name
			ExportAction(strings.SplitN(component, "/", 3)[2], packageName, maniyaml, targetManifest, projectName)

			if len(seq.Actions) > 0 {
				seq.Actions += ","
			}
			// save action in the seq list as package/action
			seq.Actions += strings.SplitN(component, "/", 3)[2]
		}

		pkg = maniyaml.Packages[packageName]
		if pkg.Sequences == nil {
			pkg.Sequences = make(map[string]parsers.Sequence)
		}

		pkg.Sequences[wskAction.Name] = *seq
	} else {
		parsedAction := *maniyaml.ComposeParsersAction(*wskAction)
		manifestDir := filepath.Dir(targetManifest)

		// store function file under action package name subdirectory in the specified manifest folder
		functionDir := filepath.Join(manifestDir, packageName)

		// save the code in file
		filename, err := saveCode(*wskAction, functionDir)
		if err != nil {
			return err
		}

		if filename != "" {
			// store function in manifest if action has code section
			parsedAction.Function = packageName + "/" + filename
		}

		pkg.Actions[wskAction.Name] = parsedAction
	}

	maniyaml.Packages[packageName] = pkg
	return nil
}

func exportProject(projectName string, targetManifest string) error {

	whisk.SetVerbose(utils.Flags.Verbose)
	whisk.SetDebug(utils.Flags.Trace)

	maniyaml := &parsers.YAML{}
	maniyaml.Project.Name = projectName

	// Get the list of packages in your namespace
	packages, _, err := client.Packages.List(&whisk.PackageListOptions{})
	if err != nil {
		return err
	}

	var bindings = make(map[string]whisk.Binding)

	// iterate over each package to find managed annotations
	// check if "managed" annotation is attached to a package
	// add to export when managed project name matches with the
	// specified project name
	for _, pkg := range packages {
		if a := pkg.Annotations.GetValue(utils.MANAGED); a != nil {
			// decode the JSON blob and retrieve __OW_PROJECT_NAME
			pa := a.(map[string]interface{})

			// we have found a package which is part of the current project
			if pa[utils.OW_PROJECT_NAME] == projectName {

				// check if the package is dependency
				if pkg.Annotations.GetValue(wski18n.BINDING) != nil {
					bindings[pkg.Name] = *pkg.Binding
					continue
				}

				if maniyaml.Packages == nil {
					maniyaml.Packages = make(map[string]parsers.Package)
				}

				maniyaml.Packages[pkg.Name] = *maniyaml.ComposeParsersPackage(pkg)
				// TODO: throw if there more than single package managed by project
				// currently will be a mess because triggers and rules managed under packages
				// instead of the top level (similar to OW model)
				//              if len(maniyaml.Packages) > 1 {
				//                  return errors.New("currently can't work with more than one package managed by one project")
				//              }

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
					if a := action.Annotations.GetValue(utils.MANAGED); a != nil {
						aa := a.(map[string]interface{})
						if aa[utils.OW_PROJECT_NAME] == projectName {
							actionName := strings.Join([]string{pkg.Name, action.Name}, "/")
							// export action to file system
							err = ExportAction(actionName, pkg.Name, maniyaml, targetManifest, projectName)
							if err != nil {
								return err
							}
						}
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
		if a := trg.Annotations.GetValue(utils.MANAGED); a != nil {
			// decode the JSON blob and retrieve __OW_PROJECT_NAME
			ta := a.(map[string]interface{})
			if ta[utils.OW_PROJECT_NAME] == projectName {

				for pkgName := range maniyaml.Packages {
					if maniyaml.Packages[pkgName].Namespace == trg.Namespace {
						if maniyaml.Packages[pkgName].Triggers == nil {
							pkg := maniyaml.Packages[pkgName]
							pkg.Triggers = make(map[string]parsers.Trigger)
							maniyaml.Packages[pkgName] = pkg
						}

						// export trigger to manifest

						if feedname, isFeed := utils.IsFeedAction(&trg); isFeed {
							// export feed input parameters
							feedAction, _, _ := client.Actions.Get(feedname, true)
							if err != nil {
								return err
							}

							params := make(map[string]interface{})
							params["authKey"] = client.Config.AuthToken
							params["lifecycleEvent"] = "READ"
							params["triggerName"] = "/" + client.Namespace + "/" + trg.Name
							res, _, err := client.Actions.Invoke(feedname, params, true, true)
							if err != nil {
								return err
							}
							feedConfig := res["config"]
							for key, val := range feedConfig.(map[string]interface{}) {

								if i := feedAction.Parameters.FindKeyValue(key); i >= 0 {
									trg.Parameters = trg.Parameters.AddOrReplace(&whisk.KeyValue{Key: key, Value: val})
								}
							}
						}

						maniyaml.Packages[pkgName].Triggers[trg.Name] = *maniyaml.ComposeParsersTrigger(trg)
					}
				}
			}
		}
	}

	// Get list of rules from OW
	rules, _, err := client.Rules.List(&whisk.RuleListOptions{})
	if err != nil {
		return err
	}

	// iterate over the list of rules to determine whether any of them is part of the manage dproject
	for _, rule := range rules {
		// get rule from OW
		wskRule, _, _ := client.Rules.Get(rule.Name)
		// rule has attached managed annotation
		if a := wskRule.Annotations.GetValue(utils.MANAGED); a != nil {
			// decode the JSON blob and retrieve __OW_PROJECT_NAME
			ta := a.(map[string]interface{})
			if ta[utils.OW_PROJECT_NAME] == projectName {

				for pkgName := range maniyaml.Packages {
					if maniyaml.Packages[pkgName].Namespace == wskRule.Namespace {
						if maniyaml.Packages[pkgName].Rules == nil {
							pkg := maniyaml.Packages[pkgName]
							pkg.Rules = make(map[string]parsers.Rule)
							maniyaml.Packages[pkgName] = pkg
						}

						// export rule to manifest
						maniyaml.Packages[pkgName].Rules[wskRule.Name] = *maniyaml.ComposeParsersRule(*wskRule)
					}
				}
			}
		}

	}

	// API Gateway is an optional component. Export APIs only when ApigwAccessToken is configured
	if len(client.ApigwAccessToken) == 0 {
		warningString := wski18n.T(wski18n.ID_MSG_CONFIG_MISSING_APIGW_ACCESS_TOKEN)
		wskprint.PrintOpenWhiskWarning(warningString)
	} else {

		// List API request query parameters
		apiListReqOptions := new(whisk.ApiListRequestOptions)
		apiListReqOptions.SpaceGuid = strings.Split(client.Config.AuthToken, ":")[0]
		apiListReqOptions.AccessToken = client.Config.ApigwAccessToken

		// Get list of APIs from OW
		retApiList, _, err := client.Apis.List(apiListReqOptions)
		if err != nil {
			return err
		}

		// iterate over the list of APIs to determine whether any of them part of the managed project
		retApiArray := (*whisk.RetApiArray)(retApiList)
		for _, api := range retApiArray.Apis {

			apiName := api.ApiValue.Swagger.Info.Title
			apiBasePath := strings.TrimPrefix(api.ApiValue.Swagger.BasePath, "/")

			// run over api paths looking for one pointing to an action belonging to the given project
			for path := range api.ApiValue.Swagger.Paths {
				for op, opv := range api.ApiValue.Swagger.Paths[path].MakeOperationMap() {
					if len(opv.XOpenWhisk.Package) > 0 {
						pkgName := opv.XOpenWhisk.Package

						if pkg, ok := maniyaml.Packages[pkgName]; ok {
							if pkg.Namespace == opv.XOpenWhisk.Namespace {

								// now adding the api to the maniyaml
								if pkg.Apis == nil {
									pkg.Apis = make(map[string]map[string]map[string]map[string]parsers.APIMethodResponse)
								}

								path = strings.TrimPrefix(path, "/")

								apiMethodResponse := *new(parsers.APIMethodResponse)
								splitApiUrl := strings.Split(opv.XOpenWhisk.ApiUrl, ".")
								responseType := splitApiUrl[len(splitApiUrl)-1]

								apiMethodResponse.Method = op
								apiMethodResponse.Response = responseType

								if pkgApi, ok := pkg.Apis[apiName]; ok {
									if pkgApiBasePath, ok := pkgApi[apiBasePath]; ok {
										if _, ok := pkgApiBasePath[path]; ok {
											pkg.Apis[apiName][apiBasePath][path][opv.XOpenWhisk.ActionName] = apiMethodResponse
										} else {
											pkg.Apis[apiName][apiBasePath][path] = map[string]parsers.APIMethodResponse{}
											pkg.Apis[apiName][apiBasePath][path][opv.XOpenWhisk.ActionName] = apiMethodResponse
										}
									} else {
										pkg.Apis[apiName][apiBasePath] = map[string]map[string]parsers.APIMethodResponse{}
										pkg.Apis[apiName][apiBasePath][path] = map[string]parsers.APIMethodResponse{}
										pkg.Apis[apiName][apiBasePath][path][opv.XOpenWhisk.ActionName] = apiMethodResponse
									}
								} else {
									pkg.Apis[apiName] = map[string]map[string]map[string]parsers.APIMethodResponse{}
									pkg.Apis[apiName][apiBasePath] = map[string]map[string]parsers.APIMethodResponse{}
									pkg.Apis[apiName][apiBasePath][path] = map[string]parsers.APIMethodResponse{}
									pkg.Apis[apiName][apiBasePath][path][opv.XOpenWhisk.ActionName] = apiMethodResponse
								}

								maniyaml.Packages[pkgName] = pkg
							}
						}
					}
				}
			}
		}
	}

	// adding dependencies to the first package
	for pkgName := range maniyaml.Packages {
		for bPkg, binding := range bindings {
			if maniyaml.Packages[pkgName].Dependencies == nil {
				pkg := maniyaml.Packages[pkgName]
				pkg.Dependencies = make(map[string]parsers.Dependency)
				maniyaml.Packages[pkgName] = pkg
			}

			bPkgData, _, err := client.Packages.Get(bPkg)
			if err != nil {
				return err
			}

			maniyaml.Packages[pkgName].Dependencies[bPkg] = *maniyaml.ComposeParsersDependency(binding, *bPkgData)
		}

		break
	}

	// find exported manifest parent directory
	manifestDir := filepath.Dir(utils.Flags.ManifestPath)
	os.MkdirAll(manifestDir, os.ModePerm)

	// export manifest to file
	parsers.Write(maniyaml, targetManifest)
	fmt.Println("Manifest exported to: " + targetManifest)

	// create dependencies directory if not exists
	depDir := filepath.Join(manifestDir, "dependencies")

	if len(bindings) > 0 {
		fmt.Println("Exporting project dependencies to " + depDir)
	}

	// now export dependencies to their own manifests
	for _, binding := range bindings {
		ns := client.Config.Namespace
		client.Config.Namespace = binding.Namespace

		pkg, _, err := client.Packages.Get(binding.Name)
		if err != nil {
			return err
		}

		if a := pkg.Annotations.GetValue(utils.MANAGED); a != nil {
			// decode the JSON blob and retrieve __OW_PROJECT_NAME
			pa := a.(map[string]interface{})

			os.MkdirAll(depDir, os.ModePerm)
			depManifestPath := filepath.Join(depDir, pa[utils.OW_PROJECT_NAME].(string)+".yaml")

			// export the whole project as dependency
			err := exportProject(pa[utils.OW_PROJECT_NAME].(string), depManifestPath)
			if err != nil {
				return err
			}
		} else {
			// showing warning to notify user that exported manifest dependent on unmanaged library which can't be exported
			fmt.Println("Warning! Dependency package " + binding.Name + " currently unmanaged by any project. Unable to export this package")
		}
		client.Config.Namespace = ns
	}

	return nil
}

func ExportCmdImp(cmd *cobra.Command, args []string) error {

	config, _ = deployers.NewWhiskConfig(wskpropsPath, utils.Flags.DeploymentPath, utils.Flags.ManifestPath)
	client, _ = deployers.CreateNewClient(config)

	// Init supported runtimes and action files extensions maps
	setSupportedRuntimes(config.Host)

	return exportProject(utils.Flags.ProjectName, utils.Flags.ManifestPath)
}

const (
	JAVA_EXT = ".jar"
	ZIP_EXT  = ".zip"
	JAVA     = "java"
)

func getBinaryKindExtension(runtime string) (extension string) {
	switch strings.ToLower(runtime) {
	case JAVA:
		extension = JAVA_EXT
	default:
		extension = ZIP_EXT
	}

	return extension
}

func saveCode(action whisk.Action, directory string) (string, error) {
	var code string
	var runtime string
	var exec whisk.Exec

	exec = *action.Exec
	runtime = strings.Split(exec.Kind, ":")[0]

	if exec.Code != nil {
		code = *exec.Code
	} else {
		return "", nil
	}

	var filename = ""
	if *exec.Binary {
		decoded, _ := base64.StdEncoding.DecodeString(code)
		code = string(decoded)

		filename = action.Name + getBinaryKindExtension(runtime)
	} else {
		filename = action.Name + "." + runtimes.FileRuntimeExtensionsMap[action.Exec.Kind]
	}

	os.MkdirAll(directory, os.ModePerm)

	path := filepath.Join(directory, filename)

	if err := utils.WriteFile(path, code); err != nil {
		return "", err
	}

	return filename, nil
}

func init() {
	RootCmd.AddCommand(exportCmd)
}
