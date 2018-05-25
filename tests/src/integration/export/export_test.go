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
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/common"
	"github.com/stretchr/testify/assert"
)

const EXPORT_TEST_PATH = "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/export/"

func TestExport(t *testing.T) {
	projectName := "EXT_PROJECT"

	manifestLib1Path := os.Getenv("GOPATH") + EXPORT_TEST_PATH + "manifest_lib1.yaml"
	manifestLib2Path := os.Getenv("GOPATH") + EXPORT_TEST_PATH + "manifest_lib2.yaml"
	manifestExtPath := os.Getenv("GOPATH") + EXPORT_TEST_PATH + "manifest_ext.yaml"

	targetManifestFolder := os.Getenv("GOPATH") + EXPORT_TEST_PATH + "tmp-" + strconv.Itoa(rand.Intn(1000)) + "/"
	targetManifestPath := targetManifestFolder + "manifest-" + projectName + ".yaml"

	wskdeploy := common.NewWskdeploy()

	_, err := wskdeploy.ManagedDeploymentOnlyManifest(manifestLib1Path)
	assert.Equal(t, nil, err, "Failed to deploy the lib1 manifest file.")

	_, err = wskdeploy.ManagedDeploymentOnlyManifest(manifestLib2Path)
	assert.Equal(t, nil, err, "Failed to deploy the lib2 manifest file.")

	_, err = wskdeploy.ManagedDeploymentOnlyManifest(manifestExtPath)
	assert.Equal(t, nil, err, "Failed to deploy the ext manifest file.")

	_, err = wskdeploy.ExportProject(projectName, targetManifestPath)
	assert.Equal(t, nil, err, "Failed to export project.")

	_, err = os.Stat(targetManifestPath)
	assert.Equal(t, nil, err, "Missing exported manifest file")

	_, err = os.Stat(targetManifestFolder + "dependencies/lib1.yaml")
	assert.Equal(t, nil, err, "Missing exported dependencies lib1 manifest")

	_, err = os.Stat(targetManifestFolder + "dependencies/lib1_package/lib1_greeting1.js")
	assert.Equal(t, nil, err, "Missing exported dependencies lib1 resources")

	_, err = wskdeploy.UndeployManifestPathOnly(manifestExtPath)
	assert.Equal(t, nil, err, "Failed to undeploy the ext.")

	_, err = wskdeploy.UndeployManifestPathOnly(manifestLib2Path)
	assert.Equal(t, nil, err, "Failed to undeploy the lib1.")

	_, err = wskdeploy.UndeployManifestPathOnly(manifestLib1Path)
	assert.Equal(t, nil, err, "Failed to undeploy the lib2.")
}

func SkipTestExportHelloWorld(t *testing.T) {
	projectName := "HELLO_WORLD"
	manifestHelloWorldPath := os.Getenv("GOPATH") + EXPORT_TEST_PATH + "manifest_helloworld.yaml"
	targetManifestFolder := os.Getenv("GOPATH") + EXPORT_TEST_PATH + "tmp-" + strconv.Itoa(rand.Intn(1000)) + "/"
	targetManifestHelloWorldPath := targetManifestFolder + "manifest-" + projectName + ".yaml"

	wskdeploy := common.NewWskdeploy()

	_, err := wskdeploy.ManagedDeploymentManifestAndProject(manifestHelloWorldPath, projectName)
	assert.Equal(t, nil, err, "Failed to deploy manifest file.")

	_, err = wskdeploy.ExportProject(projectName, targetManifestHelloWorldPath)
	assert.Equal(t, nil, err, "Failed to export project.")

	_, err = os.Stat(targetManifestHelloWorldPath)
	assert.Equal(t, nil, err, "Missing exported manifest file")

	_, err = wskdeploy.UndeployManifestPathOnly(manifestHelloWorldPath)
	assert.Equal(t, nil, err, "Failed to undeploy")

	wskprops := common.GetWskpropsFromEnvVars(common.BLUEMIX_APIHOST, common.BLUEMIX_NAMESPACE, common.BLUEMIX_AUTH)
	err = common.ValidateWskprops(wskprops)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Wsk properties are not properly configured, so tests are skipped.")
	} else {
		wskdeploy := common.NewWskdeploy()
		_, err = wskdeploy.ManagedDeploymentManifestAndProject(targetManifestHelloWorldPath, projectName)
		assert.Equal(t, nil, err, "Failed to redeploy exported project.")

		_, err = wskdeploy.UndeployManifestPathOnly(targetManifestHelloWorldPath)
		assert.Equal(t, nil, err, "Failed to undeploy exported project")
	}
}

func TestExport2Pack(t *testing.T) {
	manifest2PackPath := os.Getenv("GOPATH") + EXPORT_TEST_PATH + "manifest_2pack.yaml"
	targetManifestFolder := os.Getenv("GOPATH") + EXPORT_TEST_PATH + "tmp-" + strconv.Itoa(rand.Intn(1000)) + "/"
	target2PackManifestPath := targetManifestFolder + "exported2packmanifest.yaml"
	projectName := "2pack"
	wskdeploy := common.NewWskdeploy()

	_, err := wskdeploy.ManagedDeploymentOnlyManifest(manifest2PackPath)
	assert.Equal(t, nil, err, "Failed to deploy the 2pack manifest file.")

	_, err = wskdeploy.ExportProject(projectName, target2PackManifestPath)
	assert.Equal(t, nil, err, "Failed to export project.")

	_, err = os.Stat(target2PackManifestPath)
	assert.Equal(t, nil, err, "Missing exported manifest file")

	_, err = os.Stat(targetManifestFolder + "package_1/pack1_greeting1.js")
	assert.Equal(t, nil, err, "Missing exported package_1/pack1_greeting1.js")

	_, err = os.Stat(targetManifestFolder + "package_2/pack2_greeting2.js")
	assert.Equal(t, nil, err, "Missing exported package_2/pack2_greeting2.js")

	_, err = wskdeploy.UndeployManifestPathOnly(manifest2PackPath)
	assert.Equal(t, nil, err, "Failed to undeploy")
}
