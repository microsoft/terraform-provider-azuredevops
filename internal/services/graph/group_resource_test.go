package graph_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance/checks"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance/planchecks"
	"github.com/microsoft/terraform-provider-azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/errorutil"
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

func TestAccGroup_vsts_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.vstsBasic(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "url"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "origin"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "subject_kind"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "domain"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "principal_name"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "scope"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "storage_key"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_vsts_scopeProject(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.vstsScopeProject(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "url"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "origin"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "subject_kind"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "domain"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "principal_name"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "storage_key"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_vsts_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.vstsComplete(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "url"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "origin"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "subject_kind"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "domain"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "principal_name"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "scope"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "storage_key"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_vsts_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_group", "test")
	r := GroupResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.vstsBasic(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "url"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "origin"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "subject_kind"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "domain"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "principal_name"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "scope"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "storage_key"),
			),
		},
		data.ImportStep(),
		{
			Config: r.vstsComplete(data),
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
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "storage_key"),
			),
		},
		data.ImportStep(),
		{
			Config: r.vstsBasic(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "url"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "origin"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "subject_kind"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "domain"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "principal_name"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "scope"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "storage_key"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_aad_originId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_group", "test")
	r := GroupResource{}

	if os.Getenv("ARM_TENANT_ID") == "" {
		t.Skip("AzureAD related environment variables are not specified.")
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.aadOriginId(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "url"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "origin"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "subject_kind"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "domain"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "principal_name"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "scope"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "storage_key"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroup_aad_mail(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_group", "test")
	r := GroupResource{}

	if os.Getenv("ARM_TENANT_ID") == "" {
		t.Skip("AzureAD related environment variables are not specified.")
	}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.aadMail(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "url"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "origin"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "subject_kind"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "domain"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "principal_name"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "scope"),
				resource.TestCheckResourceAttrSet(data.ResourceAddr(), "storage_key"),
			),
		},
		data.ImportStep(),
	})
}

func (r GroupResource) vstsBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_group" "test" {
  display_name = "acctest-%[1]s"
}
`, data.RandomString)
}

func (r GroupResource) vstsComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_group" "test" {
  display_name = "acctest-%[1]s-complete"
  description = "description"
  members = [
  	azuredevops_group.member1.id,
  	azuredevops_group.member2.id,
  ]
}

resource "azuredevops_group" "member1" {
  display_name = "acctest-%[1]s-member1"
}
resource "azuredevops_group" "member2" {
  display_name = "acctest-%[1]s-member2"
}
`, data.RandomString)
}

func (r GroupResource) vstsScopeProject(data acceptance.TestData) string {
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

func (r GroupResource) aadOriginId(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azuread_client_config" "current" {}

resource "azuread_group" "test" {
  display_name     = "acctest-%[1]s"
  owners           = [data.azuread_client_config.current.object_id]
  security_enabled = true
}

resource "azuredevops_group" "test" {
  origin_id = azuread_group.test.object_id
}
`, data.RandomString)
}

func (r GroupResource) aadMail(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azuread_client_config" "current" {}

resource "azuread_group" "test" {
  display_name     = "acctest-%[1]s"
  mail_enabled     = true
  mail_nickname    = "Acctest-%[1]s"
  types            = ["Unified"]
  owners           = [data.azuread_client_config.current.object_id]
  security_enabled = true
}

resource "azuredevops_group" "test" {
  mail = azuread_group.test.mail
}
`, data.RandomString)
}
