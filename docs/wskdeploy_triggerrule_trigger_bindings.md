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

## Using a Deployment file to bind Trigger parameters
This example builds on the previous [Trigger and Rule example](wskdeploy_triggerrule_basic.md#triggers-and-rules) and will demonstrate how to use a Deployment File to bind values to a Trigger’s input parameters and apply them against a compatible Manifest File.

### Manifest File
Let’s use a variant of the [Manifest file from the previous example](examples/manifest_hello_world_triggerrule.yaml); however, we will leave the parameters on the ‘```meetPerson```’ Trigger unbound and having only Type declarations for each.

#### _Example: “Hello world” Action, Trigger and Rule with no Parameter bindings_
```yaml
packages:
  hello_world_package
  ... # Package keys omitted for brevity
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
          name: string
          place: string
          children: integer
          height: float

    rules:
      meetPersonRule:
        trigger: meetPerson
        action: hello_world_triggerrule
```

### Deployment file
Let’s create a Deployment file that is designed to be applied to the Manifest file (above) which will contain the parameter bindings (i.e., the values) for the '```meetPerson```' Trigger.

#### _Example: Deployment file that binds parameters to the '```meetPerson```' Trigger_
```yaml
project:
  packages:
      hello_world_package:
        triggers:
          meetPerson:
            inputs:
              name: Elrond
              place: Rivendell
              children: 3
              height: 1.88
```
As you can see, the package name '```hello_world_package```' and the trigger name '```meetPerson```' both match the names in the corresponding Manifest file.

### Deploying
Provide the Manifest file and the Deployment file to the wskdeploy utility:
```sh
$ wskdeploy -m docs/examples/manifest_hello_world_triggerrule_unbound.yaml -d docs/examples/deployment_hello_world_triggerrule_bindings.yaml
```

### Triggering
Fire the '```meetPerson```' Trigger:
```sh
$ wsk trigger fire meetPerson
```

#### Result
Find the activation ID for the “```hello_world_triggerrule```’ Action that firing the Trigger initiated and get the results from the activation record:

```sh
$ wsk activation list

3a7c92468b4e4170bc92468b4eb170f1 hello_world_triggerrule
afb2c02bb686484cb2c02bb686084cab meetPersonRule
9dc9324c601a4ebf89324c601a1ebf4b meetPerson

$ wsk activation get 3a7c92468b4e4170bc92468b4eb170f1
```

```json
"result": {
  "details": "You have 3 children and are 1.88 m. tall.",
  "greeting": "Hello, Elrond from Rivendell"
}
```

### Discussion
- The '```hello_world_triggerrule```' Action and the '```meetPerson```' Trigger in the Manifest file both had input parameter declarations that had no values assigned to them (only Types).
- The matching '```meetPerson```' Trigger in the Deployment file had values bound its parameters.
- The ```wskdeploy``` utility applied the parameter values (after checking for Type compatibility) from the Deployment file to the matching (by name) parameters within the Manifest file.

### Source code
- [manifest_hello_world_triggerrule_unbound.yaml](examples/manifest_hello_world_triggerrule_unbound.yaml)
- [deployment_hello_world_triggerrule_bindings.yaml](examples/deployment_hello_world_triggerrule_bindings.yaml)
- [hello_plus.js](examples/src/hello_plus.js)

### Specification
For convenience, the Actions and Parameters grammar can be found here:
- **[Triggers](https://github.com/apache/openwhisk-wskdeploy/blob/master/specification/html/spec_trigger.md#triggers)**
- **[Rules](https://github.com/apache/openwhisk-wskdeploy/blob/master/specification/html/spec_rule.md#rules)**

---
<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a href="wskdeploy_triggerrule_basic.md#triggers-and-rules">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_apigateway_helloworld.md#packages">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
