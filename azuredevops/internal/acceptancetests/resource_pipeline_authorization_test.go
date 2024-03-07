//go:build (all || pipeline_authorization) && !exclude_pipeline_authorization
// +build all pipeline_authorization
// +build !exclude_pipeline_authorization

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccPipelineAuthorization_allPipeline_queue(t *testing.T) {
	node := "azuredevops_pipeline_authorization.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclAllPipelineAuthQueue(testutils.GenerateResourceName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "project_id"),
					resource.TestCheckResourceAttrSet(node, "resource_id"),
				),
			},
		},
	})
}

func TestAccPipelineAuthorization_pipeline_queue(t *testing.T) {
	node := "azuredevops_pipeline_authorization.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclPipelineAuthQueue(testutils.GenerateResourceName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "project_id"),
					resource.TestCheckResourceAttrSet(node, "resource_id"),
				),
			},
		},
	})
}

func TestAccPipelineAuthorization_multiPipeline_queue(t *testing.T) {
	node := "azuredevops_pipeline_authorization.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclMultiPipelineAuthQueue(testutils.GenerateResourceName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "project_id"),
					resource.TestCheckResourceAttrSet(node, "resource_id"),
				),
			},
		},
	})
}

func TestAccPipelineAuthorization_allPipelineWithPipeline_queue(t *testing.T) {
	node := "azuredevops_pipeline_authorization.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclAllPipelineWithPipoelineAuthQueue(testutils.GenerateResourceName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "project_id"),
					resource.TestCheckResourceAttrSet(node, "resource_id"),
				),
			},
		},
	})
}

func TestAccPipelineAuthorization_allPipeline_environment(t *testing.T) {
	node := "azuredevops_pipeline_authorization.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclAllPipelineAuthEnvironment(testutils.GenerateResourceName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "project_id"),
					resource.TestCheckResourceAttrSet(node, "resource_id"),
				),
			},
		},
	})
}

func TestAccPipelineAuthorization_pipeline_environment(t *testing.T) {
	node := "azuredevops_pipeline_authorization.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclPipelineAuthEnvironment(testutils.GenerateResourceName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "project_id"),
					resource.TestCheckResourceAttrSet(node, "resource_id"),
				),
			},
		},
	})
}

func TestAccPipelineAuthorization_allPipeline_variableGroup(t *testing.T) {
	node := "azuredevops_pipeline_authorization.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclAllPipelineAuthVariableGroup(testutils.GenerateResourceName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "project_id"),
					resource.TestCheckResourceAttrSet(node, "resource_id"),
				),
			},
		},
	})
}

func TestAccPipelineAuthorization_pipeline_variableGroup(t *testing.T) {
	node := "azuredevops_pipeline_authorization.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclPipelineAuthVariableGroup(testutils.GenerateResourceName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "project_id"),
					resource.TestCheckResourceAttrSet(node, "resource_id"),
				),
			},
		},
	})
}

func TestAccPipelineAuthorization_allPipeline_endpoint(t *testing.T) {
	node := "azuredevops_pipeline_authorization.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclAllPipelineAuthEndpoint(testutils.GenerateResourceName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "project_id"),
					resource.TestCheckResourceAttrSet(node, "resource_id"),
				),
			},
		},
	})
}

func TestAccPipelineAuthorization_pipeline_endpoint(t *testing.T) {
	node := "azuredevops_pipeline_authorization.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclPipelineAuthEndpoint(testutils.GenerateResourceName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "project_id"),
					resource.TestCheckResourceAttrSet(node, "resource_id"),
				),
			},
		},
	})
}

func TestAccPipelineAuthorization_allPipeline_repository(t *testing.T) {
	node := "azuredevops_pipeline_authorization.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclAllPipelineAuthRepository(testutils.GenerateResourceName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "project_id"),
					resource.TestCheckResourceAttrSet(node, "resource_id"),
				),
			},
		},
	})
}

func TestAccPipelineAuthorization_pipeline_repository(t *testing.T) {
	node := "azuredevops_pipeline_authorization.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclPipelineAuthRepository(testutils.GenerateResourceName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "project_id"),
					resource.TestCheckResourceAttrSet(node, "resource_id"),
				),
			},
		},
	})
}

func TestAccPipelineAuthorization_allPipeline_repository_crossProject(t *testing.T) {
	node := "azuredevops_pipeline_authorization.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclAllPipelineAuthRepositoryCrossProject(testutils.GenerateResourceName(), testutils.GenerateResourceName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "project_id"),
					resource.TestCheckResourceAttrSet(node, "resource_id"),
				),
			},
		},
	})
}

