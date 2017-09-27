// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

/*
 * Hello, world. Plus more
 */
function main(params) {
    msg = "Hello, " + params.name + " from " + params.place;
    family = "You have " + params.children + " children ";
    stats = "and are " + params.height + " m. tall.";
    return { greeting:  msg, details: family + stats };
}
