// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//
//  Post Slack Message
//
//  @param apiToken apiToken retrieved from Keystone
//
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var request = require('request');

//var whisk = { console: console };
//main({url: 'https://hooks.slack.com/services/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', channel: '#platform-alerts', username: 'WhiskBot', icon: ':exclamation:', title: 'title', content: 'content'});

function main(msg) {
    var postPromise = post(msg.url,
         msg.channel, {
         username : msg.username || 'whisk',
         icon : msg.icon || ':openwhisk:',
         status : msg.status || 'good',
         title : msg.title || '',
         content : msg.content || ''
    });

    return postPromise;
}

function colorForStatus(status) {
    if (status === 'good' || status === 'danger')
        return status;
    else if (status === 'warning')
        return '#FF9900';
    else
        return '#3333FF';
}

function post(url, channel, msg) {
    var form = {
        channel : channel,
        username : msg.username,
        icon_emoji : msg.icon,
        mrkdwn_in : [ 'fields' ],
        attachments : [ {
            fallback : msg.content,
            color : colorForStatus(msg.status),
            fields : [ {
                title : msg.title,
                value : msg.content,
                short : false
            } ]
        } ]
    };

    return new Promise(function(resolve, reject) {
      request.post({
          url : url,
          formData : {
              payload : JSON.stringify(form)
          }
      }, function(error, response, body) {
          if (!error && response.statusCode == 200 && body == 'ok') {
              console.log('posted', msg.content, 'to slack');
              resolve();
          } else {
              console.log('[error]', error, body);
              reject(error);
          }
      });
    });
}
