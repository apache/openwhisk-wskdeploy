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

	if err := deployer.UndeployProjectApis(); err != nil {
		return err
	}

	if err := deployer.UndeployProjectRules(); err != nil {
		return err
	}

	if err := deployer.UndeployProjectTriggers(); err != nil {
		return err
	}

	if err := deployer.UndeployProjectActions(); err != nil {
		return err
	}

	if err := deployer.UndeployProjectPackages(); err != nil {
		return err
	}

	if err := deployer.UndeployProjectDependencies(); err != nil {
		return err
	}

	return nil
}

func (deployer *ServiceDeployer) UndeployProjectApis() error {
	return nil
}

func (deployer *ServiceDeployer) UndeployProjectRules() error {
	listOfRules, _, err := deployer.Client.Rules.List(&whisk.RuleListOptions{})
	if err != nil {
		return nil
	}
	for _, rule := range listOfRules {
		if a := rule.Annotations.GetValue(utils.MANAGED); a != nil {
			// decode the JSON blob and retrieve __OW_PROJECT_NAME
			ta := a.(map[string]interface{})
			if ta[utils.OW_PROJECT_NAME] == utils.Flags.ProjectName {
				displayPreprocessingInfo(parsers.YAML_KEY_RULE, rule.Name, false)
				var err error
				var response *http.Response
				err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
					response, err = deployer.Client.Rules.Delete(rule.Name)
					return err
				})
				if err != nil {
					return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_RULE, false)
				}
				displayPostprocessingInfo(parsers.YAML_KEY_RULE, rule.Name, false)
			}

		}

	}
	return nil
}

func (deployer *ServiceDeployer) UndeployProjectTriggers() error {
	listOfTriggers, _, err := deployer.Client.Triggers.List(&whisk.TriggerListOptions{})
	if err != nil {
		return nil
	}
	for _, trigger := range listOfTriggers {
		if a := trigger.Annotations.GetValue(utils.MANAGED); a != nil {
			// decode the JSON blob and retrieve __OW_PROJECT_NAME
			ta := a.(map[string]interface{})
			if ta[utils.OW_PROJECT_NAME] == utils.Flags.ProjectName {
				displayPreprocessingInfo(parsers.YAML_KEY_TRIGGER, trigger.Name, false)
				var err error
				var response *http.Response
				err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
					_, response, err = deployer.Client.Triggers.Delete(trigger.Name)
					return err
				})
				if err != nil {
					return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_TRIGGER, false)
				}
				displayPostprocessingInfo(parsers.YAML_KEY_TRIGGER, trigger.Name, false)
			}

		}

	}
	return nil
}

func (deployer *ServiceDeployer) UndeployProjectActions() error {
	listOfPackages, _, err := deployer.Client.Packages.List(&whisk.PackageListOptions{})
	if err != nil {
		return err
	}
	for _, pkg := range listOfPackages {
		if a := pkg.Annotations.GetValue(utils.MANAGED); a != nil {
			// decode the JSON blob and retrieve __OW_PROJECT_NAME
			ta := a.(map[string]interface{})
			if ta[utils.OW_PROJECT_NAME] == utils.Flags.ProjectName {
				actions, _, err := deployer.Client.Actions.List(pkg.Name, &whisk.ActionListOptions{})
				if err != nil {
					return err
				}
				for _, action := range actions {
					if aa := action.Annotations.GetValue(utils.MANAGED); aa != nil {
						taa := aa.(map[string]interface{})
						if taa[utils.OW_PROJECT_NAME] == utils.Flags.ProjectName {
							displayPreprocessingInfo(parsers.YAML_KEY_ACTION, action.Name, false)
							var err error
							var response *http.Response
							err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
								response, err = deployer.Client.Actions.Delete(pkg.Name + parsers.PATH_SEPARATOR + action.Name)
								return err
							})
							if err != nil {
								return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_ACTION, false)
							}
							displayPostprocessingInfo(parsers.YAML_KEY_ACTION, action.Name, false)
						}
					}
				}
			}
		}

	}
	return nil
}

func (deployer *ServiceDeployer) UndeployProjectPackages() error {
	listOfPackages, _, err := deployer.Client.Packages.List(&whisk.PackageListOptions{})
	if err != nil {
		return nil
	}
	for _, pkg := range listOfPackages {
		if a := pkg.Annotations.GetValue(utils.MANAGED); a != nil {
			// decode the JSON blob and retrieve __OW_PROJECT_NAME
			ta := a.(map[string]interface{})
			if ta[utils.OW_PROJECT_NAME] == utils.Flags.ProjectName {
				displayPreprocessingInfo(parsers.YAML_KEY_PACKAGE, pkg.Name, false)
				var err error
				var response *http.Response
				err = retry(DEFAULT_ATTEMPTS, DEFAULT_INTERVAL, func() error {
					response, err = deployer.Client.Packages.Delete(pkg.Name)
					return err
				})
				if err != nil {
					return createWhiskClientError(err.(*whisk.WskError), response, parsers.YAML_KEY_PACKAGE, false)
				}
				displayPostprocessingInfo(parsers.YAML_KEY_PACKAGE, pkg.Name, false)
			}

		}

	}
	return nil
}
func (deployer *ServiceDeployer) UndeployProjectDependencies() error {
	return nil
}
