// +build unit

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
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var sd *ServiceDeployer
var dr *DeploymentReader
var deployment_file = "../tests/usecases/github/deployment.yaml"
var manifest_file = "../tests/usecases/github/manifest.yaml"

const (
	// local error messages
	TEST_ERROR_DEPLOYMENT_PARSE_FAILURE        = "Deployment [%s]: Failed to parse."
	TEST_ERROR_DEPLOYMENT_BIND_TRIGGER_FAILURE = "Deployment [%s]: Failed to bind Trigger."
	TEST_ERROR_DEPLOYMENT_FIND_PROJECT         = "Deployment [%s]: Failed to find Project [%s]."
	TEST_ERROR_DEPLOYMENT_FIND_PACKAGES        = "Deployment [%s]: Failed to find Packages for project [%s]."
	TEST_ERROR_DEPLOYMENT_FIND_PACKAGE         = "Deployment [%s]: Failed to find Package [%s]."
	TEST_ERROR_DEPLOYMENT_FIND_TRIGGER         = "Deployment [%s]: Failed to find Trigger [%s]."
)

func init() {
	sd = NewServiceDeployer()
	sd.DeploymentPath = deployment_file
	sd.ManifestPath = manifest_file
	sd.Check()
	dr = NewDeploymentReader(sd)
}

// Check DeploymentReader could handle deployment yaml successfully.
func TestDeploymentReader_HandleYaml(t *testing.T) {
	dr.HandleYaml()
	assert.NotNil(t, dr.DeploymentDescriptor.GetProject().Packages["GitHubCommits"], "DeploymentReader handle deployment yaml failed.")
}

func testLoadAndBindDeploymentYAML(t *testing.T, path string, triggerName string) (*ServiceDeployer, *DeploymentReader) {

	sDeployer := NewServiceDeployer()
	sDeployer.DeploymentPath = path
	// Create empty Trigger for "bind" function to use (as we would have if we read a matching Manifest)
	sDeployer.Deployment.Triggers[triggerName] = new(whisk.Trigger)

	// Create an empty Annotation under the Trigger (again to simulate a manifest)
	// type Annotations []map[string]interface{}
	var tstKeyValue whisk.KeyValue
	tstKeyValue.Key = "bbb"
	tstKeyValue.Value = "old text"
	var tstArray whisk.KeyValueArr

	b2 := append(tstArray, tstKeyValue)

	fmt.Println(tstKeyValue)
	fmt.Println(b2)
	//fmt.Println(tstKeyValArray)

	sDeployer.Deployment.Triggers[triggerName].Annotations = b2

	//parse deployment and bind triggers input and annotations
	dReader := NewDeploymentReader(sDeployer)
	err := dReader.HandleYaml()

	// test load of deployment YAML
	if err != nil {
		assert.Fail(t, fmt.Sprintf(TEST_ERROR_DEPLOYMENT_PARSE_FAILURE, sDeployer.DeploymentPath))
	}

	// Test that we can bind Triggers and Annotations
	err = dReader.bindTriggerInputsAndAnnotations()

	// TODO() XXX remove the following debug stmts
	//dbg_sDeployer:=utils.ConvertMapToJSONString("sDeployer.Deployment", sDeployer.Deployment)
	//fmt.Println(dbg_sDeployer)
	//dbg_dReader:=utils.ConvertMapToJSONString("dReader.DeploymentDescriptor", dReader.DeploymentDescriptor)
	//fmt.Println(dbg_dReader)

	// test load of deployment YAML
	if err != nil {
		fmt.Println(err)
		assert.Fail(t, fmt.Sprintf(TEST_ERROR_DEPLOYMENT_BIND_TRIGGER_FAILURE, sDeployer.DeploymentPath))
	}

	return sDeployer, dReader
}

