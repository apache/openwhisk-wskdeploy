<!--
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
-->

# Sequences

## Creating a basic Action sequence

OpenWhisk supports creating a new, named composite action from sequencing multiple, compatible Actions.

This example:
- Shows how to use a Manifest to sequence three actions together to:
  - accept raw parameter data, creating a structured record from it
  - process the initial data and enhance it
  - augment the data via the sequenced actions

### Manifest File

#### _Example: Composing 3 actions into a sequence._
```yaml
# Example: processing data in a sequence
package:
  name: fellowship_package
  ...
  actions:
    member_join:
      function: src/member_join.js
      inputs:
        name:
          type: string
          description: name of person
          default: unknown
        place:
          type: string
          description: location of person
          default: unknown
        job:
          type: string
          description: current occupation
          default: 0
      outputs:
        member:
          type: json
          description: member record
    member_process:
      function: src/member_process.js
      inputs:
        member: {}
    member_equip:
      function: src/member_equip.js
      inputs:
        member: {}
  sequences:
    fellowship_membership:
      actions: member_join, member_process, member_equip
```

#### ```member_join.js``` code snippet:
```javascript
function main(params) {

    var member = {name:"", place:"", region:"", occupation:"", joined:"", organization:"", item:"" };

    // The organization being joined is fixed
    member.organization = "fellowship";

    // Fill in a member record from parameters
    member.name = params.name;
    member.place = params.place;
    member.occupation = params.job;

    // Save the current timestamp when we created the member record
    member.joined = Date.now();

    return { member: member };
}
```

#### ```member_process.js``` code snippet:
```javascript
const regionMap = new Map([
    ['Hobbiton', 'Eriador'],
    ['Rivendell', 'Eriador'],
    ['Minas Tirith', 'Gondor'],
    ['Lake Town', 'Rhovanion'],
    ['Minas Morgul', 'Mordor'],
]);

function main(params) {

    // Augment the member (record) created in the previous Action
    member = params.member;
    member.region = regionMap.get(member.place) || "unknown";
    member.date = new Date(member.joined).toLocaleDateString();
    member.time = new Date(member.joined).toLocaleTimeString();

    return { member: member };
}
```

#### ```member_equip.js``` code snippet:
```javascript
const equipmentMap = new Map([
    ['gentleman', 'ring'],
    ['wizard', 'staff'],
    ['archer', 'bow'],
    ['knight', 'sword'],
]);

function main(params) {

    // Equip the member based upon their occupation
    member = params.member;
    member.item = equipmentMap.get(member.occupation) || "None";

    return { member: member };
}


### Deploying
```sh
$ wskdeploy -m docs/examples/manifest_sequence_basic.yaml
```

### Invoking
```sh
$ wsk action invoke fellowship_package/fellowship_membership -p name frodo -p place Hobbiton -p job gentleman  -b
```

### Result
The invocation should return a 'success' response that includes this result:
```json
"result": {
    "member": {
        "joined": 1507155846307,
        "name": "frodo",
        "occupation": "gentleman",
        "organization": "fellowship",
        "place": "Hobbiton",
        "region": "Eriador",
        "date": "10/4/2017",
        "time": "10:24:06 PM",
        "item": "ring"
    }
}
```

and with three log entries (one for each Action in the sequence):
```json
"logs": [
    "4fdb1f27c6c84ca09b1f27c6c83ca0c6",
    "038567b035b743018567b035b70301c9",
    "aa730c99319f4b8bb30c99319f9b8b3b"
]
```

### Discussion

we can inspect the logs from the first Action "```member_join```" to view its input parameters "```params```" which where passed on the command line invocation:
```sh
$ wsk activation logs 4fdb1f27c6c84ca09b1f27c6c83ca0c6
params: {
    "name": "frodo",
    "place": "Hobbiton",
    "job": "gentleman"
 }
```
the input parameters are augmented by the first Action in the sequence to produce the output "member" object:

```json
member: {
    "organization": "fellowship",
    "name": "frodo",
    "place": "Hobbiton",
    "region": "",
    "occupation": "gentleman",
    "joined": 1507155846307,  // Date() in msec.
    "item": ""
}

```

the second Action in the sequence further processes and adds to the "```member```" data:
```sh
$ wsk activation logs 038567b035b743018567b035b70301c9

```

```json
member: {
    "organization": "fellowship",
    "name": "frodo",
    "region": "Eriador",
    "place": "Hobbiton",
    "occupation": "gentleman",
    "joined": 1507155846307,
    "date": "10/4/2017",
    "time": "10:24:06 PM",
    "item": ""
}
```

Finally, the last Action in the sequence adds the "```item```" value to the "```member```" data to produce the completed record:
```sh
$ wsk activation logs aa730c99319f4b8bb30c99319f9b8b3b

```

```json
member: {
    "organization": "fellowship",
    "name": "frodo",
    "region": "Eriador",
    "place": "Hobbiton",
    "occupation": "gentleman",
    "joined": 1507155846307,
    "date": "10/4/2017",
    "time": "10:24:06 PM",
    "item": "ring"
}
```

### Source code
The source code for the manifest and JavaScript files can be found here:
- [manifest_hello_world.yaml](examples/manifest_sequence_basic.yaml)
- [member_join.js](examples/src/member_join.js)
- [member_process.js](examples/src/member_process.js)
- [member_equip.js](examples/src/member_equip.js)

### Specification
For convenience, the Actions and Parameters grammar can be found here:
- **[Sequences](../specification/html/spec_sequences.md#sequences)**
- **[Actions](../specification/html/spec_actions.md#actions)**

---
<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a href="wskdeploy_action_env_var_parms.md#actions">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_triggerrule_basic.md#triggers-and-rules">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
