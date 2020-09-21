// +build all resource_serviceendpoint_devops
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointAzureDevOps_Defaults(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_devops"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, &[]string{"AZDO_PERSONAL_ACCESS_TOKEN"}) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: azdoResourceSetupSimple(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
		},
	})
}


func TestAccServiceEndpointAzureDevOps_PersonalTokenBasic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_devops"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, &[]string{"AZDO_PERSONAL_ACCESS_TOKEN"}) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: azdoPersonTokenConfigBasic(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_personal.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", "Managed by Terraform"),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
				),
			},
		},
	})
}

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccServiceEndpointAzureDevOps_PersonalTokenUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()
	description := "Manage by Terraform Update"
	organization := "example"

	resourceType := "azuredevops_serviceendpoint_devops"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, &[]string{"AZDO_PERSONAL_ACCESS_TOKEN"}) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: azdoPersonTokenConfigBasic(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_personal.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "organization_name", organization),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			}, {
				Config: azdoPersonTokenConfigUpdate(projectName, serviceEndpointNameSecond, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "auth_personal.#", "1"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "description", description),
					resource.TestCheckResourceAttr(tfSvcEpNode, "organization_name", organization),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
				),
			},
		},
	})
}

func azdoResourceSetupSimple(projectName string, serviceEndpointName string) string {
	projectResource := testutils.HclProjectResource(projectName)
	serviceEndpointResource := testutils.HclServiceEndpointAzureDevOpsResourceSimple(serviceEndpointName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func azdoPersonTokenConfigBasic(projectName string, serviceEndpointName string) string {
	projectResource := testutils.HclProjectResource(projectName)

	serviceEndpointResource := testutils.HclServiceEndpointAzureDevOpsResource(
		serviceEndpointName,
		"test_token_basic",
		"Managed by Terraform",
	)

	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}

func azdoPersonTokenConfigUpdate(projectName string, serviceEndpointName string, updatedDescription string) string {
	projectResource := testutils.HclProjectResource(projectName)

	serviceEndpointResource := testutils.HclServiceEndpointAzureDevOpsResource(
		serviceEndpointName,
		"test_token_update",
		updatedDescription,
	)

	return fmt.Sprintf("%s\n%s", projectResource, serviceEndpointResource)
}