func TestDeploymentReader_ProjectBindTrigger(t *testing.T) {

	//init variables
	TEST_DATA := "../tests/dat/deployment_deploymentreader_project_bind_trigger.yml"
	TEST_TRIGGER := "locationUpdate"
	TEST_PROJECT := "AppWithTriggerRule"

	// Create a ServiceDeployer with a "dummy" DeploymentPlan (i.e., simulate a fake manifest parse)
	// load the deployment YAMl into dReader.DeploymentDescriptor
	// bind the deployment inputs and annotations to the named Trigger from Deployment to Manifest YAML
	sDeployer, dReader := testLoadAndBindDeploymentYAML(t, TEST_DATA, TEST_TRIGGER)

	// test Project exists with expected name
	projectNameDeploy := dReader.DeploymentDescriptor.GetProject().Name
	if projectNameDeploy != TEST_PROJECT {
		assert.Fail(t, fmt.Sprintf(TEST_ERROR_DEPLOYMENT_FIND_PROJECT, TEST_PROJECT))
	}

	// test that the Project has Packages
	if len(dReader.DeploymentDescriptor.GetProject().Packages) == 0 {
		assert.Fail(t, fmt.Sprintf(TEST_ERROR_DEPLOYMENT_FIND_PACKAGES, TEST_PROJECT))
	}

	trigger := sDeployer.Deployment.Triggers[TEST_TRIGGER]
	for _, param := range trigger.Parameters {
		dbg := utils.ConvertMapToJSONString("param", param)
		fmt.Println(dbg)
		switch param.Key {
		case "name":
			assert.Equal(t, "Bernie", param.Value, "Failed to set inputs")
		case "place":
			assert.Equal(t, "DC", param.Value, "Failed to set inputs")
		default:
			assert.Fail(t, "Failed to get inputs key")

		}
	}
	for _, annos := range trigger.Annotations {
		switch annos.Key {
		case "bbb":
			assert.Equal(t, "this is an annotation", annos.Value, "Failed to set annotations")
		default:
			assert.Fail(t, "Failed to get annotation key")

		}
	}
}

// TODO() XXX this test need to be rewritten perhaps
//func TestDeploymentReader_bindTrigger_packages(t *testing.T) {
//	//init variables
//	sDeployer := NewServiceDeployer()
//	sDeployer.DeploymentPath = "../tests/dat/deployment_deploymentreader_packages_bind_trigger.yml"
//	sDeployer.Deployment.Triggers["locationUpdate"] = new(whisk.Trigger)
//
//	//parse deployment and bind triggers input and annotation
//	dReader := NewDeploymentReader(sDeployer)
//	dReader.HandleYaml()
//	dReader.bindTriggerInputsAndAnnotations()
//
//	trigger := sDeployer.Deployment.Triggers["locationUpdate"]
//	for _, param := range trigger.Parameters {
//		switch param.Key {
//		case "name":
//			assert.Equal(t, "Bernie", param.Value, "Failed to set inputs")
//		case "place":
//			assert.Equal(t, "DC", param.Value, "Failed to set inputs")
//		default:
//			assert.Fail(t, "Failed to get inputs key")
//
//		}
//	}
//	for _, annos := range trigger.Annotations {
//		switch annos.Key {
//		case "bbb":
//			assert.Equal(t, "this is an annotation", annos.Value, "Failed to set annotations")
//		default:
//			assert.Fail(t, "Failed to get annotation key")
//
//		}
//	}
//}

// TODO(749) - rewrite test to remove "package"
func TestDeploymentReader_PackagesBindTrigger(t *testing.T) {
	//init variables
	TEST_DATA := "../tests/dat/deployment_deploymentreader_packages_bind_trigger.yml"
	TEST_PACKAGE := "triggerrule"
	TEST_TRIGGER := "locationUpdate"

	sDeployer, dReader := testLoadAndBindDeploymentYAML(t, TEST_DATA, TEST_TRIGGER)

	// Test if we can find a Package with the name we expect
	packages := dReader.DeploymentDescriptor.Packages

	if packages == nil {
		assert.Fail(t, fmt.Sprintf(TEST_ERROR_DEPLOYMENT_FIND_PACKAGES, TEST_PACKAGE))
	}

	if pkg, ok := packages[TEST_PACKAGE]; ok {

		if trigger, ok := pkg.Triggers[TEST_TRIGGER]; ok {
			dbg := utils.ConvertMapToJSONString("trigger", trigger)
			fmt.Println(dbg)

			for key, value := range trigger.Inputs {

				dbg := utils.ConvertMapToJSONString("value", value)
				fmt.Println(dbg)
				switch key {
				case "name":
					assert.Equal(t, "Bernie", value.Value, "Failed to set inputs")
				case "place":
					assert.Equal(t, "DC", value.Value, "Failed to set inputs")
				default:
					assert.Fail(t, "Failed to get inputs key")

				}
			}
			//for _, annos := range trigger.Annotations {
			//	switch annos.Key {
			//	case "bbb":
			//		assert.Equal(t, "this is an annotation", annos.Value, "Failed to set annotations")
			//	default:
			//		assert.Fail(t, "Failed to get annotation key")
			//
			//	}
			//}

		} else {
			assert.Fail(t, fmt.Sprintf(TEST_ERROR_DEPLOYMENT_FIND_TRIGGER,
				sDeployer.DeploymentPath,
				TEST_TRIGGER))
		}
	} else {
		assert.Fail(t, fmt.Sprintf(TEST_ERROR_DEPLOYMENT_FIND_PACKAGE,
			sDeployer.DeploymentPath,
			TEST_PACKAGE))
	}
}

