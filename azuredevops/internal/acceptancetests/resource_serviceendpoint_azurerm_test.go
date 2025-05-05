//go:build (all || resource_serviceendpoint_azurerm) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_azurerm
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccServiceEndpointAzureRm_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()
	serviceprincipalidFirst := uuid.New().String()
	serviceprincipalidSecond := uuid.New().String()
	serviceprincipalkeyFirst := uuid.New().String()
	serviceprincipalkeySecond := uuid.New().String()
	serviceEndpointAuthenticationScheme := "ServicePrincipal"

	resourceType := "azuredevops_serviceendpoint_azurerm"
	tfSvcEpNode := resourceType + ".serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointAzureRMResource(projectName, serviceEndpointNameFirst, serviceprincipalidFirst, serviceprincipalkeyFirst, serviceEndpointAuthenticationScheme),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "service_principal_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalid", serviceprincipalidFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalkey", serviceprincipalkeyFirst),
				),
			}, {
				Config: testutils.HclServiceEndpointAzureRMResource(projectName, serviceEndpointNameSecond, serviceprincipalidSecond, serviceprincipalkeySecond, serviceEndpointAuthenticationScheme),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "service_principal_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalid", serviceprincipalidSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalkey", serviceprincipalkeySecond),
				),
			},
		},
	})
}

func TestAccServiceEndpointAzureRm_CreateAndUpdate_WithValidate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceprincipalidFirst := uuid.New().String()
	serviceprincipalkeyFirst := uuid.New().String()
	serviceEndpointAuthenticationScheme := "ServicePrincipal"
	validateFirst := false
	validateSecond := true

	resourceType := "azuredevops_serviceendpoint_azurerm"
	tfSvcEpNode := resourceType + ".serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointAzureRMResourceWithValidate(projectName, serviceEndpointNameFirst, serviceprincipalidFirst, serviceprincipalkeyFirst, serviceEndpointAuthenticationScheme, validateFirst),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalid", serviceprincipalidFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalkey", serviceprincipalkeyFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "features.0.validate", strconv.FormatBool(validateFirst)),
				),
			}, {
				Config:      testutils.HclServiceEndpointAzureRMResourceWithValidate(projectName, serviceEndpointNameFirst, serviceprincipalidFirst, serviceprincipalkeyFirst, serviceEndpointAuthenticationScheme, validateSecond),
				ExpectError: regexp.MustCompile("Failed to obtain the Json Web Token"),
			},
		},
	})
}

func TestAccServiceEndpointAzureRm_Create_WithValidate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceprincipalid := uuid.New().String()
	serviceprincipalkey := uuid.New().String()
	serviceEndpointAuthenticationScheme := "ServicePrincipal"
	validate := true

	resourceType := "azuredevops_serviceendpoint_azurerm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config:      testutils.HclServiceEndpointAzureRMResourceWithValidate(projectName, serviceEndpointName, serviceprincipalid, serviceprincipalkey, serviceEndpointAuthenticationScheme, validate),
				ExpectError: regexp.MustCompile("Failed to obtain the Json Web Token"),
			},
		},
	})
}

func TestAccServiceEndpointAzureRm_MgmtGrpCreateAndUpdate(t *testing.T) {
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
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalkey", serviceprincipalkey),
				),
			},
		},
	})
}

func TestAccServiceEndpointAzureRm_AutomaticCreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	serviceEndpointAuthenticationScheme := "ServicePrincipal"

	subscriptionId := os.Getenv("TEST_ARM_SUBSCRIPTION_ID")
	subscriptionName := os.Getenv("TEST_ARM_SUBSCRIPTION_NAME")
	tenantId := os.Getenv("TEST_ARM_TENANT_ID")

	resourceType := "azuredevops_serviceendpoint_azurerm"
	tfSvcEpNode := resourceType + ".serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, &[]string{
				"TEST_ARM_SUBSCRIPTION_ID",
				"TEST_ARM_SUBSCRIPTION_NAME",
				"TEST_ARM_TENANT_ID",
			})
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointAzureRMAutomaticResourceWithProject(projectName, serviceEndpointName, serviceEndpointAuthenticationScheme, subscriptionId, subscriptionName, tenantId),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "service_principal_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckNoResourceAttr(tfSvcEpNode, "credentials.0"),
				),
			},
			{
				Config: testutils.HclServiceEndpointAzureRMAutomaticResourceWithProject(projectName, serviceEndpointName, serviceEndpointAuthenticationScheme, subscriptionId, subscriptionName, tenantId),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "service_principal_id"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointName),
					resource.TestCheckNoResourceAttr(tfSvcEpNode, "credentials.0"),
				),
			},
		},
	})
}

