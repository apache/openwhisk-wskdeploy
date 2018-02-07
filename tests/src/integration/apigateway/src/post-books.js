/**
 * Return success saying a book was added into the book store.
 */
function main(params) {
    return new Promise(function(resolve, reject) {
        console.log(params.name);
        console.log(params.isbn);

        if (!params.name) {
            console.error('name parameter not set.');
            reject({
                'error': 'name parameter not set.'
            });
            return;
        } else {
            var message = 'A book ' + params.name + ' was added to the book store with ISBN ' + params.isbn;
            console.log(message);
            resolve({
                result: message
            });
            return;
        }
    });
}