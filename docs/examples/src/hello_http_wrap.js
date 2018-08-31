// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

/*
 * This action take a result from the previous action and wraps it into a
 * HTTP response structure.
 */
function main(params) {
    return {
        body: params,
        statusCode: 200,
        headers: {'Content-Type': 'application/json'}
    };
}
