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

package flagstests

import (
	"github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/common"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var wskprops = common.GetWskprops()

// support only projectpath flag
func RunTestSupportProjectPath(t *testing.T) {
	wskdeploy := common.NewWskdeploy()
	projectPath := os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/flagstests"
	_, err := wskdeploy.DeployProjectPathOnly(projectPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the projectpath")
    _, err = wskdeploy.UndeployProjectPathOnly(projectPath)
    assert.Equal(t, nil, err, "Failed to undeploy based on the projectpath")
}

// support only projectpath with trailing slash
func RunTestSupportProjectPathTrailingSlash(t *testing.T) {
	wskdeploy := common.NewWskdeploy()
	projectPath := os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/flagstests" + "/"
	_, err := wskdeploy.DeployProjectPathOnly(projectPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the projectpath")
    _, err = wskdeploy.UndeployProjectPathOnly(projectPath)
    assert.Equal(t, nil, err, "Failed to undeploy based on the projectpath")
}

// only a yaml manifest
func RunTestSupportManifestYamlPath(t *testing.T) {
	wskdeploy := common.NewWskdeploy()
	manifestPath := os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/flagstests/manifest.yaml"
	_, err := wskdeploy.DeployManifestPathOnly(manifestPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifestpath")
    _, err = wskdeploy.UndeployManifestPathOnly(manifestPath)
    assert.Equal(t, nil, err, "Failed to undeploy based on the projectpath")
}

// only a yml manifest
func RunTestSupportManifestYmlPath(t *testing.T) {
	wskdeploy := common.NewWskdeploy()
	manifestPath := os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/flagstests/manifest.yml"
	_, err := wskdeploy.DeployManifestPathOnly(manifestPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifestpath")
    _, err = wskdeploy.UndeployManifestPathOnly(manifestPath)
    assert.Equal(t, nil, err, "Failed to undeploy based on the projectpath")
}

// manifest yaml and deployment yaml
func RunTestSupportManifestYamlDeployment(t *testing.T) {
	wskdeploy := common.NewWskdeploy()
	manifestPath := os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/flagstests/manifest.yaml"
	deploymentPath := os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/flagstests/deployment.yml"
	_, err := wskdeploy.Deploy(manifestPath,deploymentPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifestpath and deploymentpath.")
    _, err = wskdeploy.Undeploy(manifestPath, deploymentPath)
    assert.Equal(t, nil, err, "Failed to undeploy based on the manifest and deployment files.")
}

// manifest yml and deployment yaml
func RunTestSupportManifestYmlDeployment(t *testing.T) {
	wskdeploy := common.NewWskdeploy()
	manifestPath := os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/flagstests/manifest.yml"
	deploymentPath := os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/flagstests/deployment.yml"
	_, err := wskdeploy.Deploy(manifestPath,deploymentPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifestpath and deploymentpath.")
    _, err = wskdeploy.Undeploy(manifestPath, deploymentPath)
    assert.Equal(t, nil, err, "Failed to undeploy based on the manifest and deployment files.")
}
