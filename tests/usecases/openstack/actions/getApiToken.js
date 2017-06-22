// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//  Get Auth Token from Keystone
//
//  Will retrieve the auth token and storage URL for the object store service.  Uses V3 of the auth calls
//
//  The params object should look like:
//
//  "userId":"f45adsa0d0478c", "password":"sd2rsS^kdfsd", "projectId":"11fdseerff"}
//
//  @param userId: user id (not user name)
//  @param password: password
//  @param projectId: project/tenantId
//  @param host: hostname of keystone endpoint (don't include 'https://')
//  @param port: port of keystone endpoint
//  @param endpointName: the name of the public endpoint service you want (nova, swift, etc)
//
//
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


var request = require('request')
var url = require('url')

function main(params) {

    if (!params.userId || !params.password || !params.projectId || !params.host || !params.port || !params.endpointName) {
        return {"msg" : "Missing required parameters"}
    }

    console.log('getApiToken. Params '+JSON.stringify(params))

    var userId = params.userId
    var password = params.password
    var projectId = params.projectId
    var pathPrefix = ""
    var host = params.host
    var port = params.port
    var endpointName = params.endpointName

    // customize your path here.
    if (params.pathPrefix) {
        pathPrefix = params.pathPrefix
    }

    var endpointUrl = 'https://' + host + ':' + port + '/v3/auth/tokens'

    var context = null
    if (params.context) {
        context = params.context
    }


    var body = {
    "auth": {
        "identity": {
            "methods": [
                "password"
            ],
            "password": {
                "user": {
                    "id": userId,
                    "password": password
                }
            }
        },
        "scope": {
            "project": {
                "id": projectId
            }
        }
    }
}

    var headers = {
       'Content-Type': 'application/json',
       'Content-Length': Buffer.byteLength(JSON.stringify(body), ['utf8'])
    };

    var options = {
        method: 'POST',
        headers: headers,
        json: body
    };

    console.log("got options "+ JSON.stringify(options))

    return new Promise(function(resolve, reject) {
        request(endpointUrl, options, function (error, response, body) {
            //console.log(body)
            if (!error) {

                var authToken = response.headers['x-subject-token']

                var j = body
                var entries = j.token.catalog

                for (var i = 0; i < entries.length; i++) {
                    var entry = entries[i]
                    console.log('Comparing '+ entry.name + " with "+endpointName)
                    if (entry.name === endpointName) {
                        var endpoints = entry.endpoints

                        console.log('Got endpoints '+endpoints)
                        for (var j = 0; j < endpoints.length; j++) {
                            var endpoint = endpoints[j]
                            if (endpoint.interface === 'public') {
                                console.log('Public endpoint is ' + endpoint.url)
                                console.log('Auth token is ' + authToken)

                                var urlParts = url.parse(endpoint.url,true)

                                var jsonResponse = {apiToken: authToken,
                                    endpointUrl: endpoint.url,
                                    host: urlParts.hostname,
                                    port: urlParts.port,
                                    path: urlParts.path,
                                    protocol: urlParts.protocol}

                                if (context) {
                                    jsonResponse.context = context
                                }

                                return resolve(jsonResponse)
                            }
                        }
                    }
            }

            reject({'msg': 'Cannot find public endpoint in response from keystone'})

            } else {
                reject(error)
            }
        });

    });
}
