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
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskderrors"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
)

// Possible sources for config info (e.g., API Host, Auth Key, Namespace)
const (
	SOURCE_WSKPROPS         = ".wskprops"
	SOURCE_WHISK_PROPERTIES = "whisk.properties"
	SOURCE_DEFAULT_VALUE    = "wskdeploy default" // TODO() i18n?
)

var (
	credential        = PropertyValue{}
	namespace         = PropertyValue{}
	apiHost           = PropertyValue{}
	key               = PropertyValue{}
	cert              = PropertyValue{}
	apigwAccessToken  = PropertyValue{}
	additionalHeaders = make(http.Header)
)

type PropertyValue struct {
	Value  string
	Source string
}

var GetPropertyValue = func(prop PropertyValue, newValue string, source string) PropertyValue {
	if len(prop.Value) == 0 && len(newValue) > 0 {
		prop.Value = newValue
		prop.Source = source
	}
	return prop
}

var GetWskPropFromWskprops = func(pi whisk.Properties, proppath string) (*whisk.Wskprops, error) {
	return whisk.GetWskPropFromWskprops(pi, proppath)
}

var GetWskPropFromWhiskProperty = func(pi whisk.Properties) (*whisk.Wskprops, error) {
	return whisk.GetWskPropFromWhiskProperty(pi)
}

var GetCommandLineFlags = func() (string, string, string, string, string, string) {
	return utils.Flags.ApiHost, utils.Flags.Auth, utils.Flags.Namespace, utils.Flags.Key, utils.Flags.Cert, utils.Flags.ApigwAccessToken
}

var CreateNewClient = func(config_input *whisk.Config) (*whisk.Client, error) {
	var netClient = &http.Client{
		Timeout: time.Second * utils.DEFAULT_HTTP_TIMEOUT,
	}
	return whisk.NewClient(netClient, config_input)
}

func AddAdditionalHeader(hdrName string, hdrValue string) {
	additionalHeaders.Add(hdrName, hdrValue)
}

func resetWhiskConfig() {
	credential = PropertyValue{}
	namespace = PropertyValue{}
	apiHost = PropertyValue{}
	key = PropertyValue{}
	cert = PropertyValue{}
	apigwAccessToken = PropertyValue{}
}

func readFromCLI() {
	// read credentials, namespace, API host, key file, cert file, and APIGW access token from command line
	apihost, auth, ns, keyfile, certfile, accessToken := GetCommandLineFlags()
	credential = GetPropertyValue(credential, auth, wski18n.COMMAND_LINE)
	namespace = GetPropertyValue(namespace, ns, wski18n.COMMAND_LINE)
	apiHost = GetPropertyValue(apiHost, apihost, wski18n.COMMAND_LINE)
	key = GetPropertyValue(key, keyfile, wski18n.COMMAND_LINE)
	cert = GetPropertyValue(cert, certfile, wski18n.COMMAND_LINE)
	apigwAccessToken = GetPropertyValue(apigwAccessToken, accessToken, wski18n.COMMAND_LINE)
}

func setWhiskConfig(cred string, ns string, host string, token string, source string) {
	credential = GetPropertyValue(credential, cred, source)
	namespace = GetPropertyValue(namespace, ns, source)
	apiHost = GetPropertyValue(apiHost, host, source)
	apigwAccessToken = GetPropertyValue(apigwAccessToken, token, source)
}

func readFromDeploymentFile(deploymentPath string) {
	if len(credential.Value) == 0 || len(namespace.Value) == 0 || len(apiHost.Value) == 0 {
		if utils.FileExists(deploymentPath) {
			mm := parsers.NewYAMLParser()
			deployment, _ := mm.ParseDeployment(deploymentPath)
			p := deployment.GetProject()
			setWhiskConfig(p.Credential, p.Namespace, p.ApiHost, p.ApigwAccessToken, path.Base(deploymentPath))
		}
	}
}

func readFromManifestFile(manifestPath string) {
	if len(credential.Value) == 0 || len(namespace.Value) == 0 || len(apiHost.Value) == 0 {
		if utils.FileExists(manifestPath) {
			mm := parsers.NewYAMLParser()
			manifest, _ := mm.ParseManifest(manifestPath)
			p := manifest.GetProject()
			setWhiskConfig(p.Credential, p.Namespace, p.ApiHost, p.ApigwAccessToken, path.Base(manifestPath))
		}
	}
}

