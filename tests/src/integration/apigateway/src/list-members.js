/**
 * Return a list of members in the book store.
 */
function main(params) {
    return new Promise(function(resolve, reject) {
        var message = 'List of members in the book store: '
        console.log(message);
        resolve({
            result: {"name":"Anne Li", "name":"Bob Young"}
        });
        return;
    });
}
