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

# TriggerRule use case for wskdeploy

### Package description

This Package named `triggerrule` includes:
- An action named "greeting". It accepts two parameters "name" and "place" and will return a greeting message "Hello, name from place!"
- A trigger named as "locationUpdate"
- A rule named "myRule" to associate the trigger with the action.

### How to deploy and test

#### Step 1. Deploy the package.

```
$ wskdeploy -p tests/usecases/triggerrule
```

#### Step 2. Verify the action is installed.

```
$ wsk action list

{
  "name": "greeting",
  "publish": false,
  "annotations": [{
    "key": "exec",
    "value": "nodejs:default"
  }],
  "version": "0.0.1",
  "namespace": "<namespace>/triggerrule"
}
```

#### Step 3. Verify the action's last activation (ID) before invoking.

```
$ wsk activation list --limit 1 greeting

activations
3b80ec50d934464faedc04bc991f8673 greeting
```

#### Step 4. Invoke the trigger

```
$ wsk trigger fire locationUpdate --param name Bernie --param place "Washington, D.C."

ok: triggered /_/locationUpdate with id 5f8f26928c384afd85f47281bf85b302
```

#### Step 5. Verify the action is invoked

```
$ wsk activation list --limit 2 greeting

activations
a2b5cb1076b24c1885c018bb46f4e990 greeting
3b80ec50d934464faedc04bc991f8673 greeting
```

Indeed a new activation was listed.
