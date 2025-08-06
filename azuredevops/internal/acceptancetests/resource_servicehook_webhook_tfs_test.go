//go:build (all || resource_servicehook_webhook_tfs) && !exclude_servicehook

package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
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
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_servicehook_webhook_tfs" {
			continue
		}

		// indicates the subscription still exists - this should fail the test
		if _, err := getServiceHookSubscriptionFromResource(resource); err == nil {
			return fmt.Errorf("Unexpectedly found a service hook subscription that should be deleted")
		}
	}

	return nil
}

// Helper function to get the service hook subscription from the test resource data
func getServiceHookSubscriptionFromResource(resource *terraform.ResourceState) (*string, error) {
	subscriptionID := resource.Primary.ID

	if subscriptionID == "" {
		return nil, fmt.Errorf("No service hook subscription ID is set")
	}

	return &subscriptionID, nil
}
