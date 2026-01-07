package core_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance/checks"
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
			Check:  checks.ExistsInAzure(t, r, data.ResourceAddr()),
		},
		//data.ImportStep(),
	})
}

func (r ProjectResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "acctest-%[1]s"
  description        = "test description"
}`, data.RandomString)
}
