//go:build (all || resource_build_definition) && !exclude_resource_build_definition

package acceptancetests

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func TestAccBuildDefinition_basic(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfBuildDefNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionPath(name, `\\`),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", `\`),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
		},
	})
}

func TestAccBuildDefinition_pathUpdate(t *testing.T) {
	name := testutils.GenerateResourceName()

	pathFirst := `\\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	pathSecond := `\\` + name + `\\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfBuildDefNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionPath(name, pathFirst),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", strings.ReplaceAll(pathFirst, `\\`, `\`)),
				),
			},
			{
				Config: hclBuildDefinitionPath(name, pathSecond),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", strings.ReplaceAll(pathSecond, `\\`, `\`)),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
		},
	})
}

// Verifies a build for with variables can create and update, including secret variables
func TestAccBuildDefinition_withVariables(t *testing.T) {
	name := testutils.GenerateResourceName()
	tfNode := "azuredevops_build_definition.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionVariable(name, "foo1", "bar1"),
				Check:  checkForVariableValues(tfNode, "foo1", "bar1"),
			}, {
				Config: hclBuildDefinitionVariable(name, "foo2", "bar2"),
				Check:  checkForVariableValues(tfNode, "foo2", "bar2"),
			},
		},
	})
}

func TestAccBuildDefinition_schedules(t *testing.T) {
	name := testutils.GenerateResourceName()
	tfNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionSchedules(name),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "revision"),
					resource.TestCheckResourceAttrSet(tfNode, "repository.0.repo_id"),
					resource.TestCheckResourceAttr(tfNode, "schedules.#", "1"),
					resource.TestCheckResourceAttr(tfNode, "schedules.0.days_to_build.#", "1"),
					resource.TestCheckResourceAttr(tfNode, "name", name),
				),
			},
		},
	})
}

func TestAccBuildDefinition_buildCompletionTrigger(t *testing.T) {
	name := testutils.GenerateResourceName()
	tfNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionBuildCompletionTrigger(name),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "repository.0.repo_id"),
					resource.TestCheckResourceAttr(tfNode, "build_completion_trigger.#", "1"),
					resource.TestCheckResourceAttr(tfNode, "name", name),
				),
			},
		},
	})
}

func TestAccBuildDefinition_otherGitRepositoryAgentJob_basic(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfBuildDefNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionOtherGitRepositoryAgentJobBasic(name, `\\`),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", `\`),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.ref_name", "agent_job1"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "jobs.0.target.#"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.type", "AgentJob"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.#", "1"),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
		},
	})
}

func TestAccBuildDefinition_otherGitRepositoryAgentJob_multiConfiguration(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfBuildDefNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionOtherGitRepositoryAgentJobMultiConfiguration(name, `\\`),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", `\`),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.ref_name", "agent_job1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.condition", "succeededOrFailed()"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "jobs.0.target.#"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.type", "AgentJob"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.demands.#", "2"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.type", "Multi-Configuration"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.multipliers", "multipliers1,multipliers2"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.max_concurrency", "3"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.continue_on_error", "true"),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
		},
	})
}

func TestAccBuildDefinition_otherGitRepositoryAgentJob_multiAgent(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfBuildDefNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionOtherGitRepositoryAgentJobMultiAgent(name, `\\`),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", `\`),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.ref_name", "agent_job3"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.condition", "succeededOrFailed()"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.type", "AgentJob"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.demands.#", "2"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.type", "Multi-Agent"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.max_concurrency", "2"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.continue_on_error", "false"),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
		},
	})
}

