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
	"runtime"
	"strings"
	"path/filepath"
)

const (
	UNKNOWN = "Unknown"
	UNKNOWN_VALUE = "Unknown value"
	LINE = "line"
	PARAMETER = "Parameter"
	TYPE = "Type"
	EXPECTED = "expected"
	ACTUAL = "actual"
	YAML_FILE = "YAML file"

	ERROR_COMMAND_FAILED = "ERROR_COMMAND_FAILED"
	ERROR_FILE_READ_ERROR = "ERROR_FILE_READ_ERROR"
	ERROR_MANIFEST_FILE_NOT_FOUND = "ERROR_MANIFEST_FILE_NOT_FOUND"
	ERROR_YAML_FILE_READ_ERROR = "ERROR_YAML_FILE_READ_ERROR"
	ERROR_YAML_FILE_FORMAT_ERROR = "ERROR_YAML_FILE_FORMAT_ERROR"
	ERROR_WHISK_CLIENT_ERROR = "ERROR_WHISK_CLIENT_ERROR"
	ERROR_WHISK_CLIENT_INVALID_CONFIG = "ERROR_WHISK_CLIENT_INVALID_CONFIG"
	ERROR_YAML_PARAMETER_TYPE_MISMATCH = "ERROR_YAML_PARAMETER_TYPE_MISMATCH"
	ERROR_YAML_INVALID_PARAMETER_TYPE = "ERROR_YAML_INVALID_PARAMETER_TYPE"
)

/*
 * BaseError
 */
type BaseErr struct {
	ErrorType string
	FileName  string
	LineNum   int
	Message   string
}

func (e *BaseErr) Error() string {
	if e.ErrorType == "" {
		return fmt.Sprintf("%s [%d]: %s\n", e.FileName, e.LineNum, e.Message)
	}
	return fmt.Sprintf("%s [%d]: [%s]: %s\n", e.FileName, e.LineNum, e.ErrorType, e.Message)
}

func (e *BaseErr) SetFileName(fileName string) {
	e.FileName = filepath.Base(fileName)
}

func (e *BaseErr) SetLineNum(lineNum int) {
	e.LineNum = lineNum
}

func (e *BaseErr) SetMessage(message string) {
	e.Message = message
}

func (e *BaseErr) SetErrorType(errorType string) {
	e.ErrorType = errorType
}

// func Caller(skip int) (pc uintptr, file string, line int, ok bool)
func (e *BaseErr) SetCallerByStackFrameSkip(skip int) {
	_, fname, lineNum, _ := runtime.Caller(skip)
	e.SetFileName(fname)
	e.SetLineNum(lineNum)
}

/*
 * FileReadError
 */
type FileReadError struct {
	BaseErr
	FileName string
	FilePath string
}

func (e *FileReadError) SetFilePath(fpath string) {
	e.FilePath = filepath.Base(fpath)
}

func (e *FileReadError) Error() string {
	return e.BaseErr.Error()
}

func NewFileReadError(fpath string, errMessage string) *FileReadError {
	var err = &FileReadError{
	}
	err.SetErrorType(ERROR_FILE_READ_ERROR)
	err.SetCallerByStackFrameSkip(2)
	err.SetMessage(errMessage)
	err.SetFilePath(fpath)
	return err
}

/*
 * YAML Base Error
 */
type YAMLBaseError struct {
	BaseErr
	YAMLFileName string
	YAMLFilePath string
}

func (e *YAMLBaseError) SetYAMLFileName(fname string) {
	e.YAMLFileName = filepath.Base(fname)
}

func (e *YAMLBaseError) SetYAMLFilePath(fpath string) {
	e.YAMLFilePath = filepath.Base(fpath)
}

func (e *YAMLBaseError) Error() string {
	return e.BaseErr.Error()
}

/*
 * CommandError
 */
type CommandError struct {
	BaseErr
	Command string
}

func NewCommandError(cmd string, errorMessage string) *CommandError {
	var err = &CommandError{
		Command: cmd,
	}
	err.SetCallerByStackFrameSkip(2)
	err.SetErrorType(ERROR_COMMAND_FAILED)
	err.SetMessage(cmd + ": " + errorMessage)
	return err
}

func (e *CommandError) Error() string {
	return e.BaseErr.Error()
}

/*
 * ManifestFileNotFoundError
 */
type ErrorManifestFileNotFound struct {
	YAMLBaseError
}

func NewErrorManifestFileNotFound(fpath string, errMessage string) *ErrorManifestFileNotFound {
	var err = &ErrorManifestFileNotFound{
	}
	err.SetErrorType(ERROR_MANIFEST_FILE_NOT_FOUND)
	err.SetCallerByStackFrameSkip(2)
	err.SetMessage(errMessage)
	err.SetYAMLFilePath(fpath)
	return err
}

func (e *ErrorManifestFileNotFound) Error() string {
	return e.BaseErr.Error()
}

/*
 * YAMLFileReadError
 */
