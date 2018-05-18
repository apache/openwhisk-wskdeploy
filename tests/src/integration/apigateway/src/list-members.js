// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

/**
 * Return a list of members in the book store.
 */
function main(params) {
    return new Promise(function(resolve, reject) {
        var message = 'List of members in the book store: '
        console.log(message);
        resolve({
            result: {"name":"Anne Li", "name":"Bob Young"}
        });
        return;
    });
}
