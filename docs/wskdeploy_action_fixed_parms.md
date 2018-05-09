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

## Adding fixed input values to an Action

This example builds upon the previous “hello world” example and shows how fixed values can be supplied to the input parameters of an Action.

It shows how to:
- declare input parameters on the action ‘```hello_world```’ using a single-line grammar.
- add ‘```name```’ and ‘```place```’ as input parameters with the fixed values “```Sam```” and “```the Shire```” respectively.

### Manifest File
#### _Example: “Hello world” with fixed input values for ‘name’ and ‘place’_
```yaml
packages:
  hello_world_package:
    version: 1.0
    license: Apache-2.0
    actions:
      hello_world_fixed_parms:
        function: src/hello.js
        inputs:
          name: Sam
          place: the Shire
```

### Deploying
```sh
$ wskdeploy -m docs/examples/manifest_hello_world_fixed_parms.yaml
```

### Invoking
```sh
$ wsk action invoke hello_world_package/hello_world_fixed_parms --blocking
```

### Result
The invocation should return an 'ok' with a response that includes this result:
```json
"result": {
    "greeting": "Hello, Sam from the Shire"
},
```

### Discussion

In this example:
- The value for the ‘```name```’ input parameter would be set to “```Sam```”.
- The value for the ‘```place```’ input parameter would be set to “```the Shire```”.
- The wskdeploy utility would infer that both ‘```name```’ and ‘```place```’ input parameters to be of type ‘```string```’.

### Source code
The manifest file for this example can be found here:
- [manifest_hello_world_fixed_parms.yaml](examples/manifest_hello_world_fixed_parms.yaml)
- [hello.js](examples/src/hello.js)

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
    <td><a href="wskdeploy_action_runtime.md#actions">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_action_typed_parms.md#actions">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
