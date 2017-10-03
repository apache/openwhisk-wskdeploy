## Actions

The Action entity schema contains the necessary information to deploy an OpenWhisk function and define its deployment configurations, inputs and outputs.

### Fields
<html>
<table>
  <tr>
   <th width="80">Key Name</th>
   <th>Required</th>
   <th>Value Type</th>
   <th>Default</th>
   <th>Description</th>
  </tr>
 <tr>
  <td>version</td>
  <td>no</td>
  <td>version</td>
  <td>N/A</td>
  <td>The optional user-controlled version for the Action.</td>
 </tr>
 <tr>
  <td>function</td>
  <td>yes</td>
  <td>string</td>
  <td>N/A</td>
  <td>Required source location (path inclusive) of the Action code either:
    <ul>
      <li>Relative to the Package manifest file.</li>
      <li>Relative to the specified Repository.</li>
    </ul>
  </td>
 </tr>
 <tr>
  <td>runtime</td>
  <td>no</td>
  <td>string</td>
  <td>N/A</td>
  <td>The required runtime name (and optional version) that the Action code requires for an execution environment.
  <p><i>Note: May be optional if tooling allowed to make assumptions about file extensions.</i></p>
  </td>
 </tr>
 <tr>
  <td>inputs</td>
  <td>no</td>
  <td>list of parameter</td>
  <td>N/A</td>
  <td>The optional ordered list inputs to the Action.</td>
 </tr>
 <tr>
  <td>outputs</td>
  <td>no</td>
  <td>list of parameter</td>
  <td>N/A</td>
  <td>The optional outputs from the Action.</td>
 </tr>
 <tr>
  <td>limits</td>
  <td>no</td>
  <td>map of limit keys and values</a></td>
  <td>N/A</td>
  <td>Optional map of limit keys and their values.
  <p><i>See section '</i><a href="#TABLE_LIMIT_KEYS"><i>Valid limit keys</i></a><i>' below for a listing of recognized keys and values.</i></p>
  </td>
 </tr>
 <tr>
  <td>feed</td>
  <td>no</td>
  <td>boolean</td>
  <td>false</td>
  <td>Optional indicator that the Action supports the required parameters (and operations) to be run as a Feed Action.</td>
 </tr>
 <tr>
  <td>web-export</td>
  <td>no</td>
  <td>boolean</td>
  <td>false</td>
  <td>Optionally, turns the Action into a <a href="https://github.com/apache/incubator-openwhisk/blob/master/docs/webactions.md">&quot;<em><u>web actions</u></em>&quot;</a>
  causing it to return HTTP content without use of an API Gateway.
  </td>
 </tr>
</table>
</html>

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
  runtime: <name>[@<[range of ]version>]
  inputs:
    <list of parameter>
  outputs:
    <list of parameter>
  limits:
    <list of limit key-values>
  feed: <boolean> # default: false
  web-export: <boolean>
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

<html>
<table>

  <tr>
   <th>Runtime value</th>
   <th>OpenWhisk kind</th>
   <th>image name</th>
   <th>Description</th>
  </tr>

 <tr>
  <td>nodejs</td>
  <td>nodejs</td>
  <td>nodejsaction:latest</td>
  <td>Latest NodeJS runtime</td>
 </tr>
 <tr>
  <td>nodejs@6</td>
  <td>nodejs:6</td>
  <td>nodejs6action:latest</td>
  <td>Latest NodeJS 6 runtime</td>
 </tr>
 <tr>
  <td>java, java@8</td>
  <td>java</td>
  <td>java8action:latest</td>
  <td>Latest Java language runtime</td>
 </tr>
 <tr>
  <td>python, python@2</td>
  <td>python:2</td>
  <td>python2action:latest</td>
  <td>Latest Python 2 language runtime</td>
 </tr>
 <tr>
  <td>python@3</td>
  <td>python:3</td>
  <td>python3action:latest</td>
  <td>Latest Python 3 language runtime</td>
 </tr>
 <tr>
  <td>swift, swift@2</td>
  <td>swift</td>
  <td>swiftaction:latest</td>
  <td>Latest Swift 2 language runtime</td>
 </tr>
 <tr>
  <td>swift@3</td>
  <td>swift:3</td>
  <td>swift3action:latest</td>
  <td>Latest Swift 3 language runtime</td>
 </tr>
 <tr>
  <td>swift@3.1.1</td>
  <td>swift:3.1.1</td>
  <td>action-swift-v3.1.1:latest</td>
  <td>Latest Swift 3.1.1 language runtime</td>
 </tr>
 <tr>
  <td>php</td>
  <td>php:7.1</td>
  <td>action-php-v7.1:latest</td>
  <td>Latest PHP language runtime</td>
 </tr>
 <tr>
  <td>language:default</td>
  <td>N/A</td>
  <td>N/A</td>
  <td>Permit the OpenWhisk platform to select the correct default language runtime.</td>
 </tr>
</table>
</html>

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
<table align="center">
  <tr>
    <!-- <td><a href="">&lt;&lt;&nbsp;previous</a></td> -->
    <td><a href="spec_index.md#openwhisk-package-specification-html">Specification Index</a></td>
    <!-- <td><a href="">next&nbsp;&gt;&gt;</a></td> -->
  </tr>
</table>
</div>
</html>
