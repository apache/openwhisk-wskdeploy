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
	"fmt"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"log"
)

type DeploymentReader struct {
	serviceDeployer      *ServiceDeployer
	DeploymentDescriptor *parsers.DeploymentYAML
}

func NewDeploymentReader(serviceDeployer *ServiceDeployer) *DeploymentReader {
	var dep DeploymentReader
	dep.serviceDeployer = serviceDeployer

	return &dep
}

// Wrapper parser to handle yaml dir
func (reader *DeploymentReader) HandleYaml() error {

	dep := reader.serviceDeployer

	deploymentParser := parsers.NewYAMLParser()
	deployment := deploymentParser.ParseDeployment(dep.DeploymentPath)

	reader.DeploymentDescriptor = deployment

	return nil
}

// Update entities with deployment settings
func (reader *DeploymentReader) BindAssets() error {

	reader.bindPackageInputsAndAnnotations()
	reader.bindActionInputsAndAnnotations()
	reader.bindTriggerInputsAndAnnotations()

	return nil

}

func (reader *DeploymentReader) bindPackageInputsAndAnnotations() {

	packMap := make(map[string]parsers.Package)

	if reader.DeploymentDescriptor.Application.Packages == nil {
		// a single package is specified in deployment YAML file with "package" key
		packMap[reader.DeploymentDescriptor.Application.Package.Packagename] = reader.DeploymentDescriptor.Application.Package
		log.Println("WARNING: The package YAML key in deployment file will soon be deprecated. Please use packages instead as described in specifications.")
	} else {
		for packName, depPacks := range reader.DeploymentDescriptor.Application.Packages {
			depPacks.Packagename = packName
			packMap[packName] = depPacks
		}
	}

	for packName, pack := range packMap {

		serviceDeployPack := reader.serviceDeployer.Deployment.Packages[packName]

		if serviceDeployPack == nil {
			log.Println("Package name in deployment file " + packName + " does not match with manifest file.")
			break
		}

		keyValArr := make(whisk.KeyValueArr, 0)

		if len(pack.Inputs) > 0 {
			for name, input := range pack.Inputs {
				var keyVal whisk.KeyValue

				keyVal.Key = name

				keyVal.Value = utils.GetEnvVar(input)

				keyValArr = append(keyValArr, keyVal)
			}

			depParams := make(map[string]whisk.KeyValue)
			for _, kv := range keyValArr {
				depParams[kv.Key] = kv
			}

			for _, keyVal := range serviceDeployPack.Package.Parameters {
				if _, exists := depParams[keyVal.Key]; !exists {
					keyValArr = append(keyValArr, keyVal)
				}
			}

			serviceDeployPack.Package.Parameters = keyValArr
		}

		keyValArr = keyValArr[:0]

		if len(pack.Annotations) > 0 {
			for name, input := range pack.Annotations {
				var keyVal whisk.KeyValue

				keyVal.Key = name
				keyVal.Value = utils.GetEnvVar(input)

				keyValArr = append(keyValArr, keyVal)
			}

			serviceDeployPack.Package.Annotations = keyValArr
		}

	}
}

func (reader *DeploymentReader) bindActionInputsAndAnnotations() {

	packMap := make(map[string]parsers.Package, 1)

	if reader.DeploymentDescriptor.Application.Packages == nil {
		// a single package is specified in deployment YAML file with "package" key
		packMap[reader.DeploymentDescriptor.Application.Package.Packagename] = reader.DeploymentDescriptor.Application.Package
	} else {
		for packName, depPacks := range reader.DeploymentDescriptor.Application.Packages {
			depPacks.Packagename = packName
			packMap[packName] = depPacks
		}
	}

	for packName, pack := range packMap {

		serviceDeployPack := reader.serviceDeployer.Deployment.Packages[packName]

		if serviceDeployPack == nil {
			break
		}

		for actionName, action := range pack.Actions {

			keyValArr := make(whisk.KeyValueArr, 0)

			if len(action.Inputs) > 0 {
				for name, input := range action.Inputs {
					var keyVal whisk.KeyValue

					keyVal.Key = name

					keyVal.Value = utils.GetEnvVar(input)

					keyValArr = append(keyValArr, keyVal)
				}

				if wskAction, exists := serviceDeployPack.Actions[actionName]; exists {
					depParams := make(map[string]whisk.KeyValue)
					for _, kv := range keyValArr {
						depParams[kv.Key] = kv
					}

					for _, keyVal := range wskAction.Action.Parameters {
						if _, exists := depParams[keyVal.Key]; !exists {
							keyValArr = append(keyValArr, keyVal)
						}
					}
					wskAction.Action.Parameters = keyValArr
				}
			}

			keyValArr = keyValArr[:0]

			if len(action.Annotations) > 0 {
				for name, input := range action.Annotations {
					var keyVal whisk.KeyValue

					keyVal.Key = name
					keyVal.Value = input

					keyValArr = append(keyValArr, keyVal)
				}

				if wskAction, exists := serviceDeployPack.Actions[actionName]; exists {
					wskAction.Action.Annotations = keyValArr
				}
			}
		}

	}
}

func (reader *DeploymentReader) bindTriggerInputsAndAnnotations() {

	packMap := make(map[string]parsers.Package)

	if reader.DeploymentDescriptor.Application.Packages == nil {
		packMap[reader.DeploymentDescriptor.Application.Package.Packagename] = reader.DeploymentDescriptor.Application.Package
	} else {
		for packName, depPacks := range reader.DeploymentDescriptor.Application.Packages {
			depPacks.Packagename = packName
			packMap[packName] = depPacks
		}
	}

	for _, pack := range packMap {

		serviceDeployment := reader.serviceDeployer.Deployment

		for triggerName, trigger := range pack.Triggers {

			keyValArr := make(whisk.KeyValueArr, 0)

			if len(trigger.Inputs) > 0 {
				for name, input := range trigger.Inputs {
					var keyVal whisk.KeyValue

					keyVal.Key = name
					keyVal.Value = utils.GetEnvVar(input)

					keyValArr = append(keyValArr, keyVal)
				}

				if wskTrigger, exists := serviceDeployment.Triggers[triggerName]; exists {

					depParams := make(map[string]whisk.KeyValue)
					for _, kv := range keyValArr {
						depParams[kv.Key] = kv
					}

					for _, keyVal := range wskTrigger.Parameters {
						fmt.Println("Checking key " + keyVal.Key)
						if _, exists := depParams[keyVal.Key]; !exists {
							keyValArr = append(keyValArr, keyVal)
						}
					}
					wskTrigger.Parameters = keyValArr
				}
			}

			keyValArr = keyValArr[:0]

			if len(trigger.Annotations) > 0 {
				for name, input := range trigger.Annotations {
					var keyVal whisk.KeyValue

					keyVal.Key = name
					keyVal.Value = input

					keyValArr = append(keyValArr, keyVal)
				}

				if wskTrigger, exists := serviceDeployment.Triggers[triggerName]; exists {
					wskTrigger.Annotations = keyValArr
				}
			}
		}

	}
}
