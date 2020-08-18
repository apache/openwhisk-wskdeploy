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
	"fmt"
	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/apache/openwhisk-wskdeploy/utils"
	"github.com/apache/openwhisk-wskdeploy/wskderrors"
	"github.com/apache/openwhisk-wskdeploy/wski18n"
	"github.com/apache/openwhisk-wskdeploy/wskprint"
	"strings"
)

//for web action support, code from wsk cli with tiny adjustments
const (
	REQUIRE_WHISK_AUTH = "require-whisk-auth"
	WEB_EXPORT_ANNOT   = "web-export"
	RAW_HTTP_ANNOT     = "raw-http"
	FINAL_ANNOT        = "final"
	TRUE               = "true"
	MAX_JS_INT         = 1<<53 - 1
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

func SetWebActionAnnotations(filePath string, action string, webMode string, annotations whisk.KeyValueArr, fetch bool) (whisk.KeyValueArr, error) {
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
		return webActionAnnotations(fetch, annotations, addWebRawAnnotations)
	default:
		return nil, wskderrors.NewInvalidWebExportError(filePath, action, webMode, getValidWebExports())
	}
}

type WebActionAnnotationMethod func(annotations whisk.KeyValueArr) whisk.KeyValueArr

func webActionAnnotations(fetchAnnotations bool, annotations whisk.KeyValueArr,
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

func addWebRawAnnotations(annotations whisk.KeyValueArr) whisk.KeyValueArr {
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

func HasAnnotation(annotations *whisk.KeyValueArr, key string) bool {
	return (annotations.FindKeyValue(key) >= 0)
}

func warnWebAnnotationMissingFromActionOrSequence(apiName string, actionName string, isSequence bool) {
	nameKey := wski18n.KEY_ACTION
	i18nWarningID := wski18n.ID_WARN_API_MISSING_WEB_ACTION_X_action_X_api_X

	if isSequence {
		nameKey = wski18n.KEY_SEQUENCE
		i18nWarningID = wski18n.ID_WARN_API_MISSING_WEB_SEQUENCE_X_sequence_X_api_X
	}

	warningString := wski18n.T(i18nWarningID,
		map[string]interface{}{
			nameKey:         actionName,
			wski18n.KEY_API: apiName})
	wskprint.PrintOpenWhiskWarning(warningString)
}

func TryUpdateAPIsActionToWebAction(records []utils.ActionRecord, pkgName string, apiName string, actionName string, isSequence bool) error {

	// if records are nil; it may be that the Action already exists at target provider OR
	// this is a unit test.  If the former case, we pass through and allow provider to validate
	// and return an error.
	if records != nil {
		action := utils.GetActionFromActionRecords(records, pkgName, actionName)

		if !HasAnnotation(&action.Annotations, WEB_EXPORT_ANNOT) {
			if !utils.Flags.Strict {
				warnWebAnnotationMissingFromActionOrSequence(apiName, actionName, isSequence)
				action.Annotations = addWebAnnotations(action.Annotations)
				wskprint.PrintOpenWhiskVerbose(utils.Flags.Verbose,
					fmt.Sprintf("Web Annotations to Action; result: %v\n", action.Annotations))
			} else {
				return wskderrors.NewInvalidWebActionError(apiName, actionName, isSequence)
			}
		} else {
			// verify its web-export annotation value is "true", else error
			if !action.WebAction() {
				return wskderrors.NewInvalidWebActionError(apiName, actionName, isSequence)
			}
		}
	}

	return nil
}

func ValidateRequireWhiskAuthAnnotationValue(actionName string, value interface{}) (string, error) {
	var isValid = false
	var enabled = wski18n.FEATURE_DISABLED

	switch value.(type) {
	case string:
		secureValue := value.(string)
		// assure the user-supplied token is valid (i.e., for now a non-empty string)
		if len(secureValue) != 0 && secureValue != "<nil>" {
			isValid = true
			enabled = wski18n.FEATURE_ENABLED
		}
	case int:
		secureValue := value.(int)
		// FYI, the CLI defines MAX_JS_INT = 1<<53 - 1 (i.e.,  9007199254740991)
		// NOTE: For JS, the largest exact integral value is 253-1, or 9007199254740991.
		// In ES6, this is defined as Number MAX_SAFE_INTEGER.
		// However, in JS, the bitwise operators and shift operators operate on 32-bit ints,
		// so in that case, the max safe integer is 231-1, or 2147483647
		// We also disallow negative integers
		// NOTE: when building for 386 archs. we need to assure comparison with MAX_JS_INT does not
		// "blow up" and must allow the compiler to compare an untyped int (secureValue) to effectively
		// an int64... so for the comparison we MUST force a type conversion to avoid "int" size mismatch
		if int64(secureValue) < MAX_JS_INT && secureValue > 0 {
			isValid = true
			enabled = wski18n.FEATURE_ENABLED
		}
	case bool:
		secureValue := value.(bool)
		isValid = true
		if secureValue {
			enabled = wski18n.FEATURE_ENABLED
		}
	}

	if !isValid {
		errMsg := wski18n.T(wski18n.ID_ERR_WEB_ACTION_REQUIRE_AUTH_TOKEN_INVALID_X_action_X_key_X_value,
			map[string]interface{}{
				wski18n.KEY_ACTION: actionName,
				wski18n.KEY_KEY:    REQUIRE_WHISK_AUTH,
				wski18n.KEY_VALUE:  fmt.Sprintf("%v", value)})
		return errMsg, wskderrors.NewActionSecureKeyError(errMsg)
	}

	// Emit an affirmation that security token will be applied to the action
	msg := wski18n.T(wski18n.ID_VERBOSE_ACTION_AUTH_X_action_X_value_X,
		map[string]interface{}{
			wski18n.KEY_ACTION: actionName,
			wski18n.KEY_VALUE:  enabled})
	wskprint.PrintlnOpenWhiskVerbose(utils.Flags.Verbose, msg)

	return msg, nil
}
