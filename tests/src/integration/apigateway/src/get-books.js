/**
 * Return a simple greeting message for someone.
 *
 * @param name A person's name.
 * @param place Where the person is from.
 */
function main() {
    return {
      body: new Buffer(JSON.stringify({result:[{"name":"JavaScript: The Good Parts", "isbn":"978-0596517748"}]})).toString('base64'),
      statusCode: 200,
      headers:{ 'Content-Type': 'application/json'}
    };
}
