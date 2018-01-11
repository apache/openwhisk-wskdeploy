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
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskderrors"
)

const (
	COMMANDLINE = "wskdeploy command line"
	DEFAULTVALUE = "default value"
	WSKPROPS = ".wskprops"
	WHISKPROPERTY = "whisk.properties"
	INTERINPUT = "interactve input"
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

var GetCommandLineFlags = func() (string, string, string, string, string) {
	return utils.Flags.ApiHost, utils.Flags.Auth, utils.Flags.Namespace, utils.Flags.Key, utils.Flags.Cert
}

var CreateNewClient = func(config_input *whisk.Config) (*whisk.Client, error) {
	var netClient = &http.Client{
		Timeout: time.Second * utils.DEFAULT_HTTP_TIMEOUT,
	}
	return whisk.NewClient(netClient, config_input)
}

// we are reading openwhisk credentials (apihost, namespace, and auth) in the following precedence order:
// (1) wskdeploy command line `wskdeploy --apihost --namespace --auth`
// (2) deployment file
// (3) manifest file
// (4) .wskprops
// (5) prompt for values in interactive mode if any of them are missing
func NewWhiskConfig(proppath string, deploymentPath string, manifestPath string, isInteractive bool) (*whisk.Config, error) {
	// struct to store credential, namespace, and host with their respective source
	credential := PropertyValue{}
	namespace := PropertyValue{}
	apiHost := PropertyValue{}
	key := PropertyValue{}
	cert := PropertyValue{}

	// read credentials from command line
	apihost, auth, ns, keyfile, certfile := GetCommandLineFlags()
	credential = GetPropertyValue(credential, auth, COMMANDLINE)
	namespace = GetPropertyValue(namespace, ns, COMMANDLINE)
	apiHost = GetPropertyValue(apiHost, apihost, COMMANDLINE)
	key = GetPropertyValue(key, keyfile, COMMANDLINE)
	cert = GetPropertyValue(cert, certfile, COMMANDLINE)

	// now, read them from deployment file if not found on command line
	if len(credential.Value) == 0 || len(namespace.Value) == 0 || len(apiHost.Value) == 0 {
		if utils.FileExists(deploymentPath) {
			mm := parsers.NewYAMLParser()
			deployment, _ := mm.ParseDeployment(deploymentPath)
			credential = GetPropertyValue(credential, deployment.GetProject().Credential, path.Base(deploymentPath))
			namespace = GetPropertyValue(namespace, deployment.GetProject().Namespace, path.Base(deploymentPath))
			apiHost = GetPropertyValue(apiHost, deployment.GetProject().ApiHost, path.Base(deploymentPath))
		}
	}

	// read credentials from manifest file as didn't find them on command line and in deployment file
	if len(credential.Value) == 0 || len(namespace.Value) == 0 || len(apiHost.Value) == 0 {
		if utils.FileExists(manifestPath) {
			mm := parsers.NewYAMLParser()
			manifest, _ := mm.ParseManifest(manifestPath)
			if manifest.Package.Packagename != "" {
				credential = GetPropertyValue(credential, manifest.Package.Credential, path.Base(manifestPath))
				namespace = GetPropertyValue(namespace, manifest.Package.Namespace, path.Base(manifestPath))
				apiHost = GetPropertyValue(apiHost, manifest.Package.ApiHost, path.Base(manifestPath))
			} else if manifest.Packages != nil {
				if len(manifest.Packages) == 1 {
					for _, pkg := range manifest.Packages {
						credential = GetPropertyValue(credential, pkg.Credential, path.Base(manifestPath))
						namespace = GetPropertyValue(namespace, pkg.Namespace, path.Base(manifestPath))
						apiHost = GetPropertyValue(apiHost, pkg.ApiHost, path.Base(manifestPath))
					}
				}
			}
		}
	}

	// Third, we need to look up the variables in .wskprops file.
	pi := whisk.PropertiesImp{
		OsPackage: whisk.OSPackageImp{},
	}

	// The error raised here can be neglected, because we will handle it in the end of this function.
	wskprops, _ := GetWskPropFromWskprops(pi, proppath)
	credential = GetPropertyValue(credential, wskprops.AuthKey, WSKPROPS)
	namespace = GetPropertyValue(namespace, wskprops.Namespace, WSKPROPS)
	apiHost = GetPropertyValue(apiHost, wskprops.APIHost, WSKPROPS)
	key = GetPropertyValue(key, wskprops.Key, WSKPROPS)
	cert = GetPropertyValue(cert, wskprops.Cert, WSKPROPS)

	// TODO() see if we can split the following whisk prop logic into a separate function
	// now, read credentials from whisk.properties but this is only acceptable within Travis
	// whisk.properties will soon be deprecated and should not be used for any production deployment
	whiskproperty, _ := GetWskPropFromWhiskProperty(pi)

	var warnmsg string

	credential = GetPropertyValue(credential, whiskproperty.AuthKey, WHISKPROPERTY)
	if credential.Source == WHISKPROPERTY {
		warnmsg = wski18n.T(wski18n.ID_WARN_WHISK_PROPS_DEPRECATED,
			map[string]interface{}{"key": "authenticaton key"})
		wskprint.PrintlnOpenWhiskWarning(warnmsg)
	}
	namespace = GetPropertyValue(namespace, whiskproperty.Namespace, WHISKPROPERTY)
	if namespace.Source == WHISKPROPERTY {
		warnmsg = wski18n.T(wski18n.ID_WARN_WHISK_PROPS_DEPRECATED,
			map[string]interface{}{"key": "namespace"})
		wskprint.PrintlnOpenWhiskWarning(warnmsg)
	}
	apiHost = GetPropertyValue(apiHost, whiskproperty.APIHost, WHISKPROPERTY)
	if apiHost.Source == WHISKPROPERTY {
		warnmsg = wski18n.T(wski18n.ID_WARN_WHISK_PROPS_DEPRECATED,
			map[string]interface{}{"key": "API host"})
		wskprint.PrintlnOpenWhiskWarning(warnmsg)
	}

	// set namespace to default namespace if not yet found
	if len(apiHost.Value) != 0 && len(credential.Value) != 0 && len(namespace.Value) == 0 {
		namespace.Value = whisk.DEFAULT_NAMESPACE
		namespace.Source = DEFAULTVALUE
	}

	// TODO() See if we can split off the interactive logic into a separate function
	// If we still can not find the values we need, check if it is interactive mode.
	// If so, we prompt users for the input.
	// The namespace is set to a default value at this point if not provided.
	if len(apiHost.Value) == 0 && isInteractive == true {
		host := promptForValue(wski18n.T(wski18n.ID_MSG_PROMPT_APIHOST))
		if host == "" {
			// TODO() programmatically tell caller that we are using this default
			// TODO() make this configurable or remove
			host = "openwhisk.ng.bluemix.net"
		}
		apiHost.Value = host
		apiHost.Source = INTERINPUT
	}

	if len(credential.Value) == 0 && isInteractive == true {
		cred := promptForValue(wski18n.T(wski18n.ID_MSG_PROMPT_AUTHKEY))
		credential.Value = cred
		credential.Source = INTERINPUT

		// The namespace is always associated with the credential. Both of them should be picked up from the same source.
		if len(namespace.Value) == 0 || namespace.Value == whisk.DEFAULT_NAMESPACE {
			tempNamespace := promptForValue(wski18n.T(wski18n.ID_MSG_PROMPT_NAMESPACE))
			source := INTERINPUT

			if tempNamespace == "" {
				tempNamespace = whisk.DEFAULT_NAMESPACE
				source = DEFAULTVALUE
			}

			namespace.Value = tempNamespace
			namespace.Source = source
		}
	}

	mode := true
	if (len(cert.Value) != 0 && len(key.Value) != 0) {
		mode = false
	}

	clientConfig = &whisk.Config{
		AuthToken: credential.Value, //Authtoken
		Namespace: namespace.Value, //Namespace
		Host:      apiHost.Value,
		Version:   "v1",
		Cert:      cert.Value,
		Key:       key.Value,
		Insecure:  mode, // true if you want to ignore certificate signing
	}

	// validate we have credential, apihost and namespace
	err := validateClientConfig(credential, apiHost, namespace)

	return clientConfig, err
}

func validateClientConfig(credential PropertyValue, apiHost PropertyValue, namespace PropertyValue) (error) {

	// Display error message based upon which config value was missing
	if len(credential.Value) == 0 || len(apiHost.Value) == 0 || len(namespace.Value) == 0 {
		var errmsg string
		if len(credential.Value) == 0 {
			errmsg = wski18n.T(wski18n.ID_MSG_CONFIG_MISSING_AUTHKEY)
		}

		if len(apiHost.Value) == 0 {
			errmsg = wski18n.T(wski18n.ID_MSG_CONFIG_MISSING_APIHOST)
		}

		if len(namespace.Value) == 0 {
			errmsg = wski18n.T(wski18n.ID_MSG_CONFIG_MISSING_NAMESPACE)
		}

		return wskderrors.NewWhiskClientInvalidConfigError(errmsg)
	}

	// Show caller what final values we used for credential, apihost and namespace
	stdout := wski18n.T(wski18n.ID_MSG_CONFIG_INFO_APIHOST_X_host_X_source_X,
		map[string]interface{}{"host": apiHost.Value, "source": apiHost.Source})
	wskprint.PrintOpenWhiskStatus(stdout)

	stdout = wski18n.T(wski18n.ID_MSG_CONFIG_INFO_AUTHKEY_X_source_X,
		map[string]interface{}{"source": credential.Source})
	wskprint.PrintOpenWhiskStatus(stdout)

	stdout = wski18n.T(wski18n.ID_MSG_CONFIG_INFO_NAMESPACE_X_namespace_X_source_X,
		map[string]interface{}{"namespace": namespace.Value, "source": namespace.Source})
	wskprint.PrintOpenWhiskStatus(stdout)

	return nil
}

// TODO() move into its own package "wskread" and add support for passing in default value
var promptForValue = func(msg string) (string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(msg)

	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	return text
}
