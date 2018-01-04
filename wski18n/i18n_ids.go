/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package wski18n

const(
	// Debug / trace message prefixes
	ID_MSG_PREFIX_ERROR	= "msg_prefix_error"	// "translation": "Error"
	ID_MSG_PREFIX_SUCCESS	= "msg_prefix_success"	// "translation": "Success"
	ID_MSG_PREFIX_WARNING	= "msg_prefix_warning"	// "translation": "Warning"
	ID_MSG_PREFIX_INFO	= "msg_prefix_info"	// "translation": "Info"

	// wskdeploy Command messages
	ID_MSG_MISSING_KEY_CMD	= "Missing cmd key"	// "translation": "Missing 'cmd' input key"

	// ???
	ID_CMD_FLAG_AUTH_KEY	= "authorization `KEY`"	// "translation": "authorization `KEY`"
	ID_CMD_FLAG_NAMESPACE	= "namespace"		// "translation": "namespace"
	ID_CMD_FLAG_API_HOST	= "whisk API `HOST`"	// "translation": "whisk API `HOST`"
	ID_CMD_FLAG_API_VERSION	= "whisk API `VERSION`"	// "translation": "whisk API `VERSION`"
	ID_CMD_FLAG_KEY_FILE	= "path of the .key file"	// "translation": "path of the .key file"
	ID_CMD_FLAG_CERT_FILE	= "path of the .cert file"	// "translation": "path of the .cert file"

	// etc.
	ID_MANIFEST_FILE_NOT_FOUND_X_path_X = "manifest_not_found_at_path"
)

//{
//"id": "Using {{.manifestPath}} for deployment.\n",
//"translation": "Using {{.manifestPath}} for deployment.\n"
//},
//{
//"id": "Using {{.manifestPath}} for undeployment.\n",
//"translation": "Using {{.manifestPath}} for undeployment.\n"
//},
//{
//"id": "the runtime is not supported by Openwhisk platform.\n",
//"translation": "the runtime is not supported by Openwhisk platform.\n"
//},
//{

//var ID_WARNING_RUNTIME_MISMATCH_X_runtime_X_ext_X_action_X = "go.yaml"

