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
	ID_MSG_PREFIX_ERROR	= "msg_prefix_error"	// "Error"
	ID_MSG_PREFIX_SUCCESS	= "msg_prefix_success"	// "Success"
	ID_MSG_PREFIX_WARNING	= "msg_prefix_warning"	// "Warning"
	ID_MSG_PREFIX_INFO	= "msg_prefix_info"	// "Info"

	// wskdeploy (as an Action) JSON messages
	ID_JSON_MISSING_KEY_CMD	= "msg_json_missing_cmd_key"	// "Missing 'cmd' input key"

	// wskdeploy Command messages
	ID_CMD_FLAG_AUTH_KEY	= "msg_cmd_flag_auth_key"	// "authorization `KEY`"
	ID_CMD_FLAG_NAMESPACE	= "msg_cmd_flag_namespace"	// "namespace"
	ID_CMD_FLAG_API_HOST	= "msg_cmd_flag_api_host"	// "whisk API `HOST`"
	ID_CMD_FLAG_API_VERSION	= "msg_cmd_flag_api_version"	// "whisk API `VERSION`"
	ID_CMD_FLAG_KEY_FILE	= "msg_cmd_flag_key_file"	// "path of the .key file"
	ID_CMD_FLAG_CERT_FILE	= "msg_cmd_flag_cert_file"	// "path of the .cert file"

	// Configuration messages
	ID_MSG_CONFIG_MISSING_AUTHKEY				= "msg_config_missing_authkey"
	ID_MSG_CONFIG_MISSING_APIHOST				= "msg_config_missing_apihost"
	ID_MSG_CONFIG_MISSING_NAMESPACE				= "msg_config_missing_namespace"

	ID_MSG_CONFIG_INFO_APIHOST_X_host_X_source_X		= "msg_config_apihost_info"
	ID_MSG_CONFIG_INFO_AUTHKEY_X_source_X			= "msg_config_authkey_info"
	ID_MSG_CONFIG_INFO_NAMESPACE_X_namespace_X_source_X	= "msg_config_namespace_info"

	// YAML marshall / unmarshall
	ID_MSG_UNMARSHAL_LOCAL					= "msg_unmarshall_local"
	ID_MSG_UNMARSHAL_NETWORK				= "msg_unmarshall_network"

	// Informational
	ID_MSG_MANIFEST_FILE_NOT_FOUND_X_path_X			= "msg_manifest_not_found"
	ID_MSG_RUNTIME_MISMATCH_X_runtime_X_ext_X_action_X	= "msg_runtime_mismatch"
	ID_MSG_RUNTIME_CHANGED_X_runtime_X_action_X		= "msg_runtime_changed"
	ID_MSG_RUNTIME_UNSUPPORTED_X_runtime_X_action_X		= "msg_runtime_unsupported"

	ID_MSG_MANIFEST_DEPLOY_X_path_X				= "msg_using_manifest_deploy"	// "Using {{.path}} for deployment.\n"
	ID_MSG_MANIFEST_UNDEPLOY_X_path_X			= "msg_using_manifest_undeploy"	// "Using {{.path}} for undeployment.\n"

	ID_MSG_DEPLOYMENT_SUCCEEDED				= "msg_deployment_succeeded"
	ID_MSG_DEPLOYMENT_FAILED				= "msg_deployment_failed"
	ID_MSG_DEPLOYMENT_CANCELLED				= "msg_deployment_cancelled"

	ID_MSG_ENTITY_DEPLOYING_X_key_X_name_X 			= "msg_entity_deploying"
	ID_MSG_ENTITY_UNDEPLOYING_X_key_X_name_X		= "msg_entity_undeploying"
	ID_MSG_ENTITY_DEPLOYED_SUCCESS_X_key_X_name_X		= "msg_entity_deployed_success"
	ID_MSG_ENTITY_UNDEPLOYED_SUCCESS_X_key_X_name_X		= "msg_entity_undeployed_success"

	ID_MSG_UNDEPLOYMENT_SUCCEEDED				= "msg_undeployment_succeeded"
	ID_MSG_UNDEPLOYMENT_FAILED				= "msg_undeployment_failed"
	ID_MSG_UNDEPLOYMENT_CANCELLED				= "msg_undeployment_cancelled"

	ID_MSG_DEPENDENCY_DEPLOYING_X_name_X			= "msg_deploying_dependency"
	ID_MSG_DEPENDENCY_UNDEPLOYING_X_name_X			= "msg_undeploying_dependency"
	ID_MSG_DEPENDENCY_DEPLOYMENT_SUCCESS_X_name_X		= "msg_dependency_deployment_success"
	ID_MSG_DEPENDENCY_DEPLOYMENT_FAILURE_X_name_X		= "msg_dependency_deployment_failure"
	ID_MSG_DEPENDENCY_UNDEPLOYMENT_SUCCESS_X_name_X		= "msg_dependency_undeployment_success"
	ID_MSG_DEPENDENCY_UNDEPLOYMENT_FAILURE_X_name_X		= "msg_dependency_undeployment_failure"

	// Managed deployments
	ID_MSG_MANAGED_UNDEPLOYMENT_FAILED 			= "msg_undeployment_managed_failed"
	ID_MSG_MANAGED_FOUND_DELETED_X_key_X_name_X_project_X	= "msg_managed_found_deleted_entity"

	// Interactive (prompts)
	ID_MSG_PROMPT_DEPLOY					= "msg_prompt_deploy"
	ID_MSG_PROMPT_UNDEPLOY					= "msg_prompt_undeploy"
	ID_MSG_PROMPT_AUTHKEY					= "msg_prompt_authkey"
	ID_MSG_PROMPT_APIHOST					= "msg_prompt_apihost"
	ID_MSG_PROMPT_NAMESPACE					= "msg_prompt_namespace"

	// Action Limits
	ID_MSG_ACTION_LIMIT_IGNORED_X_limit_X			= "msg_action_limit_ignored"	// timeout, memorySize, logSize

	// warnings
	ID_WARN_DEPRECATED_KEY_REPLACED_X_oldkey_X_filetype_X_newkey_X = "msg_warn_key_deprecated_replaced"
	ID_WARN_WHISK_PROPS_DEPRECATED				= "msg_warn_whisk_properties"
	ID_WARN_MISSING_MANDATORY_KEY_X_key_X_value_X		= "msg_warn_missing_mandatory_key"
	ID_WARN_KEYVALUE_NOT_SAVED_X_key_X			= "msg_warn_key_value_not_saved"
	ID_WARN_KEYVALUE_INVALID				= "msg_warn_invalid_key_value"

	// Errors
	ID_ERR_GET_RUNTIMES_X_err_X 				= "msg_err_get_runtimes"
	ID_ERR_MISSING_MANDATORY_KEY_X_key_X			= "msg_err_missing_mandatory_key"
        ID_ERR_MISMATCH_NAME_X_key_X_dname_X_dpath_X_mname_X_moath_X = "msg_err_mismatch_name_project"
	ID_ERR_CREATE_ENTITY_X_key_X_err_X_code_X		= "msg_err_create_entity"
	ID_ERR_DELETE_ENTITY_X_key_X_err_X_code_X		= "msg_err_delete_entity"
	ID_ERR_FEED_INVOKE_X_err_X_code_X			= "msg_err_feed_invoke"

)

