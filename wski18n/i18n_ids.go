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

// descriptive key names
// DO NOT TRANSLATE
const (
	ACTION_CODE        = "Action source"
	ACTIONS            = "Actions"
	ACTIVATIONS        = "Activations"
	API_HOST           = "API host"
	AUTH_KEY           = "authentication key"
	COMMAND_LINE       = "wskdeploy command line"
	DEPLOYMENT         = "deployment"
	CONFIGURATION      = "Configuration"
	MANIFEST           = "manifest"
	NAME_ACTION        = "Action Name"
	NAME_FEED          = "Feed Name"
	NAME_RULE          = "Rule Name"
	NAME_RUNTIME       = "Runtime Name"
	NAME_TRIGGER       = "Trigger Name"
	NAMESPACES         = "Namespaces"
	PACKAGE_BINDING    = "package binding"
	PACKAGE_LICENSE    = "package license"
	PACKAGE_NAME       = "package name"
	PACKAGE_VERSION    = "package version"
	PACKAGES           = "Packages"
	PROJECT_NAME       = "project name"
	REGISTRY           = "registry"
	REGISTRY_URL       = "registry URL"
	REPOSITORY         = "repository"
	RULES              = "Rules"
	TRIGGER_FEED       = "trigger feed"
	TRIGGERS           = "Triggers"
	WHISK_PROPS        = "wskprops"
	APIGW_ACCESS_TOKEN = "API Gateway Access Token"
)

