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

- [Parameter schema](#parameter-schema)
- [Dollar Notation ($)](#dollar-notation--schema-for-values)

## Parameter schema
The Parameter schema is used to define input and/or output data to be used by OpenWhisk entities for the purposes of validation.

### Fields
| Key Name | Required | Value Type | Default | Description |
|:---|:---|:---|:---|:---|
| type       | no | _&lt;any&gt;_ | string | Optional valid type name or the parameter’s value for alidation purposes. By default, the type is string. |
| description | no | string256 | N/A | Optional description of the Parameter. |
| value      | no | _&lt;any&gt;_ | N/A | The optional user supplied value for the parameter.</br> Note: this is not the default value, but an explicit declaration which allows simple usage of the Manifest file without a Deployment file. |
| required   | no | boolean    | true | Optional indicator to declare the parameter as required (i.e., true) or optional (i.e., false). |
|  default   | no | _&lt;any&gt;_ | N/A | Optional default value for the optional parameters. This value **MUST** be type compatible with the value declared on the parameter’s type field. |
| status     | no | string     | supported | Optional status of the parameter (e.g., deprecated, experimental). By default a parameter is without a declared status is considered supported. |
| schema     | no | _&lt;schema&gt;_ | N/A | The optional schema if the ‘type‘ key has the value ‘schema‘. The value would include a **Schema** **Object** (in YAML) as defined by the [OpenAPI Specification v2.0](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#schemaObject). This object is based upon the [JSON Schema Specification.](http://json-schema.org/) |
| properties | no |  _&lt;list of parameter&gt;_ | N/A | The optional properties if the ‘type‘ key has the value ‘object‘. Its value is a listing of Parameter schema from this specification. |

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

| Status Value | Description |
|:---|:---|
| supported (default) | Indicates the parameter is supported. This is the implied default status value for all parameters. |
| experimental | Indicates the parameter MAY be removed or changed in future versions. |
| deprecated | Indicates the parameter is no longer supported in the current version and MAY be ignored. |

#### Shared Entity Schema

The Entity Schema contains fields that are common (shared) to all
OpenWhisk entities (e.g., Actions, Triggers, Rules, etc.).

##### Fields

|  Key Name | Required | Value Type | Default | Description
|:---|:---|:---|:---|:---|
| description | no | string256 | N/A | The optional description for the Entity. |
| displayName | no | string16 | N/A | This is the optional name that will be displayed on small form-factor devices. |
| annotations | no | map of _&lt;string>&gt;_ | N/A | he optional annotations for the Entity. |

##### Grammar

```yaml
description: <string256>
displayName: <string16>
annotations: <map of <string>>
```

##### Requirements

- Non-required fields MAY be stored as “annotations” within the OpenWhisk framework after they have been used for processing.
- "description" string values SHALL be limited to 256 characters.
- "displayName" string values SHALL be limited to 16 characters.
- "annotations" MAY be ignored by target consumers of the Manifest file as they are considered data non-essential to the deployment of management of OpenWhisk entities themselves.
    - Target consumers MAY preserve (persist) these values, but are not required to.
- For any OpenWhisk Entity, the maximum size of all "annotations" (values) SHALL be 256 characters.

##### Notes

- Several, non-normative "annotations" keynames and allowed values (principally for User Interface (UI) design and tooling information) may be defined in this specification or optional usage.


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
- This specification does not currently consider using this notation for other than simple data types (i.e., we support this mechanism for values such as `strings`, `integers`, `floats`, etc.) at this time.


<!--
 Bottom Navigation
-->
---
<html>
<div align="center">
<a href="../README.md#index">Index</a>
</div>
</html>
