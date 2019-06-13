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

const regionMap = new Map([
    ['Shire', 'Eriador'],
    ['the Shire', 'Eriador'],
    ['Hobbiton', 'Eriador'],
    ['Bree', 'Eriador'],
    ['Rivendell', 'Eriador'],
    ['Minas Tirith', 'Gondor'],
    ['Esgaroth', 'Rhovanion'],
    ['Dale', 'Rhovanion'],
    ['Lake Town', 'Rhovanion'],
    ['Minas Morgul', 'Mordor'],
]);

function main(params) {
    console.log("params: " + JSON.stringify(params, null, 4));

    if(params.member && typeof params.member === "object"){
        // Augment the member (record) created in the previous Action
        member = params.member;
        member.region = regionMap.get(member.place) || "unknown";
        member.date = new Date(member.joined).toLocaleDateString();
        member.time = new Date(member.joined).toLocaleTimeString();
    }
    else
        throw new Error("Invalid parameter: 'member'. type="+typeof(member) + ", expected object).")

    console.log("member: " + JSON.stringify(member, null, 4));
    return { member: member };
}