func TestDeploymentReader_BindAssets_ActionAnnotations(t *testing.T) {
	sDeployer := NewServiceDeployer()
	sDeployer.DeploymentPath = "../tests/dat/deployment_validate_action_annotations.yaml"
	sDeployer.ManifestPath = "../tests/dat/manifest_validate_action_annotations.yaml"

	//parse deployment and bind triggers input and annotation
	dReader := NewDeploymentReader(sDeployer)
	dReader.HandleYaml()
	err := dReader.bindActionInputsAndAnnotations()

	assert.Nil(t, err, "Failed to bind action annotations")

	pkg_name := "packageActionAnnotations"
	pkg := dReader.DeploymentDescriptor.Packages[pkg_name]
	assert.NotNil(t, pkg, "Could not find package with name "+pkg_name)
	action_name := "helloworld"
	action := dReader.DeploymentDescriptor.GetProject().Packages[pkg_name].Actions[action_name]
	assert.NotNil(t, action, "Could not find action with name "+action_name)
	actual_annotations := action.Annotations
	expected_annotations := map[string]interface{}{
		"action_annotation_1": "this is annotation 1",
		"action_annotation_2": "this is annotation 2",
	}
	assert.Equal(t, len(actual_annotations), len(expected_annotations), "Could not find expected number of annotations specified in manifest file")
	eq := reflect.DeepEqual(actual_annotations, expected_annotations)
	assert.True(t, eq, "Expected list of annotations does not match with actual list, expected annotations: %v actual annotations: %v", expected_annotations, actual_annotations)

	pkg_name = "packageActionAnnotationsWithWebAction"
	pkg = dReader.DeploymentDescriptor.Packages[pkg_name]
	assert.NotNil(t, pkg, "Could not find package with name "+pkg_name)
	action = dReader.DeploymentDescriptor.GetProject().Packages[pkg_name].Actions[action_name]
	assert.NotNil(t, action, "Could not find action with name "+action_name)
	actual_annotations = action.Annotations
	expected_annotations["web-export"] = true
	assert.Equal(t, len(actual_annotations), len(expected_annotations), "Could not find expected number of annotations specified in manifest file")
	eq = reflect.DeepEqual(actual_annotations, expected_annotations)
	assert.True(t, eq, "Expected list of annotations does not match with actual list, expected annotations: %v actual annotations: %v", expected_annotations, actual_annotations)

	pkg_name = "packageActionAnnotationsFromDeployment"
	pkg = dReader.DeploymentDescriptor.Packages[pkg_name]
	assert.NotNil(t, pkg, "Could not find package with name "+pkg_name)
	action = dReader.DeploymentDescriptor.GetProject().Packages[pkg_name].Actions[action_name]
	assert.NotNil(t, action, "Could not find action with name "+action_name)
	actual_annotations = action.Annotations
	expected_annotations = map[string]interface{}{
		"action_annotation_1": "this is annotation 1 from deployment",
		"action_annotation_2": "this is annotation 2 from deployment",
	}
	assert.Equal(t, len(actual_annotations), len(expected_annotations), "Could not find expected number of annotations specified in manifest file")
	eq = reflect.DeepEqual(actual_annotations, expected_annotations)
	assert.True(t, eq, "Expected list of annotations does not match with actual list, expected annotations: %v actual annotations: %v", expected_annotations, actual_annotations)
}
