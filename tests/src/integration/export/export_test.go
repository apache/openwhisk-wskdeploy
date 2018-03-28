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
	"os"
	"testing"

	"github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/common"
	"github.com/stretchr/testify/assert"
)

/* *
 * Please configure BLUEMIX_APIHOST, BLUEMIX_NAMESPACE and BLUEMIX_AUTH on your local machine in order to run this
 * integration test.
 */
func TestExport(t *testing.T) {
	projectName := "EXT_PROJECT"
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

var (
	manifestLib1Path = os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/export/manifest_lib1.yaml"
	manifestLib2Path = os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/export/manifest_lib2.yaml"
	manifestExtPath  = os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/export/manifest_ext.yaml"

	targetManifestFolder = os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/export/tmp/"
	targetManifestPath   = targetManifestFolder + "manifest.yaml"
)
