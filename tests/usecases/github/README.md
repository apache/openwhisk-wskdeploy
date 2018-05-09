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

# Using GitHub Package with `wskdeploy`

The [GitHub usecase](https://github.com/apache/incubator-openwhisk-wskdeploy/tree/master/tests/usecases/github) demonstrates how to build an OpenWhisk app to display github commit messages using `wskdeploy`.

OpenWhisk comes with a [GitHub package](https://github.com/apache/incubator-openwhisk-catalog/blob/master/packages/github/README.md) which can be used to run GitHub APIs. For our app to display github commits, we need:

- [manifest.yaml](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/tests/usecases/github/manifest.yaml)
- [deployment.yaml](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/tests/usecases/github/deployment.yaml)
- [Action File](https://github.com/apache/incubator-openwhisk-wskdeploy/blob/master/tests/usecases/github/src/print-github-commits.js)

All you have to do is add your own credentials in `deployment.yaml` to deploy this app.

```yaml
GitHubWebhookTrigger:
    inputs:
        username: <username>  # replace it with your GitHub username
        repository: <repo>    # replace it with your repo eg. apache/incubator-openwhisk-wskdeploy
        accessToken: <token>  # replace it with token which has access to the specified repo
        events: push          # push for commits
```

### Step 1: Deploy

Deploy it using `wskdeploy`:

```
wskdeploy -m tests/usecases/github/manifest.yaml -d tests/usecases/github/deployment.yaml
```

### Step 2: Verify

```
$ wsk package get GitHubCommits
$ wsk trigger get GitHubWebhookTrigger
$ wsk rule get rule-for-github-commits
```
### Step 3: Run

Push a sample commit to your github repo and you will see `print-github-commits` is activated:

```
Activation: print-github-commits (28f)
[
    "2017-09-01T19:08:38.46838161Z  stdout: Display GitHub Commit Details for GitHub repo:  https://github.com/your-repo",
    "2017-09-01T19:08:38.468756344Z stdout: Bob Smith added code changes with commit message: Updating README to appear on openwhisk",
    "2017-09-01T19:08:38.46877569Z  stdout: Commit logs are:",
    ...
]
```
