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
 * Join the fellowship
 */
function main(params) {
    console.log("params: " + JSON.stringify(params, null, 4));
    var member = {name:"", place:"", region:"", occupation:"", joined:"", organization:"", item:"" };

    // The organization being joined is fixed
    member.organization = "fellowship";

    // Fill in a member record from parameters
    member.name = params.name;
    member.place = params.place;
    member.occupation = params.job;

    // Save the current timestamp when we created the member record
    member.joined = Date.now();

    console.log("member: " + JSON.stringify(member, null, 4));
    return { member: member };
}
