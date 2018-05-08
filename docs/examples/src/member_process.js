// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.
const regionMap = new Map([
    ['Shire', 'Eriador'],
    ['the Shire', 'Eriador'],
    ['Hobbiton', 'Eriador'],
    ['Bree', 'Eriador'],
    ['Rivendell', 'Eriador'],
    ['Minas Tirith', 'Gondor'],
    ['Esgaroth', 'Rhovanion'],
    ['Dale', 'Rhovanion'],
    ['Lake Town', 'Rhovanion'],
    ['Minas Morgul', 'Mordor'],
]);

function main(params) {
    console.log("params: " + JSON.stringify(params, null, 4));

    if(params.member && typeof params.member === "object"){
        // Augment the member (record) created in the previous Action
        member = params.member;
        member.region = regionMap.get(member.place) || "unknown";
        member.date = new Date(member.joined).toLocaleDateString();
        member.time = new Date(member.joined).toLocaleTimeString();
    }
    else
        throw new Error("Invalid parameter: 'member'. type="+typeof(member) + ", expected object).")

    console.log("member: " + JSON.stringify(member, null, 4));
    return { member: member };
}
