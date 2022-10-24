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
package parsers

import (
	"encoding/json"
	"fmt"
	"github.com/apache/openwhisk-wskdeploy/utils"
	"github.com/apache/openwhisk-wskdeploy/wskderrors"
	"github.com/apache/openwhisk-wskdeploy/wskenv"
	reflect "github.com/goccy/go-reflect"
)

// TODO(): Support other valid Package Manifest types
// TODO(): i.e., timestamp, version, string256, string64, string16
// TODO(): Support JSON schema validation for type: json
// TODO(): Support OpenAPI schema validation
const (
	STRING  string = "string"
	INTEGER string = "integer"
	FLOAT   string = "float"
	BOOLEAN string = "boolean"
	JSON    string = "json"
	SLICE   string = "slice"
)

var validParameterNameMap = map[string]string{
	STRING:    STRING,
	FLOAT:     FLOAT,
	BOOLEAN:   BOOLEAN,
	INTEGER:   INTEGER,
	"int":     INTEGER,
	"bool":    BOOLEAN,
	"int8":    INTEGER,
	"int16":   INTEGER,
	"int32":   INTEGER,
	"int64":   INTEGER,
	"float32": FLOAT,
	"float64": FLOAT,
	JSON:      JSON,
	"map":     JSON,
	"slice":   SLICE,
}

var typeDefaultValueMap = map[string]interface{}{
	STRING:  "",
	INTEGER: 0,
	FLOAT:   0.0,
	BOOLEAN: false,
	JSON:    make(map[string]interface{}),
	// TODO() Support these types + their validation
	// timestamp
	// null
	// version
	// string256
	// string64
	// string16
	// scalar-unit
	// schema
	// object
}

func isValidParameterType(typeName string) bool {
	_, isValid := typeDefaultValueMap[typeName]
	return isValid
}

// TODO(): throw errors
func getTypeDefaultValue(typeName string) interface{} {

	if val, ok := typeDefaultValueMap[typeName]; ok {
		return val
	} else {
		// TODO() throw an error "type not found" InvalidParameterType
	}
	return nil
}

func IsTypeDefaultValue(typeName string, value interface{}) bool {
	defaultValue := getTypeDefaultValue(typeName)
	if defaultValue == nil {
		return false
	} else if defaultValue == value {
		return true
	}
	return false
}

/*
   ResolveParamTypeFromValue Resolves the Parameter's data type from its actual value.

   Inputs:
   - paramName: name of the parameter for error reporting
   - filepath: the path, including name, of the YAML file which contained the parameter for error reporting
   - value: the parameter value to resolve

   Returns:
   - (string) parameter type name as a string
*/
func ResolveParamTypeFromValue(paramName string, value interface{}, filePath string) (string, error) {
	// Note: 'string' is the default type if not specified and not resolvable.
	var paramType string = "string"
	var err error = nil

	if value != nil {
		actualType := reflect.TypeOf(value).Kind().String()

		// See if the actual type of the value is valid
		if normalizedTypeName, found := validParameterNameMap[actualType]; found {
			// use the full spec. name
			paramType = normalizedTypeName

		} else {
			// raise an error if parameter's value is not a known type
			err = wskderrors.NewInvalidParameterTypeError(filePath, paramName, actualType)
		}
	}
	return paramType, err
}

/*
   resolveSingleLineParameter assures that a Parameter's Type is correctly identified and set from its Value.

   Additionally, this function:

   - detects if the parameter value contains the name of a valid OpenWhisk parameter types. if so, the
     - param.Type is set to detected OpenWhisk parameter type.
     - param.Value is set to the zero (default) value for that OpenWhisk parameter type.

   Inputs:
   - filePath: the path, including name, of the YAML file which contained the parameter for error reporting
   - paramName: name of the parameter for error reporting
   - param: pointer to Parameter structure being resolved

   Returns:
   - (interface{}) the parameter's resolved value
*/
func resolveSingleLineParameter(filePath string, paramName string, param *Parameter) (interface{}, error) {
	var errorParser error

	if !param.multiline {
		// We need to identify parameter Type here for later validation
		param.Type, errorParser = ResolveParamTypeFromValue(paramName, param.Value, filePath)

		// In single-line format, the param's <value> can be a "Type name" and NOT an actual value.
		// if this is the case, we must detect it and set the value to the default for that type name.
		if param.Value != nil && param.Type == "string" {
			// The value is a <string>; now we must test if is the name of a known Type
			if isValidParameterType(param.Value.(string)) {
				// If the value is indeed the name of a Type, we must change BOTH its
				// Type to be that type and its value to that Type's default value
				param.Type = param.Value.(string)
				param.Value = getTypeDefaultValue(param.Type)
				//fmt.Printf("EXIT: Parameter [%s] type=[%v] value=[%v]\n", paramName, param.Type, param.Value)
			}
		}

	} else {
		// TODO() - move string to i18n
		return param.Value, wskderrors.NewYAMLParserErr(filePath,
			"Parameter ["+paramName+"] is not single-line format.")
	}

	return param.Value, errorParser
}

