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
	ID_MSG_PREFIX_INFO	= "msg_prefix_info"	// "Info"
	ID_MSG_PREFIX_SUCCESS	= "msg_prefix_success"	// "Success"
	ID_MSG_PREFIX_WARNING	= "msg_prefix_warning"	// "Warning"

	// wskdeploy Command messages
	ID_CMD_FLAG_API_HOST	= "msg_cmd_flag_api_host"
	ID_CMD_FLAG_API_VERSION	= "msg_cmd_flag_api_version"
	ID_CMD_FLAG_AUTH_KEY	= "msg_cmd_flag_auth_key"
	ID_CMD_FLAG_CERT_FILE	= "msg_cmd_flag_cert_file"
	ID_CMD_FLAG_CONFIG	= "msg_cmd_flag_config"
	ID_CMD_FLAG_DEFAULTS	= "msg_cmd_flag_allow_defaults"
	ID_CMD_FLAG_DEPLOYMENT	= "msg_cmd_flag_deployment"
	ID_CMD_FLAG_INTERACTIVE	= "msg_cmd_flag_interactive"
	ID_CMD_FLAG_KEY_FILE	= "msg_cmd_flag_key_file"
	ID_CMD_FLAG_MANAGED	= "msg_cmd_flag_allow_managed"
	ID_CMD_FLAG_MANIFEST	= "msg_cmd_flag_manifest"
	ID_CMD_FLAG_NAMESPACE	= "msg_cmd_flag_namespace"
	ID_CMD_FLAG_PROJECT	= "msg_cmd_flag_project"
	ID_CMD_FLAG_STRICT	= "msg_cmd_flag_strict"
	ID_CMD_FLAG_TOGGLE_HELP	= "msg_cmd_flag_toggle_help"
	ID_CMD_FLAG_VERBOSE	= "msg_cmd_flag_allow_verbose"

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

	// Action Limits (i.e., timeout, memorySize, logSize)
	ID_MSG_ACTION_LIMIT_IGNORED_X_limit_X			= "msg_action_limit_ignored"

	// warnings
	ID_WARN_CONFIG_INVALID_X_path_X				= "msg_warn_config_invalid"
	ID_WARN_KEY_DEPRECATED_X_oldkey_X_filetype_X_newkey_X	= "msg_warn_key_deprecated_replaced"
	ID_WARN_KEYVALUE_INVALID				= "msg_warn_invalid_key_value"
	ID_WARN_KEYVALUE_NOT_SAVED_X_key_X			= "msg_warn_key_value_not_saved"
	ID_WARN_LIMIT_UNCHANGEABLE_X_name_X			= "msg_warn_limit_changeable"
	ID_WARN_LIMITS_LOG_SIZE					= "msg_warn_limits_memory_log_size" 	// TODO() remove for value range
	ID_WARN_LIMITS_MEMORY_SIZE				= "msg_warn_limits_memory_size" 	// TODO() remove for value range
	ID_WARN_LIMITS_TIMEOUT					= "msg_warn_limits_timeout"  		// TODO() remove for value range
	ID_WARN_MISSING_MANDATORY_KEY_X_key_X_value_X		= "msg_warn_missing_mandatory_key"
	ID_WARN_COMMAND_RETRY					= "msg_warn_retry_command"
	ID_WARN_VALUE_RANGE_X_name_X_key_X_filetype_X_min_X_max_X = "msg_warn_value_range"
	ID_WARN_WHISK_PROPS_DEPRECATED				= "msg_warn_whisk_properties"
	ID_WARN_RUNTIME_CHANGED_X_runtime_X_action_X		= "msg_warn_runtime_changed"

	// Errors
	ID_ERR_DEPENDENCY_UNKNOWN_TYPE				= "msg_err_dependency_unknown_type"
	ID_ERR_DEPLOYMENT_NAME_NOT_FOUND_X_key_X_name_X		= "msg_err_deployment_name_not_found"
	ID_ERR_ENTITY_CREATE_X_key_X_err_X_code_X 		= "msg_err_entity_create"
	ID_ERR_ENTITY_DELETE_X_key_X_err_X_code_X 		= "msg_err_entity_delete"
	ID_ERR_FEED_INVOKE_X_err_X_code_X			= "msg_err_feed_invoke"
	ID_ERR_GET_RUNTIMES_X_err_X 				= "msg_err_get_runtimes"
	ID_ERR_INVALID_URL_X_urltype_X_url_X_filetype_X		= "msg_err_url_invalid"
	ID_ERR_MALFORMED_URL_X_urltype_X_url_X			= "msg_err_url_malformed"
	ID_ERR_MISSING_MANDATORY_KEY_X_key_X			= "msg_err_missing_mandatory_key"
	ID_ERR_RUNTIME_INVALID_X_runtime_X_action_X		= "msg_err_runtime_invalid"
	ID_ERR_RUNTIME_MISMATCH_X_runtime_X_ext_X_action_X	= "msg_err_runtime_mismatch"
        ID_ERR_MISMATCH_NAME_X_key_X_dname_X_dpath_X_mname_X_moath_X = "msg_err_mismatch_name_project"

	// Server-side Errors (wskdeploy as an Action)
	ID_ERR_JSON_MISSING_KEY_CMD = "msg_err_json_missing_cmd_key"	// "Missing 'cmd' input key"

	// Cobra command / flag descriptions
	ID_CMD_DESC_SHORT_PUBLISH				= "msg_cmd_desc_short_publish"
	ID_CMD_DESC_LONG_PUBLISH				= "msg_cmd_desc_long_publish"
	ID_CMD_DESC_SHORT_ROOT					= "msg_cmd_desc_short_root"
	ID_CMD_DESC_LONG_ROOT					= "msg_cmd_desc_long_root"
	ID_CMD_DESC_SHORT_REPORT				= "msg_cmd_desc_short_report"
	ID_CMD_DESC_LONG_REPORT					= "msg_cmd_desc_long_report"

	// Verbose (Debug/Trace) messages
	ID_DEBUG_KEY_VERIFY_X_name_X_key_X			= "msg_dbg_key_verify"
)

