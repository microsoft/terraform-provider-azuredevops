package acceptancetests

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointAzureRM_dataSource_with_serviceEndpointID(t *testing.T) {
	serviceEndpointAzureRMName := testutils.GenerateResourceName()
	serviceprincipalid := uuid.New().String()
	serviceprincipalkey := uuid.New().String()
	projectName := testutils.GenerateResourceName()
	serviceEndpointAuthenticationScheme := "ServicePrincipal"
	createServiceEndpointAzureRMWithServiceEndpointIDData := fmt.Sprintf("%s\n%s",
		testutils.HclServiceEndpointAzureRMResource(projectName, serviceEndpointAzureRMName, serviceprincipalid, serviceprincipalkey, serviceEndpointAuthenticationScheme),
		testutils.HclServiceEndpointAzureRMDataSourceWithServiceEndpointID(),
	)

	tfNode := "data.azuredevops_serviceendpoint_azurerm.serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createServiceEndpointAzureRMWithServiceEndpointIDData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_name", serviceEndpointAzureRMName),
					resource.TestCheckResourceAttrSet(tfNode, "service_endpoint_id"),
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_authentication_scheme", serviceEndpointAuthenticationScheme),
				),
			},
		},
	})
}

func TestAccServiceEndpointAzureRM_dataSource_with_serviceEndpointName(t *testing.T) {
	serviceEndpointAzureRMName := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	serviceprincipalid := uuid.New().String()
	serviceprincipalkey := uuid.New().String()
	serviceEndpointAuthenticationScheme := "ServicePrincipal"
	createServiceEndpointAzureRMWithServiceEndpointNameData := fmt.Sprintf("%s\n%s",
		testutils.HclServiceEndpointAzureRMResource(projectName, serviceEndpointAzureRMName, serviceprincipalid, serviceprincipalkey, serviceEndpointAuthenticationScheme),
		testutils.HclServiceEndpointAzureRMDataSourceWithServiceEndpointName(serviceEndpointAzureRMName),
	)

	tfNode := "data.azuredevops_serviceendpoint_azurerm.serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createServiceEndpointAzureRMWithServiceEndpointNameData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_name", serviceEndpointAzureRMName),
					resource.TestCheckResourceAttrSet(tfNode, "service_endpoint_id"),
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_authentication_scheme", serviceEndpointAuthenticationScheme),
				),
			},
		},
	})
}

func TestAccServiceEndpointAzureRM_dataSource_with_WorkloadIdentityFederation(t *testing.T) {
	serviceEndpointAzureRMName := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	serviceprincipalid := uuid.New().String()
	serviceEndpointAuthenticationScheme := "WorkloadIdentityFederation"

	azureDevOpsOrgName := "terraform-provider-azuredevops"

	if os.Getenv("AZDO_ORG_SERVICE_URL") != "" {
		azureDevOpsOrgUrl, err := url.Parse(os.Getenv("AZDO_ORG_SERVICE_URL"))
		if err != nil {
			t.Fatal(err)
		}
		azureDevOpsOrgName = path.Base(azureDevOpsOrgUrl.Path)
	}

	createServiceEndpointAzureRMWithServiceEndpointNameData := fmt.Sprintf("%s\n%s",
		testutils.HclServiceEndpointAzureRMNoKeyResource(projectName, serviceEndpointAzureRMName, serviceprincipalid, serviceEndpointAuthenticationScheme),
		testutils.HclServiceEndpointAzureRMDataSourceWithServiceEndpointName(serviceEndpointAzureRMName),
	)

	tfNode := "data.azuredevops_serviceendpoint_azurerm.serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createServiceEndpointAzureRMWithServiceEndpointNameData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_name", serviceEndpointAzureRMName),
					resource.TestCheckResourceAttrSet(tfNode, "service_endpoint_id"),
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_authentication_scheme", serviceEndpointAuthenticationScheme),
					resource.TestMatchResourceAttr(tfNode, "workload_identity_federation_issuer", regexp.MustCompile("^https://vstoken.dev.azure.com/[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					resource.TestCheckResourceAttr(tfNode, "workload_identity_federation_subject", fmt.Sprintf("sc://%s/%s/%s", azureDevOpsOrgName, projectName, serviceEndpointAzureRMName)),
				),
			},
		},
	})
}
