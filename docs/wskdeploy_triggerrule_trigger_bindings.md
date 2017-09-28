# Triggers and Rules

## Using a Deployment file to bind Trigger parameters
This example builds on the previous [Trigger-Rule example](wskdeploy_triggerrule_basic.md) and will demonstrate how to use a Deployment File to bind values for a Trigger’s input parameters when applied against a compatible Manifest File

### Manifest File
Let’s use a variant of the Manifest file from the previous Trigger Rule example; however, we will leave the parameters on the ‘meetPerson’ Trigger unbound and only with type declarations.

#### _Example: “Hello world” Action, Trigger and Rule with no Parameter bindings_
```yaml
package:
  name: hello_world_package
  ... # Package keys omitted for brevity
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

### Deployment file
Let’s create a Deployment file that is desined to be applied to the Manifest file (above) which will contain the parameter bindings (i.e., the values) for the 'meetPerson' Trigger.

#### _Example: Deployment file that binds parameters to the 'meetPerson' Trigger_
```yaml
package:
  hello_world_package:
    triggers:
      meetPerson:
        inputs:
          name: string
          place: string
          children: integer
          height: float
```
### Deploying
```sh
$ wskdeploy -m docs/examples/manifest_hello_world_triggerrule_unbound.yaml -d docs/examples/deployment_hello_world_triggerrule_bindings.yaml
```

### Triggering

Fire the '```meetPerson```' Trigger:
```sh
$ wsk trigger fire meetPerson
```

#### Result
which results in an Activation ID:
```
$ wsk activation list

d03ee729428d4f31bd7f61d8d3ecc043 hello_world_triggerrule
3e10a54cb6914b37a8abcab53596dcc9 meetPersonRule
5ff4804336254bfba045ceaa1eeb4182 meetPerson

$ wsk activation get d03ee729428d4f31bd7f61d8d3ecc043

"result": {
   "details": "You have 13 children and are 1.2 m. tall.",
   "greeting": "Hello, Sam from the Shire"
}
```

### Discussion
TBD

### Source code
- [manifest_hello_world_triggerrule.yaml](examples/manifest_hello_world_triggerrule.yaml)
- [deployment_hello_world_triggerrule_bindings.yaml](docs/examples/deployment_hello_world_triggerrule_bindings.yaml)
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
    <td><a href="wskdeploy_triggerrule_basic.md#triggers-and-rules">&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
<!--    <td><a href="">next&nbsp;&gt;&gt;</a></td> -->
  </tr>
</table>
</div>
</html>