// Known keys used for text replacement in i18n translated strings
const(
	KEY_ACTION		= "action"
	KEY_CMD			= "cmd"
	KEY_CODE		= "code"
	KEY_DEPLOYMENT_NAME	= "dname"
	KEY_DEPLOYMENT_PATH	= "dpath"
	KEY_ERR			= "err"
	KEY_EXTENTION		= "ext"
	KEY_FILE_TYPE		= "filetype"
	KEY_HOST		= "host"
	KEY_KEY			= "key"
	KEY_LIMIT		= "limit"
	KEY_MANIFEST_NAME	= "mname"
	KEY_MANIFEST_PATH	= "mpath"
	KEY_NAME		= "name"
	KEY_NAMESPACE		= "namespace"
	KEY_NEW			= "newkey"
	KEY_OLD			= "oldkey"
	KEY_PATH		= "path"
	KEY_PROJECT		= "project"
	KEY_RUNTIME		= "runtime"
	KEY_SOURCE		= "source"
	KEY_URL			= "url"
	KEY_URL_TYPE		= "urltype"
	KEY_VALUE		= "value"
	KEY_VALUE_MIN		= "min"		// TODO() attempt to use this for Limit value range errors
	KEY_VALUE_MAX		= "max"		// TODO() attempt to use this for Limit value range errors
)

