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

package wskenv

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The dollar sign test cases.
func TestInterpolateStringWithEnvVar(t *testing.T) {
	os.Setenv("NoDollar", "NO dollar")
	os.Setenv("WithDollar", "oh, dollars!")
	os.Setenv("5000", "5000")
	fmt.Println(InterpolateStringWithEnvVar("NoDollar"))
	fmt.Println(InterpolateStringWithEnvVar("$WithDollar"))
	fmt.Println(InterpolateStringWithEnvVar("$5000"))
	assert.Equal(t, "NoDollar", InterpolateStringWithEnvVar("NoDollar"), "NoDollar should be no change.")
	assert.Equal(t, "oh, dollars!", InterpolateStringWithEnvVar("$WithDollar"), "dollar sign should be handled.")
	assert.Equal(t, "5000", InterpolateStringWithEnvVar("5000"), "Should be no difference between integer and string.")
	assert.Equal(t, "", InterpolateStringWithEnvVar("$WithDollarAgain"), "if not found in environment, return empty string.")
	assert.Equal(t, "oh, dollars!.ccc.aaa", InterpolateStringWithEnvVar("${WithDollar}.ccc.aaa"), "String concatenation fail")
	assert.Equal(t, "ddd.NO dollar.aaa", InterpolateStringWithEnvVar("ddd.${NoDollar}.aaa"), "String concatenation fail")
	assert.Equal(t, "oh, dollars!.NO dollar.aaa", InterpolateStringWithEnvVar("${WithDollar}.${NoDollar}.aaa"), "String concatenation fail")
	assert.Equal(t, "ddd.ccc.oh, dollars!", InterpolateStringWithEnvVar("ddd.ccc.${WithDollar}"), "String concatenation fail")
	assert.Equal(t, "", InterpolateStringWithEnvVar("$WithDollarAgain.ccc.aaa"), "String concatenation fail")
	assert.Equal(t, "ddd..aaa", InterpolateStringWithEnvVar("ddd.${WithDollarAgain}.aaa"), "String concatenation fail")
	assert.Equal(t, "oh, dollars!NO dollar.NO dollar", InterpolateStringWithEnvVar("${WithDollar}${NoDollar}.${NoDollar}"), "String concatenation fail")
}
