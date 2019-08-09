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

package webaction

import (
	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/apache/openwhisk-wskdeploy/wskderrors"
	"strings"
)

//for web action support, code from wsk cli with tiny adjustments
const (
	WEB_EXPORT_ANNOT = "web-export"
	RAW_HTTP_ANNOT   = "raw-http"
	FINAL_ANNOT      = "final"
	TRUE             = "true"
)

var webExport map[string]string = map[string]string{
	"TRUE":  "true",
	"FALSE": "false",
	"NO":    "no",
	"YES":   "yes",
	"RAW":   "raw",
}

func deleteKey(key string, keyValueArr whisk.KeyValueArr) whisk.KeyValueArr {
	for i := 0; i < len(keyValueArr); i++ {
		if keyValueArr[i].Key == key {
			keyValueArr = append(keyValueArr[:i], keyValueArr[i+1:]...)
			break
		}
	}

	return keyValueArr
}

func addKeyValue(key string, value interface{}, keyValueArr whisk.KeyValueArr) whisk.KeyValueArr {
	keyValue := whisk.KeyValue{
		Key:   key,
		Value: value,
	}

	return append(keyValueArr, keyValue)
}

func WebAction(filePath string, action string, webMode string, annotations whisk.KeyValueArr, fetch bool) (whisk.KeyValueArr, error) {
	switch strings.ToLower(webMode) {
	case webExport["TRUE"]:
		fallthrough
	case webExport["YES"]:
		return webActionAnnotations(fetch, annotations, addWebAnnotations)
	case webExport["NO"]:
		fallthrough
	case webExport["FALSE"]:
		return webActionAnnotations(fetch, annotations, deleteWebAnnotations)
	case webExport["RAW"]:
		return webActionAnnotations(fetch, annotations, addRawAnnotations)
	default:
		return nil, wskderrors.NewInvalidWebExportError(filePath, action, webMode, getValidWebExports())
	}
}

type WebActionAnnotationMethod func(annotations whisk.KeyValueArr) whisk.KeyValueArr

func webActionAnnotations(
	fetchAnnotations bool,
	annotations whisk.KeyValueArr,
	webActionAnnotationMethod WebActionAnnotationMethod) (whisk.KeyValueArr, error) {
	if annotations != nil || !fetchAnnotations {
		annotations = webActionAnnotationMethod(annotations)
	}

	return annotations, nil
}

func addWebAnnotations(annotations whisk.KeyValueArr) whisk.KeyValueArr {
	annotations = deleteWebAnnotationKeys(annotations)
	annotations = addKeyValue(WEB_EXPORT_ANNOT, true, annotations)
	annotations = addKeyValue(RAW_HTTP_ANNOT, false, annotations)
	annotations = addKeyValue(FINAL_ANNOT, true, annotations)

	return annotations
}

func deleteWebAnnotations(annotations whisk.KeyValueArr) whisk.KeyValueArr {
	annotations = deleteWebAnnotationKeys(annotations)
	annotations = addKeyValue(WEB_EXPORT_ANNOT, false, annotations)
	annotations = addKeyValue(RAW_HTTP_ANNOT, false, annotations)
	annotations = addKeyValue(FINAL_ANNOT, false, annotations)

	return annotations
}

func addRawAnnotations(annotations whisk.KeyValueArr) whisk.KeyValueArr {
	annotations = deleteWebAnnotationKeys(annotations)
	annotations = addKeyValue(WEB_EXPORT_ANNOT, true, annotations)
	annotations = addKeyValue(RAW_HTTP_ANNOT, true, annotations)
	annotations = addKeyValue(FINAL_ANNOT, true, annotations)

	return annotations
}

func deleteWebAnnotationKeys(annotations whisk.KeyValueArr) whisk.KeyValueArr {
	annotations = deleteKey(WEB_EXPORT_ANNOT, annotations)
	annotations = deleteKey(RAW_HTTP_ANNOT, annotations)
	annotations = deleteKey(FINAL_ANNOT, annotations)

	return annotations
}

func getValidWebExports() []string {
	var validWebExports []string
	for _, v := range webExport {
		validWebExports = append(validWebExports, v)
	}
	return validWebExports
}

func IsWebAction(webexport string) bool {
	webexport = strings.ToLower(webexport)
	if len(webexport) != 0 {
		if webexport == webExport["TRUE"] || webexport == webExport["YES"] || webexport == webExport["RAW"] {
			return true
		}
	}
	return false
}

func IsWebSequence(webexport string) bool {
	return IsWebAction(webexport)
}
