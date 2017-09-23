
## Action entity

The Action entity schema contains the necessary information to deploy an OpenWhisk function and define its deployment configurations, inputs and outputs.

### Fields
<html>
<table>
  <tr>
   <th>Key Name</th>
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
- The Action name (i.e., &lt;actionName&gt; MUST be less than or equal to 256 characters.</p>
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
- The maximum code size for an Action currently must be less than 48 MB.
- The maximum payload size for an Action (i.e., POST content length or size) currently must be less than 1 MB.
- The maximum parameter size for an Action currently must be less than 1 MB.

### Grammar
```yaml
# Note: the optional [.<type>] grammar is used for naming Web Actions.
<actionName>[.<type>]:
  <Entity schema> # Common to all OpenWhisk Entities
  function: <string>
  runtime: <name>[@<[range of ]version>]
  limits:
    <list of limit key-values>
  feed: <boolean> # default: false
```

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
   <td>
   <p>Runtime value</p>
   </td>
   <td>
   <p>OpenWhisk kind</p>
   </td>
   <td>
   <p>image name</p>
   </td>
   <td>
   <p>Description</p>
   </td>
  </tr>

 <tr>
  <td>
  <p>nodejs</p>
  </td>
  <td>
  <p>nodejs</p>
  </td>
  <td>
  <p>nodejsaction:latest</p>
  </td>
  <td>
  <p>Latest NodeJS runtime</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>nodejs@6</p>
  </td>
  <td>
  <p>nodejs:6</p>
  </td>
  <td>
  <p>nodejs6action:latest</p>
  </td>
  <td>
  <p>Latest NodeJS 6 runtime</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>java, java@8</p>
  </td>
  <td>
  <p>java</p>
  </td>
  <td>
  <p>java8action:latest</p>
  </td>
  <td>
  <p>Latest Java language runtime</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>python, python@2</p>
  </td>
  <td>
  <p>python:2</p>
  </td>
  <td>
  <p>python2action:latest</p>
  <p>&nbsp;</p>
  </td>
  <td>
  <p>Latest Python 2 language runtime</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>python@3</p>
  </td>
  <td>
  <p>python:3</p>
  </td>
  <td>
  <p>python3action:latest</p>
  </td>
  <td>
  <p>Latest Python 3 language runtime</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>swift, swift@2</p>
  </td>
  <td>
  <p>swift</p>
  </td>
  <td>
  <p>swiftaction:latest</p>
  </td>
  <td>
  <p>Latest Swift 2 language runtime</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>swift@3</p>
  </td>
  <td>
  <p>swift:3</p>
  </td>
  <td>
  <p>swift3action:latest</p>
  </td>
  <td>
  <p>Latest Swift 3 language runtime</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>swift@3.1.1</p>
  </td>
  <td>
  <p>swift:3.1.1</p>
  </td>
  <td>
  <p>action-swift-v3.1.1:latest</p>
  </td>
  <td>
  <p>Latest Swift 3.1.1 language runtime</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>php</p>
  </td>
  <td>
  <p>php:7.1</p>
  </td>
  <td>
  <p>action-php-v7.1:latest</p>
  </td>
  <td>
  <p>Latest PHP language runtime</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>language:default</p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>Permit the OpenWhisk platform to select the correct default
  language runtime.</p>
  </td>
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
   <th>>Runtime used</th>
   <th>Description</th>
  </tr>

 <tr>
  <td>
  <p>.js</p>
  </td>
  <td>
  <p>nodejs</p>
  </td>
  <td>
  <p>Latest Node.js runtime.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>.java</p>
  </td>
  <td>
  <p>java</p>
  </td>
  <td>
  <p>Latest Java language runtime.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>.py</p>
  </td>
  <td>
  <p>python</p>
  </td>
  <td>
  <p>Latest Python language runtime.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>.swift</p>
  </td>
  <td>
  <p>swift</p>
  </td>
  <td>
  <p>Latest Swift language runtime.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>.php</p>
  </td>
  <td>
  <p>php</p>
  </td>
  <td>
  <p>Latest PHP language runtime.</p>
  </td>
 </tr>
</table>
</html>

#### Valid Limit keys

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
  <td>
  <p>timeout</p>
  </td>
  <td>
  <p>scalar-unit.time</p>
  </td>
  <td>
  <p>60000 ms</p>
  </td>
  <td>
  <p>[100 ms, 300000 ms]</p>
  </td>
  <td>
  <p>The per-invocation Action timeout.&nbsp; Default unit is
  assumed to be milliseconds (ms).</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>memorySize</p>
  </td>
  <td>
  <p>scalar-unit.size</p>
  </td>
  <td>
  <p>256 MB</p>
  </td>
  <td>
  <p>[128 MB, 512 MB]</p>
  </td>
  <td>
  <p>The per-Action memory. Default unit is assumed to be in
  megabytes (MB).</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>logSize</p>
  </td>
  <td>
  <p>scalar-unit.size</p>
  </td>
  <td>
  <p>10 MB</p>
  </td>
  <td>
  <p>[0 MB, 10 MB]</p>
  </td>
  <td>
  <p>The action log size. Default unit is assumed to be in
  megabytes (MB).</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>concurrentActivations</p>
  </td>
  <td>
  <p>integer</p>
  </td>
  <td>
  <p>1000</p>
  </td>
  <td>
  <p><i>See description</i></p>
  </td>
  <td>
  <p>The maximum number of concurrent Action activations
  allowed (per-namespace). </p>
  <p>&nbsp;</p>
  <p><i>Note: This value is not changeable via APIs at this
  time.</i></p>
  </td>
 </tr>
 <tr>
  <td>
  <p>userInvocationRate</p>
  </td>
  <td>
  <p>integer</p>
  </td>
  <td>
  <p>5000</p>
  </td>
  <td>
  <p><i>See description</i></p>
  </td>
  <td>
  <p>The maximum number of Action invocations allowed per user,
  per minute.&nbsp;&nbsp;&nbsp; </p>
  <p>&nbsp;</p>
  <p><i>Note: This value is not changeable via APIs at this time.</i></p>
  </td>
 </tr>
 <tr>
  <td>
  <p>codeSize</p>
  </td>
  <td>
  <p>scalar-unit.size</p>
  </td>
  <td>
  <p>48 MB</p>
  </td>
  <td>
  <p><i>See description</i></p>
  </td>
  <td>
  <p>The maximum size of the Action code.</p>
  <p>&nbsp;</p>
  <p><i>Note: This value is not changeable via APIs at this
  time.</i></p>
  </td>
 </tr>
 <tr>
  <td>
  <p>parameterSize</p>
  </td>
  <td>
  <p>scalar-unit.size</p>
  </td>
  <td>
  <p>1 MB</p>
  </td>
  <td>
  <p><i>See description</i></p>
  </td>
  <td>
  <p>The maximum size </p>
  <p>&nbsp;</p>
  <p><i>Note: This value is not changeable via APIs at this
  time.</i></p>
  </td>
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
