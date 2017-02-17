////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//  Get Servers from Nova
//
//  Will retrieve a list of filtered servers from Nova.  
//
//  @param apiToken apiToken retrieved from Keystone
//
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////


var request = require('request')

function main(params) {

    if (!params.apiToken || !params.host || !params.port || !params.path || !params.protocol) {
        return {"msg" : "Missing required params"}
    }

    console.log('Got params ' + JSON.stringify(params))

    var apiToken = params.apiToken
    var host = params.host
    var port = params.port
    var path = params.path
    var protocol = params.protocol


    var context = null
    if (params.context) {
        context = params.context
    }
   
    var headers = {
        'X-Auth-Token':apiToken
    };
 
    var apiEndpoint = protocol+'//'+host+':'+port+path
   
   var url = apiEndpoint+'/servers'

    // hardcode for now
    var options = {
        method: 'GET',
        headers: headers
    };

    
    if (context && context.queryString) {
        options.qs = context.queryString
    }

    console.log('Options are : '+ JSON.stringify(options))

    return new Promise(function(resolve, reject) {
        request(url, options, function(error, res, body) {
            
            if (error) {
                reject(error)
            }

            var j = JSON.parse(body)

            if (context) {
                j.context = context
            }

            // set this for downstream actions that want compute endpoint
            j.apiEndpoint = apiEndpoint
            j.apiToken = apiToken

            if (j.servers.length > 0) {
                j.serverId = j.servers[0].id
            } 

            console.log(j)
            resolve(j)
        });
    });
    
}
