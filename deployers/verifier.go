package deployers

import (
	"fmt"

	"github.com/openwhisk/openwhisk-client-go/whisk"
	"github.com/openwhisk/openwhisk-wskdeploy/utils"
)

// The verifier will filter the deployer against the target DeploymentApplication
// the deployer will query the OpenWhisk platform for already deployed entities.
// We assume the deployer and the manifest are targeted for the same namespace.
type Verifier struct {
}

type Filter interface {
	// Perform some filter.
	// Target is what's been deployed on OpenWhisk platform, we use A to represent it.
	// Deployer is what's needs to be undeployed, we use B to represent it.
	// Thera are several use cases we need to consider, for examples:
	// 1. A is not intersected with B, we just return.
	// 2. A is a subset of B, we undeploy A.
	// 3. A and B is partial intersected, we undeploy the intersection part of A.
	// 4. B is subset of A, we undeploy B.
	// 5. B equals A, which is the same as use case 2, we undeploy A.
	// .. Other use case not considered?
	Filter(deployer *ServiceDeployer, target *DeploymentApplication) (filtered *DeploymentApplication, err error)
	// Perform some queries.
	Query(deployer *ServiceDeployer) (da *DeploymentApplication, err error)
}

func (vf *Verifier) Query(deployer *ServiceDeployer) (da *DeploymentApplication, err error) {

	//The results needs to be return
	da = NewDeploymentApplication()
	// query what packages been deployed on OpenWhisk platform
	pkgoptions := &whisk.PackageListOptions{false, 0, 0, 0, false}
	packages, _, err := deployer.Client.Packages.List(pkgoptions)

	// query what actions been deployed on OpenWhisk platform
	acnoptions := &whisk.ActionListOptions{0, 0, false}
	for _, pkg := range packages {
		actions, _, err := deployer.Client.Actions.List(pkg.Name, acnoptions)
		utils.Check(err)
		convertAction(actions)
	}


	for _, pa := range packages {
		deppack := NewDeploymentPackage()
		deppack.Package, _ = convertPackage(&pa)
		da.Packages[pa.Name] = deppack
		// Insert the actions into packages.
	}
	return da, nil
}

func (vf *Verifier) Filter(deployer *ServiceDeployer, target *DeploymentApplication) (rs *DeploymentApplication, err error) {
	// Cover all the use cases in Filter interface
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
func convertPackage(pa *whisk.Package) (sentpackage *whisk.SentPackageNoPublish, err error) {
	sp := &whisk.SentPackageNoPublish{}
	sp.Name = pa.Name
	sp.Annotations = pa.Annotations
	sp.Parameters = pa.Parameters
	sp.Version = pa.Version
	return sp, nil
}
// Convert whisk.actions from OpenWhisk to ActionRecords.
func convertAction(ac *whisk.Action) (sentaction *whisk.SentActionNoPublish, err error) {
        return nil, nil
}

// Convert rules, triggers? Probably not necessary?
