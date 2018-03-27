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
	"net/http"
	"strings"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
)

func (deployer *ServiceDeployer) UnDeployProjectAssets() error {

	deployer.SetProjectAssets(utils.Flags.ProjectName)

	if utils.Flags.Preview {
		deployer.printDeploymentAssets(deployer.Deployment)
		return nil
	}

	if err := deployer.UndeployProjectDependencies(utils.Flags.ProjectName); err != nil {
		return err
	}

	return deployer.unDeployAssets(deployer.Deployment)

	return nil
}

func (deployer *ServiceDeployer) SetProjectAssets(projectName string) error {

	if err := deployer.SetProjectPackages(projectName); err != nil {
		return err
	}

	if err := deployer.SetPackageActionsAndSequences(projectName); err != nil {
		return err
	}

	if err := deployer.SetProjectTriggers(projectName); err != nil {
		return err
	}

	if err := deployer.SetProjectRules(projectName); err != nil {
		return err
	}

	if err := deployer.SetProjectApis(projectName); err != nil {
		return err
	}

	return nil
}

func (deployer *ServiceDeployer) isManagedEntity(a interface{}, projectName string) bool {
	if a != nil {
		ta := a.(map[string]interface{})
		if ta[utils.OW_PROJECT_NAME] == projectName {
			return true
		}
	}
	return false
}

func (deployer *ServiceDeployer) getPackage(packageName string) (*DeploymentPackage, error) {
	var err error
	var p *whisk.Package
	var response *http.Response
	err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
		p, response, err = deployer.Client.Packages.Get(packageName)
		return err
	})
	if err != nil {
		return nil, createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_PACKAGE, false)
	}
	newPack := NewDeploymentPackage()
	newPack.Package = p
	return newPack, nil

}

func (deployer *ServiceDeployer) SetProjectPackages(projectName string) error {
	// retrieve a list of all the packages available under the namespace
	listOfPackages, _, err := deployer.Client.Packages.List(&whisk.PackageListOptions{})
	if err != nil {
		return nil
	}
	for _, pkg := range listOfPackages {
		if deployer.isManagedEntity(pkg.Annotations.GetValue(utils.MANAGED), projectName) {
			p, err := deployer.getPackage(pkg.Name)
			if err != nil {
				return err
			}
			deployer.Deployment.Packages[pkg.Name] = p
		}
	}

	return nil
}

func (deployer *ServiceDeployer) getPackageActionsAndSequences(packageName string, projectName string) (map[string]utils.ActionRecord, map[string]utils.ActionRecord, error) {
	listOfActions := make(map[string]utils.ActionRecord, 0)
	listOfSequences := make(map[string]utils.ActionRecord, 0)

	actions, _, err := deployer.Client.Actions.List(packageName, &whisk.ActionListOptions{})
	if err != nil {
		return listOfActions, listOfSequences, err
	}
	for _, action := range actions {
		if deployer.isManagedEntity(action.Annotations.GetValue(utils.MANAGED), projectName) {
			var a *whisk.Action
			var response *http.Response
			err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
				a, response, err = deployer.Client.Actions.Get(packageName+parsers.PATH_SEPARATOR+action.Name, false)
				return err
			})
			if err != nil {
				return listOfActions, listOfSequences, createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_ACTION, false)
			}
			ar := utils.ActionRecord{Action: a, Packagename: packageName}
			if a.Exec.Kind == parsers.YAML_KEY_SEQUENCE {
				listOfSequences[action.Name] = ar
			} else {
				listOfActions[action.Name] = ar
			}
		}
	}
	return listOfActions, listOfSequences, err
}

func (deployer *ServiceDeployer) SetPackageActionsAndSequences(projectName string) error {
	for _, pkg := range deployer.Deployment.Packages {
		a, s, err := deployer.getPackageActionsAndSequences(pkg.Package.Name, projectName)
		if err != nil {
			return err
		}
		deployer.Deployment.Packages[pkg.Package.Name].Actions = a
		deployer.Deployment.Packages[pkg.Package.Name].Sequences = s
	}
	return nil
}

func (deployer *ServiceDeployer) getProjectTriggers(projectName string) (map[string]*whisk.Trigger, error) {
	triggers := make(map[string]*whisk.Trigger, 0)
	listOfTriggers, _, err := deployer.Client.Triggers.List(&whisk.TriggerListOptions{})
	if err != nil {
		return triggers, nil
	}
	for _, trigger := range listOfTriggers {
		if deployer.isManagedEntity(trigger.Annotations.GetValue(utils.MANAGED), projectName) {
			var t *whisk.Trigger
			var response *http.Response
			err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
				t, response, err = deployer.Client.Triggers.Get(trigger.Name)
				return err
			})
			if err != nil {
				return triggers, createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_TRIGGER, false)
			}
			triggers[trigger.Name] = t
		}
	}
	return triggers, nil
}

func (deployer *ServiceDeployer) SetProjectTriggers(projectName string) error {
	t, err := deployer.getProjectTriggers(projectName)
	if err != nil {
		return err
	}
	deployer.Deployment.Triggers = t
	return nil
}

