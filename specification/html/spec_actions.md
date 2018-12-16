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

## Actions

#### Subsections
- [Fields](#fields)
- [Requirements](#requirements)
- [Notes](#notes)
- [Grammar](#grammar)
- [Example](#example)
- [Valid Runtime names](#valid-runtime-names)
- [Recognized File extensions](#recognized-file-extensions)
- [Valid Limit keys](#valid-limit-keys)

The Action entity schema contains the necessary information to deploy an OpenWhisk function and define its deployment configurations, inputs and outputs.

### Fields
| Key Name | Required | Value Type | Default | Description |
|:---|:---|:---|:---|:---|
| version | no | [version](spec_parameter_types.md#openwhisk-types) | N/A | The optional user-controlled version for the Action. |
| function | yes | string | N/A | Required source location (path inclusive) of the Action code either:<ul><li>Relative to the Package manifest file.</li><li>Relative to the specified Repository.</li></ul> |
| code | no | string | N/A | This optional field is now replaced by the <em>“function”</em> field. |
| runtime | maybe | string | N/A | he required runtime name (and optional version) that the Action code requires for an execution environment.<p><i>Note: May be optional if tooling is allowed to make assumptions about file extensions or infer from functional code.</i></p> |
| inputs | no | list of [parameter](spec_parameters.md) | N/A | The optional ordered list inputs to the Action. |
| outputs | no | list of [parameter](spec_parameters.md) | N/A | The optional outputs from the Action. |
| limits | no | map of [limit keys and values](#valid-limit-keys) | N/A | Optional map of limit keys and their values.</br>See section "[Valid limit keys](#valid-limit-keys)" (below) for a listing of recognized keys and values. |
| feed | no | boolen | false | Optional indicator that the Action supports the required parameters (and operations) to be run as a Feed Action. |
| web \| web-export | no | boolean | yes \| no \| raw \| false | The optional flag (annotation) that makes the action accessible to REST calls without authentication.<p>For details on all types of Web Actions, see: [Web Actions](https://github.com/apache/incubator-openwhisk/blob/master/docs/webactions.md).</p>|
| raw-http | no | boolean | false | The optional flag (annotation) to indicate if a Web Action is able to consume the raw contents within the body of an HTTP request.<p><b>Note</b>: this option is ONLY valid if <em>"web"</em> or <em>"web-export"</em> is set to <em>‘true’</em>.<p> |
| docker | no | string | N/A | The optional key that references a Docker image (e.g., openwhisk/skeleton). |
| native | no | boolean | false | The optional key (flag) that indicates the Action is should use the Docker skeleton image for OpenWhisk (i.e., short-form for docker: openwhisk/skeleton). |
| final | no | boolean | false | The optional flag (annotation) which makes all of the action parameters that are already defined immutable.<p><b>Note</b>: this option is ONLY valid if <em>"web"</em> or <em>"web-export"</em> is set to <em>‘true’</em>.<p> |
| web-custom-options | no | boolean | false | The optional flag (annotation) enables a web action to respond to OPTIONS requests with customized headers, otherwise a [default CORS response](https://github.com/apache/openwhisk/blob/master/docs/webactions.md#options-requests) applies. |
| require-whisk-auth | no | boolean | false | The optional flag (annotation) protects the web action so that it is only accessible to an authenticated subject. |
| main | no | string | N/A | The optional name of the function to be aliased as a function named “main”.<p><em><b>Note</b>: by convention, Action functions are required to be called “main”; this field allows existing functions not named “main” to be aliased and accessed as if they were named “main”.</em></p>|

### Requirements

- The Action name (i.e., &lt;actionName&gt; MUST be less than or equal to 256 characters.
- The Action entity schema includes all general <a href="#SCHEMA_ENTITY">Entity Schema</a> fields in addition to any fields declared above.
- Supplying a runtime name without a version indicates that OpenWhisk SHOULD use the most current version.
- Supplying a runtime <i>major version</i> without a <i>minor version</i> (et al.) indicates OpenWhisk SHOULD use the most current <i>minor version</i>.
- Unrecognized limit keys (and their values) SHALL be ignored.
- Invalid values for known limit keys SHALL result in an error.
- If the Feed is a Feed Action (i.e., the feed key's value is set to true), it MUST support the following parameters:
  - **lifecycleEvent**: one of 'CREATE', 'DELETE', 'PAUSE',or 'UNPAUSE'. These operation names MAY be supplied in lowercase (i.e., 'create',
'delete', 'pause', etc.).
  - **triggerName**: the fully-qualified name of the trigger which contains events produced from this feed.
  - **authKey**: the Basic auth. credentials of the OpenWhisk user who owns the trigger.
- The keyname ‘kind’ is currently supported as a synonym for the key named ‘runtime’; in the future it MAY be deprecated.
- When a code is specified, runtime SHALL be a required field.


### Notes

- Input and output parameters are implemented as JSON Objects within the OpenWhisk framework.
- The maximum code size for an Action currently must be less than 48 MB.
- The maximum payload size for an Action (i.e., POST content length or size) currently must be less than 1 MB.
- The maximum parameter size for an Action currently must be less than 1 MB.
- if no value for runtime is supplied, the value ‘language:default’ will be assumed.

### Grammar

```yaml
<actionName>[.<type>]:
  <Entity schema>
  version: <version>
  function: <string>
  code: <string>
  runtime: <name>[@<[range of ]version>]
  inputs:
    <list of parameter>
  outputs:
    <list of parameter>
  limits:
    <list of limit key-values>
  feed: <boolean> # default: false
  web | web-export: <boolean> | yes | no | raw
  raw-http: <boolean>
  docker: <string>
  native: <boolean>
  final: <boolean>
  web-custom-options: <boolean>
  require-whisk-auth: <boolean>
  main: <string>
```
_**Note**: the optional [.<type>] grammar is used for naming Web Actions._

### Example

```yaml
my_awesome_action:
  version: 1.0
  description: An awesome action written for node.js
  function: src/js/action.js
  runtime: nodejs@>0.12<6.0
  inputs:
    not_awesome_input_value:
      description: Some input string that is boring
      type: string
  outputs:
    awesome_output_value:
      description: Impressive output string
      type: string
  limits:
    memorySize: 512 kB
    logSize: 5 MB
```

### Valid Runtime names

The following runtime values are currently supported by the OpenWhisk platform.

Each of these runtimes also include additional built-in packages (or libraries) that have been determined be useful for Actions surveyed and tested by the OpenWhisk platform.

These packages may vary by OpenWhisk release; examples of supported runtimes as of this specification version include:

| Runtime value | OpenWhisk kind | Docker image name | Description |
|:---|:---|:---|:---|
| nodejs@10 | nodejs:10 | openwhisk/action-nodejs-v8:latest | Latest NodeJS 10 runtime |
| nodejs@8 | nodejs:8 | openwhisk/action-nodejs-v8:latest | Latest NodeJS 8 runtime |
| nodejs@6 | nodejs:6 | openwhisk/nodejs6action:latest | Latest NodeJS 6 runtime |
| java | java | openwhisk/java8action:latest | Latest Java (8) language runtime |
| php, php@7.3 | php:7.3 | openwhisk/action-php-v7.3:latest | Latest PHP (7.3) language runtime |
| php, php@7.2 | php:7.2 | openwhisk/action-php-v7.2:latest | Latest PHP (7.2) language runtime |
| php, php@7.1 | php:7.1 | openwhisk/action-php-v7.1:latest | Latest PHP (7.1) language runtime |
| python@3 | python:3 | openwhisk/python3action:latest | Latest Python 3 language runtime |
| python, python@2 | python:2 | openwhisk/python2action:latest | Latest Python 2 language runtime |s
| ruby | ruby:2.5 | openwhisk/action-ruby-v2.5:latest | Latest Ruby 2.5 language runtime |
| swift@4.1 | swift | openwhisk/action-swift-v4.1:latest | Latest Swift 4.1 language runtime |
| swift@3.1.1 | swift | openwhisk/action-swift-v3.1.1:latest | Latest Swift 3.1.1 language runtime |
| dotnet, dotnet@2.2 | dotnet:2.2 | openwhisk/action-dotnet-v2.2:latest | Latest .NET Core 2.2 runtime |
| language:default | N/A | N/A | Permit the OpenWhisk platform to select the correct default language runtime. |

#### Notes
- If no value for runtime is supplied, the value 'language:default' will be assumed.

### Recognized File extensions

Although it is best practice to provide a runtime value when declaring an Action, it is not required. In those cases, that a runtime is not provided, the package tooling will attempt to derive the correct runtime based upon the the file extension for the Action's function (source code file). The
following file extensions are recognized and will be run on the latest version of corresponding Runtime listed below:

<html>
<table>
  <tr>
   <th>File extension</th>
   <th>Runtime used</th>
   <th>Description</th>
  </tr>
 <tr>
  <td>.js</td>
  <td>nodejs</td>
  <td>Latest Node.js runtime.</td>
 </tr>
 <tr>
  <td>.java</td>
  <td>java</td>
  <td>Latest Java language runtime.</td>
 </tr>
 <tr>
  <td>.py</td>
  <td>python</td>
  <td>Latest Python language runtime.</td>
 </tr>
 <tr>
  <td>.swift</td>
  <td>swift</td>
  <td>Latest Swift language runtime.</td>
 </tr>
 <tr>
  <td>.php</td>
  <td>php</td>
  <td>Latest PHP language runtime.</td>
 </tr>
 <tr>
  <td>.rb</td>
  <td>ruby</td>
  <td>Latest Ruby language runtime.</td>
 </tr>
</table>
</html>

### Valid Limit keys

<html>
<table id="TABLE_LIMIT_KEYS">
  <tr>
   <th>Limit Keyname</th>
   <th>Allowed values</th>
   <th>Default value</th>
   <th>Valid Range</th>
   <th>Description</th>
  </tr>
 <tr>
  <td>timeout</td>
  <td>scalar-unit.time</td>
  <td>60000 ms</td>
  <td>[100 ms, 300000 ms]</td>
  <td>The per-invocation Action timeout. Default unit is assumed to be milliseconds (ms).</td>
 </tr>
 <tr>
  <td>memorySize</td>
  <td>scalar-unit.size</td>
  <td>256 MB</td>
  <td>[128 MB, 512 MB]</td>
  <td>The per-Action memory. Default unit is assumed to be in megabytes (MB).</p>
  </td>
 </tr>
 <tr>
  <td>logSize</td>
  <td>scalar-unit.size</td>
  <td>10 MB</td>
  <td>[0 MB, 10 MB]</td>
  <td>The action log size. Default unit is assumed to be in megabytes (MB).</td>
 </tr>
 <tr>
  <td>concurrentActivations</td>
  <td>integer</td>
  <td>1000</td>
  <td><i>See description</i></td>
  <td>The maximum number of concurrent Action activations allowed (per-namespace). <p><i>Note: This value is not changeable via APIs at this time.</i></p></td>
 </tr>
 <tr>
  <td>userInvocationRate</td>
  <td>integer</td>
  <td>5000</td>
  <td><i>See description</i></td>
  <td>The maximum number of Action invocations allowed per user, per minute. <p><i>Note: This value is not changeable via APIs at this time.</i></p></td>
 </tr>
 <tr>
  <td>codeSize</td>
  <td>scalar-unit.size</td>
  <td>48 MB</td>
  <td><i>See description</i></td>
  <td>The maximum size of the Action code.<p><i>Note: This value is not changeable via APIs at this
  time.</i></p></td>
 </tr>
 <tr>
  <td>parameterSize</td>
  <td>scalar-unit.size</td>
  <td>1 MB</td>
  <td><i>See description</i></td>
  <td>The maximum size<p><i>Note: This value is not changeable via APIs at this time.</i></p></td>
 </tr>
</table>
</html>

#### Notes

- The default values and ranges for limit configurations reflect the defaults for the OpenWhisk platform (open source code).&nbsp; These values may be changed over time to reflect the open source community consensus.

### Web Actions
OpenWhisk can turn any Action into a 'web action' causing it to return HTTP content without use of an API Gateway. Simply supply a supported 'type' extension to indicate which content type is to be returned and identified in the HTTP header (e.g., _.json_, _.html_, _.text_ or _.http_).

Return values from the Action's function are used to construct the HTTP response. The following response parameters are supported in the response object.

- **headers**: a JSON object where the keys are header-names and the values are string values for those headers (default is no headers).
- **code**: a valid HTTP status code (default is 200 OK).
- **body**: a string which is either plain text or a base64 encoded string (for binary data).

<!--
 Bottom Navigation
-->
---
<html>
<div align="center">
<a href="../README.md#index">Index</a>
</div>
</html>
