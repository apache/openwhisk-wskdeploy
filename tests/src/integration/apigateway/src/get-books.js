/*
 * Licensed to the Apache Software Foundation (ASF) under one or more contributor
 * license agreements; and to You under the Apache License, Version 2.0.
 */

/**
 * Return a list of books in the book store.
 */
function main(params) {
    return new Promise(function(resolve, reject) {
        var message = 'List of books in the book store: '
        console.log(message);
        resolve({
            result: {"name":"JavaScript: The Good Parts", "ISBN":"978-0596517748"}
        });
        return;
    });
}
