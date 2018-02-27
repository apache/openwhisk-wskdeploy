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
	"fmt"
	"reflect"
)

type WskDeployFlags struct {
	WithinOpenWhisk  bool   // is this running within an OpenWhisk action?
	ApiHost          string // OpenWhisk API host
	Auth             string // OpenWhisk API key
	Namespace        string
	ApiVersion       string // OpenWhisk version
	CfgFile          string
	CliVersion       string
	CliBuild         string
	Verbose          bool
	ProjectPath      string
	DeploymentPath   string
	ManifestPath     string
	UseInteractive   bool
	Strict           bool // strict flag to support user defined runtime version.
	Key              string
	Cert             string
	Managed          bool // OpenWhisk Managed Deployments
	ApigwAccessToken string
}

func (flags *WskDeployFlags) Format() string {

	flagNames := reflect.TypeOf(*flags)
	flagValues := reflect.ValueOf(*flags)

	var name string
	var value interface{}
	//var t interface{}
	var result string

	for i := 0; i < flagValues.NumField(); i++ {
		name = flagNames.Field(i).Name
		value = flagValues.Field(i)
		// NOTE: if you need to see the Type, add this line to output
		//t = flagValues.Field(i).Type()
		line := fmt.Sprintf("      > %s: [%v]\n", name, value)
		result += line
	}

	return result
}

var Flags WskDeployFlags
