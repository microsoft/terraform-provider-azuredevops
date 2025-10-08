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

// HclServicehookWebhookTfsResourceWithGitPushEvent creates a HCL representation of a basic TFS webhook for git push events
func HclServicehookWebhookTfsResourceWithGitPushEvent(projectName, url string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_servicehook_webhook_tfs" "test" {
  project_id = azuredevops_project.project.id
  url        = "%s"
  
  git_push {
    branch = "main"
  }
}
`, projectResource, url)
}

// HclServicehookWebhookTfsResourceWithWorkItemCreatedEvent creates a HCL representation of a TFS webhook for work item creation events
func HclServicehookWebhookTfsResourceWithWorkItemCreatedEvent(projectName, url string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_servicehook_webhook_tfs" "test" {
  project_id = azuredevops_project.project.id
  url        = "%s"
  
  work_item_created {
    work_item_type = "Bug"
    area_path      = "\\%s"
  }
}
`, projectResource, url, projectName)
}

// HclServicehookWebhookTfsResourceWithHeadersAndAuth creates a HCL representation of a TFS webhook with custom headers and basic auth
func HclServicehookWebhookTfsResourceWithHeadersAndAuth(projectName, url string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_servicehook_webhook_tfs" "test" {
  project_id = azuredevops_project.project.id
  url        = "%s"
  
  git_push {}
  
  http_headers = {
    "X-Custom-Header" = "Test Value"
    "Content-Type"    = "application/json"
  }
  
  basic_auth_username = "testuser"
  basic_auth_password = "testpassword"
}
`, projectResource, url)
}

// HclServicehookWebhookTfsResourceWithResourceDetails creates a HCL representation of a TFS webhook with custom resource details settings
func HclServicehookWebhookTfsResourceWithResourceDetails(projectName, url string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_servicehook_webhook_tfs" "test" {
  project_id = azuredevops_project.project.id
  url        = "%s"
  
  git_push {}
  
  resource_details_to_send  = "minimal"
  messages_to_send          = "text"
  detailed_messages_to_send = "html"
}
`, projectResource, url)
}

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
				Config: HclServicehookWebhookTfsResourceWithGitPushEvent(projectName, url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfCheckNode, "url", url),
					resource.TestCheckResourceAttr(tfCheckNode, "git_push.#", "1"),
				),
			},
			{
				ResourceName:      tfCheckNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  checkImportProject(),
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
				Config: HclServicehookWebhookTfsResourceWithGitPushEvent(projectName, url1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfCheckNode, "url", url1),
					resource.TestCheckResourceAttr(tfCheckNode, "git_push.#", "1"),
				),
			},
			{
				Config: HclServicehookWebhookTfsResourceWithGitPushEvent(projectName, url2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfCheckNode, "url", url2),
					resource.TestCheckResourceAttr(tfCheckNode, "git_push.#", "1"),
				),
			},
			{
				ResourceName:      tfCheckNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  checkImportProject(),
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
				Config: HclServicehookWebhookTfsResourceWithGitPushEvent(projectName, url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfCheckNode, "git_push.#", "1"),
				),
			},
			{
				Config: HclServicehookWebhookTfsResourceWithWorkItemCreatedEvent(projectName, url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfCheckNode, "work_item_created.#", "1"),
					resource.TestCheckNoResourceAttr(tfCheckNode, "git_push.#"),
				),
			},
			{
				ResourceName:      tfCheckNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  checkImportProject(),
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
				Config: HclServicehookWebhookTfsResourceWithHeadersAndAuth(projectName, url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfCheckNode, "url", url),
					resource.TestCheckResourceAttr(tfCheckNode, "http_headers.%", "2"),
					resource.TestCheckResourceAttr(tfCheckNode, "http_headers.X-Custom-Header", "Test Value"),
					resource.TestCheckResourceAttr(tfCheckNode, "basic_auth_username", "testuser"),
				),
			},
			{
				ResourceName:      tfCheckNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  checkImportProject(),
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
				Config:      HclServicehookWebhookTfsResourceWithGitPushEvent(projectName, invalidUrl),
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
				Config: HclServicehookWebhookTfsResourceWithResourceDetails(projectName, url),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfCheckNode, "url", url),
					resource.TestCheckResourceAttr(tfCheckNode, "resource_details_to_send", "minimal"),
					resource.TestCheckResourceAttr(tfCheckNode, "messages_to_send", "text"),
					resource.TestCheckResourceAttr(tfCheckNode, "detailed_messages_to_send", "html"),
				),
			},
			{
				ResourceName:      tfCheckNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  checkImportProject(),
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
