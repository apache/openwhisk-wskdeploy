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

package parsers

import (
	"fmt"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
)

// Test 1: validate manifest_parser:Unmarshal() method with a sample manifest in NodeJS
// validate that manifest_parser is able to read and parse the manifest data
func TestUnmarshalForHelloNodeJS(t *testing.T) {
	data := `
package:
  name: helloworld
  actions:
    helloNodejs:
      function: actions/hello.js
      runtime: nodejs:6`
	// set the zero value of struct ManifestYAML
	m := ManifestYAML{}
	// Unmarshal reads/parses manifest data and sets the values of ManifestYAML
	// And returns an error if parsing a manifest data fails
	err := NewYAMLParser().Unmarshal([]byte(data), &m)
	if err == nil {
		// ManifestYAML.Filepath does not get set by Parsers.Unmarshal
		// as it takes manifest YAML data as a function parameter
		// instead of file name of a manifest file, therefore there is
		// no way for Unmarshal function to set ManifestYAML.Filepath field
		// (TODO) Ideally we should change this functionality so that
		// (TODO) filepath is set to the actual path of the manifest file
		expectedResult := ""
		actualResult := m.Filepath
		assert.Equal(t, expectedResult, actualResult, "Expected filepath to be an empty"+
			" string instead its set to "+actualResult+" which is invalid value")
		// package name should be "helloworld"
		expectedResult = "helloworld"
		actualResult = m.Package.Packagename
		assert.Equal(t, expectedResult, actualResult, "Expected package name "+expectedResult+" but got "+actualResult)
		// manifest should contain only one action
		expectedResult = string(1)
		actualResult = string(len(m.Package.Actions))
		assert.Equal(t, expectedResult, actualResult, "Expected 1 but got "+actualResult)
		// get the action payload from the map of actions which is stored in
		// ManifestYAML.Package.Actions with the type of map[string]Action
		actionName := "helloNodejs"
		if action, ok := m.Package.Actions[actionName]; ok {
			// location/function of an action should be "actions/hello.js"
			expectedResult = "actions/hello.js"
			actualResult = action.Function
			assert.Equal(t, expectedResult, actualResult, "Expected action function "+expectedResult+" but got "+actualResult)
			// runtime of an action should be "nodejs:6"
			expectedResult = "nodejs:6"
			actualResult = action.Runtime
			assert.Equal(t, expectedResult, actualResult, "Expected action runtime "+expectedResult+" but got "+actualResult)
		} else {
			t.Error("Action named " + actionName + " does not exist.")
		}
	}
}

// Test 2: validate manifest_parser:Unmarshal() method with a sample manifest in Java
// validate that manifest_parser is able to read and parse the manifest data
func TestUnmarshalForHelloJava(t *testing.T) {
	data := `
package:
  name: helloworld
  actions:
    helloJava:
      function: actions/hello.jar
      runtime: java
      main: Hello`
	m := ManifestYAML{}
	err := NewYAMLParser().Unmarshal([]byte(data), &m)
	// nothing to test if Unmarshal returns an err
	if err == nil {
		// get an action from map of actions where key is action name and
		// value is Action struct
		actionName := "helloJava"
		if action, ok := m.Package.Actions[actionName]; ok {
			// runtime of an action should be java
			expectedResult := "java"
			actualResult := action.Runtime
			assert.Equal(t, expectedResult, actualResult, "Expected action runtime "+expectedResult+" but got "+actualResult)
			// Main field should be set to "Hello"
			expectedResult = action.Main
			actualResult = "Hello"
			assert.Equal(t, expectedResult, actualResult, "Expected action main function "+expectedResult+" but got "+actualResult)
		} else {
			t.Error("Expected action named " + actionName + " but does not exist.")
		}
	}
}

