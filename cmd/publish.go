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
	"fmt"

	"bufio"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	SuggestFor: []string {"publicize"},
	Short: "Publish a package to a registry",
	Long:  `Publish a package to the registry set in ~/.wskprops`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get registry location

		userHome := utils.GetHomeDirectory()
		propPath := path.Join(userHome, ".wskprops")

		configs, err := utils.ReadProps(propPath)
		if err != nil {
            return err
        }

		registry, ok := configs["REGISTRY"]
		if !ok {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Registry URL not found in ~./wskprops. Must be set before publishing.\n")
			for {
				registry = utils.Ask(reader, "Registry URL", "")

				_, err := url.Parse(registry)
				if err == nil {
					// TODO: send request to registry to check it exists.
					break
				}
				fmt.Print("Malformed repository URL. Try again")

			}
			configs["REGISTRY"] = registry
			utils.WriteProps(propPath, configs)
		}

		// Get repo URL
		maniyaml, err := parsers.ReadOrCreateManifest()
        if err != nil {
            return err
        }

		if len(maniyaml.Package.Repositories) > 0 {
			repoURL := maniyaml.Package.Repositories[0].Url

			paths := strings.Split(repoURL, "/")
			l := len(paths)
			if l < 2 {
				fmt.Print("Fatal error: malformed repository URL in manifest file :" + repoURL)
				return nil
			}

			repo := paths[l-1]
			owner := paths[l-2]

			// Send HTTP request
			client := &http.Client{}
			request, err := http.NewRequest("PUT", registry+"?owner="+owner+"&repo="+repo, nil)
			if err != nil {
                return err
            }
			_, err = client.Do(request)
            if err != nil {
                return err
            }

		} else {
			fmt.Print("Fatal error: missing repository URL in manifest file.")
		}
        return nil
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)

}
