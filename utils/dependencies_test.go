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

	"github.com/stretchr/testify/assert"
)

/* NOTE: Though username and password is used here, this does not mean that wskdeploy fully supports private repos.
 * This is merely one of many changes required to support them.
 */

func TestLocationIsGithub(t *testing.T) {
	assert.True(t, LocationIsGithub("github.com/my-org/my-project"), "Does not allow github without a http/https prefix")
	assert.True(t, LocationIsGithub("github.ibm.com/my-org/my-project"), "Does not allow github.ibm.com without a http/https prefix")
	assert.True(t, LocationIsGithub("http://github.com/my-org/my-project"), "Does not allow github with a http/https prefix")
	assert.True(t, LocationIsGithub("http://github.ibm.com/my-org/my-project"), "Does not allow github.ibm.com with a http/https prefix")
}

func TestLocationIsGithub_WithUsernamePassword(t *testing.T) {
	assert.True(t, LocationIsGithub("username:password@github.com/my-org/my-project"), "Does not allow username/password and github without a http/https prefix")
	assert.True(t, LocationIsGithub("username:password@github.ibm.com/my-org/my-project"), "Does not allow username/password and github.ibm.com without a http/https prefix")
	assert.True(t, LocationIsGithub("http://username:password@github.com/my-org/my-project"), "Does not allow username/password and github with a http/https prefix")
	assert.True(t, LocationIsGithub("http://username:password@github.ibm.com/my-org/my-project"), "Does not allow username/password and github.ibm.com with a http/https prefix")
}

func TestLocationIsGithub_NonGithub(t *testing.T) {
	assert.False(t, LocationIsGithub("git.com/my-org/my-project"), "Allows non-github without a http/https prefix")
	assert.False(t, LocationIsGithub("git.ibm.com/my-org/my-project"), "Allows non-github.ibm.com without a http/https prefix")
	assert.False(t, LocationIsGithub("http://git.com/my-org/my-project"), "Allows non-github with a http/https prefix")
	assert.False(t, LocationIsGithub("http://git.ibm.com/my-org/my-project"), "Allows non-github.ibm.com with a http/https prefix")

	assert.False(t, LocationIsGithub("git.com/my-org/my-github"), "Thinks it is github because it is part of the project name")
	assert.False(t, LocationIsGithub("git.com/my-org/github"), "Thinks it is github because it is the project name")
	assert.False(t, LocationIsGithub("git.com/my-github/my-project"), "Thinks it is github because it is part of the organization name")
	assert.False(t, LocationIsGithub("git.com/github/my-project"), "Thinks it is github because it is the organization name")

	assert.False(t, LocationIsGithub("git.com"), "Allows non-github")
	assert.False(t, LocationIsGithub(""), "Allows empty location")
}