// Test 3: validate manifest_parser:Unmarshal() method with a sample manifest in Python
// validate that manifest_parser is able to read and parse the manifest data
func TestUnmarshalForHelloPython(t *testing.T) {
	data := `
package:
  name: helloworld
  actions:
    helloPython:
      function: actions/hello.py
      runtime: python`
	m := ManifestYAML{}
	err := NewYAMLParser().Unmarshal([]byte(data), &m)
	// nothing to test if Unmarshal returns an err
	if err == nil {
		// get an action from map of actions which is defined as map[string]Action{}
		actionName := "helloPython"
		if action, ok := m.Package.Actions[actionName]; ok {
			// runtime of an action should be python
			expectedResult := "python"
			actualResult := action.Runtime
			assert.Equal(t, expectedResult, actualResult, "Expected action runtime "+expectedResult+" but got "+actualResult)
		} else {
			t.Error("Expected action named " + actionName + " but does not exist.")
		}
	}
}

// Test 4: validate manifest_parser:Unmarshal() method with a sample manifest in Swift
// validate that manifest_parser is able to read and parse the manifest data
func TestUnmarshalForHelloSwift(t *testing.T) {
	data := `
package:
  name: helloworld
  actions:
    helloSwift:
      function: actions/hello.swift
      runtime: swift`
	m := ManifestYAML{}
	err := NewYAMLParser().Unmarshal([]byte(data), &m)
	// nothing to test if Unmarshal returns an err
	if err == nil {
		// get an action from map of actions which is defined as map[string]Action{}
		actionName := "helloSwift"
		if action, ok := m.Package.Actions[actionName]; ok {
			// runtime of an action should be swift
			expectedResult := "swift"
			actualResult := action.Runtime
			assert.Equal(t, expectedResult, actualResult, "Expected action runtime "+expectedResult+" but got "+actualResult)
		} else {
			t.Error("Expected action named " + actionName + " but does not exist.")
		}
	}
}

// Test 5: validate manifest_parser:Unmarshal() method for an action with parameters
// validate that manifest_parser is able to read and parse the manifest data, specially
// validate two input parameters and their values
func TestUnmarshalForHelloWithParams(t *testing.T) {
	var data = `
package:
   name: helloworld
   actions:
     helloWithParams:
       function: actions/hello-with-params.js
       runtime: nodejs:6
       inputs:
         name: Amy
         place: Paris`
	m := ManifestYAML{}
	err := NewYAMLParser().Unmarshal([]byte(data), &m)
	if err == nil {
		actionName := "helloWithParams"
		if action, ok := m.Package.Actions[actionName]; ok {
			expectedResult := "Amy"
			actualResult := action.Inputs["name"].Value.(string)
			assert.Equal(t, expectedResult, actualResult,
				"Expected input parameter "+expectedResult+" but got "+actualResult+"for name")
			expectedResult = "Paris"
			actualResult = action.Inputs["place"].Value.(string)
			assert.Equal(t, expectedResult, actualResult,
				"Expected input parameter "+expectedResult+" but got "+actualResult+"for place")
		}
	}
}

// Test 6: validate manifest_parser:Unmarshal() method for an invalid manifest
// manifest_parser should report an error when a package section is missing
func TestUnmarshalForMissingPackage(t *testing.T) {
	data := `
  actions:
    helloNodejs:
      function: actions/hello.js
      runtime: nodejs:6
    helloJava:
      function: actions/hello.java`
	// set the zero value of struct ManifestYAML
	m := ManifestYAML{}
	// Unmarshal reads/parses manifest data and sets the values of ManifestYAML
	// And returns an error if parsing a manifest data fails
	err := NewYAMLParser().Unmarshal([]byte(data), &m)
	fmt.Println("Error: ", err)
	fmt.Println("Filepath: \"", m.Filepath, "\"")
	fmt.Println("Package: ", m.Package)
	fmt.Println("PackageName: \"", m.Package.Packagename, "\"")
	fmt.Println("Number of Actions: ", len(m.Package.Actions))
	fmt.Println("Actions: ", m.Package.Actions)
	// (TODO) Unmarshal does not report any error even if manifest file is missing required section.
	// (TODO) In this test case, "Package" section is missing which is not reported,
	// (TODO) instead ManifestYAML is set to its zero values
	// assert.NotNil(t, err, "Expected some error from Unmarshal but got no error")
}

