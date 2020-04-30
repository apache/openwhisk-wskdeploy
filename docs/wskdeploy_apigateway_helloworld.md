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

# API Gateway

## The "Hello World" API

This example builds on the ["Hello World" Action](wskdeploy_action_helloworld.md#actions) example by adding an API definition on top of that action so that I can be queried via an HTTP call.

It shows how to:
- update the Action named ‘hello_world’ to expose it to the gateway.
- specify the API's endpoint that will trigger the action.

### Manifest file
#### _Example: “Hello world” action with API_
```yaml
packages:
  hello_world_package:
    version: 1.0
    license: Apache-2.0
    actions:
      hello_world:
        function: src/hello.js
        annotations:
          web-export: true
    apis:
      hello-world:
        hello:
          world:
            hello_world:
              method: GET
```

There are two key changes to this file:
- the `hello_world` action now has the `web-export` annotation set to `true`.
- a new `apis` block has been created.

The `apis` block contains a number of groups of API endpoint. Each endpoint is then defined by the hierarchy. In this case, we are creating the `hello/world` endpoint. The leaf in the structure specifies the action to trigger when the given HTTP verb is sent to that endpoint, in this case, when the HTTP verb `GET` is used on the `hello/world` endpoint, trigger the `hello_world` action.

### Deploying

You can actually deploy the "hello world API gateway" manifest from the openwhisk-wskdeploy project directory if you have downloaded it from GitHub:

```sh
$ wskdeploy -m docs/examples/manifest_hello_world_apigateway.yaml
```

### Invoking

Check the full URL of your API first:
```sh
$ wsk api list
```

This will return some information on the API, including the full URL, which
should end with `hello/world`. It can then be invoked:

```sh
$ curl <url>
```

### Result
The invocation should return a JSON response that includes this result:

```json
{
    "greeting": "Hello, undefined from undefined"
}
```

The output parameter '```greeting```' contains "_undefined_" values for the '```name```' and '```place```' input parameters as they were not provided in the manifest or the HTTP call. You can provide them as query parameters:

```sh
$ curl <url>?name=World&place=Earth
```

### Discussion

This "hello world" example represents the minimum valid Manifest file which includes only the required parts of the Package, Action and API descriptors.

### Source code
The source code for the manifest and JavaScript files can be found here:
- [manifest_hello_world_apigateway.yaml](examples/manifest_hello_world_apigateway.yaml)
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
    <td><a href="wskdeploy_triggerrule_trigger_bindings.md#triggers-and-rules">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_apigateway_sequence.md#api-gateway-sequence">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
