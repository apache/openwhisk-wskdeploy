<html>
<h2>Package</h2>
<p>The Package entity schema is used to define an OpenWhisk package within a manifest.</p>
<h3>Fields</h3>
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
  <td>
   <p>version</p>
  </td>
  <td>
   <p>yes</p>
  </td>
  <td>
   <p>version</p>
  </td>
  <td>
   <p>N/A</p>
  </td>
  <td>
   <p>The required user-controlled version for the Package. </p>
  </td>
 </tr>
 <tr>
  <td>
   <p>license</p>
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
   <p>The required value that indicates the type of license the
    Package is governed by.</p>
   <p>The value is required to be a valid Linux-SPDX value. See <a href="https://spdx.org/licenses/">https://spdx.org/licenses/</a>.</p>
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
   <p>string</p>
  </td>
  <td>
   <p>N/A</p>
  </td>
  <td>
   <p>The optional Credential used for all entities within the Package.</p>
   <p>The value contains either:</p>
   <p>A credential string.</p>
   <p>The optional name of a credential (e.g., token) that is defined elsewhere.</p>
  </td>
 </tr>
 <tr>
  <td>
   <p>dependencies</p>
  </td>
  <td>
   <p>no</p>
  </td>
  <td>
   <p>list of Dependency</p>
  </td>
  <td>
   <p>N/A</p>
  </td>
  <td>
   <p>The optional list of external OpenWhisk packages the manifest needs deployed before it can be deployed.</p>
  </td>
 </tr>
 <tr>
  <td>
   <p>repositories</p>
  </td>
  <td>
   <p>no</p>
  </td>
  <td>
   <p>list of Repository</p>
  </td>
  <td>
   <p>N/A</p>
  </td>
  <td>
   <p>The optional list of external repositories that contain functions and other artifacts that can be found by tooling.</p>
  </td>
 </tr>
 <tr>
  <td>
   <p>actions</p>
  </td>
  <td>
   <p>no</p>
  </td>
  <td>
   <p>list of Action</p>
  </td>
  <td>
   <p>N/A</p>
  </td>
  <td>
   <p>Optional list of OpenWhisk Action entity definitions.</p>
  </td>
 </tr>
 <tr>
  <td>
   <p>sequences</p>
  </td>
  <td>
   <p>no</p>
  </td>
  <td>
   <p>list of Sequence</p>
  </td>
  <td>
   <p>N/A</p>
  </td>
  <td>
   <p>Optional list of OpenWhisk Sequence entity definitions.</p>
  </td>
 </tr>
 <tr>
  <td>
   <p>triggers</p>
  </td>
  <td>
   <p>no</p>
  </td>
  <td>
   <p>list of Trigger</p>
  </td>
  <td>
   <p>N/A</p>
  </td>
  <td>
   <p>Optional list of OpenWhisk Trigger entity definitions.</p>
  </td>
 </tr>
 <tr>
  <td>
   <p>rules</p>
  </td>
  <td>
   <p>no</p>
  </td>
  <td>
   <p>list of Rule</p>
  </td>
  <td>
   <p>N/A</p>
  </td>
  <td>
   <p>Optional list of OpenWhisk Rule entity definitions.</p>
  </td>
 </tr>
 <tr>
  <td>
   <p>feeds</p>
  </td>
  <td>
   <p>no</p>
  </td>
  <td>
   <p>list of Feed</p>
  </td>
  <td>
   <p>N/A</p>
  </td>
  <td>
   <p>Optional list of OpenWhisk Feed entity definitions.</p>
  </td>
 </tr>
 <tr>
  <td>
   <p>compositions</p>
  </td>
  <td>
   <p>no</p>
  </td>
  <td>
   <p>list of Composition</p>
  </td>
  <td>
   <p>N/A</p>
  </td>
  <td>
   <p>Optional list of OpenWhisk Composition entity definitions.</p>
  </td>
 </tr>
 <tr>
  <td>
   <p>apis</p>
  </td>
  <td>
   <p>no</p>
  </td>
  <td>
   <p>list of API</p>
  </td>
  <td>
   <p>N/A</p>
  </td>
  <td>
   <p>Optional list of API entity definitions.</p>
  </td>
 </tr>
</table>

<h3>Grammar</h3>
</html>

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

<html>
<h3>Example</h3>
</html>

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

<html>
<h3>Requirements</h3>
<ul>
<li>The Package name MUST be less than or equal to 256 characters.</li>
<li>The Package entity schema includes all general <a href="#SCHEMA_ENTITY">Entity Schema</a> fields in addition to any fields declared above.</li>
<li>A valid Package license value MUST be one of the <a href="#REF_LINUX_SPDX">Linux SPDX</a> license values; for example: Apache-2.0 or GPL-2.0+, or the value 'unlicensed'.</li>
<li>Multiple (mixed) licenses MAY be described using using <a href="#REF_NPM_SPDX_SYNTAX">NPM SPDX license syntax</a>.</li>
<li>A valid Package entity MUST have one or more valid Actions defined.</li>
</ul>

<h3>Notes</h3>
<ul>
<li>Currently, the 'version' value is not stored in Apache OpenWhisk, but there are plans to support it in the future.</li>
<li>Currently, the 'license' value is not stored in Apache OpenWhisk, but there are plans to support it in the future.</li>
<li>The Trigger and API entities within the OpenWhisk programming model are considered outside the scope of the Package. This means that Trigger and API information will not be returned when using the OpenWhisk Package API:</li>
<ul>
<li><code>$ wsk package list &lt;package name&gt;</code></li>
</ul>
<li>However, their information may be retrieved using respectively:</li>
<ul>
<li><code>$ wsk trigger list -v</li></code>
<li><code>$ wsk api list -v</li>
</ul>
</ul>
</html>
