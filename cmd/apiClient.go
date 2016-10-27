/*
 * Copyright 2015-2016 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"net/http"
	"net/url"
	"fmt"
	"os"
	"bufio"
	"strings"
	wskApi "github.com/openwhisk/go-whisk/whisk"
	"github.com/spf13/cobra"
)

const DefaultAuth string = ""
const DefaultAPIHost string = ""
const DefaultAPIVersion string = "v1"
const DefaultAPIBuild string = ""
const DefaultAPIBuildNo string = ""
const DefaultNamespace string = "_"
const DefaultPropsFile string = "~/.wskprops"

var client *wskApi.Client

func setupClientConfig(cmd *cobra.Command, args []string) error {

  props, err := readProps(DefaultPropsFile)
	check(err)

	if props["AUTH"] == "" {
		props["AUTH"] = DefaultAuth
	}

	if props["NAMESPACE"] == "" {
		props["NAMESPACE"] = DefaultNamespace
	}

	if props["APIHOST"] == "" {
		props["APIHOST"] = DefaultAPIHost
	}

  baseUrl, err := getURLBase(props["APIHOST"])
	check(err)

  fmt.Println("Got props ", props)

	clientConfig := &wskApi.Config{
		AuthToken: props["AUTH"],
		Namespace: props["NAMESPACE"],
		BaseURL:   baseUrl,
		Version:   DefaultAPIVersion,
		Insecure:  false,
	}

	// Setup client
	client, err = wskApi.NewClient(http.DefaultClient, clientConfig)
	check(err)

	return nil
}

func readProps(path string) (map[string]string, error) {

    props := map[string]string{}

    file, err := os.Open(path)
    if err != nil {
        // If file does not exist, just return props
        fmt.Println("Unable to read whisk properties file '%s' (file open error: %s); falling back to default properties\n" ,path, err)
        return props, nil
    }
    defer file.Close()

    lines := []string{}
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }

    props = map[string]string{}
    for _, line := range lines {
        kv := strings.Split(line, "=")
        if len(kv) != 2 {
            // Invalid format; skip
            continue
        }
        props[kv[0]] = kv[1]
    }

    return props, nil

}

func getURLBase(host string) (*url.URL, error)  {
    urlBase := fmt.Sprintf("%s/api/", host)
    url, err := url.Parse(urlBase)

    if len(url.Scheme) == 0 || len(url.Host) == 0 {
        urlBase = fmt.Sprintf("https://%s/api/", host)
        url, err = url.Parse(urlBase)
    }

    return url, err
}
