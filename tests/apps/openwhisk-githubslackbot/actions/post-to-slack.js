/*
 * Copyright 2015-2016 IBM Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
 
 /**
  *
  * main() will be invoked when you Run This Action.
  * 
  * @param Whisk actions accept a single parameter,
  *        which must be a JSON object.
  *
  * In this case, the params variable will look like:
  *     {
  *			"prs": "xxxx",
  *			"slack_package": "xxxx",
  *		}
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
    var prs = params.prs;
	var slackPackage = params.slack_package;

	if (prs.length === 0) {
		return {
            message: "Yeay! No old PRs!"
        };
	}
	
    var messages = ["Hello Whisk Devs! It looks like we have some Pull Requests that may require a little love:\n"];
	
    prs.sort(byDescendingAge);

    for (var i = 0; i < prs.length; i++) {
    	// read pr from a document
        var doc = prs[i];
        var pr = doc.pr;
        console.log("PR:" + doc._id);
        // find out how many days old a pull request is.
        var age = getPrAge(doc, "days");
        var ageString = age + " " + (age === 1 ? "day" : "days");
        var message;
        if(doc.state === "READY") {
        	message = "has been marked \"ready\" for more than " + ageString
        } else {
            message = "has been under \"review\" without comments for more than " + ageString
        }
        messages.push("<" + pr["html_url"] + "|[" + pr.base.repo["full_name"] + "] PR #" + pr.number + "> " + message);
    }

    console.log(messages.join("\n"));

    // access namespace as environment variables
    var namespace = process.env["__OW_NAMESPACE"];
    
    // Slack package can be accessed using /namespace/package
    packageName = "/" + namespace + "/" + slackPackage;
    
    return wsk.actions.invoke({
    	actionName: packageName + "/post",
        params: {
            "text": messages.join("\n"),
        },
        blocking: true
    })
    .then(activation => {
    	console.log("Posted messages to slack");
    	return {
    		message: activation
    	};
    })
    .catch(function (err) {
		console.log("Error posting messages to slack")
		return {
			error: err
		};
    });
}

function byDescendingAge(pr1, pr2) {
    var age1 = getPrAge(pr1, "hours");
    var age2 = getPrAge(pr2, "hours");
    // sort descending
    return age2 - age1;
}

function getPrAge(pr, unit) {
	var moment = require("moment");
	// instantiate moment with last update of a pull request
    var readyMoment = moment(pr.lastUpdate);
    // difference between now and last update in "hours" or "days"
    return moment().diff(readyMoment, unit);
}
