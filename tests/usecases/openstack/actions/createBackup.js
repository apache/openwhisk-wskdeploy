// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//  Create a snapshot of a server
//
//  @param apiToken apiToken retrieved from Keystone
//
//  http://developer.openstack.org/api-ref/compute/?expanded=create-server-back-up-createbackup-action-detail
//
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


var request = require('request')

function pad(n){return (n < 10? '0' : '') + n;}

function main(params) {

    console.log('createNewSnapshot. got params '+JSON.stringify(params))

     if (!params.apiToken) {
        return {"msg" : "Missing apiToken"}
    }
     if (!params.apiEndpoint) {
        return {"msg" : "Missing apiEndpoint"}
    }
     if (!params.serverId ) {
        return {"msg" : "Missing snapshot serverId"}
    }

    if (!params.context) {
        return {"msg" : "Missing context"}
    }

    if (!params.context.backupName) {
        return {"msg" : "Missing backup name"}
    }

    if (!params.context.backupType) {
        return {"msg" : "Missing backup type"}
    }

    if (!params.context.rotation) {
        return {"msg" : "Missing rotation"}
    }

    var apiEndpoint = params.apiEndpoint
    var apiToken = params.apiToken
    var serverId = params.serverId
    var backupName = params.context.backupName
    var backupType = params.context.backupType
    var rotation = params.context.rotation



    d = new Date();
    d.setUTCHours(d.getUTCHours() - 5);
    var backupDate = d.getUTCFullYear() + '-' + pad(d.getUTCMonth() + 1) + '-' + pad(d.getUTCDate()) + 'T' + pad(d.getUTCHours()) + ':' + pad(d.getUTCMinutes())
    console.log('Backup Date : ' + backupDate)

    var backupNameDate = backupName + '.' + backupDate

    console.log('createNewSnapshot params ' + JSON.stringify(params))

    var context = null
    if (params.context) {
        context = params.context
    }

    var headers = {
        'X-Auth-Token': apiToken,
        'content-type': 'application/json'
    };

    //May be a better way to parse this
    var post_data = {
        createBackup: {
            name: backupNameDate,
            backup_type: backupType,
            rotation: rotation
        }
    };

    var url = apiEndpoint+'/servers/'+serverId+'/action'

    // hardcode for now
    var options = {
        method: 'POST',
        json: post_data,
        headers: headers
    };


    var headers = {
        'X-Auth-Token':apiToken
    };

    console.log('Options are : '+ JSON.stringify(options))

    return new Promise(function(resolve, reject) {
        request(url, options, function(error, res, body) {


            if (error) {
                var strContent = "Error occured in backing up server: " + backupName + "(" + serverId + "). Error: " + JSON.stringify(error)
                reject({ status: "danger",  title: "Backup Report", content: strContent})
            } else if (res.statusCode >= 400) {
                for (var key in body) {
                    if (body.hasOwnProperty(key)) {
                        var strCode = body[key].code
                        var strError = key
                        var strMsg = body[key].message
                    }
                }
                var strContent = "Server: " + backupName + "(" + serverId + "). \nCode: " + strCode + "\nWarning: " + strError + "\nMessage: " + strMsg
                var strStatus = "warning"
            } else {
                var strContent = "Backup complete for server: " + backupName + "(" + serverId + ")"
                var strStatus = "good"
            }

            console.log('Response Code: ' + res.statusCode + ' Body: '+ JSON.stringify(body))

            resolve({ status: strStatus,  title: "Backup Report", content: strContent })
        });

    });

}
