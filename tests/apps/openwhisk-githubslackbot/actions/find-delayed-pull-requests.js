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
  *			"cloudant_package": "xxxx",
  *			"github_username": "xxxx",
  *			"github_access_token": "xxxx",
  *		}
  *
  * @return which must be a JSON object.
  *         It will be the output of this action.
  *
  */

// require the OpenWhisk npm package
var openwhisk = require("openwhisk");

// global variable for openwhisk object and cloudant package
var wsk;
var packageName;
//global variable for GitHub username and access token
var githubUsername;
var githubAccessToken;

// predetermined threshold for pull requests duration
// pull requests needs attention if they are older than this threshold   
var limits = {
    "READY": {
        amount: 3,
        unit: "days"
    },
    "REVIEW": {
        amount: 4,
        unit: "days"
    }
};

function main(params) {
	// instantiate the openwhisk instance before you can use it
	wsk = openwhisk();
	
	// read Params
	var cloudantPackage = params["cloudant_package"];
	githubUsername = params["github_username"];
	githubAccessToken = params["github_access_token"];
	
	// validate cloudant package is set in params
    if (typeof cloudantPackage === "undefined" || cloudantPackage === null) {
    	return {
    		"error": "Cloudant package is not specified. Please set \"cloudant_package\" bound parameter."
    	};
	}

	// validate github username is set in params
    if (typeof githubUsername === "undefined" || githubUsername === null) {
    	return {
    		"error": "GitHub username is not specified. Please set \"github_username\" bound parameter."
    	};
	}

	// validate github access token is set in params
    if (typeof githubAccessToken === "undefined" || githubAccessToken === null) {
    	return {
    		"error": "GitHub access token is not specified. Please set \"github_access_token\" bound parameter."
    	};
	}

    // access namespace as environment variables
    var namespace = process.env["__OW_NAMESPACE"];

    // Cloudant package can be accessed using /namespace/package
    packageName = "/" + namespace + "/" + cloudantPackage;

	// get list of pull requests from cloudant database
    return wsk.actions.invoke({
    	actionName: packageName + "/list-documents",
        params: { "include_docs": true },
        blocking: true
    })
    .then(activation => {
    	console.log("Found " + activation.response.result.total_rows + " docs.");
    	var listOfIDs = activation.response.result.rows.map(function (row) {
    		return row.id;
    	});
        return listOfIDs;
    })
    .then(function (listOfIds) {
    	return Promise.all(listOfIds.map(getExistingDocument));
    })
    .then(function (trackedPrDocs) {
    	// filter to only PRs that are "too old"
        return trackedPrDocs.filter(prIsTooOld);
    })
    .then(function (oldPrDocs) {
    	// filter to only PRs that are still open
        // because of the undefined order we receive Github events, it is
        // possible that we are still tracking a PR that has since been closed.
        var stillOpenPromises = oldPrDocs.map(isStillOpen);
        return Promise.all(stillOpenPromises)
        	.then(function (isPrOpenArray) {
            	var delayedPRs = oldPrDocs.filter(function (prDoc, index) {
                	return isPrOpenArray[index];
                });
                return {
                	prs: delayedPRs
                };
            });
        });
}

function isStillOpen(prDoc) {
    // fetch updated record from GitHub
    return fetchPrFromGithub(prDoc)
        .then(function (latest) {
            if (latest.state === "closed") {
                console.log("PR#" + prDoc.pr.number + " is now closed, but still being tracked - deleting it from Cloudant.");
                // delete from Cloudant - don"t bother waiting for result
                stopTracking(prDoc.pr);
                return false;
            } else {
                console.log("PR#" + prDoc.pr.number + " is still open.");
                return true;
            }
        });
}

// read document from cloudant data store
// return the doc if it exists
function getExistingDocument(id) {
    return wsk.actions.invoke({
        actionName: packageName + "/read-document",
        params: { "docid": id },
        blocking: true,
    })
    .then(activation => {
          console.log("Found a document in database with ID " + id);
          return activation.response.result;
    })
    // it could be possible that the doc with this ID doesn"t exist and
    // therefore return "undefined" instead of exiting with error
    .catch(function (err) {
        console.log("Error fetching document from database for: " + id);
        console.log(err)
        return undefined;
    });
}

function stopTracking(pullRequest, ifInState) {
    var id = pullRequest["html_url"];
    // get the existing doc
    // if it matches the ifInState, then delete it and
    // stop tracking this pull request
    return getExistingDocument(id)
        .then(function (existingDoc) {
            if (existingDoc) {
                if (!ifInState || existingDoc.state === ifInState) {
                    return wsk.invoke.actions({
                        actionName: packageName + "/delete-document",
                        params: {
                            docid: existingDoc._id,
                            docrev: existingDoc._rev
                        }
                    })
                    .then (function () {
                        return {
                            message: "Sucessfully stopped tracking " + id
                        };
                    });
                } else {
                    return {
                        message: "Refusing to delete doc because it is not in state " + ifInState + " " + id
                    };
                }
            } else {
                return {
                    message: "Refusing to delete doc because it does not exist " + id
                };
            }
        });
}

function fetchPrFromGithub(prDoc) {
    var authorizationHeader = "Basic " + new Buffer(githubUsername + ":" + githubAccessToken).toString("base64");
    var options = {
        method: "GET",
        url: prDoc.pr.url,
        json: true,
        headers: {
            "Content-Type": "application/json",
            "Authorization": authorizationHeader,
            "User-Agent": githubUsername
        }
    };
    var request = require("request");
    return new Promise(function (resolve, reject) {
        request(options, function (error, response, body) {
            if (error) {
                reject({
                    response: response,
                    error: error,
                    body: body
                });            
            } else {
            	if (response.statusCode == 200) {
            		resolve(body);
            	} else {
                    reject({
                        statusCode: response.statusCode,
                        response: body
                    });
            	}
            }
        });
    });
}

function prIsTooOld(prDoc) {
    var moment = require("moment");
    // read lastUpdate from github
    var readyMoment = moment(prDoc.lastUpdate);
	// depeneding on the state of pull request, "READY" or "REVIEW"
	// read the limit amount and days 
    var limit = limits[prDoc.state];
	// moment.diff() returns difference between today and
	// when pull request was last updated (in days as limit.unit is days)
	// return true if the pull request was updated certain (limit.amount) days ago
    return (moment().diff(readyMoment, limit.unit) >= limit.amount);
}
