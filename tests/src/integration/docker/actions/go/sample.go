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

package main

import "encoding/json"
import "fmt"
import "os"

func main() {
	//program receives one argument: the JSON object as a string
	arg := os.Args[1]

	// unmarshal the string to a JSON object
	var obj map[string]interface{}
	json.Unmarshal([]byte(arg), &obj)

	// can optionally log to stdout (or stderr)
	fmt.Println("hello Go action")

	name, ok := obj["name"].(string)
	if !ok {
		name = "Stranger"
	}

	// last line of stdout is the result JSON object as a string
	msg := map[string]string{"msg": ("Hello, " + name + "!")}
	res, _ := json.Marshal(msg)
	fmt.Println(string(res))
}
