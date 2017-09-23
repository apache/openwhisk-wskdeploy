## Package

The Package entity schema is used to define an OpenWhisk package within a manifest.

### Fields
<html>
<table width="100%">
 <tr>
  <th width="16%">
   <p>Key Name</p>
  </th>
  <th width="12%">
   <p>Required</p>
  </th>
  <th width="16%">
   <p>Value Type</p>
  </th>
  <td width="14%">
   <p>Default</p>
  </th>
  <th width="40%">
   <p>Description</p>
  </th>
 </tr>
 <tr>
  <td>version</td>
  <td>yes</td>
  <td><a href="spec_types_yaml.md#yaml-types">version</a></td>
  <td>N/A</td>
  <td>The required user-controlled version for the Package.</td>
 </tr>
 <tr>
  <td>license</td>
  <td>no</td>
  <td>string</td>
  <td>N/A</td>
  <td>The required value that indicates the type of license the Package is governed by.
   <p>The value is required to be a valid Linux-SPDX value. See <a href="https://spdx.org/licenses/">https://spdx.org/licenses/</a>.</p></td>
 </tr>
 <tr>
  <td>credential</td>
  <td>no</td>
  <td>string</td>
  <td>N/A</td>
  <td>
   <p>The optional Credential used for all entities within the Package.</p>
   <p>The value contains either:</p>
   <p>A credential string.</p>
   <p>The optional name of a credential (e.g., token) that is defined elsewhere.</p>
  </td>
 </tr>
 <tr>
  <td>dependencies</td>
  <td>no</td>
  <td>list of Dependency</td>
  <td>N/A</td>
  <td>The optional list of external OpenWhisk packages the manifest needs deployed before it can be deployed.</td>
 </tr>
 <tr>
  <td>repositories</td>
  <td>no</td>
  <td>list of Repository</td>
  <td>N/A</td>
  <td>The optional list of external repositories that contain functions and other artifacts that can be found by tooling.</td>
 </tr>
 <tr>
  <td>actions</td>
  <td>no</td>
  <td>list of Action</td>
  <td>N/A</td>
  <td>Optional list of OpenWhisk Action entity definitions.</td>
 </tr>
 <tr>
  <td>sequences</td>
  <td>no</td>
  <td>list of Sequence</td>
  <td>N/A</td>
  <td>Optional list of OpenWhisk Sequence entity definitions.</td>
 </tr>
 <tr>
  <td>triggers</td>
  <td>no</td>
  <td>list of Trigger</td>
  <td>N/A</td>
  <td>Optional list of OpenWhisk Trigger entity definitions.</td>
 </tr>
 <tr>
  <td>rules</td>
  <td>no</td>
  <td>list of Rule</td>
  <td>N/A</td>
  <td>Optional list of OpenWhisk Rule entity definitions.</td>
 </tr>
 <tr>
  <td>feeds</td>
  <td>no</td>
  <td>list of Feed</td>
  <td>N/A</td>
  <td>Optional list of OpenWhisk Feed entity definitions.</td>
 </tr>
 <tr>
  <td>compositions</td>
  <td>no</td>
  <td>list of Composition</td>
  <td>N/A</td>
  <td>Optional list of OpenWhisk Composition entity definitions.</td>
 </tr>
 <tr>
  <td>apis</td>
  <td>no</td>
  <td>list of API</td>
  <td>N/A</td>
  <td>Optional list of API entity definitions.</td>
 </tr>
</table>
</html>

### Grammar

```yaml
<packageName>:
    version: <version>
    license: <string>
    repositories: <list of Repository>
    actions: <list of Action>
    sequences: <list of Sequence>
    triggers: <list of Trigger>
    rules: <list of Rule>
    feeds: <list of Feed>
    compositions: <list of Composition>
    apis: <list of API>
```

### Example

```yaml
my_whisk_package:
  description: A complete package for my awesome action to be deployed
  version: 1.2.0
  license: Apache-2.0
  actions:
    my_awsome_action:
      <Action schema>
  triggers:
    trigger_for_awesome_action:
      <Trigger schema>
  rules:
    rule_for_awesome_action>
      <Rule schema>
```

### Requirements

- The Package name MUST be less than or equal to 256 characters.
- The Package entity schema includes all general <a href="#SCHEMA_ENTITY">Entity Schema</a> fields in addition to any fields declared above.
- A valid Package license value MUST be one of the <a href="#REF_LINUX_SPDX">Linux SPDX</a> license values; for example: Apache-2.0 or GPL-2.0+, or the value 'unlicensed'.
- Multiple (mixed) licenses MAY be described using using <a href="#REF_NPM_SPDX_SYNTAX">NPM SPDX license syntax</a>.
- A valid Package entity MUST have one or more valid Actions defined.

### Notes

- Currently, the 'version' value is not stored in Apache OpenWhisk, but there are plans to support it in the future.
- Currently, the 'license' value is not stored in Apache OpenWhisk, but there are plans to support it in the future.
- The Trigger and API entities within the OpenWhisk programming model are considered outside the scope of the Package. This means that Trigger and API information will not be returned when using the OpenWhisk Package API:
  - ```$ wsk package list <package name>```
- However, their information may be retrieved using respectively:</li>
  - ```$ wsk trigger list -v```
  - ```$ wsk api list -v```