/*
   resolveMultiLineParameter assures that the values for Parameter Type and Value are properly set and are valid.

   Additionally, this function:
   - uses param.Default as param.Value if param.Value is not provided
   - uses the actual param.Value data type for param.type if param.Type is not provided

   Inputs:
   - filepath: the path, including name, of the YAML file which contained the parameter for error reporting
   - paramName: name of the parameter for error reporting
   - param: pointer to Parameter structure being resolved

   Returns:
   - (interface{}) the parameter's resolved value

*/
func resolveMultiLineParameter(filePath string, paramName string, param *Parameter) (interface{}, error) {
	var errorParser error

	if param.multiline {
		var valueType string

		// if we do not have a value, but have a default, use it for the value
		if param.Value == nil && param.Default != nil {
			param.Value = param.Default
		}

		// Note: if either the value or default is in conflict with the type then this is an error
		valueType, errorParser = ResolveParamTypeFromValue(paramName, param.Value, filePath)

		// if we have a declared parameter Type, assure that it is a known value
		if param.Type != "" {
			if !isValidParameterType(param.Type) {
				// TODO() - move string to i18n
				return param.Value, wskderrors.NewYAMLParserErr(filePath,
					"Parameter ["+paramName+"] has an invalid Type. ["+param.Type+"]")
			}
		} else {
			// if we do not have a value for the Parameter Type, use the Parameter Value's Type
			param.Type = valueType
		}

		// TODO{} if the declared and actual parameter type conflict, generate TypeMismatch error
		//if param.Type != valueType{
		//	errorParser = utils.NewParameterTypeMismatchError("", param.Type, valueType )
		//}
	} else {
		// TODO() - move string to i18n
		return param.Value, wskderrors.NewYAMLParserErr(filePath,
			"Parameter ["+paramName+"] is not multiline format.")
	}

	return param.Value, errorParser
}
func interpolateJSON(data map[string]interface{}) map[string]interface{} {
	for key, value := range data {
		if reflect.TypeOf(value).Kind() == reflect.String {
			data[key] = wskenv.InterpolateStringWithEnvVar(value)
		} else if reflect.TypeOf(value).Kind() == reflect.Map {
			data[key] = interpolateJSON(value.(map[string]interface{}))
		}
	}
	return data
}

/*
   resolveJSONParameter assure JSON data is converted to a map[string]{interface*} type.

   This function handles the forms JSON data appears in:
   1) a string containing JSON, which needs to be parsed into map[string]interface{}
   2) is a map of JSON (but not a map[string]interface{}

   Inputs:
   - paramName: name of the parameter for error reporting
   - filePath: the path, including name, of the YAML file which contained the parameter for error reporting
   - param: pointer to Parameter structure being resolved
   - value: the current actual value of the parameter being resolved

   Returns:
   - (interface{}) the parameter's resolved value
*/
func resolveJSONParameter(filePath string, paramName string, param *Parameter, value interface{}) (interface{}, error) {
	var errorParser error

	// TODO() Is the "value" function parameter really needed with the current logic (use param.Value)?
	if param.Type == "json" {
		// Case 1: if user set parameter type to 'json' and the value's type is a 'string'
		if str, ok := value.(string); ok {
			var parsed interface{}
			errParser := json.Unmarshal([]byte(str), &parsed)
			if errParser == nil {
				//fmt.Printf("EXIT: Parameter [%s] type=[%v] value=[%v]\n", paramName, param.Type, parsed)
				return parsed, errParser
			}
		}

		// Case 2: value contains a map of JSON
		// We must make sure the map type is map[string]interface{}; otherwise we cannot
		// marshall it later on to serialize in the body of an HTTP request.
		if param.Value != nil && reflect.TypeOf(param.Value).Kind() == reflect.Map {
			if _, ok := param.Value.(map[interface{}]interface{}); ok {
				var temp map[string]interface{} = utils.ConvertInterfaceMap(param.Value.(map[interface{}]interface{}))
				temp = interpolateJSON(temp)
				//fmt.Printf("EXIT: Parameter [%s] type=[%v] value=[%v]\n", paramName, param.Type, temp)
				return temp, errorParser
			}
		} else {
			errorParser = wskderrors.NewParameterTypeMismatchError(filePath, paramName, JSON, param.Type)
		}

	} else {
		// TODO() - move string to i18n
		errorParser = wskderrors.NewYAMLParserErr(filePath, "Parameter ["+paramName+"] is not JSON format.")
	}

	return param.Value, errorParser
}