// validates that a manual workload federation service endpoint can be created and updated
func TestAccServiceEndpointAzureRm_WorkloadFederation_Manual_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()
	serviceprincipalidFirst := uuid.New().String()
	serviceprincipalidSecond := uuid.New().String()
	serviceEndpointAuthenticationScheme := "WorkloadIdentityFederation"

	azureDevOpsOrgName := "terraform-provider-azuredevops"

	if os.Getenv("AZDO_ORG_SERVICE_URL") != "" {
		azureDevOpsOrgUrl, _ := url.Parse(os.Getenv("AZDO_ORG_SERVICE_URL"))
		azureDevOpsOrgName = path.Base(azureDevOpsOrgUrl.Path)
	}

	resourceType := "azuredevops_serviceendpoint_azurerm"
	tfSvcEpNode := resourceType + ".serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointAzureRMNoKeyResource(projectName, serviceEndpointNameFirst, serviceprincipalidFirst, serviceEndpointAuthenticationScheme),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalid", serviceprincipalidFirst),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_authentication_scheme", serviceEndpointAuthenticationScheme),
					resource.TestMatchResourceAttr(tfSvcEpNode, "workload_identity_federation_issuer", regexp.MustCompile("^https://vstoken.dev.azure.com/[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					resource.TestCheckResourceAttr(tfSvcEpNode, "workload_identity_federation_subject", fmt.Sprintf("sc://%s/%s/%s", azureDevOpsOrgName, projectName, serviceEndpointNameFirst)),
				),
			}, {
				Config: testutils.HclServiceEndpointAzureRMNoKeyResource(projectName, serviceEndpointNameSecond, serviceprincipalidSecond, serviceEndpointAuthenticationScheme),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "credentials.0.serviceprincipalid", serviceprincipalidSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_authentication_scheme", serviceEndpointAuthenticationScheme),
					resource.TestMatchResourceAttr(tfSvcEpNode, "workload_identity_federation_issuer", regexp.MustCompile("^https://vstoken.dev.azure.com/[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					resource.TestCheckResourceAttr(tfSvcEpNode, "workload_identity_federation_subject", fmt.Sprintf("sc://%s/%s/%s", azureDevOpsOrgName, projectName, serviceEndpointNameSecond)),
				),
			},
		},
	})
}

// validates that an automatic workload federation service endpoint can be created and updated
func TestAccServiceEndpointAzureRm_WorkloadFederation_Automatic_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()
	serviceEndpointAuthenticationScheme := "WorkloadIdentityFederation"

	subscriptionId := os.Getenv("TEST_ARM_SUBSCRIPTION_ID")
	subscriptionName := os.Getenv("TEST_ARM_SUBSCRIPTION_NAME")
	tenantId := os.Getenv("TEST_ARM_TENANT_ID")

	resourceType := "azuredevops_serviceendpoint_azurerm"
	tfSvcEpNode := resourceType + ".serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, &[]string{
				"TEST_ARM_SUBSCRIPTION_ID",
				"TEST_ARM_SUBSCRIPTION_NAME",
				"TEST_ARM_TENANT_ID",
			})
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointAzureRMAutomaticResourceWithProject(projectName, serviceEndpointNameFirst, serviceEndpointAuthenticationScheme, subscriptionId, subscriptionName, tenantId),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_authentication_scheme", serviceEndpointAuthenticationScheme),
				),
			}, {
				Config: testutils.HclServiceEndpointAzureRMAutomaticResourceWithProject(projectName, serviceEndpointNameSecond, serviceEndpointAuthenticationScheme, subscriptionId, subscriptionName, tenantId),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_authentication_scheme", serviceEndpointAuthenticationScheme),
				),
			},
		},
	})
}

// validates that an managed identity service endpoint can be created and updated
func TestAccServiceEndpointAzureRm_ManagedServiceIdentity_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()
	serviceEndpointAuthenticationScheme := "ManagedServiceIdentity"

	tenantId := "9c59cbe5-2ca1-4516-b303-8968a070edd2"
	subscriptionId := "3b0fee91-c36d-4d70-b1e9-fc4b9d608c3d"
	subscriptionName := "Microsoft Azure DEMO"

	resourceType := "azuredevops_serviceendpoint_azurerm"
	tfSvcEpNode := resourceType + ".serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointAzureRMAutomaticResourceWithProject(projectName, serviceEndpointNameFirst, serviceEndpointAuthenticationScheme, subscriptionId, subscriptionName, tenantId),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_authentication_scheme", serviceEndpointAuthenticationScheme),
				),
			}, {
				Config: testutils.HclServiceEndpointAzureRMAutomaticResourceWithProject(projectName, serviceEndpointNameSecond, serviceEndpointAuthenticationScheme, subscriptionId, subscriptionName, tenantId),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_authentication_scheme", serviceEndpointAuthenticationScheme),
				),
			},
		},
	})
}

func TestAccServiceEndpointAzureRm_azureStack(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	tfNode := "azuredevops_serviceendpoint_azurerm.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed("azuredevops_serviceendpoint_azurerm"),
		Steps: []resource.TestStep{
			{
				Config: hclAzureRMServiceEndpointEnvironmentAzureStack(projectName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceEndpointExistsWithName(tfNode, serviceEndpointName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "azurerm_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfNode, "azurerm_subscription_id"),
					resource.TestCheckResourceAttrSet(tfNode, "azurerm_subscription_name"),
					resource.TestCheckResourceAttrSet(tfNode, "server_url"),
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_name", serviceEndpointName),
				),
			},
			{
				ResourceName:            tfNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials.0.serviceprincipalkey"},
			},
		},
	})
}

func hclAzureRMServiceEndpointEnvironmentAzureStack(projectName, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_azurerm" "test" {
  project_id                = azuredevops_project.test.id
  service_endpoint_name     = "%s"
  environment               = "AzureStack"
  server_url                = "https://www.azuredevops.com"
  azurerm_spn_tenantid      = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id   = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name = "Test Sub"
  credentials {
    serviceprincipalid  = "00000000-0000-0000-0000-000000000000"
    serviceprincipalkey = "00000000-0000-0000-0000-000000000000"
  }
}
`, projectName, serviceEndpointName)
}
