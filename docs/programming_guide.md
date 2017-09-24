# ```wskdeploy``` utility by example
_A step-by-step guide for deploying Apache OpenWhisk applications using Package Manifest files._

This guide will walk you through how to describe OpenWhisk applications and packages using the [OpenWhisk Packaging Specification](https://github.com/apache/incubator-openwhisk-wskdeploy/tree/master/specification#openwhisk-packaging-specification) and deploy them through the Whisk Deploy (```wskdeploy```) utility. Please use the specification as the ultimate reference for all Manifest file grammar and syntax.

## Getting started
### Setting up your Host and Credentials
In order to deploy your OpenWhisk package, at minimum, the ```wskdeploy``` utility needs valid OpenWhisk APIHOST and AUTH variable to attempt deployment. Please read the [Configuring wskdeploy](wskdeploy_configuring.md#configuring-wskdeploy)

### Debugging your Package Manifests

In addition to the normal output the ```wskdeploy``` utility provides, you may enable additional information that may further assist you in debugging. Please read the [Debugging Whisk Deploy](wskdeploy_debugging.md#debugging-wskdeploy) document.

---

# Guided Examples

Below is the list of "guided examples" where you can start by "Creating a 'hello world' application" and traverse through each example or jump to any example that interests you.

Each example shows the "code", that is the Package Manifest, Deployment file and Actions that will be used to deploy that application or package, as well as discusses the interesting features the example is highlighting.

- [Creating a "Hello World" package](wsdeploy_hello_world.md#creating-a-hello-world-package)
- [Creating a minimal Package](wskdeploy_packages.md#packages) - creating a basic package and deploying it from your

---

<html>
<table align="center">
  <tr>
    <td><a href="wsdeploy_hello_world.md#creating-a-hello-world-package">previous</a></td>
    <td></td>
    <td><a>next</a></td>
  </tr>
</table>
</html>
