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

# API Gateway HTTP response

## HTTP status handling

The API Gateway does a lot of things for you and in particular provides a default error handling behaviour that assumes that a successful result should map to a HTTP response with status code `200: OK` and an error result should map to a response with status code `400: Bad Request`. This works for simple use cases but when building a complex API, you may want to have more control over what is returned. Some of the things you may want to do include:
- Return a more precise status code
- Set specific response headers
- Return data in a format other than JSON

The API Gateway allows you to do this by specifying the HTTP response type on individual API endpoints. Doing this gives you more control by removing the default behaviour of the API Gateway but this means that you also need to adapt what your actions return to tell the API Gateway what do do.

This example is a variation of the ["Hello World" API Gateway](wskdeploy_apigateway_helloworld.md#api-gateway) example that shows the changes that need to apply:
- to the manifest,
- and to the underlying actions.

### Manifest file
#### _Example: “Hello world” action with HTTP response_
```yaml
packages:
  hello_world_package:
    version: 1.0
    license: Apache-2.0
    actions:
      hello_world:
        function: src/hello_http.js
        web-export: true
    apis:
      hello-world:
        hello:
          world:
            hello_world:
              method: GET
              response: http
```

There are two key changes to this file:
- a `response` field with value `http` was added to the API,
- the `hello_world` action points to a different file detailed below.

### Action file
#### _Example: “Hello world” action with structured HTTP response_
```javascript
function main(params) {
    if(params.name && params.place) {
        return {
            body: {
                greeting: `Hello ${params.name} from ${params.place}`
            },
            statusCode: 200,
            headers: {'Content-Type': 'application/json'}
        };
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

Because a HTTP response disables the API Gateway default handling, you have to provide a more complex response that fills in the blanks and provides:
- a `statusCode` field with the required HTTP status code,
- a `body` field that contains your normal payload,
- an optional `headers` field that includes any HTTP header you want to set, typically `Content-Type`.

If you don't provide this structure, the API Gateway will generate a HTTP response with status code `204: No Content` and an empty body. If this occurs when it shouldn't, it's probably a sign that you have a HTTP response specified with the gateway but the underlying action doesn't return this structure.

When you want to return an error, you need to provide the same structure wrapped into an `error` object. If you don't wrap it into an `error` object, it will still work from an HTTP perspective but OpenWhisk will not recognise it as an error.

This structure will work with any language that is supported by OpenWhisk, such as python or Java. If you are using JavaScript, you can make use of `Promise.resolve` and `Promise.reject` to make your code more readable by removing the need for the `error` wrapper:

### Action file
#### _Example: “Hello world” action with HTTP response using Promise_
```javascript
function main(params) {
    if(params.name && params.place) {
        return Promise.resolve({
            body: {
                greeting: `Hello ${params.name} from ${params.place}`
            },
            statusCode: 200,
            headers: {'Content-Type': 'application/json'}
        });
    } else {
        return Promise.reject({
            body: {
                message: 'Attributes name and place are mandatory'
            },
            statusCode: 400,
            headers: {'Content-Type': 'application/json'}
        });
    }
}
```

### Deploying

You can actually deploy the "API Gateway HTTP" manifest from the *openwhisk-wskdeploy* project directory if you have downloaded it from GitHub:

```sh
$ wskdeploy -m docs/examples/manifest_hello_world_apigateway_http.yaml
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

This example shows how you can have full control over the HTTP response produced by the API Gateway. This is essential when building a good REST API. Taking this example further, you can also return different payloads, such as CSV or XML files.

### Source code
The source code for the manifest and JavaScript files can be found here:
- [manifest_hello_world_apigateway_http.yaml](examples/manifest_hello_world_apigateway_http.yaml)
- [hello_http.js](examples/src/hello_http.js)
- [hello_http_promise.js](examples/src/hello_http_promise.js)

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
    <td><a href="wskdeploy_apigateway_sequence.md#api-gateway-sequence">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_apigateway_http_sequence.md#api-gateway-http-response-and-sequence">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
