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

func PrintOpenWhiskError(message string) {
	outputStream := colorable.NewColorableStderr()
	fmsg := fmt.Sprintf( STR_PREFIXED_MESSAGE, wski18n.T(wski18n.ID_MSG_PREFIX_ERROR), message)
	fmt.Fprintf(outputStream, color.RedString(fmsg))
}

func PrintlnOpenWhiskError(message string) {
	PrintOpenWhiskError(message + "\n")
}

func PrintOpenWhiskFromError(err error) {
	PrintOpenWhiskError(err.Error())
}

func PrintOpenWhiskWarning(message string) {
	outputStream := colorable.NewColorableStdout()
	fmsg := fmt.Sprintf( STR_PREFIXED_MESSAGE, wski18n.T(wski18n.ID_MSG_PREFIX_WARNING), message)
	fmt.Fprintf(outputStream, color.YellowString(fmsg))
}

func PrintlnOpenWhiskWarning(message string) {
	PrintOpenWhiskWarning(message + "\n")
}

func PrintOpenWhiskSuccess(message string) {
	outputStream := colorable.NewColorableStdout()
	fmsg := fmt.Sprintf( STR_PREFIXED_MESSAGE, wski18n.T(wski18n.ID_MSG_PREFIX_SUCCESS), message)
	fmt.Fprintf(outputStream, color.GreenString(fmsg))
}

func PrintlnOpenWhiskSuccess(message string) {
	PrintOpenWhiskSuccess(message + "\n")
}

func PrintOpenWhiskStatus(message string) {
	outputStream := colorable.NewColorableStdout()
	fmsg := fmt.Sprintf( STR_PREFIXED_MESSAGE, wski18n.T(wski18n.ID_MSG_PREFIX_INFO), message)
	fmt.Fprintf(outputStream, color.CyanString(fmsg))
}

func PrintlnOpenWhiskStatus(message string) {
	PrintOpenWhiskStatus(message + "\n")
}

func PrintlnOpenWhiskOutput(message string) {
   	fmt.Println(message)
}

func PrintOpenWhiskVerbose(verbose bool, message string) {
	if verbose{
		PrintlnOpenWhiskOutput(message)
	}
}
