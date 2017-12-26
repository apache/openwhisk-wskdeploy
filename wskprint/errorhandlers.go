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

func PrintOpenWhiskError(err error) {
	outputStream := colorable.NewColorableStderr()
	fmt.Fprintf(outputStream, "%s%s", color.RedString(wski18n.T("Error: ")), err.Error())
}

func PrintOpenWhiskOutput(output string) {
    outputStream := colorable.NewColorableStdout()
    fmt.Fprintf(outputStream, "%s", color.GreenString(output))
}

func PrintOpenWhiskOutputln(output string) {
    fmt.Println(output)
}

func PrintOpenWhiskErrorMessage(err string) {
    outputStream := colorable.NewColorableStderr()
    fmt.Fprintf(outputStream, "%s", color.RedString(err))
}
