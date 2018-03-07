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
	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskderrors"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"github.com/davecgh/go-spew/spew"
)

const (
	// local test assert messages
	TEST_MSG_PACKAGE_NAME_MISSING                   = "Package named [%s] missing."
	TEST_MSG_ACTION_NUMBER_MISMATCH                 = "Number of Actions mismatched."
	TEST_MSG_ACTION_NAME_MISSING                    = "Action named [%s] does not exist."
	TEST_MSG_ACTION_FUNCTION_PATH_MISMATCH          = "Action function path mismatched."
	TEST_MSG_ACTION_FUNCTION_RUNTIME_MISMATCH       = "Action function runtime mismatched."
	TEST_MSG_ACTION_FUNCTION_MAIN_MISMATCH          = "Action function main name mismatch."
	TEST_MSG_ACTION_PARAMETER_TYPE_MISMATCH         = "Action parameter [%s] had a type mismatch."
	TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH        = "Action parameter [%s] had a value mismatch."
	TEST_MSG_PARAMETER_NUMBER_MISMATCH              = "Number of Paramaters mismatched."
	TEST_MSG_MANIFEST_UNMARSHALL_ERROR_EXPECTED     = "Manifest [%s]: Expected Unmarshal error."
	TEST_MSG_ACTION_FUNCTION_RUNTIME_ERROR_EXPECTED = "Manifest [%s]: Expected runtime error."

	// local error messages
	TEST_ERROR_MANIFEST_PARSE_FAILURE   = "Manifest [%s]: Failed to parse."
	TEST_ERROR_MANIFEST_READ_FAILURE    = "Manifest [%s]: Failed to ReadFile()."
	TEST_ERROR_MANIFEST_DATA_UNMARSHALL = "Manifest [%s]: Failed to Unmarshall manifest."
)

func init() {
	op, error := utils.ParseOpenWhisk("")
	if error == nil {
		utils.SupportedRunTimes = utils.ConvertToMap(op)
		utils.DefaultRunTimes = utils.DefaultRuntimes(op)
		utils.FileExtensionRuntimeKindMap = utils.FileExtensionRuntimes(op)
	}
}

func testLoadParseManifest(t *testing.T, manifestFile string) (*YAMLParser, *YAML, error) {
	// read and parse manifest.yaml file located under ../tests folder
	p := NewYAMLParser()
	m, err := p.ParseManifest(manifestFile)
	if err != nil {
		assert.Fail(t, fmt.Sprintf(TEST_ERROR_MANIFEST_PARSE_FAILURE, manifestFile))
	}

	return p, m, err
}

func testReadAndUnmarshalManifest(t *testing.T, pathManifest string) (YAML, error) {
	// Init YAML struct and attempt to Unmarshal YAML byte[] data
	m := YAML{}

	// read raw bytes of manifest.yaml file
	data, err := ioutil.ReadFile(pathManifest)

	if err != nil {
		t.Error(fmt.Sprintf(TEST_ERROR_MANIFEST_READ_FAILURE, pathManifest))
		return m, err
	}

	err = NewYAMLParser().Unmarshal([]byte(data), &m)
	return m, err
}

/*
   testUnmarshalManifestAndActionBasic

   This function validates basic Manifest Package and Action keys including
   - Package name mismatch (single "package" only)
   - Number of Actions mismatch
   - Action Function path mismatch
   - Action runtime (name) mismatch

   and optionally,
   = Action function "main" name mismatch

   Returns:
   - N/A
*/
func testUnmarshalManifestPackageAndActionBasic(t *testing.T,
	pathManifest string,
	namePackage string,
	numActions int,
	nameAction string,
	pathFunction string,
	nameRuntime string,
	nameMain string) (YAML, *Package, error) {

	// Test that we are able to read the manifest file and unmarshall into YAML struct
	m, err := testReadAndUnmarshalManifest(t, pathManifest)

	// nothing to test if Unmarshal returns an err
	if err != nil {
		assert.Fail(t, fmt.Sprintf(TEST_ERROR_MANIFEST_DATA_UNMARSHALL, pathManifest))
	} else {
		// test package (name) exists
		if pkg, ok := m.Packages[namePackage]; ok {

			// test # of actions in manifest
			expectedActionsCount := numActions
			actualActionsCount := len(pkg.Actions)
			assert.Equal(t, expectedActionsCount, actualActionsCount, TEST_MSG_ACTION_NUMBER_MISMATCH)

			// get an action from map of actions where key is action name and value is Action struct
			if action, ok := pkg.Actions[nameAction]; ok {

				// test action's function path
				assert.Equal(t, pathFunction, action.Function, TEST_MSG_ACTION_FUNCTION_PATH_MISMATCH)

				// test action's runtime
				assert.Equal(t, nameRuntime, action.Runtime, TEST_MSG_ACTION_FUNCTION_RUNTIME_MISMATCH)

				// test action's "Main" function
				if nameMain != "" {
					assert.Equal(t, nameMain, action.Main, TEST_MSG_ACTION_FUNCTION_MAIN_MISMATCH)
				}

				return m, &pkg, err

			} else {
				t.Error(fmt.Sprintf(TEST_MSG_ACTION_NAME_MISSING, nameAction))
			}

		} else {
			assert.Fail(t, fmt.Sprintf(TEST_MSG_PACKAGE_NAME_MISSING, namePackage))
		}
	}

	return m, nil, nil
}

func testUnmarshalTemporaryFile(data []byte, filename string) (p *YAMLParser, m *YAML, t string) {
	dir, _ := os.Getwd()
	tmpfile, err := ioutil.TempFile(dir, filename)
	if err == nil {
		defer os.Remove(tmpfile.Name()) // clean up
		if _, err := tmpfile.Write(data); err == nil {
			// read and parse manifest.yaml file
			p = NewYAMLParser()
			m, _ = p.ParseManifest(tmpfile.Name())
		}
	}
	t = tmpfile.Name()
	tmpfile.Close()
	return
}

// Test 1: validate manifest_parser:Unmarshal() method with a sample manifest in NodeJS
// validate that manifest_parser is able to read and parse the manifest data
func TestUnmarshalForHelloNodeJS(t *testing.T) {
	testUnmarshalManifestPackageAndActionBasic(t,
		"../tests/dat/manifest_hello_nodejs.yaml", // Manifest path
		"helloworld",                              // Package name
		1,                                         // # of Actions
		"helloNodejs",                             // Action name
		"actions/hello.js",                        // Function path
		"nodejs:6",                                // "Runtime
		"")                                        // "Main" function name
}

// Test 2: validate manifest_parser:Unmarshal() method with a sample manifest in Java
// validate that manifest_parser is able to read and parse the manifest data
func TestUnmarshalForHelloJava(t *testing.T) {
	testUnmarshalManifestPackageAndActionBasic(t,
		"../tests/dat/manifest_hello_java_jar.yaml", // Manifest path
		"helloworld",                                // Package name
		1,                                           // # of Actions
		"helloJava",                                 // Action name
		"actions/hello.jar",                         // Function path
		"java",                                      // "Runtime
		"Hello")                                     // "Main" function name
}

// Test 3: validate manifest_parser:Unmarshal() method with a sample manifest in Python
// validate that manifest_parser is able to read and parse the manifest data
func TestUnmarshalForHelloPython(t *testing.T) {
	testUnmarshalManifestPackageAndActionBasic(t,
		"../tests/dat/manifest_hello_python.yaml", // Manifest path
		"helloworld",                              // Package name
		1,                                         // # of Actions
		"helloPython",                             // Action name
		"actions/hello.py",                        // Function path
		"python",                                  // "Runtime
		"")                                        // "Main" function name
}

// Test 4: validate manifest_parser:Unmarshal() method with a sample manifest in Swift
// validate that manifest_parser is able to read and parse the manifest data
func TestUnmarshalForHelloSwift(t *testing.T) {
	testUnmarshalManifestPackageAndActionBasic(t,
		"../tests/dat/manifest_hello_swift.yaml", // Manifest path
		"helloworld",                             // Package name
		1,                                        // # of Actions
		"helloSwift",                             // Action name
		"actions/hello.swift",                    // Function path
		"swift",                                  // "Runtime
		"")                                       // "Main" function name
}