// i18n Identifiers
// DO NOT TRANSLATE
const (
	// Debug / trace message prefixes
	ID_MSG_PREFIX_ERROR   = "msg_prefix_error"   // "Error"
	ID_MSG_PREFIX_INFO    = "msg_prefix_info"    // "Info"
	ID_MSG_PREFIX_SUCCESS = "msg_prefix_success" // "Success"
	ID_MSG_PREFIX_WARNING = "msg_prefix_warning" // "Warning"

	// Cobra command descriptions
	ID_CMD_DESC_LONG_PUBLISH      = "msg_cmd_desc_long_publish"
	ID_CMD_DESC_LONG_REPORT       = "msg_cmd_desc_long_report"
	ID_CMD_DESC_LONG_ROOT         = "msg_cmd_desc_long_root"
	ID_CMD_DESC_SHORT_ADD         = "msg_cmd_desc_short_add"
	ID_CMD_DESC_SHORT_ADD_X_key_X = "msg_cmd_desc_short_add_entity"
	ID_CMD_DESC_SHORT_PUBLISH     = "msg_cmd_desc_short_publish"
	ID_CMD_DESC_SHORT_REPORT      = "msg_cmd_desc_short_report"
	ID_CMD_DESC_SHORT_ROOT        = "msg_cmd_desc_short_root"
	ID_CMD_DESC_SHORT_VERSION     = "msg_cmd_desc_short_version"

	// Cobra Flag messages
	ID_CMD_FLAG_API_HOST    = "msg_cmd_flag_api_host"
	ID_CMD_FLAG_API_VERSION = "msg_cmd_flag_api_version"
	ID_CMD_FLAG_AUTH_KEY    = "msg_cmd_flag_auth_key"
	ID_CMD_FLAG_CERT_FILE   = "msg_cmd_flag_cert_file"
	ID_CMD_FLAG_CONFIG      = "msg_cmd_flag_config"
	ID_CMD_FLAG_DEFAULTS    = "msg_cmd_flag_allow_defaults"
	ID_CMD_FLAG_DEPLOYMENT  = "msg_cmd_flag_deployment"
	ID_CMD_FLAG_INTERACTIVE = "msg_cmd_flag_interactive"
	ID_CMD_FLAG_KEY_FILE    = "msg_cmd_flag_key_file"
	ID_CMD_FLAG_MANAGED     = "msg_cmd_flag_allow_managed"
	ID_CMD_FLAG_MANIFEST    = "msg_cmd_flag_manifest"
	ID_CMD_FLAG_NAMESPACE   = "msg_cmd_flag_namespace"
	ID_CMD_FLAG_PROJECT     = "msg_cmd_flag_project"
	ID_CMD_FLAG_STRICT      = "msg_cmd_flag_strict"
	ID_CMD_FLAG_TOGGLE_HELP = "msg_cmd_flag_toggle_help"
	ID_CMD_FLAG_VERBOSE     = "msg_cmd_flag_allow_verbose"

	// Configuration messages
	ID_MSG_CONFIG_MISSING_AUTHKEY                       = "msg_config_missing_authkey"
	ID_MSG_CONFIG_MISSING_APIHOST                       = "msg_config_missing_apihost"
	ID_MSG_CONFIG_MISSING_NAMESPACE                     = "msg_config_missing_namespace"
	ID_MSG_CONFIG_MISSING_APIGW_ACCESS_TOKEN            = "msg_config_missing_apigw_access_token"
	ID_MSG_CONFIG_INFO_APIHOST_X_host_X_source_X        = "msg_config_apihost_info"
	ID_MSG_CONFIG_INFO_AUTHKEY_X_source_X               = "msg_config_authkey_info"
	ID_MSG_CONFIG_INFO_NAMESPACE_X_namespace_X_source_X = "msg_config_namespace_info"
	ID_MSG_CONFIG_INFO_APIGE_ACCESS_TOKEN_X_source_X    = "msg_config_apigw_access_token_info"

	// YAML marshal / unmarshal
	ID_MSG_UNMARSHAL_LOCAL           = "msg_unmarshal_local"
	ID_MSG_UNMARSHAL_NETWORK_X_url_X = "msg_unmarshal_network"

	// Informational
	ID_MSG_DEPLOYMENT_CANCELLED = "msg_deployment_cancelled"
	ID_MSG_DEPLOYMENT_FAILED    = "msg_deployment_failed"
	ID_MSG_DEPLOYMENT_REPORT    = "msg_deployment_report_status"
	ID_MSG_DEPLOYMENT_SUCCEEDED = "msg_deployment_succeeded"

	ID_MSG_UNDEPLOYMENT_CANCELLED = "msg_undeployment_cancelled"
	ID_MSG_UNDEPLOYMENT_FAILED    = "msg_undeployment_failed"
	ID_MSG_UNDEPLOYMENT_SUCCEEDED = "msg_undeployment_succeeded"

	ID_MSG_MANIFEST_DEPLOY_X_path_X   = "msg_manifest_using_deployment"
	ID_MSG_MANIFEST_UNDEPLOY_X_path_X = "msg_manifest_using_undeployment"

	ID_MSG_ENTITY_DEPLOYED_SUCCESS_X_key_X_name_X   = "msg_entity_deployed_success"
	ID_MSG_ENTITY_DEPLOYING_X_key_X_name_X          = "msg_entity_deploying"
	ID_MSG_ENTITY_UNDEPLOYED_SUCCESS_X_key_X_name_X = "msg_entity_undeployed_success"
	ID_MSG_ENTITY_UNDEPLOYING_X_key_X_name_X        = "msg_entity_undeploying"

	ID_MSG_DEPENDENCY_DEPLOYING_X_name_X            = "msg_dependency_deploying"
	ID_MSG_DEPENDENCY_DEPLOYMENT_FAILURE_X_name_X   = "msg_dependency_deployment_failure"
	ID_MSG_DEPENDENCY_DEPLOYMENT_SUCCESS_X_name_X   = "msg_dependency_deployment_success"
	ID_MSG_DEPENDENCY_UNDEPLOYING_X_name_X          = "msg_dependency_undeploying"
	ID_MSG_DEPENDENCY_UNDEPLOYMENT_FAILURE_X_name_X = "msg_dependency_undeployment_failure"
	ID_MSG_DEPENDENCY_UNDEPLOYMENT_SUCCESS_X_name_X = "msg_dependency_undeployment_success"

	ID_MSG_DEFAULT_PACKAGE = "msg_default_package"

	// Managed deployments
	ID_MSG_MANAGED_UNDEPLOYMENT_FAILED                    = "msg_managed_undeployment_failed"
	ID_MSG_MANAGED_FOUND_DELETED_X_key_X_name_X_project_X = "msg_managed_found_deleted_entity"

	// Interactive (prompts)
	ID_MSG_PROMPT_APIHOST   = "msg_prompt_apihost"
	ID_MSG_PROMPT_AUTHKEY   = "msg_prompt_authkey"
	ID_MSG_PROMPT_DEPLOY    = "msg_prompt_deploy"
	ID_MSG_PROMPT_NAMESPACE = "msg_prompt_namespace"
	ID_MSG_PROMPT_UNDEPLOY  = "msg_prompt_undeploy"

	// Errors
	ID_ERR_DEPENDENCY_UNKNOWN_TYPE                               = "msg_err_dependency_unknown_type"
	ID_ERR_DEPLOYMENT_NAME_NOT_FOUND_X_key_X_name_X              = "msg_err_deployment_name_not_found"
	ID_ERR_ENTITY_CREATE_X_key_X_err_X_code_X                    = "msg_err_entity_create"
	ID_ERR_ENTITY_DELETE_X_key_X_err_X_code_X                    = "msg_err_entity_delete"
	ID_ERR_FEED_INVOKE_X_err_X_code_X                            = "msg_err_feed_invoke"
	ID_ERR_KEY_MISSING_X_key_X                                   = "msg_err_key_missing_mandatory"
	ID_ERR_MANIFEST_FILE_NOT_FOUND_X_path_X                      = "msg_err_manifest_not_found"
	ID_ERR_NAME_MISMATCH_X_key_X_dname_X_dpath_X_mname_X_moath_X = "msg_err_name_mismatch"
	ID_ERR_RUNTIME_INVALID_X_runtime_X_action_X                  = "msg_err_runtime_invalid"
	ID_ERR_RUNTIME_MISMATCH_X_runtime_X_ext_X_action_X           = "msg_err_runtime_mismatch"
	ID_ERR_RUNTIMES_GET_X_err_X                                  = "msg_err_runtimes_get"
	ID_ERR_URL_INVALID_X_urltype_X_url_X_filetype_X              = "msg_err_url_invalid"
	ID_ERR_URL_MALFORMED_X_urltype_X_url_X                       = "msg_err_url_malformed"
	ID_ERR_API_MISSING_WEB_ACTION_X_action_X_api_X               = "msg_err_api_missing_web_action"
	ID_ERR_API_MISSING_ACTION_X_action_X_api_X                   = "msg_err_api_missing_action"

	// Server-side Errors (wskdeploy as an Action)
	ID_ERR_JSON_MISSING_KEY_CMD = "msg_err_json_missing_cmd_key"

	// warnings
	ID_WARN_COMMAND_RETRY                                     = "msg_warn_command_retry"
	ID_WARN_CONFIG_INVALID_X_path_X                           = "msg_warn_config_invalid"
	ID_WARN_KEY_DEPRECATED_X_oldkey_X_filetype_X_newkey_X     = "msg_warn_key_deprecated_replaced"
	ID_WARN_KEY_MISSING_X_key_X_value_X                       = "msg_warn_key_missing"
	ID_WARN_KEYVALUE_INVALID                                  = "msg_warn_key_value_invalid"
	ID_WARN_KEYVALUE_NOT_SAVED_X_key_X                        = "msg_warn_key_value_not_saved"
	ID_WARN_LIMIT_IGNORED_X_limit_X                           = "msg_warn_limit_ignored"
	ID_WARN_LIMIT_UNCHANGEABLE_X_name_X                       = "msg_warn_limit_changeable"
	ID_WARN_LIMITS_LOG_SIZE                                   = "msg_warn_limits_log_size"    // TODO() remove for value range
	ID_WARN_LIMITS_MEMORY_SIZE                                = "msg_warn_limits_memory_size" // TODO() remove for value range
	ID_WARN_LIMITS_TIMEOUT                                    = "msg_warn_limits_timeout"     // TODO() remove for value range
	ID_WARN_RUNTIME_CHANGED_X_runtime_X_action_X              = "msg_warn_runtime_changed"
	ID_WARN_VALUE_RANGE_X_name_X_key_X_filetype_X_min_X_max_X = "msg_warn_value_range" // TODO() not used, but should be used for limit ranges
	ID_WARN_WHISK_PROPS_DEPRECATED                            = "msg_warn_whisk_properties"
	ID_WARN_ENTITY_NAME_EXISTS_X_key_X_name_X                 = "msg_warn_entity_name_exists"
	ID_WARN_PACKAGES_NOT_FOUND_X_path_X                       = "msg_warn_packages_not_found"

	// Verbose (Debug/Trace) messages
	ID_DEBUG_KEY_VERIFY_X_name_X_key_X     = "msg_dbg_key_verify"
	ID_DEBUG_DEPLOYING_USING               = "msg_dbg_deploying_using"
	ID_DEBUG_UNDEPLOYING_USING             = "msg_dbg_undeploying_using"
	ID_DEBUG_PROJECT_SEARCH_X_path_X_key_X = "msg_dbg_searching_project_directory"
)

