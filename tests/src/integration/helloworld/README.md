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

# Integration Test - helloworld

### Package description

The [manifest.yaml](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/tests/src/integration/helloworld/manifest.yaml) file defines:

- a Package named `IntegrationTestHelloWorld` which contains:

    - Four actions:

        - an action named `helloNodejs`
        - an action named `helloJava`
        - an action named `helloPython`
        - an action named `helloSwift`

    - A Sequence `hello-world-series` which includes all four actions.

    - A trigger `triggerHelloworld` to invoke `hello-world-series` sequence

    - A rule `ruleMappingHelloworld` to associate the sequence `hello-world-series` with trigger `triggerHelloworld`

- `helloNodejs`:

    - accepts two parameters:
        - `name` (string) (default: Amy)
        - `place` (string) (default: Paris)
    - returns `Hello Amy from Paris`

- `helloJava`

    - accepts one parameter:
        - `name` (string) (default: Bob)
    - returns `Hello Bob!`

- `helloPython`

    - accepts one parameter:
        - `name` (string) (default: Henry)
    - returns `Hello Henry!`

- `helloSwift`

    - accepts one parameter:
        - `name` (string) (no default)
    - returns `Hello stranger!`


### How to deploy and test

#### Step 1. Deploy

Deploy it using `wskdeploy`:

```
$ wskdeploy -p tests/src/integration/helloworld
```

#### Step 2. Verify

```
$ wsk package get IntegrationTestHelloWorld #lists all four actions and a sequence
$ wsk trigger get triggerHelloworld
$ wsk rule get ruleMappingHelloworld
```

#### Step 3. Invoke

```
# invoke all four actions in a sequence hello-world-series
# results in four activation IDs and polling on displays outputs from all four actions
$ wsk trigger fire triggerHelloworld
Activation: helloSwift (5a0)
[
    "2017-09-01T17:09:51.589079299Z stdout: [\"greeting\": \"Hello stranger!\"]"
]

Activation: helloPython (219)
[
    "2017-09-01T17:09:51.554034831Z stdout: Hello Henry!"
]

Activation: helloJava (725)
[
    "2017-09-01T17:09:51.547499195Z stdout: {\"greeting\":\"Hello Bob!\"}"
]

Activation: triggerHelloworld (566)
[]

Activation: ruleMappingHelloworld (270)
[]

Activation: helloNodejs (495)
[
    "2017-09-01T17:09:51.459935244Z stdout: Hello, Amy from Paris"
]
```