// Test 5: validate manifest_parser:Unmarshal() method for an action with parameters
// validate that manifest_parser is able to read and parse the manifest data, specially
// validate two input parameters and their values
func TestUnmarshalForHelloWithParams(t *testing.T) {

	TEST_ACTION_NAME := "helloWithParams"
	TEST_PARAM_NAME_1 := "name"
	TEST_PARAM_VALUE_1 := "Amy"
	TEST_PARAM_NAME_2 := "place"
	TEST_PARAM_VALUE_2 := "Paris"

	_, pkg, _ := testUnmarshalManifestPackageAndActionBasic(t,
		"../tests/dat/manifest_hello_nodejs_with_params.yaml", // Manifest path
		"helloworld",                   // Package name
		1,                              // # of Actions
		TEST_ACTION_NAME,               // Action name
		"actions/hello-with-params.js", // Function path
		"nodejs:6",                     // "Runtime
		"")                             // "Main" function name

	if pkg != nil {
		if action, ok := pkg.Actions[TEST_ACTION_NAME]; ok {

			// test action parameters
			actualResult := action.Inputs[TEST_PARAM_NAME_1].Value.(string)
			assert.Equal(t, TEST_PARAM_VALUE_1, actualResult,
				fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, TEST_PARAM_NAME_1))

			actualResult = action.Inputs[TEST_PARAM_NAME_2].Value.(string)
			assert.Equal(t, TEST_PARAM_VALUE_2, actualResult,
				fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, TEST_PARAM_NAME_2))

		}
	}
}

// Test 6: validate manifest_parser:Unmarshal() method for an invalid manifest
// manifest_parser should report an error when a package section is missing
func TestUnmarshalForMissingPackages(t *testing.T) {
	m, err := testReadAndUnmarshalManifest(t, "../tests/dat/manifest_invalid_packages_key_missing.yaml")
	assert.NotNil(t, err, fmt.Sprintf(TEST_MSG_MANIFEST_UNMARSHALL_ERROR_EXPECTED, m.Filepath))
}

/*
 Test 7: validate manifest_parser:ParseManifest() method for multiline parameters
 manifest_parser should be able to parse all different multiline combinations of
 inputs section.
*/
func TestParseManifestForMultiLineParams(t *testing.T) {

	_, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_validate_multiline_params.yaml")

	// validate package name should be "validate"
	packageName := "validate"

	// validate this package contains one action
	expectedActionsCount := 1
	actualActionsCount := len(m.Packages[packageName].Actions)
	assert.Equal(t, expectedActionsCount, actualActionsCount, TEST_MSG_ACTION_NUMBER_MISMATCH)

	// here Package.Actions holds a map of map[string]Action
	// where string is the action name so in case you create two actions with
	// same name, will go unnoticed
	// also, the Action struct does not have name field set it to action name
	actionName := "validate_multiline_params"

	if action, ok := m.Packages[packageName].Actions[actionName]; ok {
		// test action function's path
		expectedResult := "actions/dump_params.js"
		actualResult := action.Function
		assert.Equal(t, expectedResult, actualResult, TEST_MSG_ACTION_FUNCTION_PATH_MISMATCH)

		// test action's runtime
		expectedResult = "nodejs:6"
		actualResult = action.Runtime
		assert.Equal(t, expectedResult, actualResult, TEST_MSG_ACTION_FUNCTION_RUNTIME_MISMATCH)

		// test # input params
		expectedResult = strconv.FormatInt(13, 10)
		actualResult = strconv.FormatInt(int64(len(action.Inputs)), 10)
		assert.Equal(t, expectedResult, actualResult, TEST_MSG_PARAMETER_NUMBER_MISMATCH)

		// validate inputs to this action
		for input, param := range action.Inputs {
			switch input {
			case "param_string_value_only":
				expectedResult = "foo"
				actualResult = param.Value.(string)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_int_value_only":
				expectedResult = strconv.FormatInt(123, 10)
				actualResult = strconv.FormatInt(int64(param.Value.(int)), 10)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_float_value_only":
				expectedResult = strconv.FormatFloat(3.14, 'f', -1, 64)
				actualResult = strconv.FormatFloat(param.Value.(float64), 'f', -1, 64)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_string_type_and_value_only":
				expectedResult = "foo"
				actualResult = param.Value.(string)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
				expectedResult = "string"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_string_type_only":
				expectedResult = "string"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_integer_type_only":
				expectedResult = "integer"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_float_type_only":
				expectedResult = "float"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_string_with_default":
				expectedResult = "string"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
				expectedResult = "bar"
				actualResult = param.Default.(string)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_integer_with_default":
				expectedResult = "integer"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
				expectedResult = strconv.FormatInt(-1, 10)
				actualResult = strconv.FormatInt(int64(param.Default.(int)), 10)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_float_with_default":
				expectedResult = "float"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_TYPE_MISMATCH, param))
				expectedResult = strconv.FormatFloat(2.9, 'f', -1, 64)
				actualResult = strconv.FormatFloat(param.Default.(float64), 'f', -1, 64)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			}
		}

		// validate Outputs from this action
		for output, param := range action.Outputs {
			switch output {
			case "payload":
				expectedType := "string"
				actualType := param.Type
				assert.Equal(t, expectedType, actualType, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_TYPE_MISMATCH, param))
				expectedDesc := "parameter dump"
				actualDesc := param.Description
				assert.Equal(t, expectedDesc, actualDesc, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))

			}
		}
	}
}

// Test 8: validate manifest_parser:ParseManifest() method for single line parameters
// manifest_parser should be able to parse input section with different types of values
func TestParseManifestForSingleLineParams(t *testing.T) {

	_, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_validate_singleline_params.yaml")

	// validate package name should be "validate"
	packageName := "validate"

	// validate this package contains one action
	expectedActionsCount := 1
	actualActionsCount := len(m.Packages[packageName].Actions)
	assert.Equal(t, expectedActionsCount, actualActionsCount, TEST_MSG_ACTION_NUMBER_MISMATCH)

	actionName := "validate_singleline_params"
	if action, ok := m.Packages[packageName].Actions[actionName]; ok {
		// test Action function's path
		expectedResult := "actions/dump_params.js"
		actualResult := action.Function
		assert.Equal(t, expectedResult, actualResult, TEST_MSG_ACTION_FUNCTION_PATH_MISMATCH)

		// test Action runtime
		expectedResult = "nodejs:6"
		actualResult = action.Runtime
		assert.Equal(t, expectedResult, actualResult, TEST_MSG_ACTION_FUNCTION_RUNTIME_MISMATCH)

		// test # of inputs
		expectedResult = strconv.FormatInt(22, 10)
		actualResult = strconv.FormatInt(int64(len(action.Inputs)), 10)
		assert.Equal(t, expectedResult, actualResult, TEST_MSG_PARAMETER_NUMBER_MISMATCH)

		// validate Inputs to this action
		for input, param := range action.Inputs {
			switch input {
			case "param_simple_string":
				expectedResult = "foo"
				actualResult = param.Value.(string)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_simple_integer_1":
				expectedResult = strconv.FormatInt(1, 10)
				actualResult = strconv.FormatInt(int64(param.Value.(int)), 10)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_simple_integer_2":
				expectedResult = strconv.FormatInt(0, 10)
				actualResult = strconv.FormatInt(int64(param.Value.(int)), 10)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_simple_integer_3":
				expectedResult = strconv.FormatInt(-1, 10)
				actualResult = strconv.FormatInt(int64(param.Value.(int)), 10)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_simple_integer_4":
				expectedResult = strconv.FormatInt(99999, 10)
				actualResult = strconv.FormatInt(int64(param.Value.(int)), 10)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_simple_integer_5":
				expectedResult = strconv.FormatInt(-99999, 10)
				actualResult = strconv.FormatInt(int64(param.Value.(int)), 10)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_simple_float_1":
				expectedResult = strconv.FormatFloat(1.1, 'f', -1, 64)
				actualResult = strconv.FormatFloat(param.Value.(float64), 'f', -1, 64)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_simple_float_2":
				expectedResult = strconv.FormatFloat(0.0, 'f', -1, 64)
				actualResult = strconv.FormatFloat(param.Value.(float64), 'f', -1, 64)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_simple_float_3":
				expectedResult = strconv.FormatFloat(-1.1, 'f', -1, 64)
				actualResult = strconv.FormatFloat(param.Value.(float64), 'f', -1, 64)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_simple_env_var_1":
				expectedResult = "$GOPATH"
				actualResult = param.Value.(string)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_simple_invalid_env_var":
				expectedResult = "$DollarSignNotInEnv"
				actualResult = param.Value.(string)
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			case "param_simple_implied_empty":
				assert.Nil(t, param.Value, "Expected nil")
			case "param_simple_explicit_empty_1":
				actualResult = param.Value.(string)
				assert.Empty(t, actualResult)
			case "param_simple_explicit_empty_2":
				actualResult = param.Value.(string)
				assert.Empty(t, actualResult)
			}
		}

		// validate Outputs from this action
		for output, param := range action.Outputs {
			switch output {
			case "payload":
				expectedResult = "string"
				actualResult = param.Type
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_TYPE_MISMATCH, param))

				expectedResult = "parameter dump"
				actualResult = param.Description
				assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, param))
			}
		}
	}
}