/*
 Test 7: validate manifest_parser:ParseManifest() method for multiline parameters
 manifest_parser should be able to parse all different mutliline combinations of
 inputs section including:

 case 1: value only
 param:
    value: <value>
 case 2: type only
 param:
    type: <type>
 case 3: type and value only
 param:
    type: <type>
    value: <value>
 case 4: default value
 param:
    type: <type>
    default: <default value>
*/
func TestParseManifestForMultiLineParams(t *testing.T) {
	// manifest file is located under ../tests folder
	manifestFile := "../tests/dat/manifest_validate_multiline_params.yaml"
	// read and parse manifest.yaml file
	m := NewYAMLParser().ParseManifest(manifestFile)

	// validate package name should be "validate"
	expectedPackageName := m.Package.Packagename
	actualPackageName := "validate"
	assert.Equal(t, expectedPackageName, actualPackageName,
		"Expected "+expectedPackageName+" but got "+actualPackageName)

	// validate this package contains one action
	expectedActionsCount := 1
	actualActionsCount := len(m.Package.Actions)
	assert.Equal(t, expectedActionsCount, actualActionsCount,
		"Expected "+string(expectedActionsCount)+" but got "+string(actualActionsCount))

	// here Package.Actions holds a map of map[string]Action
	// where string is the action name so in case you create two actions with
	// same name, will go unnoticed
	// also, the Action struct does not have name field set it to action name
	actionName := "validate_multiline_params"
	if action, ok := m.Package.Actions[actionName]; ok {
		// validate location/function of an action to be "actions/dump_params.js"
		expectedResult := "actions/dump_params.js"
		actualResult := action.Function
		assert.Equal(t, expectedResult, actualResult, "Expected action function "+expectedResult+" but got "+actualResult)

		// validate runtime of an action to be "nodejs:6"
		expectedResult = "nodejs:6"
		actualResult = action.Runtime
		assert.Equal(t, expectedResult, actualResult, "Expected action runtime "+expectedResult+" but got "+actualResult)

		// validate the number of inputs to this action
		expectedResult = strconv.FormatInt(10, 10)
		actualResult = strconv.FormatInt(int64(len(action.Inputs)), 10)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// validate inputs to this action
		for input, param := range action.Inputs {
			switch input {
			case "param_string_value_only":
				expectedResult = "foo"
				actualResult = param.Value.(string)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_int_value_only":
				expectedResult = strconv.FormatInt(123, 10)
				actualResult = strconv.FormatInt(int64(param.Value.(int)), 10)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_float_value_only":
				expectedResult = strconv.FormatFloat(3.14, 'f', -1, 64)
				actualResult = strconv.FormatFloat(param.Value.(float64), 'f', -1, 64)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_string_type_and_value_only":
				expectedResult = "foo"
				actualResult = param.Value.(string)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
				expectedResult = "string"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_string_type_only":
				expectedResult = "string"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_integer_type_only":
				expectedResult = "integer"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_float_type_only":
				expectedResult = "float"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_string_with_default":
				expectedResult = "string"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
				expectedResult = "bar"
				actualResult = param.Default.(string)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_integer_with_default":
				expectedResult = "integer"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
				expectedResult = strconv.FormatInt(-1, 10)
				actualResult = strconv.FormatInt(int64(param.Default.(int)), 10)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_float_with_default":
				expectedResult = "float"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
				expectedResult = strconv.FormatFloat(2.9, 'f', -1, 64)
				actualResult = strconv.FormatFloat(param.Default.(float64), 'f', -1, 64)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			}
		}

		// validate outputs
		// output payload is of type string and has a description
		if payload, ok := action.Outputs["payload"]; ok {
			p := payload.(map[interface{}]interface{})
			expectedResult = "string"
			actualResult = p["type"].(string)
			assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			expectedResult = "parameter dump"
			actualResult = p["description"].(string)
			assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
		}
	}
}

