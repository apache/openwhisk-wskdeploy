#!/usr/bin/env node

var fs = require('fs')

var input = process.argv.slice(2)

var loop = function() {
  if (!input.length) return
  var next = input.shift()
  var s = next === '-' ? process.stdin : fs.createReadStream(next)
  s.on('end', loop).pipe(process.stdout)
}

loop()