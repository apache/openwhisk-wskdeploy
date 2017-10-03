# Sequences

## Creating a basic Action sequence

This example

This example:
- TBD

### Manifest File

#### _Example: input and output parameters with explicit types and descriptions_
```yaml

```

```javascript
function main(params) {
    var member = {name:"", place:"", occupation:"", height:0.0, joined:""};
    name = params.name;
    place = params.place;
    occupation = params.job;
    height = params.height;
    join_date = Date.now();
    return { joined: member };
}
```

### Deploying
```sh
$ wskdeploy -m docs/examples/manifest_sequence_basic.yaml
```

### Invoking
```sh
$ wsk action invoke
```

### Result
The invocation should return an 'ok' with a response that includes this result:
```json

```

### Discussion
-

### Source code
The manifest file for this example can be found here:
-

### Specification
For convenience, the Actions and Parameters grammar can be found here:
- **[Sequences](../specification/html/spec_sequences.md#sequences)**
- **[Actions](../specification/html/spec_actions.md#actions)**

---
<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a href="">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
