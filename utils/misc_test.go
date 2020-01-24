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

package utils

import (
	"github.com/apache/openwhisk-wskdeploy/dependencies"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var contentReader = new(ContentReader)
var testfile = "../tests/dat/deployment_data_action_with_inputs_outputs.yaml"

func TestLocalReader_ReadLocal(t *testing.T) {
	b, _ := contentReader.ReadLocal(testfile)
	if b == nil {
		t.Error("get local centent failed")
	}
}

func TestURLReader_ReadUrl(t *testing.T) {
	var exampleUrl = "https://www.apache.org/"
	b, _ := contentReader.ReadUrl(exampleUrl)
	if b == nil {
		t.Error("get web content failed")
	}
}

func TestDependencies(t *testing.T) {
	var record = dependencies.NewDependencyRecord("projectPath", "packageName", "http://github.com/user/repo", "master", nil, nil, false)
	assert.Equal(t, "projectPath", record.ProjectPath, "ProjectPath is wrong")
	assert.Equal(t, "http://github.com/user/repo", record.Location, "URL is wrong")
	assert.Equal(t, "http://github.com/user/repo", record.BaseRepo, "BaseRepo is wrong")
	assert.Equal(t, "", record.SubFolder, "SubFolder is wrong")

	record = dependencies.NewDependencyRecord("projectPath", "packageName", "http://github.com/user/repo/subfolder1/subfolder2", "master", nil, nil, false)
	assert.Equal(t, "projectPath", record.ProjectPath, "ProjectPath is wrong")
	assert.Equal(t, "http://github.com/user/repo/subfolder1/subfolder2", record.Location, "URL is wrong")
	assert.Equal(t, "http://github.com/user/repo", record.BaseRepo, "BaseRepo is wrong")
	assert.Equal(t, "/subfolder1/subfolder2", record.SubFolder, "SubFolder is wrong")
}

func TestNewZipWriter(t *testing.T) {
	filePath := "../tests/src/integration/zipaction/actions/cat"
	zipName := filePath + ".zip"
	err := NewZipWritter(filePath, zipName, make([][]string, 0), make([]string, 0), "").Zip()
	defer os.Remove(zipName)
	assert.Equal(t, nil, err, "zip folder error happened.")
}
