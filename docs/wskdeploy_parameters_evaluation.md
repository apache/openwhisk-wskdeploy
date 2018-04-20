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
# Whisk Deploy Parameters Evaluation

A project is generally defined as a collection of packages. A package is generally defined as a collection of OpenWhisk entities such as
actions, sequences, triggers, rules, and apis in manifest/deployment file. These OpenWhisk entities in manifest/deployment files
generally need data from users/environment for their successful deployment. This data includes (1) default values of action parameters
specially sensitive information such as passwords, (2) package bindings which are created outside of an existing deployment,
(3) GitHub credentials in case of deploying private GitHub repo, and many more.

We have designed functionality where a user can provide a list of such parameters needed per project and per package, such as:

```yaml
project:
    name: myproject
    parameters:
        SLACK_USERNAME:
            type: string
            description: "Slack User Name"
            required: true
            default: ${SLACK_USERNAME}
        SLACK_WEBHOOK_URL:
            type: string
            description: "Slack Webhook URL"
            required: true
            default: https://hooks.slack.com/services/${SLACK_WEBHOOK_URL}
        SLACK_CHANNEL:
            type: string
            description: "Slack Channel"
            required: true
            default: #general
        SLACK_TOKEN:
            type: string
            description: "Slack Token"
            required: true
            default: ${SLACK_TOKEN}
        GITHUB_USERNAME:
            type: string
            description: "GitHub Username"
            required: true
            default: ${GITHUB_USERNAME}
        GITHUB_REPOSITORY:
            type: string
            description: "GitHub Repository"
            required: true
            default: ${GITHUB_REPOSITORY}
        GITHUB_ACCESS_TOKEN:
            type: string
            description: "GitHub Access Token"
            required: true
            default: ${GITHUB_ACCESS_TOKEN}
        packages:
            github-slack-trigger:
                parameters:
                    SLACK_CHANNEL:
                        type: string
                        description: "Slack Channel"
                        required: true
                        default: #dev
                    RULE_NAME:
                        type: string
                        description: "Rule Name"
                        required: true
                        default: ${RULE_NAME}
                    TRIGGER_NAME:
                        type: string
                        description: "Trigger Name"
                        required: true
                        default: ${TRIGGER_NAME}
            slack-notifications:
                ...
```

Now, when you run `wskdeploy` to generate a report of these parameters with:

```
wskdeploy report -m manifest.yaml
{
    "github-slack-trigger": [
        {
            "name": "SLACK_USERNAME",
            "type": "string",
            "value": ""
            "description": "Slack User Name",
            "required": true,
        },
        {
            "name": "SLACK_WEBHOOK_URL",
            "type": "string",
            "value": " https://hooks.slack.com/services/"
            "description": "Slack Webhook URL",
            "required": true,
        },
        {
            "name": "SLACK_CHANNEL",
            "type": "string",
            "value": "#dev"
            "description": "Slack Channel",
            "required": true,
        },
        {
            "name": "SLACK_TOKEN",
            "type": "string",
            "value": ""
            "description": "Slack Token",
            "required": true,
        },
        {
            "name": "GITHUB_USERNAME",
            "type": string,
            "value": "",
            "description": "GitHub Username",
            "required": true,
        },
        {
            "name": "GITHUB_REPOSITORY",
            "type": string
            "value": "",
            "description": "GitHub Repository"
            "required": true,
        },
        {
            "name": "GITHUB_ACCESS_TOKEN",
            "type": string,
            "value": "",
            "description": "GitHub Access Token",
            "required": true,
        },
        {
            "name": "RULE_NAME",
            "type": string,
            "value": "",
            "description": "Rule Name"
            "required": true,
        },
        {
            "name": "TRIGGER_NAME"
            "type": string
            "value": "",
            "description": "Trigger Name"
            "required": true
        }
    ]
    "slack-notifications": [
        {
            "name": "SLACK_USERNAME",
            "type": "string",
            "value": ""
            "description": "Slack User Name",
            "required": true,
        },
        {
            "name": "SLACK_WEBHOOK_URL",
            "type": "string",
            "value": " https://hooks.slack.com/services/"
            "description": "Slack Webhook URL",
            "required": true,
        },
        {
            "name": "SLACK_CHANNEL",
            "type": "string",
            "value": "#general"
            "description": "Slack Channel",
            "required": true,
        },
        {
            "name": "SLACK_TOKEN",
            "type": "string",
            "value": ""
            "description": "Slack Token",
            "required": true,
        },
        {
            "name": "GITHUB_USERNAME",
            "type": string,
            "value": "",
            "description": "GitHub Username",
            "required": true,
        },
        {
            "name": "GITHUB_REPOSITORY",
            "type": string
            "value": "",
            "description": "GitHub Repository"
            "required": true,
        },
        {
            "name": "GITHUB_ACCESS_TOKEN",
            "type": string,
            "value": "",
            "description": "GitHub Access Token",
            "required": true,
        },
        {
            "name": "RULE_NAME",
            "type": string,
            "value": "",
            "description": "Rule Name"
            "required": true,
        },
        {
            "name": "TRIGGER_NAME"
            "type": string
            "value": "",
            "description": "Trigger Name"
            "required": true
        }
    ]
}
```