// Test 8: validate manifest_parser:ParseManifest() method for single line parameters
// manifest_parser should be able to parse input section with different types of values
func TestParseManifestForSingleLineParams(t *testing.T) {
	// manifest file is located under ../tests folder
	manifestFile := "../tests/dat/manifest_validate_singleline_params.yaml"
	// read and parse manifest.yaml file
	m := NewYAMLParser().ParseManifest(manifestFile)

	// validate package name should be "validate"
	expectedPackageName := m.Package.Packagename
	actualPackageName := "validate"
	assert.Equal(t, expectedPackageName, actualPackageName,
		"Expected "+expectedPackageName+" but got "+actualPackageName)

	// validate this package contains one action
	expectedActionsCount := 1
	actualActionsCount := len(m.Package.Actions)
	assert.Equal(t, expectedActionsCount, actualActionsCount,
		"Expected "+string(expectedActionsCount)+" but got "+string(actualActionsCount))

	actionName := "validate_singleline_params"
	if action, ok := m.Package.Actions[actionName]; ok {
		// validate location/function of an action to be "actions/dump_params.js"
		expectedResult := "actions/dump_params.js"
		actualResult := action.Function
		assert.Equal(t, expectedResult, actualResult, "Expected action function "+expectedResult+" but got "+actualResult)

		// validate runtime of an action to be "nodejs:6"
		expectedResult = "nodejs:6"
		actualResult = action.Runtime
		assert.Equal(t, expectedResult, actualResult, "Expected action runtime "+expectedResult+" but got "+actualResult)

		// validate the number of inputs to this action
		expectedResult = strconv.FormatInt(17, 10)
		actualResult = strconv.FormatInt(int64(len(action.Inputs)), 10)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// validate inputs to this action
		for input, param := range action.Inputs {
			switch input {
			case "param_simple_string":
				expectedResult = "foo"
				actualResult = param.Value.(string)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_simple_integer_1":
				expectedResult = strconv.FormatInt(1, 10)
				actualResult = strconv.FormatInt(int64(param.Value.(int)), 10)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_simple_integer_2":
				expectedResult = strconv.FormatInt(0, 10)
				actualResult = strconv.FormatInt(int64(param.Value.(int)), 10)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_simple_integer_3":
				expectedResult = strconv.FormatInt(-1, 10)
				actualResult = strconv.FormatInt(int64(param.Value.(int)), 10)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_simple_integer_4":
				expectedResult = strconv.FormatInt(99999, 10)
				actualResult = strconv.FormatInt(int64(param.Value.(int)), 10)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_simple_integer_5":
				expectedResult = strconv.FormatInt(-99999, 10)
				actualResult = strconv.FormatInt(int64(param.Value.(int)), 10)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_simple_float_1":
				expectedResult = strconv.FormatFloat(1.1, 'f', -1, 64)
				actualResult = strconv.FormatFloat(param.Value.(float64), 'f', -1, 64)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_simple_float_2":
				expectedResult = strconv.FormatFloat(0.0, 'f', -1, 64)
				actualResult = strconv.FormatFloat(param.Value.(float64), 'f', -1, 64)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_simple_float_3":
				expectedResult = strconv.FormatFloat(-1.1, 'f', -1, 64)
				actualResult = strconv.FormatFloat(param.Value.(float64), 'f', -1, 64)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_simple_env_var_1":
				expectedResult = "$GOPATH"
				actualResult = param.Value.(string)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_simple_invalid_env_var":
				expectedResult = "$DollarSignNotInEnv"
				actualResult = param.Value.(string)
				assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			case "param_simple_implied_empty":
				assert.Nil(t, param.Value, "Expected nil")
			case "param_simple_explicit_empty_1":
				actualResult = param.Value.(string)
				assert.Empty(t, actualResult, "Expected empty string but got "+actualResult)
			case "param_simple_explicit_empty_2":
				actualResult = param.Value.(string)
				assert.Empty(t, actualResult, "Expected empty string but got "+actualResult)
			}
		}

		// validate outputs
		// output payload is of type string and has a description
		if payload, ok := action.Outputs["payload"]; ok {
			p := payload.(map[interface{}]interface{})
			expectedResult = "string"
			actualResult = p["type"].(string)
			assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
			expectedResult = "parameter dump"
			actualResult = p["description"].(string)
			assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
		}
	}
}

