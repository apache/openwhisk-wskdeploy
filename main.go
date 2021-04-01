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

package main

import (
	"github.com/apache/openwhisk-wskdeploy/cmd"
	"github.com/apache/openwhisk-wskdeploy/utils"
)

func main() {
	cmd.Execute()
}

// Struct used to hold tagged (release) build information
// Which is displayed by the `version` command.
// Added automatically for CI/CD Travis builds in GitHub
var (
	// Apache OpenWhisk Whisk Deploy release/build version
	Version = "unset"
	// Git source code commit # that initiated the build
	GitCommit = "unset"
	// Date stamp indicating when build was initiated
	BuildDate = "unset"
)

func init() {
	utils.Flags.CliVersion = Version
	utils.Flags.CliGitCommit = GitCommit
	utils.Flags.CliBuildDate = BuildDate
}
