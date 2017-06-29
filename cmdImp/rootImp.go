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

package cmdImp

import (
	"errors"
	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-client-go/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/deployers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

type DeployParams struct {
	Verbose        bool
	ProjectPath    string
	ManifestPath   string
	DeploymentPath string
	UseDefaults    bool
	UseInteractive bool
}

func Deploy(params DeployParams) error {

	whisk.SetVerbose(params.Verbose)

	projectPath, err := filepath.Abs(params.ProjectPath)
	utils.Check(err)

	if params.ManifestPath == "" {
		if ok, _ := regexp.Match(deployers.ManifestFileNameYml, []byte(params.ManifestPath)); ok {
			params.ManifestPath = path.Join(projectPath, deployers.ManifestFileNameYml)
		} else {
			params.ManifestPath = path.Join(projectPath, deployers.ManifestFileNameYaml)
		}
	} else {
		if _, err := os.Stat(path.Join(projectPath, "manifest.yaml")); err == nil {
			params.ManifestPath = path.Join(projectPath, deployers.ManifestFileNameYaml)
		} else if _, err := os.Stat(path.Join(projectPath, "manifest.yml")); err == nil {
			params.ManifestPath = path.Join(projectPath, deployers.ManifestFileNameYml)
		}
	}

	if params.DeploymentPath == "" {
		if ok, _ := regexp.Match(deployers.DeploymentFileNameYml, []byte(params.ManifestPath)); ok {
			params.DeploymentPath = path.Join(projectPath, deployers.DeploymentFileNameYml)
		} else {
			params.DeploymentPath = path.Join(projectPath, deployers.DeploymentFileNameYaml)
		}
	} else {
		if _, err := os.Stat(path.Join(projectPath, "deployment.yaml")); err == nil {
			params.DeploymentPath = path.Join(projectPath, deployers.DeploymentFileNameYaml)
		} else if _, err := os.Stat(path.Join(projectPath, "deployment.yml")); err == nil {
			params.DeploymentPath = path.Join(projectPath, deployers.DeploymentFileNameYml)
		}
	}

	if utils.MayExists(params.ManifestPath) {

		var deployer = deployers.NewServiceDeployer()
		deployer.ProjectPath = projectPath
		deployer.ManifestPath = params.ManifestPath
		deployer.DeploymentPath = params.DeploymentPath
		// perform some quick check here.
		go func() {
			deployer.Check()
		}()
		deployer.IsDefault = params.UseDefaults

		deployer.IsInteractive = params.UseInteractive

		// master record of any dependency that has been downloaded
		deployer.DependencyMaster = make(map[string]utils.DependencyRecord)

		propPath := ""
		if !utils.Flags.WithinOpenWhisk {
			userHome := utils.GetHomeDirectory()
			propPath = path.Join(userHome, ".wskprops")
		}
		whiskClient, clientConfig := deployers.NewWhiskClient(propPath, params.DeploymentPath, deployer.IsInteractive)
		deployer.Client = whiskClient
		deployer.ClientConfig = clientConfig

		err := deployer.ConstructDeploymentPlan()

		if err != nil {
			utils.Check(err)
			return err
		}

		err = deployer.Deploy()
		if err != nil {
			utils.Check(err)
			return err
		} else {
			return nil
		}

	} else {
		if utils.Flags.WithinOpenWhisk {
			utils.PrintOpenWhiskError(wski18n.T("missing manifest.yaml file"))
			return errors.New("missing manifest.yaml file")
		} else {
			log.Println("missing manifest.yaml file")
			return errors.New("missing manifest.yaml file")
		}
	}

}