// Test 9: validate manifest_parser.ComposeActions() method for implicit runtimes
// when a runtime of an action is not provided, manifest_parser determines the runtime
// based on the file extension of an action file
func TestComposeActionsForImplicitRuntimes(t *testing.T) {
	data :=
		`package:
  name: helloworld
  actions:
    helloNodejs:
      function: ../tests/usecases/helloworld/actions/hello.js
    helloJava:
      function: ../tests/usecases/helloworld/actions/hello.jar
      main: Hello
    helloPython:
      function: ../tests/usecases/helloworld/actions/hello.py
    helloSwift:
      function: ../tests/usecases/helloworld/actions/hello.swift`

	dir, _ := os.Getwd()
	tmpfile, err := ioutil.TempFile(dir, "manifest_parser_validate_runtimes_")
	if err == nil {
		defer os.Remove(tmpfile.Name()) // clean up
		if _, err := tmpfile.Write([]byte(data)); err == nil {
			// read and parse manifest.yaml file
			p := NewYAMLParser()
			m := p.ParseManifest(tmpfile.Name())
			actions, _, err := p.ComposeActions(m, tmpfile.Name())
			var expectedResult string
			if err == nil {
				for i := 0; i < len(actions); i++ {
					if actions[i].Action.Name == "helloNodejs" {
						expectedResult = "nodejs:default"
						// (TODO) change expectedResult in the following condition
						// (TODO) once issue #306 is fixed as runtime is set to
						// (TODO) nodejs:default for any kind of action file except
						// (TODO) files with extension .js, .py, and .swift
					} else if actions[i].Action.Name == "helloJava" {
						expectedResult = "nodejs:default"
					} else if actions[i].Action.Name == "helloPython" {
						expectedResult = "python"
					} else if actions[i].Action.Name == "helloSwift" {
						expectedResult = "swift:default"
					}
					actualResult := actions[i].Action.Exec.Kind
					assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
				}
			}

		}
		tmpfile.Close()
	}
}

// Test 10: validate manifest_parser.ComposeActions() method for invalid runtimes
// when a runtime of an action is set to some garbage, manifest_parser should
// report an error for that action
func TestComposeActionsForInvalidRuntime(t *testing.T) {
	data :=
		`package:
   name: helloworld
   actions:
     helloInvalidRuntime:
       function: ../tests/usecases/helloworld/actions/hello.js
       runtime: invalid`
	dir, _ := os.Getwd()
	tmpfile, err := ioutil.TempFile(dir, "manifest_parser_validate_runtime_")
	if err == nil {
		defer os.Remove(tmpfile.Name()) // clean up
		if _, err := tmpfile.Write([]byte(data)); err == nil {
			// read and parse manifest.yaml file
			p := NewYAMLParser()
			m := p.ParseManifest(tmpfile.Name())
			_, _, err := p.ComposeActions(m, tmpfile.Name())
			// (TODO) uncomment the following test case after issue #307 is fixed
			// (TODO) its failing right now as we are lacking check on invalid runtime
			// assert.NotNil(t, err, "Invalid runtime, ComposeActions should report an error")
			// (TODO) remove this print statement after uncommenting above test case
			fmt.Println(err)
		}
		tmpfile.Close()
	}
}

