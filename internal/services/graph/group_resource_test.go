package graph_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance/checks"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance/planchecks"
	"github.com/microsoft/terraform-provider-azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/internal/errorutil"
)

type GroupResource struct{}

func (p GroupResource) Exists(ctx context.Context, client *client.Client, state *terraform.InstanceState) (bool, error) {
	_, err := client.GraphClient.GetGroup(ctx, graph.GetGroupArgs{GroupDescriptor: &state.ID})
	if err == nil {
		return true, nil
	}
	if errorutil.WasNotFound(err) {
		return false, nil
	}
	return false, err
}

func TestAccGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "url"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "origin"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "subject_kind"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "domain"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "principal_name"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "scope"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_scopeProject(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.scopeProject(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "url"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "origin"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "subject_kind"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "domain"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "principal_name"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "url"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "origin"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "subject_kind"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "domain"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "principal_name"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "scope"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "url"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "origin"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "subject_kind"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "domain"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "principal_name"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "scope"),
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
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "url"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "origin"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "subject_kind"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "domain"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "principal_name"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "scope"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "url"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "origin"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "subject_kind"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "domain"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "principal_name"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "scope"),
			),
		},
		data.ImportStep(),
	})
}

func (r GroupResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_group" "test" {
  display_name = "acctest-%[1]s"
}
`, data.RandomString)
}

func (r GroupResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_group" "test" {
  display_name = "acctest-%[1]s-complete"
  description = "description"
}
`, data.RandomString)
}

func (r GroupResource) scopeProject(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "acctest-%[1]s"
}

resource "azuredevops_group" "test" {
  scope = azuredevops_project.test.id
  display_name = "acctest-%[1]s"
}
`, data.RandomString)
}

// func (r GroupResource) update(data acceptance.TestData) string {
// 	return fmt.Sprintf(`
// resource "azuredevops_project" "test" {
//   name               = "acctest-%[1]s-update"
//   description        = "test description"
//   version_control    = "Git"
// }`, data.RandomString)
// }

// func (r GroupResource) complete(data acceptance.TestData) string {
// 	return fmt.Sprintf(`
// resource "azuredevops_project" "test" {
//   name               = "acctest-%[1]s"
//   description        = "test description"
//   version_control    = "Tfvc"
//   work_item_template = "Agile"
// }`, data.RandomString)
// }
