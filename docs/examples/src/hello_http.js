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
 * This action returns a greeting wrapped into a HTTP response that can be used
 * by the API Gateway. It also checks for mandatory attributes and returns a
 * HTTP error response when any attribute is missing.
 */
function main(params) {
    if(params.name && params.place) {
        return {
            body: {
                greeting: `Hello ${params.name} from ${params.place}!`
            },
            statusCode: 200,
            headers: {'Content-Type': 'application/json'}
        };
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
