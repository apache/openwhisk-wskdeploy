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

# Using Slack Package with `wskdeploy`

The [Slack usecase](https://github.com/apache/incubator-openwhisk-wskdeploy/tree/master/tests/usecases/slack) demonstrates how to build an OpenWhisk app to post a sample message to slack every hour using `wskdeploy`.

OpenWhisk comes with a [Slack package](https://github.com/apache/incubator-openwhisk-catalog/blob/master/packages/slack/README.md) which can be used to post messages to slack. For our app to post hourly messages, we need:

- [manifest.yaml](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/tests/usecases/slack/manifest.yaml)
- [deployment.yaml](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/tests/usecases/slack/deployment.yaml)
- [Action File](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/tests/usecases/slack/src/post-to-slack.js)

All you have to do is export few environment variables with your slack webhook settings in `deployment.yaml` to deploy this app. You can create a new incoming webhook by following step by step instructions from [here](https://github.com/apache/incubator-openwhisk-GitHubSlackBot/blob/master/docs/add-webhook-to-slack.md).

```yaml
    dependencies:
        slack-package-to-post-messages:
            location: /whisk.system/slack
            inputs:
                username: $SLACK_USERNAME
                url: $SLACK_URL
                channel: $SLACK_CHANNEL
```

### Step 1: Deploy

Deploy it using `wskdeploy`:

```
wskdeploy -m tests/usecases/slack/manifest.yaml -d tests/usecases/slack/deployment.yaml
```

### Step 2: Verify

```
$ wsk package get SlackPackage
$ wsk package get slack-package-to-post-messages
$ wsk trigger get everyhour
$ wsk rule get post-to-slack-every-hour
```
### Step 3: Run

Fire the `everyhour` trigger and notice a new message on your slack channel:

```
Activation: post-to-slack (9909dd5229e84526bff9902a2cd860df)
[
    "2017-09-12T23:05:17.17872899Z  stdout: Hello from WskDeploy!",
    "2017-09-12T23:05:17.549177677Z stdout: Posted message to slack"
]
```
