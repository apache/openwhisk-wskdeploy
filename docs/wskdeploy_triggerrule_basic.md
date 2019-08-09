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

# Triggers and Rules

## Creating a Trigger for an Action
This example shows how to create a Trigger that is compatible with the previous, more [advanced "Hello world" Action example](wskdeploy_action_advanced_parms.md#actions), which has multiple input parameters of different types, and connect them together using a Rule.

### Manifest File
#### _Example: “Hello world” Action with a compatible Trigger and Rule_
```yaml
package:
  name: hello_world_package
  version: 1.0
  license: Apache-2.0
  actions:
    hello_world_triggerrule:
      function: src/hello_plus.js
      runtime: nodejs
      inputs:
        name: string
        place: string
        children: integer
        height: float
      outputs:
        greeting: string
        details: string

  triggers:
    meetPerson:
      inputs:
        name: Sam
        place: the Shire
        children: 13
        height: 1.2

  rules:
    meetPersonRule:
      trigger: meetPerson
      action: hello_world_triggerrule
```

### Deploying
```sh
$ wskdeploy -m docs/examples/manifest_hello_world_triggerrule.yaml
```

### Invoking
First, let's try _"invoking"_ the '```hello_world_triggerrule```' Action directly without the Trigger.
```sh
$ wsk action invoke hello_world_package/hello_world_triggerrule --blocking
```

#### Result
```json
"result": {
  "details": "You have 0 children and are 0 m. tall.",
  "greeting": "Hello,  from "
},
```
As you can see, the results verify that the default values (i.e., empty strings and zeros) for the input parameters on the '```hello_world_triggerrule```' Action were used to compose the '```greeting```' and '```details```' output parameters. This result is expected since we did not bind any values or provide any defaults when we defined the '```hello_world_triggerrule```' Action in the manifest file.

### Triggering

Instead of invoking the Action, here try _"firing"_ the '```meetPerson```' Trigger:
```sh
$ wsk trigger fire meetPerson
```

#### Result
which results in an Activation ID:
```sh
ok: triggered /_/meetPerson with id a8e9246777a7499b85c4790280318404
```

The '```meetPerson```' Trigger is associated with '```hello_world_triggerrule```' Action the via the '```meetPersonRule```' Rule. We can verify that firing the Trigger indeed cause the Rule to be activated which in turn causes the Action to be invoked:
```sh
$ wsk activation list

d03ee729428d4f31bd7f61d8d3ecc043 hello_world_triggerrule
3e10a54cb6914b37a8abcab53596dcc9 meetPersonRule
5ff4804336254bfba045ceaa1eeb4182 meetPerson
```

we can then use the '```hello_world_triggerrule```' Action's Activation ID to see the result:
```sh
$ wsk activation get d03ee729428d4f31bd7f61d8d3ecc043
```
to view the actual results from the action:
```json
"result": {
   "details": "You have 13 children and are 1.2 m. tall.",
   "greeting": "Hello, Sam from the Shire"
}
```

which verifies that the parameter bindings of the values (i.e, _"Sam"_ (name), _"the Shire"_ (place), _'13'_ (age) and _'1.2'_ (height)) on the Trigger were passed to the Action's corresponding input parameters correctly.

### Discussion
- Firing the '```meetPerson```' Trigger correctly causes a series of non-blocking "activations" of the associated '```meetPersonRule```' Rule and subsequently the '```hello_world_triggerrule```' Action.
- The Trigger's parameter bindings were correctly passed to the corresponding input parameters on the '```hello_world_triggerrule```' Action.

### Source code
- [manifest_hello_world_triggerrule.yaml](examples/manifest_hello_world_triggerrule.yaml)
- [hello_plus.js](examples/src/hello_plus.js)

### Specification
For convenience, the Actions and Parameters grammar can be found here:
- **[Triggers and Rules](https://github.com/apache/openwhisk-wskdeploy/blob/master/specification/html/spec_trigger_rule.md#triggers-and-rules)**

---
<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a href="wskdeploy_sequence_basic.md#sequences">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_triggerrule_trigger_bindings.md#triggers-and-rules">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
