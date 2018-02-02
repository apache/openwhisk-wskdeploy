/**
 * Return a simple greeting message for someone.
 *
 * @param name A person's name.
 * @param place Where the person is from.
 */
function main({name:name='Serverless API Gateway'}) {
    return {
      body: new Buffer(JSON.stringify({payload:`Hello world ${name}`})).toString('base64'),
      statusCode: 200,
      headers:{ 'Content-Type': 'application/json'}
    };
}