// Test 11: validate manfiest_parser.ComposeActions() method for single line parameters
// manifest_parser should be able to parse input section with different types of values
func TestComposeActionsForSingleLineParams(t *testing.T) {
	// manifest file is located under ../tests folder
	manifestFile := "../tests/dat/manifest_validate_singleline_params.yaml"
	// read and parse manifest.yaml file
	p := NewYAMLParser()
	m := p.ParseManifest(manifestFile)
	actions, _, err := p.ComposeActions(m, manifestFile)

	if err == nil {
		// assert that the actions variable has only one action
		assert.Equal(t, 1, len(actions), "We have defined only one action but we got "+string(len(actions)))

		action := actions[0]

		// param_simple_string should value "foo"
		expectedResult := "foo"
		actualResult := action.Action.Parameters.GetValue("param_simple_string").(string)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_simple_integer_1 should have value 1
		expectedResult = strconv.FormatInt(1, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue("param_simple_integer_1").(int)), 10)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_simple_integer_2 should have value 0
		expectedResult = strconv.FormatInt(0, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue("param_simple_integer_2").(int)), 10)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_simple_integer_3 should have value -1
		expectedResult = strconv.FormatInt(-1, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue("param_simple_integer_3").(int)), 10)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_simple_integer_4 should have value 99999
		expectedResult = strconv.FormatInt(99999, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue("param_simple_integer_4").(int)), 10)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_simple_integer_5 should have value -99999
		expectedResult = strconv.FormatInt(-99999, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue("param_simple_integer_5").(int)), 10)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_simple_float_1 should have value 1.1
		expectedResult = strconv.FormatFloat(1.1, 'f', -1, 64)
		actualResult = strconv.FormatFloat(action.Action.Parameters.GetValue("param_simple_float_1").(float64), 'f', -1, 64)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_simple_float_2 should have value 0.0
		expectedResult = strconv.FormatFloat(0.0, 'f', -1, 64)
		actualResult = strconv.FormatFloat(action.Action.Parameters.GetValue("param_simple_float_2").(float64), 'f', -1, 64)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_simple_float_3 should have value -1.1
		expectedResult = strconv.FormatFloat(-1.1, 'f', -1, 64)
		actualResult = strconv.FormatFloat(action.Action.Parameters.GetValue("param_simple_float_3").(float64), 'f', -1, 64)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_simple_env_var_1 should have value of env. variable $GOPATH
		expectedResult = os.Getenv("GOPATH")
		actualResult = action.Action.Parameters.GetValue("param_simple_env_var_1").(string)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_simple_invalid_env_var should have value of ""
		expectedResult = ""
		actualResult = action.Action.Parameters.GetValue("param_simple_invalid_env_var").(string)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_simple_implied_empty should be ""
		actualResult = action.Action.Parameters.GetValue("param_simple_implied_empty").(string)
		assert.Empty(t, actualResult, "Expected empty string but got "+actualResult)

		// param_simple_explicit_empty_1 should be ""
		actualResult = action.Action.Parameters.GetValue("param_simple_explicit_empty_1").(string)
		assert.Empty(t, actualResult, "Expected empty string but got "+actualResult)

		// param_simple_explicit_empty_2 should be ""
		actualResult = action.Action.Parameters.GetValue("param_simple_explicit_empty_2").(string)
		assert.Empty(t, actualResult, "Expected empty string but got "+actualResult)
	}
}

// Test 12: validate manfiest_parser.ComposeActions() method for multi line parameters
// manifest_parser should be able to parse input section with different types of values
func TestComposeActionsForMultiLineParams(t *testing.T) {
	// manifest file is located under ../tests folder
	manifestFile := "../tests/dat/manifest_validate_multiline_params.yaml"
	// read and parse manifest.yaml file
	p := NewYAMLParser()
	m := p.ParseManifest(manifestFile)
	actions, _, err := p.ComposeActions(m, manifestFile)

	if err == nil {
		// assert that the actions variable has only one action
		assert.Equal(t, 1, len(actions), "We have defined only one action but we got "+string(len(actions)))

		action := actions[0]

		fmt.Println(action.Action.Parameters)

		// param_string_value_only should be "foo"
		expectedResult := "foo"
		actualResult := action.Action.Parameters.GetValue("param_string_value_only").(string)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_int_value_only should be 123
		expectedResult = strconv.FormatInt(123, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue("param_int_value_only").(int)), 10)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_float_value_only should be 3.14
		expectedResult = strconv.FormatFloat(3.14, 'f', -1, 64)
		actualResult = strconv.FormatFloat(action.Action.Parameters.GetValue("param_float_value_only").(float64), 'f', -1, 64)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_string_type_and_value_only should be foo
		expectedResult = "foo"
		actualResult = action.Action.Parameters.GetValue("param_string_type_and_value_only").(string)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_string_type_only should be ""
		actualResult = action.Action.Parameters.GetValue("param_string_type_only").(string)
		assert.Empty(t, actualResult, "Expected empty string but got "+actualResult)

		// param_integer_type_only should be 0
		expectedResult = strconv.FormatInt(0, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue("param_integer_type_only").(int)), 10)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_float_type_only should be 0
		expectedResult = strconv.FormatFloat(0.0, 'f', -1, 64)
		actualResult = strconv.FormatFloat(action.Action.Parameters.GetValue("param_float_type_only").(float64), 'f', -1, 64)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_string_with_default should be "bar"
		expectedResult = "bar"
		actualResult = action.Action.Parameters.GetValue("param_string_with_default").(string)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_integer_with_default should be -1
		expectedResult = strconv.FormatInt(-1, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue("param_integer_with_default").(int)), 10)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)

		// param_float_with_default should be 2.9
		expectedResult = strconv.FormatFloat(2.9, 'f', -1, 64)
		actualResult = strconv.FormatFloat(action.Action.Parameters.GetValue("param_float_with_default").(float64), 'f', -1, 64)
		assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
	}
}

