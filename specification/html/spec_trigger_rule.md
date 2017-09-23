## Trigger

The Trigger entity schema contains the necessary information to describe the stream of events that it represents. For more information, see the document "[Creating Triggers and Rules](https://github.com/apache/incubator-openwhisk/blob/master/docs/triggers_rules.md)".

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
  <p>The optional credential used to acces the feed service.</p>
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

---

## Rule

The Rule entity schema contains the information necessary to associates one trigger with one action, with every firing of the trigger causing the corresponding action to be invoked with the trigger event as input. For more information, see the document "[Creating Triggers and Rules](https://github.com/apache/incubator-openwhisk/blob/master/docs/triggers_rules.md)".

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
