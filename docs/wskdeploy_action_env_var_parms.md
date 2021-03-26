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

## Setting values from environment variables

This example shows how to set input parameter values using environment variables. “hello world” example and shows how fixed values can be supplied to the input parameters of an Action.

It shows how to:
- declare input parameters on the action ‘```hello_world```’ whose values are set (i.e., bound) from values taken from environment variables.
  - using both single-line and multi-line grammars.
- concatenate string parameter values with values provided from environment variables.
- use this feature using both single-line and multi-line grammars.

### Manifest File
#### _Example: “Hello world” with input values set from environment variables_
```yaml
packages:
  hello_world_package:
    version: 1.0
    license: Apache-2.0
    actions:
      hello_world_env_var_parms:
        function: src/hello.js
        inputs:
          name: $FIRSTNAME
          place: ${TOWN}, ${COUNTRY}
```

### Deploying
Set the values for the three environment variables expected by the Action and deploy:
```sh
$ FIRSTNAME=Sam TOWN="the Shire" COUNTRY="Middle-earth" wskdeploy -m docs/examples/manifest_hello_world_env_var_parms.yaml
```

### Invoking
```sh
$ wsk action invoke hello_world_package/hello_world_env_var_parms --blocking
```

### Result
The invocation should return an 'ok' with a response that includes this result:
```json
"result": {
   "greeting": "Hello, Sam from the Shire, Middle-earth"
},
```

### Invoking
if we modify the three environment variables to different values, update the action and invoke it again:
```sh
$ export FIRSTNAME=Elrond TOWN="Rivendell" COUNTRY="M.E."
$ wskdeploy -m docs/examples/manifest_hello_world_env_var_parms.yaml
$ wsk action invoke hello_world_package/hello_world_env_var_parms --blocking
```

### Result
the result will reflect the changes to the environment variables:
```json
"result": {
   "greeting": "Hello, Elrond from Rivendell, M.E."
},
```

### Discussion

In this example:
- it was shown how values provided by environment variables, within the execution of the wskdeploy utility, could be bound to input parameter values within a Manifest file.
- we further demonstrated how string values from environment variables could be concatenated with other strings within a Manifest file

### Notes:
- These methods for binding environment variables to input parameters are also available within Deployment files.

### Source code
The manifest file for this example can be found here:
- [manifest_hello_world_env_var_parms.yaml](examples/manifest_hello_world_env_var_parms.yaml)
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
    <td><a href="wskdeploy_action_advanced_parms.md#actions">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_sequence_basic.md#sequences">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