// Test 9: validate manifest_parser.ComposeActions() method for implicit runtimes
// when a runtime of an action is not provided, manifest_parser determines the runtime
// based on the file extension of an action file
func TestComposeActionsForImplicitRuntimes(t *testing.T) {

	p, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_data_compose_runtimes_implicit.yaml")

	actions, err := p.ComposeActionsFromAllPackages(m, m.Filepath, whisk.KeyValue{})
	var expectedResult string
	if err == nil {
		for i := 0; i < len(actions); i++ {
			if actions[i].Action.Name == "helloNodejs" {
				expectedResult = utils.DefaultRunTimes[utils.FileExtensionRuntimeKindMap["js"]]
			} else if actions[i].Action.Name == "helloJava" {
				expectedResult = utils.DefaultRunTimes[utils.FileExtensionRuntimeKindMap["jar"]]
			} else if actions[i].Action.Name == "helloPython" {
				expectedResult = utils.DefaultRunTimes[utils.FileExtensionRuntimeKindMap["py"]]
			} else if actions[i].Action.Name == "helloSwift" {
				expectedResult = utils.DefaultRunTimes[utils.FileExtensionRuntimeKindMap["swift"]]
			}
			actualResult := actions[i].Action.Exec.Kind
			assert.Equal(t, expectedResult, actualResult, TEST_MSG_ACTION_FUNCTION_RUNTIME_MISMATCH)
		}
	}

}

// Test 10(1): validate manifest_parser.ComposeActions() method for invalid runtimes
// when the action has a source file written in unsupported runtimes, manifest_parser should
// report an error for that action
// TODO() rewrite
func TestComposeActionsForInvalidRuntime_1(t *testing.T) {
	data := `packages:
    helloworld:
        actions:
            helloInvalidRuntime:
                function: ../tests/src/integration/common/wskdeploy.go`
	p, m, tmpfile := testUnmarshalTemporaryFile([]byte(data), "manifest_parser_validate_runtime_")
	_, err := p.ComposeActionsFromAllPackages(m, tmpfile, whisk.KeyValue{})
	assert.NotNil(t, err, fmt.Sprintf(TEST_MSG_ACTION_FUNCTION_RUNTIME_ERROR_EXPECTED, tmpfile))
}

// Test 10(2): validate manifest_parser.ComposeActions() method for invalid runtimes
// when a runtime of an action is missing for zip action, manifest_parser should
// report an error for that action
func TestComposeActionsForInvalidRuntime_2(t *testing.T) {
	data := `packages:
    helloworld:
        actions:
            helloInvalidRuntime:
                function: ../tests/src/integration/runtimetests/src/helloworld/`
	p, m, tmpfile := testUnmarshalTemporaryFile([]byte(data), "manifest_parser_validate_runtime_")
	_, err := p.ComposeActionsFromAllPackages(m, tmpfile, whisk.KeyValue{})
	assert.NotNil(t, err, fmt.Sprintf(TEST_MSG_ACTION_FUNCTION_RUNTIME_ERROR_EXPECTED, tmpfile))
}

// Test 10(3): validate manifest_parser.ComposeActions() method for invalid runtimes
// when a runtime of an action is missing for zip action, manifest_parser should
// report an error for that action
func TestComposeActionsForInvalidRuntime_3(t *testing.T) {
	data := `packages:
    helloworld:
        actions:
            helloInvalidRuntime:
                function: ../tests/src/integration/runtimetests/src/helloworld/helloworld.zip`
	p, m, tmpfile := testUnmarshalTemporaryFile([]byte(data), "manifest_parser_validate_runtime_")
	_, err := p.ComposeActionsFromAllPackages(m, tmpfile, whisk.KeyValue{})
	assert.NotNil(t, err, fmt.Sprintf(TEST_MSG_ACTION_FUNCTION_RUNTIME_ERROR_EXPECTED, tmpfile))
}

// Test 10(3): validate manifest_parser.ComposeActions() method for valid runtimes with zip action
// when a runtime of a zip action is set to one of the supported runtimes, manifest_parser should
// return a valid actionRecord with specified runtime
func TestComposeActionsForValidRuntime_ZipAction(t *testing.T) {
	data := `packages:
    helloworld:
        actions:
            hello:
                function: ../tests/src/integration/runtimetests/src/helloworld/helloworld.zip
                runtime: nodejs:6`
	p, m, tmpfile := testUnmarshalTemporaryFile([]byte(data), "manifest_parser_validate_runtime_")
	actions, _ := p.ComposeActionsFromAllPackages(m, tmpfile, whisk.KeyValue{})
	for _, action := range actions {
		if action.Action.Name == "hello" {
			assert.Equal(t, action.Action.Exec.Kind, "nodejs:6", fmt.Sprintf(TEST_MSG_ACTION_FUNCTION_RUNTIME_MISMATCH, action))
		}

	}
}

// Test 11: validate manifest_parser.ComposeActions() method for single line parameters
// manifest_parser should be able to parse input section with different types of values
func TestComposeActionsForSingleLineParams(t *testing.T) {

	p, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_validate_singleline_params.yaml")

	// Call the method we are testing
	actions, err := p.ComposeActionsFromAllPackages(m, m.Filepath, whisk.KeyValue{})

	if err == nil {
		// test # actions
		assert.Equal(t, 1, len(actions), TEST_MSG_ACTION_NUMBER_MISMATCH)

		action := actions[0]

		/*
		 * Simple 'string' value tests
		 */

		// param_simple_string should value "foo"
		paramName := "param_simple_string"
		expectedResult := "foo"
		actualResult := action.Action.Parameters.GetValue(paramName).(string)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		/*
		 * Simple 'integer' value tests
		 */

		// param_simple_integer_1 should have value 1
		paramName = "param_simple_integer_1"
		expectedResult = strconv.FormatInt(1, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue(paramName).(int)), 10)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_integer_2 should have value 0
		paramName = "param_simple_integer_2"
		expectedResult = strconv.FormatInt(0, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue(paramName).(int)), 10)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_integer_3 should have value -1
		paramName = "param_simple_integer_3"
		expectedResult = strconv.FormatInt(-1, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue(paramName).(int)), 10)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_integer_4 should have value 99999
		paramName = "param_simple_integer_4"
		expectedResult = strconv.FormatInt(99999, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue(paramName).(int)), 10)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_integer_5 should have value -99999
		paramName = "param_simple_integer_5"
		expectedResult = strconv.FormatInt(-99999, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue(paramName).(int)), 10)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		/*
		 * Simple 'float' value tests
		 */

		// param_simple_float_1 should have value 1.1
		paramName = "param_simple_float_1"
		expectedResult = strconv.FormatFloat(1.1, 'f', -1, 64)
		actualResult = strconv.FormatFloat(action.Action.Parameters.GetValue(paramName).(float64), 'f', -1, 64)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_float_2 should have value 0.0
		paramName = "param_simple_float_2"
		expectedResult = strconv.FormatFloat(0.0, 'f', -1, 64)
		actualResult = strconv.FormatFloat(action.Action.Parameters.GetValue(paramName).(float64), 'f', -1, 64)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_float_3 should have value -1.1
		paramName = "param_simple_float_3"
		expectedResult = strconv.FormatFloat(-1.1, 'f', -1, 64)
		actualResult = strconv.FormatFloat(action.Action.Parameters.GetValue(paramName).(float64), 'f', -1, 64)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		/*
		 * Environment Variable / dollar ($) notation tests
		 */

		// param_simple_env_var_1 should have value of env. variable $GOPATH
		paramName = "param_simple_env_var_1"
		expectedResult = os.Getenv("GOPATH")
		actualResult = action.Action.Parameters.GetValue(paramName).(string)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_env_var_2 should have value of env. variable $GOPATH
		paramName = "param_simple_env_var_2"
		expectedResult = os.Getenv("GOPATH")
		actualResult = action.Action.Parameters.GetValue(paramName).(string)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_env_var_3 should have value of env. variable "${}"
		paramName = "param_simple_env_var_3"
		expectedResult = "${}"
		actualResult = action.Action.Parameters.GetValue(paramName).(string)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_invalid_env_var should have value of ""
		paramName = "param_simple_invalid_env_var"
		expectedResult = ""
		actualResult = action.Action.Parameters.GetValue(paramName).(string)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		/*
		 * Environment Variable concatenation tests
		 */

		// param_simple_env_var_concat_1 should have value of env. variable "$GOPTH/test" empty string
		paramName = "param_simple_env_var_concat_1"
		expectedResult = os.Getenv("GOPATH") + "/test"
		actualResult = action.Action.Parameters.GetValue(paramName).(string)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_env_var_concat_2 should have value of env. variable "" empty string
		// as the "/test" is treated as part of the environment var. and not concatenated.
		paramName = "param_simple_env_var_concat_2"
		expectedResult = ""
		actualResult = action.Action.Parameters.GetValue(paramName).(string)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_env_var_concat_3 should have value of env. variable "" empty string
		paramName = "param_simple_env_var_concat_3"
		expectedResult = "ddd.ccc." + os.Getenv("GOPATH")
		actualResult = action.Action.Parameters.GetValue(paramName).(string)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		/*
		 * Empty string tests
		 */

		// param_simple_implied_empty should be ""
		paramName = "param_simple_implied_empty"
		actualResult = action.Action.Parameters.GetValue(paramName).(string)
		assert.Empty(t, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_explicit_empty_1 should be ""
		paramName = "param_simple_explicit_empty_1"
		actualResult = action.Action.Parameters.GetValue(paramName).(string)
		assert.Empty(t, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_explicit_empty_2 should be ""
		paramName = "param_simple_explicit_empty_2"
		actualResult = action.Action.Parameters.GetValue(paramName).(string)
		assert.Empty(t, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		/*
		 * Test values that contain "Type names" (e.g., "string", "integer", "float, etc.)
		 */

		// param_simple_type_string should be "" when value set to "string"
		paramName = "param_simple_type_string"
		expectedResult = ""
		actualResult = action.Action.Parameters.GetValue(paramName).(string)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_type_integer should be 0.0 when value set to "integer"
		paramName = "param_simple_type_integer"
		expectedResult = strconv.FormatInt(0, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue(paramName).(int)), 10)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_simple_type_float should be 0 when value set to "float"
		paramName = "param_simple_type_float"
		expectedResult = strconv.FormatFloat(0.0, 'f', -1, 64)
		actualResult = strconv.FormatFloat(action.Action.Parameters.GetValue(paramName).(float64), 'f', -1, 64)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

	}
}

