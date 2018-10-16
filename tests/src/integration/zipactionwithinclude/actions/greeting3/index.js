// Licensed to the Apache Software Foundation (ASF) under one or more contributor
// license agreements; and to You under the Apache License, Version 2.0.

/**
 * Return a simple greeting message for someone.
 *
 * @param name A person's name.
 * @param place Where the person is from.
 */


var common1 = require('./common/utils.js')
var common2 = require('./common/common1/utils.js')

function main(params) {
    var name = params.name || params.payload || 'stranger';
    var place = params.place || 'somewhere';
    var hello = common1.hello || common2.hello || 'Hello';
    return {payload:  hello + ', ' + name + ' from ' + place + '!'};
}
exports.main = main;