/*
   ResolveParameter assures that the Parameter structure's values are correctly filled out for
   further processing.  This includes special processing for

   - single-line format parameters
     - deriving missing param.Type from param.Value
     - resolving case where param.Value contains a valid Parameter type name
   - multi-line format parameters:
     - assures that param.Value is set while taking into account param.Default
     - validating param.Type

   Note: parameter values may set later (overridden) by an (optional) Deployment file

   Inputs:
   - paramName: name of the parameter for error reporting
   - filepath: the path, including name, of the YAML file which contained the parameter for error reporting
   - param: pointer to Parameter structure being resolved

   Returns:
   - (interface{}) the parameter's resolved value
*/
func ResolveParameter(paramName string, param *Parameter, filePath string) (interface{}, error) {

	var errorParser error
	// default resolved parameter value to empty string
	var value interface{} = ""

	// Trace Parameter struct before any resolution
	//dumpParameter(paramName, param, "BEFORE")

	// Parameters can be single OR multi-line declarations which must be processed/validated differently
	// Regardless, the following functions will assure that param.Value and param.Type are correctly set
	if !param.multiline {
		value, errorParser = resolveSingleLineParameter(filePath, paramName, param)

	} else {
		value, errorParser = resolveMultiLineParameter(filePath, paramName, param)
	}

	// String value pre-processing (interpolation)
	// See if we have any Environment Variable replacement within the parameter's value

	// Make sure the parameter's value is a valid, non-empty string
	if param.Value != nil && param.Type == "string" {
		// perform $ notation replacement on string if any exist
		value = wskenv.InterpolateStringWithEnvVar(param.Value)
	}

	// JSON - Handle both cases, where value 1) is a string containing JSON, 2) is a map of JSON
	if param.Value != nil && param.Type == "json" {
		value, errorParser = resolveJSONParameter(filePath, paramName, param, value)
	}

	if param.Value != nil && param.Type == "slice" {
		value = wskenv.InterpolateStringWithEnvVar(param.Value)
		value = utils.ConvertInterfaceValue(value)
	}

	// Default value to zero value for the Type
	// Do NOT error/terminate as Value may be provided later by a Deployment file.
	if value == nil {
		value = getTypeDefaultValue(param.Type)
		// @TODO(): Need warning message here to warn of default usage, support for warnings (non-fatal)
		//msgs := []string{"Parameter [" + paramName + "] is not multiline format."}
		//return param.Value, utils.NewParserErr(filePath, nil, msgs)
	}

	// Trace Parameter struct after resolution
	//dumpParameter(paramName, param, "AFTER")
	//fmt.Printf("EXIT: Parameter [%s] type=[%v] value=[%v]\n", paramName, param.Type, value)
	return value, errorParser
}

// Provide custom Parameter marshalling and unmarshalling
type ParsedParameter Parameter

func (n *Parameter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux ParsedParameter

	// Attempt to unmarshal the multi-line schema
	if err := unmarshal(&aux); err == nil {
		n.multiline = true
		n.Type = aux.Type
		n.Description = aux.Description
		n.Value = aux.Value
		n.Required = aux.Required
		n.Default = aux.Default
		n.Status = aux.Status
		n.Schema = aux.Schema
		return nil
	}

	// If we did not find the multi-line schema, assume in-line (or single-line) schema
	var inline interface{}
	if err := unmarshal(&inline); err != nil {
		return err
	}

	n.Value = inline
	n.multiline = false
	return nil
}

func (n *Parameter) MarshalYAML() (interface{}, error) {
	if _, ok := n.Value.(string); len(n.Type) == 0 && len(n.Description) == 0 && ok {
		if !n.Required && len(n.Status) == 0 && n.Schema == nil {
			return n.Value.(string), nil
		}
	}

	return n, nil
}

// Provides debug/trace support for Parameter type
func dumpParameter(paramName string, param *Parameter, separator string) {

	fmt.Printf("%s:\n", separator)
	fmt.Printf("\t%s: (%T)\n", paramName, param)
	if param != nil {
		fmt.Printf("\t\tParameter.Description: [%s]\n", param.Description)
		fmt.Printf("\t\tParameter.Type: [%s]\n", param.Type)
		fmt.Printf("\t\t--> Actual Type: [%T]\n", param.Value)
		fmt.Printf("\t\tParameter.Value: [%v]\n", param.Value)
		fmt.Printf("\t\tParameter.Default: [%v]\n", param.Default)
	}
}
