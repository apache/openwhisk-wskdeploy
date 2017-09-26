# Triggers and Rules

## Creating a Trigger for an Action
This example shows how to create a Trigger that is compatible with the previous, more advanced "Hello world" Action, which has multiple input parameters of different types, and connect them together using a Rule.

### Manifest File
#### _Example: “Hello world” with fixed input values for ‘name’ and ‘place’_
```yaml
package:
  name: hello_world
  version: 1.0
  license: Apache-2.0
  actions:
    hello_world_triggerrule:
      function: src/hello_plus.js
      runtime: nodejs
      inputs:
        name: string
        place: string
        children: integer
        height: float
      outputs:
        greeting: string
        details: string

  triggers:
    meetPerson:
      inputs:
        name: string
        place: string
        children: integer
        height: float

  rules:
    meetPersonRule:
      trigger: meetPerson
      action: hello_world_triggerrule
```

### Deploying
```sh
$ wskdeploy -m
```

### Invoking
```sh
$ wsk action invoke
```

### Result
```sh
"result": {

},
```

### Discussion
TODO

### Source code
TODO

### Specification
For convenience, the Actions and Parameters grammar can be found here:
- **[Triggers and Rules](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/specification/html/spec_trigger_rule.md#triggers-and-rules)**

---
<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a href="wskdeploy_action_advanced_parms.md#actions">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
<!--    <td><a href="">next&nbsp;&gt;&gt;</a></td> -->
  </tr>
</table>
</div>
</html>
