// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

/*
 * Join the fellowship
 */
function main(params) {
    var member = {name:"", place:"", occupation:"", height:0.0, joined:""};
    name = params.name;
    place = params.place;
    occupation = params.job;
    height = params.height;
    join_date = Date.now();
    return { joined: join_date };
}
