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
	"testing"

	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
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
	var exampleUrl = "http://www.baidu.com"
	b, _ := contentReader.ReadUrl(exampleUrl)
	if b == nil {
		t.Error("get web content failed")
	}
}

// The dollar sign test cases.
func TestGetEnvVar(t *testing.T) {
	os.Setenv("NoDollar", "NO dollar")
	os.Setenv("WithDollar", "oh, dollars!")
	os.Setenv("5000", "5000")
	fmt.Println(GetEnvVar("NoDollar"))
	fmt.Println(GetEnvVar("$WithDollar"))
	fmt.Println(GetEnvVar("$5000"))
	assert.Equal(t, "NoDollar", GetEnvVar("NoDollar"), "NoDollar should be no change.")
	assert.Equal(t, "oh, dollars!", GetEnvVar("$WithDollar"), "dollar sign should be handled.")
	assert.Equal(t, "5000", GetEnvVar("5000"), "Should be no difference between integer and string.")
	assert.Equal(t, "", GetEnvVar("$WithDollarAgain"), "if not found in environemnt, return empty string.")
	assert.Equal(t, "oh, dollars!.ccc.aaa", GetEnvVar("${WithDollar}.ccc.aaa"), "String concatenation fail")
	assert.Equal(t, "ddd.NO dollar.aaa", GetEnvVar("ddd.${NoDollar}.aaa"), "String concatenation fail")
	assert.Equal(t, "oh, dollars!.NO dollar.aaa", GetEnvVar("${WithDollar}.${NoDollar}.aaa"), "String concatenation fail")
	assert.Equal(t, "ddd.ccc.oh, dollars!", GetEnvVar("ddd.ccc.${WithDollar}"), "String concatenation fail")
	assert.Equal(t, "", GetEnvVar("$WithDollarAgain.ccc.aaa"), "String concatenation fail")
	assert.Equal(t, "ddd..aaa", GetEnvVar("ddd.${WithDollarAgain}.aaa"), "String concatenation fail")
	assert.Equal(t, "oh, dollars!NO dollar.NO dollar", GetEnvVar("${WithDollar}${NoDollar}.${NoDollar}"), "String concatenation fail")
}

func TestDependencies(t *testing.T) {
	var record = NewDependencyRecord("projectPath","packageName","http://github.com/user/repo","master",nil,nil,false)
	assert.Equal(t, "projectPath", record.ProjectPath,"ProjectPath is wrong")
	assert.Equal(t, "http://github.com/user/repo", record.Location, "URL is wrong")
	assert.Equal(t, "http://github.com/user/repo", record.BaseRepo, "BaseRepo is wrong")
	assert.Equal(t, "", record.SubFolder, "SubFolder is wrong")

	record = NewDependencyRecord("projectPath","packageName","http://github.com/user/repo/subfolder1/subfolder2","master",nil,nil,false)
	assert.Equal(t, "projectPath", record.ProjectPath,"ProjectPath is wrong")
	assert.Equal(t, "http://github.com/user/repo/subfolder1/subfolder2", record.Location, "URL is wrong")
	assert.Equal(t, "http://github.com/user/repo", record.BaseRepo, "BaseRepo is wrong")
	assert.Equal(t, "/subfolder1/subfolder2", record.SubFolder, "SubFolder is wrong")
}

func TestParseOpenWhisk(t *testing.T) {
	openwhiskHost := "https://openwhisk.ng.bluemix.net"
    openwhisk, err := ParseOpenWhisk(openwhiskHost)
    assert.Equal(t, nil, err, "parse openwhisk info error happened.")
    converted := ConvertToMap(openwhisk)
    assert.Equal(t, 1, len(converted["nodejs"]), "not expected length")
    assert.Equal(t, 1, len(converted["php"]),  "not expected length")
    assert.Equal(t, 1, len(converted["java"]), "not expected length")
    assert.Equal(t, 3, len(converted["python"]), "not expected length")
    assert.Equal(t, 2, len(converted["swift"]), "not expected length")
}

func TestNewZipWritter(t *testing.T) {
    filePath := "../tests/src/integration/zipaction/actions/cat"
    zipName := filePath + ".zip"
    err := NewZipWritter(filePath, zipName).Zip()
    defer os.Remove(zipName)
    assert.Equal(t, nil, err, "zip folder error happened.")
}
