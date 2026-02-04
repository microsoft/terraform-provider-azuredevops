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
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/errorutil"
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
			Config: r.complete(data),
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

func TestAccProject_features(t *testing.T) {
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
			Config: r.features(data, true),
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
		{
			Config: r.features(data, false),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "process_template_id"),
			),
		},
		data.ImportStep(),
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

func TestAccProject_migrateFromV0(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_project", "test")
	r := ProjectResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.migrateV0(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
			),
		},
		data.MigratePlanStep(r.migrateNow(data)),
	})
}

func (r ProjectResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "acctest-%[1]s"
}`, data.RandomString)
}

func (r ProjectResource) features(data acceptance.TestData, v bool) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "acctest-%[1]s"
  features = {
    boards     = %[2]t
    repositories      = %[2]t
    pipelines  = %[2]t
    testplans  = %[2]t
    artifacts  = %[2]t
  }
}`, data.RandomString, v)
}

func (r ProjectResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "acctest-%[1]s"
  description        = "test description"
  version_control    = "Git"
  work_item_template = "Basic"
  features = {
    boards     = false
    repositories      = false
    pipelines  = false
    testplans  = false
    artifacts  = false
  }
}`, data.RandomString)
}

func (r ProjectResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_project" "import" {
  name = azuredevops_project.test.name
}`, r.basic(data))
}

func (r ProjectResource) migrateV0(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  provider = azuredevops-v1

  name               = "acctest-%[1]s"
  description        = "test description"
  version_control    = "Git"
  work_item_template = "Basic"
  features = {
    boards     		= "disabled"
    repositories    = "enabled"
  }
}

# This is to accomodate the eventual consistency issue of the ADO API,
# which is not handled by the v1 azuredevops provider.
resource "time_sleep" "wait_30_seconds" {
  depends_on = [azuredevops_project.test]
  create_duration = "30s"
}

`, data.RandomString)
}

func (r ProjectResource) migrateNow(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "acctest-%[1]s"
  description        = "test description"
  version_control    = "Git"
  work_item_template = "Basic"
  features = {
    boards     		= false
    repositories    = true
  }
}`, data.RandomString)
}
