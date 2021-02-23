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

## Rules

The Rule entity schema contains the information necessary to associates one trigger with one action, with every firing of the trigger causing the corresponding action to be invoked with the trigger event as input. For more information, see the document "[Creating Triggers and Rules](https://github.com/apache/openwhisk/blob/master/docs/triggers_rules.md)".

#### Subsections

- [Fields](#fields)
- [Requirements](#requirements)
- [Notes](#notes)
- [Grammar](#grammar)
- [Example](#example)

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
  <p>trigger</p>
  </td>
  <td>
  <p>yes</p>
  </td>
  <td>
  <p>string</p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>Required name of the Trigger the Rule applies to.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>action</p>
  </td>
  <td>
  <p>yes</p>
  </td>
  <td>
  <p>string</p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>Required name of the Action the Rule applies to.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>rule</p>
  </td>
  <td>
  <p>no</p>
  </td>
  <td>
  <p>regex</p>
  </td>
  <td>
  <p>true</p>
  </td>
  <td>
  <p>The optional regular expression that determines if the Action is fired.</p>
  <p><i>Note: In this version of the specification, only the expression 'true' is currently supported.</i></p>
  </td>
 </tr>
</table>
</html>

### Requirements
- The Rule name (i.e., <ruleName>) MUST be less than or equal to 256 characters.
- The Rule entity schema includes all general [Entity Schem](#TBD) fields in addition to any fields
declared above.

### Notes
- OpenWhisk only supports a value of '```true```' for the '```rule```' key's value at this time.

### Grammar
```yaml
<ruleName>:
  description: <string>
  trigger: <string>
  action: <string>
  rule: <regex>
```

### Example

```yaml
my_rule:
  description: Enable events for my Action
  trigger: my_trigger
  action: my_action
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
