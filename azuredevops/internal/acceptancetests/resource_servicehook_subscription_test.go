//go:build (all || resource_servicehook_subscription) && !exclude_subscriptions

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func TestAccServicehookSubscription_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	resourceType := "azuredevops_servicehook_subscription"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServicehookSubscriptionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclServicehookSubscriptionResourceBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "publisher_id", "tfs"),
					resource.TestCheckResourceAttr(tfCheckNode, "event_type", "workitem.created"),
					resource.TestCheckResourceAttr(tfCheckNode, "consumer_id", "webHooks"),
					resource.TestCheckResourceAttr(tfCheckNode, "consumer_action_id", "httpRequest"),
					resource.TestCheckResourceAttr(tfCheckNode, "status", "enabled"),
				),
			},
		},
	})
}

func TestAccServicehookSubscription_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	resourceType := "azuredevops_servicehook_subscription"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServicehookSubscriptionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclServicehookSubscriptionResourceBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "status", "enabled"),
				),
			},
			{
				Config: hclServicehookSubscriptionResourceUpdated(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "status", "disabled"),
				),
			},
		},
	})
}

func checkServicehookSubscriptionDestroyed(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azuredevops_servicehook_subscription" {
			continue
		}

		// Check if subscription actually exists in Azure DevOps
		if _, err := getServicehookSubscriptionFromResource(rs); err == nil {
			return fmt.Errorf("Unexpectedly found a service hook subscription that should be deleted")
		}
	}

	return nil
}

// given a resource from the state, return a servicehook subscription (and error)
func getServicehookSubscriptionFromResource(resource *terraform.ResourceState) (*servicehooks.Subscription, error) {
	subscriptionID, err := uuid.Parse(resource.Primary.ID)
	if err != nil {
		return nil, err
	}

	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
	return clients.ServiceHooksClient.GetSubscription(clients.Ctx, servicehooks.GetSubscriptionArgs{
		SubscriptionId: &subscriptionID,
	})
}

func hclServicehookSubscriptionResourceBasic(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%s"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
}

resource "azuredevops_servicehook_subscription" "test" {
  project_id         = azuredevops_project.test.id
  publisher_id       = "tfs"
  event_type         = "workitem.created"
  consumer_id        = "webHooks"
  consumer_action_id = "httpRequest"

  publisher_inputs = {
    workItemType = "Bug"
  }

  consumer_inputs = {
    url = "https://example.com/webhook"
  }

  status = "enabled"
}
`, projectName)
}

func hclServicehookSubscriptionResourceUpdated(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%s"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
}

resource "azuredevops_servicehook_subscription" "test" {
  project_id         = azuredevops_project.test.id
  publisher_id       = "tfs"
  event_type         = "workitem.created"
  consumer_id        = "webHooks"
  consumer_action_id = "httpRequest"

  publisher_inputs = {
    workItemType = "Task"
  }

  consumer_inputs = {
    url = "https://example.com/updated-webhook"
  }

  status = "disabled"
}
`, projectName)
}
