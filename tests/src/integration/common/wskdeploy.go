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

const cmd = "wskdeploy"

type Wskdeploy struct {
    Path string
    Dir  string
}

func NewWskdeploy() *Wskdeploy {
    return NewWskWithPath(os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/")
}

func NewWskWithPath(path string) *Wskdeploy {
    var dep Wskdeploy
    dep.Path = cmd
    dep.Dir = path
    return &dep
}

func (wskdeploy *Wskdeploy) RunCommand(s ...string) ([]byte, error) {
    command := exec.Command(wskdeploy.Path, s...)
    command.Dir = wskdeploy.Dir
    return command.CombinedOutput()
}

func (wskdeploy *Wskdeploy) Deploy(manifestPath string, deploymentPath string) ([]byte, error) {
    return wskdeploy.RunCommand("-m", manifestPath, "-d", deploymentPath)
}

func (wskdeploy *Wskdeploy) Undeploy(manifestPath string, deploymentPath string) ([]byte, error) {
    return wskdeploy.RunCommand("undeploy", "-m", manifestPath, "-d", deploymentPath)
}
