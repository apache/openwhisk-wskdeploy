// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/openwhisk/openwhisk-wskdeploy/parsers"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init helps you create a manifest file on OpenWhisk",
	Run: func(cmd *cobra.Command, args []string) {
		maniyaml := readOrCreateManifest()

		reader := bufio.NewReader(os.Stdin)

		maniyaml.Package.Packagename = askName(reader, maniyaml.Package.Packagename)
		maniyaml.Package.Version = askVersion(reader, maniyaml.Package.Version)
		maniyaml.Package.License = askLicense(reader, maniyaml.Package.License)

		yamlparser := parsers.NewYAMLParser()
		output, err := yamlparser.Marshal(maniyaml)
		utils.Check(err)

		f, err := os.Create("manifest.yaml")
		utils.Check(err)
		defer f.Close()

		f.Write(output)

		// Create directory structure
		os.Mkdir("actions", 0777)
		os.Mkdir("feeds", 0777)
	},
}

// Read existing manifest file or create new if none exists
func readOrCreateManifest() *parsers.ManifestYAML {
	maniyaml := parsers.ManifestYAML{}

	if _, err := os.Stat("manifest.yaml"); err == nil {
		dat, _ := ioutil.ReadFile("manifest.yaml")
		err := parsers.NewYAMLParser().Unmarshal(dat, &maniyaml)
		utils.Check(err)
	}
	return &maniyaml
}

func ask(reader *bufio.Reader, question string, def string) string {
	fmt.Print(question + " (" + def + "): ")
	answer, _ := reader.ReadString('\n')
	len := len(answer)
	if len == 1 {
		return def
	}
	return answer[:len-1]
}

func askName(reader *bufio.Reader, def string) string {
	if len(def) == 0 {
		abspath, err := filepath.Abs(projectPath)
		utils.Check(err)
		def = filepath.Base(abspath)
	}
	return ask(reader, "Name", def)
}

func askVersion(reader *bufio.Reader, def string) string {
	if len(def) == 0 {
		def = "0.0.1"
	}
	return ask(reader, "Version", def)
}

func askLicense(reader *bufio.Reader, def string) string {
	if len(def) == 0 {
		def = "Apache-2.0"
	}
	return ask(reader, "License", def)
}

// init initializes this package
func init() {
	RootCmd.AddCommand(initCmd)
}
