package graph_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance/checks"
	"github.com/microsoft/terraform-provider-azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/errorutil"
	"github.com/microsoft/terraform-provider-azuredevops/internal/utils/pointer"
)

type GroupMembershipResource struct{}

func (p GroupMembershipResource) Exists(ctx context.Context, client *client.Client, state *terraform.InstanceState) (bool, error) {
	err := client.GraphClient.CheckMembershipExistence(ctx, graph.CheckMembershipExistenceArgs{
		SubjectDescriptor:   pointer.From(state.Attributes["member_id"]),
		ContainerDescriptor: pointer.From(state.Attributes["group_id"]),
	})
	if err == nil {
		return true, nil
	}
	if errorutil.WasNotFound(err) {
		return false, nil
	}
	return false, err
}

func TestAccGroupMembership_group(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_group_membership", "test")
	r := GroupMembershipResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.group(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
			),
		},
		data.ImportStep(),
	})
}

func TestAccGroupMembership_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azuredevops_group_membership", "test")
	r := GroupMembershipResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.group(data),
			Check: resource.ComposeTestCheckFunc(
				checks.ExistsInAzure(t, r, data.ResourceAddr()),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.ResourceExistError(),
		},
	})
}

func (r GroupMembershipResource) group(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_group" "container" {
  display_name = "acctest-%[1]s"
  lifecycle {
	  ignore_changes = [members]
  }
}

resource "azuredevops_group" "member" {
  display_name = "acctest-member-%[1]s"
}

resource "azuredevops_group_membership" "test" {
  group_id = azuredevops_group.container.id
  member_id = azuredevops_group.member.id
}
`, data.RandomString)
}

func (r GroupMembershipResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_group_membership" "import" {
  group_id = azuredevops_group_membership.test.group_id
  member_id = azuredevops_group_membership.test.member_id
}
`, r.group(data))
}