// Known keys used for text replacement in i18n translated strings
const (
	KEY_ACTION             = "action"
	KEY_CMD                = "cmd"
	KEY_CODE               = "code"
	KEY_DEPLOYMENT_NAME    = "dname"
	KEY_DEPLOYMENT_PATH    = "dpath"
	KEY_ERR                = "err"
	KEY_EXTENTION          = "ext"
	KEY_FILE_TYPE          = "filetype"
	KEY_HOST               = "host"
	KEY_KEY                = "key"
	KEY_LIMIT              = "limit"
	KEY_MANIFEST_NAME      = "mname"
	KEY_MANIFEST_PATH      = "mpath"
	KEY_NAME               = "name"
	KEY_NAMESPACE          = "namespace"
	KEY_NEW                = "newkey"
	KEY_OLD                = "oldkey"
	KEY_PATH               = "path"
	KEY_PROJECT            = "project"
	KEY_RUNTIME            = "runtime"
	KEY_SOURCE             = "source"
	KEY_URL                = "url"
	KEY_URL_TYPE           = "urltype"
	KEY_VALUE              = "value"
	KEY_VALUE_MIN          = "min" // TODO() attempt to use this for Limit value range errors
	KEY_VALUE_MAX          = "max" // TODO() attempt to use this for Limit value range errors
	KEY_API                = "api"
	KEY_APIGW_ACCESS_TOKEN = "apigw_access_token"
)

