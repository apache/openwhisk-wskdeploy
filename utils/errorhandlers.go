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
	"errors"
	"fmt"
	"log"
	"os"
)

// Check is a util function to panic when there is an error.
func Check(e error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Runtime panic : %v", err)
		}
	}()

	if e != nil {
		log.Printf("%v", e)
		erro := errors.New("Error happened during execution, please type 'wskdeploy -h' for help messages.")
		log.Printf("%v", erro)
		if Flags.WithinOpenWhisk {
			PrintOpenWhiskError(e.Error())
		} else {
			os.Exit(1)
		}

	}
}

func PrintOpenWhiskError(err string) {
	fmt.Print(`{"error":"` + err + `"}`)
}
