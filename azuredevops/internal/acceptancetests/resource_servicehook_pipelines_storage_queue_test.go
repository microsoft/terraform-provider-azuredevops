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

func TestAccServicehookStorageQueue_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	queueName := "testqueue"
	publishedEvent := "RunStateChanged"
	stateFilter := "Completed"
	resultFilter := "Succeeded"
	accountKey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	resourceType := "azuredevops_servicehook_storage_queue_pipelines"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookStorageQueuePipelinesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResource(projectName, accountKey, queueName, stateFilter, resultFilter, publishedEvent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName),
					resource.TestCheckResourceAttr(tfCheckNode, "event_config.0.run_result_filter", resultFilter),
				),
			},
		},
	})
}

func TestAccServicehookStorageQueuePipelines_Update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	publishedEvent := "StageStateChanged"
	queueName1 := "testqueue"
	stateFilter1 := "Completed"
	resultFilter1 := "Succeeded"
	accountKey1 := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	queueName2 := "testqueue"
	stateFilter2 := "Canceling"
	resultFilter2 := "Canceled"
	accountKey2 := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

	resourceType := "azuredevops_servicehook_storage_queue_pipelines"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookStorageQueuePipelinesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResource(projectName, accountKey1, queueName1, stateFilter1, resultFilter1, publishedEvent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName1),
					resource.TestCheckResourceAttr(tfCheckNode, "event_config.0.run_result_filter", resultFilter1),
				),
			},
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResource(projectName, accountKey2, queueName2, stateFilter2, resultFilter2, publishedEvent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName2),
					resource.TestCheckResourceAttr(tfCheckNode, "event_config.0.run_result_filter", resultFilter2),
				),
			},
		},
	})
}

func TestAccServicehookStorageQueuePipelines_InvalidEventConfig(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	queueName := "testqueue"
	publishedEvent := "RunStateChanged"
	stateFilter := "Completed"
	resultFilter := "Succeeded"
	accountKey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookStorageQueuePipelinesDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      testutils.HclServicehookStorageQeueuePipelinesResource(projectName, accountKey, queueName, stateFilter, resultFilter, publishedEvent),
				ExpectError: regexp.MustCompile("Unknown subscription input \"stageResultId\""),
			},
		},
	})
}

func TestAccServicehookStorageQueue_accountKeyError(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookStorageQueuePipelinesDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      testutils.HclServicehookStorageQeueuePipelinesResource(projectName, "accountkey", "testqueue", "Canceled", "Canceled", "RunStateChanged"),
				ExpectError: regexp.MustCompile("expected length of account_key to be in the range \\(64 - 100\\)"),
			},
		},
	})
}

func TestAccServicehookStorageQueuePipelines_NoEventConfig_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	queueName := "testqueue"
	publishedEvent1 := "RunStateChanged"
	publishedEvent2 := "StageStateChanged"
	accountKey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	resourceType := "azuredevops_servicehook_storage_queue_pipelines"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookStorageQueuePipelinesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResourceWithoutEventConfig(projectName, accountKey, queueName, publishedEvent1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName),
					resource.TestCheckResourceAttr(tfCheckNode, "published_event", publishedEvent1),
					resource.TestCheckResourceAttr(tfCheckNode, "event_config.0.pipeline_id", ""),
				),
			},
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResourceWithoutEventConfig(projectName, accountKey, queueName, publishedEvent2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName),
					resource.TestCheckResourceAttr(tfCheckNode, "published_event", publishedEvent1),
					resource.TestCheckResourceAttr(tfCheckNode, "event_config.0.pipeline_id", ""),
				),
			},
		},
	})
}

func TestAccServicehookStorageQueuePipelines_AddEventConfig(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	queueName := "testqueue"
	publishedEvent := "RunStateChanged"
	accountKey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	runStateFilter := "Completed"
	runResultFilter := "Succeeded"

	resourceType := "azuredevops_servicehook_storage_queue_pipelines"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookStorageQueuePipelinesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResourceWithoutEventConfig(projectName, accountKey, queueName, publishedEvent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName),
					resource.TestCheckResourceAttr(tfCheckNode, "published_event", publishedEvent),
					resource.TestCheckResourceAttr(tfCheckNode, "event_config.0.pipeline_id", ""),
				),
			},
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResource(projectName, accountKey, queueName, runStateFilter, runResultFilter, publishedEvent),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName),
					resource.TestCheckResourceAttr(tfCheckNode, "published_event", publishedEvent),
					resource.TestCheckResourceAttr(tfCheckNode, "event_config.0.pipeline_id", ""),
				),
			},
		},
	})
}

func CheckServicehookStorageQueuePipelinesDestroyed(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_servicehook_storage_queue_pipelines" {
			continue
		}

		// indicates the build definition still exists - this should fail the test
		if _, err := getBuildDefinitionFromResource(resource); err == nil {
			return fmt.Errorf("Unexpectedly found a build definition that should be deleted")
		}
	}

	return nil
}
