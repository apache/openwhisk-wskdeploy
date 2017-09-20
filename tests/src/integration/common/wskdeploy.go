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

package common

import (
	"bytes"
	"fmt"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
    "github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"os"
	"os/exec"
	"strings"
    "errors"
)

const (
    cmd = "wskdeploy"
    BLUEMIX_APIHOST = "BLUEMIX_APIHOST"
    BLUEMIX_NAMESPACE = "BLUEMIX_NAMESPACE"
    BLUEMIX_AUTH = "BLUEMIX_AUTH"
)

type Wskdeploy struct {
	Path string
	Dir  string
}

func NewWskdeploy() *Wskdeploy {
	return NewWskWithPath(os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/")
}

func GetWskpropsFromEnvVars(apiHost string, namespace string, authKey string) *whisk.Wskprops {
    return GetWskpropsFromValues(os.Getenv(apiHost), os.Getenv(namespace), os.Getenv(authKey), "v1")
}

func GetWskpropsFromValues(apiHost string, namespace string, authKey string, version string) *whisk.Wskprops {
    dep := whisk.Wskprops {
        APIHost: apiHost,
        AuthKey: authKey,
        Namespace: namespace,
        AuthAPIGWKey: "",
        APIGWSpaceSuid: "",
        Apiversion: version,
        Key: "",
        Cert: "",
        Source: "",
    }
    return &dep
}

func ValidateWskprops(wskprops *whisk.Wskprops) error {
    if len(wskprops.APIHost) == 0 {
        return errors.New("Missing APIHost for wskprops.")
    }
    if len(wskprops.Namespace) == 0 {
        return errors.New("Missing Namespace for wskprops.")
    }
    if len(wskprops.AuthKey) == 0 {
        return errors.New("Missing AuthKey for wskprops.")
    }
    return nil
}


func NewWskWithPath(path string) *Wskdeploy {
	var dep Wskdeploy
	dep.Path = cmd
	dep.Dir = path
	return &dep
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err string) {
	outputStream := colorable.NewColorableStderr()
	if len(err) > 0 {
		fmt.Fprintf(outputStream, "==> Error: %s.\n", color.RedString(err))
	} else {
		fmt.Fprintf(outputStream, "==> Error: %s.\n", color.RedString("No error message"))
	}
}

func printOutput(outs string) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s.\n", outs)
	}
}

func (wskdeploy *Wskdeploy) RunCommand(s ...string) (string, error) {
	command := exec.Command(wskdeploy.Path, s...)
	command.Dir = wskdeploy.Dir

	printCommand(command)

	var outb, errb bytes.Buffer
	command.Stdout = &outb
	command.Stderr = &errb
	err := command.Run()

	var returnError error = nil
	if err != nil {
		if len(errb.String()) > 0 {
			returnError = utils.NewTestCaseError(errb.String())
		} else {
            returnError = err
        }
	}
	printOutput(outb.String())
	if returnError != nil {
		printError(returnError.Error())
	}
	return outb.String(), returnError
}

func (wskdeploy *Wskdeploy) Deploy(manifestPath string, deploymentPath string) (string, error) {
	return wskdeploy.RunCommand("-m", manifestPath, "-d", deploymentPath)
}

func (wskdeploy *Wskdeploy) DeployWithCredentials(manifestPath string, deploymentPath string, wskprops *whisk.Wskprops) (string, error) {
    return wskdeploy.RunCommand("-m", manifestPath, "-d", deploymentPath, "--auth", wskprops.AuthKey,
        "--namespace", wskprops.Namespace, "--apihost", wskprops.APIHost, "--apiversion", wskprops.Apiversion)
}

func (wskdeploy *Wskdeploy) Undeploy(manifestPath string, deploymentPath string) (string, error) {
	return wskdeploy.RunCommand("undeploy", "-m", manifestPath, "-d", deploymentPath)
}

func (wskdeploy *Wskdeploy) UndeployWithCredentials(manifestPath string, deploymentPath string, wskprops *whisk.Wskprops) (string, error) {
    return wskdeploy.RunCommand("undeploy", "-m", manifestPath, "-d", deploymentPath, "--auth", wskprops.AuthKey,
        "--namespace", wskprops.Namespace, "--apihost", wskprops.APIHost, "--apiversion", wskprops.Apiversion)
}

func (wskdeploy *Wskdeploy) DeployProjectPathOnly(projectPath string) (string, error) {
	return wskdeploy.RunCommand("-p", projectPath)
}

func (wskdeploy *Wskdeploy) UndeployProjectPathOnly(projectPath string) (string, error) {
	return wskdeploy.RunCommand("undeploy", "-p", projectPath)

}

func (wskdeploy *Wskdeploy) DeployManifestPathOnly(manifestpath string) (string, error) {
	return wskdeploy.RunCommand("-m", manifestpath)
}

func (wskdeploy *Wskdeploy) UndeployManifestPathOnly(manifestpath string) (string, error) {
	return wskdeploy.RunCommand("undeploy", "-m", manifestpath)

}
