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
    "github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
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
)

type TestCaseError struct {
    errorMessage string
}

func NewTestCaseError(errorMessage string) *TestCaseError {
    return &TestCaseError{
        errorMessage: errorMessage,
    }
}

func (e *TestCaseError) Error() string {
    return e.errorMessage
}

type BaseErr struct {
    FileName string
    LineNum  int
    Message  string
}

func (e *BaseErr) Error() string {
    return fmt.Sprintf("%s [%d]: %s\n", e.FileName, e.LineNum, e.Message)
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

type ErrorManifestFileNotFound struct {
    BaseErr
    errorType string
}

func NewErrorManifestFileNotFound(errMessage string) *ErrorManifestFileNotFound {
    _, fn, lineNum, _ := runtime.Caller(1)
    var err = &ErrorManifestFileNotFound{
        errorType: wski18n.T(MANIFEST_NOT_FOUND),
    }
    //err.SetFileName(fn)
    err.SetFileName(filepath.Base(fn))
    err.SetLineNum(lineNum)
    err.SetMessage(errMessage)
    return err
}

func (e *ErrorManifestFileNotFound) Error() string {
    if e.errorType == "" {
        return fmt.Sprintf("%s [%d]: %s\n", e.FileName, e.LineNum, e.Message)
    }
    return fmt.Sprintf("%s [%d]: %s ==> %s\n", e.FileName, e.LineNum, e.errorType, e.Message)
}

type InputYamlFileError struct {
    BaseErr
    errorType string
}

func NewInputYamlFileError(errMessage string) *InputYamlFileError {
    _, fn, lineNum, _ := runtime.Caller(1)
    var err = &InputYamlFileError{
        errorType: wski18n.T(INVALID_YAML_INPUT),
    }
    err.SetFileName(filepath.Base(fn))
    err.SetLineNum(lineNum)
    err.SetMessage(errMessage)
    return err
}

func (e *InputYamlFileError) SetErrorType(errorType string) {
    e.errorType = errorType
}

func (e *InputYamlFileError) Error() string {
    if e.errorType == "" {
        return fmt.Sprintf("%s [%d]: %s\n", e.FileName, e.LineNum, e.Message)
    }
    return fmt.Sprintf("%s [%d]: %s %s\n", e.FileName, e.LineNum, e.errorType, e.Message)
}

type InputYamlFormatError struct {
    InputYamlFileError
}

func NewInputYamlFormatError(errMessage string) *InputYamlFormatError {
    _, fn, lineNum, _ := runtime.Caller(1)
    var err = &InputYamlFormatError{}
    err.SetErrorType(wski18n.T(INVALID_YAML_FORMAT))
    err.SetFileName(fn)
    err.SetLineNum(lineNum)
    err.SetMessage(errMessage)
    return err
}

type WhiskClientError struct {
    BaseErr
    errorType string
    errorCode int
}

func NewWhiskClientError(errMessage string, code int) *WhiskClientError {
    _, fn, lineNum, _ := runtime.Caller(1)
    var err = &WhiskClientError{
        errorType: wski18n.T(OPENWHISK_CLIENT_ERROR),
        errorCode: code,
    }
    err.SetFileName(fn)
    err.SetLineNum(lineNum)
    err.SetMessage(errMessage)
    return err
}

func (e *WhiskClientError) Error() string {
    return fmt.Sprintf("%s [%d]: %s =====> %s Error code: %d.\n", e.FileName, e.LineNum, e.errorType, e.Message, e.errorCode)
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
