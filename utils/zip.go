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
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

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

	var zipFile *os.File
	var err error
	var fileInfo os.FileInfo

	// create zip file e.g. greeting.zip
	if zipFile, err = os.Create(zw.des); err != nil {
		return err
	}
	defer zipFile.Close()

	// creating a new zip writter for greeting.zip
	zw.zipWritter = zip.NewWriter(zipFile)

	// walk file system rooted at the directory specified in "function"
	// walk over each file and dir under root directory e.g. function: actions/greeting
	// add actions/greeting/index.js and actions/greeting/package.json to zip file
	if err = filepath.Walk(zw.src, zw.zipFile); err != nil {
		return nil
	}

	// maintain a list of included files and/or directories with their destination
	var includeInfo []Include

	// iterate over set of included files specified in manifest YAML e.g.
	// include:
	// - ["source"]
	// - ["source", "destination"]
	for _, includeData := range zw.include {
		var i Include
		// if "destination" is not specified, its considered same as "source"
		// "source" is relative to where manifest.yaml file is located
		// relative source path is converted to absolute path by appending manifest path
		// since the relative source path might not be accessible from where wskdeploy is invoked
		// "destination" is relative to the action directory, the one specified in function
		// relative path is converted to absolute path by appending function directory
		if len(includeData) == 1 {
			i.source = filepath.Join(zw.manifestFilePath, includeData[0])
			i.destination = i.source
		} else if len(includeData) == 2 {
			i.source = filepath.Join(zw.manifestFilePath, includeData[0])
			i.destination = filepath.Join(zw.src, includeData[1])
		} else {
			continue
		}

		// append just parsed include info to the list for further processing
		includeInfo = append(includeInfo, i)

		// now determine whether the included item is file or dir
		if fileInfo, err = os.Stat(i.source); err != nil {
			return err
		}

		// if the included item is a directory, call a function to copy the
		// entire directory recursively including its subdirectories and files
		if fileInfo.IsDir() {
			if err = copyDir(i.source, i.destination); err != nil {
				return err
			}
			// if the included item is a file, call a function to copy the file
			// along with its path by creating the parent directories
		} else {
			if err = copyFile(i.source, i.destination); err != nil {
				return err
			}
		}

		// add included item into zip file greeting.zip
		if err = filepath.Walk(i.destination, zw.zipFile); err != nil {
			return nil
		}
	}

	// now close the zip file greeting.zip as all the included items
	// are added into the zip file along with the action root dir
	if err = zw.zipWritter.Close(); err != nil {
		return err
	}

	// and its safe to delete the files/directories which we copied earlier
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
