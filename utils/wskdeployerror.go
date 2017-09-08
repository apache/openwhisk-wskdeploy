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
)

const (
    INVALID_YAML_INPUT = "Invalid input of Yaml file"
    INVALID_YAML_FORMAT = "Invalid input of Yaml format"
    OPENWHISK_CLIENT_ERROR = "OpenWhisk Client Error"
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

type InputYamlFileError struct {
    BaseErr
    errorType string
}

func NewInputYamlFileError(errMessage string) *InputYamlFileError {
    _, fn, lineNum, _ := runtime.Caller(1)
    var err = &InputYamlFileError{
        errorType: wski18n.T(INVALID_YAML_INPUT),
    }
    err.SetFileName(fn)
    err.SetLineNum(lineNum)
    err.SetMessage(errMessage)
    return err
}

func (e *InputYamlFileError) SetErrorType(errorType string) {
    e.errorType = errorType
}

func (e *InputYamlFileError) Error() string {
    return fmt.Sprintf("%s [%d]: %s =====> %s\n", e.FileName, e.LineNum, e.errorType, e.Message)
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
}

func NewParserErr(yamlFile string, msg string) *ParserErr {
    _, fn, line, _ := runtime.Caller(1)
    var err = &ParserErr{
        YamlFile: yamlFile,
    }
    err.SetFileName(fn)
    err.SetLineNum(line)
    err.SetMessage(msg)
    return err
}

func (e *ParserErr) Error() string {
    return fmt.Sprintf("%s [%d]: Failed to parse the yaml file %s =====> %s\n", e.FileName, e.LineNum, e.YamlFile, e.Message)
}
