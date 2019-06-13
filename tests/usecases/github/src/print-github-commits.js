/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
  *
  * main() will be invoked when you Run This Action.
  *
  * @param OpenWhisk actions accept a single parameter,
  *         which must be a JSON object.
  *
  * In this case, the params variable will look like:
  *         { "message":  Webhook POST payload}
  *
  * @return which must be a JSON object.
  *         It will be the output of this action.
  *         returns commit history
  *
  */
function main(params) {

    console.log("Display GitHub Commit Details for GitHub repo: ", params.repository.url);
    for (var commit of params.commits) {
        console.log(params.head_commit.author.name + " added code changes with commit message: " + commit.message);
    }

    console.log("Commit logs are: ")
    console.log(params.commits)

    return { message: params };
}
