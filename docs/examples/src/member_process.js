// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

function formatDate(date) {
  var hours = date.getHours();
  var minutes = date.getMinutes();
  var ampm = hours >= 12 ? 'pm' : 'am';
  hours = hours % 12;
  hours = hours ? hours : 12; // the hour '0' should be '12'
  minutes = minutes < 10 ? '0'+minutes : minutes;
  var strTime = hours + ':' + minutes + ' ' + ampm;
  return date.getMonth()+1 + "/" + date.getDate() + "/" + date.getFullYear() + "  " + strTime;
}

function main(params) {
  console.log("params: " + JSON.stringify(params, null, 4));

  let regionMap = new Map([
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

  console.log("typeof=" + typeof params.member )

  if(params.member && typeof params.member === "object"){
    member = params.member;
    member.region = regionMap.get(member.place) || "unknown";
    member.date = formatDate(new Date(member.joined))
  }
  else
    throw new Error("Invalid parameter: 'member'. type="+typeof(member)+" (expected object).")

  console.log("member: " + JSON.stringify(member, null, 4));
  return { member: member };
}
