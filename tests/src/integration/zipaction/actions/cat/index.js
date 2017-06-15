function myAction(args) {
    var cat = require('cat');
    cat('https://baidu.com', console.log); 
}
exports.main = myAction;
