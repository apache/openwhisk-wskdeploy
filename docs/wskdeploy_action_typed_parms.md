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

# Actions

## Typed parameters

This example shows the 'Hello world' example with typed input and output Parameters.

It shows how to:
- declare input and output parameters on the action '```hello_world```' using a simple, single-line format.
- add two input parameters, '```name```' and '```place```', both of type '```string```' to the '```hello_world```' action.
- add an '```integer```' parameter, '```age```', to the action.
- add a '```float```' parameter, '```height```', to the action.
- add two output parameters, '```greeting```' and '```details```', both of type '```string```', to the action.

### Manifest File
#### _Example: 'Hello world' with typed input and output parameter declarations_
```yaml
packages:
  hello_world_package:
    ... # Package keys omitted for brevity
    actions:
      hello_world_typed_parms:
        function: src/hello_plus.js
        inputs:
          name: string
          place: string
          children: integer
          height: float
        outputs:
          greeting: string
          details: string
```
where the function '```hello_plus.js```', within the package-relative subdirectory named ‘```src```’, is updated to use the new parameters:
```javascript
function main(params) {
    msg = "Hello, " + params.name + " from " + params.place;
    family = "You have " + params.children + " children ";
    stats = "and are " + params.height + " m. tall.";
    return { greeting:  msg, details: family + stats };
}
```

### Deploying
```sh
$ wskdeploy -m docs/examples/manifest_hello_world_typed_parms.yaml
```

### Invoking
```sh
$ wsk action invoke hello_world_package/hello_world_typed_parms --blocking
```

### Result
The invocation should return an 'ok' with a response that includes this result:
```json
"result": {
  "details": "You have 0 children and are 0 m. tall.",
  "greeting": "Hello,  from "
},
```

### Discussion

In this example:

- The default value for the '```string```' type is the empty string (i.e., \"\"); it was assigned to the '```name```' and '```place```' input parameters.
- The default value for the '```integer```' type is zero (0); it was assigned to the '```age```' input parameter.
- The default value for the '```float```' type is zero (0.0f); it was assigned to the '```height```' input parameter.

### Source code
The manifest file for this example can be found here:
- [manifest_hello_world_typed_parms.yaml](examples/manifest_hello_world_typed_parms.yaml)
- [hello_plus.js](examples/src/hello_plus.js)

### Specification
For convenience, the Actions and Parameters grammar can be found here:
- **[Actions](../specification/html/spec_actions.md#actions)**
- **[Parameters](../specification/html/spec_parameters.md#parameters)**

---
<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a href="wskdeploy_action_fixed_parms.md#actions">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_action_advanced_parms.md#actions">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