We can define parameters at the project level and also at the package level.
Project parameters are generally defined for a project which has multiple packages
and those packages have common parameters which are shared among those packages.
Package parameters are a collection of project parameters and parameters defined in that package.
Package parameters always takes higher precedence over project parameters i.e. package parameters
also defined at the project level takes value specified in the package. In the above example, `SLACK_CHANNEL`
is defined in `github-slack-trigger` and also listed at the project level. In this case, `SLACK_CHANNEL` is
assigned `#dev` which is specified in `github-slack-trigger` vs `#general` at the project level.


`wskdeploy report` mode interpolates parameter values in the manifest file and produces
the list of package parameters. In the above example, none of the referenced environment
variables were set and therefore it returned empty values for missing environment variables.
Let's look at how the values are calculated with environment variables set:

```
export SLACK_USERNAME=slack_username
export SLACK_WEBHOOK_URL=slack_webhook_url
export SLACK_TOKEN=slack_token
export GITHUB_USERNAME=github_username
export GITHUB_REPOSITORY=github_repository
export GITHUB_ACCESS_TOKEN=github_access_token
export RULE_NAME=rule_name
export TRIGGER_NAME=trigger_name
wskdeploy report -m manifest.yaml
{
    "github-slack-trigger": [
        {
            "name": "SLACK_USERNAME",
            "type": "string",
            "value": "slack_username"
            "description": "Slack User Name",
            "required": true,
        },
        {
            "name": "SLACK_WEBHOOK_URL",
            "type": "string",
            "value": " https://hooks.slack.com/services/slack_webhook_url"
            "description": "Slack Webhook URL",
            "required": true,
        },
        {
            "name": "SLACK_CHANNEL",
            "type": "string",
            "value": "#dev"
            "description": "Slack Channel",
            "required": true,
        },
        {
            "name": "SLACK_TOKEN",
            "type": "string",
            "value": "slack_token"
            "description": "Slack Token",
            "required": true,
        },
        {
            "name": "GITHUB_USERNAME",
            "type": string,
            "value": "github_username",
            "description": "GitHub Username",
            "required": true,
        },
        {
            "name": "GITHUB_REPOSITORY",
            "type": string
            "value": "github_repository",
            "description": "GitHub Repository"
            "required": true,
        },
        {
            "name": "GITHUB_ACCESS_TOKEN",
            "type": string,
            "value": "github_access_token",
            "description": "GitHub Access Token",
            "required": true,
        },
        {
            "name": "RULE_NAME",
            "type": string,
            "value": "rule_name",
            "description": "Rule Name"
            "required": true,
        },
        {
            "name": "TRIGGER_NAME"
            "type": string
            "value": "trigger_name",
            "description": "Trigger Name"
            "required": true
        }
    ]
    "slack-notifications": [
        {
            "name": "SLACK_USERNAME",
            "type": "string",
            "value": "slack_username"
            "description": "Slack User Name",
            "required": true,
        },
        {
            "name": "SLACK_WEBHOOK_URL",
            "type": "string",
            "value": " https://hooks.slack.com/services/slack_webhook_url"
            "description": "Slack Webhook URL",
            "required": true,
        },
        {
            "name": "SLACK_CHANNEL",
            "type": "string",
            "value": "#general"
            "description": "Slack Channel",
            "required": true,
        },
        {
            "name": "SLACK_TOKEN",
            "type": "string",
            "value": "slack_token"
            "description": "Slack Token",
            "required": true,
        },
        {
            "name": "GITHUB_USERNAME",
            "type": string,
            "value": "github_username",
            "description": "GitHub Username",
            "required": true,
        },
        {
            "name": "GITHUB_REPOSITORY",
            "type": string
            "value": "github_repository",
            "description": "GitHub Repository"
            "required": true,
        },
        {
            "name": "GITHUB_ACCESS_TOKEN",
            "type": string,
            "value": "github_access_token",
            "description": "GitHub Access Token",
            "required": true,
        },
        {
            "name": "RULE_NAME",
            "type": string,
            "value": "rule_name",
            "description": "Rule Name"
            "required": true,
        },
        {
            "name": "TRIGGER_NAME"
            "type": string
            "value": "trigger_name",
            "description": "Trigger Name"
            "required": true
        }
    ]
}
```

