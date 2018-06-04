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

## Sequences

Actions can be composed into sequences to, in effect, form a new Action. The Sequence entity allows for a simple, convenient way to describe them in the Package Manifest.

### Fields
<html>
<table>
  <tr>
   <th>Key Name</th>
   <th>Required</th>
   <th>Value Type</th>
   <th>Default</th>
   <td>Description</th>
  </tr>
 <tr>
  <td>
  <p>actions</p>
  </td>
  <td>
  <p>yes</p>
  </td>
  <td>
  <p>list of Action</p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>The required list of two or more actions</p>
  </td>
 </tr>
</table>
</html>

### Requirements

- The comma separated list of Actions on the actions key SHALL imply the order of the sequence (from left, to right).
- There MUST be two (2) or more actions declared in the sequence.

### Notes

- The sequences key exists for convenience; however, it is just one possible instance of a composition of Actions. The composition entity is provided for not only describing sequences, but also for other (future) compositions and additional information needed to compose them.&nbsp; For example, the composition entity allows for more complex mappings of input and output parameters between Actions.

### Grammar

```yaml
sequences:
  <sequence name>:
     <Entity schema>
     actions: <ordered list of action names>
  ...
```

### Example
```yaml
sequences:
  newbot:
    actions: newbot-create, newbot-select-persona, newbot-greeting
```

<!--
 Bottom Navigation
-->
---
<html>
<div align="center">
<a href="../README.md#index">Index</a>
</div>
</html>
