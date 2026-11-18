//go:build (all || resource_serviceendpoint_workload_identity) && !exclude_serviceendpoints

package acceptancetests

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

// These tests require an existing identity with organization access and a caller
// that can configure its federated credential. They create real connections and
// may cause Azure DevOps to create credentials on that identity.
// AZDO_WIF_OBJECT_ID is the app registration object ID for Application, or the
// managed identity's principal object ID for ManagedIdentity.
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
	envVars := []string{"AZDO_ORG_SERVICE_URL", "AZDO_WIF_TENANT_ID", "AZDO_WIF_CLIENT_ID", "AZDO_WIF_OBJECT_ID", "AZDO_WIF_SUBJECT_DESCRIPTOR"}
	if principalType == "ManagedIdentity" {
		envVars = append(envVars, "AZDO_WIF_MANAGED_IDENTITY_RESOURCE_ID")
	}
	resource.ParallelTest(t, resource.TestCase{
		// Allow any supported provider authentication, including Azure CLI and
		// app-only credentials. A PAT is not required for these tests.
		PreCheck: func() {
			for _, name := range envVars {
				if os.Getenv(name) == "" {
					t.Fatalf("%s must be set for this acceptance test", name)
				}
			}
		},
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testAccCheckServiceEndpointWorkloadIdentityDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointWorkloadIdentity(projectName, endpointName, principalType, true),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(node, endpointName),
					resource.TestCheckResourceAttr(node, "service_principal_type", principalType),
					resource.TestCheckResourceAttr(node, "authorization.scheme", "WorkloadIdentityFederation"),
					resource.TestCheckResourceAttrSet(node, "workload_identity_federation_issuer"),
					resource.TestCheckResourceAttrSet(node, "workload_identity_federation_subject"),
				),
			},
			{
				Config: hclServiceEndpointWorkloadIdentity(projectName, endpointName+"-updated", principalType, false),
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

func testAccCheckServiceEndpointWorkloadIdentityDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azuredevops_serviceendpoint_workload_identity" {
			continue
		}
		id, err := uuid.Parse(rs.Primary.ID)
		if err != nil {
			return err
		}
		project := rs.Primary.Attributes["project_id"]
		ep, err := clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, serviceendpoint.GetServiceEndpointDetailsArgs{Project: &project, EndpointId: &id})
		if utils.ResponseWasNotFound(err) {
			continue
		}
		if err != nil {
			return err
		}
		if ep != nil && *ep != (serviceendpoint.ServiceEndpoint{}) {
			return fmt.Errorf("workload identity service endpoint %s still exists", id)
		}
	}
	return nil
}

func hclServiceEndpointWorkloadIdentity(projectName, endpointName, principalType string, draft bool) string {
	managedIdentity := ""
	if principalType == "ManagedIdentity" {
		managedIdentity = fmt.Sprintf("managed_identity_resource_id = %q", os.Getenv("AZDO_WIF_MANAGED_IDENTITY_RESOURCE_ID"))
	}
	return testutils.HclProjectResource(projectName) + fmt.Sprintf(`
data "azuredevops_group" "readers" {
  project_id = azuredevops_project.project.id
  name       = "Readers"
}

resource "azuredevops_group_membership" "identity" {
  group   = data.azuredevops_group.readers.descriptor
  members = [%q]
  mode    = "add"
}

resource "azuredevops_serviceendpoint_workload_identity" "test" {
  project_id             = azuredevops_project.project.id
  service_endpoint_name  = %q
  org_url                = %q
  tenant_id              = %q
  service_principal_id   = %q
  object_id              = %q
  subject_descriptor     = %q
  service_principal_type = %q
  draft                  = %t
  %s

  depends_on = [azuredevops_group_membership.identity]
}
`, os.Getenv("AZDO_WIF_SUBJECT_DESCRIPTOR"), endpointName, os.Getenv("AZDO_ORG_SERVICE_URL"), os.Getenv("AZDO_WIF_TENANT_ID"), os.Getenv("AZDO_WIF_CLIENT_ID"), os.Getenv("AZDO_WIF_OBJECT_ID"), os.Getenv("AZDO_WIF_SUBJECT_DESCRIPTOR"), principalType, draft, managedIdentity)
}
