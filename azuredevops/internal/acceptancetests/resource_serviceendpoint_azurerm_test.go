//go:build (all || resource_serviceendpoint_azurerm) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_azurerm
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccServiceEndpointAzureRm_CreateAndUpdate(t *testing.T) {
	t.Skip("Skipping test TestAccServiceEndpointAzureRm_CreateAndUpdate: test resource limit")
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()
	serviceprincipalidFirst := uuid.New().String()
	serviceprincipalidSecond := uuid.New().String()
	serviceprincipalkeyFirst := uuid.New().String()
	serviceprincipalkeySecond := uuid.New().String()
	serviceEndpointType := "ServicePrincipal"

	resourceType := "azuredevops_serviceendpoint_azurerm"
	tfSvcEpNode := resourceType + ".serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointAzureRMResource(projectName, serviceEndpointNameFirst, serviceprincipalidFirst, serviceprincipalkeyFirst, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalid", serviceprincipalidFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "credentials.0.serviceprincipalkey_hash"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalkey", serviceprincipalkeyFirst),
				),
			}, {
				Config: testutils.HclServiceEndpointAzureRMResource(projectName, serviceEndpointNameSecond, serviceprincipalidSecond, serviceprincipalkeySecond, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalid", serviceprincipalidSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "credentials.0.serviceprincipalkey_hash"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalkey", serviceprincipalkeySecond),
				),
			},
		},
	})
}

func TestAccServiceEndpointAzureRm_MgmtGrpCreateAndUpdate(t *testing.T) {
	t.Skip("Skipping test TestAccServiceEndpointAzureRm_MgmtGrpCreateAndUpdate: test resource limit")
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceprincipalid := uuid.New().String()
	serviceprincipalkey := uuid.New().String()

	resourceType := "azuredevops_serviceendpoint_azurerm"
	tfSvcEpNode := resourceType + ".serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointAzureRMResourceWithMG(projectName, serviceEndpointName, serviceprincipalid, serviceprincipalkey),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_management_group_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_management_group_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalid", serviceprincipalid),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "credentials.0.serviceprincipalkey_hash"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalkey", serviceprincipalkey),
				),
			},
		},
	})
}

func TestAccServiceEndpointAzureRm_AutomaticCreateAndUpdate(t *testing.T) {
	t.Skip("Skipping test TestAccServiceEndpointAzureRm_AutomaticCreateAndUpdate: test resource limit")

	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointType := "ServicePrincipal"

	resourceType := "azuredevops_serviceendpoint_azurerm"
	tfSvcEpNode := resourceType + ".serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointAzureRMAutomaticResourceWithProject(projectName, serviceEndpointName, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckNoResourceAttr(tfSvcEpNode, "credentials.0"),
				),
			},
			{
				Config: testutils.HclServiceEndpointAzureRMAutomaticResourceWithProject(projectName, serviceEndpointName, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckNoResourceAttr(tfSvcEpNode, "credentials.0"),
				),
			},
		},
	})
}

// validates that a manual workload federation service endpoint can be created and updated
func TestAccServiceEndpointAzureRm_WorkloadFederation_Manual_CreateAndUpdate(t *testing.T) {
	t.Skip("Skipping test TestAccServiceEndpointAzureRm_CreateAndUpdate: test resource limit")
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()
	serviceprincipalidFirst := uuid.New().String()
	serviceprincipalidSecond := uuid.New().String()
	serviceprincipalkeyFirst := uuid.New().String()
	serviceprincipalkeySecond := uuid.New().String()
	serviceEndpointType := "WorkloadIdentityFederation"

	resourceType := "azuredevops_serviceendpoint_azurerm"
	tfSvcEpNode := resourceType + ".serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointAzureRMResource(projectName, serviceEndpointNameFirst, serviceprincipalidFirst, serviceprincipalkeyFirst, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalid", serviceprincipalidFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "azurerm_service_endpoint_type", serviceEndpointType),
				),
			}, {
				Config: testutils.HclServiceEndpointAzureRMResource(projectName, serviceEndpointNameSecond, serviceprincipalidSecond, serviceprincipalkeySecond, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalid", serviceprincipalidSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "azurerm_service_endpoint_type", serviceEndpointType),
				),
			},
		},
	})
}

// validates that an automatic workload federation service endpoint can be created and updated
func TestAccServiceEndpointAzureRm_WorkloadFederation_Automatic_CreateAndUpdate(t *testing.T) {
	t.Skip("Skipping test TestAccServiceEndpointAzureRm_CreateAndUpdate: test resource limit")
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()
	serviceEndpointType := "WorkloadIdentityFederation"

	resourceType := "azuredevops_serviceendpoint_azurerm"
	tfSvcEpNode := resourceType + ".serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointAzureRMAutomaticResourceWithProject(projectName, serviceEndpointNameFirst, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "azurerm_service_endpoint_type", serviceEndpointType),
				),
			}, {
				Config: testutils.HclServiceEndpointAzureRMAutomaticResourceWithProject(projectName, serviceEndpointNameSecond, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "azurerm_service_endpoint_type", serviceEndpointType),
				),
			},
		},
	})
}

// validates that an managed identity service endpoint can be created and updated
func TestAccServiceEndpointAzureRm_ManagedIdentity_CreateAndUpdate(t *testing.T) {
	t.Skip("Skipping test TestAccServiceEndpointAzureRm_CreateAndUpdate: test resource limit")
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()
	serviceEndpointType := "ManagedIdentity"

	resourceType := "azuredevops_serviceendpoint_azurerm"
	tfSvcEpNode := resourceType + ".serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointAzureRMAutomaticResourceWithProject(projectName, serviceEndpointNameFirst, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "azurerm_service_endpoint_type", serviceEndpointType),
				),
			}, {
				Config: testutils.HclServiceEndpointAzureRMAutomaticResourceWithProject(projectName, serviceEndpointNameSecond, serviceEndpointType),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "azurerm_service_endpoint_type", serviceEndpointType),
				),
			},
		},
	})
}