// Test 12: validate manifest_parser.ComposeActions() method for multi line parameters
// manifest_parser should be able to parse input section with different types of values
func TestComposeActionsForMultiLineParams(t *testing.T) {
	os.Setenv("USERNAME", "MY_USERNAME")
	os.Setenv("PASSWORD", "MY_PASSWORD")
	defer os.Unsetenv("USERNAME")
	defer os.Unsetenv("PASSWORD")

	p, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_validate_multiline_params.yaml")

	// call the method we are testing
	actions, err := p.ComposeActionsFromAllPackages(m, m.Filepath, whisk.KeyValue{})

	if err == nil {
		// test # actions
		assert.Equal(t, 1, len(actions), TEST_MSG_ACTION_NUMBER_MISMATCH)

		action := actions[0]

		// param_string_value_only should be "foo"
		paramName := "param_string_value_only"
		expectedResult := "foo"
		actualResult := action.Action.Parameters.GetValue(paramName).(string)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_int_value_only should be 123
		paramName = "param_int_value_only"
		expectedResult = strconv.FormatInt(123, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue(paramName).(int)), 10)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_float_value_only should be 3.14
		paramName = "param_float_value_only"
		expectedResult = strconv.FormatFloat(3.14, 'f', -1, 64)
		actualResult = strconv.FormatFloat(action.Action.Parameters.GetValue(paramName).(float64), 'f', -1, 64)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_string_type_and_value_only should be foo
		paramName = "param_string_type_and_value_only"
		expectedResult = "foo"
		actualResult = action.Action.Parameters.GetValue(paramName).(string)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_string_type_only should be ""
		paramName = "param_string_type_only"
		actualResult = action.Action.Parameters.GetValue(paramName).(string)
		assert.Empty(t, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_integer_type_only should be 0
		paramName = "param_integer_type_only"
		expectedResult = strconv.FormatInt(0, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue(paramName).(int)), 10)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_float_type_only should be 0
		paramName = "param_float_type_only"
		expectedResult = strconv.FormatFloat(0.0, 'f', -1, 64)
		actualResult = strconv.FormatFloat(action.Action.Parameters.GetValue(paramName).(float64), 'f', -1, 64)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_string_with_default should be "bar"
		paramName = "param_string_with_default"
		expectedResult = "bar"
		actualResult = action.Action.Parameters.GetValue(paramName).(string)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_integer_with_default should be -1
		paramName = "param_integer_with_default"
		expectedResult = strconv.FormatInt(-1, 10)
		actualResult = strconv.FormatInt(int64(action.Action.Parameters.GetValue(paramName).(int)), 10)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_float_with_default should be 2.9
		paramName = "param_float_with_default"
		expectedResult = strconv.FormatFloat(2.9, 'f', -1, 64)
		actualResult = strconv.FormatFloat(action.Action.Parameters.GetValue(paramName).(float64), 'f', -1, 64)
		assert.Equal(t, expectedResult, actualResult, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_json_type_and_value_only_1 should be { "name": "Sam", "place": "Shire" }
		paramName = "param_json_type_and_value_only_1"
		expectedResult1 := map[string]interface{}{"name": "Sam", "place": "Shire"}
		actualResult1 := action.Action.Parameters.GetValue(paramName)
		assert.Equal(t, expectedResult1, actualResult1, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_json_type_and_value_only_2 should be { "name": "MY_USERNAME", "password": "MY_PASSWORD" }
		paramName = "param_json_type_and_value_only_2"
		expectedResult2 := map[string]interface{}{"name": "MY_USERNAME", "password": "MY_PASSWORD"}
		actualResult2 := action.Action.Parameters.GetValue(paramName)
		assert.Equal(t, expectedResult2, actualResult2, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

		// param_json_type_and_value_only_3 should be { "name": "${USERNAME}", "password": "${PASSWORD}" }
		paramName = "param_json_type_and_value_only_3"
		expectedResult3 := map[string]interface{}{"name": "${USERNAME}", "password": "${PASSWORD}"}
		actualResult3 := action.Action.Parameters.GetValue(paramName)
		assert.Equal(t, expectedResult3, actualResult3, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))
	}
}

// Test 13: validate manifest_parser.ComposeActions() method
func TestComposeActionsForFunction(t *testing.T) {

	p, m, err := testLoadParseManifest(t, "../tests/dat/manifest_data_compose_actions_for_function.yaml")
	assert.Nil(t, err, "Failed to parse manifest file.")

	actions, err := p.ComposeActionsFromAllPackages(m, m.Filepath, whisk.KeyValue{})
	assert.Nil(t, err, "Failed to compose actions from manifest file.")

	var expectedResult, actualResult string
	for i := 0; i < len(actions); i++ {
		if actions[i].Action.Name == "hello1" {
			expectedResult, _ = filepath.Abs("../tests/src/integration/helloworld/actions/hello.js")
			actualResult, _ = filepath.Abs(actions[i].Filepath)
			assert.Equal(t, expectedResult, actualResult, "Expected "+expectedResult+" but got "+actualResult)
		} else if actions[i].Action.Name == "hello2" {
			spew.Dump(actions[i].Action.Exec.Code)
			assert.NotNil(t, actions[i].Action.Exec.Code, "Expected source code from an action file but found it empty")
		}
	}
}

// validate manifest_parser.ComposeActions() method
func TestComposeActionsForFunctionAndCode(t *testing.T) {
	p, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_data_compose_actions_for_function_and_code.yaml")
	_, err := p.ComposeActionsFromAllPackages(m, m.Filepath, whisk.KeyValue{})
	assert.NotNil(t, err, "Compose actions should have exited with error when code and function both exist.")
}

// validate manifest_parser.ComposeActions() method
func TestComposeActionsForCodeWithMissingRuntime(t *testing.T) {
	p, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_data_compose_actions_for_missing_runtime_With_code.yaml")
	_, err := p.ComposeActionsFromAllPackages(m, m.Filepath, whisk.KeyValue{})
	assert.NotNil(t, err, "Compose actions should have exited with error when code is specified but runtime is missing.")
}

// Test 14: validate manifest_parser.ComposeActions() method
func TestComposeActionsForLimits(t *testing.T) {

	p, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_data_compose_actions_for_limits.yaml")

	actions, err := p.ComposeActionsFromAllPackages(m, m.Filepath, whisk.KeyValue{})

	if err == nil {
		for i := 0; i < len(actions); i++ {
			if actions[i].Action.Name == "hello1" {
				assert.Nil(t, actions[i].Action.Limits, "Expected limit section to be empty but got %s", actions[i].Action.Limits)
			} else if actions[i].Action.Name == "hello2" {
				assert.NotNil(t, actions[i].Action.Limits, "Expected limit section to be not empty but found it empty")
				assert.Equal(t, 180, *actions[i].Action.Limits.Timeout, "Failed to get Timeout")
				assert.Equal(t, 128, *actions[i].Action.Limits.Memory, "Failed to get Memory")
				assert.Equal(t, 1, *actions[i].Action.Limits.Logsize, "Failed to get Logsize")
			}
		}
	}

}

// Test 15: validate manifest_parser.ComposeActions() method
func TestComposeActionsForWebActions(t *testing.T) {

	p, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_data_compose_actions_for_web.yaml")

	actions, err := p.ComposeActionsFromAllPackages(m, m.Filepath, whisk.KeyValue{})
	if err == nil {
		for i := 0; i < len(actions); i++ {
			if actions[i].Action.Name == "hello1" {
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
			} else if actions[i].Action.Name == "hello2" {
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
			} else if actions[i].Action.Name == "hello3" {
				for _, a := range actions[i].Action.Annotations {
					switch a.Key {
					case "web-export":
						assert.Equal(t, true, a.Value, "Expected true for web-export but got "+strconv.FormatBool(a.Value.(bool)))
					case "raw-http":
						assert.Equal(t, true, a.Value, "Expected false for raw-http but got "+strconv.FormatBool(a.Value.(bool)))
					case "final":
						assert.Equal(t, true, a.Value, "Expected true for final but got "+strconv.FormatBool(a.Value.(bool)))
					}
				}
			} else if actions[i].Action.Name == "hello4" {
				for _, a := range actions[i].Action.Annotations {
					switch a.Key {
					case "web-export":
						assert.Equal(t, false, a.Value, "Expected true for web-export but got "+strconv.FormatBool(a.Value.(bool)))
					case "raw-http":
						assert.Equal(t, false, a.Value, "Expected false for raw-http but got "+strconv.FormatBool(a.Value.(bool)))
					case "final":
						assert.Equal(t, false, a.Value, "Expected true for final but got "+strconv.FormatBool(a.Value.(bool)))
					}
				}
			} else if actions[i].Action.Name == "hello5" {
				for _, a := range actions[i].Action.Annotations {
					switch a.Key {
					case "web-export":
						assert.Equal(t, false, a.Value, "Expected true for web-export but got "+strconv.FormatBool(a.Value.(bool)))
					case "raw-http":
						assert.Equal(t, false, a.Value, "Expected false for raw-http but got "+strconv.FormatBool(a.Value.(bool)))
					case "final":
						assert.Equal(t, false, a.Value, "Expected true for final but got "+strconv.FormatBool(a.Value.(bool)))
					}
				}
			}
		}
	}
}

// Test 15-1: validate manifest_parser.ComposeActions() method
func TestComposeActionsForInvalidWebActions(t *testing.T) {

	p, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_data_compose_actions_for_invalid_web.yaml")
	_, err := p.ComposeActionsFromAllPackages(m, m.Filepath, whisk.KeyValue{})
	assert.NotNil(t, err, "Expected error for invalid web-export.")
}

// Test 16: validate manifest_parser.ResolveParameter() method
func TestResolveParameterForMultiLineParams(t *testing.T) {
	paramName := "name"
	v := "foo"
	y := reflect.TypeOf(v).Name() // y := string
	d := "default_name"

	// type string - value only param
	param1 := Parameter{Value: v, multiline: true}
	r1, _ := ResolveParameter(paramName, &param1, "")
	assert.Equal(t, v, r1, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))
	assert.IsType(t, v, r1, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_TYPE_MISMATCH, paramName))

	// type string - type and value only param
	param2 := Parameter{Type: y, Value: v, multiline: true}
	r2, _ := ResolveParameter(paramName, &param2, "")
	assert.Equal(t, v, r2, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))
	assert.IsType(t, v, r2, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_TYPE_MISMATCH, paramName))

	// type string - type, no value, but default value param
	param3 := Parameter{Type: y, Default: d, multiline: true}
	r3, _ := ResolveParameter(paramName, &param3, "")
	assert.Equal(t, d, r3, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))
	assert.IsType(t, d, r3, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_TYPE_MISMATCH, paramName))

	// type string - type and value only param
	// type is "string" and value is of type "int"
	// ResolveParameter matches specified type with the type of the specified value
	// it fails if both types don't match
	// ResolveParameter determines type from the specified value
	// in this case, ResolveParameter returns value of type int
	v1 := 11
	param4 := Parameter{Type: y, Value: v1, multiline: true}
	r4, _ := ResolveParameter(paramName, &param4, "")
	assert.Equal(t, v1, r4, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))
	assert.IsType(t, v1, r4, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_TYPE_MISMATCH, paramName))

	// type invalid - type only param
	param5 := Parameter{Type: "invalid", multiline: true}
	_, err := ResolveParameter(paramName, &param5, "")
	assert.NotNil(t, err, "Expected error saying Invalid type for parameter")
	switch errorType := err.(type) {
	default:
		assert.Fail(t, "Wrong error type received: We are expecting ParserErr.")
	case *wskderrors.YAMLParserError:
		assert.Equal(t, "Parameter [name] has an invalid Type. [invalid]", errorType.Message)
	}

	// type none - param without type, without value, and without default value
	param6 := Parameter{multiline: true}
	paramName = "none"
	r6, _ := ResolveParameter(paramName, &param6, "")
	assert.Empty(t, r6, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, paramName))

}