Now, project/package parameters can also be specified in deployment file and takes precedence over manifest file. For example:

Deployment file:

```yaml
project:
    name: myproject
    packages:
       github-slack-trigger:
            parameters:
                SLACK_CHANNEL:
                    type: string
                    description: "Slack Channel"
                    required: true
                    value: #dev-pr
```

Running `wskdeploy` with deployment file:

```
export SLACK_USERNAME=slack_username
...
wskdeploy report -m manifest.yaml -d deployment.yaml
{
    "github-slack-trigger": [
        ...
        {
            "name": "SLACK_CHANNEL",
            "type": "string",
            "value": "#dev-pr"
            "description": "Slack Channel",
            "required": true,
        },
        ...
    ]
    "slack-notifications": [
        ...
        {
            "name": "SLACK_CHANNEL",
            "type": "string",
            "value": "#dev-pr"
            "description": "Slack Channel",
            "required": true,
        },
        ...
    ]
}
```
In the end, `wskdeploy` supports a flag `--param` which takes the highest precedence over values specified in manifest/deployment file, for example:

```
export SLACK_USERNAME=slack_username
...
wskdeploy report -m manifest.yaml -d deployment.yaml --param SLACK_CHANNEL "#dev-push" --param GITHUB_REPOSITORY https://github.com
{
    "github-slack-trigger": [
        ...
        {
            "name": "SLACK_CHANNEL",
            "type": "string",
            "value": "#dev-push"
            "description": "Slack Channel",
            "required": true,
        },
        {
            "name": "GITHUB_REPOSITORY",
            "type": string
            "value": "https://github.com",
            "description": "GitHub Repository"
            "required": true,
        },
        ...
    ]
    "slack-notifications": [
        ...
        {
            "name": "SLACK_CHANNEL",
            "type": "string",
            "value": "#dev-push"
            "description": "Slack Channel",
            "required": true,
        },
        {
            "name": "GITHUB_REPOSITORY",
            "type": string
            "value": "https://github.com",
            "description": "GitHub Repository"
            "required": true,
        },
        ...
    ]
}
```

`--param` also supports specifying environment variables:

```
export SLACK_USERNAME=slack_username
export SLACK_CHANNEL=#wskdeploy
export GITHUB_REPOSITORY=https://github.com
...
wskdeploy report -m manifest.yaml -d deployment.yaml --param SLACK_CHANNEL "#dev-push" --param GITHUB_REPOSITORY $GITHUB_REPOSITORY 
{
    "github-slack-trigger": [
        ...
        {
            "name": "SLACK_CHANNEL",
            "type": "string",
            "value": "#wskdeploy"
            "description": "Slack Channel",
            "required": true,
        },
        {
            "name": "GITHUB_REPOSITORY",
            "type": string
            "value": "https://github.com",
            "description": "GitHub Repository"
            "required": true,
        },
        ...
    ]
    "slack-notifications": [
        ...
        {
            "name": "SLACK_CHANNEL",
            "type": "string",
            "value": "#wskdeploy"
            "description": "Slack Channel",
            "required": true,
        },
        {
            "name": "GITHUB_REPOSITORY",
            "type": string
            "value": "https://github.com",
            "description": "GitHub Repository"
            "required": true,
        },
        ...
    ]
}
```

Note that, the precedence order of reading and evaluating parameter values is:

1. `wskdeploy` CLI with `--param` with string interpolated using environment variables.
2. Deployment file with string interpolated using environment variables.
3. Manifest file with string interpolated using environment variables.






