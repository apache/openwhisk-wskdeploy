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

# Test Case of message hub

This is a test case for message hub. Before running this use case, please make sure you have set the following
environment variables on your machine: MESSAGEHUB_ADMIN_HOST, KAFKA_BROKERS_SASL, SOURCE_TOPIC and DESTINATION_TOPIC.

The environment variables, SOURCE_TOPIC and DESTINATION_TOPIC, are two topic names in message hub service. Both of them
must be available in the message hub service you are about to use. The variable MESSAGEHUB_ADMIN_HOST specifies the url
link of the admin host for the message hub. The variable KAFKA_BROKERS_SASL specifies the array of the kafka brokers, e.g.
[kafka01-prod01.messagehub.services.net:9093 kafka02-prod01.messagehub.services.net:9093].

It can be deployed and tested with:

```bash
$ wskdeploy -p tests/src/integration/message-hub
```
