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
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-client-go/whisk"
)

const (
	DEPLOYMENT_HOST = "sample.deployment.openwhisk.org"
	DEPLOYMENT_AUTH = "sample-deployment-credential"
	DEPLOYMENT_NAMESPACE = "sample-deployment-namespace"

	MANIFEST_HOST = "sample.manifest.openwhisk.org"
	MANIFEST_AUTH = "sample-manifest-credential"
	MANIFEST_NAMESPACE = "sample-manifest-namespace"

	CLI_HOST = "sample.cli.openwhisk.org"
	CLI_AUTH = "sample-cli-credential"
	CLI_NAMESPACE = "sample-cli-namespace"

	WSKPROPS_HOST = "openwhisk.ng.bluemix.net"
	WSKPROPS_AUTH = "a4f8c502:123zO3xZCLrMN6v2BKK"
	WSKPROPS_NAMESPACE = "guest"
)

func TestNewWhiskConfig(t *testing.T) {
	propPath := ""
	manifestPath := ""
	deploymentPath := ""
	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	if err == nil {
		pi := whisk.PropertiesImp {
			OsPackage: whisk.OSPackageImp{},
		}
		wskprops, err := whisk.GetDefaultWskProp(pi)
		if err == nil {
			assert.Equal(t, config.Namespace, wskprops.Namespace, "")
			assert.Equal(t, config.Host, wskprops.APIHost, "")
			assert.Equal(t, config.AuthToken, wskprops.AuthKey, "")
		}
		return
	}
	assert.NotNil(t, err, "Failed to produce an error when credentials could not be retrieved")
}

func TestNewWhiskConfigCommandLine(t *testing.T) {
	propPath := ""
	manifestPath := ""
	deploymentPath := ""
	utils.Flags.ApiHost = CLI_HOST
	utils.Flags.Auth = CLI_AUTH
	utils.Flags.Namespace = CLI_NAMESPACE

	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from wskdeploy command line")
	assert.Equal(t, config.Host, CLI_HOST, "Failed to get host name from wskdeploy command line")
	assert.Equal(t, config.AuthToken, CLI_AUTH, "Failed to get auth token from wskdeploy command line")
	assert.Equal(t, config.Namespace, CLI_NAMESPACE, "Failed to get namespace from wskdeploy command line")

	utils.Flags.Auth = ""
	utils.Flags.Namespace = ""
	utils.Flags.ApiHost = ""
}

func TestNewWhiskConfigDeploymentFile(t *testing.T) {
	propPath := ""
	manifestPath := ""
	deploymentPath := "../tests/dat/deployment_validate_credentials.yaml"
	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from deployment file")
	assert.Equal(t, config.Host, DEPLOYMENT_HOST, "Failed to get host name from deployment file")
	assert.Equal(t, config.AuthToken, DEPLOYMENT_AUTH, "Failed to get auth token from deployment file")
	assert.Equal(t, config.Namespace, DEPLOYMENT_NAMESPACE, "Failed to get namespace from deployment file")
}

func TestNewWhiskConfigManifestFile(t *testing.T) {
	propPath := ""
	manifestPath := "../tests/dat/manifest_validate_credentials.yaml"
	deploymentPath := ""
	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from manifest file")
	assert.Equal(t, config.Host, MANIFEST_HOST, "Failed to get host name from manifest file")
	assert.Equal(t, config.AuthToken, MANIFEST_AUTH, "Failed to get auth token from manifest file")
	assert.Equal(t, config.Namespace, MANIFEST_NAMESPACE, "Failed to get namespace from manifest file")
}


func TestNewWhiskConfigWithWskProps(t *testing.T) {
	propPath := "../tests/dat/wskprops"
	manifestPath := ""
	deploymentPath := ""
	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from wskprops")
	assert.Equal(t, config.Host, WSKPROPS_HOST, "Failed to get host name from wskprops")
	assert.Equal(t, config.AuthToken, WSKPROPS_AUTH, "Failed to get auth token from wskprops")
	assert.Equal(t, config.Namespace, WSKPROPS_NAMESPACE, "Failed to get namespace from wskprops")
}

// (TODO) add the following test
/*func TestNewWhiskConfigInteractiveMode(t *testing.T) {
	propPath := ""
	manifestPath := ""
	deploymentPath := ""
	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, true)
	assert.Nil(t, err, "Failed to read credentials in interactive mode")
}*/

