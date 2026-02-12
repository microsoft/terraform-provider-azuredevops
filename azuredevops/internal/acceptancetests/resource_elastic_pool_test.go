package acceptancetests

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/elastic"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func TestAccElasticPool_basic(t *testing.T) {
	poolName := testutils.GenerateResourceName()
	tfNode := "azuredevops_elastic_pool.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, &[]string{"TEST_SPN_ID", "TEST_SPN_SECRET", "TEST_TENANT_ID", "TEST_SUB_ID", "TEST_SUB_NAME", "TEST_AZURE_VMSS_ID"})
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkElasticPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclElasticPoolBasic(poolName,
					os.Getenv("TEST_SPN_ID"), os.Getenv("TEST_SPN_SECRET"),
					os.Getenv("TEST_TENANT_ID"), os.Getenv("TEST_SUB_ID"),
					os.Getenv("TEST_SUB_NAME"), os.Getenv("TEST_AZURE_VMSS_ID")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", poolName),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccElasticPool_update(t *testing.T) {
	poolName := testutils.GenerateResourceName()
	tfNode := "azuredevops_elastic_pool.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, &[]string{"TEST_SPN_ID", "TEST_SPN_SECRET", "TEST_TENANT_ID", "TEST_SUB_ID", "TEST_SUB_NAME", "TEST_AZURE_VMSS_ID"})
		}, Providers: testutils.GetProviders(),
		CheckDestroy: checkElasticPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclElasticPoolBasic(poolName,
					os.Getenv("TEST_SPN_ID"), os.Getenv("TEST_SPN_SECRET"),
					os.Getenv("TEST_TENANT_ID"), os.Getenv("TEST_SUB_ID"),
					os.Getenv("TEST_SUB_NAME"), os.Getenv("TEST_AZURE_VMSS_ID")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", poolName),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: hclElasticPoolUpdate(poolName,
					os.Getenv("TEST_SPN_ID"),
					os.Getenv("TEST_SPN_SECRET"),
					os.Getenv("TEST_TENANT_ID"),
					os.Getenv("TEST_SUB_ID"),
					os.Getenv("TEST_SUB_NAME"),
					os.Getenv("TEST_AZURE_VMSS_ID")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", poolName),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccElasticPool_requiresImportErrorStep(t *testing.T) {
	poolName := testutils.GenerateResourceName()
	tfNode := "azuredevops_elastic_pool.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, &[]string{"TEST_SPN_ID", "TEST_SPN_SECRET", "TEST_TENANT_ID", "TEST_SUB_ID", "TEST_SUB_NAME", "TEST_AZURE_VMSS_ID"})
		}, Providers: testutils.GetProviders(),
		CheckDestroy: checkElasticPoolDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclElasticPoolBasic(poolName,
					os.Getenv("TEST_SPN_ID"), os.Getenv("TEST_SPN_SECRET"), os.Getenv("TEST_TENANT_ID"),
					os.Getenv("TEST_SUB_ID"), os.Getenv("TEST_SUB_NAME"), os.Getenv("TEST_AZURE_VMSS_ID")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", poolName),
				),
			},

			{
				Config: hclElasticPoolResourceRequiresImport(poolName,
					os.Getenv("TEST_SPN_ID"), os.Getenv("TEST_SPN_SECRET"), os.Getenv("TEST_TENANT_ID"),
					os.Getenv("TEST_SUB_ID"), os.Getenv("TEST_SUB_NAME"), os.Getenv("TEST_AZURE_VMSS_ID")),
				ExpectError: requiresElasticPoolImportError(poolName),
			},
		},
	})
}

func checkElasticPoolDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_elastic_pool" {
			continue
		}

		id, err := strconv.Atoi(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Elastic Pool ID=%d cannot be parsed!. Error=%v", id, err)
		}

		if _, err := clients.ElasticClient.GetElasticPool(clients.Ctx, elastic.GetElasticPoolArgs{PoolId: &id}); err == nil {
			return fmt.Errorf("Elastic Pool ID %d should not exist", id)
		}
	}
	return nil
}

func requiresElasticPoolImportError(resourceName string) *regexp.Regexp {
	message := " creating Elastic Pool: Agent pool %[1]s already exists."
	return regexp.MustCompile(fmt.Sprintf(message, resourceName))
}

func hclElasticPoolTemplate(name, spnId, spnSecret, tenantId, subId, subName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_azurerm" "test" {
  project_id                             = azuredevops_project.test.id
  service_endpoint_name                  = "%[1]s"
  description                            = "Managed by Terraform"
  service_endpoint_authentication_scheme = "ServicePrincipal"
  credentials {
    serviceprincipalid  = "%[2]s"
    serviceprincipalkey = "%[3]s"
  }
  azurerm_spn_tenantid      = "%[4]s"
  azurerm_subscription_id   = "%[5]s"
  azurerm_subscription_name = "%[6]s"
}
`, name, spnId, spnSecret, tenantId, subId, subName)
}

func hclElasticPoolBasic(name, spnId, spnSecret, tenantId, subId, subName, vmssId string) string {
	template := hclElasticPoolTemplate(name, spnId, spnSecret, tenantId, subId, subName)
	return fmt.Sprintf(`


%[1]s

resource "azuredevops_elastic_pool" "test" {
  name                   = "%[2]s"
  service_endpoint_id    = azuredevops_serviceendpoint_azurerm.test.id
  service_endpoint_scope = azuredevops_project.test.id
  desired_idle           = 3
  max_capacity           = 3
  azure_resource_id      = "%[3]s"
}`, template, name, vmssId)
}

func hclElasticPoolUpdate(name, spnId, spnSecret, tenantId, subId, subName, vmssId string) string {
	template := hclElasticPoolTemplate(name, spnId, spnSecret, tenantId, subId, subName)
	return fmt.Sprintf(`

%[1]s

resource "azuredevops_elastic_pool" "test" {
  name = "%[2]s"

  service_endpoint_id    = azuredevops_serviceendpoint_azurerm.test.id
  service_endpoint_scope = azuredevops_project.test.id
  desired_idle           = 3
  max_capacity           = 3

  recycle_after_each_use = true
  agent_interactive_ui   = true
  time_to_live_minutes   = 40

  auto_provision = true
  auto_update    = false


  azure_resource_id = "%[3]s"
}`, template, name, vmssId)
}

func hclElasticPoolResourceRequiresImport(name, spnId, spnSecret, tenantId, subId, subName, vmssId string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_elastic_pool" "import" {
  name                   = azuredevops_elastic_pool.test.name
  service_endpoint_id    = azuredevops_elastic_pool.test.service_endpoint_id
  service_endpoint_scope = azuredevops_elastic_pool.test.service_endpoint_scope
  desired_idle           = azuredevops_elastic_pool.test.desired_idle
  max_capacity           = azuredevops_elastic_pool.test.max_capacity
  recycle_after_each_use = azuredevops_elastic_pool.test.recycle_after_each_use
  agent_interactive_ui   = azuredevops_elastic_pool.test.agent_interactive_ui
  time_to_live_minutes   = azuredevops_elastic_pool.test.time_to_live_minutes
  auto_provision         = azuredevops_elastic_pool.test.auto_provision
  auto_update            = azuredevops_elastic_pool.test.auto_update
  azure_resource_id      = azuredevops_elastic_pool.test.azure_resource_id
}`, hclElasticPoolBasic(name, spnId, spnSecret, tenantId, subId, subName, vmssId))
}
