package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServicehookStorageQueuePipelines_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	queueName := "testqueue"
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
				Config: testutils.HclServicehookStorageQeueuePipelinesResourceWithStageEvent(projectName, accountKey, queueName, stateFilter, resultFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName),
					resource.TestCheckResourceAttr(tfCheckNode, "account_key", accountKey),
					resource.TestCheckResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_result_filter", resultFilter),
					resource.TestCheckResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_state_filter", stateFilter),
				),
			},
		},
	})
}

func TestAccServicehookStorageQueuePipelines_Update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	queueName1 := "testqueue"
	stateFilter1 := "Completed"
	resultFilter1 := "Succeeded"
	accountKey1 := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	queueName2 := "testqueue"
	stateFilter2 := "Completed"
	resultFilter2 := "Failed"
	accountKey2 := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

	resourceType := "azuredevops_servicehook_storage_queue_pipelines"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookStorageQueuePipelinesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResourceWithStageEvent(projectName, accountKey1, queueName1, stateFilter1, resultFilter1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName1),
					resource.TestCheckResourceAttr(tfCheckNode, "account_key", accountKey1),
					resource.TestCheckResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_result_filter", resultFilter1),
					resource.TestCheckResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_state_filter", stateFilter1),
				),
			},
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResourceWithStageEvent(projectName, accountKey2, queueName2, stateFilter2, resultFilter2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName2),
					resource.TestCheckResourceAttr(tfCheckNode, "account_key", accountKey2),
					resource.TestCheckResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_result_filter", resultFilter2),
					resource.TestCheckResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_state_filter", stateFilter2),
				),
			},
		},
	})
}

func TestAccServicehookStorageQueuePipelines_accountKeyError(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookStorageQueuePipelinesDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      testutils.HclServicehookStorageQeueuePipelinesResourceWithStageEvent(projectName, "accountkey", "testqueue", "Canceled", "Canceled"),
				ExpectError: regexp.MustCompile(`expected length of account_key to be in the range (64 - 100)`),
			},
		},
	})
}

func TestAccServicehookStorageQueuePipelines_NoEventConfig_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	queueName := "testqueue"
	eventType1 := "run_state_changed_event"
	eventType2 := "stage_state_changed_event"
	accountKey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	resourceType := "azuredevops_servicehook_storage_queue_pipelines"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookStorageQueuePipelinesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResourceWithoutEventConfig(projectName, accountKey, queueName, eventType1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
				),
			},
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResourceWithoutEventConfig(projectName, accountKey, queueName, eventType2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
				),
			},
		},
	})
}

func TestAccServicehookStorageQueuePipelines_AddEventConfig(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	queueName := "testqueue"
	eventType := "stage_state_changed_event"
	accountKey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	stateFilter := "Completed"
	resultFilter := "Succeeded"

	resourceType := "azuredevops_servicehook_storage_queue_pipelines"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookStorageQueuePipelinesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResourceWithoutEventConfig(projectName, accountKey, queueName, eventType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName),
					resource.TestCheckNoResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_result_filter"),
				),
			},
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResourceWithStageEvent(projectName, accountKey, queueName, stateFilter, resultFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName),
					resource.TestCheckResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_result_filter", resultFilter),
					resource.TestCheckResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_state_filter", stateFilter),
				),
			},
		},
	})
}

func TestAccServicehookStorageQueuePipelines_RemoveEventConfigAndChangeEvent(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	queueName := "testqueue"
	eventType2 := "run_state_changed_event"
	accountKey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	stateFilter := "Completed"
	resultFilter := "Succeeded"

	resourceType := "azuredevops_servicehook_storage_queue_pipelines"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookStorageQueuePipelinesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResourceWithStageEvent(projectName, accountKey, queueName, stateFilter, resultFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName),
					resource.TestCheckResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_result_filter", resultFilter),
					resource.TestCheckResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_state_filter", stateFilter),
				),
			},
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResourceWithoutEventConfig(projectName, accountKey, queueName, eventType2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName),
					resource.TestCheckNoResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_result_filter"),
					resource.TestCheckNoResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_state_filter"),
				),
			},
		},
	})
}

func TestAccServicehookStorageQueuePipelines_RemoveEventConfig(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	queueName := "testqueue"
	eventType := "stage_state_changed_event"
	accountKey := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	stateFilter := "Completed"
	resultFilter := "Succeeded"

	resourceType := "azuredevops_servicehook_storage_queue_pipelines"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: CheckServicehookStorageQueuePipelinesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResourceWithStageEvent(projectName, accountKey, queueName, stateFilter, resultFilter),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName),
					resource.TestCheckResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_result_filter", resultFilter),
					resource.TestCheckResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_state_filter", stateFilter),
				),
			},
			{
				Config: testutils.HclServicehookStorageQeueuePipelinesResourceWithoutEventConfig(projectName, accountKey, queueName, eventType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "queue_name", queueName),
					resource.TestCheckNoResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_result_filter"),
					resource.TestCheckNoResourceAttr(tfCheckNode, "stage_state_changed_event.0.stage_state_filter"),
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
