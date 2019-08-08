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

// qualifiedname.go
package utils

import (
	"encoding/json"
	"github.com/apache/openwhisk-client-go/whisk"
	"github.com/apache/openwhisk-wskdeploy/wski18n"
	"github.com/apache/openwhisk-wskdeploy/wskprint"
	"strings"
)

var LocalLicenseRecords = map[string][]string{
	"APACHE": {"1.0", "1.1", "2.0"},
	"LGPL":   {"2.0", "2.1", "3.0"},
	"GPL":    {"1.0", "2.0", "3.0"},
	"MPL":    {"1.0", "1.1", "2.0", "2.0-no-copyleft-exception"},
	"BSD":    {"2-Clause", "2-Clause-FreeBSD", "2-Clause-NetBSD", "3-Clause", "3-Clause-Clear", "3-Clause-No-Nuclear-License", "3-Clause-No-Nuclear-License-2014", "3-Clause-No-Nuclear-Warranty", "3-Clause-Attribution", "4-Clause", "Protection", "Source-Code", "4-Clause-UC", "3-Clause-LBNL"},
	"0BSD":   {""},
	"GFDL":   {"1.0", "1.2", "1.3"},
	"AGPL":   {"3.0"},
	"MIT":    {"", "feh", "enna", "advertising", "CMU"},
	"AFL":    {"1.1", "1.2", "2.0", "2.1", "3.0"},
	"APSL":   {"1.0", "1.1", "1.2", "2.0"},
	"EPL":    {"1.0"},
	"OSL":    {"1.0", "1.1", "2.0", "2.1", "3.0"},
	"PHP":    {"3.0", "3.01"},
}

var RemoteLicenseURL = "https://raw.githubusercontent.com/spdx/license-list-data/master/json/licenses.json"

type LicenseJSON struct {
	Description string        `json:"licenseListVersion"`
	Licenses    []LicenseItem `json:"licenses"`
	ReleaseDate string        `json:"releaseDate"`
}

type LicenseItem struct {
	Name      string `json:"name"`
	LicenseID string `json:"licenseId"`
}

var license_json = LicenseJSON{}

//Check if the license is valid
//Check local data record at first
//Then check remote json data
func CheckLicense(license string) bool {
	// TODO(#673) Strict flag should cause an error to be generatd
	if !LicenseLocalValidation(license) && !LicenseRemoteValidation(license) {
		warningString := wski18n.T(
			wski18n.ID_WARN_KEYVALUE_INVALID,
			map[string]interface{}{
				wski18n.KEY_KEY: license})
		wskprint.PrintlnOpenWhiskWarning(warningString)
		return false
	}
	return true
}

// Check if the license is in local license record
// If it is, return ture
// If it is not, return false
func LicenseLocalValidation(license string) bool {
	license_upper := strings.ToUpper(license) //change it to upper case
	parts := strings.SplitN(license_upper, "-", 2)
	license_type := parts[0]
	license_version := ""
	if len(parts) > 1 {
		license_version = parts[1]
	}

	version_list := LocalLicenseRecords[license_type]
	if version_list == nil {
		return false
	} else {
		for _, version := range LocalLicenseRecords[license_type] {
			if license_version == strings.ToUpper(version) {
				return true
			}
		}
		return false
	}
}

//Check if license is one of validate license in
//https://github.com/spdx/license-list-data/blob/master/json/licenses.json
// If it is, return ture
// If it is not, return false
func LicenseRemoteValidation(license string) bool {
	license_upper := strings.ToUpper(license) //change it to upper case

	if len(license_json.Licenses) == 0 {
		json_data, err := Read(RemoteLicenseURL)
		if err != nil {
			// TODO() i18n
			errString := wski18n.T("Failed to get the remote license json.\n")
			whisk.Debug(whisk.DbgError, errString)
			return false
		}

		//parse json
		err = json.Unmarshal(json_data, &license_json)
		if err != nil || len(license_json.Licenses) == 0 {
			// TODO() i18n
			errString := wski18n.T("Failed to parse the remote license json.\n")
			whisk.Debug(whisk.DbgError, errString)
			return false
		}
	}

	//check license
	for _, licenseobj := range license_json.Licenses {
		if strings.ToUpper(licenseobj.LicenseID) == license_upper {
			return true
		}
	}
	return false
}

//if valid or nil, true
//or else, false
func LimitsTimeoutValidation(timeout *int) bool {
	if timeout == nil {
		return true
	}
	if *timeout < 100 {
		// Do not allow invalid limit to be added to API
		wskprint.PrintlnOpenWhiskWarning(wski18n.T(wski18n.ID_WARN_LIMITS_TIMEOUT))
		return false
	} else if *timeout > 600000 {
		// Emit a warning, but allow to pass through to provider
		wskprint.PrintlnOpenWhiskWarning(wski18n.T(wski18n.ID_WARN_LIMITS_TIMEOUT))
	}
	return true
}

//if valid or nil, true
//or else, false
func LimitsMemoryValidation(memory *int) bool {
	if memory == nil {
		return true
	}
	if *memory < 128 {
		// Do not allow invalid limit to be added to API
		wskprint.PrintlnOpenWhiskWarning(wski18n.T(wski18n.ID_WARN_LIMITS_MEMORY_SIZE))
		return false
	} else if *memory > 2048 {
		// Emit a warning, but allow to pass through to provider
		wskprint.PrintlnOpenWhiskWarning(wski18n.T(wski18n.ID_WARN_LIMITS_MEMORY_SIZE))
	}
	return true
}

//if valid or nil, true
//or else, false
func LimitsLogsizeValidation(logsize *int) bool {
	if logsize == nil {
		return true
	}
	if *logsize < 0 {
		// Do not allow invalid limit to be added to API
		wskprint.PrintlnOpenWhiskWarning(wski18n.T(wski18n.ID_WARN_LIMITS_LOG_SIZE))
		return false
	} else if *logsize > 10 {
		// Emit a warning, but allow to pass through to provider
		wskprint.PrintlnOpenWhiskWarning(wski18n.T(wski18n.ID_WARN_LIMITS_LOG_SIZE))
	}
	return true
}

func NotSupportLimits(value *int, name string) {
	if value != nil {
		warningString := wski18n.T(
			wski18n.ID_WARN_LIMIT_UNCHANGEABLE_X_name_X,
			map[string]interface{}{wski18n.KEY_NAME: name})
		wskprint.PrintlnOpenWhiskWarning(warningString)
	}
}
