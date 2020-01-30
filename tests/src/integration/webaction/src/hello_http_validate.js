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

/*
 * This action validates the parameters passed to the hello world action and
 * returns an error if any of them is missing.
 */
function main(params) {
    if(params.name && params.place) {
        return params;
    } else {
        return {
            error: {
                body: {
                    message: 'Attributes name and place are mandatory'
                },
                statusCode: 400,
                headers: {'Content-Type': 'application/json'}
            }
        }
    }
}
