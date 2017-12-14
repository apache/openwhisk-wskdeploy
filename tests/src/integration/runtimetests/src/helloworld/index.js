function helloworld(params) {
    var format = require('string-format');
    var name = params.name || 'Stranger';
    payload = format('Hello, {}!', name)
    return { message: payload };
}

exports.main = helloworld;
