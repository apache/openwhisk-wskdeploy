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
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"github.com/spf13/cobra"
	"github.com/apache/incubator-openwhisk-wskdeploy/parsers"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/apache/incubator-openwhisk-wskdeploy/wski18n"
	"github.com/apache/incubator-openwhisk-wskdeploy/wskprint"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:		"publish",
	SuggestFor:	[]string{"publicize"},
	Short:		wski18n.T(wski18n.ID_CMD_DESC_SHORT_PUBLISH),
	Long:		wski18n.T(wski18n.ID_CMD_DESC_LONG_PUBLISH),
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
			wskprint.PrintOpenWhiskError(
				wski18n.T(wski18n.ID_ERR_INVALID_URL_X_urltype_X_url_X_filetype_X,
					map[string]interface{}{
						wski18n.KEY_URL_TYPE: parsers.REGISTRY,
						wski18n.KEY_URL: "",
						wski18n.KEY_FILE_TYPE: parsers.WHISK_PROPS}))

			// TODO() should only read if interactive mode is on
			reader := bufio.NewReader(os.Stdin)
			for {
				registry = utils.Ask(reader, parsers.REGISTRY_URL, "")

				_, err := url.Parse(registry)
				if err == nil {
					// TODO() send request to registry to check if it exists.
					break
				}

				// Tell user the URL they entered was invalid, try again...
				wskprint.PrintOpenWhiskError(
					wski18n.T(wski18n.ID_ERR_MALFORMED_URL_X_urltype_X_url_X,
						map[string]interface{}{
							wski18n.KEY_URL_TYPE: parsers.REGISTRY,
							wski18n.KEY_URL: registry}))
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
				wskprint.PrintOpenWhiskError(
					wski18n.T(wski18n.ID_ERR_INVALID_URL_X_urltype_X_url_X_filetype_X,
						map[string]interface{}{
							wski18n.KEY_URL_TYPE: parsers.REPOSITORY,
							wski18n.KEY_URL: repoURL,
							wski18n.KEY_FILE_TYPE: parsers.MANIFEST}))
				return nil
			}

			repo := paths[l - 1]
			owner := paths[l - 2]

			// Send HTTP request
			client := &http.Client{}
			request, err := http.NewRequest("PUT", registry + "?owner=" + owner + "&repo=" + repo, nil)
			if err != nil {
				return err
			}
			_, err = client.Do(request)
			if err != nil {
				return err
			}

		} else {
			wskprint.PrintOpenWhiskError(
				wski18n.T(wski18n.ID_ERR_INVALID_URL_X_urltype_X_url_X_filetype_X,
					map[string]interface{}{
						wski18n.KEY_URL_TYPE: parsers.REPOSITORY,
						wski18n.KEY_URL: "",
						wski18n.KEY_FILE_TYPE: parsers.MANIFEST}))
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)

}
