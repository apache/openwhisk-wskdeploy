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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOpenWhisk(t *testing.T) {
	openwhiskHost := "https://openwhisk.ng.bluemix.net"
	openwhisk, err := ParseOpenWhisk(openwhiskHost)
	assert.Equal(t, nil, err, "parse openwhisk info error happened.")
	println(openwhisk.Runtimes)
	converted := ConvertToMap(openwhisk)
	println(converted["nodejs"])
	println(converted["python"])
	//println(converted["go"])
	//one larger than in api response because of explicit :default
	assert.Equal(t, 4, len(converted["nodejs"]), "not expected length")
	assert.Equal(t, 4, len(converted["php"]), "not expected length")
	assert.Equal(t, 2, len(converted["java"]), "not expected length")
	assert.Equal(t, 7, len(converted["python"]), "not expected length")
	assert.Equal(t, 2, len(converted["ruby"]), "not expected length")
	assert.Equal(t, 4, len(converted["swift"]), "not expected length")
	// cannot pass this test until it is deployed in OpenWhisk
	//assert.Equal(t, 2, len(converted["go"]), "not expected length")
}
