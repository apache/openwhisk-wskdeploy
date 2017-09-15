// +build unit

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
    "testing"
    "github.com/stretchr/testify/assert"
    "os"
    "io/ioutil"
    "bufio"
    "fmt"
)

func CreateTmpfile(filename string) (f *os.File, err error) {
    dir, _ := os.Getwd()
    tmpfile, err := ioutil.TempFile(dir, filename)
    if err != nil {
        return nil, err
    }
    return tmpfile, nil
}

func CreateTmpDir(dir string) error {
    return os.Mkdir(dir, 0777)
}

func CreateFile(lines []string, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    for _, line := range lines {
        fmt.Fprintln(w, line)
    }
    return w.Flush()
}

func TestFileExists(t *testing.T) {
    testfile := "testFile"
    tmpFile, err := CreateTmpfile(testfile)
    assert.Nil(t, err, "Failed to create the temporary file.")

    // The file exists.
    fileName := tmpFile.Name()
    result := FileExists(fileName)
    assert.True(t, result)

    // After the file is removed, the file does not exist.
    tmpFile.Close()
    os.Remove(tmpFile.Name())
    result = FileExists(fileName)
    assert.False(t, result)
}

func TestIsDirectory(t *testing.T) {
    testDir := "testFile"
    err := CreateTmpDir(testDir)
    assert.Nil(t, err, "Failed to create the temporary file.")

    // The directory exists.
    result := IsDirectory(testDir)
    assert.True(t, result)

    // After the file is removed, the file does not exist.
    os.Remove(testDir)
    result = IsDirectory(testDir)
    assert.False(t, result)

    // If this function is tested against a valid file, we will get false.
    testfile := "testFile"
    tmpFile, err := CreateTmpfile(testfile)
    assert.Nil(t, err, "Failed to create the temporary file.")
    fileName := tmpFile.Name()
    result = IsDirectory(fileName)
    assert.False(t, result)
    tmpFile.Close()
    os.Remove(tmpFile.Name())

    // If this function is tested against an invalid path, we will get false.
    result = IsDirectory("NonExistent")
    assert.False(t, result)
}

func TestReadProps(t *testing.T) {
    testfile := "testWskPropsRead"
    testKey := "testKey"
    testValue := "testValue"
    testKeySec := "testKeySec"
    testValueSec := "testValueSec"
    lines := []string{ testKey + "=" + testValue, testKeySec + "=" + testValueSec }
    CreateFile(lines, testfile)
    props, err := ReadProps(testfile)
    assert.Nil(t, err, "Failed to read the test prop file.")
    assert.Equal(t, testValue, props[testKey])
    assert.Equal(t, testValueSec, props[testKeySec])
    err = os.Remove(testfile)
    assert.Nil(t, err, "Failed to delete the test prop file.")

    // Failed to read wskprops file if it does not exist.
    props, err = ReadProps(testfile)
    assert.NotNil(t, err)
    assert.Equal(t, 0, len(props))
}

func TestWriteProps(t *testing.T) {
    testfile := "testWskPropsWrite"
    testKey := "testKeyWrite"
    testValue := "testValueWrite"
    testKeySec := "testKeyWriteSec"
    testValueSec := "testValueWriteSec"
    props := map[string]string{
        testKey: testValue,
        testKeySec:   testValueSec,
    }
    err := WriteProps(testfile, props)
    assert.Nil(t, err, "Failed to write the test prop file.")

    propsResult, error := ReadProps(testfile)
    assert.Nil(t, error, "Failed to read the test prop file.")
    assert.Equal(t, testValue, propsResult[testKey])
    assert.Equal(t, testValueSec, propsResult[testKeySec])
    err = os.Remove(testfile)
    assert.Nil(t, err, "Failed to delete the test prop file.")
}
