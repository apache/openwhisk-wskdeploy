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
	"strings"
	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
    "github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
    "path"
)

const (
	COMMANDLINE = "wskdeploy command line"
	DEFAULTVALUE = "default value"
	WSKPROPS = ".wskprops"
	WHISKPROPERTY = "whisk.properties"
	INTERINPUT = "interactve input"
)

type PropertyValue struct {
	Value string
	Source string
}

var GetPropertyValue = func (prop PropertyValue, newValue string, source string) PropertyValue {
	if len(prop.Value) == 0 && len(newValue) > 0 {
		prop.Value = newValue
		prop.Source = source
	}
	return prop
}

var GetWskPropFromWskprops = func (pi whisk.Properties, proppath string) (*whisk.Wskprops, error) {
	return whisk.GetWskPropFromWskprops(pi, proppath)
}

var GetWskPropFromWhiskProperty = func (pi whisk.Properties) (*whisk.Wskprops, error) {
	return whisk.GetWskPropFromWhiskProperty(pi)
}

var GetCommandLineFlags = func () (string, string, string) {
	return utils.Flags.ApiHost, utils.Flags.Auth, utils.Flags.Namespace
}

var CreateNewClient = func (httpClient *http.Client, config_input *whisk.Config) (*whisk.Client, error) {
	return whisk.NewClient(http.DefaultClient, clientConfig)
}

func NewWhiskConfig(proppath string, deploymentPath string, manifestPath string, isInteractive bool) (*whisk.Config, error) {
    credential := PropertyValue {}
    namespace := PropertyValue {}
    apiHost := PropertyValue {}

    // First, we look up the above variables in the deployment file.
    if utils.FileExists(deploymentPath) {
        mm := parsers.NewYAMLParser()
        deployment := mm.ParseDeployment(deploymentPath)
        credential.Value = deployment.Application.Credential
        credential.Source = path.Base(deploymentPath)
        namespace.Value = deployment.Application.Namespace
        namespace.Source = path.Base(deploymentPath)
        apiHost.Value = deployment.Application.ApiHost
        apiHost.Source = path.Base(deploymentPath)
    }

    if len(credential.Value) == 0 || len(namespace.Value) == 0 || len(apiHost.Value) == 0 {
        if utils.FileExists(manifestPath) {
            mm := parsers.NewYAMLParser()
            manifest := mm.ParseManifest(manifestPath)
            credential = GetPropertyValue(credential, manifest.Package.Credential, path.Base(manifestPath))
            namespace = GetPropertyValue(namespace, manifest.Package.Namespace, path.Base(manifestPath))
            apiHost = GetPropertyValue(apiHost, manifest.Package.ApiHost, path.Base(manifestPath))
        }
    }

    // If the variables are not correctly assigned, we look up auth key and api host in the command line. The namespace
    // is currently not available in command line, which can be added later.
    apihost, auth, ns := GetCommandLineFlags()
    credential = GetPropertyValue(credential, auth, COMMANDLINE)
    namespace = GetPropertyValue(namespace, ns, COMMANDLINE)
    apiHost = GetPropertyValue(apiHost, apihost, COMMANDLINE)

    // Third, we need to look up the variables in .wskprops file.
    pi := whisk.PropertiesImp {
        OsPackage: whisk.OSPackageImp{},
    }

    // The error raised here can be neglected, because we will handle it in the end of this function.
    wskprops, _ := GetWskPropFromWskprops(pi, proppath)
    credential = GetPropertyValue(credential, wskprops.AuthKey, WSKPROPS)
    namespace = GetPropertyValue(namespace, wskprops.Namespace, WSKPROPS)
    apiHost = GetPropertyValue(apiHost, wskprops.APIHost, WSKPROPS)

    // Fourth, we look up the variables in whisk.properties on a local openwhisk deployment.
    if len(credential.Value) == 0 || len(apiHost.Value) == 0 {
        // No need to keep the default value for namespace, since both of auth and apihost are not set after .wskprops.
        // whisk.property will set the default value as well.
        apiHost.Value = ""
        apiHost.Source = DEFAULTVALUE
    }

    // The error raised here can be neglected, because we will handle it in the end of this function.
    whiskproperty, _ := GetWskPropFromWhiskProperty(pi)
    credential = GetPropertyValue(credential, whiskproperty.AuthKey, WHISKPROPERTY)
    namespace = GetPropertyValue(namespace, whiskproperty.Namespace, WHISKPROPERTY)
    apiHost = GetPropertyValue(apiHost, whiskproperty.APIHost, WHISKPROPERTY)

    // If we still can not find the variables we need, check if it is interactive mode. If so, we accept the input
    // from the user. The namespace will be set to a default value, when the code reaches this line, because WSKPROPS
    // has a default value for namespace.
    if len(apiHost.Value) == 0 && isInteractive == true {
        host, err := promptForValue("\nPlease provide the hostname for OpenWhisk [default value is openwhisk.ng.bluemix.net]: ")
        utils.Check(err)
        if host == "" {
            host = "openwhisk.ng.bluemix.net"
        }
        apiHost.Value = host
        apiHost.Source = INTERINPUT
    }

    if len(credential.Value) == 0 && isInteractive == true {
        cred, err := promptForValue("\nPlease provide an authentication token: ")
        utils.Check(err)
        credential.Value = cred
        credential.Source = INTERINPUT

        // The namespace is always associated with the credential. Both of them should be picked up from the same source.
        if len(namespace.Value) == 0 || namespace.Value == whisk.DEFAULT_NAMESPACE {
            ns, err := promptForValue("\nPlease provide a namespace [default value is guest]: ")
            utils.Check(err)

            source := INTERINPUT
            if ns == "" {
                ns = whisk.DEFAULT_NAMESPACE
                source = DEFAULTVALUE
            }
            namespace.Value = ns
            namespace.Source = source
        }
    }

    clientConfig = &whisk.Config{
        AuthToken: credential.Value, //Authtoken
        Namespace: namespace.Value,  //Namespace
        Host: apiHost.Value,
        Version:   "v1",
        Insecure:  true, // true if you want to ignore certificate signing
    }

    if len(credential.Value) == 0 {
        errStr := wski18n.T("The authentication key is not configured.\n")
        whisk.Debug(whisk.DbgError, errStr)
        return clientConfig, utils.NewInvalidWskpropsError(errStr)
    }

    if len(apiHost.Value) == 0 {
        errStr := wski18n.T("The API host is not configured.\n")
        whisk.Debug(whisk.DbgError, errStr)
        return clientConfig, utils.NewInvalidWskpropsError(errStr)
    }

    if len(namespace.Value) == 0 {
        errStr := wski18n.T("The namespace is not configured.\n")
        whisk.Debug(whisk.DbgError, errStr)
        return clientConfig, utils.NewInvalidWskpropsError(errStr)
    }

    stdout := wski18n.T("The API host is {{.apihost}}, from {{.apisource}}.\n",
        map[string]interface{}{"apihost": apiHost.Value, "apisource": apiHost.Source})
    whisk.Debug(whisk.DbgInfo, stdout)

    stdout = wski18n.T("The auth key is set, from {{.authsource}}.\n",
        map[string]interface{}{"authsource": credential.Source})
    whisk.Debug(whisk.DbgInfo, stdout)

    stdout = wski18n.T("The namespace is {{.namespace}}, from {{.namespacesource}}.\n",
        map[string]interface{}{"namespace": namespace.Value, "namespacesource": namespace.Source})
    whisk.Debug(whisk.DbgInfo, stdout)
    return clientConfig, nil
}

var promptForValue = func (msg string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(msg)

	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	return text, nil

}
