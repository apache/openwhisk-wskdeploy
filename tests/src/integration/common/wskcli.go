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

package common

import (
	"os"
	"os/exec"
)

const wsk_cli_cmd = "wsk"

type WskCLI struct {
	Path string
	Dir  string
}

func NewWskCLI() *WskCLI {
	return NewWskCLIWithPath(os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/")
}

func NewWskCLIWithPath(path string) *WskCLI {
	var dep WskCLI
	dep.Path = wsk_cli_cmd
	dep.Dir = path
	return &dep
}

func (wskcli *WskCLI) runCommand(s ...string) ([]byte, error) {
	command := exec.Command(wskcli.Path, s...)
	command.Dir = wskcli.Dir

	printCommand(command)
	output, err := command.CombinedOutput()
	printOutput(output)
	printError(err)

	return output, err
}

func (wskcli *WskCLI) Invoke(s ...string) ([]byte, error) {
	args := []string{"action", "invoke", "--result"}
	args = append(args, s...)
	return wskcli.runCommand(args...)
}

