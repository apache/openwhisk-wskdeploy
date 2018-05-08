// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

/*
 * This action is meant to capture the output from the basic Hello, world
 * action and improve the greeting. It expects params to contain a greeting
 * and just adds more to it.
 */
function main(params) {
    msg = params.greeting + ", have a good day!";
    return { greeting:  msg };
}
