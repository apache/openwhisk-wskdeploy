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
	"strings"
	"os"
	"reflect"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
)

// Test if a string
func isValidEnvironmentVar(value string) bool {

	// A valid Env. variable should start with or contain '$' (dollar) char.
	//
	// If the value is a single Env. variable, it should start with a '$' (dollar) char
	// and have at least 1 additional character after it, e.g. $ENV_VAR
	// If the value is a concatenation of a string and a Env. variable, it should contain '$' (dollar)
	// and have a string following which is surrounded with '{' and '}', e.g. xxx${ENV_VAR}xxx.
	if value != "" && strings.HasPrefix(value, "$") && len(value) > 1 {
		return true
	}
	if value != "" && strings.Contains(value, "${") && strings.Count(value, "{") == strings.Count(value, "}") {
		return true
	}
	return false
}

// Get the env variable value by key.
// Get the env variable if the key is start by $
func GetEnvVar(key interface{}) interface{} {
	// Assure the key itself is not nil
	if key == nil {
		return nil
	}

	if reflect.TypeOf(key).String() == "string" {
		keystr := key.(string)
		if isValidEnvironmentVar(keystr) {
			// retrieve the value of the env. var. from the host system.
			var thisValue string
			//split the string with ${}
			//test if the substr is a environment var
			//if it is, replace it with the value
			f := func(c rune) bool {
				return c == '$' || c == '{' || c == '}'
			}
			for _, substr := range strings.FieldsFunc(keystr, f) {
				//if the substr is a $ENV_VAR
				if strings.Contains(keystr, "$"+substr) {
					thisValue = os.Getenv(substr)
					if thisValue == "" {
						wskprint.PrintlnOpenWhiskOutput("WARNING: Missing Environment Variable " + substr + ".")
					}
					keystr = strings.Replace(keystr, "$"+substr, thisValue, -1)
					//if the substr is a ${ENV_VAR}
				} else if strings.Contains(keystr, "${"+substr+"}") {
					thisValue = os.Getenv(substr)
					if thisValue == "" {
						wskprint.PrintlnOpenWhiskOutput("WARNING: Missing Environment Variable " + substr + ".")
					}
					keystr = strings.Replace(keystr, "${"+substr+"}", thisValue, -1)
				}
			}
			return keystr
		}

		// The key was not a valid env. variable, simply return it as the value itself (of type string)
		return keystr
	}
	return key
}
