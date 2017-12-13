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

package wskderrors

import (
	"fmt"
	"runtime"
	"strings"
	"path/filepath"
)

const (
	// Error message compositional strings
	STR_UNKNOWN_VALUE = "Unknown value"
	STR_COMMAND = "Command"
	STR_ERROR_CODE = "Error code"
	STR_FILE = "File"
	STR_PARAMETER = "Parameter"
	STR_TYPE = "Type"
	STR_EXPECTED = "Expected"
	STR_ACTUAL = "Actual"
	STR_NEWLINE = "\n"

	// Formatting
	STR_INDENT_1 = "==>"

	// Error Types
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
type WskDeployBaseErr struct {
	ErrorType 	string
	FileName  	string
	LineNum   	int
	Message   	string
	MessageFormat	string
}

func NewWskDeployBaseError(typ string, fn string, ln int, msg string) *WskDeployBaseErr {
	var err = &WskDeployBaseErr{
		ErrorType: typ,
		FileName:  fn,
		LineNum:   ln,
	}
	err.SetMessage(msg)
	return err
}

func (e *WskDeployBaseErr) Error() string {
	return fmt.Sprintf("%s [%d]: [%s]: %s\n", e.FileName, e.LineNum, e.ErrorType, e.Message)
}

func (e *WskDeployBaseErr) SetFileName(fileName string) {
	e.FileName = filepath.Base(fileName)
}

func (e *WskDeployBaseErr) SetLineNum(lineNum int) {
	e.LineNum = lineNum
}

func (e *WskDeployBaseErr) SetErrorType(errorType string) {
	e.ErrorType = errorType
}

func (e *WskDeployBaseErr) SetMessageFormat(fmt string) {
	e.MessageFormat = fmt
}

func (e *WskDeployBaseErr) GetMessage()(string) {
	return e.Message
}

func (e *WskDeployBaseErr) GetMessageFormat()(string) {
	return e.MessageFormat
}

func (e *WskDeployBaseErr) SetMessage(message interface{}) {

	if message != nil{
		switch message.(type) {
		case string:
			e.Message = message.(string)
		case error:
			err := message.(error)
			e.appendErrorDetails(err)
		}
	}
}

func (e *WskDeployBaseErr) appendDetail(detail string){
	fmt := fmt.Sprintf("\n%s %s", STR_INDENT_1, detail)
	e.Message = e.Message + fmt
}

func (e *WskDeployBaseErr) appendErrorDetails(err error){
	if err != nil {
		errorMsg := err.Error()
		var detailMsg string
		msgs := strings.Split(errorMsg, STR_NEWLINE)
		for i := 0; i < len(msgs); i++ {
			detailMsg = msgs[i]
			e.appendDetail(strings.TrimSpace(detailMsg))
		}
	}
}

// func Caller(skip int) (pc uintptr, file string, line int, ok bool)
func (e *WskDeployBaseErr) SetCallerByStackFrameSkip(skip int) {
	_, fname, lineNum, _ := runtime.Caller(skip)
	e.SetFileName(fname)
	e.SetLineNum(lineNum)
}

/*
 * CommandError
 */
type CommandError struct {
	WskDeployBaseErr
	Command string
}

func NewCommandError(cmd string, errorMessage string) *CommandError {
	var err = &CommandError{
		Command: cmd,
	}
	err.SetErrorType(ERROR_COMMAND_FAILED)
	err.SetCallerByStackFrameSkip(2)
	err.SetMessageFormat("%s: [%s]: %s")
	str := fmt.Sprintf(err.MessageFormat, STR_COMMAND, cmd, errorMessage)
	err.SetMessage(str)
	return err
}

/*
 * WhiskClientError
 */
type WhiskClientError struct {
	WskDeployBaseErr
	ErrorCode int
}

func NewWhiskClientError(errorMessage string, code int) *WhiskClientError {
	var err = &WhiskClientError{
		ErrorCode: code,
	}
	err.SetErrorType(ERROR_WHISK_CLIENT_ERROR)
	err.SetCallerByStackFrameSkip(2)
	err.SetMessageFormat("%s: %d: %s")
	str := fmt.Sprintf(err.MessageFormat, STR_ERROR_CODE, code, errorMessage)
	err.SetMessage(str)
	return err
}

/*
 * WhiskClientInvalidConfigError
 */
type WhiskClientInvalidConfigError struct {
	WskDeployBaseErr
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
	WskDeployBaseErr
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
	return fmt.Sprintf("%s [%d]: [%s]: " + STR_FILE + ": [%s]: %s\n",
		e.FileName,
		e.LineNum,
		e.ErrorType,
		e.ErrorFileName,
		e.Message)
}

/*
 * FileReadError
 */

type FileReadError struct {
	FileError
}

func NewFileReadError(fpath string, errMessage interface{}) *FileReadError {
	var err = &FileReadError{
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

func NewErrorManifestFileNotFound(fpath string, errMessage interface{}) *ErrorManifestFileNotFound {
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

func NewYAMLFileFormatError(fpath string, errorMessage interface{}) *YAMLFileFormatError {
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
	err.SetMessageFormat("%s [%s]: %s %s: [%s], %s: [%s]")
	str := fmt.Sprintf(err.MessageFormat,
		STR_PARAMETER, param,
		STR_TYPE,
		STR_EXPECTED, expectedType,
		STR_ACTUAL, actualType)
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
	err.SetMessageFormat("%s [%s]: %s [%s]")
	str := fmt.Sprintf(err.MessageFormat,
		STR_PARAMETER, param,
		STR_TYPE, actualType)
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

func NewYAMLParserErr(fpath string, msg interface{}) *YAMLParserError {
	var err = &YAMLParserError{

	}
	err.SetErrorType(ERROR_YAML_PARSER_ERROR)
	err.SetErrorFilePath(fpath)
	err.SetCallerByStackFrameSkip(2)
        err.SetMessage(msg)
	return err
}

func IsCustomError( err error ) bool {

	switch err.(type) {

	case *CommandError:
	case *WhiskClientError:
	case *WhiskClientInvalidConfigError:
	case *FileError:
	case *FileReadError:
	case *ErrorManifestFileNotFound:
	case *YAMLFileFormatError:
	case *ParameterTypeMismatchError:
	case *InvalidParameterTypeError:
	case *YAMLParserError:
		return true
	}
	return false
}
