//go:build (all || resource_servicehook_webhook_tfs) && !exclude_servicehook

package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/servicehooks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func TestAccServicehookWebhookTfs_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	url := "https://example.com/webhook"

	resourceType := "azuredevops_servicehook_webhook_tfs"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookWebhookTfsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookWebhookTfsResourceWithGitPushEvent(projectName, url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "url", url),
					resource.TestCheckResourceAttr(tfCheckNode, "git_push.#", "1"),
				),
			},
		},
	})
}

func TestAccServicehookWebhookTfs_Update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	url1 := "https://example.com/webhook"
	url2 := "https://example.org/webhook2"

	resourceType := "azuredevops_servicehook_webhook_tfs"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookWebhookTfsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookWebhookTfsResourceWithGitPushEvent(projectName, url1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "url", url1),
					resource.TestCheckResourceAttr(tfCheckNode, "git_push.#", "1"),
				),
			},
			{
				Config: testutils.HclServicehookWebhookTfsResourceWithGitPushEvent(projectName, url2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "url", url2),
					resource.TestCheckResourceAttr(tfCheckNode, "git_push.#", "1"),
				),
			},
		},
	})
}

func TestAccServicehookWebhookTfs_ChangeEventType(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	url := "https://example.com/webhook"

	resourceType := "azuredevops_servicehook_webhook_tfs"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookWebhookTfsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookWebhookTfsResourceWithGitPushEvent(projectName, url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "git_push.#", "1"),
				),
			},
			{
				Config: testutils.HclServicehookWebhookTfsResourceWithWorkItemCreatedEvent(projectName, url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "work_item_created.#", "1"),
					resource.TestCheckNoResourceAttr(tfCheckNode, "git_push.#"),
				),
			},
		},
	})
}

func TestAccServicehookWebhookTfs_WithHeaders(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	url := "https://example.com/webhook"

	resourceType := "azuredevops_servicehook_webhook_tfs"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookWebhookTfsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookWebhookTfsResourceWithHeadersAndAuth(projectName, url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "url", url),
					resource.TestCheckResourceAttr(tfCheckNode, "http_headers.%", "2"),
					resource.TestCheckResourceAttr(tfCheckNode, "http_headers.X-Custom-Header", "Test Value"),
					resource.TestCheckResourceAttr(tfCheckNode, "basic_auth_username", "testuser"),
				),
			},
		},
	})
}

func TestAccServicehookWebhookTfs_InvalidUrl(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	invalidUrl := "not-a-valid-url"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookWebhookTfsDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      testutils.HclServicehookWebhookTfsResourceWithGitPushEvent(projectName, invalidUrl),
				ExpectError: regexp.MustCompile("expected \"url\" to have a host"),
			},
		},
	})
}

func TestAccServicehookWebhookTfs_WithResourceDetails(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	url := "https://example.com/webhook"

	resourceType := "azuredevops_servicehook_webhook_tfs"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookWebhookTfsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookWebhookTfsResourceWithResourceDetails(projectName, url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "url", url),
					resource.TestCheckResourceAttr(tfCheckNode, "resource_details_to_send", "minimal"),
					resource.TestCheckResourceAttr(tfCheckNode, "messages_to_send", "text"),
					resource.TestCheckResourceAttr(tfCheckNode, "detailed_messages_to_send", "html"),
				),
			},
		},
	})
}

func CheckServicehookWebhookTfsDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_servicehook_webhook_tfs" {
			continue
		}

		subscriptionID := resource.Primary.ID
		if subscriptionID == "" {
			return fmt.Errorf("No service hook subscription ID is set")
		}

		// Parse the subscription ID as UUID
		subscriptionUUID, err := uuid.Parse(subscriptionID)
		if err != nil {
			return fmt.Errorf("Invalid subscription ID format: %v", err)
		}

		// Try to get the subscription from the service - if it exists, this should fail the test
		// If it's been deleted, this should return an error (which is expected)
		_, err = clients.ServiceHooksClient.GetSubscription(clients.Ctx, servicehooks.GetSubscriptionArgs{
			SubscriptionId: &subscriptionUUID,
		})

		if err == nil {
			return fmt.Errorf("Unexpectedly found a service hook subscription that should be deleted")
		}
	}

	return nil
}
