<!--
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
-->

# Whisk Deploy Inputs

In this programming guide, we are going to discuss `inputs` section of manifest and
deployment file along with command line options `--param` and `--param-file`. First some background.

We define a project in manifest/deployment file which is a collection of packages.
A package in turn is a collection of OpenWhisk entities such as actions, sequences,
triggers, rules, apis, etc. These OpenWhisk entities in manifest/deployment files
generally need data from users/environment for its successful deployment. This data includes:

* Default values of action parameters including sensitive information such as credentials.
* Shared package bindings which are created outside of an existing deployment, for example, Cloudant, Slack, etc.
* Service credentials, for example, cloudant credentials, slack token, etc

Inputs can be specified at different levels:

 * Action Inputs
 * Trigger Inputs
 * Package Inputs
 * Project Inputs

And can be specified in multiple different ways:

* Manifest file
* Deployment file
* CLI using `--param` and `--param-file`

Before we dive into details of each level of inputs with all different ways, `wskdeploy` follows a particular order in which the values are read:

* Input values specified using `--param` and/or `--param-file` takes the highest precedence order. The values specified on CLI are taken to the server.
* Next, input values are read from deployment file
* Last, input values are read from manifest file

### Action Inputs:

Let's start with a simple example of a `helloworld` action which has two inputs `name` and `place`.

```yaml
packages:
    helloworldapp:
        actions:
            hello:
                inputs:
                    name:
                        type: string
                        description: "your first name"
                        required: false
                        value: Amy
                    place:
                        type: string
                        description: "The city name"
                        required: false
                        value: Paris
                code: |
                          function main(params) {
                              return {payload:  'Hello, ' + params.name + ' from ' + params.place};
                          }
                runtime: nodejs:default
```

Or (Single Line Inputs):

```yaml
packages:
    helloworldapp:
        actions:
            hello:
                inputs:
                    name: Amy
                    place: Paris
                code: |
                          function main(params) {
                              return {payload:  'Hello, ' + params.name + ' from ' + params.place};
                          }
                runtime: nodejs:default
```
Whisk deploy creates bindings at the action level with two parameters `name` and `place`:

```bash
./wskdeploy --preview -m tests/dat/manifest_validate_package_inputs_1.yaml

Packages:
Name: helloworldapp
    bindings:
    annotation:

  * action: hello
    bindings:
        - name : "Amy"
        - place : "Paris"
    annotation:
```

This is how two inputs `name` and `place` are stored in `hello` action on OpenWhisk server:

```
"parameters": [
        {
            "key": "name",
            "value": "Amy"
        }
        {
            "key": "place",
            "value": "Paris"
        },
    ],
```

### Action Inputs with Env. Variables

```yaml
packages:
    helloworldapp:
        actions:
            hello:
                inputs:
                    name:
                        type: string
                        description: "your first name"
                        required: true
                        value: $FIRST_NAME
                    place:
                        type: string
                        description: "The city name"
                        required: true
                        value: $CITY_NAME
                code: |
                          function main(params) {
                              return {payload:  'Hello, ' + params.name + ' from ' + params.place};
                          }
                runtime: nodejs:default
```

Deployment of this kind of manifest file results in following failure as inputs
`name` and `place` are marked `required` but their values `$FIRST_NAME` and `$CITY_NAME` could not be determined.

```bash
Error: manifestreader.go [92]: [ERROR_YAML_FILE_FORMAT_ERROR]: File: [manifest_validate_package_inputs_2.yaml]:
==> manifest_parser.go [148]: [ERROR_YAML_FILE_FORMAT_ERROR]: File: [manifest_validate_package_inputs_2.yaml]: Required inputs are missing values even after applying interpolation using env. variables. Please set missing env. variables and/or input values in manifest/deployment file or on CLI for following inputs: name, place
```

On the other side, single line inputs does not support this kind of validation as
they are not marked as `required` by default.

```yaml
packages:
    helloworldapp:
        actions:
            hello:
                inputs:
                    name: $FIRST_NAME
                    place: $CITY_NAME
                code: |
                          function main(params) {
                              return {payload:  'Hello, ' + params.name + ' from ' + params.place};
                          }
                runtime: nodejs:default
```

Action is created with two bindings `name` and `place` set to `""`:

```bash
./wskdeploy --preview -m tests/dat/manifest_validate_package_inputs_2.yaml
Packages:
Name: helloworldapp
    bindings:
    annotation:

  * action: hello
    bindings:
        - name : ""
        - place : ""
    annotation:
```

Now, after setting env. variables, `wskdeploy` creates an action with bindings similar to previous example:

```bash
export FIRST_NAME=Amy
export CITY_NAME=Paris

./wskdeploy --preview -m tests/dat/manifest_validate_package_inputs_2.yaml

Packages:
Name: helloworldapp
    bindings:
    annotation:

  * action: hello
    bindings:
        - name : "Amy"
        - place : "Paris"
    annotation:
```

### Action Inputs with `--param`

