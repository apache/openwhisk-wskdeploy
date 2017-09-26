# Triggers and Rules

## Adding fixed input parameters
TODO

### Manifest File

#### _Example: “Hello world” with fixed input values for ‘name’ and ‘place’_
```yaml
package:
  name: hello_world_package
  version: 1.0
  license: Apache-2.0
  actions:

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
