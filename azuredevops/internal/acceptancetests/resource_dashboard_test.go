package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/dashboard"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func TestAccDashboard_project_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	name := testutils.GenerateResourceName()

	tfNode := "azuredevops_dashboard.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkDashboardDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDashboardProjectBasic(projectName, name),
				Check: resource.ComposeTestCheckFunc(
					checkDashboardExist(name),
					resource.TestCheckResourceAttr(tfNode, "name", name),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "owner_id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: importDashboardId(tfNode),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDashboard_project_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	name := testutils.GenerateResourceName()

	tfNode := "azuredevops_dashboard.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkDashboardDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDashboardProjectBasic(projectName, name),
				Check: resource.ComposeTestCheckFunc(
					checkDashboardExist(name),
					resource.TestCheckResourceAttr(tfNode, "name", name),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "owner_id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: importDashboardId(tfNode),
				ImportStateVerify: true,
			},
			{
				Config: hclDashboardProjectUpdate(projectName, name),
				Check: resource.ComposeTestCheckFunc(
					checkDashboardExist(name+"update"),
					resource.TestCheckResourceAttr(tfNode, "name", name+"update"),
					resource.TestCheckResourceAttr(tfNode, "description", "description"),
					resource.TestCheckResourceAttr(tfNode, "refresh_interval", "5"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: importDashboardId(tfNode),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDashboard_project_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	name := testutils.GenerateResourceName()

	tfNode := "azuredevops_dashboard.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkDashboardDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDashboardProjectComplete(projectName, name),
				Check: resource.ComposeTestCheckFunc(
					checkDashboardExist(name),
					resource.TestCheckResourceAttr(tfNode, "name", name),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "owner_id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: importDashboardId(tfNode),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDashboard_team_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	name := testutils.GenerateResourceName()

	tfNode := "azuredevops_dashboard.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkDashboardDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDashboardTeamBasic(projectName, name),
				Check: resource.ComposeTestCheckFunc(
					checkDashboardExist(name),
					resource.TestCheckResourceAttr(tfNode, "name", name),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "owner_id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: importDashboardId(tfNode),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDashboard_team_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	name := testutils.GenerateResourceName()

	tfNode := "azuredevops_dashboard.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkDashboardDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDashboardTeamBasic(projectName, name),
				Check: resource.ComposeTestCheckFunc(
					checkDashboardExist(name),
					resource.TestCheckResourceAttr(tfNode, "name", name),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "team_id"),
					resource.TestCheckResourceAttrSet(tfNode, "owner_id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: importDashboardId(tfNode),
				ImportStateVerify: true,
			},
			{
				Config: hclDashboardTeamUpdate(projectName, name),
				Check: resource.ComposeTestCheckFunc(
					checkDashboardExist(name+"update"),
					resource.TestCheckResourceAttr(tfNode, "name", name+"update"),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "team_id"),
					resource.TestCheckResourceAttrSet(tfNode, "owner_id"),
					resource.TestCheckResourceAttr(tfNode, "description", "description"),
					resource.TestCheckResourceAttr(tfNode, "refresh_interval", "5"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: importDashboardId(tfNode),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDashboard_team_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	name := testutils.GenerateResourceName()

	tfNode := "azuredevops_dashboard.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkDashboardDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDashboardTeamComplete(projectName, name),
				Check: resource.ComposeTestCheckFunc(
					checkDashboardExist(name),
					resource.TestCheckResourceAttr(tfNode, "name", name),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "team_id"),
					resource.TestCheckResourceAttrSet(tfNode, "description"),
					resource.TestCheckResourceAttrSet(tfNode, "refresh_interval"),
					resource.TestCheckResourceAttrSet(tfNode, "owner_id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: importDashboardId(tfNode),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDashboard_team_requireImportError(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	name := testutils.GenerateResourceName()

	tfNode := "azuredevops_dashboard.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkDashboardDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDashboardTeamBasic(projectName, name),
				Check: resource.ComposeTestCheckFunc(
					checkDashboardExist(name),
					resource.TestCheckResourceAttr(tfNode, "name", name),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "team_id"),
					resource.TestCheckResourceAttrSet(tfNode, "owner_id"),
				),
			},
			{
				ExpectError: regexp.MustCompile("Creating dashboard in Azure DevOps: VS403345: Duplicate name on dashboard. Each dashboard held by a team must use a distinct name"),
				Config:      hclDashboardTeamRequireImport(projectName, name),
			},
		},
	})
}

func checkDashboardDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	// verify that every project referenced in the state does not exist in AzDO
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_dashboard" {
			continue
		}

		id := resource.Primary.ID
		dashboardId, err := uuid.Parse(id)
		if err != nil {
			return fmt.Errorf("Parsing dashboard ID: %+v", err)
		}

		args := dashboard.GetDashboardArgs{
			Project:     converter.String(resource.Primary.Attributes["project_id"]),
			DashboardId: &dashboardId,
		}

		if v, ok := resource.Primary.Attributes["team_id"]; ok {
			args.Team = converter.String(v)
		}
		// indicates the project still exists - this should fail the test
		if _, err = clients.DashboardClient.GetDashboard(clients.Ctx, args); err == nil {
			return fmt.Errorf("Dashboard with ID %s should not exist", id)
		}
	}
	return nil
}

