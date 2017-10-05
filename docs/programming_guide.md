# ```wskdeploy``` utility by example
_A step-by-step guide for deploying Apache OpenWhisk applications using Package Manifest files._

This guide will walk you through how to describe OpenWhisk applications and packages using the [OpenWhisk Packaging Specification](https://github.com/apache/incubator-openwhisk-wskdeploy/tree/master/specification#openwhisk-packaging-specification) and deploy them through the Whisk Deploy (```wskdeploy```) utility. Please use the specification as the ultimate reference for all Manifest file grammar and syntax.

## Getting started

### Setting up your Host and Credentials
In order to deploy your OpenWhisk package, at minimum, the ```wskdeploy``` utility needs valid OpenWhisk APIHOST and AUTH variable to attempt deployment. Please read the [Configuring wskdeploy](wskdeploy_configuring.md#configuring-wskdeploy)

### Debugging your Package Manifests
In addition to the normal output the ```wskdeploy``` utility provides, you may enable additional information that may further assist you in debugging. Please read the [Debugging Whisk Deploy](wskdeploy_debugging.md#debugging-wskdeploy) document.

### FAQ for ```wskdeploy```
Answers to Frequently Asked Questions may be found in the [wskdeploy utility FAQ](wskdeploy_faq.md).

---

# Guided Examples

Below is the list of "guided examples" where you can start by "Creating a 'hello world' application" and traverse through each example or jump to any example that interests you.

Each example shows the "code", that is the Package Manifest, Deployment file and Actions that will be used to deploy that application or package, as well as discusses the interesting features the example is highlighting.

- Package examples
  - [Creating a minimal Package](wskdeploy_package_minimal.md#packages) - creating a basic package manifest and deploying it.
- Action examples
  - [The "Hello World" Action](wskdeploy_action_helloworld.md#actions) - deploy a "hello world" JavaScript function using a manifest.
  - [Setting your Function's runtime](wskdeploy_action_runtime.md#actions) - explicitly set the runtime language and version to deplly your action onto.
  - [Adding fixed input parameters](wskdeploy_action_fixed_parms.md#actions) - bind fixed values to the input parameters of "hello world".
  - [Typed Parameters](wskdeploy_action_typed_parms.md#actions) - declare named input and output parameters on an Action with their types.
  - [Advanced Parameters](wskdeploy_action_advanced_parms.md#actions) - input and output parameter declarations with types, descriptions, defaults and more.
  - [Using Environment Variables](wskdeploy_action_env_var_parms.md#actions) - reference values from environment variables and bind them to an Action's input parameters.
_ Sequences examples:
  - [Sequencing Actions together](wskdeploy_sequence_basic.md#sequences) - sequence three actions together to process and augment data.
- Trigger and Rule examples.
  - [Basic Trigger and Rule](wskdeploy_triggerrule_basic.md#triggers-and-rules) - adding a basic trigger and rule to the advanced Parameter "hello world".
  - [Binding parameters in a Deployment file](wskdeploy_triggerrule_trigger_bindings.md#triggers-and-rules) - using a deployment file to bind values to a Trigger’s parameters and applying them to a compatible manifest file.

---
<!--
 Bottom Navigation
-->
<html>
<div align="center">
<table align="center">
  <tr>
    <td><a>&lt;&lt;&nbsp;previous</a></td>
    <td><a href="programming_guide.md#guided-examples">Example Index</a></td>
    <td><a href="wskdeploy_package_minimal.md#packages">next&nbsp;&gt;&gt;</a></td>
  </tr>
</table>
</div>
</html>
