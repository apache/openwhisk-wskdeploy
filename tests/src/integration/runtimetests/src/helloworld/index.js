// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

function helloworld(params) {
    var format = require('string-format');
    var name = params.name || 'Stranger';
    payload = format('Hello, {}!', name)
    return { message: payload };
}

exports.main = helloworld;
