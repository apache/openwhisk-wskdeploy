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
    //"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
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

    ERROR_COMMAND_FAILED = "ERROR_COMMAND_FAILED"
    ERROR_MANIFEST_FILE_NOT_FOUND = "ERROR_MANIFEST_FILE_NOT_FOUND"
    ERROR_YAML_FILE_READ_ERROR = "ERROR_YAML_FILE_READ_ERROR"
    ERROR_YAML_FORMAT_ERROR = "ERROR_YAML_FORMAT_ERROR"
    ERROR_WHISK_CLIENT_ERROR = "ERROR_WHISK_CLIENT_ERROR"
    ERROR_WHISK_CLIENT_INVALID_CONFIG = "ERROR_WHISK_CLIENT_INVALID_CONFIG"
    ERROR_YAML_PARAMETER_TYPE_MISMATCH = "ERROR_YAML_PARAMETER_TYPE_MISMATCH"
)

/*
 * BaseError
 */
type BaseErr struct {
    ErrorType   string
    FileName    string
    LineNum     int
    Message     string
}

func (e *BaseErr) Error() string {
    if e.ErrorType == "" {
        return fmt.Sprintf("%s [%d]: %s\n", e.FileName, e.LineNum, e.Message)
    }
    return fmt.Sprintf("%s [%d]: [%s]: %s\n", e.FileName, e.LineNum, e.ErrorType, e.Message)
}

func (e *BaseErr) SetFileName(fileName string) {
    e.FileName = fileName
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
    _, fn, lineNum, _ := runtime.Caller(skip)
    e.SetFileName(filepath.Base(fn))
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
    BaseErr
    YAMLFilename string
}

func NewErrorManifestFileNotFound(errMessage string) *ErrorManifestFileNotFound {
    var err = &ErrorManifestFileNotFound{
    }
    err.SetErrorType(ERROR_MANIFEST_FILE_NOT_FOUND)
    err.SetCallerByStackFrameSkip(2)
    err.SetMessage(errMessage)
    return err
}

func (e *ErrorManifestFileNotFound) Error() string {
    return e.BaseErr.Error()
}

/*
 * YAMLFileReadError
 */
type YAMLFileReadError struct {
    BaseErr
}

func NewYAMLFileReadError(errorMessage string) *YAMLFileReadError {
    var err = &YAMLFileReadError{
        //errorType: wski18n.T(INVALID_YAML_INPUT),
    }
    err.SetErrorType(ERROR_YAML_FILE_READ_ERROR)
    err.SetCallerByStackFrameSkip(2)
    err.SetMessage(errorMessage)
    return err
}

func (e *YAMLFileReadError) Error() string {
    return e.BaseErr.Error()
}

/*
 * YAMLFormatError
 */
type YAMLFormatError struct {
    YAMLFileReadError
}

func NewYAMLFormatError(errorMessage string) *YAMLFormatError {
    var err = &YAMLFormatError{
        //err.SetErrorType(wski18n.T(INVALID_YAML_FORMAT))
    }
    err.SetErrorType(ERROR_YAML_FORMAT_ERROR)
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
        //errorType: wski18n.T(OPENWHISK_CLIENT_ERROR),
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
    BaseErr
    YamlFile string
    lines []string
    msgs []string
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
        } else{
            s = fmt.Sprintf("====> Line [%v]: %s", e.lines[index], msg)
        }
        result[index] = s
    }
    return fmt.Sprintf("\n==> %s [%d]: Failed to parse the yaml file: %s: \n%s", e.FileName, e.LineNum, e.YamlFile, strings.Join(result, "\n"))
}

/*
 * ParameterTypeMismatchError
 */
type ParameterTypeMismatchError struct {
    BaseErr
    Parameter       string
    ExpectedType    string
    ActualType      string
}

func NewParameterTypeMismatchError(param string, expectedType string, actualType string) *ParameterTypeMismatchError {
    var err = &ParameterTypeMismatchError{
        ExpectedType: expectedType,
        ActualType: actualType,
    }
    err.SetErrorType(ERROR_YAML_PARAMETER_TYPE_MISMATCH)
    err.SetCallerByStackFrameSkip(2)
    str := fmt.Sprintf( "%s [%s]: %s %s: [%s], %s: [%s]", PARAMETER, param, TYPE, EXPECTED, expectedType, ACTUAL, actualType)
    err.SetMessage(str)
    return err
}

func (e *ParameterTypeMismatchError) Error() string {
    return e.BaseErr.Error()
}

//func TestCustomErrors(){
//    err1 := NewCommandError("Deploy", "Bad error")
//    str1 := err1.Error()
//    fmt.Printf(str1)
//
//    err2 := NewErrorManifestFileNotFound("Not Found")
//    str2 := err2.Error()
//    fmt.Printf(str2)
//
//    //err3 := utils.NewWhiskClientError("Some error message", 3201 )
//    //str3 := err3.Error()
//    //fmt.Printf(str3)
//
//    err6 := NewParameterTypeMismatchError("order_info", "json", "integer")
//    str6 := err6.Error()
//    fmt.Printf(str6)
//}