The input values can be overwritten using `--param` on CLI. Sample manifest in [Action Inputs](#Action Inputs) can be deployed using:

```bash
./wskdeploy --preview -m tests/dat/manifest_validate_package_inputs_1.yaml --param name Bob

Packages:
Name: helloworldapp
    bindings:
    annotation:

  * action: hello
    bindings:
        - name : "Bob"
        - place : "Paris"
    annotation:
```

### Project Inputs:

This example shows how env. variables `FIRST_NAME` and `CITY_NAME` are defined under project.
These two variables are needed for deployment and need not be propagated to OpenWhisk server.
This is the best and common practice of specifying inputs at the action level and env.
variables at the project level. Here, both `FIRST_NAME` and `CITY_NAME` have default
values but reads from environment if they are specified.

```yaml
project:
    name: helloworld
    inputs:
        FIRST_NAME:
            type: string
            description: "your first name"
            required: true
            value: Amy
        CITY_NAME:
            type: string
            description: "The city name"
            required: true
            value: Paris
    packages:
        helloworldapp:
            actions:
                hello:
                    inputs:
                        name:
                            type: string
                            description: "your first name"
                            required: true
                            value: $FIRST_NAME
                        place:
                            type: string
                            description: "The city name"
                            required: true
                            value: $CITY_NAME
                    code: |
                          function main(params) {
                              return {payload:  'Hello, ' + params.name + ' from ' + params.place};
                          }
                    runtime: nodejs:default
```

This is how bindings are created only under action:

```bash
./wskdeploy --preview -m tests/dat/manifest_validate_package_inputs_3.yaml
Packages:
Name: helloworldapp
    bindings:
    annotation:

  * action: hello
    bindings:
        - name : "Amy"
        - place : "Paris"
    annotation:
```

With env. variables:

```bash
export FIRST_NAME=Bob
export CITY_NAME=San Francisco
./wskdeploy --preview -m tests/dat/manifest_validate_package_inputs_3.yaml
Packages:
Name: helloworldapp
    bindings:
    annotation:

  * action: hello
    bindings:
        - name : "Bob"
        - place : "San Francisco"
    annotation:
```

### Package Inputs

```yaml
project:
    name: helloworld
    inputs:
        FIRST_NAME:
            type: string
            description: "your first name"
            required: true
            value: Amy
        CITY_NAME:
            type: string
            description: "The city name"
            required: true
            value: Paris
    packages:
        helloworldapp:
            inputs:
                name:
                    type: string
                    description: "your first name"
                    required: true
                    value: $FIRST_NAME
                place:
                    type: string
                    description: "The city name"
                    required: true
                    value: $CITY_NAME
            actions:
                hello:
                    code: |
                          function main(params) {
                              return {payload:  'Hello, ' + params.name + ' from ' + params.place};
                          }
                    runtime: nodejs:default
```

Now, bindings are created under Package:

```bash
./wskdeploy --preview -m tests/dat/manifest_validate_package_inputs_4.yaml
Packages:
Name: helloworldapp
    bindings:
        - name : "Amy"
        - place : "Paris"
    annotation:

  * action: hello
    bindings:
    annotation:
```

And can be overwritten using env. variables:

```bash
export FIRST_NAME=Bob
export CITY_NAME=San Francisco
./wskdeploy --preview -m tests/dat/manifest_validate_package_inputs_4.yaml
Packages:
Name: helloworldapp
    bindings:
        - name : "Bob"
        - place : "San Francisco"
    annotation:

  * action: hello
    bindings:
    annotation:
```


#### Package Inputs with Action Inputs:

Here, package `helloworldapp` has three inputs defined which are available to all
three actions `helloWithMorning`, `helloWithEvening`, and `helloWithNight` during
invocation. But two of the actions have redefined the same input `message` which
is created as an action binding on the server.

```yaml
project:
    name: helloworld
    inputs:
        FIRST_NAME:
            type: string
            description: "your first name"
            required: true
            value: Amy
        CITY_NAME:
            type: string
            description: "The city name"
            required: true
            value: Paris
    packages:
        helloworldapp:
            inputs:
                name:
                    type: string
                    description: "your first name"
                    required: true
                    value: $FIRST_NAME
                place:
                    type: string
                    description: "The city name"
                    required: true
                    value: $CITY_NAME
                message:
                    type: string
                    description: "The Message"
                    required: true
                    value: "Good Night"
            actions:
                helloWithMorning:
                    inputs:
                        message:
                            type: string
                            description: "The Message"
                            required: true
                            value: "Good Morning"
                    code: |
                          function main(params) {
                              return {payload:  'Hello, ' + params.message + ' ' + params.name + ' from ' + params.place};
                          }
                    runtime: nodejs:default
                helloWithEvening:
                    inputs:
                        message:
                            type: string
                            description: "The Message"
                            required: true
                            value: "Good Evening"
                    code: |
                          function main(params) {
                              return {payload:  'Hello, ' + params.message + ' ' + params.name + ' from ' + params.place};
                          }
                    runtime: nodejs:default
                helloWithNight:
                    code: |
                          function main(params) {
                              return {payload:  'Hello, ' + params.message + ' ' + params.name + ' from ' + params.place};
                          }
                    runtime: nodejs:default
```

Now, invoking `helloWithMorning` returns `Good Morning` and invoking `helloWithEvening`
returns `Good Evening` whereas invoking `helloWithNight` returns `Good Night` which
is stored as the package binding.

```bash
./wskdeploy --preview -m tests/dat/manifest_validate_package_inputs_5.yaml
Packages:
Name: helloworldapp
    bindings:
        - name : "Amy"
        - place : "Paris"
        - message : "Good Night"
    annotation:

  * action: helloWithMorning
    bindings:
        - message : "Good Morning"
    annotation:
  * action: helloWithEvening
    bindings:
        - message : "Good Evening"
    annotation:
  * action: helloWithNight
    bindings:
    annotation:
```
