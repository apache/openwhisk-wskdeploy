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
)

func TestInvalidDeploymentYaml(t *testing.T) {
    data :=`application:
  invalidKey: test
  name: wskdeploy-samples`
    tmpfile, err := _createTmpfile(data, "deployment_parser_test_")
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
    assert.Contains(t, err.Error(), "line 2: field invalidKey not found in struct parsers.Application")
}

func TestMappingValueDeploymentYaml(t *testing.T) {
    data :=`application:
  name: wskdeploy-samples
    packages: test`
    tmpfile, err := _createTmpfile(data, "deployment_parser_test_")
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
    assert.Contains(t, err.Error(), "line 2: mapping values are not allowed in this context")
}
