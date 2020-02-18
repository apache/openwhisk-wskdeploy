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

## Advanced parameters

This example builds on the previous [“Hello world" with typed input and output parameters](wskdeploy_action_typed_parms.md#actions) example with more robust input and output parameter declarations by using a multi-line format for declaration.

This example:
- shows how to declare input and output parameters on the action ‘```hello_world```’ using a multi-line grammar.

### Manifest File
If we want to do more than declare the type (i.e., ‘string’, ‘integer’, ‘float’, etc.) of the input parameter, we can use the multi-line grammar.

#### _Example: input and output parameters with explicit types and descriptions_
```yaml
packages:
  hello_world_package:
    ... # Package keys omitted for brevity
    actions:
      hello_world_advanced_parms:
        function: src/hello_plus.js
        runtime: nodejs@10
        inputs:
          name:
            type: string
            description: name of person
            default: unknown person
          place:
            type: string
            description: location of person
            value: the Shire
          children:
            type: integer
            description: Number of children
            default: 0
          height:
            type: float
            description: height in meters
            default: 0.0
        outputs:
          greeting:
            type: string
            description: greeting string
          details:
            type: string
            description: detailed information about the person
```

### Deploying
```sh
$ wskdeploy -m docs/examples/manifest_hello_world_advanced_parms.yaml
```

### Invoking
```sh
$ wsk action invoke hello_world_package/hello_world_advanced_parms --blocking
```

### Result
The invocation should return an 'ok' with a response that includes this result:
```json
"result": {
    "details": "You have 0 children and are 0 m. tall.",
    "greeting": "Hello, unknown person from the Shire"
},
```

### Discussion
- Describing the input and output parameter types, descriptions, defaults and other data:
  - enables tooling to validate values users may input and prompt for missing values using the descriptions provided.
  - allows verification that outputs of an Action are compatible with the expected inputs of another Action so that they can be composed in a sequence.
- The '```name```' input parameter was assigned the '```default```' key's value "```unknown person```".
- The '```place```' input parameter was assigned the '```value```' key's value "```the Shire```".

### Source code
The manifest file for this example can be found here:
- [manifest_hello_world_advanced_parms.yaml](examples/manifest_hello_world_advanced_parms.yaml)
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
    <td><a href="wskdeploy_action_typed_parms.md#actions">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_action_env_var_parms.md#actions">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
