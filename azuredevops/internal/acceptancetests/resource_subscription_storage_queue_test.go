//go:build (all || resource_check_required_template) && !exclude_approvalsandchecks
// +build all resource_check_required_template
// +build !exclude_approvalsandchecks

package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccSubscriptionStorageQueue_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	queueName := "testqueue"
	resultFilter := "Succeeded"
	accountKey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	resourceType := "azuredevops_subscription_storage_queue"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckSubscriptionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclSubscriptionStorageQeueueResourceWithPipelinesPublisher(projectName, accountKey, queueName, resultFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName),
					resource.TestCheckResourceAttr(tfCheckNode, "publisher.0.stage_state_changed.0.result_filter", resultFilter),
				),
			},
		},
	})
}

func TestAccSubscriptionStorageQueue_accountKeyError(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckSubscriptionDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      testutils.HclSubscriptionStorageQeueueResourceWithPipelinesPublisher(projectName, "accountkey", "testqueue", "Canceled"),
				ExpectError: regexp.MustCompile("expected length of account_key to be in the range \\(64 - 100\\)"),
			},
		},
	})
}

func CheckSubscriptionDestroyed(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_subscription_storage_queue" {
			continue
		}

		// indicates the build definition still exists - this should fail the test
		if _, err := getBuildDefinitionFromResource(resource); err == nil {
			return fmt.Errorf("Unexpectedly found a build definition that should be deleted")
		}
	}

	return nil
}