func TestAccBuildDefinition_otherGitRepositoryAgentJob_update(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfBuildDefNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionOtherGitRepositoryAgentJobBasic(name, `\\`),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", `\`),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.ref_name", "agent_job1"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "jobs.0.target.#"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.type", "AgentJob"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.#", "1"),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
			{
				Config: hclBuildDefinitionOtherGitRepositoryAgentJobMultiConfiguration(name, `\\`),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", `\`),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.ref_name", "agent_job1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.condition", "succeededOrFailed()"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "jobs.0.target.#"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.type", "AgentJob"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.demands.#", "2"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.type", "Multi-Configuration"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.multipliers", "multipliers1,multipliers2"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.max_concurrency", "3"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.continue_on_error", "true"),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
			{
				Config: hclBuildDefinitionOtherGitRepositoryAgentJobMultiAgent(name, `\\`),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", `\`),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.ref_name", "agent_job3"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.condition", "succeededOrFailed()"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.type", "AgentJob"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.demands.#", "2"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.type", "Multi-Agent"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.max_concurrency", "2"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.continue_on_error", "false"),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
		},
	})
}

func TestAccBuildDefinition_otherGitRepositoryAgentJob_complete(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfBuildDefNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionOtherGitRepositoryAgentJobComplete(name, `\\`),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", `\`),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.#", "3"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.ref_name", "agent_job2"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.dependencies.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.job_timeout_in_minutes", "60"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.job_cancel_timeout_in_minutes", "5"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.target.0.execution_options.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.target.0.execution_options.0.type", "Multi-Configuration"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.target.0.execution_options.0.max_concurrency", "3"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.target.0.demands.#", "2"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.2.dependencies.#", "2"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.2.ref_name", "agent_job3"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.2.condition", "succeededOrFailed()"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.2.job_timeout_in_minutes", "60"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.2.job_cancel_timeout_in_minutes", "5"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.2.target.0.type", "AgentJob"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.2.target.0.demands.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.2.target.0.execution_options.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.2.target.0.execution_options.0.type", "Multi-Agent"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.2.target.0.execution_options.0.max_concurrency", "2"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.2.target.0.execution_options.0.continue_on_error", "false"),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
		},
	})
}

func TestAccBuildDefinition_otherGitRepositoryAgentlessJob_basic(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfBuildDefNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionOtherGitRepositoryAgentlessJobBasic(name, `\\`),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", `\`),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.ref_name", "agentless_job1"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "jobs.0.target.#"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.type", "AgentlessJob"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.type", "None"),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
		},
	})
}

