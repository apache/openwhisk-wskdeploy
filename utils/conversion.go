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
	"encoding/json"
	"fmt"
)

func convertInterfaceArray(in []interface{}) []interface{} {
	res := make([]interface{}, len(in))
	for i, v := range in {
		res[i] = ConvertInterfaceValue(v)
	}
	return res
}

func ConvertInterfaceMap(mapIn map[interface{}]interface{}) map[string]interface{} {
	mapOut := make(map[string]interface{})
	for k, v := range mapIn {
		mapOut[fmt.Sprintf("%v", k)] = ConvertInterfaceValue(v)
	}
	return mapOut
}

func ConvertInterfaceValue(value interface{}) interface{} {
	switch typedVal := value.(type) {
	case []interface{}:
		return convertInterfaceArray(typedVal)
	case map[string]interface{}:
	case map[interface{}]interface{}:
		return ConvertInterfaceMap(typedVal)
	case bool:
	case int:
	case int8:
	case int16:
	case int32:
	case int64:
	case float32:
	case float64:
	case string:
		return typedVal
	default:
		return fmt.Sprintf("%v", typedVal)
	}
	return value
}

// TODO() add a Print function to wskprint that calls this and adds the label
// TODO add prettyjson formatting as an option
func ConvertMapToJSONString(name string, mapIn interface{}) string {
	strMapOut, _ := json.MarshalIndent(mapIn, "", "  ")
	return fmt.Sprintf("%s: %s", name, string(strMapOut))
}
