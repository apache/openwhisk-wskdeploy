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
	"strings"
	"fmt"
	"runtime"
	"path/filepath"
)

/*
 * TestCustomErrorOutputFormat
 */
func TestCustomErrorOutputFormat(t *testing.T) {

	_, fn, _, _ := runtime.Caller(0)
	packageName := filepath.Base(fn)
	const TEST_DEFAULT_ERROR_MESSAGE = "Some bad error"
	const TEST_COMMAND string = "test"
	const TEST_ERROR_CODE = 400  // Bad request
	const TEST_EXISTANT_MANIFEST_FILE = "tests/dat/manifest_validate_multiline_params.yaml"
	const TEST_NONEXISTANT_MANIFEST_FILE = "tests/dat/missing_manifest.yaml"
	const TEST_INVALID_YAML_MANIFEST_FILE = "tests/dat/manifest_bad_yaml_invalid_comment.yaml"
	const TEST_PARAM_NAME = "Age"
	const TEST_PARAM_TYPE_INT = "integer"
	const TEST_PARAM_TYPE_FLOAT = "float"
	const TEST_PARAM_TYPE_FOO = "foo"

	/*
	 * CommandError
	 */
	err1 := NewCommandError(TEST_COMMAND, TEST_DEFAULT_ERROR_MESSAGE)
	actualResult :=  strings.TrimSpace(err1.Error())
	expectedResult := fmt.Sprintf("%s [%d]: [%s]: %s: [%s]: %s",
		packageName,
		err1.LineNum,
		ERROR_COMMAND_FAILED,
		STR_COMMAND,
		TEST_COMMAND,
		TEST_DEFAULT_ERROR_MESSAGE )
	assert.Equal(t, expectedResult, actualResult)

	/*
	 * WhiskClientError
	 */
	err2 := NewWhiskClientError(TEST_DEFAULT_ERROR_MESSAGE, TEST_ERROR_CODE)
	actualResult =  strings.TrimSpace(err2.Error())
	expectedResult = fmt.Sprintf("%s [%d]: [%s]: %s: %d: %s",
		packageName,
		err2.LineNum,
		ERROR_WHISK_CLIENT_ERROR,
		STR_ERROR_CODE,
		TEST_ERROR_CODE,
		TEST_DEFAULT_ERROR_MESSAGE )
	assert.Equal(t, expectedResult, actualResult)

	/*
	 * WhiskClientInvalidConfigError
	 */
	err3 := NewWhiskClientInvalidConfigError(TEST_DEFAULT_ERROR_MESSAGE)
	actualResult =  strings.TrimSpace(err3.Error())
	expectedResult = fmt.Sprintf("%s [%d]: [%s]: %s",
		packageName,
		err3.LineNum,
		ERROR_WHISK_CLIENT_INVALID_CONFIG,
		TEST_DEFAULT_ERROR_MESSAGE )
	assert.Equal(t, expectedResult, actualResult)

	/*
 	 * FileReadError
 	 */
	err4 := NewFileReadError(TEST_NONEXISTANT_MANIFEST_FILE, TEST_DEFAULT_ERROR_MESSAGE)
	actualResult =  strings.TrimSpace(err4.Error())
	expectedResult = fmt.Sprintf("%s [%d]: [%s]: " + STR_FILE + ": [%s]: %s",
		packageName,
		err4.LineNum,
		ERROR_FILE_READ_ERROR,
		filepath.Base(TEST_NONEXISTANT_MANIFEST_FILE),
		TEST_DEFAULT_ERROR_MESSAGE )
	assert.Equal(t, expectedResult, actualResult)

	/*
 	 * ManifestFileNotFoundError
 	 */
	err5 := NewErrorManifestFileNotFound(TEST_NONEXISTANT_MANIFEST_FILE, TEST_DEFAULT_ERROR_MESSAGE)
	actualResult =  strings.TrimSpace(err5.Error())
	expectedResult = fmt.Sprintf("%s [%d]: [%s]: %s: [%s]: %s",
		packageName,
		err5.LineNum,
		ERROR_MANIFEST_FILE_NOT_FOUND,
		STR_FILE,
		filepath.Base(TEST_NONEXISTANT_MANIFEST_FILE),
		TEST_DEFAULT_ERROR_MESSAGE )
	assert.Equal(t, expectedResult, actualResult)

	/*
         * YAMLFileFormatError
         */
	err6 := NewYAMLFileFormatError(TEST_INVALID_YAML_MANIFEST_FILE, TEST_DEFAULT_ERROR_MESSAGE)
	actualResult =  strings.TrimSpace(err6.Error())
	expectedResult = fmt.Sprintf("%s [%d]: [%s]: %s: [%s]: %s",
		packageName,
		err6.LineNum,
		ERROR_YAML_FILE_FORMAT_ERROR,
		STR_FILE,
		filepath.Base(TEST_INVALID_YAML_MANIFEST_FILE),
		TEST_DEFAULT_ERROR_MESSAGE )
	assert.Equal(t, expectedResult, actualResult)

	/*
	 * ParameterTypeMismatchError
	 */
	err8 := NewParameterTypeMismatchError(
		TEST_EXISTANT_MANIFEST_FILE,
		TEST_PARAM_NAME,
		TEST_PARAM_TYPE_INT,
		TEST_PARAM_TYPE_FLOAT)
	actualResult =  strings.TrimSpace(err8.Error())
	msg8 := fmt.Sprintf("%s [%s]: %s %s: [%s], %s: [%s]",
		STR_PARAMETER, TEST_PARAM_NAME,
		STR_TYPE,
		STR_EXPECTED, TEST_PARAM_TYPE_INT,
		STR_ACTUAL, TEST_PARAM_TYPE_FLOAT)
	expectedResult = fmt.Sprintf("%s [%d]: [%s]: %s: [%s]: %s",
		packageName,
		err8.LineNum,
		ERROR_YAML_PARAMETER_TYPE_MISMATCH,
		STR_FILE,
		filepath.Base(TEST_EXISTANT_MANIFEST_FILE),
		msg8 )
	assert.Equal(t, expectedResult, actualResult)

	/*
	 * InvalidParameterType
	 */
	err9 := NewInvalidParameterTypeError(TEST_EXISTANT_MANIFEST_FILE, TEST_PARAM_NAME, TEST_PARAM_TYPE_FOO)
	actualResult =  strings.TrimSpace(err9.Error())
	msg9 := fmt.Sprintf("%s [%s]: %s [%s]",
		STR_PARAMETER, TEST_PARAM_NAME,
		STR_TYPE, TEST_PARAM_TYPE_FOO)
	expectedResult = fmt.Sprintf("%s [%d]: [%s]: " + STR_FILE + ": [%s]: %s",
		packageName,
		err9.LineNum,
		ERROR_YAML_INVALID_PARAMETER_TYPE,
		filepath.Base(TEST_EXISTANT_MANIFEST_FILE),
		msg9 )
	assert.Equal(t, expectedResult, actualResult)

	/*
	 * YAMLParserErr
	 */

	// TODO add a unit test once we re-factor error related modules into a new package
	// to avoid cyclic deplendency errors in GoLang.
	//manifestFile := "../tests/dat/manifest_bad_yaml_invalid_comment.yaml"
	//// read and parse manifest.yaml file
	//m, err := parsers.NewYAMLParser().ParseManifest(manifestFile)
	//fmt.Println(err)

	// TODO - use actual YAML files to generate actual errors for comparison with expected error output
	//var TEST_LINES    = []string{"40", STR_UNKNOWN, "123"}
	//var TEST_MESSAGES = []string{"did not find expected key", "did not find expected ',' or ']'", "found duplicate %YAML directive"}
	//
	//err10 := NewYAMLParserErr(TEST_EXISTANT_MANIFEST_FILE, TEST_LINES, TEST_MESSAGES)
	//actualResult =  strings.TrimSpace(err10.Error())
	//
	//msgs := "\n==> Line [40]: did not find expected key" +
	//	"\n==> Line [Unknown]: did not find expected ',' or ']'" +
	//	"\n==> Line [123]: found duplicate %YAML directive"
	//
	//expectedResult = fmt.Sprintf("%s [%d]: [%s]: " + STR_FILE + ": [%s]: %s",
	//	packageName,
	//	err10.LineNum,
	//	ERROR_YAML_PARSER_ERROR,
	//	filepath.Base(TEST_EXISTANT_MANIFEST_FILE),
	//	msgs)
	//assert.Equal(t, expectedResult, actualResult)
}
