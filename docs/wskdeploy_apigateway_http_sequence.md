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

# API Gateway HTTP response and sequence

## HTTP status handling in sequences

The previous example in [API Gateway HTTP response](wskdeploy_apigateway_http.md#api-gateway-http-response) showed how to handle custom HTTP responses. The downside of the approach used in that example is that it requires the action to be updated to return a HTTP response structure. What if we want to re-use actions that have no knowledge of HTTP? One simple options is to use a sequence. The HTTP response structure needs to be included at any point where the sequence terminates. A sequence typically terminates when the result of the last action is returned. It will terminate early if an error occurs in the sequence. Therefore a simple rule of thumb for combining sequences and HTTP responses is:
- Wrap the final result in a HTTP response structure,
- Wrap any error in a HTTP response structure.

The first step is to create a manifest with a sequence that validates input, calls the existing `hello.js` action and wraps the result.

### Manifest file
#### _Example: “Hello world” action with HTTP response and sequence_
```yaml
packages:
  hello_world_package:
    version: 1.0
    license: Apache-2.0
    actions:
      hello_validate:
        function: src/hello_http_validate.js
      hello:
        function: src/hello.js
      hello_wrap:
        function: src/hello_http_wrap.js
    sequences:
      hello_world:
        actions: hello_validate, hello, hello_wrap
        web: true
    apis:
      hello-world:
        hello:
          world:
            hello_world:
              method: GET
              response: http
```

The key points of this file are:
- multiple actions are defined to validate input, run the main action and wrap the output,
- a sequence is defined that combines all three actions and that has the `web` field set to `true`,
- a `response` field with value `http` is set on the the API.

### Validate action file
#### _Example: validation action with structured HTTP error_
```javascript
function main(params) {
    if(params.name && params.place) {
        return params;
    } else {
        return {
            error: {
                body: {
                    message: 'Attributes name and place are mandatory'
                },
                statusCode: 400,
                headers: {'Content-Type': 'application/json'}
            }
        }
    }
}
```

The first action validates inputs. If everything is fine, it returns the parameters unchanged for those to be passed to the next action in the sequence. If validation fails, it create an error HTTP structure. The presence of the `error` field in the result will make OpenWhisk abort the sequence and return what is contained in that field, which in turn is transformed into an HTTP response by the API Gateway.

The second action in the chain is the simple greeting action from [The "Hello World" Action](wskdeploy_action_helloworld.md#actions).

### Wrap action file
#### _Example: action that wraps a simple result into a HTTP response_
```javascript
function main(params) {
    return {
        body: params,
        statusCode: 200,
        headers: {'Content-Type': 'application/json'}
    };
}
```

The last action in the chain simply wraps the result from the previous action into a HTTP response structure with a status code and headers.

### Deploying

You can actually deploy the "API Gateway HTTP" manifest from the openwhisk-wskdeploy project directory if you have downloaded it from GitHub:

```sh
$ wskdeploy -m docs/examples/manifest_hello_world_apigateway_http_sequence.yaml
```

### Invoking

Check the full URL of your API first:
```sh
$ wsk api list
```

This will return some information on the API, including the full URL, which
should end with `hello/world`. It can then be invoked:

```sh
$ curl -i <url>
```

### Result
The invocation should return a JSON response with status code `400` that includes this result:

```json
{
    "message": "Attributes name and place are mandatory"
}
```

This shows our error handling code working. To get a valid response, we need to provide the `name` and `place` parameters:

```sh
$ curl -i <url>?name=World&place=Earth
```

You should then see a JSON response with status code `200` and the following result:

```json
{
    "greeting": "Hello World from Earth!"
}
```

### Discussion

By combining HTTP responses and sequences, you can re-use existing actions that are not designed to return HTTP responses by adding the necessary wrapper to the final result. You need to be careful how errors are handled as they will short-circuit the sequence execution and return early.

### Source code
The source code for the manifest and JavaScript files can be found here:
- [manifest_hello_world_apigateway_http_sequence.yaml](examples/manifest_hello_world_apigateway_http_sequence.yaml)
- [hello_http_validate.js](examples/src/hello_http_validate.js)
- [hello.js](examples/src/hello.js)
- [hello_http_wrap.js](examples/src/hello_http_wrap.js)

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
    <td><a href="wskdeploy_apigateway_http.md#api-gateway-http-response">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <!--<td><a href="">next&nbsp;&gt;&gt;</a></td>-->
  </tr>
</table>
</div>
</html>
