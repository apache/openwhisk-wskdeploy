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

function main(params) {
    let step = params.$step || 0
    delete params.$step
    package_name = "conductorPackage1"
    switch (step) {
        case 0: return { action: package_name+'/triple', params, state: { $step: 1 } }
        case 1: return { action: package_name+'/increment', params, state: { $step: 2 } }
        case 2: return { params }
    }
}
