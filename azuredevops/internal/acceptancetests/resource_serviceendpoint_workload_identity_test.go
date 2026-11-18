//go:build (all || resource_serviceendpoint_workload_identity) && !exclude_serviceendpoints

package acceptancetests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// These tests require an existing identity with organization access and a caller
// that can configure its federated credential. They create real connections and
// may cause Azure DevOps to create credentials on that identity.
func TestAccServiceEndpointWorkloadIdentity_application(t *testing.T) {
	testAccServiceEndpointWorkloadIdentity(t, "Application")
}

func TestAccServiceEndpointWorkloadIdentity_managedIdentity(t *testing.T) {
	testAccServiceEndpointWorkloadIdentity(t, "ManagedIdentity")
}

func testAccServiceEndpointWorkloadIdentity(t *testing.T, principalType string) {
	projectName := testutils.GenerateResourceName()
	endpointName := testutils.GenerateResourceName()
	resourceType := "azuredevops_serviceendpoint_workload_identity"
	node := resourceType + ".test"
	envVars := []string{"AZDO_WIF_TENANT_ID", "AZDO_WIF_CLIENT_ID", "AZDO_WIF_OBJECT_ID", "AZDO_WIF_SUBJECT_DESCRIPTOR"}
	if principalType == "ManagedIdentity" {
		envVars = append(envVars, "AZDO_WIF_MANAGED_IDENTITY_RESOURCE_ID")
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, &envVars) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointWorkloadIdentity(projectName, endpointName, principalType),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(node, endpointName),
					resource.TestCheckResourceAttr(node, "service_principal_type", principalType),
					resource.TestCheckResourceAttr(node, "authorization.scheme", "WorkloadIdentityFederation"),
					resource.TestCheckResourceAttr(node, "is_ready", "true"),
					resource.TestCheckResourceAttrSet(node, "workload_identity_federation_issuer"),
					resource.TestCheckResourceAttrSet(node, "workload_identity_federation_subject"),
				),
			},
			{
				Config: hclServiceEndpointWorkloadIdentity(projectName, endpointName+"-updated", principalType),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(node, endpointName+"-updated"),
					resource.TestCheckResourceAttr(node, "is_ready", "true"),
				),
			},
			{
				ResourceName:            node,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(node),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"draft"},
			},
		},
	})
}

func hclServiceEndpointWorkloadIdentity(projectName, endpointName, principalType string) string {
	managedIdentity := ""
	if principalType == "ManagedIdentity" {
		managedIdentity = fmt.Sprintf("managed_identity_resource_id = %q", os.Getenv("AZDO_WIF_MANAGED_IDENTITY_RESOURCE_ID"))
	}
	return testutils.HclProjectResource(projectName) + fmt.Sprintf(`
resource "azuredevops_serviceendpoint_workload_identity" "test" {
  project_id             = azuredevops_project.project.id
  service_endpoint_name  = %q
  org_url                = %q
  tenant_id              = %q
  service_principal_id   = %q
  object_id              = %q
  subject_descriptor     = %q
  service_principal_type = %q
  %s
}
`, endpointName, os.Getenv("AZDO_ORG_SERVICE_URL"), os.Getenv("AZDO_WIF_TENANT_ID"), os.Getenv("AZDO_WIF_CLIENT_ID"), os.Getenv("AZDO_WIF_OBJECT_ID"), os.Getenv("AZDO_WIF_SUBJECT_DESCRIPTOR"), principalType, managedIdentity)
}