// Test 17: validate JSON parameters
func TestParseManifestForJSONParams(t *testing.T) {

	_, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_validate_json_params.yaml")

	// validate package name should be "validate"
	packageName := "validate_json"
	actionName := "validate_json_params"

	// validate this package contains one action
	actualActionsCount := len(m.Packages[packageName].Actions)
	assert.Equal(t, 1, actualActionsCount, TEST_MSG_ACTION_NUMBER_MISMATCH)

	if action, ok := m.Packages[packageName].Actions[actionName]; ok {
		// test Action function's path
		expectedResult := "actions/dump_params.js"
		actualResult := action.Function
		assert.Equal(t, expectedResult, actualResult, TEST_MSG_ACTION_FUNCTION_PATH_MISMATCH)

		// validate runtime of an action to be "nodejs:6"
		expectedResult = "nodejs:6"
		actualResult = action.Runtime
		assert.Equal(t, expectedResult, actualResult, TEST_MSG_ACTION_FUNCTION_RUNTIME_MISMATCH)

		// validate the number of inputs to this action
		expectedResult = strconv.FormatInt(8, 10)
		actualResult = strconv.FormatInt(int64(len(action.Inputs)), 10)
		assert.Equal(t, expectedResult, actualResult, TEST_MSG_PARAMETER_NUMBER_MISMATCH)

		// validate inputs to this action
		for input, param := range action.Inputs {
			// Trace to help debug complex values:
			// utils.PrintTypeInfo(input, param.Value)
			switch input {
			case "member1":
				actualResult1 := param.Value.(string)
				expectedResult1 := "{ \"name\": \"Sam\", \"place\": \"Shire\" }"
				assert.Equal(t, expectedResult1, actualResult1, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, input))
			case "member2":
				actualResult2 := param.Value.(map[interface{}]interface{})
				expectedResult2 := map[interface{}]interface{}{"name": "Sam", "place": "Shire"}
				assert.Equal(t, expectedResult2, actualResult2, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, input))
			case "member3":
				actualResult3 := param.Value.(map[interface{}]interface{})
				expectedResult3 := map[interface{}]interface{}{"name": "Elrond", "place": "Rivendell"}
				assert.Equal(t, expectedResult3, actualResult3, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, input))
			case "member4":
				actualResult4 := param.Value.(map[interface{}]interface{})
				expectedResult4 := map[interface{}]interface{}{"name": "Gimli", "place": "Gondor", "age": 139, "children": map[interface{}]interface{}{"<none>": "<none>"}}
				assert.Equal(t, expectedResult4, actualResult4, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, input))
			case "member5":
				actualResult5 := param.Value.(map[interface{}]interface{})
				expectedResult5 := map[interface{}]interface{}{"name": "Gloin", "place": "Gondor", "age": 235, "children": map[interface{}]interface{}{"Gimli": "Son"}}
				assert.Equal(t, expectedResult5, actualResult5, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, input))
			case "member6":
				actualResult6 := param.Value.(map[interface{}]interface{})
				expectedResult6 := map[interface{}]interface{}{"name": "Frodo", "place": "Undying Lands", "items": []interface{}{"Sting", "Mithril mail"}}
				assert.Equal(t, expectedResult6, actualResult6, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, input))
			case "member7":
				actualResult7 := param.Value.(map[interface{}]interface{})
				expectedResult7 := map[interface{}]interface{}{"name": "${USERNAME}", "password": "${PASSWORD}"}
				assert.Equal(t, expectedResult7, actualResult7, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, input))
			case "member8":
				actualResult8 := param.Value.(map[interface{}]interface{})
				expectedResult8 := map[interface{}]interface{}{"name": "$${USERNAME}", "password": "$${PASSWORD}"}
				assert.Equal(t, expectedResult8, actualResult8, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_VALUE_MISMATCH, input))
			}
		}

		// validate Outputs from this action
		for output, param := range action.Outputs {
			switch output {
			case "fellowship":
				expectedType := "json"
				actualType := param.Type
				assert.Equal(t, expectedType, actualType, fmt.Sprintf(TEST_MSG_ACTION_PARAMETER_TYPE_MISMATCH, output))
			}
		}
	}
}

