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

	WHISK_PROPERTY_HOST = "172.17.0.1"
	WHISK_PROPERTY_AUTH = ""
	WHISK_PROPERTY_NAMESPACE = "guest"

	DEFAULT_NAMESPACE = "guest"

)

// this test generates varied results depending on the environment in which its running
// credentials are retrieved (for wskdeploy running in non-interactive mode) in the following precedence order:
// (1) deployment file
// (2) manifest file
// (3) wskdeploy command line
// (4) .wskprops
// (5) whisk.properties
// in this test, credentials could be coming from
// whisk.properties if you have OPENWHISK_HOME set to valid path
// $HOME/.wskprops if you have set openwhisk properties using wsk CLI
// this test works great in Travis as unit tests are run before deploying openwhisk
// but might fail in developer's environment and therefore commenting out
/*func TestNewWhiskConfig(t *testing.T) {
	propPath := ""
	manifestPath := ""
	deploymentPath := ""
	_, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.NotNil(t, err, "Failed to produce an error when credentials could not be retrieved")
}*/

func TestNewWhiskConfigWithOnlyDeploymentFile(t *testing.T) {
	propPath := ""
	manifestPath := ""
	deploymentPath := "../tests/dat/deployment_validate_credentials.yaml"
	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from deployment file")
	assert.Equal(t, config.Host, DEPLOYMENT_HOST, "Failed to get host name from deployment file")
	assert.Equal(t, config.AuthToken, DEPLOYMENT_AUTH, "Failed to get auth token from deployment file")
	assert.Equal(t, config.Namespace, DEPLOYMENT_NAMESPACE, "Failed to get namespace from deployment file")
}

func TestNewWhiskConfigWithOnlyManifestFile(t *testing.T) {
	propPath := ""
	manifestPath := "../tests/dat/manifest_validate_credentials.yaml"
	deploymentPath := ""
	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from manifest file")
	assert.Equal(t, config.Host, MANIFEST_HOST, "Failed to get host name from manifest file")
	assert.Equal(t, config.AuthToken, MANIFEST_AUTH, "Failed to get auth token from manifest file")
	assert.Equal(t, config.Namespace, MANIFEST_NAMESPACE, "Failed to get namespace from manifest file")
}

func TestNewWhiskConfigCommandLine(t *testing.T) {
	propPath := ""
	manifestPath := ""
	deploymentPath := ""
	utils.Flags.ApiHost = CLI_HOST
	utils.Flags.Auth = CLI_AUTH
	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from wskdeploy command line")
	assert.Equal(t, config.Host, CLI_HOST, "Failed to get host name from wskdeploy command line")
	assert.Equal(t, config.AuthToken, CLI_AUTH, "Failed to get auth token from wskdeploy command line")
	assert.Equal(t, config.Namespace, DEFAULT_NAMESPACE, "Failed to get namespace from defaults")

	utils.Flags.Namespace = CLI_NAMESPACE
	config_with_ns, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from wskdeploy command line")
	assert.Equal(t, config_with_ns.Namespace, CLI_NAMESPACE, "Failed to get namespace from wskdeploy command line")

	utils.Flags.Auth = ""
	utils.Flags.Namespace = ""
	utils.Flags.ApiHost = ""
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

func TestNewWhiskConfigWithManifestAndCLI(t *testing.T) {
	propPath := ""
	manifestPath := "../tests/dat/manifest_validate_credentials.yaml"
	deploymentPath := ""

	utils.Flags.ApiHost = CLI_HOST
	utils.Flags.Auth = CLI_AUTH
	utils.Flags.Namespace = CLI_NAMESPACE

	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from manifest file")
	assert.Equal(t, config.Host, MANIFEST_HOST, "Failed to get host name from manifest file")
	assert.Equal(t, config.AuthToken, MANIFEST_AUTH, "Failed to get auth token from manifest file")
	assert.Equal(t, config.Namespace, MANIFEST_NAMESPACE, "Failed to get namespace from manifest file")

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
