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
	FILE = "File"
	LINE = "Line"
	PARAMETER = "Parameter"
	TYPE = "Type"
	EXPECTED = "expected"
	ACTUAL = "actual"

	ERROR_COMMAND_FAILED = "ERROR_COMMAND_FAILED"
	ERROR_WHISK_CLIENT_ERROR = "ERROR_WHISK_CLIENT_ERROR"
	ERROR_WHISK_CLIENT_INVALID_CONFIG = "ERROR_WHISK_CLIENT_INVALID_CONFIG"
	ERROR_FILE_READ_ERROR = "ERROR_FILE_READ_ERROR"
	ERROR_MANIFEST_FILE_NOT_FOUND = "ERROR_MANIFEST_FILE_NOT_FOUND"
	ERROR_YAML_FILE_FORMAT_ERROR = "ERROR_YAML_FILE_FORMAT_ERROR"
	ERROR_YAML_PARSER_ERROR = "ERROR_YAML_PARSER_ERROR"
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
	err.SetErrorType(ERROR_COMMAND_FAILED)
	err.SetCallerByStackFrameSkip(2)
	err.SetMessage(cmd + ": " + errorMessage)
	return err
}

func (e *CommandError) Error() string {
	return e.BaseErr.Error()
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
 * FileError
 */
type FileError struct {
	BaseErr
	ErrorFileName string
	ErrorFilePath string
}

func (e *FileError) SetErrorFilePath(fpath string) {
	e.ErrorFilePath = fpath
	e.ErrorFileName = filepath.Base(fpath)
}

func (e *FileError) SetErrorFileName(fname string) {
	e.ErrorFilePath = fname
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%s [%d]: [%s]: " + FILE + ": [%s]: %s\n",
		e.FileName,
		e.LineNum,
		e.ErrorType,
		e.ErrorFileName,
		e.Message)
}

/*
 * FileReadError
 */
func NewFileReadError(fpath string, errMessage string) *FileError {
	var err = &FileError{
	}
	err.SetErrorType(ERROR_FILE_READ_ERROR)
	err.SetCallerByStackFrameSkip(2)
	err.SetErrorFilePath(fpath)
	err.SetMessage(errMessage)
	return err
}


/*
 * ManifestFileNotFoundError
 */
type ErrorManifestFileNotFound struct {
	FileError
}

func NewErrorManifestFileNotFound(fpath string, errMessage string) *ErrorManifestFileNotFound {
	var err = &ErrorManifestFileNotFound{
	}
	err.SetErrorType(ERROR_MANIFEST_FILE_NOT_FOUND)
	err.SetCallerByStackFrameSkip(2)
	err.SetErrorFilePath(fpath)
	err.SetMessage(errMessage)
	return err
}

/*
 * YAMLFileFormatError
 */
type YAMLFileFormatError struct {
	FileError
}

func NewYAMLFileFormatError(fpath string, errorMessage string) *YAMLFileFormatError {
	var err = &YAMLFileFormatError{
	}
	err.SetErrorType(ERROR_YAML_FILE_FORMAT_ERROR)
	err.SetCallerByStackFrameSkip(2)
	err.SetErrorFilePath(fpath)
	err.SetMessage(errorMessage)
	return err
}

/*
 * ParameterTypeMismatchError
 */
type ParameterTypeMismatchError struct {
	FileError
	Parameter    string
	ExpectedType string
	ActualType   string
}

func NewParameterTypeMismatchError(fpath string, param string, expectedType string, actualType string) *ParameterTypeMismatchError {
	var err = &ParameterTypeMismatchError{
		ExpectedType: expectedType,
		ActualType: actualType,
	}

	err.SetErrorType(ERROR_YAML_PARAMETER_TYPE_MISMATCH)
	err.SetCallerByStackFrameSkip(2)
	err.SetErrorFilePath(fpath)
	str := fmt.Sprintf("%s [%s]: %s %s: [%s], %s: [%s]",
		PARAMETER, param,
		TYPE,
		EXPECTED, expectedType,
		ACTUAL, actualType)
	err.SetMessage(str)
	return err
}

/*
 * InvalidParameterType
 */
type InvalidParameterTypeError struct {
	FileError
	Parameter    string
	ActualType   string
}

func NewInvalidParameterTypeError(fpath string, param string, actualType string) *ParameterTypeMismatchError {
	var err = &ParameterTypeMismatchError{
		ActualType: actualType,
	}
	err.SetErrorFilePath(fpath)
	err.SetErrorType(ERROR_YAML_INVALID_PARAMETER_TYPE)
	err.SetCallerByStackFrameSkip(2)
	str := fmt.Sprintf("%s [%s]: %s [%s]",
		PARAMETER, param,
		TYPE, actualType)
	err.SetMessage(str)
	return err
}

/*
 * YAMLParserErr
 */
type YAMLParserError struct {
	FileError
	lines    []string
	msgs     []string
}

func NewYAMLParserErr(fpath string, lines []string, msgs []string) *YAMLParserError {
	var err = &YAMLParserError{
		lines: lines,
		msgs: msgs,
	}
	err.SetErrorType(ERROR_YAML_PARSER_ERROR)
	err.SetErrorFilePath(fpath)
	err.SetCallerByStackFrameSkip(2)
	return err
}


func (e *YAMLParserError) Error() string {
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

	e.SetMessage( "\n" + strings.Join(result, "\n"))
	return e.FileError.Error()
}
