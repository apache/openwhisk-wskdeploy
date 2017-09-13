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

package parsers

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
    "io/ioutil"
)

func createTmpfile(data string, filename string) (f *os.File, err error) {
    dir, _ := os.Getwd()
    tmpfile, err := ioutil.TempFile(dir, filename)
    if err != nil {
        return nil, err
    }
    _, err = tmpfile.Write([]byte(data))
    if err != nil {
        return tmpfile, err
    }
    return tmpfile, nil
}

func TestInvalidKeyDeploymentYaml(t *testing.T) {
    data :=`application:
  name: wskdeploy-samples
  invalidKey: test`
    tmpfile, err := createTmpfile(data, "deployment_parser_test_")
    if err != nil {
        assert.Fail(t, "Failed to create temp file")
    }
    defer func() {
        tmpfile.Close()
        os.Remove(tmpfile.Name())
    }()
    p := NewYAMLParser()
    _, err = p.ParseDeployment(tmpfile.Name())
    assert.NotNil(t, err)
    // go-yaml/yaml prints the wrong line number for mapping values. It should be 3.
    assert.Contains(t, err.Error(), "field invalidKey not found in struct parsers.Application: Line 2, its neighbour lines, or the lines on the same level")
}

func TestMappingValueDeploymentYaml(t *testing.T) {
    data :=`application:
  name: wskdeploy-samples
    packages: test`
    tmpfile, err := createTmpfile(data, "deployment_parser_test_")
    if err != nil {
        assert.Fail(t, "Failed to create temp file")
    }
    defer func() {
        tmpfile.Close()
        os.Remove(tmpfile.Name())
    }()
    p := NewYAMLParser()
    _, err = p.ParseDeployment(tmpfile.Name())
    assert.NotNil(t, err)
    // go-yaml/yaml prints the wrong line number for mapping values. It should be 3.
    assert.Contains(t, err.Error(), "mapping values are not allowed in this context: Line 2, its neighbour lines, or the lines on the same level")
}

func TestMissingRootNodeDeploymentYaml(t *testing.T) {
    data :=`name: wskdeploy-samples`
    tmpfile, err := createTmpfile(data, "deployment_parser_test_")
    if err != nil {
        assert.Fail(t, "Failed to create temp file")
    }
    defer func() {
        tmpfile.Close()
        os.Remove(tmpfile.Name())
    }()
    p := NewYAMLParser()
    _, err = p.ParseDeployment(tmpfile.Name())
    assert.NotNil(t, err)
    // go-yaml/yaml prints the wrong line number for mapping values. It should be 3.
    assert.Contains(t, err.Error(), "field name not found in struct parsers.DeploymentYAML: Line 1, its neighbour lines, or the lines on the same level")
}
