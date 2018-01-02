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

	ID_I18B_PREFIX_ERROR 	= "msg_prefix_error"
	ID_I18B_PREFIX_WARNING	= "msg_prefix_warning"
	ID_I18B_PREFIX_SUCCESS	= "msg_prefix_success"
	ID_I18B_PREFIX_INFO	= "msg_prefix_info"
)

func PrintOpenWhiskError(message string) {
	outputStream := colorable.NewColorableStderr()
	fmsg := fmt.Sprintf( STR_PREFIXED_MESSAGE, wski18n.T(ID_I18B_PREFIX_ERROR), message)
	fmt.Fprintf(outputStream, color.RedString(fmsg))
}

func PrintOpenWhiskFromError(err error) {
	PrintOpenWhiskError(err.Error())
}

func PrintOpenWhiskWarning(message string) {
	outputStream := colorable.NewColorableStdout()
	fmsg := fmt.Sprintf( STR_PREFIXED_MESSAGE, wski18n.T(ID_I18B_PREFIX_WARNING), message)
	fmt.Fprintf(outputStream, color.YellowString(fmsg))
}

func PrintOpenWhiskSuccess(message string) {
	outputStream := colorable.NewColorableStdout()
	fmsg := fmt.Sprintf( STR_PREFIXED_MESSAGE, wski18n.T(ID_I18B_PREFIX_SUCCESS), message)
	fmt.Fprintf(outputStream, color.GreenString(fmsg))
}

func PrintlnOpenWhiskOutput(output string) {
    fmt.Println(output)
}