func TestAccBuildDefinition_otherGitRepositoryAgentlessJob_multiConfiguration(t *testing.T) {
	name := testutils.GenerateResourceName()
	tfBuildDefNode := "azuredevops_build_definition.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionOtherGitRepositoryAgentlessJobMultiConfiguration(name, `\\`),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", `\`),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.ref_name", "agentless_job1"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "jobs.0.target.#"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.type", "AgentlessJob"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.type", "Multi-Configuration"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.multipliers", "multiplierstest"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.continue_on_error", "true"),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
		},
	})

}

func TestAccBuildDefinition_otherGitRepositoryAgentlessJob_update(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfBuildDefNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionOtherGitRepositoryAgentlessJobBasic(name, `\\`),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", `\`),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.ref_name", "agentless_job1"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "jobs.0.target.#"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.type", "AgentlessJob"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.type", "None"),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
			{
				Config: hclBuildDefinitionOtherGitRepositoryAgentlessJobMultiConfiguration(name, `\\`),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", `\`),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.ref_name", "agentless_job1"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "jobs.0.target.#"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.type", "AgentlessJob"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.type", "Multi-Configuration"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.multipliers", "multiplierstest"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.continue_on_error", "true"),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
		},
	})
}

func TestAccBuildDefinition_otherGitRepositoryAgentlessJob_complete(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfBuildDefNode := "azuredevops_build_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclBuildDefinitionOtherGitRepositoryAgentlessJobComplete(name, `\\`),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(name),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", name),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", `\`),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.#", "2"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.ref_name", "agentless_job1"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "jobs.0.target.#"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.type", "AgentlessJob"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.job_timeout_in_minutes", "60"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.job_cancel_timeout_in_minutes", "5"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.0.target.0.execution_options.0.type", "None"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.ref_name", "agentless_job2"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.job_timeout_in_minutes", "60"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.job_cancel_timeout_in_minutes", "5"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.target.0.execution_options.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.target.0.execution_options.0.type", "Multi-Configuration"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.dependencies.#", "1"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.target.0.execution_options.0.multipliers", "multiplierstest"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "jobs.1.target.0.execution_options.0.continue_on_error", "true"),
				),
			}, {
				ResourceName:            tfBuildDefNode,
				ImportStateIdFunc:       testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_first_run"},
			},
		},
	})
}

// Checks that the expected variable values exist in the state
func checkForVariableValues(tfNode string, expectedVals ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rootModule := s.RootModule()
		resource, ok := rootModule.Resources[tfNode]
		if !ok {
			return fmt.Errorf("Did not find resource in TF state")
		}

		is := resource.Primary
		if is == nil {
			return fmt.Errorf("No primary instance: %s in %s", tfNode, rootModule.Path)
		}

		for _, expectedVal := range expectedVals {
			found := false
			for _, value := range is.Attributes {
				if value == expectedVal {
					found = true
				}
			}

			if !found {
				return fmt.Errorf("Did not find variable with value %s", expectedVal)
			}

		}

		return nil
	}
}

func checkBuildDefinitionExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		buildDef, ok := s.RootModule().Resources["azuredevops_build_definition.test"]
		if !ok {
			return fmt.Errorf("Did not find a build definition in the TF state")
		}

		buildDefinition, err := getBuildDefinitionFromResource(buildDef)
		if err != nil {
			return err
		}

		if *buildDefinition.Name != expectedName {
			return fmt.Errorf("Build Definition has Name=%s, but expected Name=%s", *buildDefinition.Name, expectedName)
		}

		return nil
	}
}

// verifies that all build definitions referenced in the state are destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func checkBuildDefinitionDestroyed(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_build_definition" {
			continue
		}

		// indicates the build definition still exists - this should fail the test
		if _, err := getBuildDefinitionFromResource(resource); err == nil {
			return fmt.Errorf("Unexpectedly found a build definition that should be deleted")
		}
	}

	return nil
}

// given a resource from the state, return a build definition (and error)
func getBuildDefinitionFromResource(resource *terraform.ResourceState) (*build.BuildDefinition, error) {
	buildDefID, err := strconv.Atoi(resource.Primary.ID)
	if err != nil {
		return nil, err
	}

	projectID := resource.Primary.Attributes["project_id"]
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
	return clients.BuildClient.GetDefinition(clients.Ctx, build.GetDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &buildDefID,
	})
}

func hclBuildDefinitionTemplate(name string) string {
	return fmt.Sprintf(`

resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "%[1]s-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "acc-%[1]s"
  initialization {
    init_type = "Clean"
  }
}`, name)
}

func hclBuildDefinitionPath(name, path string) string {
	template := hclBuildDefinitionTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  path       = "%[3]s"

  ci_trigger {
    override {
      batch = true
      branch_filter {
        include = ["master"]
      }
      path_filter {
        include = ["*/**.ts"]
      }
      max_concurrent_builds_per_branch = 2
      polling_interval                 = 0
    }
  }

  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.test.id
    branch_name = azuredevops_git_repository.test.default_branch
    yml_path    = "azure-pipelines.yml"
  }
}
`, template, name, path)
}

func hclBuildDefinitionVariable(name, varVal, secretVarVal string) string {
	template := hclBuildDefinitionTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.test.id
    branch_name = azuredevops_git_repository.test.default_branch
    yml_path    = "azure-pipelines.yml"
  }

  variable {
    name  = "FOO_VAR"
    value = "%[3]s"
  }

  variable {
    name         = "BAR_VAR"
    secret_value = "%[4]s"
    is_secret    = true
  }
}`, template, name, varVal, secretVarVal)
}

func hclBuildDefinitionSchedules(name string) string {
	template := hclBuildDefinitionTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  path       = "\\ExampleFolder"

  ci_trigger {
    override {
      batch = true
      branch_filter {
        include = ["master"]
      }
      path_filter {
        include = ["*/**.ts"]
      }
      max_concurrent_builds_per_branch = 2
      polling_interval                 = 0
    }
  }

  schedules {
    branch_filter {
      include = ["master"]
    }

    days_to_build              = ["Mon"]
    schedule_only_with_changes = true
    start_hours                = 0
    start_minutes              = 0
    time_zone                  = "(UTC) Coordinated Universal Time"
  }

  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.test.id
    branch_name = azuredevops_git_repository.test.default_branch
    yml_path    = "azure-pipelines.yml"
  }
}
`, template, name)
}

