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

# API Gateway sequence

## Sequencing actions

This example builds on the ["Hello World" API Gateway](wskdeploy_apigateway_helloworld.md#api-gateway) example by combining multiple actions in a sequence to create a complex result and making this sequence available through the API Gateway.

It shows how to:
- create a sequence that chains multiple actions.
- export the sequence via the gateway.

### Manifest file
#### _Example: “Hello world” sequence with API_
```yaml
packages:
  hello_world_package:
    version: 1.0
    license: Apache-2.0
    actions:
      hello_basic:
        function: src/hello.js
      hello_goodday:
        function: src/hello_goodday.js
    sequences:
      hello_world:
        actions: hello_basic, hello_goodday
        web: true
    apis:
      hello-world:
        hello:
          world:
            hello_world:
              method: GET
```

There are two key changes to this file:
- we now have two actions define and neither of them is web available.
- we have a new sequence that is web available.
- the `apis` block now refers to the sequence rather than the action.

### Deploying

You can actually deploy the "API Gateway sequence" manifest from the openwhisk-wskdeploy project directory if you have downloaded it from GitHub:

```sh
$ wskdeploy -m docs/examples/manifest_hello_world_apigateway_sequence.yaml
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
    "greeting": "Hello, undefined from undefined, have a good day!"
}
```

The output parameter '```greeting```' contains "_undefined_" values for the '```name```' and '```place```' input parameters as they were not provided in the manifest or the HTTP call. You can provide them as query parameters:

```sh
$ curl <url>?name=World&place=Earth
```

### Discussion

This example is very basic but it shows the power of OpenWhisk in creating complex processes by chaining together simple actions and making the resulting sequence available via the API Gateway. This allows you to build small modular actions rather than big bulky ones. And because OpenWhisk allows actions to be written with different programming languages, you can mix and match, for example chaining a combination of JavaScript and Python actions in the same sequence.

### Source code
The source code for the manifest and JavaScript files can be found here:
- [manifest_hello_world_apigateway_sequence.yaml](examples/manifest_hello_world_apigateway_sequence.yaml)
- [hello.js](examples/src/hello.js)
- [hello_goodday.js](examples/src/hello_goodday.js)

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
    <td><a href="wskdeploy_apigateway_helloworld.md#api-gateway">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_apigateway_http.md#api-gateway-http-response">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
