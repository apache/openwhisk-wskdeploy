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
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/apache/openwhisk-wskdeploy/wskprint"
)

// check if the path represents file path or dir path
func isFilePath(path string) bool {
	// when split returns dir and file, splitting path on the final "/"
	// check if file is not empty to classify that path as a file path
	_, file := filepath.Split(path)
	if len(file) == 0 {
		return false
	}
	return true
}

// check if the given path exists as a file
func isFile(path string) (bool, error) {
	var err error
	var info os.FileInfo
	// run stat on the file and if the it returns no error,
	// read the fileInfo to check if its a file or not
	if info, err = os.Stat(path); err == nil {
		if info.Mode().IsRegular() {
			return true, nil
		}
	}
	// stat returned an error and here we are chekcking if it was os.PathError
	if !os.IsNotExist(err) {
		return false, nil
	}
	// after running through all the possible checks, return false and an err
	return false, err
}

// copy one single source file to the destination path
func copyFile(src, dst string) error {
	var err error
	var sourceFD *os.File
	var destFD *os.File
	var srcInfo os.FileInfo
	var srcDirInfo os.FileInfo

	wskprint.PrintlnOpenWhiskVerbose(Flags.Verbose, "Source File ["+src+"] is being copied to ["+dst+"]")

	// open source file to read it from disk and defer to close the file
	// once this function is done executing
	if sourceFD, err = os.Open(src); err != nil {
		return err
	}
	defer sourceFD.Close()

	// running stat on Dir(src), Dir returns all but the last element of the src
	// this info is needed in case when a destination path has a directory structure which does not exist
	if srcDirInfo, err = os.Stat(filepath.Dir(src)); err != nil {
		return err
	}

	// check if the parent directory exist before creating a destination file
	// create specified path along with creating any parent directory
	// e.g. when destination is greeting/common/utils.js and parent dir common
	// doesn't exist, its getting created here at greeting/common
	if _, err = os.Stat(filepath.Dir(dst)); os.IsNotExist(err) {
		wskprint.PrintlnOpenWhiskVerbose(Flags.Verbose, "Creating directory pattern ["+filepath.Dir(dst)+"] before creating destination file")
		if err = os.MkdirAll(filepath.Dir(dst), srcDirInfo.Mode()); err != nil {
			return err
		}
	}

	// create destination file before copying source content
	// defer closing the destination file until the function is done executing
	if destFD, err = os.Create(dst); err != nil {
		return err
	}
	defer destFD.Close()

	// now, actually copy the source file content into destination file
	if _, err = io.Copy(destFD, sourceFD); err != nil {
		return err
	}

	// retrieve the file mode bits of the source file
	// so that the bits can be set to the destination file
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// recursively copy the entire source directory to destination path
func copyDir(src, dst string) error {
	var err error
	var fileDescriptors []os.FileInfo
	var srcInfo os.FileInfo

	// retrieve os.fileInfo of the source directory
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	// create destination directory with parent directories
	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// now, retrieve all the directory/file entries under the source directory
	if fileDescriptors, err = ioutil.ReadDir(src); err != nil {
		return err
	}

	// iterating over the entire list of files/directories under the destination path
	// run copyFile or recursive copyDir based on if its file or dir
	for _, fd := range fileDescriptors {
		srcFilePath := path.Join(src, fd.Name())
		dstFilePath := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = copyDir(srcFilePath, dstFilePath); err != nil {
				return err
			}
		} else {
			if err = copyFile(srcFilePath, dstFilePath); err != nil {
				return err
			}
		}
	}
	return nil
}