func hclBuildDefinitionBuildCompletionTrigger(name string) string {
	template := hclBuildDefinitionTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_build_definition" "build_trigger" {
  project_id = azuredevops_project.test.id
  name       = "trigger%[2]s"
  path       = "\\ExampleFolder"


  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.test.id
    branch_name = azuredevops_git_repository.test.default_branch
    yml_path    = "azure-pipelines.yml"
  }
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  path       = "\\ExampleFolder"

  build_completion_trigger {
    build_definition_id = azuredevops_build_definition.build_trigger.id
    branch_filter {
      include = ["main"]
      exclude = ["test", "regression"]
    }
  }
  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.test.id
    branch_name = azuredevops_git_repository.test.default_branch
    yml_path    = "azure-pipelines.yml"
  }
}
`, template, name)
}

func hclBuildDefinitionOtherGitRepositoryAgentJobBasic(name, path string) string {
	template := hclBuildDefinitionTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_git" "test" {
  project_id            = azuredevops_project.test.id
  repository_url        = "https://dev.azure.com/org/project/_git/test"
  username              = "username"
  password              = "password"
  service_endpoint_name = "Generic Git"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  path       = "%[3]s"

  ci_trigger {
    override {
      batch = true
      branch_filter {
        include = ["master"]
      }
      path_filter {
        include = ["*/**.ts"]
      }
      max_concurrent_builds_per_branch = 2
      polling_interval                 = 0
    }
  }

  agent_specification = "windows-latest"

  repository {
    repo_type             = "Git"
    repo_id               = azuredevops_serviceendpoint_generic_git.test.repository_url
    branch_name           = "refs/heads/main"
    url                   = azuredevops_serviceendpoint_generic_git.test.repository_url
    service_connection_id = azuredevops_serviceendpoint_generic_git.test.id
  }

  jobs {
    name      = "Agent Job1"
    ref_name  = "agent_job1"
    condition = "succeeded()"
    target {
      type = "AgentJob"
      execution_options {
        type = "None"
      }
    }
  }
}
`, template, name, path)
}

func hclBuildDefinitionOtherGitRepositoryAgentJobMultiConfiguration(name, path string) string {
	template := hclBuildDefinitionTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_git" "test" {
  project_id            = azuredevops_project.test.id
  repository_url        = "https://dev.azure.com/org/project/_git/test"
  username              = "username"
  password              = "password"
  service_endpoint_name = "Generic Git"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  path       = "%[3]s"

  ci_trigger {
    override {
      batch = true
      branch_filter {
        include = ["master"]
      }
      path_filter {
        include = ["*/**.ts"]
      }
      max_concurrent_builds_per_branch = 2
      polling_interval                 = 0
    }
  }

  repository {
    repo_type             = "Git"
    repo_id               = azuredevops_serviceendpoint_generic_git.test.repository_url
    branch_name           = "refs/heads/main"
    url                   = azuredevops_serviceendpoint_generic_git.test.repository_url
    service_connection_id = azuredevops_serviceendpoint_generic_git.test.id
  }

  agent_specification = "windows-latest"

  jobs {
    name      = "Agent Job1"
    ref_name  = "agent_job1"
    condition = "succeededOrFailed()"
    target {
      type    = "AgentJob"
      demands = ["git -equals git", "git"]
      execution_options {
        type              = "Multi-Configuration"
        continue_on_error = true
        multipliers       = "multipliers1,multipliers2"
        max_concurrency   = 3
      }
    }
  }
}
`, template, name, path)
}

func hclBuildDefinitionOtherGitRepositoryAgentJobMultiAgent(name, path string) string {
	template := hclBuildDefinitionTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_git" "test" {
  project_id            = azuredevops_project.test.id
  repository_url        = "https://dev.azure.com/org/project/_git/test"
  username              = "username"
  password              = "password"
  service_endpoint_name = "Generic Git"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  path       = "%[3]s"

  ci_trigger {
    override {
      batch = true
      branch_filter {
        include = ["master"]
      }
      path_filter {
        include = ["*/**.ts"]
      }
      max_concurrent_builds_per_branch = 2
      polling_interval                 = 0
    }
  }

  repository {
    repo_type             = "Git"
    repo_id               = azuredevops_serviceendpoint_generic_git.test.repository_url
    branch_name           = "refs/heads/main"
    url                   = azuredevops_serviceendpoint_generic_git.test.repository_url
    service_connection_id = azuredevops_serviceendpoint_generic_git.test.id
  }

  agent_specification = "windows-latest"

  jobs {
    name      = "Agent Job3"
    ref_name  = "agent_job3"
    condition = "succeededOrFailed()"
    target {
      type    = "AgentJob"
      demands = ["git -equals git", "git"]
      execution_options {
        type              = "Multi-Agent"
        continue_on_error = false
        max_concurrency   = 2
      }
    }
  }
}
`, template, name, path)
}

