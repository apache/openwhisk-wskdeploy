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

	// wskdeploy Command messages
	ID_CMD_FLAG_MANIFEST	= "msg_cmd_flag_manifest"
	ID_CMD_FLAG_DEPLOYMENT	= "msg_cmd_flag_deployment"
	ID_CMD_FLAG_STRICT	= "msg_cmd_flag_strict"
	ID_CMD_FLAG_INTERACTIVE	= "msg_cmd_flag_interactive"
	ID_CMD_FLAG_DEFAULTS	= "msg_cmd_flag_allow_defaults"
	ID_CMD_FLAG_VERBOSE	= "msg_cmd_flag_allow_verbose"
	ID_CMD_FLAG_AUTH_KEY	= "msg_cmd_flag_auth_key"
	ID_CMD_FLAG_NAMESPACE	= "msg_cmd_flag_namespace"
	ID_CMD_FLAG_API_HOST	= "msg_cmd_flag_api_host"
	ID_CMD_FLAG_API_VERSION	= "msg_cmd_flag_api_version"
	ID_CMD_FLAG_KEY_FILE	= "msg_cmd_flag_key_file"
	ID_CMD_FLAG_CERT_FILE	= "msg_cmd_flag_cert_file"
	ID_CMD_FLAG_MANAGED	= "msg_cmd_flag_allow_managed"
	ID_CMD_FLAG_PROJECT	= "msg_cmd_flag_project"
	ID_CMD_FLAG_TOGGLE_HELP	= "msg_cmd_flag_toggle_help"
	ID_CMD_FLAG_CONFIG	= "msg_cmd_flag_config"

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

	ID_MSG_MANIFEST_DEPLOY_X_path_X				= "msg_using_manifest_deploy"
	ID_MSG_MANIFEST_UNDEPLOY_X_path_X			= "msg_using_manifest_undeploy"

	ID_MSG_DEPLOYMENT_SUCCEEDED				= "msg_deployment_succeeded"
	ID_MSG_DEPLOYMENT_FAILED				= "msg_deployment_failed"
	ID_MSG_DEPLOYMENT_CANCELLED				= "msg_deployment_cancelled"
	ID_MSG_DEPLOYMENT_REPORT				= "msg_deployment_report_status"

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
	ID_WARN_LIMITS_TIMEOUT					= "msg_warn_limits_timeout"
	ID_WARN_LIMITS_MEMORY_SIZE				= "msg_warn_limits_memory_size"
	ID_WARN_LIMITS_LOG_SIZE					= "msg_warn_limits_memory_log_size"
	ID_WARN_LIMIT_UNCHANGEABLE_X_name_X			= "msg_warn_limit_changeable"
	ID_WARN_RETRY_COMMAND					= "msg_warn_retry_command"
	ID_WARN_CONFIG_INVALID_X_path_X				= "msg_warn_config_invalid"

	// Errors
	ID_ERR_GET_RUNTIMES_X_err_X 				= "msg_err_get_runtimes"
	ID_ERR_MISSING_MANDATORY_KEY_X_key_X			= "msg_err_missing_mandatory_key"
        ID_ERR_MISMATCH_NAME_X_key_X_dname_X_dpath_X_mname_X_moath_X = "msg_err_mismatch_name_project"
	ID_ERR_CREATE_ENTITY_X_key_X_err_X_code_X		= "msg_err_create_entity"
	ID_ERR_DELETE_ENTITY_X_key_X_err_X_code_X		= "msg_err_delete_entity"
	ID_ERR_FEED_INVOKE_X_err_X_code_X			= "msg_err_feed_invoke"
	ID_ERR_INVALID_URL_X_urltype_X_url_X_filetype_X		= "msg_err_url_invalid"
	ID_ERR_MALFORMED_URL_X_urltype_X_url_X			= "msg_err_url_malformed"

	// wskdeploy (as an Action) JSON messages
	ID_ERR_JSON_MISSING_KEY_CMD = "msg_err_json_missing_cmd_key"	// "Missing 'cmd' input key"

	// Misc
	ID_CMD_PUBLISH_DESC_SHORT				= "msg_cmd_publish_short"
	ID_CMD_PUBLISH_DESC_LONG				= "msg_cmd_publish_long"
)

// Known keys used for text replacement in i18n translated strings
const(
	KEY_KEY			= "key"
	KEY_VALUE		= "value"
	KEY_NAME		= "name"
	KEY_CODE		= "code"
	KEY_CMD			= "cmd"
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
	KEY_URL			= "url"
	KEY_URL_TYPE		= "urltype"
)

var I18N_ID_SET = [](string){
	ID_MSG_PREFIX_ERROR,
	ID_MSG_PREFIX_SUCCESS,
	ID_MSG_PREFIX_WARNING,
	ID_MSG_PREFIX_INFO,
	ID_CMD_FLAG_MANIFEST,
	ID_CMD_FLAG_DEPLOYMENT,
	ID_CMD_FLAG_STRICT,
	ID_CMD_FLAG_INTERACTIVE,
	ID_CMD_FLAG_DEFAULTS,
	ID_CMD_FLAG_VERBOSE,
	ID_CMD_FLAG_AUTH_KEY,
	ID_CMD_FLAG_NAMESPACE,
	ID_CMD_FLAG_API_HOST,
	ID_CMD_FLAG_API_VERSION,
	ID_CMD_FLAG_KEY_FILE,
	ID_CMD_FLAG_CERT_FILE,
	ID_CMD_FLAG_MANAGED,
	ID_CMD_FLAG_PROJECT,
	ID_CMD_FLAG_TOGGLE_HELP,
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
	ID_MSG_DEPLOYMENT_REPORT,
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
	ID_WARN_LIMITS_TIMEOUT,
	ID_WARN_LIMITS_MEMORY_SIZE,
	ID_WARN_LIMITS_LOG_SIZE,
	ID_WARN_LIMIT_UNCHANGEABLE_X_name_X,
	ID_ERR_GET_RUNTIMES_X_err_X,
	ID_ERR_MISSING_MANDATORY_KEY_X_key_X,
	ID_ERR_MISMATCH_NAME_X_key_X_dname_X_dpath_X_mname_X_moath_X,
	ID_ERR_CREATE_ENTITY_X_key_X_err_X_code_X,
	ID_ERR_DELETE_ENTITY_X_key_X_err_X_code_X,
	ID_ERR_INVALID_URL_X_urltype_X_url_X_filetype_X,
	ID_ERR_MALFORMED_URL_X_urltype_X_url_X,
	ID_ERR_JSON_MISSING_KEY_CMD,
	ID_WARN_CONFIG_INVALID_X_path_X,
	ID_CMD_FLAG_CONFIG,
}
