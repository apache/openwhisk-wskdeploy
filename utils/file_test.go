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
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const IS_FILE_PATH_FAILURE_FOR_FILE = "Failed to classify the file path [%s] as file."
const IS_FILE_PATH_FAILURE_FOR_DIR = "Failed to classify the dir path [%s] as dir."
const IS_FILE_PATH_FAILURE_FOR_EMPTY_PATH = "Failed to classify empty path [%s] as dir."
const IS_FILE_FAILURE_FOR_FILE = "Failed to detect if the path [%s] is a valid file."
const IS_FILE_FAILURE_FOR_DIR = "Failed to detect if the path [%s] is a dir."
const TEST_ERROR_IS_FILE_FAILURE = "Failed to run isFile on [%s] Error: %s"

func TestIsFilePath(t *testing.T) {
	paths := []string{
		"",
		"/home/arnie/amelia.jpg",
		"/mnt/photos/",
		"rabbit.jpg",
		"/usr/local//go",
		"common/actions/*.js/",
		"common/*/actions/*/code/*.js",
		"common/actions/",
	}
	for _, p := range paths {
		switch p {
		case "":
			assert.False(t, isFilePath(p), fmt.Sprintf(IS_FILE_PATH_FAILURE_FOR_EMPTY_PATH, p))
		case "/home/arnie/amelia.jpg":
			assert.True(t, isFilePath(p), fmt.Sprintf(IS_FILE_PATH_FAILURE_FOR_FILE, p))
		case "/mnt/photos/":
			assert.False(t, isFilePath(p), fmt.Sprintf(IS_FILE_PATH_FAILURE_FOR_DIR, p))
		case "rabbit.jpg":
			assert.True(t, isFilePath(p), fmt.Sprintf(IS_FILE_PATH_FAILURE_FOR_FILE, p))
		case "/usr/local//go":
			assert.True(t, isFilePath(p), fmt.Sprintf(IS_FILE_PATH_FAILURE_FOR_FILE, p))
		case "common/actions/*.js/":
			assert.False(t, isFilePath(p), fmt.Sprintf(IS_FILE_PATH_FAILURE_FOR_DIR, p))
		case "common/*/actions/*/code/*.js":
			assert.True(t, isFilePath(p), fmt.Sprintf(IS_FILE_PATH_FAILURE_FOR_FILE, p))
		case "common/actions/":
			assert.False(t, isFilePath(p), fmt.Sprintf(IS_FILE_PATH_FAILURE_FOR_DIR, p))
		}
	}
}

func TestIsFile(t *testing.T) {
	paths := []string{
		"../tests/dat/",
		"../tests/dat/manifest_hello_swift.yaml",
		"../tests/data",
		"../tests/data/manifest_hello_swift.yaml",
	}
	for _, p := range paths {
		f, err := isFile(p)
		if err != nil {
			assert.Fail(t, fmt.Sprintf(TEST_ERROR_IS_FILE_FAILURE, p, err.Error()))
		}
		switch p {
		case "tests/dat/":
			assert.False(t, f, fmt.Sprintf(IS_FILE_FAILURE_FOR_FILE, p))
		case "tests/dat/manifest_hello_swift.yaml":
			assert.True(t, f, fmt.Sprintf(IS_FILE_FAILURE_FOR_FILE, p))
		case "../tests/data":
			assert.False(t, f, fmt.Sprintf(IS_FILE_FAILURE_FOR_DIR, p))
		case "../tests/data/manifest_hello_swift.yaml":
			assert.False(t, f, fmt.Sprintf(IS_FILE_FAILURE_FOR_DIR, p))
		}
	}
}