func checkDashboardExist(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources["azuredevops_dashboard.test"]
		if !ok {
			return fmt.Errorf("Did not find a `azuredevops_dashboard` in the Terraform state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

		id := res.Primary.ID
		dashboardId, err := uuid.Parse(id)
		if err != nil {
			return fmt.Errorf("Parsing dashboard ID %s", id)
		}

		args := dashboard.GetDashboardArgs{
			Project:     converter.String(res.Primary.Attributes["project_id"]),
			DashboardId: &dashboardId,
		}

		if v, ok := res.Primary.Attributes["team_id"]; ok {
			args.Team = converter.String(v)
		}

		resp, err := clients.DashboardClient.GetDashboard(clients.Ctx, args)
		if err != nil {
			return fmt.Errorf("Dashboard with ID: %s cannot be found!. Error: %v", id, err)
		}

		if *resp.Name != expectedName {
			return fmt.Errorf("Dashboard with ID: %s has Name: %s, but expected Name: %s", id, *resp.Name, expectedName)
		}
		return nil
	}
}

func importDashboardId(resourceType string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		if res, ok := s.RootModule().Resources[resourceType]; ok {
			projectId := res.Primary.Attributes["project_id"]
			if teamId, ok := res.Primary.Attributes["team_id"]; ok {
				return fmt.Sprintf("%s/%s/%s", projectId, teamId, res.Primary.ID), nil
			}
			return fmt.Sprintf("%s/%s", projectId, res.Primary.ID), nil
		}
		return "", fmt.Errorf("Not found: %s", resourceType)
	}
}

func hclDashboardProjectBasic(projectName, name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_dashboard" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
}
`, projectName, name)
}

func hclDashboardProjectUpdate(projectName, name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_dashboard" "test" {
  project_id       = azuredevops_project.test.id
  name             = "%supdate"
  description      = "description"
  refresh_interval = 5
}
`, projectName, name)
}

func hclDashboardProjectComplete(projectName, name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_dashboard" "test" {
  project_id       = azuredevops_project.test.id
  name             = "%s"
  description      = "description"
  refresh_interval = 5
}
`, projectName, name)
}

func hclDashboardTeamBasic(projectName, name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%[1]s"
}

resource "azuredevops_team" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s dashboard"
}

resource "azuredevops_dashboard" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  team_id    = azuredevops_team.test.id
}
`, projectName, name)
}

func hclDashboardTeamUpdate(projectName, name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%[1]s"
}

resource "azuredevops_team" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s dashboard"
}

resource "azuredevops_dashboard" "test" {
  project_id       = azuredevops_project.test.id
  name             = "%[2]supdate"
  team_id          = azuredevops_team.test.id
  description      = "description"
  refresh_interval = 5
}
`, projectName, name)
}

func hclDashboardTeamComplete(projectName, name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%[1]s"
}

resource "azuredevops_team" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s dashboard"
}

resource "azuredevops_dashboard" "test" {
  project_id       = azuredevops_project.test.id
  name             = "%[2]s"
  team_id          = azuredevops_team.test.id
  description      = "description"
  refresh_interval = 5
}
`, projectName, name)
}

func hclDashboardTeamRequireImport(projectName, name string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_dashboard" "import" {
  project_id = azuredevops_dashboard.test.project_id
  name       = azuredevops_dashboard.test.name
  team_id    = azuredevops_dashboard.test.team_id
}
`, hclDashboardTeamBasic(projectName, name))
}
