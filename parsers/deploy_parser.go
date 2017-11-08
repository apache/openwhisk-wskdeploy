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
    "strings"
)

func (dm *YAMLParser) UnmarshalDeployment(input []byte, deploy *YAML) error {
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
        return &dplyyaml, utils.NewYAMLFileReadError(err.Error())
    }
	err = dm.UnmarshalDeployment(content, &dplyyaml)
    if err != nil {
        lines, msgs := dm.convertErrorToLinesMsgs(err.Error())
        return &dplyyaml, utils.NewYAMLParserErr(deploymentPath, lines, msgs)
    }
	dplyyaml.Filepath = deploymentPath
    dplyyamlEnvVar := ReadEnvVariable(&dplyyaml)
	return dplyyamlEnvVar, nil
}

func (dm *YAMLParser) convertErrorToLinesMsgs(errorString string) (lines []string, msgs []string) {
    strs := strings.Split(errorString, "\n")
    for i := 0; i < len(strs); i++ {
        var errorMsg string
	if strings.Contains(strs[i], utils.LINE) {
		errorMsg = strings.Replace(strs[i], utils.LINE, "(on or near) "+utils.LINE, 1)
	} else {
		errorMsg = strs[i]
	}
        lines = append(lines, utils.UNKNOWN)
        msgs = append(msgs, strings.TrimSpace(errorMsg))
    }
    return
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
