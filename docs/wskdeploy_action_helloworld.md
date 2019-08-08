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

## The "Hello World" Action

As with most language introductions, in this first example we encode a simple "hello world" action, written in JavaScript, using an OpenWhisk Package Manifest YAML file.

It shows how to:
- declare a single Action named ‘hello_world’ within the ‘hello_world_package’ Package.
- associate the JavaScript function’s source code, stored in the file ‘src/hello.js’, to the ‘hello_world’ Action.

### Manifest file
#### _Example: “Hello world” using a NodeJS (JavaScript) action_
```yaml
packages:
  hello_world_package:
    version: 1.0
    license: Apache-2.0
    actions:
      hello_world:
        function: src/hello.js
```

where "hello.js", within the package-relative subdirectory named ‘src’, contains the following JavaScript code:

```javascript
function main(params) {
    msg = "Hello, " + params.name + " from " + params.place;
    return { greeting:  msg };
}
```

### Deploying

You can actually deploy the "hello world" manifest from the openwhisk-wskdeploy project directory if you have downloaded it from GitHub:

```sh
$ wskdeploy -m docs/examples/manifest_hello_world.yaml
```

### Invoking
```sh
$ wsk action invoke hello_world_package/hello_world --blocking
```

### Result
The invocation should return an 'ok' with a response that includes this result:

```json
"result": {
    "greeting": "Hello, undefined from undefined"
},
```

The output parameter '```greeting```' contains "_undefined_" values for the '```name```' and '```place```' input parameters as they were not provided in the manifest.

### Discussion

This "hello world" example represents the minimum valid Manifest file which includes only the required parts of the Package and Action descriptors.

In the above example,
- The Package and its Action were deployed to the user’s default namespace using the ‘package’ and 'action' names from the manifest.
  - ```/<default namespace>/hello_world_package/hello_world```
- The NodeJS default runtime (i.e., ```runtime: nodejs```) was selected automatically based upon the file extension '```.js```'' of the function's source file '```hello.js```'.

### Source code
The source code for the manifest and JavaScript files can be found here:
- [manifest_hello_world.yaml](examples/manifest_hello_world.yaml)
- [hello.js](examples/src/hello.js)

### Specification
For convenience, the Packages and Actions grammar can be found here:
- **[Packages](../specification/html/spec_packages.md#packages)**
- **[Actions](../specification/html/spec_actions.md#actions)**

---
<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a href="wskdeploy_package_minimal.md#packages">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_action_runtime.md#actions">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
