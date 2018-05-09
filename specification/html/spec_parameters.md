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

# Parameters

- [Dollar Notation ($)](#dollar-notation--schema-for-values)

## Parameter schema
The Parameter schema is used to define input and/or output data to be used by OpenWhisk entities for the purposes of validation.

### Fields
<html>
<table width="100%">
  <tr>
   <th>Key Name</th>
   <th>Required</th>
   <th>Value Type</th>
   <th>Default</th>
   <th>Description</th>
  </tr>
 <tr>
  <td>
  <p>type</p>
  </td>
  <td>
  <p>no</p>
  </td>
  <td>
  <p>&lt;any&gt;</p>
  </td>
  <td>
  <p>string</p>
  </td>
  <td>
  <p>Optional valid type name or the parameter's value for validation purposes. By default, the type is string.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>description</p>
  </td>
  <td>
  <p>no</p>
  </td>
  <td>
  <p>string256</p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>Optional description of the Parameter.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>value</p>
  </td>
  <td>
  <p>no </p>
  </td>
  <td>
  <p>&lt;any&gt;</p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>The optional user supplied value for the parameter.</p>
  <p>Note: this is not the default value, but an explicit declaration which allows simple usage of the Manifest file without a Deployment file.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>required</p>
  </td>
  <td>
  <p>no</p>
  </td>
  <td>
  <p>boolean</p>
  </td>
  <td>
  <p>true</p>
  </td>
  <td>
  <p>Optional indicator to declare the parameter as required (i.e., ```true```) or optional (i.e., ```false```).</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>default</p>
  </td>
  <td>
  <p>no</p>
  </td>
  <td>
  <p>&lt;any&gt;</p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>Optional default value for the optional parameters. This value <b>MUST</b> be type compatible with the value declared on the parameter's type field.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>status</p>
  </td>
  <td>
  <p>no</p>
  </td>
  <td>
  <p>string</p>
  </td>
  <td>
  <p>supported</p>
  </td>
  <td>
  <p>Optional status of the parameter (e.g., deprecated, experimental). By default a parameter is without a declared status is considered supported.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>schema</p>
  </td>
  <td>
  <p>no</p>
  </td>
  <td>
  <p>&lt;schema&gt;</p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>The optional schema if the 'type' key has the value 'schema'. The value would include a <b>Schema</b> <b>Object</b> (in YAML) as defined by the <a href="See%20https:/github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#schemaObject">OpenAPI Specifcation v2.0</a>. This object is based upon the href="http://json-schema.org/">JSON Schema Specification.</a></p>
  </td>
 </tr>
 <tr>
  <td>
  <p>properties</p>
  </td>
  <td>
  <p>no</p>
  </td>
  <td>
  <p>&lt;list of parameter&gt;</p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>The optional properties if the 'type' key has the value 'object'. Its value is a listing of Parameter schema from this specification.</p>
  </td>
 </tr>
</table>
</html>

### Requirements

The 'schema' key's value MUST be compatible with the value provided on both the 'type'  and 'value' keys; otherwise, it is considered an error.

### Notes

The 'type' key acknowledges some popular schema (e.g., JSON) to use when validating the value of the parameter. In the future additional (schema) types may be added for convenience.

### Grammar

#### Single-line
```yaml
```

Where <YAML type> is inferred to be a YAML type as shown in the YAML Types section above (e.g., string, integer, float, boolean, etc.).

If you wish the parser to validate against a different schema, then the multi-line grammar MUST be used where the value would be supplied on the keyname 'value' and the type (e.g., 'json') and/or schema (e.g., OpenAPI) can be supplied.

#### Multi-line
```yaml
```

### Status values

<table width="100%">
  <tr>
   <td>
   <p>Status Value</p>
   </td>
   <td>
   <p>Description</p>
   </td>
  </tr>

 <tr>
  <td>
  <p>supported (default)</p>
  </td>
  <td>
  <p>Indicates the parameter is supported.&nbsp; This is the
  implied default status value for all parameters.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>experimental</p>
  </td>
  <td>
  <p>Indicates the parameter MAY be removed or changed in
  future versions.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>deprecated</p>
  </td>
  <td>
  <p>Indicates the parameter is no longer supported in the
  current version and MAY be ignored.</p>
  </td>
 </tr>
</table>

## Dollar Notation ($) schema for values

In a Manifest or Deployment file, a parameter value may be set from the local execution environment by using the dollar ($) notation to denote names of local environment variables which supply the value to be inserted at execution time.

### Syntax
```yaml
<parameter>: $<local environment variable name>
```

### Example
```yaml
...
  inputs:
    userName: $DEFAULT_USERNAME
```

### Requirements

- Processors or tooling that encounter ($) Dollar notation and are unable to locate the value in the execution environment SHOULD resolve the value to be the default value for the type (e.g., an empty string ("") for type 'string').

### Notes

- Processors or tooling that encounter ($) Dollar notation for values should attempt to locate the corresponding named variables set into the local execution environment (e.g., where the tool was invoked) and assign its value to the named input parameter for the OpenWhisk entity.
- This specification does not currently consider using this notation for other than simple data types (i.e., we support this mechanism for values such as strings, integers, floats, etc.) at this time.


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
