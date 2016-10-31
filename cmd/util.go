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
	"io/ioutil"
	"net/http"
	"net/url"
)

// ServerlessBinaryCommand is the CLI name to run serverless
const ServerlessBinaryCommand = "serverless"

// ManifestProvider is a provider description in the manifest
type ManifestProvider struct {
	Name    string
	Runtime string
}

// Manifest is the main manifest file
type Manifest struct {
	Service  string
	Provider ManifestProvider
}

// errors
type serverlessErr struct {
	msg string
}

func (e *serverlessErr) Error() string {
	return e.msg
}

// Check is a util function to panic when there is an error.
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

type URLReader interface {
	ReadUrl(url *url.URL) (b []byte, err error)
}

func (urlReader *URLReader) ReadUrl(url *url.URL) (content []byte, err error) {
	resp, err := http.Get(url)
	Check(err)
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	Check(err)
	return b, nil
}

type LocalReader interface {
	ReadLocal(path string) (content []byte, err error)
}

func (localReader *LocalReader) ReadLocal(path string) (content []byte, err error) {
	cont, err := ioutil.ReadFile(path)
	Check(err)
	content = append(content, cont)
	return
}

// agnostic util reader to fetch content from web or local path or potentially other places.
type ContentReader struct {
	URLReader
	LocalReader
}