func TestComposePackage(t *testing.T) {

	p, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_data_compose_packages.yaml")

	pkg, err := p.ComposeAllPackages(m, m.Filepath, whisk.KeyValue{})
	if err == nil {
		n := "helloworld"
		assert.NotNil(t, pkg[n], "Failed to get the whole package")
		assert.Equal(t, n, pkg[n].Name, "Failed to get package name")
		assert.Equal(t, "default", pkg[n].Namespace, "Failed to get package namespace")
	} else {
		assert.Fail(t, "Failed to compose package")
	}
}

func TestComposeSequences(t *testing.T) {

	p, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_data_compose_sequences.yaml")

	// Note: set first param (namespace) to empty string
	seqList, err := p.ComposeSequencesFromAllPackages("", m, whisk.KeyValue{})
	if err != nil {
		assert.Fail(t, "Failed to compose sequences")
	}
	assert.Equal(t, 2, len(seqList), "Failed to get sequences")
	for _, seq := range seqList {
		wsk_action := seq.Action
		switch wsk_action.Name {
		case "sequence1":
			assert.Equal(t, "sequence", wsk_action.Exec.Kind, "Failed to set sequence exec kind")
			assert.Equal(t, 2, len(wsk_action.Exec.Components), "Failed to set sequence exec components")
			assert.Equal(t, "/helloworld/action1", wsk_action.Exec.Components[0], "Failed to set sequence 1st exec components")
			assert.Equal(t, "/helloworld/action2", wsk_action.Exec.Components[1], "Failed to set sequence 2nd exec components")
		case "sequence2":
			assert.Equal(t, "sequence", wsk_action.Exec.Kind, "Failed to set sequence exec kind")
			assert.Equal(t, 3, len(wsk_action.Exec.Components), "Failed to set sequence exec components")
			assert.Equal(t, "/helloworld/action3", wsk_action.Exec.Components[0], "Failed to set sequence 1st exec components")
			assert.Equal(t, "/helloworld/action4", wsk_action.Exec.Components[1], "Failed to set sequence 2nd exec components")
			assert.Equal(t, "/helloworld/action5", wsk_action.Exec.Components[2], "Failed to set sequence 3rd exec components")
		}
	}
}

func TestComposeTriggers(t *testing.T) {
	// set env variables needed for the trigger feed
	os.Setenv("KAFKA_INSTANCE", "kafka-broker")
	os.Setenv("SRC_TOPIC", "topic")

	p, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_data_compose_triggers.yaml")

	triggerList, err := p.ComposeTriggersFromAllPackages(m, m.Filepath, whisk.KeyValue{})

	if err != nil {
		assert.Fail(t, "Failed to compose trigger")
	}

	assert.Equal(t, 3, len(triggerList), "Failed to get trigger list")
	for _, trigger := range triggerList {
		switch trigger.Name {
		case "trigger1":
			assert.Equal(t, 2, len(trigger.Parameters), "Failed to set trigger parameters")
		case "trigger2":
			assert.Equal(t, "feed", trigger.Annotations[0].Key, "Failed to set trigger annotation")
			assert.Equal(t, "myfeed", trigger.Annotations[0].Value, "Failed to set trigger annotation")
			assert.Equal(t, 2, len(trigger.Parameters), "Failed to set trigger parameters")
		case "message-trigger":
			assert.Equal(t, 2, len(trigger.Parameters), "Failed to set trigger parameters")
			assert.Equal(t, "feed", trigger.Annotations[0].Key, "Failed to set trigger annotation")
			assert.Equal(t, "Bluemix_kafka-broker_Credentials-1/messageHubFeed", trigger.Annotations[0].Value, "Failed to set trigger annotation")
		}
	}
}

func TestComposeRules(t *testing.T) {

	p, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_data_compose_rules.yaml")

	ruleList, err := p.ComposeRulesFromAllPackages(m, whisk.KeyValue{})
	if err != nil {
		assert.Fail(t, "Failed to compose rules")
	}
	assert.Equal(t, 2, len(ruleList), "Failed to get rules")
	for _, rule := range ruleList {
		switch rule.Name {
		case "rule1":
			assert.Equal(t, "locationUpdate", rule.Trigger, "Failed to set rule trigger")
			assert.Equal(t, "helloworld/greeting", rule.Action, "Failed to set rule action")
		case "rule2":
			assert.Equal(t, "trigger1", rule.Trigger, "Failed to set rule trigger")
			assert.Equal(t, "helloworld/action1", rule.Action, "Failed to set rule action")
		}
	}
}

// TODO(752) We SHOULD automatically add "web-export" to each Action referenced in the "apis" section
// as this is implied.  The user should not have to do this manually
func TestComposeApiRecords(t *testing.T) {

	p, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_data_compose_api_records.yaml")

	// create a fake configuration
	config := whisk.Config{
		Namespace:        "test",
		AuthToken:        "user:pass",
		Host:             "host",
		ApigwAccessToken: "token",
	}

	apiList, err := p.ComposeApiRecordsFromAllPackages(&config, m)
	if err != nil {
		assert.Fail(t, "Failed to compose api records: "+err.Error())
	}
	assert.Equal(t, 6, len(apiList), "Failed to get api records")
	for _, apiRecord := range apiList {
		apiDoc := apiRecord.ApiDoc
		action := apiDoc.Action
		switch action.Name {
		case "apiTest/putBooks":
			assert.Equal(t, "book-club", apiDoc.ApiName, "Failed to set api name")
			assert.Equal(t, "/club", apiDoc.GatewayBasePath, "Failed to set api base path")
			assert.Equal(t, "/books", apiDoc.GatewayRelPath, "Failed to set api rel path")
			assert.Equal(t, "put", action.BackendMethod, "Failed to set api backend method")
		case "apiTest/deleteBooks":
			assert.Equal(t, "book-club", apiDoc.ApiName, "Failed to set api name")
			assert.Equal(t, "/club", apiDoc.GatewayBasePath, "Failed to set api base path")
			assert.Equal(t, "/books", apiDoc.GatewayRelPath, "Failed to set api rel path")
			assert.Equal(t, "delete", action.BackendMethod, "Failed to set api backend method")
		case "apiTest/listMembers":
			assert.Equal(t, "book-club", apiDoc.ApiName, "Failed to set api name")
			assert.Equal(t, "/club", apiDoc.GatewayBasePath, "Failed to set api base path")
			assert.Equal(t, "/members", apiDoc.GatewayRelPath, "Failed to set api rel path")
			assert.Equal(t, "get", action.BackendMethod, "Failed to set api backend method")
		case "apiTest/getBooks2":
			assert.Equal(t, "book-club2", apiDoc.ApiName, "Failed to set api name")
			assert.Equal(t, "/club2", apiDoc.GatewayBasePath, "Failed to set api base path")
			assert.Equal(t, "/books2", apiDoc.GatewayRelPath, "Failed to set api rel path")
			assert.Equal(t, "get", action.BackendMethod, "Failed to set api backend method")
		case "apiTest/postBooks2":
			assert.Equal(t, "book-club2", apiDoc.ApiName, "Failed to set api name")
			assert.Equal(t, "/club2", apiDoc.GatewayBasePath, "Failed to set api base path")
			assert.Equal(t, "/books2", apiDoc.GatewayRelPath, "Failed to set api rel path")
			assert.Equal(t, "post", action.BackendMethod, "Failed to set api backend method")
		case "apiTest/listMembers2":
			assert.Equal(t, "book-club2", apiDoc.ApiName, "Failed to set api name")
			assert.Equal(t, "/club2", apiDoc.GatewayBasePath, "Failed to set api base path")
			assert.Equal(t, "/members2", apiDoc.GatewayRelPath, "Failed to set api rel path")
			assert.Equal(t, "get", action.BackendMethod, "Failed to set api backend method")
		default:
			assert.Fail(t, "Failed to get api action name")
		}
	}
}