func hclBuildDefinitionOtherGitRepositoryAgentJobComplete(name, path string) string {
	template := hclBuildDefinitionTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_git" "test" {
  project_id            = azuredevops_project.test.id
  repository_url        = "https://dev.azure.com/org/project/_git/test"
  username              = "username"
  password              = "password"
  service_endpoint_name = "Generic Git"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  path       = "%[3]s"

  ci_trigger {
    override {
      batch = true
      branch_filter {
        include = ["master"]
      }
      path_filter {
        include = ["*/**.ts"]
      }
      max_concurrent_builds_per_branch = 2
      polling_interval                 = 0
    }
  }

  repository {
    repo_type             = "Git"
    repo_id               = azuredevops_serviceendpoint_generic_git.test.repository_url
    branch_name           = "refs/heads/main"
    url                   = azuredevops_serviceendpoint_generic_git.test.repository_url
    service_connection_id = azuredevops_serviceendpoint_generic_git.test.id
  }

  agent_specification = "windows-latest"

  jobs {
    name      = "Agent Job1"
    ref_name  = "agent_job1"
    condition = "succeeded()"
    target {
      type = "AgentJob"
      execution_options {
        type = "None"
      }
    }
  }

  jobs {
    name                          = "Agent Job2"
    ref_name                      = "agent_job2"
    condition                     = "succeededOrFailed()"
    job_timeout_in_minutes        = 60
    job_cancel_timeout_in_minutes = 5
    dependencies {
      scope = "agent_job1"
    }
    target {
      type    = "AgentJob"
      demands = ["git -equals git", "git"]
      execution_options {
        type              = "Multi-Configuration"
        continue_on_error = true
        multipliers       = "multipliers1,multipliers2"
        max_concurrency   = 3
      }
    }
  }

  jobs {
    name      = "Agent Job3"
    ref_name  = "agent_job3"
    condition = "succeededOrFailed()"
    dependencies {
      scope = "agent_job1"
    }
    dependencies {
      scope = "agent_job2"
    }
    job_timeout_in_minutes        = 60
    job_cancel_timeout_in_minutes = 5
    target {
      type    = "AgentJob"
      demands = ["git -equals git"]
      execution_options {
        type              = "Multi-Agent"
        continue_on_error = false
        max_concurrency   = 2
      }
    }
  }
}
`, template, name, path)
}

func hclBuildDefinitionOtherGitRepositoryAgentlessJobBasic(name, path string) string {
	template := hclBuildDefinitionTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_git" "test" {
  project_id            = azuredevops_project.test.id
  repository_url        = "https://dev.azure.com/org/project/_git/test"
  username              = "username"
  password              = "password"
  service_endpoint_name = "Generic Git"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  path       = "%[3]s"

  ci_trigger {
    override {
      batch = true
      branch_filter {
        include = ["master"]
      }
      path_filter {
        include = ["*/**.ts"]
      }
      max_concurrent_builds_per_branch = 2
      polling_interval                 = 0
    }
  }

  repository {
    repo_type             = "Git"
    repo_id               = azuredevops_serviceendpoint_generic_git.test.repository_url
    branch_name           = "refs/heads/main"
    url                   = azuredevops_serviceendpoint_generic_git.test.repository_url
    service_connection_id = azuredevops_serviceendpoint_generic_git.test.id
  }

  agent_specification = "windows-latest"

  jobs {
    name      = "Agentless Job1"
    ref_name  = "agentless_job1"
    condition = "succeeded()"
    target {
      type = "AgentlessJob"
      execution_options {
        type = "None"
      }
    }
  }
}
`, template, name, path)
}

