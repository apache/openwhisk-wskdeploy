/*
 * Copyright 2015-2016 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

// constants
const ServerlessBinaryCommand = "serverless"

// ManifestProvder is a provider description in the manifest
type ManifestProvider struct {
	Name    string
	Runtime string
}

// Manifest is the main manifest file
type Manifest struct {
	Service  string
	Provider ManifestProvider
}

// errors
type serverlessErr struct {
	msg string
}

func (e *serverlessErr) Error() string {
	return e.msg
}

// util functions
func check(e error) {
	if e != nil {
		panic(e)
	}
}