type YAMLFileReadError struct {
	YAMLBaseError
}

func NewYAMLFileReadError(fname string, errorMessage string) *YAMLFileReadError {
	var err = &YAMLFileReadError{
	}
	err.SetErrorType(ERROR_YAML_FILE_READ_ERROR)
	err.SetCallerByStackFrameSkip(2)
	err.SetMessage(errorMessage)
	err.SetYAMLFileName(fname)
	return err
}

func (e *YAMLFileReadError) Error() string {
	return e.BaseErr.Error()
}

/*
 * YAMLFileFormatError
 */
type YAMLFileFormatError struct {
	YAMLFileReadError
}

func NewYAMLFileFormatError(fname string, errorMessage string) *YAMLFileFormatError {
	var err = &YAMLFileFormatError{
	}
	err.SetErrorType(ERROR_YAML_FILE_FORMAT_ERROR)
	err.SetCallerByStackFrameSkip(2)
	err.SetMessage(errorMessage)
	return err
}

/*
 * WhiskClientError
 */
type WhiskClientError struct {
	BaseErr
	ErrorCode int
}

func NewWhiskClientError(errorMessage string, code int) *WhiskClientError {
	var err = &WhiskClientError{
		ErrorCode: code,
	}
	err.SetErrorType(ERROR_WHISK_CLIENT_ERROR)
	err.SetCallerByStackFrameSkip(2)
	str := fmt.Sprintf("Error Code: %d: %s", code, errorMessage)
	err.SetMessage(str)
	return err
}

func (e *WhiskClientError) Error() string {
	return e.BaseErr.Error()
}

/*
 * WhiskClientInvalidConfigError
 */
type WhiskClientInvalidConfigError struct {
	BaseErr
}

func NewWhiskClientInvalidConfigError(errorMessage string) *WhiskClientInvalidConfigError {
	var err = &WhiskClientInvalidConfigError{
	}
	err.SetErrorType(ERROR_WHISK_CLIENT_INVALID_CONFIG)
	err.SetCallerByStackFrameSkip(2)
	err.SetMessage(errorMessage)
	return err
}

/*
 * YAMLParserErr
 */
type YAMLParserErr struct {
	YAMLBaseError
	YamlFile string
	lines    []string
	msgs     []string
}

func NewYAMLParserErr(yamlFile string, lines []string, msgs []string) *YAMLParserErr {
	var err = &YAMLParserErr{
		YamlFile: yamlFile,
		lines: lines,
		msgs: msgs,
	}
	err.SetCallerByStackFrameSkip(2)
	return err
}


func (e *YAMLParserErr) Error() string {
	result := make([]string, len(e.msgs))

	for index, msg := range e.msgs {
		var s string
		if e.lines == nil || e.lines[index] == UNKNOWN {
			s = fmt.Sprintf("====> %s", msg)
		} else {
			s = fmt.Sprintf("====> %s [%v]: %s", LINE, e.lines[index], msg)
		}
		result[index] = s
	}
	return fmt.Sprintf("\n==> %s [%d]: %s: %s: \n%s", e.FileName, e.LineNum, YAML_FILE, e.YamlFile, strings.Join(result, "\n"))
}

/*
 * ParameterTypeMismatchError
 */
type ParameterTypeMismatchError struct {
	YAMLParserErr
	Parameter    string
	ExpectedType string
	ActualType   string
}

func NewParameterTypeMismatchError(yamlFile string, param string, expectedType string, actualType string) *ParameterTypeMismatchError {
	var err = &ParameterTypeMismatchError{
		ExpectedType: expectedType,
		ActualType: actualType,
	}
	err.SetYAMLFileName(yamlFile)
	err.SetErrorType(ERROR_YAML_PARAMETER_TYPE_MISMATCH)
	err.SetCallerByStackFrameSkip(2)
	str := fmt.Sprintf("%s [%s]: %s %s: [%s], %s: [%s]", PARAMETER, param, TYPE, EXPECTED, expectedType, ACTUAL, actualType)
	err.SetMessage(str)
	return err
}

func (e *ParameterTypeMismatchError) Error() string {
	return e.BaseErr.Error()
}

/*
 * InvalidParameterType
 */
type InvalidParameterTypeError struct {
	YAMLParserErr
	Parameter    string
	ActualType   string
}

func NewInvalidParameterTypeError(yamlFile string, param string, actualType string) *ParameterTypeMismatchError {
	var err = &ParameterTypeMismatchError{
		ActualType: actualType,
	}
	err.SetYAMLFileName(yamlFile)
	err.SetErrorType(ERROR_YAML_INVALID_PARAMETER_TYPE)
	err.SetCallerByStackFrameSkip(2)
	str := fmt.Sprintf("%s [%s]: %s [%s]", PARAMETER, param, TYPE, actualType)
	err.SetMessage(str)
	return err
}

func (e *InvalidParameterTypeError) Error() string {
	return e.BaseErr.Error()
}
