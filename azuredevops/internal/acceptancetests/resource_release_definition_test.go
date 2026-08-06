package acceptancetests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/release"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

const releaseCmdLineTaskID = "d9bafed4-0b18-4f58-968d-86655b4d2ce9"

func TestAccReleaseDefinition_basic(t *testing.T) {
	name := testutils.GenerateResourceName()
	tfNode := "azuredevops_release_definition.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkReleaseDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclReleaseDefinitionBasic(name),
				Check: resource.ComposeTestCheckFunc(
					checkReleaseDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "revision"),
					resource.TestCheckResourceAttr(tfNode, "name", name),
					resource.TestCheckResourceAttr(tfNode, "stage.#", "1"),
					resource.TestCheckResourceAttr(tfNode, "stage.0.name", "Stage 1"),
					resource.TestCheckResourceAttr(tfNode, "stage.0.deploy_phase.0.task.0.task_id", releaseCmdLineTaskID),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccReleaseDefinition_update(t *testing.T) {
	name := testutils.GenerateResourceName()
	nameUpdated := testutils.GenerateResourceName()
	tfNode := "azuredevops_release_definition.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkReleaseDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclReleaseDefinitionBasic(name),
				Check: resource.ComposeTestCheckFunc(
					checkReleaseDefinitionExists(name),
					resource.TestCheckResourceAttr(tfNode, "name", name),
					resource.TestCheckResourceAttr(tfNode, "path", `\`),
				),
			},
			{
				Config: hclReleaseDefinitionUpdated(nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					checkReleaseDefinitionExists(nameUpdated),
					resource.TestCheckResourceAttr(tfNode, "name", nameUpdated),
					resource.TestCheckResourceAttr(tfNode, "path", `\ExampleFolder`),
					resource.TestCheckResourceAttr(tfNode, "stage.#", "2"),
				),
			},
		},
	})
}

func TestAccReleaseDefinition_variables(t *testing.T) {
	name := testutils.GenerateResourceName()
	tfNode := "azuredevops_release_definition.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkReleaseDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclReleaseDefinitionVariables(name),
				Check: resource.ComposeTestCheckFunc(
					checkReleaseDefinitionExists(name),
					resource.TestCheckResourceAttr(tfNode, "variable.#", "2"),
				),
			},
		},
	})
}

func TestAccReleaseDefinition_complete(t *testing.T) {
	name := testutils.GenerateResourceName()
	tfNode := "azuredevops_release_definition.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkReleaseDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclReleaseDefinitionComplete(name),
				Check: resource.ComposeTestCheckFunc(
					checkReleaseDefinitionExists(name),
					resource.TestCheckResourceAttr(tfNode, "stage.#", "2"),
					resource.TestCheckResourceAttr(tfNode, "stage.0.execution_policy.0.concurrency_count", "2"),
					resource.TestCheckResourceAttr(tfNode, "stage.1.environment_trigger.0.trigger_type", "rollbackRedeploy"),
					resource.TestCheckResourceAttr(tfNode, "stage.1.pre_deployment_gates.#", "1"),
				),
			},
		},
	})
}

func checkReleaseDefinitionExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources["azuredevops_release_definition.test"]
		if !ok {
			return fmt.Errorf("Did not find a release definition in the TF state")
		}

		releaseDefinition, err := getReleaseDefinitionFromResource(res)
		if err != nil {
			return err
		}

		if *releaseDefinition.Name != expectedName {
			return fmt.Errorf("Release Definition has Name=%s, but expected Name=%s", *releaseDefinition.Name, expectedName)
		}

		return nil
	}
}

func checkReleaseDefinitionDestroyed(s *terraform.State) error {
	for _, res := range s.RootModule().Resources {
		if res.Type != "azuredevops_release_definition" {
			continue
		}

		if _, err := getReleaseDefinitionFromResource(res); err == nil {
			return fmt.Errorf("Unexpectedly found a release definition that should be deleted")
		}
	}

	return nil
}

func getReleaseDefinitionFromResource(res *terraform.ResourceState) (*release.ReleaseDefinition, error) {
	releaseDefID, err := strconv.Atoi(res.Primary.ID)
	if err != nil {
		return nil, err
	}

	projectID := res.Primary.Attributes["project_id"]
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
	return clients.ReleaseClient.GetReleaseDefinition(clients.Ctx, release.GetReleaseDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &releaseDefID,
	})
}

func hclReleaseDefinitionProject(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "%[1]s-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_agent_queue" "test" {
  project_id = azuredevops_project.test.id
  name       = "Azure Pipelines"
}`, name)
}

func hclReleaseDefinitionBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_release_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"

  stage {
    name = "Stage 1"

    deploy_phase {
      name          = "Agent job"
      agent_queue_id = data.azuredevops_agent_queue.test.id

      task {
        task_id = "%[3]s"
        version = "2.*"
        inputs = {
          script = "echo hello"
        }
      }
    }
  }
}`, hclReleaseDefinitionProject(name), name, releaseCmdLineTaskID)
}

func hclReleaseDefinitionUpdated(name string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_release_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  path       = "\\ExampleFolder"

  stage {
    name = "Stage 1"

    deploy_phase {
      name          = "Agent job"
      agent_queue_id = data.azuredevops_agent_queue.test.id

      task {
        task_id = "%[3]s"
        version = "2.*"
        inputs = {
          script = "echo hello"
        }
      }
    }
  }

  stage {
    name = "Stage 2"

    condition {
      condition_type = "environmentState"
      name           = "Stage 1"
      value          = "4"
    }

    deploy_phase {
      name          = "Agent job"
      agent_queue_id = data.azuredevops_agent_queue.test.id

      task {
        task_id = "%[3]s"
        version = "2.*"
        inputs = {
          script = "echo world"
        }
      }
    }
  }
}`, hclReleaseDefinitionProject(name), name, releaseCmdLineTaskID)
}