func TestAccPipelineAuthorization_pipeline_repository_crossProject(t *testing.T) {
	node := "azuredevops_pipeline_authorization.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, nil)
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclPipelineAuthRepositoryCrossPeojct(testutils.GenerateResourceName(), testutils.GenerateResourceName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(node, "project_id"),
					resource.TestCheckResourceAttrSet(node, "resource_id"),
				),
			},
		},
	})
}

func hclAllPipelineAuthQueue(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_agent_pool" "test" {
  name           = "%[1]s"
  auto_provision = false
  auto_update    = false
}

resource "azuredevops_agent_queue" "test" {
  project_id    = azuredevops_project.test.id
  agent_pool_id = azuredevops_agent_pool.test.id
}

resource "azuredevops_pipeline_authorization" "test" {
  project_id  = azuredevops_project.test.id
  resource_id = azuredevops_agent_queue.test.id
  type        = "queue"
}
`, name)
}

func hclPipelineAuthQueue(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_agent_pool" "test" {
  name           = "%[1]s"
  auto_provision = false
  auto_update    = false
}

resource "azuredevops_agent_queue" "test" {
  project_id    = azuredevops_project.test.id
  agent_pool_id = azuredevops_agent_pool.test.id
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"

  repository {
    repo_type = "TfsGit"
    repo_id   = data.azuredevops_git_repository.test.id
    yml_path  = "azure-pipelines.yml"
  }
}

resource "azuredevops_pipeline_authorization" "test" {
  project_id  = azuredevops_project.test.id
  resource_id = azuredevops_agent_queue.test.id
  pipeline_id = azuredevops_build_definition.test.id
  type        = "queue"
}
`, name)
}

func hclMultiPipelineAuthQueue(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_agent_pool" "test" {
  name           = "%[1]s"
  auto_provision = false
  auto_update    = false
}

resource "azuredevops_agent_queue" "test" {
  project_id    = azuredevops_project.test.id
  agent_pool_id = azuredevops_agent_pool.test.id
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"

  repository {
    repo_type = "TfsGit"
    repo_id   = data.azuredevops_git_repository.test.id
    yml_path  = "azure-pipelines.yml"
  }
}

resource "azuredevops_build_definition" "test2" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s2"

  repository {
    repo_type = "TfsGit"
    repo_id   = data.azuredevops_git_repository.test.id
    yml_path  = "azure-pipelines.yml"
  }
}

resource "azuredevops_pipeline_authorization" "test" {
  project_id  = azuredevops_project.test.id
  resource_id = azuredevops_agent_queue.test.id
  pipeline_id = azuredevops_build_definition.test.id
  type        = "queue"
}

resource "azuredevops_pipeline_authorization" "test2" {
  project_id  = azuredevops_project.test.id
  resource_id = azuredevops_agent_queue.test.id
  pipeline_id = azuredevops_build_definition.test2.id
  type        = "queue"
}
`, name)
}

func hclAllPipelineWithPipoelineAuthQueue(name string) string {
	{
		return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_agent_pool" "test" {
  name           = "%[1]s"
  auto_provision = false
  auto_update    = false
}

resource "azuredevops_agent_queue" "test" {
  project_id    = azuredevops_project.test.id
  agent_pool_id = azuredevops_agent_pool.test.id
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"

  repository {
    repo_type = "TfsGit"
    repo_id   = data.azuredevops_git_repository.test.id
    yml_path  = "azure-pipelines.yml"
  }
}

resource "azuredevops_pipeline_authorization" "test" {
  project_id  = azuredevops_project.test.id
  resource_id = azuredevops_agent_queue.test.id
  pipeline_id = azuredevops_build_definition.test.id
  type        = "queue"
}

resource "azuredevops_pipeline_authorization" "test_all" {
  project_id  = azuredevops_project.test.id
  resource_id = azuredevops_agent_queue.test.id
  type        = "queue"
}

`, name)
	}
}

