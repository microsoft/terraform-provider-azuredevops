// +build all resource_servicehook_webhook
// +build !exclude_servicehook

package acceptancetests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceHookWebhook_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	eventType := "git.push"
	url := "https://webhooks.org"

	resourceType := "azuredevops_servicehook_webhook"
	tfSvcHookNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceHookWebhookDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcHookWebhookResourceBasic(projectName, eventType, url),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceHookWebhookExistsWithEventTypeAndUrl(tfSvcHookNode, eventType, url),
					resource.TestCheckResourceAttrSet(tfSvcHookNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcHookNode, "event_type", eventType),
					resource.TestCheckResourceAttr(tfSvcHookNode, "url", url),
				),
			},
		},
	})
}

func TestAccServiceHookWebhook_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repositoryName := testutils.GenerateResourceName()
	eventType := "git.push"
	url := "https://webhooks.org"
	basicAuth := BasicAuth{
		Username: "some_username",
		Password: "some_password",
	}
	httpHeaders := map[string]string{
		"Authorization":     "Bearer bearing",
		"X-My-Secret-Token": "whatever",
	}

	resourceType := "azuredevops_servicehook_webhook"
	tfSvcHookNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceHookWebhookDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcHookWebhookResourceComplete(projectName, repositoryName, eventType, url, basicAuth, httpHeaders),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceHookWebhookExistsWithEventTypeAndUrl(tfSvcHookNode, eventType, url),
					resource.TestCheckResourceAttrSet(tfSvcHookNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcHookNode, "event_type", eventType),
					resource.TestCheckResourceAttr(tfSvcHookNode, "url", url),
				),
			},
		},
	})
}

func TestAccServiceHookWebhook_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	eventType := "git.push"
	url := "https://webhooks.org"

	updatedEventType := "git.pullrequest.created"
	updatedUrl := "https://webhooks2.org"

	resourceType := "azuredevops_servicehook_webhook"
	tfSvcHookNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceHookWebhookDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclSvcHookWebhookResourceBasic(projectName, eventType, url),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceHookWebhookExistsWithEventTypeAndUrl(tfSvcHookNode, eventType, url),
					resource.TestCheckResourceAttrSet(tfSvcHookNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcHookNode, "event_type", eventType),
					resource.TestCheckResourceAttr(tfSvcHookNode, "url", url),
				),
			},
			{
				Config: hclSvcHookWebhookResourceUpdate(projectName, updatedEventType, updatedUrl),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckServiceHookWebhookExistsWithEventTypeAndUrl(tfSvcHookNode, updatedEventType, updatedUrl),
					resource.TestCheckResourceAttrSet(tfSvcHookNode, "project_id"),
					resource.TestCheckResourceAttr(tfSvcHookNode, "event_type", updatedEventType),
					resource.TestCheckResourceAttr(tfSvcHookNode, "url", updatedUrl),
				),
			},
		},
	})
}

func hclSvcHookWebhookResourceBasic(projectName string, eventType string, url string) string {
	serviceHookResource := fmt.Sprintf(`
resource "azuredevops_servicehook_webhook" "test" {
	project_id             = azuredevops_project.project.id
	event_type = "%s"
	url        = "%s"
}`, eventType, url)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceHookResource)
}

type BasicAuth struct {
	Username string
	Password string
}

func hclSvcHookWebhookResourceComplete(projectName string, repositoryName string, eventType string, url string, basicAuth BasicAuth, httpHeaders map[string]string) string {
	headers := []string{}
	for key, val := range httpHeaders {
		headers = append(headers, fmt.Sprintf("%s = \"%s\"", key, val))
	}

	serviceHookResource := fmt.Sprintf(`
resource "azuredevops_servicehook_webhook" "test" {
	project_id             = azuredevops_project.project.id
	event_type = "%s"
	url        = "%s"

	basic_auth {
		username = "%s"
		password = "%s"
	}


	filters = {
		repository = azuredevops_git_repository.repository.id
	}

	http_headers = {
		%s
	}
}`, eventType, url, basicAuth.Username, basicAuth.Password, strings.Join(headers, "\n"))

	projectAndRepositoryResource := testutils.HclGitRepoResource(projectName, repositoryName, "Clean")
	return fmt.Sprintf("%s\n%s", projectAndRepositoryResource, serviceHookResource)
}

func hclSvcHookWebhookResourceUpdate(projectName string, eventType string, url string) string {
	serviceHookResource := fmt.Sprintf(`
resource "azuredevops_servicehook_webhook" "test" {
	project_id             = azuredevops_project.project.id
	event_type = "%s"
	url        = "%s"
}`, eventType, url)

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, serviceHookResource)
}