func hclReleaseDefinitionVariables(name string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_release_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"

  variable {
    name  = "FOO_VAR"
    value = "foo"
  }

  variable {
    name         = "BAR_VAR"
    is_secret    = true
    secret_value = "bar"
  }

  stage {
    name = "Stage 1"

    deploy_phase {
      name          = "Agent job"
      agent_queue_id = data.azuredevops_agent_queue.test.id

      task {
        task_id = "%[3]s"
        version = "2.*"
        inputs = {
          script = "echo $(FOO_VAR)"
        }
      }
    }
  }
}`, hclReleaseDefinitionProject(name), name, releaseCmdLineTaskID)
}

func hclReleaseDefinitionComplete(name string) string {
	return fmt.Sprintf(`
%s

data "azuredevops_group" "test" {
  project_id = azuredevops_project.test.id
  name       = "Contributors"
}

resource "azuredevops_serviceendpoint_generic" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%[2]s-gate"
  server_url            = "https://api.example.com"
}

resource "azuredevops_release_definition" "test" {
  project_id          = azuredevops_project.test.id
  name                = "%[2]s"
  release_name_format = "Release-$(rev:r)"

  variable {
    name  = "app_name"
    value = "myapp"
  }

  stage {
    name = "Dev"

    condition {
      condition_type = "event"
      name           = "ReleaseStarted"
    }

    execution_policy {
      concurrency_count = 2
      queue_depth_count = 0
    }

    variable {
      name  = "target_url"
      value = "https://dev.example.com"
    }

    retention_policy {
      days_to_keep     = 15
      releases_to_keep = 5
      retain_build     = true
    }

    deploy_phase {
      name          = "Agent job"
      agent_queue_id = data.azuredevops_agent_queue.test.id

      task {
        task_id = "%[3]s"
        version = "2.*"
        inputs = {
          script = "echo Deploying $(app_name) to $(target_url)"
        }
      }
    }

    deploy_phase {
      name       = "Agentless job"
      phase_type = "runOnServer"

      task {
        task_id = "9C3E8943-130D-4C78-AC63-8AF81DF62DFB"
        version = "1.*"
        inputs = {
          connectionType       = "connectedServiceName"
          connectedServiceName = azuredevops_serviceendpoint_generic.test.id
          method               = "GET"
        }
      }
    }
  }

  stage {
    name = "Prod"

    condition {
      condition_type = "environmentState"
      name           = "Dev"
      value          = "4"
    }

    environment_trigger {
      trigger_type = "rollbackRedeploy"
      action       = "LatestSuccessfulDeployment"
      event_types  = ["MS.TF.DistributedTask.DeploymentFailed"]
    }

    pre_deployment_gates {
      sampling_interval        = 5
      stabilization_time       = 0
      minimum_success_duration = 0
      timeout                  = 120

      gate {
        task {
          task_id = "9C3E8943-130D-4C78-AC63-8AF81DF62DFB"
          version = "1.*"
          inputs = {
            connectionType       = "connectedServiceName"
            connectedServiceName = azuredevops_serviceendpoint_generic.test.id
            method               = "GET"
          }
        }
      }
    }

    pre_deploy_approval {
      approval {
        approver_id        = data.azuredevops_group.test.group_id
        is_automated       = false
        is_notification_on = true
      }

      approval_options {
        required_approver_count         = 1
        release_creator_can_be_approver = false
        enforce_identity_revalidation   = false
        timeout_in_minutes              = 1440
        execution_order                 = "beforeGates"
      }
    }

    deploy_phase {
      name          = "Agent job"
      agent_queue_id = data.azuredevops_agent_queue.test.id

      task {
        task_id = "%[3]s"
        version = "2.*"
        inputs = {
          script = "echo Deploying $(app_name) to $(target_url)"
        }
      }
    }
  }
}`, hclReleaseDefinitionProject(name), name, releaseCmdLineTaskID)
}
