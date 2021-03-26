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

## Triggers

The Trigger entity schema contains the necessary information to describe the stream of events that it represents. For more information, see the document "[Creating Triggers and Rules](https://github.com/apache/openwhisk/blob/master/docs/triggers_rules.md)".

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
  <p>feed</p>
  </td>
  <td>
  <p>no</p>
  </td>
  <td>
  <p>string</p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>The optional name of the Feed associated with the Trigger.
  </p>
  </td>
 </tr>
 <tr>
  <td>
  <p>credential</p>
  </td>
  <td>
  <p>no</p>
  </td>
  <td>
  <p>Credential</p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>The optional credential used to access the feed service.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>inputs</p>
  </td>
  <td>
  <p>no</p>
  </td>
  <td>
  <p>list of</p>
  <p>parameter</p>
  </td>
  <td>
  <p>N/A</p>
  </td>
  <td>
  <p>The optional ordered list inputs to the feed.</p>
  </td>
 </tr>
 <tr>
  <td>
  <p>events</p>
  <p><i>&nbsp;</i></p>
  </td>
  <td>
  <p>no</p>
  </td>
  <td>list of Event</td>
  <td>
  <p>N/A</p>
  </td>
  <td>The optional list of valid Event schema the trigger supports. OpenWhisk would validate incoming Event data for conformance against any Event schema declared under this key.
  <p><b><i>Note</i></b><i>: This feature is <u>not supported at
  this time</u>. This is viewed as a possible feature that may be
  implemented along with configurable options for handling of invalid events.</i></p></td>
 </tr>
</table>
</html>

### Requirements

The Trigger name (i.e., <triggerName> MUST be less than or equal to 256 characters.

The Trigger entity schema includes all general [Entity Schema](#TBD) fields in addition to any fields
declared above.

### Notes

- The 'events' key name is not supported at this time.</p>
- The Trigger entity within the OpenWhisk programming model is considered outside the scope of the Package (although there are discussions about changing this in the future). This means that Trigger and API information will not be returned when using the OpenWhisk Package API:
  -  ```$ wsk package list <package name>```
- However, it may be obtained using the Trigger API:
  - ```$ wsk trigger list -v

### Grammar

```yaml
<triggerName>:
  <Entity schema>
  feed: <feed name>
  credential: <Credential>
  inputs:
    <list of parameter>
```

### Example

```yaml
triggers:
  everyhour:
    feed: /whisk.system/alarms/alarm
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
