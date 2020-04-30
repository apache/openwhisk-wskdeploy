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

## APIs

The API entity schema is used to define an OpenWhisk API within a manifest.

### Fields
<html>
<table width="100%">
 <tr>
  <th width="16%">
   <p>Key Name</p>
  </th>
  <th width="12%">
   <p>Required</p>
  </th>
  <th width="16%">
   <p>Value Type</p>
  </th>
  <td width="14%">
   <p>Default</p>
  </th>
  <th width="40%">
   <p>Description</p>
  </th>
 </tr>
 <tr>
  <td>method</td>
  <td>yes</td>
  <td>string</td>
  <td>N/A</td>
  <td>The HTTP method for the endpoint. All valid HTTP methods are supported, but a response type value of <b>http</b> may be required to correctly process the associated request.</td>
 </tr>
 <tr>
  <td>response</td>
  <td>no</td>
  <td>string</td>
  <td>http</td>
  <td>The response type or <i>content extension</i> used when the API Gateway invokes the web action. See <a href="https://github.com/apache/openwhisk/blob/master/docs/webactions.md#content-extensions">https://github.com/apache/openwhisk/blob/master/docs/webactions.md#content-extensions</a>.</p></td>
 </tr>
</table>
</html>

### Grammar

```yaml
<apiName>:
    <basePath>:
        <relativePath>:
            <actionName>:
                method: <get | delete | put | post | ...>
                response: <http | json | text | html>
```

### Example

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
              response: json
      goodbye-world:
        hello:
          world:
            hello_world:
              method: DELETE
              response: json
```

### Requirements

- The API name MUST be less than or equal to 256 characters.
- The API `basePath` value MUST start with a `/` character.
- The APi `relativePath` value MUST start with a `/` character.
- The API entity schema includes all required fields declared above.
- Only web actions, actions having `web-export` set to `true`, can be used as an API endpoint's action.
    - If needed, the action will be automatically converted to a web action during deployment.
- A valid API entity MUST have one or more valid endpoints defined.


### Notes

- When an API endpoint is being added to an existing API, the `apiName` in the manifest is ignored.
- See <a href="https://github.com/apache/openwhisk/blob/master/docs/webactions.md">https://github.com/apache/openwhisk/blob/master/docs/webactions.md</a> for the complete set of supported `response` values, also known as <i>content extensions</i>.
- Using a `response` value of `http` will give you the most control over the API request and response handling.

<!--
 Bottom Navigation
-->
---
<html>
<div align="center">
<a href="../README.md#index">Index</a>
</div>
</html>
