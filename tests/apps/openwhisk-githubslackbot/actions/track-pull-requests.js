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
  *            "cloudant_package": "xxxx",
  *            "action": "xxxx",
  *            "pull_request": {
  *                "html_url": "xxxx",
  *                "state": "xxxx",
  *                "number": "xxxx",
  *                "updated_at": "xxxx",
  *                "base": {
  *                     "repo": {
  *                         "full_name": "xxxx"
  *                     }
  *                }
  *            },
  *            "label": {
  *                 "name": "xxxx"
  *            }
  *      }
  * except cloudant_package, rest of the params are
  * sent with webhook POST request.
  * cloudant_package should be set as a bound parameter while
  * deploying this action.
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

function main(params) {
    // instantiate the openwhisk instance before you can use it
    wsk = openwhisk();

    // read Params
    var cloudantPackage = params["cloudant_package"];
    var pullRequest = params["pull_request"];
    var action = params["action"];

    // validate cloudant package is set in params
    if (typeof cloudantPackage === 'undefined' || cloudantPackage === null) {
        return "Cloudant package is not specified. Please set \"cloudant_package\" bound parameter.";
	}

    // access namespace as environment variables
    var namespace = process.env["__OW_NAMESPACE"];

    // Cloudant package can be accessed using /namespace/package
    packageName = "/" + namespace + "/" + cloudantPackage;

    console.log("Action: " + action);
    console.log("Cloudant Package Name: " + packageName);

    // pull request labels
    var ready = "ready";
    var review = "review";

    // action is set to "closed" for closed pull request
    // stop tracking pull request if its closed
    // delete it from the datastore
    if (action === "closed") {
        return stopTracking(pullRequest);
    }

    // make sure pull request is still open
    if (pullRequest.state !== "closed") {
        // when pull request is labeled to either "ready" or "review"
        if (action === "labeled") {
            if (params.label.name === ready) {
                console.log("PR#" + pullRequest.number + " is now " + ready + ".");
                return track(pullRequest, ready.toUpperCase());
            } else if (params.label.name === review) {
                console.log("PR#" + pullRequest.number + " is now in " + review + ".");
                return track(pullRequest, review.toUpperCase());
            }
        // when pull request label ("ready" or "review") is removed
        } else if (action === "unlabeled") {
            if (params.label.name === ready) {
                console.log("PR#" + pullRequest.number + " is no longer " + ready + ".");
                return stopTracking(pullRequest, ready.toUpperCase());
            } else if (params.label.name === review) {
                console.log("PR#" + pullRequest.number + " is no longer in " + review + ".");
                return stopTracking(pullRequest, review.toUpperCase());
            }
        } else if (params.action === "synchronize") {
            return refreshLastUpdate(pullRequest);
        }
    } else {
        console.log("Received event for closed PR#" + pullRequest.number + ". That\'s curious...\n" + params);
    }

    return {
        message: "[" + pullRequest.base.repo["full_name"] + "] PR#" + pullRequest.number + " - No interesting changes"
    };}

// read document from cloudant data store
// return the doc if it exists
function getExistingDocument(id) {
       return wsk.actions.invoke({
        actionName: packageName + "/read-document",
        params: { "docid": id },
        blocking: true,
    })
    .then(activation => {
          console.log("Found pull request in database with ID " + id);
          return activation.response.result;
    })
    // it could be possible that the doc with this ID doesn't exist and
    // therefore return "undefined" instead of exiting with error
    .catch(function (err) {
        console.log("Error fetching pull request from database for: " + id);
        console.log(err)
        return undefined;
    });
}

function refreshLastUpdate(pullRequest){
    var id = pullRequest["html_url"]
    // get the existing doc, set lastUpdate to pullRequest["updated_at"]
    return getExistingDocument(id)
        .then(function (existingDoc) {
            if (existingDoc && existingDoc.pr) {
                existingDoc.lastUpdate = id;
                console.log("Refreshing lastUpdate: " + id);
                return wsk.actions.invoke({
                    actionName: packageName + "/update-document",
                    params: {
                        doc: existingDoc
                    }
                });
            } else {
                return {
                    message: "Not refreshing lastUpdate because PR is not tracked"
                };
            }
        })
        .catch(function (err) {
            console.log("Error fetching pull request from DB: " + id)
            return err;
        });
}
// write a document in cloudant data store
// updates an existing doc if it exist
function writeDocument (doc) {
       return wsk.actions.invoke({
        actionName: packageName + "/create-document",
        params: { "doc": doc },
        blocking: true,
    })
    .then(activation => {
          console.log("Created new document with ID " + activation.response.result.id);
          return activation.response.result;
    })
    .catch(function (err) {
        console.log("Error creating document");
        return err;
    });
}

function track(pullRequest, state) {
    var doc = undefined;
    return getExistingDocument(pullRequest["html_url"])
        .then(function (existingDoc) {
            if (existingDoc) {
                // update pull request if it exists
                existingDoc.state = state;
                existingDoc.pr = pullRequest;
                existingDoc.lastUpdate = pullRequest["updated_at"];
                // createDocument updates doc if _rev and _id exist
                doc = existingDoc;
            } else {
                // the doc does not exist, create one
                doc = {
                    "_id": pullRequest["html_url"],
                    "pr": pullRequest,
                    "state": state,
                    "lastUpdate": pullRequest["updated_at"]
                };
                console.log("Here is the new document")
                console.log(doc)
            }
            return writeDocument(doc);
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
                    return wsk.actions.invoke({
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
