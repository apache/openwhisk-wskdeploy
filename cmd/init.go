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

package cmd

import (
	"bufio"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init helps you create a manifest file on OpenWhisk",
	Run: func(cmd *cobra.Command, args []string) {
		maniyaml := parsers.ReadOrCreateManifest()

		reader := bufio.NewReader(os.Stdin)

		maniyaml.Package.Packagename = askName(reader, maniyaml.Package.Packagename)
		maniyaml.Package.Version = askVersion(reader, maniyaml.Package.Version)
		maniyaml.Package.License = askLicense(reader, maniyaml.Package.License)

		parsers.Write(maniyaml, "manifest.yaml")

		// Create directory structure
		os.Mkdir("actions", 0777)
	},
}

func askName(reader *bufio.Reader, def string) string {
	if len(def) == 0 {
		abspath, err := filepath.Abs(utils.Flags.ProjectPath)
		utils.Check(err)
		def = filepath.Base(abspath)
	}
	return utils.Ask(reader, "Name", def)
}

func askVersion(reader *bufio.Reader, def string) string {
	if len(def) == 0 {
		def = "0.0.1"
	}
	return utils.Ask(reader, "Version", def)
}

func askLicense(reader *bufio.Reader, def string) string {
	if len(def) == 0 {
		def = "Apache-2.0"
	}
	return utils.Ask(reader, "License", def)
}

// init initializes this package
func init() {
	RootCmd.AddCommand(initCmd)
}
