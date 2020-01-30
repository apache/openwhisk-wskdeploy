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
	"github.com/apache/openwhisk-wskdeploy/wski18n"
	"github.com/apache/openwhisk-wskdeploy/wskprint"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// Error message compositional strings
	STR_UNKNOWN_VALUE         = "Unknown value"
	STR_COMMAND               = "Command"
	STR_ERROR_CODE            = "Error code"
	STR_FILE                  = "File"
	STR_PARAMETER             = "Parameter"
	STR_TYPE                  = "Type"
	STR_EXPECTED              = "Expected"
	STR_ACTUAL                = "Actual"
	STR_NEWLINE               = "\n"
	STR_ACTION                = "Action"
	STR_RUNTIME               = "Runtime"
	STR_SUPPORTED_RUNTIMES    = "Supported Runtimes"
	STR_HTTP_STATUS           = "HTTP Response Status"
	STR_HTTP_BODY             = "HTTP Response Body"
	STR_SUPPORTED_WEB_EXPORTS = "Supported Web Exports"
	STR_WEB_EXPORT            = "web-export"
	STR_API                   = "API"
	STR_API_METHOD            = "API gateway method"
	STR_API_SUPPORTED_METHODS = "API gateway supported methods"

	// Formatting
	STR_INDENT_1 = "==>"

	// Error Types
	ERROR_COMMAND_FAILED                  = "ERROR_COMMAND_FAILED"
	ERROR_WHISK_CLIENT_ERROR              = "ERROR_WHISK_CLIENT_ERROR"
	ERROR_WHISK_CLIENT_INVALID_CONFIG     = "ERROR_WHISK_CLIENT_INVALID_CONFIG"
	ERROR_FILE_READ_ERROR                 = "ERROR_FILE_READ_ERROR"
	ERROR_MANIFEST_FILE_NOT_FOUND         = "ERROR_MANIFEST_FILE_NOT_FOUND"
	ERROR_YAML_FILE_FORMAT_ERROR          = "ERROR_YAML_FILE_FORMAT_ERROR"
	ERROR_YAML_PARSER_ERROR               = "ERROR_YAML_PARSER_ERROR"
	ERROR_YAML_PARAMETER_TYPE_MISMATCH    = "ERROR_YAML_PARAMETER_TYPE_MISMATCH"
	ERROR_YAML_INVALID_PARAMETER_TYPE     = "ERROR_YAML_INVALID_PARAMETER_TYPE"
	ERROR_YAML_INVALID_RUNTIME            = "ERROR_YAML_INVALID_RUNTIME"
	ERROR_YAML_INVALID_WEB_EXPORT         = "ERROR_YAML_INVALID_WEB_EXPORT"
	ERROR_YAML_INVALID_API                = "ERROR_YAML_INVALID_API"
	ERROR_YAML_INVALID_API_GATEWAY_METHOD = "ERROR_YAML_INVALID_API_GATEWAY_METHOD"
	ERROR_RUNTIME_PARSER_FAILURE          = "ERROR_RUNTIME_PARSER_FAILURE"
	ERROR_ACTION_ANNOTATION               = "ERROR_ACTION_ANNOTATION"
)

/*
 * BaseError
 */
