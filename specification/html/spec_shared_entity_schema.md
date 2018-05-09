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

## Shared Entity Schema

The Entity Schema contains fields that are common (shared) to all OpenWhisk entities (e.g., Actions, Triggers, Rules, etc.).

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
  <p>The optional description for the Entity.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>displayName</p>
  </td>
  <td>
  <p>no</p>
  </td>
  <td>
  <p>string16</p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>This is the optional name that will be displayed on small form-factor devices.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>annotations</p>
  </td>
  <td>
  <p>no</p>
  </td>
  <td>
  <p>map of
  <string></p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>The optional annotations for the Entity.</p>
  </td>
 </tr>
</table>
</html>

### Grammar
```yaml
description: <string256>
displayName: <string16>
annotations: <map of <string>>
```

### Requirements
- Non-required fields MAY be stored as ÒannotationsÓ within the OpenWhisk framework after they have been used for processing.
- Description string values SHALL be limited to 256 characters.
- DisplayName string values SHALL be limited to 16 characters.
- Annotations MAY be ignored by target consumers of the Manifest file as they are considered data non-essential to the deployment of management of OpenWhisk entities themselves.
- Target consumers MAY preserve (persist) these values, but are not required to.
- For any OpenWhisk Entity, the maximum size of all Annotations SHALL be 256 characters.

### Notes
- Several, non-normative Annotation keynames and allowed values for (principally for User Interface (UI) design) may be defined below for optional usage.

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
