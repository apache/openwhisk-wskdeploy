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

# API Open API Spec

## The "Hello World" API

This example builds on the ["Hello World" Action](wskdeploy_action_helloworld.md#actions) example and ["Hello World" API](wdkdeploy_apigateway_helloworld.md) example by adding an Open API Specification to define an API that is invoked via an API call.

It shows how to:
- update the Action named ‘hello_world’ to expose it to the gateway.
- specify the API's endpoint that will trigger the action via an Open API Specification.

### Manifest file
#### _Example: “Hello world” action with project config to point to the API spec_
```yaml
project:
  config: open_api_spec.json
  packages:
    hello_world_package:
      version: 1.0
      license: Apache-2.0
      actions:
        hello_world:
          function: src/hello.js
        annotations:
          web-export: true
```
### Open API Specification
#### _Example: "Hello world" open api specification_
```json
{
    "swagger": "2.0",
    "info": {
        "version": "1.0",
        "title": "Hello World API"
    },
    "basePath": "/hello",
    "schemes": [
        "https"
    ],
    "consumes": [
        "application/json"
    ],
    "produces": [
        "application/json"
    ],
    "paths": {
        "/world": {
            "get": {
                "description": "Returns a greeting to the user!",
                "operationId": "hello_world",
                "responses": {
                    "200": {
                        "description": "Returns the greeting.",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }

    },
    "x-gateway-configuration": {
    "assembly": {
      "execute": [
        {
          "operation-switch": {
            "case": [
              {
                "operations": [
                  "getHello"
                ],
                "execute": [
                  {
                    "invoke": {
                      "target-url": "https://openwhisk.ng.bluemix.net/api/some/action/path.http",
                      "verb": "keep"
                    }
                  }
                ]
              }
            ],
            "otherwise": []
          }
        }
      ]
    }
  }
}
```
*NOTE*: Some providers such as IBM may have changed these keys. For example, IBM uses `x-ibm-configuration` instead of `x-gateway-configuration`.

There are two major differences from _"Hello World" API_ example:
- the root key is now project as the open api specification is a project wide concept.
- a new `config` key specifying where the Open API Specification is located.

The `config` key under `project` in the manifest file specifies where the Open API Specification is located. The keyword `config` was chosen to remain consistent with the `config-file` terminology in OpenWhisk CLI flag option. The Open API Specification describes in a JSON document the the base path, endpoint, HTTP verb, and other details describing the API. For example, the document above describes a GET endpoint at `/hello/world` that recieves JSON as input and returns JSON as output.

### Deploying

You cannot deploy the "hello world API gateway Open API Specification" manifest from the openwhisk-wskdeploy project directory directly. You need to update the `target-url`. This valued will be specific to your deployment/provider. An example of such a URL would be: `https://us-south.functions.cloud.ibm.com/api/v1/web/jdoe@ibm.com/hello_world_package/hello_world.json`. After filling out that value, you would deploy via:

```sh
$ wskdeploy -m docs/examples/manifest_hello_world_apigateway_open_api_spec.yaml
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

This "hello world" example represents the minimum valid Manifest file which includes only the required parts of the Project, Package, Action and Open API Specification location.

### Source code
The source code for the manifest and JavaScript files can be found here:
- [manifest_hello_world_apigateway.yaml](examples/manifest_hello_world_apigateway_open_api_spec.yaml)
- [open_api_spec.json](examples/open_api_spec.json)
- [hello.js](examples/src/hello.js)

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
