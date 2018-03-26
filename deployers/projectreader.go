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

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
)

func (deployer *ServiceDeployer) UnDeployProjectAssets() error {

	if err := deployer.SetProjectPackages(); err != nil {
		return err
	}

	if err := deployer.SetProjectActionsAndSequences(); err != nil {
		return err
	}

	if err := deployer.SetProjectTriggers(); err != nil {
		return err
	}

	if err := deployer.SetProjectRules(); err != nil {
		return err
	}

	if err := deployer.SetProjectApis(); err != nil {
		return err
	}

	if err := deployer.SetProjectDependencies(); err != nil {
		return err
	}

	if utils.Flags.Preview {
		deployer.printDeploymentAssets(deployer.Deployment)
		return nil
	}

	return deployer.unDeployAssets(deployer.Deployment)

	return nil
}

func (deployer *ServiceDeployer) isManagedEntity(a interface{}) bool {
	if a != nil {
		ta := a.(map[string]interface{})
		if ta[utils.OW_PROJECT_NAME] == utils.Flags.ProjectName {
			return true
		}
	}
	return false
}

func (deployer *ServiceDeployer) SetProjectPackages() error {
	// retrieve a list of all the packages available under the namespace
	listOfPackages, _, err := deployer.Client.Packages.List(&whisk.PackageListOptions{})
	if err != nil {
		return nil
	}
	for _, pkg := range listOfPackages {
		if deployer.isManagedEntity(pkg.Annotations.GetValue(utils.MANAGED)) {
			var p *whisk.Package
			var response *http.Response
			err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
				p, response, err = deployer.Client.Packages.Get(pkg.Name)
				return err
			})
			if err != nil {
				return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_PACKAGE, false)
			}
			newPack := NewDeploymentPackage()
			newPack.Package = p
			deployer.Deployment.Packages[pkg.Name] = newPack
		}
	}

	return nil
}

func (deployer *ServiceDeployer) SetProjectActionsAndSequences() error {
	for _, pkg := range deployer.Deployment.Packages {
		actions, _, err := deployer.Client.Actions.List(pkg.Package.Name, &whisk.ActionListOptions{})
		if err != nil {
			return err
		}
		for _, action := range actions {
			if deployer.isManagedEntity(action.Annotations.GetValue(utils.MANAGED)) {
				var a *whisk.Action
				var response *http.Response
				err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
					a, response, err = deployer.Client.Actions.Get(pkg.Package.Name+parsers.PATH_SEPARATOR+action.Name, false)
					return err
				})
				if err != nil {
					return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_ACTION, false)
				}
				ar := utils.ActionRecord{Action: a, Packagename: pkg.Package.Name}
				if a.Exec.Kind == parsers.YAML_KEY_SEQUENCE {
					deployer.Deployment.Packages[pkg.Package.Name].Sequences[action.Name] = ar
				} else {
					deployer.Deployment.Packages[pkg.Package.Name].Actions[action.Name] = ar
				}
			}
		}
	}
	return nil
}

func (deployer *ServiceDeployer) SetProjectTriggers() error {
	listOfTriggers, _, err := deployer.Client.Triggers.List(&whisk.TriggerListOptions{})
	if err != nil {
		return nil
	}
	for _, trigger := range listOfTriggers {
		if deployer.isManagedEntity(trigger.Annotations.GetValue(utils.MANAGED)) {
			var t *whisk.Trigger
			var response *http.Response
			err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
				t, response, err = deployer.Client.Triggers.Get(trigger.Name)
				return err
			})
			if err != nil {
				return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_TRIGGER, false)
			}
			deployer.Deployment.Triggers[trigger.Name] = t
		}
	}
	return nil
}

func (deployer *ServiceDeployer) SetProjectRules() error {
	listOfRules, _, err := deployer.Client.Rules.List(&whisk.RuleListOptions{})
	if err != nil {
		return nil
	}
	for _, rule := range listOfRules {
		if deployer.isManagedEntity(rule.Annotations.GetValue(utils.MANAGED)) {
			var r *whisk.Rule
			var response *http.Response
			err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
				r, response, err = deployer.Client.Rules.Get(rule.Name)
				return err
			})
			if err != nil {
				return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_RULE, false)
			}
			deployer.Deployment.Rules[rule.Name] = r
		}
	}
	return nil
}

func (deployer *ServiceDeployer) SetProjectApis() error {
	return nil
}

func (deployer *ServiceDeployer) SetProjectDependencies() error {
	return nil
}
