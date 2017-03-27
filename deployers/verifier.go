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

package deployers

import (
	"fmt"

	"github.com/openwhisk/openwhisk-client-go/whisk"
)

// The verifier will filter the deployer against the target DeploymentApplication
// the deployer will query the OpenWhisk platform for already deployed entities.
// We assume the deployer and the manifest are targeted for the same namespace.
type Verifier struct {
}

type Filter interface {
	// Perform some filter.
	Filter(deployer *ServiceDeployer, target *DeploymentApplication) (filtered *DeploymentApplication, err error)
	// Perform some queries.
	Query(deployer *ServiceDeployer) (da *DeploymentApplication, err error)
}

func (vf *Verifier) Query(deployer *ServiceDeployer) (da *DeploymentApplication, err error) {
	pkgoptions := &whisk.PackageListOptions{false, 0, 0, 0, false}
	packages, _, err := deployer.Client.Packages.List(pkgoptions)

	da = NewDeploymentApplication()
	for _, pa := range packages {
		deppack := NewDeploymentPackage()
		deppack.Package, _ = convert(&pa)
		da.Packages[pa.Name] = deppack
	}
	return da, nil
}

func (vf *Verifier) Filter(deployer *ServiceDeployer, target *DeploymentApplication) (rs *DeploymentApplication, err error) {
	//substract
	for _, pa := range target.Packages {
		for _, dpa := range deployer.Deployment.Packages {
			if pa.Package.Name == dpa.Package.Name {
				delete(target.Packages, dpa.Package.Name)
			}
		}
	}

	depApp := NewDeploymentApplication()
	fmt.Printf("Target Packages are %#v\n", target.Packages)
	depApp.Packages = target.Packages
	return depApp, nil
}

// Convert whisk.package to whisk.SentPackageNoPublish
func convert(pa *whisk.Package) (sentpackage *whisk.Package, err error) {
	sp := &whisk.Package{}
	sp.Name = pa.Name
	sp.Annotations = pa.Annotations
	sp.Parameters = pa.Parameters
	sp.Version = pa.Version
	return sp, nil
}
