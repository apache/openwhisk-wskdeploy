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
# Package Parameters Evaluation

A package is generally defined as a collection of actions, sequences, triggers, rules, and apis in manifest/deployment file.
These OpenWhisk entities in manifest/deployment files generally need data from users/environment for their successful deployment.
This data includes default values of action parameters specially sensitive information such as passwords, package bindings which
are created outside of an existing deployment, GitHub credentials in case of deploying private GitHub repo, and many more.

We have designed functionality where a user can provide a list of such parameters needed per package, such as:

```yaml
project:
    name: myproject
    packages:
        mypackage:
            parameters:
                USERNAME:
                    type: string
                    description: "The User Name"
                    required: true
                    default: test 
                PASSWORD:
                    type: string
                    description: "The Password"
                    required: true
                    default: ${PASSWORD}
                DATABASE:
                    type: string
                    description: "The Database"
                    required: true
                OPENWHISK_DATABASE_FEED:
                    type: string
                    description: "The Database feed instance"
                    required: true
                    default: ${DATABASE_HOST}_instance
```

Now, when you run `wskdeploy` to generate a report of these parameters with:

```
wskdeploy report -m manifest.yaml
{
    "parameters": [
        {
            "name": "USERNAME",
            "type": "string",
            "value": "test"
            "description": "The User Name",
            "required": true,
        },
        {
            "name": "PASSWORD",
            "type": "string",
            "value": ""
            "description": "The Password",
            "required": true,
        },
        {
            "name": "DATABASE",
            "type": "string",
            "value": ""
            "description": "The Database",
            "required": true,
        },
        {
            "name": "OPENWHISK_DATABASE_FEED",
            "type": "string",
            "value": ""
            "description": "The Database feed instance",
            "required": true,
        }
    ]
}
```

`wskdeploy` interpolates parameter values in the manifest file and produces this list. In the above example, none of the referenced environment
variables were set and therefore it returned empty values for missing enviornment variables. Let's look at how the values are calculated with environment variables set:

```
export PASSWORD=password
export DATABASE_HOST=localhost
wskdeploy report -m manifest.yaml
{
    "parameters": [
        {
            "name": "USERNAME",
            "type": "string",
            "value": "test"
            "description": "The User Name",
            "required": true,
        },
        {
            "name": "PASSWORD",
            "type": "string",
            "value": "password"
            "description": "The Password",
            "required": true,
        },
        {
            "name": "DATABASE",
            "type": "string",
            "value": ""
            "description": "The Database",
            "required": true,
        },
        {
            "name": "OPENWHISK_DATABASE_FEED",
            "type": "string",
            "value": "localhost_instance"
            "description": "The Database feed instance",
            "required": true,
        }
    ]
}
```

Note that, `DATABASE` is still not set which can be fixed with deployment file:

Deployment file:

```yaml
project:
    name: myproject
    packages:
        mypackage:
            parameters:
                DATABASE:
                    type: string
                    description: "The Database"
                    required: true
                    value: mydepdatabase
```

Running `wskdeploy` with deployment file:

```
export PASSWORD=password
export DATABASE_HOST=localhost
wskdeploy report -m manifest.yaml -d deployment.yaml
{
    "parameters": [
        {
            "name": "USERNAME",
            "type": "string",
            "value": "test"
            "description": "The User Name",
            "required": true,
        },
        {
            "name": "PASSWORD",
            "type": "string",
            "value": "password"
            "description": "The Password",
            "required": true,
        },
        {
            "name": "DATABASE",
            "type": "string",
            "value": "mydepdatabase"
            "description": "The Database",
            "required": true,
        },
        {
            "name": "OPENWHISK_DATABASE_FEED",
            "type": "string",
            "value": "localhost_instance"
            "description": "The Database feed instance",
            "required": true,
        }
    ]
}
```

In the end, `wskdeploy` supports a flag `--param` which takes the highest precedence over values specified in manifest/deployment file for example:


```
export PASSWORD=password
export DATABASE_HOST=localhost
wskdeploy report -m manifest.yaml -d deployment.yaml --param USERNAME myusername --param DATABASE mydatabase --param OPENWHISK_DATABASE_FEED database_feed
{
    "parameters": [
        {
            "name": "USERNAME",
            "type": "string",
            "value": "myusername"
            "description": "The User Name",
            "required": true,
        },
        {
            "name": "PASSWORD",
            "type": "string",
            "value": "password"
            "description": "The Password",
            "required": true,
        },
        {
            "name": "DATABASE",
            "type": "string",
            "value": "mydatabase"
            "description": "The Database",
            "required": true,
        },
        {
            "name": "OPENWHISK_DATABASE_FEED",
            "type": "string",
            "value": "database_feed"
            "description": "The Database feed instance",
            "required": true,
        }
    ]
}
```
`--param` also supports specifying environment variable:

```
export USERNAME=myusername
export PASSWORD=password
export DATABASE_HOST=localhost
wskdeploy report -m manifest.yaml -d deployment.yaml --param USERNAME $USERNAME --param DATABASE mydatabase --param OPENWHISK_DATABASE_FEED database_feed
{
    "parameters": [
        {
            "name": "USERNAME",
            "type": "string",
            "value": "myusername"
            "description": "The User Name",
            "required": true,
        },
        {
            "name": "PASSWORD",
            "type": "string",
            "value": "password"
            "description": "The Password",
            "required": true,
        },
        {
            "name": "DATABASE",
            "type": "string",
            "value": "mydatabase"
            "description": "The Database",
            "required": true,
        },
        {
            "name": "OPENWHISK_DATABASE_FEED",
            "type": "string",
            "value": "database_feed"
            "description": "The Database feed instance",
            "required": true,
        }
    ]
}
```

Note that, the precedence order of reading and evaluating parameter values is:

1. `wskdeploy` CLI with `--param` with string interpolated using environment variables.
2. Deployment file with string interpolated using environment variables.
3. Manifest file with string interpolated using environment variables.






