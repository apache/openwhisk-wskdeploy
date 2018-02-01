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
	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	DEPLOYMENT_HOST      = "sample.deployment.openwhisk.org"
	DEPLOYMENT_AUTH      = "sample-deployment-credential"
	DEPLOYMENT_NAMESPACE = "sample-deployment-namespace"

	MANIFEST_HOST      = "sample.manifest.openwhisk.org"
	MANIFEST_AUTH      = "sample-manifest-credential"
	MANIFEST_NAMESPACE = "sample-manifest-namespace"

	CLI_HOST      = "sample.cli.openwhisk.org"
	CLI_AUTH      = "sample-cli-credential"
	CLI_NAMESPACE = "sample-cli-namespace"

	WSKPROPS_HOST      = "openwhisk.ng.bluemix.net"
	WSKPROPS_AUTH      = "a4f8c502:123zO3xZCLrMN6v2BKK"
	WSKPROPS_NAMESPACE = "guest"

	WSKPROPS_KEY  = "test_key_file"
	WSKPROPS_CERT = "test_cert_file"
)

func initializeFlags() {
	utils.Flags.Auth = ""
	utils.Flags.Namespace = ""
	utils.Flags.ApiHost = ""
	utils.Flags.Key = ""
	utils.Flags.Cert = ""
}

func TestNewWhiskConfig(t *testing.T) {
	propPath := ""
	manifestPath := ""
	deploymentPath := ""
	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	if err == nil {
		pi := whisk.PropertiesImp{
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
	assert.Equal(t, CLI_HOST, config.Host, "Failed to get host name from wskdeploy command line")
	assert.Equal(t, CLI_AUTH, config.AuthToken, "Failed to get auth token from wskdeploy command line")
	assert.Equal(t, CLI_NAMESPACE, config.Namespace, "Failed to get namespace from wskdeploy command line")
	assert.True(t, config.Insecure, "Config should set insecure to true")

	utils.Flags.Key = WSKPROPS_KEY
	utils.Flags.Cert = WSKPROPS_CERT
	config, err = NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from wskdeploy command line")
	assert.Equal(t, CLI_HOST, config.Host, "Failed to get host name from wskdeploy command line")
	assert.Equal(t, CLI_AUTH, config.AuthToken, "Failed to get auth token from wskdeploy command line")
	assert.Equal(t, CLI_NAMESPACE, config.Namespace, "Failed to get namespace from wskdeploy command line")
	assert.Equal(t, WSKPROPS_KEY, config.Key, "Failed to get key file from wskdeploy command line")
	assert.Equal(t, WSKPROPS_CERT, config.Cert, "Failed to get cert file from wskdeploy command line")
	assert.False(t, config.Insecure, "Config should set insecure to false")

	initializeFlags()
}

func TestNewWhiskConfigDeploymentFile(t *testing.T) {
	propPath := ""
	manifestPath := ""
	deploymentPath := "../tests/dat/deployment_validate_credentials.yaml"
	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from deployment file")
	assert.Equal(t, DEPLOYMENT_HOST, config.Host, "Failed to get host name from deployment file")
	assert.Equal(t, DEPLOYMENT_AUTH, config.AuthToken, "Failed to get auth token from deployment file")
	assert.Equal(t, DEPLOYMENT_NAMESPACE, config.Namespace, "Failed to get namespace from deployment file")
	assert.True(t, config.Insecure, "Config should set insecure to true")
}

func TestNewWhiskConfigManifestFile(t *testing.T) {
	propPath := ""
	manifestPath := "../tests/dat/manifest_validate_credentials.yaml"
	deploymentPath := ""
	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from manifest file")
	assert.Equal(t, MANIFEST_HOST, config.Host, "Failed to get host name from manifest file")
	assert.Equal(t, MANIFEST_AUTH, config.AuthToken, "Failed to get auth token from manifest file")
	assert.Equal(t, MANIFEST_NAMESPACE, config.Namespace, "Failed to get namespace from manifest file")
	assert.True(t, config.Insecure, "Config should set insecure to true")
}

func TestNewWhiskConfigWithWskProps(t *testing.T) {
	propPath := "../tests/dat/wskprops"
	manifestPath := ""
	deploymentPath := ""
	config, err := NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from wskprops")
	assert.Equal(t, WSKPROPS_HOST, config.Host, "Failed to get host name from wskprops")
	assert.Equal(t, WSKPROPS_AUTH, config.AuthToken, "Failed to get auth token from wskprops")
	assert.Equal(t, WSKPROPS_NAMESPACE, config.Namespace, "Failed to get namespace from wskprops")
	assert.Equal(t, WSKPROPS_KEY, config.Key, "Failed to get key file from wskprops")
	assert.Equal(t, WSKPROPS_CERT, config.Cert, "Failed to get cert file from wskprops")
	assert.False(t, config.Insecure, "Config should set insecure to false")

	propPath = "../tests/dat/wskpropsnokeycert"
	config, err = NewWhiskConfig(propPath, deploymentPath, manifestPath, false)
	assert.Nil(t, err, "Failed to read credentials from wskprops")
	assert.Equal(t, WSKPROPS_HOST, config.Host, "Failed to get host name from wskprops")
	assert.Equal(t, WSKPROPS_AUTH, config.AuthToken, "Failed to get auth token from wskprops")
	assert.Equal(t, WSKPROPS_NAMESPACE, config.Namespace, "Failed to get namespace from wskprops")
	assert.Empty(t, config.Key, "Failed to get key file from wskprops")
	assert.Empty(t, config.Cert, "Failed to get cert file from wskprops")
	assert.True(t, config.Insecure, "Config should set insecure to true")
}

// TODO(#693) add the following test
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
	assert.True(t, config.Insecure, "Config should set insecure to true")

	initializeFlags()
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
	assert.True(t, config.Insecure, "Config should set insecure to true")

	initializeFlags()
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
	assert.True(t, config.Insecure, "Config should set insecure to true")

	initializeFlags()
}

func TestNewWhiskConfigWithCLIAndWskProps(t *testing.T) {
	propPath := "../tests/dat/wskpropsnokeycert"
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
	assert.True(t, config.Insecure, "Config should set insecure to true")

	initializeFlags()
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
	assert.True(t, config.Insecure, "Config should set insecure to true")
}