var I18N_ID_SET = [](string){
	ID_CMD_DESC_LONG_PUBLISH,
	ID_CMD_DESC_LONG_REPORT,
	ID_CMD_DESC_LONG_ROOT,
	ID_CMD_DESC_SHORT_PUBLISH,
	ID_CMD_DESC_SHORT_REPORT,
	ID_CMD_DESC_SHORT_ROOT,
	ID_CMD_FLAG_API_HOST,
	ID_CMD_FLAG_API_VERSION,
	ID_CMD_FLAG_AUTH_KEY,
	ID_CMD_FLAG_CERT_FILE,
	ID_CMD_FLAG_CONFIG,
	ID_CMD_FLAG_DEFAULTS,
	ID_CMD_FLAG_DEPLOYMENT,
	ID_CMD_FLAG_INTERACTIVE,
	ID_CMD_FLAG_KEY_FILE,
	ID_CMD_FLAG_MANAGED,
	ID_CMD_FLAG_MANIFEST,
	ID_CMD_FLAG_NAMESPACE,
	ID_CMD_FLAG_PROJECT,
	ID_CMD_FLAG_STRICT,
	ID_CMD_FLAG_TOGGLE_HELP,
	ID_CMD_FLAG_VERBOSE,
	ID_DEBUG_KEY_VERIFY_X_name_X_key_X,
	ID_ERR_DEPENDENCY_UNKNOWN_TYPE,
	ID_ERR_ENTITY_CREATE_X_key_X_err_X_code_X,
	ID_ERR_ENTITY_DELETE_X_key_X_err_X_code_X,
	ID_ERR_GET_RUNTIMES_X_err_X,
	ID_ERR_INVALID_URL_X_urltype_X_url_X_filetype_X,
	ID_ERR_JSON_MISSING_KEY_CMD,
	ID_ERR_MALFORMED_URL_X_urltype_X_url_X,
	ID_ERR_MISMATCH_NAME_X_key_X_dname_X_dpath_X_mname_X_moath_X,
	ID_ERR_MISSING_MANDATORY_KEY_X_key_X,
	ID_ERR_RUNTIME_INVALID_X_runtime_X_action_X,
	ID_ERR_RUNTIME_MISMATCH_X_runtime_X_ext_X_action_X,
	ID_MSG_ACTION_LIMIT_IGNORED_X_limit_X,
	ID_MSG_CONFIG_INFO_APIHOST_X_host_X_source_X,
	ID_MSG_CONFIG_INFO_AUTHKEY_X_source_X,
	ID_MSG_CONFIG_INFO_NAMESPACE_X_namespace_X_source_X,
	ID_MSG_CONFIG_MISSING_APIHOST,
	ID_MSG_CONFIG_MISSING_AUTHKEY,
	ID_MSG_CONFIG_MISSING_NAMESPACE,
	ID_MSG_DEPENDENCY_DEPLOYING_X_name_X,
	ID_MSG_DEPENDENCY_DEPLOYMENT_FAILURE_X_name_X,
	ID_MSG_DEPENDENCY_DEPLOYMENT_SUCCESS_X_name_X,
	ID_MSG_DEPENDENCY_UNDEPLOYING_X_name_X,
	ID_MSG_DEPENDENCY_UNDEPLOYMENT_FAILURE_X_name_X,
	ID_MSG_DEPENDENCY_UNDEPLOYMENT_SUCCESS_X_name_X,
	ID_MSG_DEPLOYMENT_CANCELLED,
	ID_MSG_DEPLOYMENT_FAILED,
	ID_MSG_DEPLOYMENT_REPORT,
	ID_MSG_DEPLOYMENT_SUCCEEDED,
	ID_MSG_ENTITY_DEPLOYED_SUCCESS_X_key_X_name_X,
	ID_MSG_ENTITY_DEPLOYING_X_key_X_name_X,
	ID_MSG_ENTITY_UNDEPLOYED_SUCCESS_X_key_X_name_X,
	ID_MSG_ENTITY_UNDEPLOYING_X_key_X_name_X,
	ID_MSG_MANAGED_FOUND_DELETED_X_key_X_name_X_project_X,
	ID_MSG_MANAGED_UNDEPLOYMENT_FAILED,
	ID_MSG_MANIFEST_DEPLOY_X_path_X,
	ID_MSG_MANIFEST_FILE_NOT_FOUND_X_path_X,
	ID_MSG_MANIFEST_UNDEPLOY_X_path_X,
	ID_MSG_PREFIX_ERROR,
	ID_MSG_PREFIX_INFO,
	ID_MSG_PREFIX_SUCCESS,
	ID_MSG_PREFIX_WARNING,
	ID_MSG_PROMPT_APIHOST,
	ID_MSG_PROMPT_AUTHKEY,
	ID_MSG_PROMPT_DEPLOY,
	ID_MSG_PROMPT_NAMESPACE,
	ID_MSG_PROMPT_UNDEPLOY,
	ID_MSG_UNDEPLOYMENT_CANCELLED,
	ID_MSG_UNDEPLOYMENT_FAILED,
	ID_MSG_UNDEPLOYMENT_SUCCEEDED,
	ID_MSG_UNMARSHAL_LOCAL,
	ID_MSG_UNMARSHAL_NETWORK,
	ID_WARN_CONFIG_INVALID_X_path_X,
	ID_WARN_KEY_DEPRECATED_X_oldkey_X_filetype_X_newkey_X,
	ID_WARN_KEYVALUE_INVALID,
	ID_WARN_KEYVALUE_NOT_SAVED_X_key_X,
	ID_WARN_LIMIT_UNCHANGEABLE_X_name_X,
	ID_WARN_LIMITS_LOG_SIZE,
	ID_WARN_LIMITS_MEMORY_SIZE,
	ID_WARN_LIMITS_TIMEOUT,
	ID_WARN_MISSING_MANDATORY_KEY_X_key_X_value_X,
	ID_WARN_RUNTIME_CHANGED_X_runtime_X_action_X,
	ID_WARN_WHISK_PROPS_DEPRECATED,
}
