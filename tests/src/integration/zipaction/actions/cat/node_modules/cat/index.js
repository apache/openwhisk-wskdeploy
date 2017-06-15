var fs = require('fs');
var parse = require('url').parse;

var catter = function(lib) {
	var cat = function(url, callback) {	
		if (typeof url === 'string') {
			url = parse(url);
		}
		lib.get({host:url.hostname, port:url.port, path:url.pathname}, function(response) {
			if (/3\d\d/.test(response.statusCode) && response.headers.location) {
				cat(parse(response.headers.location), callback);
				return;
			}
			if (!(/2\d\d/).test(response.statusCode)) {
				callback(new Error('non 2xx status code: ' + response.statusCode));
				return;
			}
			var buffer = '';
		
			response.setEncoding('utf-8');
			response.on('data', function(data) {
				buffer += data;
			});
			response.on('close', function() {
				callback(new Error('unexpected close of response'));
			});
			response.on('end', function() {
				callback(null, buffer);
			});
		}).on('error', callback);
	};
	return cat;
};

var http = catter(require('http'));
var https = catter(require('https'));

module.exports = function(location, callback) {
	var protocol = (location.match(/^(\w+):\/\//) || [])[1] || 'file';

	if (protocol === 'file') {
		fs.readFile(location.replace(/^(file:\/\/localhost|file:\/\/)/, ''), 'utf-8', callback);
		return;
	}
	if (protocol === 'http') {
		http(location, callback);
		return;
	}
	if (protocol === 'https') {
		https(location, callback);
		return;
	}
	throw new Error('protocol '+protocol+' currently not supported :(');
};