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

This Package named `relationships` includes:
- A manifest.yaml which represents library project LIB
- An ext_manifest.yml
- relationships.yml that describes projects relationships


### How to deploy and test

#### Step 1. Deploy the LIB package.

```
$ wskdeploy -m tests/usecases/relationships/manifest.yml --managed
```

#### Step 2. Verify the assets been installed.

e.g. 
```
$ wsk package get lib_package
{
    "name": "lib_package",
    "annotations": [
        {
            "key": "managed",
            "value": {
                "__OW_FILE": "tests/usecases/relationships/manifest.yml",
                "__OW_PROJECT_HASH": "34d92ecfcf71c5ec2d3c2936ae9415c755b05158",
                "__OW_PROJECT_NAME": "LIB"
            }
        }
    ],
    "actions": [
        {
            "name": "lib_greeting_2",
            "annotations": [
                {
                    "key": "managed",
                    "value": {
                        "__OW_FILE": "tests/usecases/relationships/manifest.yml",
                        "__OW_PROJECT_HASH": "34d92ecfcf71c5ec2d3c2936ae9415c755b05158",
                        "__OW_PROJECT_NAME": "LIB"
                    }
                },
                {
                    "key": "exec",
                    "value": "nodejs:6"
                }
            ]
        },
        {
            "name": "lib_greeting_1",
            "version": "0.0.1",
            "annotations": [
                {
                    "key": "managed",
                    "value": {
                        "__OW_FILE": "tests/usecases/relationships/manifest.yml",
                        "__OW_PROJECT_HASH": "34d92ecfcf71c5ec2d3c2936ae9415c755b05158",
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
$ wskdeploy -m tests/usecases/relationships/ext_manifest.yml -r tests/usecases/relationships/relationships.yml --managed
```

#### Step 4. Verify the assets relationships been updated.

e.g. 
```
$ wsk package get lib_package
ok: got package lib_package
{
    "annotations": [
        {
            "key": "managed",
            "value": {
                "__OW_FILE": "tests/usecases/relationships/manifest.yml",
                "__OW_PROJECT_HASH": "34d92ecfcf71c5ec2d3c2936ae9415c755b05158",
                "__OW_PROJECT_NAME": "LIB"
            }
        },
        {
            "key": "managed_list",
            "value": [
                {
                    "__OW_FILE": "tests/usecases/relationships/relationships.yml",
                    "__OW_PROJECT_HASH": "512917c4df9a0acaa0aa339801526fb810b0ee65",
                    "__OW_PROJECT_NAME": "EXT"
                }
            ]
        }
    ],
    "actions": [
        {
            "name": "lib_greeting_2",
            "annotations": [
                {
                    "key": "managed",
                    "value": {
                        "__OW_FILE": "tests/usecases/relationships/manifest.yml",
                        "__OW_PROJECT_HASH": "34d92ecfcf71c5ec2d3c2936ae9415c755b05158",
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
                            "__OW_PROJECT_HASH": "512917c4df9a0acaa0aa339801526fb810b0ee65",
                            "__OW_PROJECT_NAME": "EXT"
                        }
                    ]
                }
            ]
        },
        {
            "name": "lib_greeting_1",
            "annotations": [
                {
                    "key": "managed",
                    "value": {
                        "__OW_FILE": "tests/usecases/relationships/manifest.yml",
                        "__OW_PROJECT_HASH": "34d92ecfcf71c5ec2d3c2936ae9415c755b05158",
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

#### Step 5. Export EXT project

```
$ wskdeploy export -p EXT -m exported_ext.yaml --managed

explore the exported_ext.yaml manifest file and notice both EXT and LIB project assets are there without any notion of LIB project
```