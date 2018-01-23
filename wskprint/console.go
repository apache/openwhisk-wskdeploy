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

package wskprint

import (
	"fmt"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

const(
	STR_PREFIXED_MESSAGE = "%s: %s"
)

var(
	COLOR_INFO = color.New(color.FgCyan)
	COLOR_WARNING = color.New(color.FgYellow)
	COLOR_ERROR = color.New(color.FgRed)
	COLOR_SUCCESS = color.New(color.FgGreen)
	COLOR_TITLE = color.New(color.FgRed).Add(color.Underline)
)

func PrintOpenWhiskError(message string) {
	outputStream := colorable.NewColorableStderr()
	fmt.Fprintf(outputStream,COLOR_ERROR.Sprintf( STR_PREFIXED_MESSAGE,
		wski18n.T(wski18n.ID_MSG_PREFIX_ERROR), message))
}

func PrintOpenWhiskFromError(err error) {
	PrintOpenWhiskError(err.Error())
}

func PrintOpenWhiskWarning(message string) {
	outputStream := colorable.NewColorableStdout()
	fmt.Fprintf(outputStream,COLOR_WARNING.Sprintf( STR_PREFIXED_MESSAGE,
		wski18n.T(wski18n.ID_MSG_PREFIX_WARNING), message))
}

func PrintlnOpenWhiskWarning(message string) {
	PrintOpenWhiskWarning(message + "\n")
}

func PrintOpenWhiskSuccess(message string) {
	outputStream := colorable.NewColorableStdout()
	fmt.Fprintf(outputStream,COLOR_SUCCESS.Sprintf( STR_PREFIXED_MESSAGE,
		wski18n.T(wski18n.ID_MSG_PREFIX_SUCCESS), message))
}

func PrintlnOpenWhiskSuccess(message string) {
	PrintOpenWhiskSuccess(message + "\n")
}

func PrintOpenWhiskInfo(message string) {
	outputStream := colorable.NewColorableStdout()
	fmt.Fprintf(outputStream,COLOR_INFO.Sprintf( STR_PREFIXED_MESSAGE,
		wski18n.T(wski18n.ID_MSG_PREFIX_INFO), message))
}


func PrintlnOpenWhiskInfo(message string) {
	PrintOpenWhiskInfo(message + "\n")
}

func PrintlnOpenWhiskInfoTitle(message string) {
	outputStream := colorable.NewColorableStdout()
	fmt.Fprintf(outputStream,COLOR_TITLE.Sprintf( STR_PREFIXED_MESSAGE,
		wski18n.T(wski18n.ID_MSG_PREFIX_INFO), message))
}

func PrintlnOpenWhiskOutput(message string) {
   	fmt.Println(message)
}

func PrintOpenWhiskVerbose(verbose bool, message string) {
	if verbose{
		PrintOpenWhiskInfo(message)
	}
}

func PrintlnOpenWhiskVerbose(verbose bool, message string) {
	PrintOpenWhiskVerbose(verbose, message + "\n")
}
