# cat

cat will read the contents of an url. it's available through npm

	npm install cat
	
it will read your files

```js
var cat = require('cat');

cat('myfile.txt', console.log);             // reads the file as utf-8 and returns it output
cat('file://myfile.txt', console.log);      // same as above
```

and your `http` / `https` urls

```js
cat('http://google.com', console.log);      // cat also follows any redirects
cat('https://github.com', console.log);     // and cat read https
cat('http://fail.google.com', console.log); // if a status code != 2xx is received it will 
                                            // call the callback with an error.

```