func hclBuildDefinitionOtherGitRepositoryAgentlessJobMultiConfiguration(name, path string) string {
	template := hclBuildDefinitionTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_git" "test" {
  project_id            = azuredevops_project.test.id
  repository_url        = "https://dev.azure.com/org/project/_git/test"
  username              = "username"
  password              = "password"
  service_endpoint_name = "Generic Git"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  path       = "%[3]s"

  ci_trigger {
    override {
      batch = true
      branch_filter {
        include = ["master"]
      }
      path_filter {
        include = ["*/**.ts"]
      }
      max_concurrent_builds_per_branch = 2
      polling_interval                 = 0
    }
  }

  repository {
    repo_type             = "Git"
    repo_id               = azuredevops_serviceendpoint_generic_git.test.repository_url
    branch_name           = "refs/heads/main"
    url                   = azuredevops_serviceendpoint_generic_git.test.repository_url
    service_connection_id = azuredevops_serviceendpoint_generic_git.test.id
  }

  agent_specification = "windows-latest"

  jobs {
    name      = "Agentless Job1"
    ref_name  = "agentless_job1"
    condition = "succeeded()"
    target {
      type = "AgentlessJob"
      execution_options {
        type              = "Multi-Configuration"
        continue_on_error = true
        multipliers       = "multiplierstest"
      }
    }
  }
}
`, template, name, path)
}

func hclBuildDefinitionOtherGitRepositoryAgentlessJobComplete(name, path string) string {
	template := hclBuildDefinitionTemplate(name)
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_git" "test" {
  project_id            = azuredevops_project.test.id
  repository_url        = "https://dev.azure.com/org/project/_git/test"
  username              = "username"
  password              = "password"
  service_endpoint_name = "Generic Git"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  path       = "%[3]s"

  ci_trigger {
    override {
      batch = true
      branch_filter {
        include = ["master"]
      }
      path_filter {
        include = ["*/**.ts"]
      }
      max_concurrent_builds_per_branch = 2
      polling_interval                 = 0
    }
  }

  repository {
    repo_type             = "Git"
    repo_id               = azuredevops_serviceendpoint_generic_git.test.repository_url
    branch_name           = "refs/heads/main"
    url                   = azuredevops_serviceendpoint_generic_git.test.repository_url
    service_connection_id = azuredevops_serviceendpoint_generic_git.test.id
  }

  agent_specification = "windows-latest"

  jobs {
    name                          = "Agentless Job1"
    ref_name                      = "agentless_job1"
    condition                     = "succeeded()"
    job_timeout_in_minutes        = 60
    job_cancel_timeout_in_minutes = 5
    target {
      type = "AgentlessJob"
      execution_options {
        type = "None"
      }
    }
  }

  jobs {
    name                          = "Agentless Job2"
    ref_name                      = "agentless_job2"
    condition                     = "succeeded()"
    job_timeout_in_minutes        = 60
    job_cancel_timeout_in_minutes = 5
    dependencies {
      scope = "agentless_job1"
    }
    target {
      type = "AgentlessJob"
      execution_options {
        type              = "Multi-Configuration"
        continue_on_error = true
        multipliers       = "multiplierstest"
      }
    }
  }
}
`, template, name, path)
}