func hclAllPipelineAuthEnvironment(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_environment" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_pipeline_authorization" "test" {
  project_id  = azuredevops_project.test.id
  resource_id = azuredevops_environment.test.id
  type        = "environment"
}
`, name)
}

func hclPipelineAuthEnvironment(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_environment" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"

  repository {
    repo_type = "TfsGit"
    repo_id   = data.azuredevops_git_repository.test.id
    yml_path  = "azure-pipelines.yml"
  }
}

resource "azuredevops_pipeline_authorization" "test" {
  project_id  = azuredevops_project.test.id
  resource_id = azuredevops_environment.test.id
  type        = "environment"
  pipeline_id = azuredevops_build_definition.test.id
}
`, name)
}

func hclAllPipelineAuthVariableGroup(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_variable_group" "test" {
  project_id   = azuredevops_project.test.id
  name         = "%[1]s"
  allow_access = true

  variable {
    name  = "key1"
    value = "val1"
  }
}

resource "azuredevops_pipeline_authorization" "test" {
  project_id  = azuredevops_project.test.id
  resource_id = azuredevops_variable_group.test.id
  type        = "variablegroup"
}
`, name)
}

func hclPipelineAuthVariableGroup(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_variable_group" "test" {
  project_id   = azuredevops_project.test.id
  name         = "%[1]s"
  allow_access = true

  variable {
    name  = "key1"
    value = "val1"
  }
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"

  repository {
    repo_type = "TfsGit"
    repo_id   = data.azuredevops_git_repository.test.id
    yml_path  = "azure-pipelines.yml"
  }
}

resource "azuredevops_pipeline_authorization" "test" {
  project_id  = azuredevops_project.test.id
  resource_id = azuredevops_variable_group.test.id
  type        = "variablegroup"
  pipeline_id = azuredevops_build_definition.test.id
}
`, name)
}

func hclAllPipelineAuthEndpoint(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_github" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%[1]s"

  auth_personal {
    personal_access_token = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
}

resource "azuredevops_pipeline_authorization" "test" {
  project_id  = azuredevops_project.test.id
  resource_id = azuredevops_serviceendpoint_github.test.id
  type        = "endpoint"
}
`, name)
}

func hclPipelineAuthEndpoint(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_github" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%[1]s"

  auth_personal {
    personal_access_token = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"

  repository {
    repo_type = "TfsGit"
    repo_id   = data.azuredevops_git_repository.test.id
    yml_path  = "azure-pipelines.yml"
  }
}

resource "azuredevops_pipeline_authorization" "test" {
  project_id  = azuredevops_project.test.id
  resource_id = azuredevops_serviceendpoint_github.test.id
  type        = "endpoint"
  pipeline_id = azuredevops_build_definition.test.id
}
`, name)
}

func hclAllPipelineAuthRepository(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_pipeline_authorization" "test" {
  project_id  = azuredevops_project.test.id
  resource_id = "${azuredevops_project.test.id}.${data.azuredevops_git_repository.test.id}"
  type        = "repository"
}
`, name)
}

func hclPipelineAuthRepository(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"

  repository {
    repo_type = "TfsGit"
    repo_id   = data.azuredevops_git_repository.test.id
    yml_path  = "azure-pipelines.yml"
  }
}

resource "azuredevops_pipeline_authorization" "test" {
  project_id  = azuredevops_project.test.id
  resource_id = "${azuredevops_project.test.id}.${data.azuredevops_git_repository.test.id}"
  pipeline_id = azuredevops_build_definition.test.id
  type        = "repository"
}
`, name)
}

func hclAllPipelineAuthRepositoryCrossProject(name, name2 string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_project" "test2" {
  name               = "%[2]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

resource "azuredevops_pipeline_authorization" "test" {
  project_id  = azuredevops_project.test2.id
  resource_id = "${azuredevops_project.test.id}.${data.azuredevops_git_repository.test.id}"
  type        = "repository"
}
`, name, name2)
}

func hclPipelineAuthRepositoryCrossPeojct(name, name2 string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_project" "test2" {
  name               = "%[2]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}

data "azuredevops_git_repository" "test2" {
  project_id = azuredevops_project.test2.id
  name       = "%[2]s"
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test2.id
  name       = "%[1]s"

  repository {
    repo_type = "TfsGit"
    repo_id   = data.azuredevops_git_repository.test2.id
    yml_path  = "azure-pipelines.yml"
  }
}

resource "azuredevops_pipeline_authorization" "test" {
  project_id  = azuredevops_project.test2.id
  resource_id = "${azuredevops_project.test.id}.${data.azuredevops_git_repository.test.id}"
  pipeline_id = azuredevops_build_definition.test.id
  type        = "repository"
}
`, name, name2)
}
