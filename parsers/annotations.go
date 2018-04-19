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
	"fmt"
	"reflect"
)

func encodeBoolAnnotations(annotation interface{}) interface{} {
	return annotation.(bool)
}

func encodeIntAnnotations(annotation interface{}) interface{} {
	return annotation.(int)
}

func encodeFloatAnnotations(annotation interface{}) interface{} {
	return annotation.(float64)
}

func encodeListAnnotations(annotation interface{}) interface{} {
	list := make([]interface{}, 0)
	annotationList := annotation.([]interface{})
	for _, s := range annotationList {
		l := EncodeAnnotations(s)
		list = append(list, l)
	}
	return list
}

func encodeMapAnnotations(annotation interface{}) interface{} {
	mapValue := make(map[string]interface{})
	for k, v := range annotation.(map[interface{}]interface{}) {
		key := fmt.Sprintf("%v", k)
		value := EncodeAnnotations(v)
		mapValue[key] = value
	}
	return mapValue
}

func EncodeAnnotations(annotation interface{}) interface{} {
	var annotationValue interface{}
	if reflect.ValueOf(annotation).Kind() == reflect.Bool {
		annotationValue = encodeBoolAnnotations(annotation)
	} else if reflect.ValueOf(annotation).Kind() == reflect.Int {
		annotationValue = encodeIntAnnotations(annotation)
	} else if reflect.ValueOf(annotation).Kind() == reflect.Float32 {
		annotationValue = encodeFloatAnnotations(annotation)
	} else if reflect.ValueOf(annotation).Kind() == reflect.Float64 {
		annotationValue = encodeFloatAnnotations(annotation)
	} else if reflect.ValueOf(annotation).Kind() == reflect.Slice {
		annotationValue = encodeListAnnotations(annotation)
	} else if reflect.ValueOf(annotation).Kind() == reflect.Map {
		annotationValue = encodeMapAnnotations(annotation)
	} else {
		annotationValue = fmt.Sprintf("%v", annotation)
	}
	return annotationValue
}
