//go:build (all || resource_serviceendpoint_azurecr) && !exclude_serviceendpoints

package acceptancetests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointAzureCR_spn_basic(t *testing.T) {
	if os.Getenv("TEST_ARM_SUBSCRIPTION_ID") == "" || os.Getenv("TEST_ARM_SUBSCRIPTION_NAME") == "" ||
		os.Getenv("TEST_ARM_TENANT_ID") == "" || os.Getenv("TEST_ARM_RESOURCE_GROUP") == "" || os.Getenv("TEST_ARM_ACR_NAME") == "" {
		t.Skip("Skip test as `TEST_ARM_SUBSCRIPTION_ID` or `TEST_ARM_SUBSCRIPTION_NAME` or `TEST_ARM_TENANT_ID` or `TEST_ARM_RESOURCE_GROUP` or `TEST_ARM_ACR_NAME` is not set")
	}

	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_azurecr"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclAzureCRSpn(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			}, {
				ResourceName:      tfSvcEpNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceEndpointAzureCR_spn_update(t *testing.T) {
	if os.Getenv("TEST_ARM_SUBSCRIPTION_ID") == "" || os.Getenv("TEST_ARM_SUBSCRIPTION_NAME") == "" ||
		os.Getenv("TEST_ARM_TENANT_ID") == "" || os.Getenv("TEST_ARM_RESOURCE_GROUP") == "" || os.Getenv("TEST_ARM_ACR_NAME") == "" {
		t.Skip("Skip test as `TEST_ARM_SUBSCRIPTION_ID` or `TEST_ARM_SUBSCRIPTION_NAME` or `TEST_ARM_TENANT_ID` or `TEST_ARM_RESOURCE_GROUP` or `TEST_ARM_ACR_NAME` is not set")
	}

	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_azurecr"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclAzureCRSpn(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			},
			{
				ResourceName:      tfSvcEpNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: hclAzureCRSpn(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
				),
			},
			{
				ResourceName:      tfSvcEpNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceEndpointAzureCR_workLoadIdentity_basic(t *testing.T) {
	if os.Getenv("TEST_ARM_SUBSCRIPTION_ID") == "" || os.Getenv("TEST_ARM_SUBSCRIPTION_NAME") == "" ||
		os.Getenv("TEST_ARM_TENANT_ID") == "" || os.Getenv("TEST_ARM_RESOURCE_GROUP") == "" || os.Getenv("TEST_ARM_ACR_NAME") == "" {
		t.Skip("Skip test as `TEST_ARM_SUBSCRIPTION_ID` or `TEST_ARM_SUBSCRIPTION_NAME` or `TEST_ARM_TENANT_ID` or `TEST_ARM_RESOURCE_GROUP` or `TEST_ARM_ACR_NAME` is not set")
	}

	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_azurecr"
	tfSvcEpNode := resourceType + ".test"
	// fmt.Println(hclAzureCRWorkLoadIdentity(projectName, serviceEndpointNameFirst))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclAzureCRWorkLoadIdentity(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			}, {
				ResourceName:      tfSvcEpNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceEndpointAzureCR_workLoadIdentity_update(t *testing.T) {
	if os.Getenv("TEST_ARM_SUBSCRIPTION_ID") == "" || os.Getenv("TEST_ARM_SUBSCRIPTION_NAME") == "" ||
		os.Getenv("TEST_ARM_TENANT_ID") == "" || os.Getenv("TEST_ARM_RESOURCE_GROUP") == "" || os.Getenv("TEST_ARM_ACR_NAME") == "" {
		t.Skip("Skip test as `TEST_ARM_SUBSCRIPTION_ID` or `TEST_ARM_SUBSCRIPTION_NAME` or `TEST_ARM_TENANT_ID` or `TEST_ARM_RESOURCE_GROUP` or `TEST_ARM_ACR_NAME` is not set")
	}

	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_azurecr"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclAzureCRWorkLoadIdentity(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			},
			{
				ResourceName:      tfSvcEpNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: hclAzureCRWorkLoadIdentity(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
				),
			},
			{
				ResourceName:      tfSvcEpNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclAzureCRSpn(projectName, serviceConnectionName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_azurecr" "test" {
  project_id                             = azuredevops_project.test.id
  service_endpoint_authentication_scheme = "ServicePrincipal"
  service_endpoint_name                  = "%s"
  azurecr_spn_tenantid                   = "%s"
  azurecr_subscription_id                = "%s"
  azurecr_subscription_name              = "%s"
  resource_group                         = "%s"
  azurecr_name                           = "%s"
}
`, projectName, serviceConnectionName, os.Getenv("TEST_ARM_TENANT_ID"), os.Getenv("TEST_ARM_SUBSCRIPTION_ID"),
		os.Getenv("TEST_ARM_SUBSCRIPTION_NAME"), os.Getenv("TEST_ARM_RESOURCE_GROUP"), os.Getenv("TEST_ARM_ACR_NAME"))
}

func hclAzureCRWorkLoadIdentity(projectName, serviceConnectionName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_azurecr" "test" {
  project_id                             = azuredevops_project.test.id
  service_endpoint_authentication_scheme = "WorkloadIdentityFederation"
  service_endpoint_name                  = "%s"
  azurecr_spn_tenantid                   = "%s"
  azurecr_subscription_id                = "%s"
  azurecr_subscription_name              = "%s"
  resource_group                         = "%s"
  azurecr_name                           = "%s"
}
`, projectName, serviceConnectionName, os.Getenv("TEST_ARM_TENANT_ID"), os.Getenv("TEST_ARM_SUBSCRIPTION_ID"),
		os.Getenv("TEST_ARM_SUBSCRIPTION_NAME"), os.Getenv("TEST_ARM_RESOURCE_GROUP"), os.Getenv("TEST_ARM_ACR_NAME"))
}