// Used to unit test that translations exist with these IDs
var I18N_ID_SET = [](string){
	ID_CMD_DESC_LONG_PUBLISH,
	ID_CMD_DESC_LONG_REPORT,
	ID_CMD_DESC_LONG_ROOT,
	ID_CMD_DESC_SHORT_ADD,
	ID_CMD_DESC_SHORT_ADD_X_key_X,
	ID_CMD_DESC_SHORT_PUBLISH,
	ID_CMD_DESC_SHORT_REPORT,
	ID_CMD_DESC_SHORT_ROOT,
	ID_CMD_DESC_SHORT_VERSION,
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
	ID_DEBUG_DEPLOYING_USING,
	ID_DEBUG_KEY_VERIFY_X_name_X_key_X,
	ID_DEBUG_PROJECT_SEARCH_X_path_X_key_X,
	ID_DEBUG_UNDEPLOYING_USING,
	ID_ERR_DEPENDENCY_UNKNOWN_TYPE,
	ID_ERR_DEPLOYMENT_NAME_NOT_FOUND_X_key_X_name_X,
	ID_ERR_ENTITY_CREATE_X_key_X_err_X_code_X,
	ID_ERR_ENTITY_DELETE_X_key_X_err_X_code_X,
	ID_ERR_JSON_MISSING_KEY_CMD,
	ID_ERR_JSON_MISSING_KEY_CMD,
	ID_ERR_KEY_MISSING_X_key_X,
	ID_ERR_MANIFEST_FILE_NOT_FOUND_X_path_X,
	ID_ERR_NAME_MISMATCH_X_key_X_dname_X_dpath_X_mname_X_moath_X,
	ID_ERR_RUNTIME_INVALID_X_runtime_X_action_X,
	ID_ERR_RUNTIME_MISMATCH_X_runtime_X_ext_X_action_X,
	ID_ERR_RUNTIMES_GET_X_err_X,
	ID_ERR_URL_INVALID_X_urltype_X_url_X_filetype_X,
	ID_ERR_URL_MALFORMED_X_urltype_X_url_X,
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
	ID_MSG_UNMARSHAL_NETWORK_X_url_X,
	ID_WARN_COMMAND_RETRY,
	ID_WARN_CONFIG_INVALID_X_path_X,
	ID_WARN_ENTITY_NAME_EXISTS_X_key_X_name_X,
	ID_WARN_KEY_DEPRECATED_X_oldkey_X_filetype_X_newkey_X,
	ID_WARN_KEY_MISSING_X_key_X_value_X,
	ID_WARN_KEYVALUE_INVALID,
	ID_WARN_KEYVALUE_NOT_SAVED_X_key_X,
	ID_WARN_LIMIT_IGNORED_X_limit_X,
	ID_WARN_LIMIT_UNCHANGEABLE_X_name_X,
	ID_WARN_LIMITS_LOG_SIZE,
	ID_WARN_LIMITS_MEMORY_SIZE,
	ID_WARN_LIMITS_TIMEOUT,
	ID_WARN_PACKAGES_NOT_FOUND_X_path_X,
	ID_WARN_RUNTIME_CHANGED_X_runtime_X_action_X,
	ID_WARN_WHISK_PROPS_DEPRECATED,
}