//"id": "WARNING: Runtime specified in manifest YAML {{.runtime}} does not match with action source file extension {{.ext}} for action {{.action}}.\n",
//"translation": "WARNING: Runtime specified in manifest YAML {{.runtime}} does not match with action source file extension {{.ext}} for action {{.action}}.\n"
//},
//{
//"id": "WARNING: Whisk Deploy has chosen appropriate runtime {{.runtime}} based on the action source file extension for that action {{.action}}.\n",
//"translation": "WARNING: Whisk Deploy has chosen appropriate runtime {{.runtime}} based on the action source file extension for that action {{.action}}.\n"
//},
//{
//"id": "WARNING: Runtime specified in manifest YAML {{.runtime}} is not supported by OpenWhisk server for the action {{.action}}.\n",
//"translation": "WARNING: Runtime specified in manifest YAML {{.runtime}} is not supported by OpenWhisk server for the action {{.action}}.\n"
//},
//{
//"id": "Unsupported runtime type, set to nodejs",
//"translation": "Unsupported runtime type, set to nodejs"
//},
//{
//"id": "The authentication key is not configured.\n",
//"translation": "The authentication key is not configured.\n"
//},
//{
//"id": "The API host is not configured.\n",
//"translation": "The API host is not configured.\n"
//},
//{
//"id": "The namespace is not configured.\n",
//"translation": "The namespace is not configured.\n"
//},
//{
//"id": "The API host is {{.apihost}}, from {{.apisource}}.\n",
//"translation": "The API host is {{.apihost}}, from {{.apisource}}.\n"
//},
//{
//"id": "The auth key is set, from {{.authsource}}.\n",
//"translation": "The auth key is set, from {{.authsource}}.\n"
//},
//{
//"id": "The namespace is {{.namespace}}, from {{.namespacesource}}.\n",
//"translation": "The namespace is {{.namespace}}, from {{.namespacesource}}.\n"
//},
//{
//"id": "Failed to get the supported runtimes from OpenWhisk service: {{.err}}.\n",
//"translation": "Failed to get the supported runtimes from OpenWhisk service: {{.err}}.\n"
//},
//{
//"id": "Start to unmarshal Openwhisk info from local values.\n",
//"translation": "Start to unmarshal Openwhisk info from local values.\n"
//},
//{
//"id": "Unmarshal Openwhisk info from internet.\n",
//"translation": "Unmarshal Openwhisk info from internet.\n"
//},
//{
//"id": "Deployment completed successfully.\n",
//"translation": "Deployment completed successfully.\n"
//},
//{
//"id": "Got error creating package with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error creating package with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "Got error creating action with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error creating package with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "Got error creating api with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error creating api with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "Got error creating rule with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error creating rule with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "Got error setting the status of rule with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error setting the status of rule with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "Got error creating trigger with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error creating trigger with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "Got error creating trigger feed with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error creating trigger feed with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "Got error creating package binding with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error creating package binding with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "Deployment of dependency {{.depName}} did not complete sucessfully. Run `wskdeploy undeploy` to remove partially deployed assets.\n",
//"translation": "Deployment of dependency {{.depName}} did not complete sucessfully. Run `wskdeploy undeploy` to remove partially deployed assets.\n"
//},
//{
//"id": "Deploying action {{.output}} ...",
//"translation": "Deploying action {{.output}} ..."
//},
//{
//"id": "Deploying rule {{.output}} ...",
//"translation": "Deploying rule {{.output}} ..."
//},
//{
//"id": "Deploying trigger feed {{.output}} ...",
//"translation": "Deploying trigger feed {{.output}} ..."
//},
//{
//"id": "Deploying package {{.output}} ...",
//"translation": "Deploying package {{.output}} ..."
//},
//{
//"id": "Deploying package binding {{.output}} ...",
//"translation": "Deploying package binding {{.output}} ..."
//},
//{
//"id": "Deploying dependency {{.output}} ...",
//"translation": "Deploying dependency {{.output}} ..."
//},
//{
//"id": "OK. Cancelling deployment.\n",
//"translation": "OK. Cancelling deployment.\n"
//},
//{
//"id": "OK. Canceling undeployment.\n",
//"translation": "OK. Canceling undeployment.\n"
//},
//{
//"id": "Undeployment did not complete sucessfully.\n",
//"translation": "Undeployment did not complete sucessfully.\n"
//},
//{
//"id": "Deployment removed successfully.\n",
//"translation": "Deployment removed successfully.\n"
//},
//{
//"id": "Undeployment did not complete sucessfully.\n",
//"translation": "Undeployment did not complete sucessfully.\n"
//},
//{
//"id": "Undeploying dependency {{.depName}} ...",
//"translation": "Undeploying dependency {{.depName}} ..."
//},
//{
//"id": "Undeployment of dependency {{.depName}} did not complete sucessfully.\n",
//"translation": "Undeployment of dependency {{.depName}} did not complete sucessfully.\n"
//},
//{
//"id": "Got error deleting action with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error deleting action with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "Got error deleting rule with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error deleting rule with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "Got error setting the status of rule with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error setting the status of rule with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "Got error deleting trigger with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error deleting trigger with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "Got error deleting trigger feed with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error deleting trigger feed with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "Got error deleting package with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error deleting package with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "WARNING: The 'source' YAML key in trigger entity is deprecated. Please use 'feed' instead as described in specifications.\n",
//"translation": "WARNING: The 'source' YAML key in trigger entity is deprecated. Please use 'feed' instead as described in specifications.\n"
//},
//{
//"id": "Got error deleting binding package with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Got error deleting binding package with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "Dependency {{.output}} has been successfully deployed.\n",
//"translation": "Dependency {{.output}} has been successfully deployed.\n"
//},
//{
//"id": "Package binding {{.output}} has been successfully deployed.\n",
//"translation": "Package binding {{.output}} has been successfully deployed.\n"
//},
//{
//"id": "Package {{.output}} has been successfully deployed.\n",
//"translation": "Package {{.output}} has been successfully deployed.\n"
//},
//{
//"id": "Trigger {{.output}} has been successfully deployed.\n",
//"translation": "Trigger {{.output}} has been successfully deployed.\n"
//},
//{
//"id": "Trigger feed {{.output}} has been successfully deployed.\n",
//"translation": "Trigger feed {{.output}} has been successfully deployed.\n"
//},
//{
//"id": "Rule {{.output}} has been successfully deployed.\n",
//"translation": "Rule {{.output}} has been successfully deployed.\n"
//},
//{
//"id": "Action {{.output}} has been successfully deployed.\n",
//"translation": "Action {{.output}} has been successfully deployed.\n"
//},
//{
//"id": "Dependency {{.depName}} has been successfully undeployed.\n",
//"translation": "Dependency {{.depName}} has been successfully undeployed.\n"
//},
//{
//"id": "Trigger {{.trigger}} has been removed.\n",
//"translation": "Trigger {{.trigger}} has been removed.\n"
//},
//{
//"id": "Rule {{.rule}} has been removed.\n",
//"translation": "Rule {{.rule}} has been removed.\n"
//},
//{
//"id": "Action {{.action}} has been removed.\n",
//"translation": "Action {{.action}} has been removed.\n"
//},
//{
//"id": "Failed to invoke the feed when deleting trigger feed with error message: {{.err}} and error code: {{.code}}.\n",
//"translation": "Failed to invoke the feed when deleting trigger feed with error message: {{.err}} and error code: {{.code}}.\n"
//},
//{
//"id": "WARNING: Mandatory field Package Version must be set.\n",
//"translation": "WARNING: Mandatory field Package Version must be set.\n"
//},
//{
//"id": "WARNING: Package Version is not saved in the current wskdeploy version.\n",
//"translation": "WARNING: Package Version is not saved in the current wskdeploy version.\n"
//},
//{
//"id": "WARNING: Mandatory field Package License must be set.\n",
//"translation": "WARNING: Mandatory field Package License must be set.\n"
//},
//{
//"id": "WARNING: Package License is not saved in the current wskdeploy version.\n",
//"translation": "WARNING: Package License is not saved in the current wskdeploy version.\n"
//},
//{
//"id": "WARNING: License {{.licenseID}} is not a valid one.\n",
//"translation": "WARNING: License {{.licenseID}} is not a valid one.\n"
//},
//{
//"id": "memorySize of limits in manifest should be an integer between 128 and 512.\n",
//"translation": "memorySize of limits in manifest should be an integer between 128 and 512.\n"
//},
//{
//"id": "timeout of limits in manifest should be an integer between 100 and 300000.\n",
//"translation": "timeout of limits in manifest should be an integer between 100 and 300000.\n"
//},
//{
//"id": "logSize of limits in manifest should be an integer between 0 and 10.\n",
//"translation": "logSize of limits in manifest should be an integer between 0 and 10.\n"
//},
//{
//"id": "WARNING: Invalid limitation 'timeout' of action in manifest is ignored. Please check errors.\n",
//"translation": "WARNING: Invalid limitation 'timeout' of action in manifest is ignored. Please check errors.\n"
//},
//{
//"id": "WARNING: Invalid limitation 'memorySize' of action in manifest is ignored. Please check errors.\n",
//"translation": "WARNING: Invalid limitation 'memorySize' of action in manifest is ignored. Please check errors.\n"
//},
//{
//"id": "WARNING: Invalid limitation 'logSize' of action in manifest is ignored. Please check errors.\n",
//"translation": "WARNING: Invalid limitation 'logSize' of action in manifest is ignored. Please check errors.\n"
//},
//{
//"id": "WARNING: Limits  {{.limitname}}  is not changable as to now, which will be ignored.\n",
//"translation": "WARNING: Limits  {{.limitname}}  is not changable as to now, which will be ignored.\n"
//},