// Known keys used for text replacement in i18n translated strings
const(
	KEY_KEY			= "key"
	KEY_VALUE		= "value"
	KEY_NAME		= "name"
	KEY_CODE		= "code"
	KEY_ERR			= "err"
	KEY_PROJECT		= "project"
	KEY_ACTION		= "action"
	KEY_LIMIT		= "limit"
	KEY_HOST		= "host"
	KEY_SOURCE		= "source"
	KEY_NAMESPACE		= "namespace"
	KEY_PATH		= "path"
	KEY_EXTENTION		= "ext"
	KEY_RUNTIME		= "runtime"
	KEY_DEPLOYMENT_NAME	= "dname"
	KEY_DEPLOYMENT_PATH	= "dpath"
	KEY_MANIFEST_NAME	= "mname"
	KEY_MANIFEST_PATH	= "mpath"
	KEY_OLD			= "oldkey"
	KEY_NEW			= "newkey"
	KEY_FILE_TYPE		= "filetype"
)

var I18N_ID_SET = [](string){
	ID_MSG_PREFIX_ERROR,
	ID_MSG_PREFIX_SUCCESS,
	ID_MSG_PREFIX_WARNING,
	ID_MSG_PREFIX_INFO,
	ID_JSON_MISSING_KEY_CMD,
	ID_CMD_FLAG_AUTH_KEY,
	ID_CMD_FLAG_NAMESPACE,
	ID_CMD_FLAG_API_HOST,
	ID_CMD_FLAG_API_VERSION,
	ID_CMD_FLAG_KEY_FILE,
	ID_CMD_FLAG_CERT_FILE,
	ID_MSG_CONFIG_MISSING_AUTHKEY,
	ID_MSG_CONFIG_MISSING_APIHOST,
	ID_MSG_CONFIG_MISSING_NAMESPACE,
	ID_MSG_CONFIG_INFO_APIHOST_X_host_X_source_X,
	ID_MSG_CONFIG_INFO_AUTHKEY_X_source_X,
	ID_MSG_CONFIG_INFO_NAMESPACE_X_namespace_X_source_X,
	ID_MSG_UNMARSHAL_LOCAL,
	ID_MSG_UNMARSHAL_NETWORK,
	ID_MSG_MANIFEST_FILE_NOT_FOUND_X_path_X,
	ID_MSG_RUNTIME_MISMATCH_X_runtime_X_ext_X_action_X,
	ID_MSG_RUNTIME_CHANGED_X_runtime_X_action_X,
	ID_MSG_RUNTIME_UNSUPPORTED_X_runtime_X_action_X,
	ID_MSG_MANIFEST_DEPLOY_X_path_X,
	ID_MSG_MANIFEST_UNDEPLOY_X_path_X,
	ID_MSG_DEPLOYMENT_SUCCEEDED,
	ID_MSG_DEPLOYMENT_FAILED,
	ID_MSG_DEPLOYMENT_CANCELLED,
	ID_MSG_ENTITY_DEPLOYING_X_key_X_name_X,
	ID_MSG_ENTITY_UNDEPLOYING_X_key_X_name_X,
	ID_MSG_ENTITY_DEPLOYED_SUCCESS_X_key_X_name_X,
	ID_MSG_ENTITY_UNDEPLOYED_SUCCESS_X_key_X_name_X,
	ID_MSG_UNDEPLOYMENT_SUCCEEDED,
	ID_MSG_UNDEPLOYMENT_FAILED,
	ID_MSG_UNDEPLOYMENT_CANCELLED,
	ID_MSG_DEPENDENCY_DEPLOYING_X_name_X,
	ID_MSG_DEPENDENCY_UNDEPLOYING_X_name_X,
	ID_MSG_DEPENDENCY_DEPLOYMENT_SUCCESS_X_name_X,
	ID_MSG_DEPENDENCY_DEPLOYMENT_FAILURE_X_name_X,
	ID_MSG_DEPENDENCY_UNDEPLOYMENT_SUCCESS_X_name_X,
	ID_MSG_DEPENDENCY_UNDEPLOYMENT_FAILURE_X_name_X,
	ID_MSG_MANAGED_UNDEPLOYMENT_FAILED,
	ID_MSG_MANAGED_FOUND_DELETED_X_key_X_name_X_project_X,
	ID_MSG_PROMPT_DEPLOY,
	ID_MSG_PROMPT_UNDEPLOY,
	ID_MSG_PROMPT_AUTHKEY,
	ID_MSG_PROMPT_APIHOST,
	ID_MSG_PROMPT_NAMESPACE,
	ID_MSG_ACTION_LIMIT_IGNORED_X_limit_X,
	ID_WARN_DEPRECATED_KEY_REPLACED_X_oldkey_X_filetype_X_newkey_X,
	ID_WARN_WHISK_PROPS_DEPRECATED,
	ID_WARN_MISSING_MANDATORY_KEY_X_key_X_value_X,
	ID_WARN_KEYVALUE_NOT_SAVED_X_key_X,
	ID_WARN_KEYVALUE_INVALID,
	ID_ERR_GET_RUNTIMES_X_err_X,
	ID_ERR_MISSING_MANDATORY_KEY_X_key_X,
	ID_ERR_MISMATCH_NAME_X_key_X_dname_X_dpath_X_mname_X_moath_X,
	ID_ERR_CREATE_ENTITY_X_key_X_err_X_code_X,
	ID_ERR_DELETE_ENTITY_X_key_X_err_X_code_X,
}
