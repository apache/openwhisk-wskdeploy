/*
 * Licensed to the Apache Software Foundation (ASF) under one or more contributor
 * license agreements; and to You under the Apache License, Version 2.0.
 */

var openwhisk = require('openwhisk');

/**
 * Analyze incoming message and generate a summary as a response
 */
function transform(events) {
    var average = 0;
    for (var i = 0; i < events.length; i++) {
        average += events[i].payload.category;
    }
    average = average / events.length;
    var result = {
        "agent": "OpenWhisk action",
        "events_count": events.length,
        "avg_category": average
    };
    return result;
}

/**
 * Process incoming message from the receive-messages action earlier
 * in the sequence and publish a new message to Message Hub.
 */
function main(params) {
    console.log("DEBUG: Received message as input: " + JSON.stringify(params));

    return new Promise(function(resolve, reject) {
        if (!params.topic || !params.messagehub_instance || !params.events || !params.events[0]) {
            reject("Error: Invalid arguments. Must include topic, events[], message hub service name.");
        }

        var transformedMessage = JSON.stringify(transform(params.events));
        console.log("DEBUG: Message to be published: " + transformedMessage);

        openwhisk().actions.invoke({
            name: params.messagehub_instance + '/messageHubProduce',
            blocking: true,
            result: true,
            params: {
                value: transformedMessage,
                topic: params.topic
            }
        }).then(result => {
            resolve({
                "result": "Success: Message was sent to Message Hub."
            });
        }).catch(error => {
            reject(error);
        });

    });
}
