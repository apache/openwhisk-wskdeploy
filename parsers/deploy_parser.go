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
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"gopkg.in/yaml.v2"
)

func (dm *YAMLParser) unmarshalDeployment(input []byte, deploy *YAML) error {
	err := yaml.UnmarshalStrict(input, deploy)
	if err != nil {
		return err
	}
	return nil
}

func (dm *YAMLParser) ParseDeployment(deploymentPath string) (*YAML, error) {
	dplyyaml := YAML{}
	content, err := new(utils.ContentReader).LocalReader.ReadLocal(deploymentPath)
    if err != nil {
        return &dplyyaml, utils.NewFileReadError(deploymentPath, err.Error())
    }
	err = dm.unmarshalDeployment(content, &dplyyaml)
    if err != nil {
        return &dplyyaml, utils.NewYAMLParserErr(deploymentPath, err)
    }
	dplyyaml.Filepath = deploymentPath
    dplyyamlEnvVar := ReadEnvVariable(&dplyyaml)
	return dplyyamlEnvVar, nil
}


//********************Project functions*************************//
//This is for parse the deployment yaml file.
func (app *Project) GetPackageList() []Package {
	var s1 []Package = make([]Package, 0)
	for _, pkg := range app.Packages {
		pkg.Packagename = pkg.Packagename
		s1 = append(s1, pkg)
	}
	return s1
}
