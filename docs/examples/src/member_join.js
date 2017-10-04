// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

/*
 * Join the fellowship
 */
function main(params) {
    // var scripts = document.getElementsByTagName('script');
    // var lastScript = scripts[scripts.length-1];
    // var scriptName = lastScript.src;
    var member = {name:"", city:"", region:"", occupation:"", joined:"", organization:"", item:"" };

    // Fill in a member record from parameters
    member.name = params.name;
    member.city = params.place;
    member.occupation = params.job;

    // Note the current timestamp when we created the member record
    member.joined = Date.now();

    // The organization being joined is fixed
    member.organization = "fellowship"

    console.log("member=" + member);
    return { member: member };
}
