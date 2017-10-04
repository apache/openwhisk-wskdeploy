// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

// function formatDate(date) {
//   var hours = date.getHours();
//   var minutes = date.getMinutes();
//   var ampm = hours >= 12 ? 'pm' : 'am';
//   hours = hours % 12;
//   hours = hours ? hours : 12; // the hour '0' should be '12'
//   minutes = minutes < 10 ? '0'+minutes : minutes;
//   var strTime = hours + ':' + minutes + ' ' + ampm;
//   return date.getMonth()+1 + "/" + date.getDate() + "/" + date.getFullYear() + "  " + strTime;
// }

function main(params) {

  let regionMap = new Map([
    ['Shire', 'Eriador'],
    ['the Shire', 'Eriador'],
    ['Hobbiton', 'Eriador']
  ]);

  member = params.member;
  if(member)
  {
    if( member.city)
    {
     // the Western regions of M.E. contained the lands of Eriador, Gondor, the Misty Mountains,
     // and the vales of the river Anduin
     if( member.city === "Shire" || member.city === "the Shire")
        member.region = "Eriador";
    }
  }
  console.log(JSON.stringify(member, null, 4));
  return {member: member, region: member.region};
}
