/*
 * Licensed to the Apache Software Foundation (ASF) under one or more contributor
 * license agreements; and to You under the Apache License, Version 2.0.
 */

/**
 * Return success saying a book was updated into the book store.
 */
function main(params) {
    return new Promise(function(resolve, reject) {
        console.log(params.name);
        console.log(params.isbn);

        if (!params.name) {
            console.error('name parameter not set.');
            reject({
                'error': 'name parameter not set.'
            });
            return;
        } else if (!params.isbn) {
            console.error('isbn parameter not set.');
            reject({
                'error': 'isbn parameter not set.'
            });
            return;
        } else {
            var message = 'A book ' + params.name + ' was updated to a new ISBN ' + params.isbn;
            console.log(message);
            resolve({
                result: message
            });
            return;
        }
    });
}
