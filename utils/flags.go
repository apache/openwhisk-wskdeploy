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

var Flags struct {
	WithinOpenWhisk bool   // is this running within an OpenWhisk action?
	ApiHost         string // OpenWhisk API host
	Auth            string // OpenWhisk API key
	Namespace       string
	ApiVersion      string // OpenWhisk version
	CfgFile         string
	CliVersion      string
	CliBuild        string
	Verbose         bool
	ProjectPath     string
	DeploymentPath  string
	ManifestPath    string
	UseDefaults     bool
	UseInteractive  bool

	//action flag definition
	//from go cli
	action struct {
		docker   bool
		copy     bool
		pipe     bool
		web      string
		sequence bool
		timeout  int
		memory   int
		logsize  int
		result   bool
		kind     string
		main     string
	}
}