// Test 13: validate manfiest_parser.ComposeActions() method
func TestComposeActionsForFunction(t *testing.T) {
	data :=
		`package:
  name: helloworld
  actions:
    hello1:
      function: ../tests/usecases/helloworld/actions/hello.js`
	// (TODO) uncomment this after we add support for action file content from URL
	// hello2:
	//  function: https://raw.githubusercontent.com/apache/incubator-openwhisk-wskdeploy/master/tests/usecases/helloworld/manifest.yaml`
	dir, _ := os.Getwd()
	tmpfile, err := ioutil.TempFile(dir, "manifest_parser_validate_locations_")
	if err == nil {
		defer os.Remove(tmpfile.Name()) // clean up
		if _, err := tmpfile.Write([]byte(data)); err == nil {
			// read and parse manifest.yaml file
			p := NewYAMLParser()
			m := p.ParseManifest(tmpfile.Name())
			actions, _, err := p.ComposeActions(m, tmpfile.Name())
			var expectedResult, actualResult string
			if err == nil {
				for i := 0; i < len(actions); i++ {
					if actions[i].Action.Name == "hello1" {
						expectedResult, _ = filepath.Abs("../tests/usecases/helloworld/actions/hello.js")
						actualResult, _ = filepath.Abs(actions[i].Filepath)
						assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
						// (TODO) Uncomment the following condition, hello2
						// (TODO) after issue # 311 is fixed
						//} else if actions[i].Action.Name == "hello2" {
						//  assert.NotNil(t, actions[i].Action.Exec.Code, "Expected source code from an action file but found it empty")
					}
				}
			}

		}
		tmpfile.Close()
	}

}

// (TODO) uncomment this test case after issue # 312 is fixed
// Test 14: validate manfiest_parser.ComposeActions() method
//func TestComposeActionsForLimits (t *testing.T) {
//  data :=
//`package:
//  name: helloworld
//  actions:
//    hello1:
//      function: ../tests/usecases/helloworld/actions/hello.js
//    hello2:
//      function: ../tests/usecases/helloworld/actions/hello.js
//      limits:
//        timeout: 60
//        memorySize: 128
//        logSize: 1
//        concurrentActivations: 10
//        userInvocationRate: 50
//        codeSize: 1024
//        parameterSize: 128`
//  dir, _ := os.Getwd()
//  tmpfile, err := ioutil.TempFile(dir, "manifest_parser_validate_limits_")
//  if err == nil {
//      defer os.Remove(tmpfile.Name()) // clean up
//      if _, err := tmpfile.Write([]byte(data)); err == nil {
//          // read and parse manifest.yaml file
//          p := NewYAMLParser()
//          m := p.ParseManifest(tmpfile.Name())
//          actions, _, err := p.ComposeActions(m, tmpfile.Name())
//          //var expectedResult, actualResult string
//          if err == nil {
//              for i:=0; i<len(actions); i++ {
//                  if actions[i].Action.Name == "hello1" {
//                      fmt.Println(actions[i].Action.Limits)
//                      //assert.Nil(t, actions[i].Action.Limits, "Expected limit section to be empty but got " + actions[i].Action.Limits)
//                  } else if actions[i].Action.Name == "hello2" {
//                      fmt.Println(actions[i].Action.Limits)
//                      //assert.NotNil(t, actions[i].Action.Limits, "Expected limit section to be not empty but found it empty")
//                  }
//              }
//          }
//
//      }
//      tmpfile.Close()
//  }
//}

