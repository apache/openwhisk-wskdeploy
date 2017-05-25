/*
 * Copyright 2017 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * Analyze incoming message and generate a summary as a response
 */
function transform(events) {
  var average = 0;
  for (var i = 0; i < events.length; i++) {
    average += events[i].payload.velocity;
  }
  average = average / events.length;
  var result = {
    "agent": "OpenWhisk action",
    "events_count": events.length,
    "avg_velocity": average
  };
  return result;
}

/**
 * Process incoming message and publish it to Message Hub or Kafka.
 * This code is used as the OpenWhisk Action implementation and linked to a trigger via a rule.
 */
function mhpost(args) {
  console.log("DEBUG: Received message as input: " + JSON.stringify(args));

  return new Promise(function(resolve, reject) {
    if (!args.topic || !args.events || !args.events[0] || !args.kafka_rest_url || !args.api_key)
      reject("Error: Invalid arguments. Must include topic, events[], kafka_rest_url, api_key.");

    // construct CF-style VCAP services JSON
    var vcap_services = {
      "messagehub": [{
        "credentials": {
          "kafka_rest_url": args.kafka_rest_url,
          "api_key": args.api_key
        }
      }]
    };

    var MessageHub = require('message-hub-rest');
    var kafka = new MessageHub(vcap_services);
    var transformedMessage = transform(args.events);
    console.log("DEBUG: Message to be published: " + JSON.stringify(transformedMessage));

    kafka.produce(args.topic, transformedMessage)
      .then(function() {
        resolve({
          "result": "Success: Message was sent to IBM Message Hub."
        });
      })
      .fail(function(error) {
        reject(error);
      });
  });
}

exports.main = mhpost;