func readFromWskprops(pi whisk.PropertiesImp, proppath string) {
	// The error raised here can be neglected, because we will handle it in the end of its calling function.
	wskprops, _ := GetWskPropFromWskprops(pi, proppath)
	credential = GetPropertyValue(credential, wskprops.AuthKey, SOURCE_WSKPROPS)
	namespace = GetPropertyValue(namespace, wskprops.Namespace, SOURCE_WSKPROPS)
	apiHost = GetPropertyValue(apiHost, wskprops.APIHost, SOURCE_WSKPROPS)
	key = GetPropertyValue(key, wskprops.Key, SOURCE_WSKPROPS)
	cert = GetPropertyValue(cert, wskprops.Cert, SOURCE_WSKPROPS)
	apigwAccessToken = GetPropertyValue(apigwAccessToken, wskprops.AuthAPIGWKey, SOURCE_WSKPROPS)
}

func readFromWhiskProperty(pi whisk.PropertiesImp) {
	// now, read credentials from whisk.properties but this is only acceptable within Travis
	// whisk.properties will soon be deprecated and should not be used for any production deployment
	whiskproperty, _ := GetWskPropFromWhiskProperty(pi)
	var warnMsg string
	credential = GetPropertyValue(credential, whiskproperty.AuthKey, SOURCE_WHISK_PROPERTIES)
	if credential.Source == SOURCE_WHISK_PROPERTIES {
		warnMsg = wski18n.T(wski18n.ID_WARN_WHISK_PROPS_DEPRECATED,
			map[string]interface{}{wski18n.KEY_KEY: wski18n.AUTH_KEY})
		wskprint.PrintlnOpenWhiskWarning(warnMsg)
	}
	namespace = GetPropertyValue(namespace, whiskproperty.Namespace, SOURCE_WHISK_PROPERTIES)
	if namespace.Source == SOURCE_WHISK_PROPERTIES {
		warnMsg = wski18n.T(wski18n.ID_WARN_WHISK_PROPS_DEPRECATED,
			map[string]interface{}{wski18n.KEY_KEY: parsers.YAML_KEY_NAMESPACE})
		wskprint.PrintlnOpenWhiskWarning(warnMsg)
	}
	apiHost = GetPropertyValue(apiHost, whiskproperty.APIHost, SOURCE_WHISK_PROPERTIES)
	if apiHost.Source == SOURCE_WHISK_PROPERTIES {
		warnMsg = wski18n.T(wski18n.ID_WARN_WHISK_PROPS_DEPRECATED,
			map[string]interface{}{wski18n.KEY_KEY: wski18n.API_HOST})
		wskprint.PrintlnOpenWhiskWarning(warnMsg)
	}
	apigwAccessToken = GetPropertyValue(apigwAccessToken, whiskproperty.AuthAPIGWKey, SOURCE_WHISK_PROPERTIES)
	if apigwAccessToken.Source == SOURCE_WHISK_PROPERTIES {
		warnMsg = wski18n.T(wski18n.ID_WARN_WHISK_PROPS_DEPRECATED,
			map[string]interface{}{wski18n.KEY_KEY: wski18n.APIGW_ACCESS_TOKEN})
		wskprint.PrintlnOpenWhiskWarning(warnMsg)
	}
}

