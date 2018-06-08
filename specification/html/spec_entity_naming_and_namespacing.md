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

## Naming and Namespacing

### Namespacing

Every OpenWhisk entity (i.e., Actions, Feeds, Triggers), including
packages, belongs in a *namespace.*

The fully qualified name of any entity has the format:

```
/<namespaceName>[/<packageName>]/<entityName>
```

The namespace is typically provided at bind-time by the user deploying the package to their chosen OpenWhisk platform provider.

#### Requirements

-   The “/whisk.system“ namespace is reserved for entities that are distributed with the OpenWhisk system.

### Entity Names

The names of all entities, including actions, triggers, rules, packages,
and namespaces, are a sequence of characters that follow the following
format:

-   The first character SHALL be an alphanumeric character, a digit, or an underscore.
-   The subsequent characters MAY be alphanumeric, digits, spaces, or any of the following:
```
_, @, ., -
```
- The last character SHALL NOT be a space.
- The maximum name length of any entity name is 256 characters (i.e., ENTITY_NAME_MAX_LENGTH = 256).

Valid entity names are described with the following regular expression (Java metacharacter syntax):

```
"\A([\w]|[\w][\w@ .-]{0,${ENTITY_NAME_MAX_LENGTH - 2}}[\w@.-])\z"
```

<!--
 Bottom Navigation
-->
---
<html>
<div align="center">
<a href="../README.md#index">Index</a>
</div>
</html>
