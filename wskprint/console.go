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
	"github.com/apache/openwhisk-wskdeploy/wski18n"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	STR_PREFIXED_MESSAGE = "%s: %s"
)

var (
	clrInfo      = color.New(color.FgCyan)
	clrWarning   = color.New(color.FgYellow)
	clrError     = color.New(color.FgRed)
	clrSuccess   = color.New(color.FgGreen)
	clrTitleInfo = color.New(color.FgCyan).Add(color.Underline)
)

func PrintOpenWhiskError(message string) {
	outputStream := colorable.NewColorableStderr()
	fmt.Fprintf(outputStream, clrError.Sprintf(STR_PREFIXED_MESSAGE,
		wski18n.T(wski18n.ID_MSG_PREFIX_ERROR), message))
}

func PrintOpenWhiskFromError(err error) {
	PrintOpenWhiskError(err.Error())
}

func PrintOpenWhiskWarning(message string) {
	if DetectVerbose() {
		outputStream := colorable.NewColorableStdout()
		fmt.Fprintf(outputStream, clrWarning.Sprintf(STR_PREFIXED_MESSAGE,
			wski18n.T(wski18n.ID_MSG_PREFIX_WARNING), message))
	}
}

func PrintlnOpenWhiskWarning(message string) {
	if DetectVerbose() {
		PrintOpenWhiskWarning(message + "\n")
	}
}

func PrintOpenWhiskSuccess(message string) {
	outputStream := colorable.NewColorableStdout()
	fmt.Fprintf(outputStream, clrSuccess.Sprintf(STR_PREFIXED_MESSAGE,
		wski18n.T(wski18n.ID_MSG_PREFIX_SUCCESS), message))
}

func PrintlnOpenWhiskSuccess(message string) {
	PrintOpenWhiskSuccess(message + "\n")
}

func PrintOpenWhiskInfo(message string) {
	outputStream := colorable.NewColorableStdout()
	fmt.Fprintf(outputStream, clrInfo.Sprintf(STR_PREFIXED_MESSAGE,
		wski18n.T(wski18n.ID_MSG_PREFIX_INFO), message))
}

func PrintlnOpenWhiskInfo(message string) {
	PrintOpenWhiskInfo(message + "\n")
}

func PrintlnOpenWhiskInfoTitle(message string) {
	outputStream := colorable.NewColorableStdout()
	fmt.Fprintf(outputStream, clrTitleInfo.Sprintf(STR_PREFIXED_MESSAGE,
		wski18n.T(wski18n.ID_MSG_PREFIX_INFO), message))
}

func PrintlnOpenWhiskOutput(message string) {
	fmt.Println(message)
}

func PrintOpenWhiskVerboseTitle(verbose bool, message string) {
	if verbose {
		PrintlnOpenWhiskInfoTitle(message)
	}
}

func PrintOpenWhiskVerbose(verbose bool, message string) {
	if verbose {
		PrintOpenWhiskInfo(message)
	}
}

func PrintlnOpenWhiskVerbose(verbose bool, message string) {
	PrintOpenWhiskVerbose(verbose, message+"\n")
}

func PrintlnOpenWhiskTrace(trace bool, message string) {
	if trace {
		_, fname, lineNum, _ := runtime.Caller(2)
		out := fmt.Sprintf("%s [%v]: %s\n", filepath.Base(fname), lineNum, message)
		PrintOpenWhiskVerbose(trace, out)
	}
}

// Display "trace" output if either param is true OR we are running Go test verbose (i.e., "go test -v")
// Typical Args for "go test" looks as follows:
// arg[0] = [/var/folders/nj/<uuid>/T/<build-id>/github.com/apache/openwhisk-wskdeploy/deployers/_test/deployers.test
// arg[1] = -test.v=true
// arg[2] = -test.run=TestDeploymentReader_PackagesBindTrigger]
func DetectGoTestVerbose() bool {
	arguments := os.Args
	for i := range arguments {
		if strings.HasPrefix(arguments[i], "-test.v=true") {
			return true
		}
	}
	return false
}

func DetectVerbose() bool {
	arguments := os.Args
	for i := range arguments {
		if strings.HasPrefix(arguments[i], "-v") ||
			strings.HasPrefix(arguments[i], "--verbose") {
			return true
		}
	}
	return false
}

//func PrintOpenWhiskBanner(verbose bool) {
//	if verbose {
//		PrintlnOpenWhiskOutput("         ____      ___                   _    _ _     _     _\n        /\\   \\    / _ \\ _ __   ___ _ __ | |  | | |__ (_)___| | __\n   /\\  /__\\   \\  | | | | '_ \\ / _ \\ '_ \\| |  | | '_ \\| / __| |/ /\n  /  \\____ \\  /  | |_| | |_) |  __/ | | | |/\\| | | | | \\__ \\   <\n  \\   \\  /  \\/    \\___/| .__/ \\___|_| |_|__/\\__|_| |_|_|___/_|\\_\\ \n   \\___\\/              |_|\n")
//	}
//}