func TestComposeDependencies(t *testing.T) {

	p, m, _ := testLoadParseManifest(t, "../tests/dat/manifest_data_compose_dependencies.yaml")

	depdList, err := p.ComposeDependenciesFromAllPackages(m, "/project_folder", m.Filepath)

	if err != nil {
		assert.Fail(t, "Failed to compose rules")
	}
	assert.Equal(t, 2, len(depdList), "Failed to get rules")
	for depdy_name, depdy := range depdList {
		assert.Equal(t, "helloworld", depdy.Packagename, "Failed to set dependecy isbinding")
		assert.Equal(t, "/project_folder/Packages", depdy.ProjectPath, "Failed to set dependecy isbinding")
		d := strings.Split(depdy_name, ":")
		assert.NotEqual(t, d[1], "", "Failed to get dependency name")
		switch d[1] {
		case "myhelloworld":
			assert.Equal(t, "https://github.com/user/repo/folder", depdy.Location, "Failed to set dependecy location")
			assert.Equal(t, false, depdy.IsBinding, "Failed to set dependecy isbinding")
			assert.Equal(t, "https://github.com/user/repo", depdy.BaseRepo, "Failed to set dependecy base repo url")
			assert.Equal(t, "/folder", depdy.SubFolder, "Failed to set dependecy sub folder")
		case "myCloudant":
			assert.Equal(t, "/whisk.system/cloudant", depdy.Location, "Failed to set rule trigger")
			assert.Equal(t, true, depdy.IsBinding, "Failed to set dependecy isbinding")
			assert.Equal(t, 1, len(depdy.Parameters), "Failed to set dependecy parameter")
			assert.Equal(t, 1, len(depdy.Annotations), "Failed to set dependecy annotation")
			assert.Equal(t, "myAnnotation", depdy.Annotations[0].Key, "Failed to set dependecy parameter key")
			assert.Equal(t, "Here it is", depdy.Annotations[0].Value, "Failed to set dependecy parameter value")
			assert.Equal(t, "dbname", depdy.Parameters[0].Key, "Failed to set dependecy annotation key")
			assert.Equal(t, "myGreatDB", depdy.Parameters[0].Value, "Failed to set dependecy annotation value")
		default:
			assert.Fail(t, "Failed to get dependency name")
		}
	}
}

func TestBadYAMLInvalidPackageKeyInManifest(t *testing.T) {
	// read and parse manifest.yaml file located under ../tests folder
	p := NewYAMLParser()
	_, err := p.ParseManifest("../tests/dat/manifest_bad_yaml_invalid_package_key.yaml")

	assert.NotNil(t, err)
	// NOTE: go-yaml/yaml gets the line # wrong; testing only for the invalid key message
	assert.Contains(t, err.Error(), "field invalidKey not found in struct parsers.Package")
}

func TestBadYAMLInvalidKeyMappingValueInManifest(t *testing.T) {
	// read and parse manifest.yaml file located under ../tests folder
	p := NewYAMLParser()
	_, err := p.ParseManifest("../tests/dat/manifest_bad_yaml_invalid_key_mapping_value.yaml")

	assert.NotNil(t, err)
	// go-yaml/yaml prints the wrong line number for mapping values. It should be 5.
	assert.Contains(t, err.Error(), "mapping values are not allowed in this context")
}

func TestBadYAMLMissingRootKeyInManifest(t *testing.T) {
	// read and parse manifest.yaml file located under ../tests folder
	p := NewYAMLParser()
	_, err := p.ParseManifest("../tests/dat/manifest_bad_yaml_missing_root_key.yaml")

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "field actions not found in struct parsers.YAML")
}

func TestBadYAMLInvalidCommentInManifest(t *testing.T) {
	// read and parse manifest.yaml file located under ../tests folder
	p := NewYAMLParser()
	_, err := p.ParseManifest("../tests/dat/manifest_bad_yaml_invalid_comment.yaml")

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "could not find expected ':'")
}

