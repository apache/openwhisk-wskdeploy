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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/apache/openwhisk-wskdeploy/wskderrors"
	"github.com/apache/openwhisk-wskdeploy/wski18n"
)

const (
	FLAG_CONFIG           = "config"
	FLAG_PROJECT          = "project"
	FLAG_PROJECT_SHORT    = "p"
	FLAG_MANIFEST         = "manifest"
	FLAG_MANIFEST_SHORT   = "m"
	FLAG_DEPLOYMENT       = "deployment"
	FLAG_DEPLOYMENT_SHORT = "d"
	FLAG_STRICT           = "strict"
	FLAG_STRICT_SHORT     = "s"
	FLAG_PREVIEW          = "preview"
	FLAG_VERBOSE          = "verbose"
	FLAG_VERBOSE_SHORT    = "v"
	FLAG_API_HOST         = "apihost"
	FLAG_NAMESPACE        = "namespace"
	FLAG_NAMESPACE_SHORT  = "n"
	FLAG_AUTH             = "auth"
	FLAG_AUTH_SHORT       = "u"
	FLAG_APIVERSION       = "apiversion"
	FLAG_KEY              = "key"
	FLAG_KEY_SHORT        = "k"
	FLAG_CERT             = "cert"
	FLAG_CERT_SHORT       = "c"
	FLAG_MANAGED          = "managed"
	FLAG_PROJECTNAME      = "projectname"
	FLAG_TRACE            = "trace"
	FLAG_TRACE_SHORT      = "t"
	FLAG_PARAM            = "param"
	FLAG_PARAMFILE        = "param-file"
	FLAG_PARAMFILE_SHORT  = "P"
	SHORT_CMD             = "-"
	LONG_CMD              = SHORT_CMD + SHORT_CMD
)

// parse key value pairs from --param
// read parameters from param JSON file specified in --param-file
func parseArgsForParams(args []string) ([]string, []string, error) {
	var paramArgs []string
	var err error

	i := 0

	// iterate over the list of input arguments
	// append key value pair of --param <parameter name> <parameter value> as a string
	// in case of --param-file, read parameter file, append each JSON "key: value" pair as a string
	for i < len(args) {
		// when arg is -P or --param-file, open and read the specified JSON file
		if args[i] == SHORT_CMD+FLAG_PARAMFILE_SHORT || args[i] == LONG_CMD+FLAG_PARAMFILE {
			// Command line parser library Cobra, assigns value of --param-file to utils.Flags.ParamFile
			// but at this point of execution, we still don't have utils.Flags.ParamFile assigned any value
			// and that's the reason, we have to explicitly parse the argument list to get the file name
			paramArgs, args, err = getValueFromArgs(args, i, paramArgs)
			if err != nil {
				return nil, nil, err
			}
			filename := paramArgs[len(paramArgs)-1]
			// drop the argument (--param-file) and its value from the argument list after retrieving filename
			// read file content as a single string and append it to the list of params
			file, readErr := ioutil.ReadFile(filename)
			if readErr != nil {
				err = wskderrors.NewCommandError(FLAG_PARAMFILE+"/"+FLAG_PARAMFILE_SHORT,
					wski18n.T(wski18n.ID_ERR_INVALID_PARAM_FILE_X_file_X,
						map[string]interface{}{
							wski18n.KEY_PATH: filename,
							wski18n.KEY_ARG:  FLAG_PARAMFILE + "/" + FLAG_PARAMFILE_SHORT,
							wski18n.KEY_ERR:  readErr}))
				return nil, nil, err
			}
			paramArgs[len(paramArgs)-1] = string(file)
			// --param can appear multiple times in a single invocation of whisk deploy
			// for example, wskdeploy -m manifest.yaml --param key1 value1 --param key2 value2
			// parse key value map for each --param from the argument list
			// drop each --param and key value map from the argument list after reading the map
			// append key value as a string to the list of params
		} else if args[i] == LONG_CMD+FLAG_PARAM {
			paramArgs, args, err = getKeyValueArgs(args, i, paramArgs)
			if err != nil {
				return nil, nil, err
			}
		} else {
			i++
		}
	}
	return args, paramArgs, nil
}

func getValueFromArgs(args []string, argIndex int, parsedArgs []string) ([]string, []string, error) {
	var err error
	// check if parameter argument has any value
	if len(args)-1 >= argIndex+1 {
		// update arguments list and parameter list
		// append param value (JSON file name) to list of parameters
		parsedArgs = append(parsedArgs, args[argIndex+1])
		// drop the parameter (--param-file) and its value (JSON file name) from the argument list
		args = append(args[:argIndex], args[argIndex+2:]...)
	} else {
		err = wskderrors.NewCommandError(args[argIndex],
			wski18n.T(wski18n.ID_ERR_ARG_MISSING_VALUE_X_arg_X,
				map[string]interface{}{wski18n.KEY_ARG: args[argIndex]}))
	}
	return parsedArgs, args, err
}

func getKeyValueArgs(args []string, argIndex int, parsedArgs []string) ([]string, []string, error) {
	var err error
	var key string
	var value string
	// check if parameter argument has a key/value pair specified
	if len(args)-1 >= argIndex+2 {
		key = args[argIndex+1]
		value = args[argIndex+2]
		// append key/value pairs to list of parameters
		// drop those pairs from the argument list
		parsedArgs = append(parsedArgs, getFormattedJSON(key, value))
		args = append(args[:argIndex], args[argIndex+3:]...)
	} else {
		err = wskderrors.NewCommandError(args[argIndex],
			wski18n.T(wski18n.ID_ERR_ARG_MISSING_KEY_VALUE_X_arg_X,
				map[string]interface{}{wski18n.KEY_ARG: args[argIndex]}))
	}
	return parsedArgs, args, err
}

func getFormattedJSON(key string, value string) string {
	var res string
	key = getEscapedJSON(key)
	if isValidJSON(value) {
		res = fmt.Sprintf("{\"%s\": %s}", key, value)
	} else {
		res = fmt.Sprintf("{\"%s\": \"%s\"}", key, value)
	}
	return res
}

func isValidJSON(value string) bool {
	var jsonInterface interface{}
	err := json.Unmarshal([]byte(value), &jsonInterface)
	return err == nil
}

func getEscapedJSON(value string) string {
	value = strings.Replace(value, "\\", "\\\\", -1)
	value = strings.Replace(value, "\"", "\\\"", -1)
	return value
}