// we are reading openwhisk credentials (apihost, namespace, and auth) in the following precedence order:
// (1) wskdeploy command line `wskdeploy --apihost --namespace --auth`
// (2) deployment file
// (3) manifest file
// (4) .wskprops
// we are following the same precedence order for APIGW_ACCESS_TOKEN
// but as a separate thread as APIGW_ACCESS_TOKEN only needed for APIs
func NewWhiskConfig(proppath string, deploymentPath string, manifestPath string) (*whisk.Config, error) {
	// reset credential, apiHost, namespace, etc to avoid any conflicts as they initialized globally
	resetWhiskConfig()

	// initialize APIGW_ACCESS_TOKEN to "DUMMY TOKEN" for Travis builds
	if strings.ToLower(os.Getenv("TRAVIS")) == "true" {
		apigwAccessToken.Value = "DUMMY TOKEN"
		apigwAccessToken.Source = SOURCE_DEFAULT_VALUE
	}

	// read from command line
	readFromCLI()

	// TODO() i18n
	// Print all flags / values if verbose
	wskprint.PrintlnOpenWhiskVerbose(utils.Flags.Verbose, wski18n.CONFIGURATION+":\n"+utils.Flags.Format())

	// now, read them from deployment file if not found on command line
	readFromDeploymentFile(deploymentPath)

	// read credentials from manifest file as didn't find them on command line and in deployment file
	readFromManifestFile(manifestPath)

	// Third, we need to look up the variables in .wskprops file.
	pi := whisk.PropertiesImp{
		OsPackage: whisk.OSPackageImp{},
	}

	readFromWskprops(pi, proppath)

	// TODO() whisk.properties should be deprecated
	readFromWhiskProperty(pi)

	// set namespace to default namespace if not yet found
	if len(apiHost.Value) != 0 && len(credential.Value) != 0 && len(namespace.Value) == 0 {
		namespace.Value = whisk.DEFAULT_NAMESPACE
		namespace.Source = SOURCE_DEFAULT_VALUE
	}

	mode := true
	if len(cert.Value) != 0 && len(key.Value) != 0 {
		mode = false
	}

	clientConfig = &whisk.Config{
		AuthToken:         credential.Value, //Authtoken
		Namespace:         namespace.Value,  //Namespace
		Host:              apiHost.Value,
		Version:           "v1", // TODO() should not be hardcoded, should warn user of default
		Cert:              cert.Value,
		Key:               key.Value,
		Insecure:          mode, // true if you want to ignore certificate signing
		ApigwAccessToken:  apigwAccessToken.Value,
		AdditionalHeaders: additionalHeaders,
	}

	// validate we have credential, apihost and namespace
	err := validateClientConfig(credential, apiHost, namespace)
	return clientConfig, err
}

func validateClientConfig(credential PropertyValue, apiHost PropertyValue, namespace PropertyValue) error {

	// Display error message for each config value found missing
	if len(credential.Value) == 0 || len(apiHost.Value) == 0 || len(namespace.Value) == 0 {

		var errorMsg string = ""
		if len(credential.Value) == 0 {
			errorMsg = wskderrors.AppendDetailToErrorMessage(
				errorMsg, wski18n.T(wski18n.ID_MSG_CONFIG_MISSING_AUTHKEY), 1)
		}

		if len(apiHost.Value) == 0 {
			errorMsg = wskderrors.AppendDetailToErrorMessage(
				errorMsg, wski18n.T(wski18n.ID_MSG_CONFIG_MISSING_APIHOST), 1)

		}

		if len(namespace.Value) == 0 {
			errorMsg = wskderrors.AppendDetailToErrorMessage(
				errorMsg, wski18n.T(wski18n.ID_MSG_CONFIG_MISSING_NAMESPACE), 1)

		}

		if len(errorMsg) > 0 {
			return wskderrors.NewWhiskClientInvalidConfigError(errorMsg)
		}
	}

	// Show caller what final values we used for credential, apihost and namespace
	stdout := wski18n.T(wski18n.ID_MSG_CONFIG_INFO_APIHOST_X_host_X_source_X,
		map[string]interface{}{wski18n.KEY_HOST: apiHost.Value, wski18n.KEY_SOURCE: apiHost.Source})
	wskprint.PrintOpenWhiskInfo(stdout)

	stdout = wski18n.T(wski18n.ID_MSG_CONFIG_INFO_AUTHKEY_X_source_X,
		map[string]interface{}{wski18n.KEY_SOURCE: credential.Source})
	wskprint.PrintOpenWhiskInfo(stdout)

	stdout = wski18n.T(wski18n.ID_MSG_CONFIG_INFO_NAMESPACE_X_namespace_X_source_X,
		map[string]interface{}{wski18n.KEY_NAMESPACE: namespace.Value, wski18n.KEY_SOURCE: namespace.Source})
	wskprint.PrintOpenWhiskInfo(stdout)

	if len(apigwAccessToken.Value) != 0 {
		stdout = wski18n.T(wski18n.ID_MSG_CONFIG_INFO_APIGE_ACCESS_TOKEN_X_source_X,
			map[string]interface{}{wski18n.KEY_SOURCE: apigwAccessToken.Source})
		wskprint.PrintOpenWhiskInfo(stdout)
	}

	return nil
}
