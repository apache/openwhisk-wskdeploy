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

# Parameter Types

- [YAML Types](#yaml-types)
- [OpenWhisk Types](#openwhisk-types)
  - [scalar-unit Types](#scalar-unit-types)

---

#### YAML Types

Many of the types we use in this profile are *built-in* types from the [YAML 1.2 specification](http://www.yaml.org/spec/1.2/spec.html) (i.e., those identified by the “tag:yaml.org,2002” version tag).

The following table declares the valid YAML type URIs and aliases that SHALL be used when defining parameters or properties within an OpenWhisk package manifest:

| Type Name | Type URI | Notes |
|:---|:---|:---|
| string    | tag:yaml.org,2002:str (default) | Default type if no type provided |
| integer   | tag:yaml.org,2002:int | Signed. Includes large integers (i.e., long type) |
| float     | tag:yaml.org,2002:float | Signed. Includes large floating point values (i.e., double type) |
| boolean   | tag:yaml.org,2002:bool | This specification uses lowercase ‘true’ and lowercase ‘false’ |
| timestamp | tag:yaml.org,2002:timestamp | ISO 8601 compatible timestamp. See [YAML-TS-1.1](spec_normative_refs.md) |
| null      | tag:yaml.org,2002:null | Different meaning than an empty string, map, list, etc. |

#### Requirements

- The ‘string’ type SHALL be the default type when not specified on a parameter or property declaration.
- All ‘boolean’ values SHALL be lowercased (i.e., ‘true’ or ‘false’).
- The case (lower, upper) of the Unit’s value is considered significant and SHALL be preserved.  For example, the Unit value “kB” is not the same as “KB” or “kb”.

### OpenWhisk Types

In addition to the YAML built-in types, OpenWhisk supports the types listed in the table below. A complete description of each of these types is provided below.

| Type Name | Description | Notes |
|:---|:---|:---|
| version | tag:maven.apache.org:version | Typically found in modern tooling (i.e., 'package@version' or 'package:version' format). Aligns with Maven format principles, but is a simplification of Maven spec. considerations. See [Maven version](spec_normative_refs.md#normative-references"). |
| string256 | long length strings (e.g., descriptions) | A string type limited to 256 characters. |
| string64  | medium length strings (e.g., abstracts, hover text) | A string type limited to 64 characters. |
| string16  | short length strings (e.g., small form-factor list displays) | A string type limited to 16 characters. |
| json      | The parameter value represents a JavaScript Object Notation (JSON) data object. | The deploy tool will validate the corresponding parameter value against JSON schema.</br></br>Note: The implied schema for JSON the JSON Schema (see http://json-schema.org/). |
| scalar-unit | Convenience type for declaring common scalars that have an associated unit. For example, “10 msec.”, “2 Gb”, etc.) | Currently, the following scalar-unit subtypes are supported: ```scalar-unit.size```, ```scalar-unit.time```. |
| schema    | The parameter itself is an OpenAPI Specification v2.0 **Schema Object** (in YAML format) with self-defining schema. | The schema declaration follows the [OpenAPI](spec_normative_refs.md#normative-references") v2.0 specification for Schema Objects (YAML format). |
| object   | The parameter itself is an object with the associated defined Parameters (schemas). |  Parameters of this type would include a declaration of its constituting Parameter schema. |

#### scalar-unit types

Scalar-unit types can be used to define scalar values along with a unit
from the list of recognized units (a subset of GNU units) provided
below.

##### Grammar

```
<scalar> <unit>
```

In the above grammar, the pseudo values that appear in angle brackets
have the following meaning:
-   scalar: is a *required* scalar value (e.g., integer).
-   unit: is a *required* unit value. The unit value MUST be
    type-compatible with the scalar value.

##### Example

```yaml
inputs:
  max_storage_size:
    type: scalar-unit.size
    default: 10 GB
  archive_period:
    type: scalar-unit.time
    default: 30 d
```

##### Requirements
-   Whitespace: any number of spaces (including zero or none) SHALL be allowed between the scalar value and the unit value.
-   It SHALL be considered an error if either the scalar or unit portion is missing on a property or attribute declaration derived from any scalar-unit type.

##### Recognized units for sizes (i.e., scalar-unit.size)


| Unit | Description |
|:---|:---|
| B   | byte |
| kB  | kilobyte (1000 bytes) |
| MB  | megabyte (1000000 bytes) |
| GB  | gigabyte (1000000000 bytes) |
| TB  | terabyte (1000000000000 bytes) |

##### Example

```yaml
inputs:
  memory_size:
    type: scalar-unit.size
    value: 256 MB
```

##### Recognized units for times (i.e., scalar-unit.time)

| Unit | Description |
|:---|:---|
| d  | days |
| h  | hours |
| m  | minutes |
| s  | seconds |
| ms | milliseconds |
| us | microseconds |

##### Example

```yaml
inputs:
  max_execution|time:
    type: scalar-unit.time
    value: 600 s
```

#### Object type example

The Object type allows for complex objects to be declared as parameters
with an optional validateable schema.

```yaml
inputs:
  person:
    type: object
    parameters:
      <Parameter schema>
```

<!--
 Bottom Navigation
-->
---
[Index](../README.md#index)
