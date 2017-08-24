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
	"context"
	"time"
	"bytes"
	"fmt"
	"strings"
	"syscall"
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

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func (wskdeploy *Wskdeploy) RunCommand(args ...string) (bool) {
	// (TODO) hardcoding 30 secounds timeout here, make this configurable
	// (TODO) so that timeout can be also be read from env. variable
	timeout := 30*time.Second
	// create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// The cancel should be deferred so that resources are cleaned up
	defer cancel()

	// Create the command with our context
	cmd := exec.CommandContext(ctx, wskdeploy.Path, args...)
	// set the cwd of the command
	cmd.Dir = wskdeploy.Dir

	// print command
	printCommand(cmd)

	// capture standard error in this buffer
	var b bytes.Buffer
	cmd.Stderr = &b

	// create pipe for standard input
	stdin, err := cmd.StdinPipe()
	printError(err)
	if err != nil {
		fmt.Println("==> Error: failed to create pipe for STDIN")
		return false
	}

	// send "y" in response to wskdeploy interacting mode to continue with deployment
	// (TODO) this has to be revised as right now, we are just sending "y" in stdin
	// (TODO) without any check on stdout, check whether wskdeploy is waiting for an
	// (TODO) yes from user or not.
	go func() {
		stdin.Write([]byte("y"))
		defer stdin.Close()
	}()

	out,err := cmd.Output()

	printOutput(out)
	printOutput(b.Bytes())

	// We want to check the context error to see if the timeout was executed.
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("==> Error: wskdeploy timed out before finishing its execution after " + timeout.String() + " seconds")
		printError(ctx.Err())
		return false
	}

	var waitStatus syscall.WaitStatus

	// If there's no context error, we know the command completed (or errored).
	if err != nil {
		printError(err)
		// Did the command fail because of an unsuccessful exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		}
		return false
	}

	// wskdeploy was successful
	waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
	printOutput([]byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))

	return true


}

func (wskdeploy *Wskdeploy) Deploy(manifestPath string, deploymentPath string) (bool) {
	return wskdeploy.RunCommand("-m", manifestPath, "-d", deploymentPath)
}

func (wskdeploy *Wskdeploy) Undeploy(manifestPath string, deploymentPath string) (bool) {
	return wskdeploy.RunCommand("undeploy", "-m", manifestPath, "-d", deploymentPath)
}

func (wskdeploy *Wskdeploy) DeployProjectPathOnly(projectPath string) (bool) {
        return wskdeploy.RunCommand("-p", projectPath)
}

func (wskdeploy *Wskdeploy) DeployManifestPathOnly(manifestpath string) (bool) {
	return wskdeploy.RunCommand("-p", manifestpath)
}
