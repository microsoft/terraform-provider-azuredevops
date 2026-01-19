package core_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance/checks"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance/planchecks"
	"github.com/microsoft/terraform-provider-azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/internal/errorutil"
)

type ProjectResource struct{}

func (p ProjectResource) Exists(ctx context.Context, client *client.Client, state *terraform.InstanceState) (bool, error) {
	_, err := client.CoreClient.GetProject(ctx, core.GetProjectArgs{ProjectId: &state.ID})
	if err == nil {
		return true, nil
	}
	if errorutil.WasNotFound(err) {
		return false, nil
	}
	return false, err
}

func TestAccProject_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_project", "test")
	r := ProjectResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "process_template_id"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccProject_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_project", "test")
	r := ProjectResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "process_template_id"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccProject_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_project", "test")
	r := ProjectResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "process_template_id"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					planchecks.IsNotResourceAction(data.ResourceAddr(), plancheck.ResourceActionReplace),
				},
			},
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "process_template_id"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccProject_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_project", "test")
	r := ProjectResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config:      r.requiresImport(data),
			ExpectError: regexp.MustCompile(`TF200019: The following project already exists`),
		},
	})
}

func (r ProjectResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "acctest-%[1]s"
  description        = "test description"
}`, data.RandomString)
}

func (r ProjectResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "acctest-%[1]s-update"
  description        = "test description update"
  version_control    = "Git"
}`, data.RandomString)
}

func (r ProjectResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "acctest-%[1]s"
  description        = "test description complete"
  version_control    = "Tfvc"
  work_item_template = "Agile"
}`, data.RandomString)
}

func (r ProjectResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_project" "import" {
  name               = azuredevops_project.test.name
  description        = azuredevops_project.test.description
}`, r.basic(data))
}
