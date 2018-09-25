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

package utils

import (
	"archive/zip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"fmt"
	"github.com/apache/incubator-openwhisk-client-go/whisk"
	"github.com/davecgh/go-spew/spew"
	"github.com/hokaccha/go-prettyjson"
)

const (
	DEFAULT_HTTP_TIMEOUT = 60
	DEFAULT_PROJECT_PATH = "."
	HTTP_FILE_EXTENSION  = "http"
	// name of manifest and deployment files
	ManifestFileNameYaml   = "manifest.yaml"
	ManifestFileNameYml    = "manifest.yml"
	DeploymentFileNameYaml = "deployment.yaml"
	DeploymentFileNameYml  = "deployment.yml"
)

// ActionRecord is a container to keep track of
// a whisk action struct and a location filepath we use to
// map files and manifest declared actions
type ActionRecord struct {
	Action      *whisk.Action
	Packagename string
	Filepath    string
}

type TriggerRecord struct {
	Trigger     *whisk.Trigger
	Packagename string
}

type RuleRecord struct {
	Rule        *whisk.Rule
	Packagename string
}

func GetHomeDirectory() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}

	return usr.HomeDir
}

// Potentially complex structures(such as DeploymentProject, DeploymentPackage)
// could implement those interface which is convenient for put, get subtract in
// containers etc.
type Comparable interface {
	HashCode() uint32
	Equals() bool
}

func IsFeedAction(trigger *whisk.Trigger) (string, bool) {
	for _, annotation := range trigger.Annotations {
		if annotation.Key == "feed" {
			return annotation.Value.(string), true
		}
	}

	return "", false
}

func PrettyJSON(j interface{}) (string, error) {
	formatter := prettyjson.NewFormatter()
	bytes, err := formatter.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func NewZipWritter(src string, des string, include [][]string, manifestFilePath string) *ZipWritter {
	zw := &ZipWritter{src: src, des: des, include: include, manifestFilePath: manifestFilePath}
	return zw
}

type ZipWritter struct {
	src              string
	des              string
	include          [][]string
	manifestFilePath string
	zipWritter       *zip.Writer
}

type Include struct {
	source      string
	destination string
}

func (zw *ZipWritter) zipFile(path string, f os.FileInfo, err error) error {
	spew.Println("++++++++ path ++++++++")
	spew.Dump(path)
	spew.Println("+++++++++ src +++++++")
	spew.Dump(zw.src)
	if err != nil {
		return err
	}
	if !f.Mode().IsRegular() || f.Size() == 0 {
		return nil
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	fileName := strings.TrimPrefix(path, zw.src+"/")
	spew.Println("++++++++file name is+++++++")
	spew.Dump(fileName)
	wr, err := zw.zipWritter.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(wr, file)
	if err != nil {
		return err
	}
	return nil
}

func (zw *ZipWritter) Zip() error {
	// create zip file
	zipFile, err := os.Create(zw.des)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zw.zipWritter = zip.NewWriter(zipFile)
	err = filepath.Walk(zw.src, zw.zipFile)
	if err != nil {
		return nil
	}

	spew.Println("***********src*********")
	spew.Dump(zw.src)

	spew.Println("***********dest********")
	spew.Dump(zw.des)

	spew.Println("***********include******")
	spew.Dump(zw.include)

	var includeInfo []Include

	for _, includeData := range zw.include {

		var i Include

		if len(includeData) == 1 {
			i.source = filepath.Join(zw.manifestFilePath, includeData[0])
			i.destination = i.source
		} else if len(includeData) == 2 {
			i.source = filepath.Join(zw.manifestFilePath, includeData[0])
			i.destination = filepath.Join(zw.src, includeData[1])
		}
		includeInfo = append(includeInfo, i)

		s, err := os.Stat(i.source)
		if err != nil {
			return err
		}

		if s.IsDir() {
			spew.Println("is a directory")
			err = copyDir(i.source, i.destination)
			if err != nil {
				return err
			}
		} else {
			spew.Println("is not a directory")
			err = copyFile(i.source, i.destination)
			if err != nil {
				return err
			}
		}

		err = filepath.Walk(i.destination, zw.zipFile)
		if err != nil {
			return nil
		}
	}

	err = zw.zipWritter.Close()
	if err != nil {
		return err
	}

	for _, i := range includeInfo {
		spew.Println("Deleting destination")
		spew.Dump(i.destination)
		os.RemoveAll(i.destination)
	}

	return nil
}

func copyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo
	var srcDirInfo os.FileInfo

	spew.Println("############src##############")
	spew.Dump(src)
	spew.Println("###############dst#############")
	spew.Dump(dst)

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	spew.Println("Done opening src")

	spew.Dump(filepath.Dir(dst))

	if srcDirInfo, err = os.Stat(filepath.Dir(src)); err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Dir(dst)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(dst), srcDirInfo.Mode())
		spew.Println("Done creating a dir")
		if err != nil {
			spew.Println("Failed to create a dir")
			return err
		}
	}

	if dstfd, err = os.Create(dst); err != nil {
		spew.Println("Failed to create dst")
		spew.Dump(err)
		return err
	}
	defer dstfd.Close()

	spew.Println("done creating dst")

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}

	spew.Println("done copying src to dst")

	spew.Println("src info")
	spew.Dump(srcinfo)
	return os.Chmod(dst, srcinfo.Mode())
}

func copyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = copyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func GetManifestFilePath(projectPath string) string {
	if _, err := os.Stat(path.Join(projectPath, ManifestFileNameYaml)); err == nil {
		return path.Join(projectPath, ManifestFileNameYaml)
	} else if _, err := os.Stat(path.Join(projectPath, ManifestFileNameYml)); err == nil {
		return path.Join(projectPath, ManifestFileNameYml)
	} else {
		return ""
	}
}

func GetDeploymentFilePath(projectPath string) string {
	if _, err := os.Stat(path.Join(projectPath, DeploymentFileNameYaml)); err == nil {
		return path.Join(projectPath, DeploymentFileNameYaml)
	} else if _, err := os.Stat(path.Join(projectPath, DeploymentFileNameYml)); err == nil {
		return path.Join(projectPath, DeploymentFileNameYml)
	} else {
		return ""
	}
}

// agnostic util reader to fetch content from web or local path or potentially other places.
type ContentReader struct {
	URLReader
	LocalReader
}

type URLReader struct {
}

func (urlReader *URLReader) ReadUrl(url string) (content []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return content, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return content, err
	} else {
		defer resp.Body.Close()
	}
	return b, nil
}

type LocalReader struct {
}

func (localReader *LocalReader) ReadLocal(path string) ([]byte, error) {
	cont, err := ioutil.ReadFile(path)
	return cont, err
}

func Read(url string) ([]byte, error) {
	if strings.HasPrefix(url, HTTP_FILE_EXTENSION) {
		return new(ContentReader).URLReader.ReadUrl(url)
	} else {
		return new(ContentReader).LocalReader.ReadLocal(url)
	}
}

func GetJSONFromStrings(content []string, keyValueFormat bool) (interface{}, error) {
	var data map[string]interface{}
	var res interface{}

	for i := 0; i < len(content); i++ {
		dc := json.NewDecoder(strings.NewReader(content[i]))
		dc.UseNumber()
		if err := dc.Decode(&data); err != nil {
			return whisk.KeyValueArr{}, err
		}
	}

	if keyValueFormat {
		res = getKeyValueFormattedJSON(data)
	} else {
		res = data
	}

	return res, nil
}

func getKeyValueFormattedJSON(data map[string]interface{}) whisk.KeyValueArr {
	var keyValueArr whisk.KeyValueArr
	for key, value := range data {
		keyValue := whisk.KeyValue{
			Key:   key,
			Value: value,
		}
		keyValueArr = append(keyValueArr, keyValue)
	}
	return keyValueArr
}
