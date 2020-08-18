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

# Whisk Deploy Interpolation

Here is the list of YAML key/value in manifest/deployment file for which Whisk Deploy
supports interpolation including evaluating strings using environment variables.
For example, `$USERNAME` and `${USERNAME}` evaluates to environment variable `USERNAME`.
It also supports double `$` notation, for example, `$${USERNAME}` evaluates to `${USERNAME}`.
## Manifest File

#### Package Name

```yaml
project:
    name: helloworld
    packages:
        $PACKAGE_NAME_1:
            actions:
                ...
        ${PACKAGE_NAME_2}:
            actions:
                ...
        ${PACKAGE_NAME_3}_package:
            actions:
                ...
```

#### Annotation (under Action, Sequence, Trigger, Rule, and Dependency)

```yaml
project:
    name: helloworld
    packages:
        helloworld:
            actions:
                hello:
                    function: hello.js
                    annotation:
                        username: $USERNAME
                        password: ${PASSWORD}
                        host: http://${USERNAME}@${PASSWORD}/${URL}.com
```

#### Action Function

```yaml
project:
    name: helloworld
    packages:
        helloworld:
            actions:
                hello1:
                    function: $OPENWHISK_FUNCTION_FILE
                    runtime: nodejs:default
                hello2:
                    function: ${OPENWHISK_FUNCTION_FILE}
                    runtime: nodejs:default
                hello3:
                    function: ${OPENWHISK_FUNCTION_PYTHON}.py
                hello4:
                    function: https://${OPENWHISK_FUNCTION_GITHUB_DIR}.js                    function: github.com/apache/openwhisk-test/packages/helloworlds
                hello5:
                    function: $OPENWHISK_FUNCTION_FILE
                    docker: $DOCKER_HUB_IMAGE
```

#### Trigger Feed

```yaml
project:
    name: helloworld
    packages:
        helloworld:
            triggers:
                everyhour:
                    feed: /whisk.system/alarms/alarm
                message-trigger:
                    feed: Cloud_Functions_${KAFKA_INSTANCE}_Credentials-1/messageHubFeed
                github-trigger:
                    feed: ${GITHUB_PACKAGE}/webhook
```

#### Inputs (under Package, Action, Dependency, Trigger, and Trigger Feed)

```yaml
project:
    name: helloworld
    packages:
        helloworld:
            actions:
                hello:
                    inputs:
                        username: $USERNAME
                        password: ${PASSWORD}
                        host: https://${USERNAME}@${PASSWORD}/github.com
```

## Deployment File

#### Inputs (under Package, Action, and Trigger)

```yaml
project:
    name: helloworld
    packages:
        helloworld:
            actions:
                hello:
                    inputs:
                        username: $USERNAME
                        password: ${PASSWORD}
                        host: https://${USERNAME}@${PASSWORD}/github.com
```