type WskDeployBaseErr struct {
	ErrorType     string
	FileName      string
	LineNum       int
	Message       string
	MessageFormat string
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

func (e *WskDeployBaseErr) GetMessage() string {
	return e.Message
}

func (e *WskDeployBaseErr) GetMessageFormat() string {
	return e.MessageFormat
}

func (e *WskDeployBaseErr) SetMessage(message interface{}) {

	if message != nil {
		switch message.(type) {
		case string:
			e.Message = message.(string)
		case error:
			err := message.(error)
			e.appendErrorDetails(err)
		}
	}
}

func (e *WskDeployBaseErr) AppendDetail(detail string) {
	e.appendDetail(detail)
}

func (e *WskDeployBaseErr) appendDetail(detail string) {
	fmt := fmt.Sprintf("\n%s %s", STR_INDENT_1, detail)
	e.Message = e.Message + fmt
}

func (e *WskDeployBaseErr) appendErrorDetails(err error) {
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

func NewWhiskClientError(errorMessage string, code int, response *http.Response) *WhiskClientError {
	var err = &WhiskClientError{
		ErrorCode: code,
	}
	err.SetErrorType(ERROR_WHISK_CLIENT_ERROR)
	err.SetCallerByStackFrameSkip(2)
	err.SetMessageFormat("%s: %d: %s")
	var str = fmt.Sprintf(err.MessageFormat, STR_ERROR_CODE, code, errorMessage)
	if response != nil {
		responseData, _ := ioutil.ReadAll(response.Body)
		err.SetMessageFormat("%s: %d: %s: %s: %s %s: %s")
		// do not add body in case of a success
		// when response.Status is 200, response.Body contains the entire action source code
		// we should not expose the action source when the HTTP request was successful
		if strings.Contains(response.Status, "200 OK") {
			str = fmt.Sprintf(err.MessageFormat, STR_ERROR_CODE, code, errorMessage, STR_HTTP_STATUS, response.Status, "", "")
		} else {
			str = fmt.Sprintf(err.MessageFormat, STR_ERROR_CODE, code, errorMessage, STR_HTTP_STATUS, response.Status, STR_HTTP_BODY, string(responseData))
		}
	}
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
	var err = &WhiskClientInvalidConfigError{}
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
	return fmt.Sprintf("%s [%d]: [%s]: "+STR_FILE+": [%s]: %s\n",
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
	var err = &FileReadError{}
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
	var err = &ErrorManifestFileNotFound{}
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
	var err = &YAMLFileFormatError{}
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
		ActualType:   actualType,
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
	Parameter  string
	ActualType string
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
	lines []string
	msgs  []string
}

func NewYAMLParserErr(fpath string, msg interface{}) *YAMLParserError {
	var err = &YAMLParserError{}
	err.SetErrorType(ERROR_YAML_PARSER_ERROR)
	err.SetErrorFilePath(fpath)
	err.SetCallerByStackFrameSkip(2)
	err.SetMessage(msg)
	return err
}

/*
 * InvalidRuntime
 */
type InvalidRuntimeError struct {
	FileError
	Runtime           string
	SupportedRuntimes []string
}

func NewInvalidRuntimeError(errMessage string, fpath string, action string, runtime string, supportedRuntimes []string) *InvalidRuntimeError {
	var err = &InvalidRuntimeError{
		SupportedRuntimes: supportedRuntimes,
	}
	err.SetErrorFilePath(fpath)
	err.SetErrorType(ERROR_YAML_INVALID_RUNTIME)
	err.SetCallerByStackFrameSkip(2)
	str := fmt.Sprintf("%s %s [%s]: %s [%s]: %s [%s]",
		errMessage,
		STR_ACTION, action,
		STR_RUNTIME, runtime,
		STR_SUPPORTED_RUNTIMES, strings.Join(supportedRuntimes, ", "))
	err.SetMessage(str)
	return err
}

/*
 * InvalidWebExport
 */
type InvalidWebExportError struct {
	FileError
	Webexport           string
	SupportedWebexports []string
}

func NewInvalidWebExportError(fpath string, action string, webexport string, supportedWebexports []string) *InvalidWebExportError {
	var err = &InvalidWebExportError{
		SupportedWebexports: supportedWebexports,
	}
	err.SetErrorFilePath(fpath)
	err.SetErrorType(ERROR_YAML_INVALID_WEB_EXPORT)
	err.SetCallerByStackFrameSkip(2)
	str := fmt.Sprintf("%s [%s]: %s [%s]: %s [%s]",
		STR_ACTION, action,
		STR_WEB_EXPORT, webexport,
		STR_SUPPORTED_WEB_EXPORTS, strings.Join(supportedWebexports, ", "))
	err.SetMessage(str)
	return err
}

/*
 * Invalid API Gateway Method
 */
type InvalidAPIGatewayMethodError struct {
	FileError
	method           string
	SupportedMethods []string
}

func NewInvalidAPIGatewayMethodError(fpath string, api string, method string, supportedMethods []string) *InvalidAPIGatewayMethodError {
	var err = &InvalidAPIGatewayMethodError{
		SupportedMethods: supportedMethods,
	}
	err.SetErrorFilePath(fpath)
	err.SetErrorType(ERROR_YAML_INVALID_API_GATEWAY_METHOD)
	err.SetCallerByStackFrameSkip(2)
	str := fmt.Sprintf("%s [%s]: %s [%s]: %s [%s]",
		STR_API, api,
		STR_API_METHOD, method,
		STR_API_SUPPORTED_METHODS, strings.Join(supportedMethods, ", "))
	err.SetMessage(str)
	return err
}

/*
 * Invalid Web Action API
 */
type InvalidWebActionAPIError struct {
	WskDeployBaseErr
}

func NewInvalidWebActionError(apiName string, actionName string, isSequence bool) *InvalidWebActionAPIError {
	var err = &InvalidWebActionAPIError{
	}

	i18nErrorID := wski18n.ID_ERR_API_MISSING_WEB_ACTION_X_action_X_api_X

	if isSequence {
		i18nErrorID = wski18n.ID_ERR_API_MISSING_WEB_SEQUENCE_X_sequence_X_api_X
	}

	errString := wski18n.T(i18nErrorID,
		map[string]interface{}{
			wski18n.KEY_SEQUENCE: actionName,
			wski18n.KEY_API:      apiName})
	wskprint.PrintOpenWhiskWarning(errString)

	err.SetErrorType(ERROR_YAML_INVALID_API)
	err.SetCallerByStackFrameSkip(2)
	err.SetMessage(errString)
	return err
}

/*
 * Failed to Retrieve/Parse Runtime
 */
type RuntimeParserError struct {
	WskDeployBaseErr
}

func NewRuntimeParserError(errorMsg string) *RuntimeParserError {
	var err = &RuntimeParserError{}
	err.SetErrorType(ERROR_RUNTIME_PARSER_FAILURE)
	err.SetCallerByStackFrameSkip(2)
	err.SetMessage(errorMsg)
	return err
}

/*
 * Failed to Retrieve/Parse Runtime
 */
type DeployError struct {
	WskDeployBaseErr
}

func NewActionSecureKeyError(errorMsg string) *DeployError {
	var err = &DeployError{}
	err.SetErrorType(ERROR_ACTION_ANNOTATION)
	err.SetCallerByStackFrameSkip(2)
	err.SetMessage(errorMsg)
	return err
}

func IsCustomError(err error) bool {

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

func AppendDetailToErrorMessage(detail string, add string, location int) string {

	if len(detail) == 0 {
		detail = "\n"
	}
	_, fname, lineNum, _ := runtime.Caller(location)
	detail += fmt.Sprintf("  >> %s [%v]: %s", filepath.Base(fname), lineNum, add)
	return detail
}
