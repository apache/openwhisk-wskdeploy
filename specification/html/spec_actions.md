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
| feed | no | boolean | false | Optional indicator that the Action supports the required parameters (and operations) to be run as a Feed Action. |
| web | no | string | true&nbsp;&#124; false&nbsp;&#124; yes&nbsp;&#124; no&nbsp;&#124; raw | The optional flag that makes the action accessible to REST calls without authentication.<p>For details on all types of Web Actions, see: [Web Actions](https://github.com/apache/openwhisk/blob/master/docs/webactions.md).</p>|
| raw-http | no | boolean | false | The optional flag (annotation) to indicate if a Web Action is able to consume the raw contents within the body of an HTTP request.<p><b>Note</b>: this option is ONLY valid if <em>"web"</em> or <em>"web-export"</em> is set to <em>‘true’</em>.<p> |
| docker | no | string | N/A | The optional key that references a Docker image (e.g., openwhisk/skeleton). |
| native | no | boolean | false | The optional key (flag) that indicates the Action is should use the Docker skeleton image for OpenWhisk (i.e., short-form for docker: openwhisk/skeleton). |
| final | no | boolean | false | The optional flag (annotation) which makes all of the action parameters that are already defined immutable.<p><b>Note</b>: this option is ONLY valid if <em>"web"</em> or <em>"web-export"</em> is set to <em>‘true’</em>.<p> |
| main | no | string | N/A | The optional name of the function to be aliased as a function named “main”.<p><em><b>Note</b>: by convention, Action functions are required to be called “main”; this field allows existing functions not named “main” to be aliased and accessed as if they were named “main”.</em></p>|
| annotations | no | N/A | The optional map of annotation key-values. See below for [Action annotations](#action-annotations) on actions. |

#### Action Annotations

The following annotations have special meanings for Actions:

| Key Name | Required | Value Type | Default | Description |
|:---|:---|:---|:---|:---|
| final              | no | boolean | not set (false) | Parameters are protected and treated as immutable. This is required for web actions (i.e., `web` or `web-export` set to `true`. |
| web-export         | no | boolean&nbsp;&#124; yes&nbsp;&#124; no&nbsp;&#124; raw  | not set (false) | The optional annotation used to export an action as a `web action` which is accessible through an API REST interface (url). |
| web-custom-options | no | boolean | not set (false) | The optional annotation that enables a web action to respond to OPTIONS requests with customized headers, otherwise a [default CORS response](https://github.com/apache/openwhisk/blob/master/docs/webactions.md#options-requests) applies. |
| require-whisk-auth | no | string&nbsp;&#124; integer&nbsp;&#124; boolean | not set (false) | The optional annotation that can secure a `web action` so that it is only accessible to an authenticated subject.<p>See [Securing web actions](https://github.com/apache/openwhisk/blob/master/docs/webactions.md#Securing-web-actions)</p> |

### Requirements

- The Action entity schema **SHALL** include all general <a href="#SCHEMA_ENTITY">Entity Schema</a> fields in addition to any fields declared above.
- The Action name (i.e., &lt;actionName&gt; **MUST** be less than or equal to 256 characters.
- Supplying a runtime name without a version indicates that OpenWhisk **SHOULD** use the current default version.
- Supplying a runtime <i>major version</i> without a <i>minor version</i> (et al.) indicates OpenWhisk **SHOULD** use the most current <i>minor version</i>.
- Unrecognized limit keys (and their values) **SHALL** be ignored.
- Invalid values for known limit keys **SHALL** result in an error.
- If the Feed is a Feed Action (i.e., the feed key's value is set to true), it **MUST** support the following parameters:
  - **lifecycleEvent**: one of `CREATE`, `DELETE`, `PAUSE`,or `UNPAUSE`. These operation names **MAY** be supplied in lowercase (i.e., `create`,
`delete`, `pause`, etc.).
  - **triggerName**: the fully-qualified name of the trigger which contains events produced from this feed.
  - **authKey**: the Basic auth. credentials of the OpenWhisk user who owns the trigger.
- The keyname `kind` is currently supported as a synonym for the key named ‘`runtime`’; in the future it **MAY** be deprecated.
- When the `code` key-value is specified, the `runtime` **SHALL** be a required field.

#### Annotation requirements
- The annotation `require-whisk-auth` **SHALL** only be valid for web actions (i.e., if the `web` key or `web-export` annotation is set to `true`).
- If the value of the `require-whisk-auth` annotation is an `integer` its value **MUST** be a positive integer less than or equal to the `MAX_INT` value of `9007199254740991`.
- When the `web` or `web-export` key is present and set to `true` the web action's **MUST** also be marked `final`.  This happens automatically when the `web` or `web-export` keys are present and set to `true`.

### Notes

- Input and output parameters are implemented as JSON Objects within the CLI client framework.
- The maximum code size for an Action currently must be less than 48 MB.
- The maximum payload size for an Action (i.e., POST content length or size) currently must be less than 1 MB.
- The maximum parameter size for an Action currently must be less than 1 MB.
- if no value for runtime is supplied, the value `language:default` will be assumed.

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
  web: <boolean> | yes | no | raw
  raw-http: <boolean>
  docker: <string>
  native: <boolean>
  final: <boolean>
  main: <string>
  annotations:
    <map of annotation key-values>
    web-export: <boolean> | yes | no | raw # optional
    web-custom-options: <boolean> # optional, only valid when `web-export` enabled
    require-whisk-auth: <boolean> | <string> | <positive integer> # optional, only valid when `web-export` enabled
```
_**Note**: the optional [.<type>] grammar is used for naming Web Actions._

### Example

```yaml
my_awesome_action:
  version: 1.0
  description: An awesome action written for node.js
  function: src/js/action.js
  runtime: nodejs@>0.12<6.0
  web: true
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
  annotations:
    require-whisk-auth: "my-auth-token"
```

### Valid Runtime names

The following runtime values are currently supported by the OpenWhisk platform "out-of-box" at around the time of the Openwhisk platform release 1.0.

| Runtime value | OpenWhisk kind | Docker image | Tag | Description |
|:---|:---|:---|:---|:---|
| go&nbsp;&#124; go:1.11 (default)| go:1.11 | [openwhisk/action-golang-v1.11](https://hub.docker.com/r/openwhisk/action-golang-v1.11) | nightly | Go 1.11 runtime |
| nodejs@12 | nodejs:12 | [openwhisk/nodejs12action](https://hub.docker.com/r/openwhisk/action-nodejs-v12) | nightly | NodeJS 12 runtime |
| nodejs&nbsp;&#124; nodejs@10 (default)| nodejs:10 | [openwhisk/action-nodejs-v10](https://hub.docker.com/r/openwhisk/action-nodejs-v10) | nightly |NodeJS 10 runtime |
| nodejs@8 | nodejs:8 | [openwhisk/action-nodejs-v8](https://hub.docker.com/r/openwhisk/action-nodejs-v8) | nightly | NodeJS 8 runtime |
| nodejs@6 **(deprecated)**| nodejs:6 | [openwhisk/nodejs6action](https://hub.docker.com/r/openwhisk/nodejs6action) | nightly | NodeJS 6 runtime |
| java&nbsp;&#124; java8 (default) | java:8 | [openwhisk/java8action](https://hub.docker.com/r/openwhisk/java8action) | nightly | Java (8) language runtime |
| php&nbsp;&#124; php@7.4 (default) | php:7.4 | [openwhisk/action-php-v7.4](https://hub.docker.com/r/openwhisk/action-php-v7.4) | nightly | PHP (7.3) language runtime |
| php@7.3 | php:7.3 | [openwhisk/action-php-v7.3](https://hub.docker.com/r/openwhisk/action-php-v7.3) | nightly | PHP (7.3) language runtime |
| php@7.2 **(deprecated)** | php:7.2 | [openwhisk/action-php-v7.2](https://hub.docker.com/r/openwhisk/action-php-v7.2) | nightly | PHP (7.2) language runtime |
| php@7.1 **(deprecated)** | php:7.1 | [openwhisk/action-php-v7.1](https://hub.docker.com/r/openwhisk/action-php-v7.1) | nightly | PHP (7.1) language runtime |
| python&nbsp;&#124; python@3 (default) | python:3 | [openwhisk/python3action](https://hub.docker.com/r/openwhisk/python3action) | nightly | Python 3 (3.6) language runtime |
| python@2 | python:2 | [openwhisk/python2action](https://hub.docker.com/r/openwhisk/python2action) | 1.13.0-incubating | Python 2 (2.7) language runtime |
| ruby&nbsp;&#124; (default) | ruby:2.5 | [openwhisk/action-ruby-v2.5](https://hub.docker.com/repository/docker/openwhisk/action-ruby-v2.5) | nightly | Ruby 2.5 language runtime |
| swift&nbsp;&#124; swift@4.2 (default) | swift:4.2 | [openwhisk/action-swift-v4.2](https://hub.docker.com/r/openwhisk/action-swift-v4.2) | nightly | Swift 4.2 language runtime |
| swift@4.1 | swift:4.1 | [openwhisk/action-swift-v4.1](https://hub.docker.com/r/openwhisk/action-swift-v4.1) | nightly | Swift 4.1 language runtime |
| swift@3.1.1 **(deprecated)** | swift:3.1.1 | [openwhisk/action-swift-v3.1.1](https://hub.docker.com/r/openwhisk/action-swift-v3.1.1) | nightly | Swift 3.1.1 language runtime |
| dotnet&nbsp;&#124; dotnet@2.2 (default) | dotnet:2.2 | [openwhisk/action-dotnet-v2.2](https://hub.docker.com/r/openwhisk/action-dotnet-v2.2) | nightly | .NET Core 2.2 runtime |
| dotnet@3.1 | dotnet:3.1 | [openwhisk/action-dotnet-v3.1](https://hub.docker.com/r/openwhisk/action-dotnet-v3.1) | nightly | .NET Core 3.1 runtime |
| language:default | N/A | N/A | N/A | Permit the OpenWhisk platform to select the correct default language runtime. |

See the file [runtimes.json](https://github.com/apache/openwhisk/blob/master/ansible/files/runtimes.json) in
the main [apache/openwhisk](https://github.com/apache/openwhisk) repository for the latest supported runtimes nad versions.

#### Notes
- **WARNING**: _For OpenWhisk project builds, the Docker image used is tagged `nightly` in Docker Hub (e.g, for GitHub pull
requests). Production uses of OpenWhisk code may use different images and tagged (released) image versions._
- If no value for `runtime` is supplied, the value `language:default` will be assumed.
- OpenWhisk runtimes may also include additional built-in packages (or libraries) that have been determined be useful for Actions surveyed and tested by the OpenWhisk platform.


### Recognized File extensions

Although it is best practice to provide a runtime value when declaring an Action, it is not required. In those cases, that a runtime is not provided, the package tooling will attempt to derive the correct runtime based upon the the file extension for the Action's function (source code file). The
following file extensions are recognized and will be run on the version of corresponding Runtime listed below:

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
  <th>Type</th>
  <th>Default value <sup><a href="#limit-notes">1</a></sup></th>
  <th>Default Range  <sup><a href="#limit-notes">2</a></sup></th>
  <th>Description</th>
</tr>
<tr>
  <td>codeSize</td>
  <td>scalar-unit.size</td>
  <td>48 MB</td>
  <td>[1, 48] MB<sup><a href="#limit-notes">3</a></sup></td>
  <td>The maximum size of the Action code.</td>
</tr>
<tr>
  <td>concurrentActivations</td>
  <td>integer</td>
  <td>1000</td>
  <td>[1, 1000] <sup><a href="#limit-notes">3</a></sup></td>
  <td>The maximum number of concurrent Action activations allowed (per-namespace).</td>
</tr>
<tr>
  <td>logSize</td>
  <td>scalar-unit.size</td>
  <td>10 MB</td>
  <td>[0, 10] MB</td>
  <td>The action log size. Default unit is assumed to be in megabytes (MB).</td>
</tr>
<tr>
  <td>memorySize</td>
  <td>scalar-unit.size</td>
  <td>256 MB</td>
  <td>[128, 2048] MB</td>
  <td>The per-Action memory. Default unit is assumed to be in megabytes (MB).</p>
  </td>
</tr>
<tr>
  <td>parameterSize</td>
  <td>scalar-unit.size</td>
  <td>5 MB</td>
  <td>[0, 5] MB <sup><a href="#limit-notes">3, 4</a></sup></td>
  <td>The maximum size of all parameters (total) for an Action.</td>
</tr>
<tr>
  <td>timeout</td>
  <td>scalar-unit.time</td>
  <td>60000 ms</td>
  <td>[100, 600000] ms</td>
  <td>The per-invocation Action timeout. Default unit is assumed to be milliseconds (ms).</td>
</tr>
<tr>
  <td>userInvocationRate</td>
  <td>integer</td>
  <td>5000</td>
  <td>[1, 5000] <sup><a href="#limit-notes">3</a></sup></td>
  <td>The maximum number of Action invocations allowed per user, per minute.</td>
</tr>
</table>
</html>

#### Limit Notes

1. The default values and ranges for limit configurations reflect the defaults for the OpenWhisk platform (open source code).&nbsp; These values may be changed over time to reflect the open source community consensus.
2. Serverless providers that use Apache OpenWhisk MAY choose to enforce different defaults and value ranges for limits.
3. This limit is not currently user configurable.
4. The parameter size limit also applies to Triggers and Packages.

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