//{
//"id": "The name of the application {{.appNameDeploy}} in deployment file at [{{.deploymentFile}}] does not match the name of the application {{.appNameManifest}}} in manifest file at [{{.manifestFile}}].",
//"translation": "The name of the application {{.appNameDeploy}} in deployment file at [{{.deploymentFile}}] does not match the name of the application {{.appNameManifest}}} in manifest file at [{{.manifestFile}}]."
//},
//{
//"id": "WARNING: application in deployment file will soon be deprecated, please use project instead.\n",
//"translation": "WARNING: application in deployment file will soon be deprecated, please use project instead.\n"
//},
//{
//"id": "WARNING: application in manifest file will soon be deprecated, please use project instead.\n",
//"translation": "WARNING: application in manifest file will soon be deprecated, please use project instead.\n"
//},
//{
//"id": "Undeployment of deleted entities did not complete sucessfully during managed deployment. Run `wskdeploy undeploy` to remove partially deployed assets.\n",
//"translation": "Undeployment of deleted entities did not complete sucessfully during managed deployment. Run `wskdeploy undeploy` to remove partially deployed assets.\n"
//},
//{
//"id": "Found the action {{.action}} which is deleted from the current project {{.project}} in manifest file which is being undeployed.\n",
//"translation": "Found the action {{.action}} which is deleted from the current project {{.project}} in manifest file which is being undeployed.\n"
//},
//{
//"id": "Found the trigger {{.trigger}} which is deleted from the current project {{.project}} in manifest file which is being undeployed.\n",
//"translation": "Found the trigger {{.trigger}} which is deleted from the current project {{.project}} in manifest file which is being undeployed.\n"
//},
//{
//"id": "Found the package {{.package}} which is deleted from the current project {{.project}} in manifest file which is being undeployed.\n",
//"translation": "Found the package {{.package}} which is deleted from the current project {{.project}} in manifest file which is being undeployed.\n"
//}