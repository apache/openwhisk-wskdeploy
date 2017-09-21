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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLicenseLocalValidation(t *testing.T) {
	assert.False(t, LicenseLocalValidation("---"))
	assert.True(t, LicenseLocalValidation("Apache-2.0"))
	assert.True(t, LicenseLocalValidation("0BSD"))
	assert.True(t, LicenseLocalValidation("0bsd"))
	assert.True(t, LicenseLocalValidation("MIT"))
	assert.True(t, LicenseLocalValidation("MIT-feh"))
	assert.True(t, LicenseLocalValidation("BSD-3-Clause-LBNL"))
	assert.False(t, LicenseLocalValidation("GPL-3.0+"))
	assert.False(t, LicenseLocalValidation("Zimbra-1.3"))
	assert.False(t, LicenseLocalValidation("xpp"))
}

func TestLicenseRemoteValidation(t *testing.T) {
	assert.False(t, LicenseRemoteValidation("---"))
	assert.True(t, LicenseRemoteValidation("Apache-2.0"))
	assert.True(t, LicenseRemoteValidation("0BSD"))
	assert.True(t, LicenseRemoteValidation("0bsd"))
	assert.True(t, LicenseRemoteValidation("MIT"))
	assert.True(t, LicenseRemoteValidation("MIT-feh"))
	assert.True(t, LicenseRemoteValidation("BSD-3-Clause-LBNL"))
	assert.True(t, LicenseRemoteValidation("GPL-3.0+"))
	assert.True(t, LicenseRemoteValidation("Zimbra-1.3"))
	assert.True(t, LicenseRemoteValidation("xpp"))
}

func TestCheckLicense(t *testing.T) {
	assert.False(t, CheckLicense("---"))
	assert.True(t, CheckLicense("Apache-2.0"))
	assert.True(t, CheckLicense("0BSD"))
	assert.True(t, CheckLicense("0bsd"))
	assert.True(t, CheckLicense("MIT"))
	assert.True(t, CheckLicense("MIT-feh"))
	assert.True(t, CheckLicense("BSD-3-Clause-LBNL"))
	assert.True(t, CheckLicense("GPL-3.0+"))
	assert.True(t, CheckLicense("Zimbra-1.3"))
	assert.True(t, CheckLicense("xpp"))
}
