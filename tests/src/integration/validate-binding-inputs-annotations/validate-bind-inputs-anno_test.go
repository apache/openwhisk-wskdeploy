// +build integration

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

package tests

import (
	"github.com/apache/openwhisk-wskdeploy/tests/src/integration/common"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var PATH = "/src/github.com/apache/openwhisk-wskdeploy/tests/src/integration/validate-binding-inputs-annotations/"

func TestBindingInputsAnnotations(t *testing.T) {

	manifestPath := os.Getenv("GOPATH") + PATH + "manifest.yaml"
	deploymentPath := os.Getenv("GOPATH") + PATH + "deployment.yaml"

	wskdeploy := common.NewWskdeploy()
	// verify the inputs & annotations are set
	deploymentObjects, err := wskdeploy.GetDeploymentObjects(manifestPath, deploymentPath)
	assert.NoError(t, err, "Failed to get deployment object.")

	// verify the inputs & annotations of package
	pkgobj := deploymentObjects.Packages["packagebinding"]
	wskpkg := pkgobj.Package
	for _, param := range wskpkg.Parameters {
		switch param.Key {
		case "name":
			assert.Equal(t, "daisy", param.Value, "Failed to get package inputs")
		case "city":
			assert.Equal(t, "Beijing", param.Value, "Failed to get package inputs")
		default:
			assert.Fail(t, "Failed to get package inputs key")
		}
	}
	for _, annos := range wskpkg.Annotations {
		switch annos.Key {
		case "package_annotation_1":
			assert.Equal(t, "this is a package annotation 1", annos.Value, "Failed to get package annotations")
		case "package_annotation_2":
			assert.Equal(t, "this is a package annotation 2", annos.Value, "Failed to get package annotations")
		default:
			assert.Fail(t, "Failed to get package annotation key")
		}
	}

	// verify the inputs & annotations of action
	wskaction := pkgobj.Actions["helloworld"].Action
	for _, param := range wskaction.Parameters {
		switch param.Key {
		case "name":
			assert.Equal(t, "Amy", param.Value, "Failed to get action inputs")
		case "place":
			assert.Equal(t, "Paris", param.Value, "Failed to get action inputs")
		default:
			assert.Fail(t, "Failed to get action inputs key")
		}
	}
	for _, annos := range wskaction.Annotations {
		switch annos.Key {
		case "action_annotation_1":
			assert.Equal(t, "this is an action annotation 1", annos.Value, "Failed to get action annotations")
		case "action_annotation_2":
			assert.Equal(t, "this is an action annotation 2", annos.Value, "Failed to get action annotations")
		default:
			assert.Fail(t, "Failed to get action annotation key")
		}
	}

	// verify the inputs & annotations of trigger
	wsktrigger := deploymentObjects.Triggers["dbtrigger"]
	for _, param := range wsktrigger.Parameters {
		switch param.Key {
		case "dbname":
			assert.Equal(t, "cats", param.Value, "Failed to get trigger inputs")
		case "docid":
			assert.Equal(t, 1234567, param.Value, "Failed to get trigger inputs")
		default:
			assert.Fail(t, "Failed to get trigger inputs key")
		}
	}
	for _, annos := range wsktrigger.Annotations {
		switch annos.Key {
		case "trigger_annotation_1":
			assert.Equal(t, "this is a trigger annotation 1", annos.Value, "Failed to get action annotations")
		case "trigger_annotation_2":
			assert.Equal(t, "this is a trigger annotation 2", annos.Value, "Failed to get action annotations")
		default:
			assert.Fail(t, "Failed to get annotation key")
		}
	}

	// verify the annotations of a package
	pkgobj = deploymentObjects.Packages["packageBindingFromDeployment"]
	wskpkg = pkgobj.Package
	assert.Equal(t, 3, len(wskpkg.Annotations), "Failed to get package annotations")
	for _, annos := range wskpkg.Annotations {
		switch annos.Key {
		case "package_annotation_1":
			assert.Equal(t, "this is a package annotation 1 from deployment", annos.Value, "Failed to get package annotations")
		case "package_annotation_2":
			assert.Equal(t, "this is a package annotation 2", annos.Value, "Failed to get package annotations")
		case "package_annotation_3":
			assert.Equal(t, "this is a package annotation 3 from deployment", annos.Value, "Failed to get package annotations")
		default:
			assert.Fail(t, "Failed to get package annotation key")
		}
	}
	// verify the annotations of action
	wskaction = pkgobj.Actions["helloworld"].Action
	assert.Equal(t, 3, len(wskaction.Annotations), "Failed to get action annotations")
	for _, annos := range wskaction.Annotations {
		switch annos.Key {
		case "action_annotation_1":
			assert.Equal(t, "this is an action annotation 1 from deployment", annos.Value, "Failed to get action annotations")
		case "action_annotation_2":
			assert.Equal(t, "this is an action annotation 2", annos.Value, "Failed to get action annotations")
		case "action_annotation_3":
			assert.Equal(t, "this is an action annotation 3 from deployment", annos.Value, "Failed to get action annotations")
		default:
			assert.Fail(t, "Failed to get action annotation key")
		}
	}
	// verify the inputs & annotations of trigger
	wsktrigger = deploymentObjects.Triggers["dbtrigger1"]
	assert.Equal(t, 3, len(wsktrigger.Annotations), "Failed to get trigger annotations")
	for _, annos := range wsktrigger.Annotations {
		switch annos.Key {
		case "trigger_annotation_1":
			assert.Equal(t, "this is a trigger annotation 1 from deployment", annos.Value, "Failed to get action annotations")
		case "trigger_annotation_2":
			assert.Equal(t, "this is a trigger annotation 2", annos.Value, "Failed to get action annotations")
		case "trigger_annotation_3":
			assert.Equal(t, "this is a trigger annotation 3 from deployment", annos.Value, "Failed to get action annotations")
		default:
			assert.Fail(t, "Failed to get annotation key")
		}
	}

	// testing deploy and undeploy
	_, err = wskdeploy.Deploy(manifestPath, deploymentPath)
	assert.NoError(t, err, "Failed to deploy based on the manifest and deployment files.")
	_, err = wskdeploy.Undeploy(manifestPath, deploymentPath)
	assert.NoError(t, err, "Failed to undeploy based on the manifest and deployment files.")
}
