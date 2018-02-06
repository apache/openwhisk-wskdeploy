<!--
#
# Licensed to the Apache Software Foundation (ASF) under one or more contributor
# license agreements.  See the NOTICE file distributed with this work for additional
# information regarding copyright ownership.  The ASF licenses this file to you
# under the Apache License, Version 2.0 (the # "License"); you may not use this
# file except in compliance with the License.  You may obtain a copy of the License
# at:
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software distributed
# under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
# CONDITIONS OF ANY KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations under the License.
#
-->

# GitHub Slack Bot

[Github Slack Bot](https://github.com/apache/incubator-openwhisk-GitHubSlackBot)
is an application designed to post updates to Slack when a GitHub pull request is
ready to merge or a list of pull requests are under review for certain days and
haven't merged.

You can find detailed Architecture and Usage at
[incubator-openwhisk-githubslackbot](https://github.com/apache/incubator-openwhisk-GitHubSlackBot).

Github Slack Bot application is dependent on three major components:

1. [Cloudant Package](https://github.com/apache/incubator-openwhisk-package-cloudant)
2. [GitHub Package](https://github.com/apache/incubator-openwhisk-catalog/tree/master/packages/github)
3. [Slack Package](https://github.com/apache/incubator-openwhisk-catalog/tree/master/packages/slack)


#### `manifest.yaml` for Cloudant Package

```yaml
    dependencies:
        cloudant-package:
            location: /whisk.system/cloudant
            inputs:
                username: $CLOUDANT_USERNAME
                password: $CLOUDANT_PASSWORD
                host: ${CLOUDANT_USERNAME}.cloudant.com
```

#### `manifest.yaml` for Github Package

```yaml
    dependencies:
            github-package:
                location: /whisk.system/github
                inputs:
                    username: $GITHUB_USERNAME
                    repository: $GITHUB_REPOSITORY
                    accessToken: $GITHUB_ACCESSTOKEN
```

#### `manifest.yaml` for Slack Package


```yaml
    dependencies:
        slack-package:
            location: /whisk.system/slack
            inputs:
                username: $SLACK_USERNAME
                url: $SLACK_URL
                channel: $SLACK_CHANNEL
```
### Step 1: Deploy


Export the following env. variables before running `wskdeploy`:

```
CLOUDANT_USERNAME
CLOUDANT_PASSWORD
CLOUDANT_DATABASE
GITHUB_USERNAME
GITHUB_REPOSITORY
GITHUB_ACCESSTOKEN
SLACK_USERNAME
SLACK_URL
SLACK_CHANNEL
```
Deploy it using `wskdeploy`:

```
wskdeploy -p tests/apps/openwhisk-githubslackbot
```

### Step 2: Verify

```
$ wsk package get TrackPRsInCloudant
$ wsk package get GitHubWebHook
$ wsk package get PostPRToSlack
$ wsk action get track-pull-requests
$ wsk action get find-delayed-pull-requests
$ wsk action get post-to-slack
$ wsk trigger get GitHubWebHookTrigger
$ wsk trigger get Every12Hours
$ wsk rule get RuleToTrackPullRequests
$ wsk rule get RuleToPostGitHubPRsToSlack
```