func (deployer *ServiceDeployer) getProjectRules(projectName string) (map[string]*whisk.Rule, error) {
	rules := make(map[string]*whisk.Rule, 0)
	listOfRules, _, err := deployer.Client.Rules.List(&whisk.RuleListOptions{})
	if err != nil {
		return rules, nil
	}
	for _, rule := range listOfRules {
		if deployer.isManagedEntity(rule.Annotations.GetValue(utils.MANAGED), projectName) {
			var r *whisk.Rule
			var response *http.Response
			err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
				r, response, err = deployer.Client.Rules.Get(rule.Name)
				return err
			})
			if err != nil {
				return rules, createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_RULE, false)
			}
			rules[rule.Name] = r
		}
	}
	return rules, nil
}

func (deployer *ServiceDeployer) SetProjectRules(projectName string) error {
	r, err := deployer.getProjectRules(projectName)
	if err != nil {
		return err
	}
	deployer.Deployment.Rules = r
	return nil
}

// "whisk-manged" annotation stores package name with namespace such as
// /<namespace>/<package name>
// parse this kind of structure to determine package name
func (deployer *ServiceDeployer) filterPackageName(name string) string {
	s := strings.SplitAfterN(name, "/", 3)
	if len(s) == 3 && len(s[2]) != 0 {
		return s[2]
	}
	return ""
}

// determine if any other package on the server is using the dependent package
func (deployer *ServiceDeployer) isPackageUsedByOtherPackages(projectName string, depPackageName string) bool {
	// retrieve a list of packages on the server
	listOfPackages, _, err := deployer.Client.Packages.List(&whisk.PackageListOptions{})
	if err != nil {
		return false
	}
	for _, pkg := range listOfPackages {
		if a := pkg.Annotations.GetValue(utils.MANAGED); a != nil {
			ta := a.(map[string]interface{})
			// compare project names of the given package and other packages from server
			// we want to skip comparing packages from the same project
			if ta[utils.OW_PROJECT_NAME] != projectName {
				d := a.(map[string]interface{})[utils.OW_PROJECT_DEPS]
				listOfDeps := d.([]interface{})
				// iterate over a list of dependencies of a package
				// to determine whether it has listed the dependent package as its dependency as well
				// in such case, we dont want to undeploy dependent package if its used by any other package
				for _, dep := range listOfDeps {
					name := deployer.filterPackageName(dep.(map[string]interface{})[wski18n.KEY_KEY].(string))
					if name == depPackageName {
						return true
					}

				}

			}

		}
	}
	return false
}

// derive a map of dependent packages using "whisk-managed" annotation
// "whisk-managed" annotation has a list of dependent packages in "projectDeps"
// projectDeps is a list of key value pairs with its own "whisk-managed" annotation
// for a given package, get the list of dependent packages
// for each dependent package, determine whether any other package is using it or not
// if not, collect a list of actions, sequences, triggers, and rules of the dependent package
// and delete them in order following by deleting the package itself
func (deployer *ServiceDeployer) UndeployProjectDependencies(projectName string) error {
	// iterate over each package in a given project
	for _, pkg := range deployer.Deployment.Packages {
		// get the "whisk-managed" annotation
		if a := pkg.Package.Annotations.GetValue(utils.MANAGED); a != nil {
			// read the list of dependencies from "projectDeps"
			d := a.(map[string]interface{})[utils.OW_PROJECT_DEPS]
			listOfDeps := d.([]interface{})
			// iterate over a list of dependencies
			for _, dep := range listOfDeps {
				// dependent package name is in form of "/<namespace>/<package-name>
				// filter it to derive the package name
				name := deployer.filterPackageName(dep.(map[string]interface{})[wski18n.KEY_KEY].(string))
				// undeploy dependent package if its not used by any other package
				if !deployer.isPackageUsedByOtherPackages(projectName, name) {
					// get the *whisk.Package object for the given dependent package
					p, err := deployer.getPackage(name)
					if err != nil {
						return err
					}
					// construct a new DeploymentProject for each dependency
					depProject := NewDeploymentProject()
					depProject.Packages[p.Package.Name] = p
					// Now, get the project name of dependent package so that
					// we can get other entities from that project
					pa := p.Package.Annotations.GetValue(utils.MANAGED)
					depProjectName := (pa.(map[string]interface{})[utils.OW_PROJECT_NAME]).(string)
					// get a list of actions and sequences of a dependent package
					actions, sequences, err := deployer.getPackageActionsAndSequences(p.Package.Name, depProjectName)
					if err != nil {
						return err
					}
					depProject.Packages[p.Package.Name].Actions = actions
					depProject.Packages[p.Package.Name].Sequences = sequences
					// get a list of triggers of a dependent project
					t, err := deployer.getProjectTriggers(depProjectName)
					if err != nil {
						return err
					}
					depProject.Triggers = t
					// get a list of rules of a dependent project
					r, err := deployer.getProjectRules(depProjectName)
					if err != nil {
						return err
					}
					depProject.Rules = r
					// now, we have collected list of actions, sequences, rules, and triggers
					// of a dependent project and ready to undeploy the  whole project
					err = deployer.unDeployAssets(depProject)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (deployer *ServiceDeployer) SetProjectApis(projectName string) error {
	return nil
}