// validate manifest_parser:Unmarshal() method for package in manifest YAML
// validate that manifest_parser is able to read and parse the manifest data
func TestUnmarshalForPackages(t *testing.T) {

	//manifestFile := "../tests/dat/manifest_data_unmarshal_packages.yaml"
	m, err := testReadAndUnmarshalManifest(t, "../tests/dat/manifest_data_unmarshal_packages.yaml")

	// Unmarshal reads/parses manifest data and sets the values of YAML
	// And returns an error if parsing a manifest data fails
	if err == nil {
		expectedResult := string(2)
		actualResult := string(len(m.Packages))
		assert.Equal(t, expectedResult, actualResult, "Expected 2 packages but got "+actualResult)
		// we have two packages
		// package name should be "helloNodejs" and "helloPython"
		for k, v := range m.Packages {
			switch k {
			case "package1":
				assert.Equal(t, "package1", k, "Expected package name package1 but got "+k)
				expectedResult = string(1)
				actualResult = string(len(v.Actions))
				assert.Equal(t, expectedResult, actualResult, "Expected 1 but got "+actualResult)
				// get the action payload from the map of actions which is stored in
				// YAML.Package.Actions with the type of map[string]Action
				actionName := "helloNodejs"
				if action, ok := v.Actions[actionName]; ok {
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
			case "package2":
				assert.Equal(t, "package2", k, "Expected package name package2 but got "+k)
				expectedResult = string(1)
				actualResult = string(len(v.Actions))
				assert.Equal(t, expectedResult, actualResult, "Expected 1 but got "+actualResult)
				// get the action payload from the map of actions which is stored in
				// YAML.Package.Actions with the type of map[string]Action
				actionName := "helloPython"
				if action, ok := v.Actions[actionName]; ok {
					// location/function of an action should be "actions/hello.js"
					expectedResult = "actions/hello.py"
					actualResult = action.Function
					assert.Equal(t, expectedResult, actualResult, "Expected action function "+expectedResult+" but got "+actualResult)
					// runtime of an action should be "python"
					expectedResult = "python"
					actualResult = action.Runtime
					assert.Equal(t, expectedResult, actualResult, "Expected action runtime "+expectedResult+" but got "+actualResult)
				} else {
					t.Error("Action named " + actionName + " does not exist.")
				}
			}
		}
	}
}

func TestParseYAML_trigger(t *testing.T) {
	data, err := ioutil.ReadFile("../tests/dat/manifest_validate_triggerfeed.yaml")
	if err != nil {
		panic(err)
	}

	var manifest YAML
	err = NewYAMLParser().Unmarshal(data, &manifest)
	if err != nil {
		panic(err)
	}

	packageName := "manifest3"

	assert.Equal(t, 2, len(manifest.Packages[packageName].Triggers), "Get trigger list failed.")
	for trigger_name := range manifest.Packages[packageName].Triggers {
		var trigger = manifest.Packages[packageName].Triggers[trigger_name]
		switch trigger_name {
		case "trigger1":
		case "trigger2":
			assert.Equal(t, "myfeed", trigger.Feed, "Get trigger feed name failed.")
		default:
			t.Error("Get trigger name failed")
		}
	}
}

func TestParseYAML_rule(t *testing.T) {
	data, err := ioutil.ReadFile("../tests/dat/manifest_validate_rule.yaml")
	if err != nil {
		panic(err)
	}

	var manifest YAML
	err = NewYAMLParser().Unmarshal(data, &manifest)
	if err != nil {
		panic(err)
	}

	packageName := "manifest4"

	assert.Equal(t, 1, len(manifest.Packages[packageName].Rules), "Get trigger list failed.")
	for rule_name := range manifest.Packages[packageName].Rules {
		var rule = manifest.Packages[packageName].Rules[rule_name]
		switch rule_name {
		case "rule1":
			assert.Equal(t, "trigger1", rule.Trigger, "Get trigger name failed.")
			assert.Equal(t, "hellpworld", rule.Action, "Get action name failed.")
			assert.Equal(t, "true", rule.Rule, "Get rule expression failed.")
		default:
			t.Error("Get rule name failed")
		}
	}
}

func TestParseYAML_feed(t *testing.T) {
	data, err := ioutil.ReadFile("../tests/dat/manifest_validate_feed.yaml")
	if err != nil {
		panic(err)
	}

	var manifest YAML
	err = NewYAMLParser().Unmarshal(data, &manifest)
	if err != nil {
		panic(err)
	}

	packageName := "manifest5"

	assert.Equal(t, 1, len(manifest.Packages[packageName].Feeds), "Get feed list failed.")
	for feed_name := range manifest.Packages[packageName].Feeds {
		var feed = manifest.Packages[packageName].Feeds[feed_name]
		switch feed_name {
		case "feed1":
			assert.Equal(t, "https://my.company.com/services/eventHub", feed.Location, "Get feed location failed.")
			assert.Equal(t, "my_credential", feed.Credential, "Get feed credential failed.")
			assert.Equal(t, 2, len(feed.Operations), "Get operations number failed.")
			for operation_name := range feed.Operations {
				switch operation_name {
				case "operation1":
				case "operation2":
				default:
					t.Error("Get feed operation name failed")
				}
			}
		default:
			t.Error("Get feed name failed")
		}
	}
}

func TestParseYAML_param(t *testing.T) {
	data, err := ioutil.ReadFile("../tests/dat/manifest_validate_params.yaml")
	if err != nil {
		panic(err)
	}

	var manifest YAML
	err = NewYAMLParser().Unmarshal(data, &manifest)
	if err != nil {
		panic(err)
	}

	packageName := "validateParams"

	assert.Equal(t, 1, len(manifest.Packages[packageName].Actions), "Get action list failed.")
	for action_name := range manifest.Packages[packageName].Actions {
		var action = manifest.Packages[packageName].Actions[action_name]
		switch action_name {
		case "action1":
			for param_name := range action.Inputs {
				var param = action.Inputs[param_name]
				switch param_name {
				case "inline1":
					assert.Equal(t, "{ \"key\": true }", param.Value, "Get param value failed.")
				case "inline2":
					assert.Equal(t, "Just a string", param.Value, "Get param value failed.")
				case "inline3":
					assert.Equal(t, nil, param.Value, "Get param value failed.")
				case "inline4":
					assert.Equal(t, true, param.Value, "Get param value failed.")
				case "inline5":
					assert.Equal(t, 42, param.Value, "Get param value failed.")
				case "inline6":
					assert.Equal(t, -531, param.Value, "Get param value failed.")
				case "inline7":
					assert.Equal(t, 432.432E-43, param.Value, "Get param value failed.")
				case "inline8":
					assert.Equal(t, "[ true, null, \"boo\", { \"key\": 0 }]", param.Value, "Get param value failed.")
				case "inline9":
					assert.Equal(t, false, param.Value, "Get param value failed.")
				case "inline0":
					assert.Equal(t, 456.423, param.Value, "Get param value failed.")
				case "inlin10":
					assert.Equal(t, nil, param.Value, "Get param value failed.")
				case "inlin11":
					assert.Equal(t, true, param.Value, "Get param value failed.")
				case "expand1":
					assert.Equal(t, nil, param.Value, "Get param value failed.")
				case "expand2":
					assert.Equal(t, true, param.Value, "Get param value failed.")
				case "expand3":
					assert.Equal(t, false, param.Value, "Get param value failed.")
				case "expand4":
					assert.Equal(t, 15646, param.Value, "Get param value failed.")
				case "expand5":
					assert.Equal(t, "{ \"key\": true }", param.Value, "Get param value failed.")
				case "expand6":
					assert.Equal(t, "[ true, null, \"boo\", { \"key\": 0 }]", param.Value, "Get param value failed.")
				case "expand7":
					assert.Equal(t, nil, param.Value, "Get param value failed.")
				default:
					t.Error("Get param name failed")
				}
			}
		default:
			t.Error("Get action name failed")
		}
	}
}

func TestPackageName_Env_Var(t *testing.T) {
	testPackage := "test_package"
	os.Setenv("package_name", testPackage)
	testPackageSec := "test_package_second"
	os.Setenv("package_name_second", testPackageSec)
	mm := NewYAMLParser()
	manifestfile := "../tests/dat/manifest_validate_package_grammar_env_var.yaml"
	manifest, _ := mm.ParseManifest(manifestfile)
	assert.Equal(t, 4, len(manifest.Packages), "Get package list failed.")
	expectedPackages := [4]string{testPackage, testPackageSec, testPackage + "suffix", testPackage + "-" + testPackageSec}
	for _, pkg_name := range expectedPackages {
		var pkg = manifest.Packages[pkg_name]
		assert.Equal(t, "1.0", pkg.Version, "Get the wrong package version.")
		assert.Equal(t, "Apache-2.0", pkg.License, "Get the wrong license.")
	}
}

func TestRuleName_Env_Var(t *testing.T) {
	// read and parse manifest file with env var for rule name, and rule trigger and action
	testRule := "test_rule"
	os.Setenv("rule_name", testRule)
	testTrigger := "test_trigger"
	os.Setenv("trigger_name", testTrigger)
	testAction := "test_actions"
	os.Setenv("action_name", testAction)
	mm := NewYAMLParser()
	manifestfile := "../tests/dat/manifest_data_rule_env_var.yaml"
	manifest, _ := mm.ParseManifest(manifestfile)
	rules, err := mm.ComposeRulesFromAllPackages(manifest, whisk.KeyValue{})
	if err != nil {
		assert.Fail(t, "Failed to compose rules")
	}
	packageName := "manifest1"

	assert.Equal(t, 1, len(manifest.Packages[packageName].Rules), "Get rule list failed.")
	for _, rule := range rules {
		wskprint.PrintlnOpenWhiskVerbose(false, fmt.Sprintf("ruleName: %v", rule))
		switch rule.Name {
		case testRule:
			assert.Equal(t, "test_trigger", rule.Trigger, "Get trigger name failed.")
			assert.Equal(t, packageName+"/"+testAction, rule.Action, "Get action name failed.")
		//assert.Equal(t, "true", rule.Rule, "Get rule expression failed.")
		default:
			t.Error("Get rule name failed")
		}
	}
}

func TestComposeActionForAnnotations(t *testing.T) {
	manifestFile := "../tests/dat/manifest_validate_action_annotations.yaml"
	mm := NewYAMLParser()
	manifest, _ := mm.ParseManifest(manifestFile)
	pkg_name := "packageActionAnnotations"
	pkg := manifest.Packages[pkg_name]
	assert.NotNil(t, pkg, "Could not find package with name "+pkg_name)
	action_name := "helloworld"
	action := pkg.Actions[action_name]
	assert.NotNil(t, action, "Could not find action with name "+action_name)
	actual_annotations := action.Annotations
	expected_annotations := map[string]interface{}{
		"action_annotation_1": "this is annotation 1",
		"action_annotation_2": "this is annotation 2",
		"action_annotation_3": "this is annotation 3",
		"action_annotation_4": "this is annotation 4",
	}
	assert.Equal(t, len(actual_annotations), len(expected_annotations), "Could not find expected number of annotations specified in manifest file")
	eq := reflect.DeepEqual(actual_annotations, expected_annotations)
	assert.True(t, eq, "Expected list of annotations does not match with actual list, expected annotations: %v actual annotations: %v", expected_annotations, actual_annotations)

	pkg_name = "packageActionAnnotationsWithWebAction"
	pkg = manifest.Packages[pkg_name]
	assert.NotNil(t, pkg, "Could not find package with name "+pkg_name)
	action = pkg.Actions[action_name]
	assert.NotNil(t, action, "Could not find action with name "+action_name)
	actual_annotations = action.Annotations
	expected_annotations["web-export"] = true
	assert.Equal(t, len(actual_annotations), len(expected_annotations), "Could not find expected number of annotations specified in manifest file")
	eq = reflect.DeepEqual(actual_annotations, expected_annotations)
	assert.True(t, eq, "Expected list of annotations does not match with actual list, expected annotations: %v actual annotations: %v", expected_annotations, actual_annotations)
}
