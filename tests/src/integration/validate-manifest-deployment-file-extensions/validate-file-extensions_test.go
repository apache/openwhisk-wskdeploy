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
	"github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/common"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var projectPath = "/src/github.com/apache/incubator-openwhisk-wskdeploy/tests/src/integration/validate-manifest-deployment-file-extensions/"

func TestYAMLExtension(t *testing.T) {
	manifestPath := os.Getenv("GOPATH") + projectPath + "manifest.yaml"
	deploymentPath := os.Getenv("GOPATH") + projectPath + "deployment.yaml"
	wskdeploy := common.NewWskdeploy()
	_, err := wskdeploy.Deploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifest and deployment files with .yaml extension.")
	_, err = wskdeploy.Undeploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to undeploy based on the manifest and deployment files with .yaml extension.")
}

func TestYMLExtension(t *testing.T) {
	manifestPath := os.Getenv("GOPATH") + projectPath + "manifest.yml"
	deploymentPath := os.Getenv("GOPATH") + projectPath + "deployment.yml"
	wskdeploy := common.NewWskdeploy()
	_, err := wskdeploy.Deploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifest and deployment files with .yml extension.")
	_, err = wskdeploy.Undeploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to undeploy based on the manifest and deployment files with .yml extension.")
}

func TestNonStandardFileNames(t *testing.T) {
	manifestPath := os.Getenv("GOPATH") + projectPath + "not-standard-manifest.yaml"
	deploymentPath := os.Getenv("GOPATH") + projectPath + "not-standard-deployment.yaml"
	wskdeploy := common.NewWskdeploy()
	_, err := wskdeploy.Deploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifest and deployment files with non standard names.")
	_, err = wskdeploy.Undeploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to undeploy based on the manifest and deployment files with non standard names.")
}

func TestRandomFileNames(t *testing.T) {
	manifestPath := os.Getenv("GOPATH") + projectPath + "random-name-1.yaml"
	deploymentPath := os.Getenv("GOPATH") + projectPath + "random-name-2.yaml"
	wskdeploy := common.NewWskdeploy()
	_, err := wskdeploy.Deploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifest and deployment files with random names.")
	_, err = wskdeploy.Undeploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to undeploy based on the manifest and deployment files with random names.")
}

func TestYAMLManifestWithYMLDeployment(t *testing.T) {
	manifestPath := os.Getenv("GOPATH") + projectPath + "yaml-manifest-with-yml-deployment.yaml"
	deploymentPath := os.Getenv("GOPATH") + projectPath + "yml-deployment-with-yaml-manifest.yml"
	wskdeploy := common.NewWskdeploy()
	_, err := wskdeploy.Deploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifest and deployment files with mix of .yaml and .yml extensions.")
	_, err = wskdeploy.Undeploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to undeploy based on the manifest and deployment files with mix of .yaml and .yml extension.")
}

func TestYMLManifestWithYAMLDeployment(t *testing.T) {
	manifestPath := os.Getenv("GOPATH") + projectPath + "yml-manifest-with-yaml-deployment.yml"
	deploymentPath := os.Getenv("GOPATH") + projectPath + "yaml-deployment-with-yml-manifest.yaml"
	wskdeploy := common.NewWskdeploy()
	_, err := wskdeploy.Deploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to deploy based on the manifest and deployment files with .yml manifest and .yaml deployment file.")
	_, err = wskdeploy.Undeploy(manifestPath, deploymentPath)
	assert.Equal(t, nil, err, "Failed to undeploy based on the manifest and deployment files with .yml manifest and .yaml deployment file.")
}
