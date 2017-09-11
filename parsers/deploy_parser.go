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

package parsers

import (
	"fmt"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"gopkg.in/yaml.v2"
)

func (dm *YAMLParser) UnmarshalDeployment(input []byte, deploy *DeploymentYAML) error {
	err := yaml.UnmarshalStrict(input, deploy)
	if err != nil {
		fmt.Printf("error happened during unmarshal :%v", err)
		return err
	}
	return nil
}

func (dm *YAMLParser) MarshalDeployment(deployment *DeploymentYAML) (output []byte, err error) {
	data, err := yaml.Marshal(deployment)
	if err != nil {
		fmt.Printf("err happened during marshal :%v", err)
		return nil, err
	}
	return data, nil
}

func (dm *YAMLParser) ParseDeployment(deploymentPath string) (*DeploymentYAML, error) {
	dplyyaml := DeploymentYAML{}
	content, err := new(utils.ContentReader).LocalReader.ReadLocal(deploymentPath)
    if err != nil {
        return &dplyyaml, utils.NewInputYamlFileError(err.Error())
    }
	err = dm.UnmarshalDeployment(content, &dplyyaml)
    if err != nil {
        if err != nil {
            return &dplyyaml, utils.NewParserErr(deploymentPath, err.Error())
        }
    }
	dplyyaml.Filepath = deploymentPath
	return &dplyyaml, nil
}

//********************Application functions*************************//
//This is for parse the deployment yaml file.
func (app *Application) GetPackageList() []Package {
	var s1 []Package = make([]Package, 0)
	for _, pkg := range app.Packages {
		pkg.Packagename = pkg.Packagename
		s1 = append(s1, pkg)
	}
	return s1
}
