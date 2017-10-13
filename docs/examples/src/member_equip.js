// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.
const equipmentMap = new Map([
    ['gentleman', 'ring'],
    ['wizard', 'staff'],
    ['archer', 'bow'],
    ['knight', 'sword'],
    ['barbarian', 'club'],
    ['warrior', 'axe'],
    ['thief', 'dagger'],
    ['gardener', 'rope'],
    ['squire', 'shortsword'],
    ['scout', 'horn'],
]);

function main(params) {
    console.log("params: " + JSON.stringify(params, null, 4));

    if(params.member && typeof params.member === "object"){
        // Equip the member based upon their occupation
        member = params.member;
        member.item = equipmentMap.get(member.occupation) || "None";
    }
    else
        throw new Error("Invalid parameter: 'member'. type="+typeof(member) + ", expected object).")

    console.log("member: " + JSON.stringify(member, null, 4));
    return { member: member };
}
