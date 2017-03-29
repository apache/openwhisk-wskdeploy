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

// shared.go
package cmdImp

// name of manifest and deployment files
const ManifestFileNameYaml = "manifest.yaml"
const ManifestFileNameYml = "manifest.yml"
const DeploymentFileNameYaml = "deployment.yaml"
const DeploymentFileNameYml = "deployment.yml"

var CfgFile string
var CliVersion string
var CliBuild string

// used to configure service deployer for various commands
// TODO: should move this into utils.Flags
var Verbose bool
var ProjectPath string
var DeploymentPath string
var ManifestPath string
var UseDefaults bool
var UseInteractive bool
