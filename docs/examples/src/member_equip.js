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

const equipmentMap = new Map([
    ['gentleman', 'ring'],
    ['wizard', 'staff'],
    ['archer', 'bow'],
    ['knight', 'sword'],
    ['barbarian', 'club'],
    ['warrior', 'axe'],
    ['thief', 'dagger'],
    ['gardener', 'rope'],
    ['squire', 'shortsword'],
    ['scout', 'horn'],
]);

function main(params) {
    console.log("params: " + JSON.stringify(params, null, 4));

    if(params.member && typeof params.member === "object"){
        // Equip the member based upon their occupation
        member = params.member;
        member.item = equipmentMap.get(member.occupation) || "None";
    }
    else
        throw new Error("Invalid parameter: 'member'. type="+typeof(member) + ", expected object).")

    console.log("member: " + JSON.stringify(member, null, 4));
    return { member: member };
}
