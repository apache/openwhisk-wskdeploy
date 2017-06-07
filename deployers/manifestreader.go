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
	manifest := manifestParser.ParseManifest(dep.ManifestPath)

	return manifest, manifestParser, nil
}

func (reader *ManifestReader) InitRootPackage(manifestParser *parsers.YAMLParser, manifest *parsers.ManifestYAML) error {
	packg, err := manifestParser.ComposePackage(manifest)
	utils.Check(err)
	reader.SetPackage(packg)

	return nil
}

// Wrapper parser to handle yaml dir
func (deployer *ManifestReader) HandleYaml(sdeployer *ServiceDeployer, manifestParser *parsers.YAMLParser, manifest *parsers.ManifestYAML) error {

	deps, err := manifestParser.ComposeDependencies(manifest, deployer.serviceDeployer.ProjectPath)

	actions, aubindings, err := manifestParser.ComposeActions(manifest, deployer.serviceDeployer.ManifestPath)
	utils.Check(err)

	sequences, err := manifestParser.ComposeSequences(deployer.serviceDeployer.ClientConfig.Namespace, manifest)
	utils.Check(err)

	triggers, err := manifestParser.ComposeTriggers(manifest)
	utils.Check(err)

	rules, err := manifestParser.ComposeRules(manifest)
	utils.Check(err)

	err = deployer.SetDependencies(deps)
	utils.Check(err)

	err = deployer.SetActions(actions)
	utils.Check(err)

	err = deployer.SetSequences(sequences)
	utils.Check(err)

	err = deployer.SetTriggers(triggers)
	utils.Check(err)

	err = deployer.SetRules(rules)
	utils.Check(err)

	//only set api if aubindings
	if len(aubindings) != 0 {
		err = deployer.SetApis(sdeployer, aubindings)
	}

	return nil
}

func (reader *ManifestReader) SetDependencies(deps map[string]utils.DependencyRecord) error {
	for depName, dep := range deps {
		if !dep.IsBinding && !reader.IsUndeploy {
			if _, exists := reader.serviceDeployer.DependencyMaster[depName]; !exists {
				// dependency
				gitReader := utils.NewGitReader(depName, dep)
				err := gitReader.CloneDependency()
				utils.Check(err)
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

func (reader *ManifestReader) SetPackage(pkg *whisk.Package) error {

	dep := reader.serviceDeployer

	dep.mt.Lock()
	defer dep.mt.Unlock()
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
				utils.Check(err)

			} else {
				// Action exists, but references two different sources
				return errors.New("manifestReader. Error: Conflict detected for action named " + existAction.Action.Name + ". Found two locations for source file: " + existAction.Filepath + " and " + manifestAction.Filepath)
			}
		} else {
			// not a new action so to actions in package

			err := reader.checkAction(manifestAction)
			utils.Check(err)
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
	return reader.SetActions(actions)
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

func (reader *ManifestReader) SetApis(deployer *ServiceDeployer, aubs []*utils.ActionExposedURLBinding) error {
	dep := reader.serviceDeployer
	var apis []*whisk.ApiCreateRequest = make([]*whisk.ApiCreateRequest, 0)

	dep.mt.Lock()
	defer dep.mt.Unlock()

	for _, aub := range aubs {
		api := createApiEntity(deployer, aub)
		apis = append(apis, api)
	}

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

// create the api entity according to the action definition and deployer.
func createApiEntity(dp *ServiceDeployer, au *utils.ActionExposedURLBinding) *whisk.ApiCreateRequest {
	sendapi := new(whisk.ApiCreateRequest)
	api := new(whisk.Api)
	//Compose the api
	bindingInfo := strings.Split(au.ExposedUrl, "/")
	api.Namespace = dp.Client.Namespace
	//api.ApiName = ""
	api.GatewayBasePath = bindingInfo[1]
	api.GatewayRelPath = bindingInfo[2]
	api.GatewayMethod = strings.ToUpper(bindingInfo[0])
	api.Id = "API" + ":" + dp.ClientConfig.Namespace + ":" + "/" + api.GatewayBasePath
	//api.GatewayFullPath = ""
	//api.Swagger = ""
	//compose the api action
	api.Action = new(whisk.ApiAction)
	api.Action.Name = au.ActionName
	api.Action.Namespace = dp.ClientConfig.Namespace
	api.Action.BackendMethod = "POST"
	api.Action.BackendUrl = "https://" + dp.ClientConfig.Host + "/api/v1/namespaces/" + dp.ClientConfig.Namespace + "/actions/" + au.ActionName
	api.Action.Auth = dp.Client.Config.AuthToken
	sendapi.ApiDoc = api
	return sendapi
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
