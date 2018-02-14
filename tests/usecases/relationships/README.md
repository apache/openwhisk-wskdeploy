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

# Relationships use case for wskdeploy

### Package description

This Package named `relatinships` includes:
- A package (manifest) based on the one in triggerrule usecase. This package it represents the library project LIB
- A manifest ext_manifest.yml extending the LIB project
- relationships.yml describing projects relationships

### How to deploy and test

#### Step 1. Deploy the LIB package.

```
$ wskdeploy -m tests/usecases/triggerrule/manifest.yml --managed
```

#### Step 2. Verify the assets been installed.

e.g. 
```
$ wsk package get triggerrule
{
    "name": "triggerrule", 
    "annotations": [
        {
            "key": "managed",
            "value": {
                "__OW_FILE": "tests/usecases/relationships/manifest.yml",
                "__OW_PROJECT_HASH": "420e6ac07c0965a488a40770ecff663f0880ad78",
                "__OW_PROJECT_NAME": "LIB"
            }
        }
    ],
    "actions": [
        {
            "name": "greeting",
            "annotations": [
                {
                    "key": "managed",
                    "value": {
                        "__OW_FILE": "tests/usecases/relationships/manifest.yml",
                        "__OW_PROJECT_HASH": "420e6ac07c0965a488a40770ecff663f0880ad78",
                        "__OW_PROJECT_NAME": "LIB"
                    }
                },
                {
                    "key": "exec",
                    "value": "nodejs:6"
                }
            ]
        }
    ]
}
```

#### Step 3. Deploy the EXT package.

```
$ wskdeploy -m tests/usecases/triggerrule/ext_manifest.yml -r tests/usecases/triggerrule/relationships.yml --managed
```

#### Step 4. Verify the assets relationships been updated.

e.g. 
```
$ wsk package get triggerrule
root@hslt1:~# wsk package get triggerrule
ok: got package triggerrule
{
    "name": "triggerrule",
    "annotations": [
        {
            "key": "managed",
            "value": {
                "__OW_FILE": "tests/usecases/relationships/manifest.yml",
                "__OW_PROJECT_HASH": "420e6ac07c0965a488a40770ecff663f0880ad78",
                "__OW_PROJECT_NAME": "LIB"
            }
        },
        {
            "key": "managed_list",
            "value": [
                {
                    "__OW_FILE": "tests/usecases/relationships/relationships.yml",
                    "__OW_PROJECT_HASH": "cb27e577748016b28d9a84b04a1702fe1422c260",
                    "__OW_PROJECT_NAME": "EXT"
                }
            ]
        }
    ],
    "actions": [
        {
            "name": "greeting",
            "annotations": [
                {
                    "key": "managed",
                    "value": {
                        "__OW_FILE": "tests/usecases/relationships/manifest.yml",
                        "__OW_PROJECT_HASH": "420e6ac07c0965a488a40770ecff663f0880ad78",
                        "__OW_PROJECT_NAME": "LIB"
                    }
                },
                {
                    "key": "exec",
                    "value": "nodejs:6"
                },
                {
                    "key": "managed_list",
                    "value": [
                        {
                            "__OW_FILE": "tests/usecases/relationships/relationships.yml",
                            "__OW_PROJECT_HASH": "cb27e577748016b28d9a84b04a1702fe1422c260",
                            "__OW_PROJECT_NAME": "EXT"
                        }
                    ]
                }
            ]
        }
    ]
}

```

#### Step 5. Export EXT project

```
$ wskdeploy export -p EXT -m exported_ext.yaml --managed

explore the exported_ext.yaml manifest file and notice all the required EXT and LIB project assets are there
```