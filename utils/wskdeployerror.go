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
    INVALID_YAML_INPUT = "Invalid input of Yaml file"
    INVALID_YAML_FORMAT = "Invalid input of Yaml format"
    OPENWHISK_CLIENT_ERROR = "OpenWhisk Client Error"
    PARAMETER_TYPE_MISMATCH = "Parameter type mismatch error"
    MANIFEST_NOT_FOUND = INVALID_YAML_INPUT  // TODO{} This should be a unique message.
    UNKNOWN = "Unknown"
    UNKNOWN_VALUE = "Unknown value"
    LINE = "line"

    ERROR_COMMAND_FAILED = "ERROR_COMMAND_FAILED"
    ERROR_MANIFEST_FILE_NOT_FOUND = "ERROR_MANIFEST_FILE_NOT_FOUND"
    ERROR_YAML_FILE_ERROR = "ERROR_YAML_FILE_ERROR"
    ERROR_YAML_FORMAT_ERROR = "ERROR_YAML_FORMAT_ERROR"
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
 * YAMLFileError
 */
type InputYamlFileError struct {
    BaseErr
}

func NewInputYamlFileError(errMessage string) *InputYamlFileError {
    var err = &InputYamlFileError{
        //errorType: wski18n.T(INVALID_YAML_INPUT),
    }
    err.SetErrorType(ERROR_YAML_FILE_ERROR)
    err.SetCallerByStackFrameSkip(2)
    err.SetMessage(errMessage)
    return err
}

func (e *InputYamlFileError) Error() string {
    return e.BaseErr.Error()
}

/*
 * YAMLFormatError
 */
type InputYamlFormatError struct {
    InputYamlFileError
}

func NewInputYamlFormatError(errMessage string) *InputYamlFormatError {
    var err = &InputYamlFormatError{
        //err.SetErrorType(wski18n.T(INVALID_YAML_FORMAT))
    }
    err.SetErrorType(ERROR_YAML_FORMAT_ERROR)
    err.SetCallerByStackFrameSkip(2)
    err.SetMessage(errMessage)
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
    err.SetErrorType(ERROR_YAML_FORMAT_ERROR)
    err.SetCallerByStackFrameSkip(2)
    str := fmt.Sprintf("Code: %d: %s", code, errorMessage)
    err.SetMessage(str)
    return err
}

func (e *WhiskClientError) Error() string {
    return e.BaseErr.Error()
}

type InvalidWskpropsError struct {
    BaseErr
}

func NewInvalidWskpropsError(errMessage string) *InvalidWskpropsError {
    _, fn, lineNum, _ := runtime.Caller(1)
    var err = &InvalidWskpropsError{}
    err.SetFileName(fn)
    err.SetLineNum(lineNum)
    err.SetMessage(errMessage)
    return err
}

type ParserErr struct {
    BaseErr
    YamlFile string
    lines []string
    msgs []string
}

func NewParserErr(yamlFile string, lines []string, msgs []string) *ParserErr {
    _, fn, line, _ := runtime.Caller(1)
    var err = &ParserErr{
        YamlFile: yamlFile,
        lines: lines,
        msgs: msgs,
    }
    err.SetFileName(fn)
    err.SetLineNum(line)
    return err
}

func (e *ParserErr) Error() string {
    result := make([]string, len(e.msgs))
    var fn = filepath.Base(e.FileName)

    for index, msg := range e.msgs {
        var s string
        if e.lines == nil || e.lines[index] == UNKNOWN {
            s = fmt.Sprintf("====> %s", msg)
        } else{
            s = fmt.Sprintf("====> Line [%v]: %s", e.lines[index], msg)
        }
        result[index] = s
    }
    return fmt.Sprintf("\n==> %s [%d]: Failed to parse the yaml file: %s: \n%s", fn, e.LineNum, e.YamlFile, strings.Join(result, "\n"))
}

type ParameterTypeMismatchError struct {
    BaseErr
    errorType string
    expectedType string
    actualType string
}

func (e *ParameterTypeMismatchError) Error() string {
    if e.errorType == "" {
        return fmt.Sprintf("%s [%d]: %s\n", e.FileName, e.LineNum, e.Message)
    }
    return fmt.Sprintf("%s [%d]: %s ==> %s\n", e.FileName, e.LineNum, e.errorType, e.Message)
}

func NewParameterTypeMismatchError(errMessage string, expectedType string, actualType string) *ParameterTypeMismatchError {
    _, fn, lineNum, _ := runtime.Caller(1)
    var err = &ParameterTypeMismatchError{
        // TODO{} add i18n
        //errorType: wski18n.T(PARAMETER_TYPE_MISMATCH),
        errorType: PARAMETER_TYPE_MISMATCH,
        expectedType: expectedType,
        actualType: actualType,
    }

    err.SetFileName(filepath.Base(fn))
    err.SetLineNum(lineNum)
    err.SetMessage(errMessage)
    return err
}
