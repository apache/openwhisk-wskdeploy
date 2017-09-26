# Triggers and Rules

## Creating a Trigger for an Action
This example shows how to create a Trigger that is compatible with the previous, more advanced "Hello world" Action, which has multiple input parameters of different types, and connect them together using a Rule.

### Manifest File
#### _Example: “Hello world” Action with a compatible Trigger and Rule_
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
        name: person
        place: the Shire
        children: integer
        height: float

  rules:
    meetPersonRule:
      trigger: meetPerson
      action: hello_world_triggerrule
```

### Deploying
```sh
$ wskdeploy -m docs/examples/manifest_hello_world_triggerrule.yaml
```

### Invoking
```sh
$ wsk action invoke hello_world_package/hello_world_triggerrule --blocking
```

### Result
```sh
"result": {
  "details": "You have 0 children and are 0 m. tall.",
  "greeting": "Hello,  from "
},
```

### Discussion
TODO

### Source code
- [manifest_hello_world_triggerrule.yaml](examples/manifest_hello_world_triggerrule.yaml)
- [deployment_hello_world_triggerrule.yaml](examples/deployment_hello_world_triggerrule.yaml)
- [hello_plus.js](examples/src/hello_plus.js)

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
