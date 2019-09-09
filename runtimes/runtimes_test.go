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

package runtimes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test only verifies that the runtimes list (i.e., "kinds") supports 1 or more
// languages we have more granular tests for (plus 1 default).
// NOTE: We do not intend for this testcase to comprehensively test for all runtime
// language versions.  Individual tests that require specific versions will have
// better, more localized failures and error messages.
func TestParseOpenWhisk(t *testing.T) {
	openwhiskHost := "https://openwhisk.ng.bluemix.net"
	openwhisk, err := ParseOpenWhisk(openwhiskHost)
	assert.Equal(t, nil, err, "parse openwhisk info error happened.")
	println(openwhisk.Runtimes)
	converted := ConvertToMap(openwhisk)
	supported := []string{"go", "java", "nodejs", "php", "python", "ruby", "swift"}
	fmt.Printf("Testing runtime kind support for: %v\n", supported)

	// TODO (see GitHub issue: #1069): add tests for newer runtime kinds (e.g., .NET, ballerina, etc.)
	for _, language := range supported {
		fmt.Printf("Testing runtime kind: [%s]...\n", language)
		assert.GreaterOrEqual(t, len(converted[language]), 2, "Runtime kind [%s] not found at [%s]", language, openwhiskHost)
	}
}
