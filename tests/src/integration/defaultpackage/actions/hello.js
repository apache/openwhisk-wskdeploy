/*
 * Licensed to the Apache Software Foundation (ASF) under one or more contributor
 * license agreements; and to You under the Apache License, Version 2.0.
 */

/*
 * Return a simple greeting message for the whole world.
 */
function main(params) {
    msg = "Hello, " + params.name + " from " + params.place;
    console.log(msg)
    return { payload:  msg };
}
