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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func isFilePath(path string) bool {
	_, file := filepath.Split(path)
	if len(file) == 0 {
		return false
	}
	return true
}

func isFile(path string) (bool, error) {
	var err error
	var info os.FileInfo
	if info, err = os.Stat(path); err == nil {
		if info.Mode().IsRegular() {
			return true, nil
		}
	}
	if !os.IsNotExist(err) {
		return false, err
	}
	return false, nil
}

func copyFile(src, dst string) error {
	var err error
	var sourceFD *os.File
	var destFD *os.File
	var srcInfo os.FileInfo
	var srcDirInfo os.FileInfo

	if sourceFD, err = os.Open(src); err != nil {
		return err
	}
	defer sourceFD.Close()

	if srcDirInfo, err = os.Stat(filepath.Dir(src)); err != nil {
		return err
	}

	if _, err = os.Stat(filepath.Dir(dst)); os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(dst), srcDirInfo.Mode()); err != nil {
			return err
		}
	}

	if destFD, err = os.Create(dst); err != nil {
		return err
	}
	defer destFD.Close()

	if _, err = io.Copy(destFD, sourceFD); err != nil {
		return err
	}

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

func copyDir(src string, dst string) error {
	var err error
	var fileDescriptors []os.FileInfo
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fileDescriptors, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fileDescriptors {
		srcFilePath := path.Join(src, fd.Name())
		dstFilePath := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDir(srcFilePath, dstFilePath); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = copyFile(srcFilePath, dstFilePath); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
