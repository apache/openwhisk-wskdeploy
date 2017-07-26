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
var testfile = "../tests/dat/deployment.yaml"

func TestLocalReader_ReadLocal(t *testing.T) {
	b, err := contentReader.ReadLocal(testfile)
	Check(err)
	if b == nil {
		t.Error("get local centent failed")
	}
}

func TestURLReader_ReadUrl(t *testing.T) {
	var exampleUrl = "http://www.baidu.com"
	b, err := contentReader.ReadUrl(exampleUrl)
	Check(err)
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
}
