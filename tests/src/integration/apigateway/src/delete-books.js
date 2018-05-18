/*
 * Licensed to the Apache Software Foundation (ASF) under one or more contributor
 * license agreements; and to You under the Apache License, Version 2.0.
 */

/**
 * Return success saying a book was deleted from the book store.
 */
function main(params) {
    return new Promise(function(resolve, reject) {
        console.log(params.name);

        if (!params.name) {
            console.error('name parameter not set.');
            reject({
                'error': 'name parameter not set.'
            });
            return;
        } else {
            var message = 'A book ' + params.name + ' was deleted from the book store.';
            console.log(message);
            resolve({
                result: message
            });
            return;
        }
    });
}
