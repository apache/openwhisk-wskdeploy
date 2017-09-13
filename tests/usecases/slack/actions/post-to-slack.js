/**
  *
  * main() will be invoked when you Run This Action.
  *
  * @param Whisk actions accept a single parameter,
  *        which must be a JSON object.
  *
  * In this case, the params variable will look like:
  *     {
  *         "message": "xxxx",
  *         "slack_package": "xxxx",
  *     }
  *
  * @return which must be a JSON object.
  *         It will be the output of this action.
  *
  */


function main(params) {
    // require the OpenWhisk npm package
    var openwhisk = require("openwhisk");

    // instantiate the openwhisk instance before you can use it
    wsk = openwhisk();

    //read Params
    var message = params.message;
    var slackPackage = params.slack_package;

    console.log(message);

    // access namespace as environment variables
    var namespace = process.env["__OW_NAMESPACE"];

    // Slack package can be accessed using /namespace/package
    packageName = "/" + namespace + "/" + slackPackage;

    return wsk.actions.invoke({
        actionName: packageName + "/post",
        params: {
            "text": message,
        },
        blocking: true
    })
    .then(activation => {
        console.log("Posted message to slack");
        return {
            message: activation
        };
    })
    .catch(function (err) {
        console.log("Error posting message to slack")
        return {
            error: err
        };
    });
}