func TestNewWhiskConfigWithCLIDeploymentAndManifestFile(t *testing.T) {
	propPath := ""
	manifestPath := "../tests/dat/manifest_validate_credentials.yaml"
	deploymentPath := "../tests/dat/deployment_validate_credentials.yaml"

	utils.Flags.ApiHost = CLI_HOST
	utils.Flags.Auth = CLI_AUTH
	utils.Flags.Namespace = CLI_NAMESPACE

	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from CLI or deployment or manifest file")
	assert.Equal(t, config.Host, CLI_HOST, "Failed to get host name from wskdeploy CLI")
	assert.Equal(t, config.AuthToken, CLI_AUTH, "Failed to get auth token from wskdeploy CLI")
	assert.Equal(t, config.Namespace, CLI_NAMESPACE, "Failed to get namespace from wskdeploy CLI")

	utils.Flags.Auth = ""
	utils.Flags.Namespace = ""
	utils.Flags.ApiHost = ""
}


func TestNewWhiskConfigWithCLIAndDeployment(t *testing.T) {
	propPath := ""
	manifestPath := "../tests/dat/deployment_validate_credentials.yaml"
	deploymentPath := ""

	utils.Flags.ApiHost = CLI_HOST
	utils.Flags.Auth = CLI_AUTH
	utils.Flags.Namespace = CLI_NAMESPACE

	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from wskdeploy CLI")
	assert.Equal(t, config.Host, CLI_HOST, "Failed to get host name from wskdeploy CLI")
	assert.Equal(t, config.AuthToken, CLI_AUTH, "Failed to get auth token from wskdeploy CLI")
	assert.Equal(t, config.Namespace, CLI_NAMESPACE, "Failed to get namespace from wskdeploy CLI")

	utils.Flags.Auth = ""
	utils.Flags.Namespace = ""
	utils.Flags.ApiHost = ""
}

func TestNewWhiskConfigWithCLIAndManifest(t *testing.T) {
	propPath := ""
	manifestPath := "../tests/dat/manifest_validate_credentials.yaml"
	deploymentPath := ""

	utils.Flags.ApiHost = CLI_HOST
	utils.Flags.Auth = CLI_AUTH
	utils.Flags.Namespace = CLI_NAMESPACE

	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from manifest file")
	assert.Equal(t, config.Host, CLI_HOST, "Failed to get host name from wskdeploy CLI")
	assert.Equal(t, config.AuthToken, CLI_AUTH, "Failed to get auth token from wskdeploy CLI")
	assert.Equal(t, config.Namespace, CLI_NAMESPACE, "Failed to get namespace from wskdeploy CLI")

	utils.Flags.Auth = ""
	utils.Flags.Namespace = ""
	utils.Flags.ApiHost = ""
}

func TestNewWhiskConfigWithCLIAndWskProps(t *testing.T) {
	propPath := "../tests/dat/wskprops"
	manifestPath := ""
	deploymentPath := ""

	utils.Flags.ApiHost = CLI_HOST
	utils.Flags.Auth = CLI_AUTH
	utils.Flags.Namespace = CLI_NAMESPACE
	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from wskdeploy command line")
	assert.Equal(t, config.Host, CLI_HOST, "Failed to get host name from wskdeploy command line")
	assert.Equal(t, config.AuthToken, CLI_AUTH, "Failed to get auth token from wskdeploy command line")
	assert.Equal(t, config.Namespace, CLI_NAMESPACE, "Failed to get namespace from wskdeploy command line")

	utils.Flags.Auth = ""
	utils.Flags.Namespace = ""
	utils.Flags.ApiHost = ""
}

func TestNewWhiskConfigWithDeploymentAndManifestFile(t *testing.T) {
	propPath := ""
	manifestPath := "../tests/dat/manifest_validate_credentials.yaml"
	deploymentPath := "../tests/dat/deployment_validate_credentials.yaml"
	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from manifest or deployment file")
	assert.Equal(t, config.Host, DEPLOYMENT_HOST, "Failed to get host name from deployment file")
	assert.Equal(t, config.AuthToken, DEPLOYMENT_AUTH, "Failed to get auth token from deployment file")
	assert.Equal(t, config.Namespace, DEPLOYMENT_NAMESPACE, "Failed to get namespace from deployment file")
}