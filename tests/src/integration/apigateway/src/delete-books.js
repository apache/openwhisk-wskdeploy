/**
 * Return success saying a book was deleted from the book store.
 */
function main(params) {
    return new Promise(function(resolve, reject) {
        console.log(params.name);

        if (!params.name) {
            console.error('name parameter not set.');
            reject({
                'error': 'name parameter not set.'
            });
            return;
        } else {
            var message = 'A book ' + params.name + ' was deleted from the book store.';
            console.log(message);
            resolve({
                result: message
            });
            return;
        }
    });
}