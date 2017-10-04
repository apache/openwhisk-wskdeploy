// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

/*
 * Join the fellowship
 */
function main(params) {
    console.log("params: " + JSON.stringify(params, null, 4));
    var member = {name:"", place:"", region:"", occupation:"", joined:"", organization:"", item:"" };

    // Fill in a member record from parameters
    member.name = params.name;
    member.place = params.place;
    member.occupation = params.job;

    // Note the current timestamp when we created the member record
    member.joined = Date.now();

    // The organization being joined is fixed
    member.organization = "fellowship";

    console.log("member: " + JSON.stringify(member, null, 4));
    return { member: member };
}