// Test 15: validate manfiest_parser.ComposeActions() method
func TestComposeActionsForWebActions(t *testing.T) {
	data :=
		`package:
  name: helloworld
  actions:
    hello:
      function: ../tests/usecases/helloworld/actions/hello.js
      annotations:
        foo: bar
      web-export: true`
	// (TODO) remove annotations from manifest file after we fix issue # 313
	dir, _ := os.Getwd()
	tmpfile, err := ioutil.TempFile(dir, "manifest_parser_validate_web_actions_")
	if err == nil {
		defer os.Remove(tmpfile.Name()) // clean up
		if _, err := tmpfile.Write([]byte(data)); err == nil {
			// read and parse manifest.yaml file
			p := NewYAMLParser()
			m := p.ParseManifest(tmpfile.Name())
			actions, _, err := p.ComposeActions(m, tmpfile.Name())
			if err == nil {
				for i := 0; i < len(actions); i++ {
					if actions[i].Action.Name == "hello" {
						for _, a := range actions[i].Action.Annotations {
							switch a.Key {
							case "web-export":
								assert.Equal(t, true, a.Value, "Expected true for web-export but got "+strconv.FormatBool(a.Value.(bool)))
							case "raw-http":
								assert.Equal(t, false, a.Value, "Expected false for raw-http but got "+strconv.FormatBool(a.Value.(bool)))
							case "final":
								assert.Equal(t, true, a.Value, "Expected true for final but got "+strconv.FormatBool(a.Value.(bool)))
							}
						}
					}
				}
			}

		}
		tmpfile.Close()
	}
}

// Test 16: validate manifest_parser.ResolveParameter() method
func TestResolveParameterForMultiLineParams(t *testing.T) {
	p := "name"
	v := "foo"
	y := reflect.TypeOf(v).Name() // y := string
	d := "default_name"

	// type string - value only param
	param1 := Parameter{Value: v, multiline: true}
	r1, _ := ResolveParameter(p, &param1)
	assert.Equal(t, v, r1, "Expected value "+v+" but got "+r1.(string))
	assert.IsType(t, v, r1, "Expected parameter %v of type %T but found %T", p, v, r1)

	// type string - type and value only param
	param2 := Parameter{Type: y, Value: v, multiline: true}
	r2, _ := ResolveParameter(p, &param2)
	assert.Equal(t, v, r2, "Expected value "+v+" but got "+r2.(string))
	assert.IsType(t, v, r2, "Expected parameter %v of type %T but found %T", p, v, r2)

	// type string - type, no value, but default value param
	param3 := Parameter{Type: y, Default: d, multiline: true}
	r3, _ := ResolveParameter(p, &param3)
	assert.Equal(t, d, r3, "Expected value "+d+" but got "+r3.(string))
	assert.IsType(t, d, r3, "Expected parameter %v of type %T but found %T", p, d, r3)

	// type string - type and value only param
	// type is "string" and value is of type "int"
	// ResolveParameter matches specified type with the type of the specified value
	// it fails if both types don't match
	// ResolveParameter determines type from the specified value
	// in this case, ResolveParameter returns value of type int
	v1 := 11
	param4 := Parameter{Type: y, Value: v1, multiline: true}
	r4, _ := ResolveParameter(p, &param4)
	assert.Equal(t, v1, r4, "Expected value "+strconv.FormatInt(int64(v1), 10)+" but got "+strconv.FormatInt(int64(r4.(int)), 10))
	assert.IsType(t, v1, r4, "Expected parameter %v of type %T but found %T", p, v1, r4)

	// type invalid - type only param
	param5 := Parameter{Type: "invalid", multiline: true}
	_, err := ResolveParameter(p, &param5)
	assert.NotNil(t, err, "Expected error saying Invalid type for parameter")
	expectedErr := utils.NewParserErr("", -1, "Invalid Type for parameter. [invalid]")
	assert.Equal(t, err, expectedErr, "Expected error "+expectedErr.Error()+" but found "+err.Error())

	// type none - param without type, without value, and without default value
	param6 := Parameter{multiline: true}
	r6, _ := ResolveParameter("none", &param6)
	assert.Empty(t, r6, "Expected default value of empty string but found "+r6.(string))

}
