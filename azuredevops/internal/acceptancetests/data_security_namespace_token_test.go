package acceptancetests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// isAzureDevOpsServices checks if the tests are running against Azure DevOps Services (dev.azure.com)
// vs Azure DevOps Server (on-premises)
func isAzureDevOpsServices() bool {
	orgURL := os.Getenv("AZDO_ORG_SERVICE_URL")
	return strings.Contains(strings.ToLower(orgURL), "dev.azure.com")
}

// TestAccDataSecurityNamespaceToken_collection tests token generation for Collection namespace
func TestAccDataSecurityNamespaceToken_collection(t *testing.T) {
	tfNode := "data.azuredevops_security_namespace_token.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_collection(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "token", "NAMESPACE:"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_collection() string {
	return `
data "azuredevops_security_namespace_token" "test" {
  namespace_name = "Collection"
}
`
}

// TestAccDataSecurityNamespaceToken_project tests token generation for Project namespace
func TestAccDataSecurityNamespaceToken_project(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_project(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckOutput("token_matches", "true"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_project(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_security_namespace_token" "test" {
  namespace_name = "Project"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

output "token_matches" {
  value = data.azuredevops_security_namespace_token.test.token == "$PROJECT:vstfs:///Classification/TeamProject/${azuredevops_project.test.id}"
}
`, projectName)
}

// TestAccDataSecurityNamespaceToken_gitRepositories_project tests token generation for Git Repositories namespace (project level)
func TestAccDataSecurityNamespaceToken_gitRepositories_project(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_gitRepositories_project(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckOutput("token_matches", "true"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_gitRepositories_project(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_security_namespace_token" "test" {
  namespace_name = "Git Repositories"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

output "token_matches" {
  value = data.azuredevops_security_namespace_token.test.token == "repoV2/${azuredevops_project.test.id}"
}
`, projectName)
}

// TestAccDataSecurityNamespaceToken_gitRepositories_repository tests token generation for Git Repositories namespace (repository level)
func TestAccDataSecurityNamespaceToken_gitRepositories_repository(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_gitRepositories_repository(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckOutput("token_matches", "true"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_gitRepositories_repository(projectName, repoName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  initialization {
    init_type = "Clean"
  }
}

data "azuredevops_security_namespace_token" "test" {
  namespace_name = "Git Repositories"
  identifiers = {
    project_id    = azuredevops_project.test.id
    repository_id = azuredevops_git_repository.test.id
  }
}

output "token_matches" {
  value = data.azuredevops_security_namespace_token.test.token == "repoV2/${azuredevops_project.test.id}/${azuredevops_git_repository.test.id}"
}
`, projectName, repoName)
}

// TestAccDataSecurityNamespaceToken_build_project tests token generation for Build namespace (project level)
func TestAccDataSecurityNamespaceToken_build_project(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_build_project(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "token"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_build_project(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_security_namespace_token" "test" {
  namespace_name = "Build"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}
`, projectName)
}

// TestAccDataSecurityNamespaceToken_build_path tests token generation for Build namespace (with path)
func TestAccDataSecurityNamespaceToken_build_path(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	folderName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_build_path(projectName, folderName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "token"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_build_path(projectName, folderName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_build_folder" "test" {
  project_id  = azuredevops_project.test.id
  path        = "\\%[2]s"
  description = "Test folder"
}

data "azuredevops_security_namespace_token" "test" {
  namespace_name = "Build"
  identifiers = {
    project_id = azuredevops_project.test.id
    path       = azuredevops_build_folder.test.path
  }
}
`, projectName, folderName)
}

// TestAccDataSecurityNamespaceToken_build_definition tests token generation for Build namespace (with definition)
func TestAccDataSecurityNamespaceToken_build_definition(t *testing.T) {
	if !isAzureDevOpsServices() {
		t.Skip("Skipping test because it requires Azure Pipelines agent pool which may not be available in Azure DevOps Server environments")
	}
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()
	buildName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_build_definition(projectName, repoName, buildName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "token"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_build_definition(projectName, repoName, buildName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_build_definition" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[3]s"

  ci_trigger {
    use_yaml = true
  }

  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.test.id
    branch_name = azuredevops_git_repository.test.default_branch
    yml_path    = "azure-pipelines.yml"
  }
}

data "azuredevops_security_namespace_token" "test" {
  namespace_name = "Build"
  identifiers = {
    project_id    = azuredevops_project.test.id
    definition_id = azuredevops_build_definition.test.id
  }
}
`, projectName, repoName, buildName)
}

// TestAccDataSecurityNamespaceToken_css tests token generation for CSS (Areas) namespace
func TestAccDataSecurityNamespaceToken_css(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_css(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckOutput("token_matches", "true"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_css(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_security_namespace_token" "test" {
  namespace_name = "CSS"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

output "token_matches" {
  value = startswith(data.azuredevops_security_namespace_token.test.token, "vstfs:///Classification/Node/") ? "true" : "false"
}
`, projectName)
}

// TestAccDataSecurityNamespaceToken_iteration tests token generation for Iteration namespace
func TestAccDataSecurityNamespaceToken_iteration(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_iteration(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckOutput("token_matches", "true"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_iteration(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_security_namespace_token" "test" {
  namespace_name = "Iteration"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

output "token_matches" {
  value = startswith(data.azuredevops_security_namespace_token.test.token, "vstfs:///Classification/Node/") ? "true" : "false"
}
`, projectName)
}

// TestAccDataSecurityNamespaceToken_tagging_project tests token generation for Tagging namespace (project level)
func TestAccDataSecurityNamespaceToken_tagging_project(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_tagging_project(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "token"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_tagging_project(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_security_namespace_token" "test" {
  namespace_name = "Tagging"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}
`, projectName)
}

// TestAccDataSecurityNamespaceToken_tagging_collection tests token generation for Tagging namespace (collection level)
func TestAccDataSecurityNamespaceToken_tagging_collection(t *testing.T) {
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_tagging_collection(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "token", ""),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_tagging_collection() string {
	return `
data "azuredevops_security_namespace_token" "test" {
  namespace_name = "Tagging"
}
`
}

// TestAccDataSecurityNamespaceToken_serviceHooks_project tests token generation for Service Hooks namespace (project level)
func TestAccDataSecurityNamespaceToken_serviceHooks_project(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_serviceHooks_project(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckOutput("token_matches", "true"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_serviceHooks_project(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_security_namespace_token" "test" {
  namespace_name = "ServiceHooks"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

output "token_matches" {
  value = data.azuredevops_security_namespace_token.test.token == "PublisherSecurity/${azuredevops_project.test.id}"
}
`, projectName)
}

// TestAccDataSecurityNamespaceToken_serviceHooks_collection tests token generation for Service Hooks namespace (collection level)
func TestAccDataSecurityNamespaceToken_serviceHooks_collection(t *testing.T) {
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_serviceHooks_collection(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "token", "PublisherSecurity"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_serviceHooks_collection() string {
	return `
data "azuredevops_security_namespace_token" "test" {
  namespace_name = "ServiceHooks"
}
`
}

// TestAccDataSecurityNamespaceToken_workItemQueryFolders tests token generation for Work Item Query Folders namespace
func TestAccDataSecurityNamespaceToken_workItemQueryFolders(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_workItemQueryFolders(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckOutput("token_matches", "true"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_workItemQueryFolders(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_security_namespace_token" "test" {
  namespace_name = "WorkItemQueryFolders"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

output "token_matches" {
  value = data.azuredevops_security_namespace_token.test.token == "$/${azuredevops_project.test.id}"
}
`, projectName)
}

// TestAccDataSecurityNamespaceToken_analytics tests token generation for Analytics namespace
func TestAccDataSecurityNamespaceToken_analytics(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_analytics(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckOutput("token_matches", "true"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_analytics(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_security_namespace_token" "test" {
  namespace_name = "Analytics"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

output "token_matches" {
  value = data.azuredevops_security_namespace_token.test.token == "$/${azuredevops_project.test.id}"
}
`, projectName)
}

// TestAccDataSecurityNamespaceToken_analyticsViews tests token generation for AnalyticsViews namespace
func TestAccDataSecurityNamespaceToken_analyticsViews(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_analyticsViews(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckOutput("token_matches", "true"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_analyticsViews(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_security_namespace_token" "test" {
  namespace_name = "AnalyticsViews"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

output "token_matches" {
  value = data.azuredevops_security_namespace_token.test.token == "$/Shared/${azuredevops_project.test.id}"
}
`, projectName)
}

// TestAccDataSecurityNamespaceToken_process tests token generation for Process namespace
func TestAccDataSecurityNamespaceToken_process(t *testing.T) {
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_process(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "token", "$PROCESS:"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_process() string {
	return `
data "azuredevops_security_namespace_token" "test" {
  namespace_name = "Process"
}
`
}

// TestAccDataSecurityNamespaceToken_auditLog tests token generation for AuditLog namespace
func TestAccDataSecurityNamespaceToken_auditLog(t *testing.T) {
	if !isAzureDevOpsServices() {
		t.Skip("Skipping test because AuditLog namespace is only available in Azure DevOps Services (dev.azure.com), not Azure DevOps Server")
	}
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_auditLog(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "token", "AllPermissions"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_auditLog() string {
	return `
data "azuredevops_security_namespace_token" "test" {
  namespace_name = "AuditLog"
}
`
}

// TestAccDataSecurityNamespaceToken_buildAdministration tests token generation for BuildAdministration namespace
func TestAccDataSecurityNamespaceToken_buildAdministration(t *testing.T) {
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_buildAdministration(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "token", "BuildPrivileges"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_buildAdministration() string {
	return `
data "azuredevops_security_namespace_token" "test" {
  namespace_name = "BuildAdministration"
}
`
}

// TestAccDataSecurityNamespaceToken_server tests token generation for Server namespace
func TestAccDataSecurityNamespaceToken_server(t *testing.T) {
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_server(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "token", "FrameworkGlobalSecurity"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_server() string {
	return `
data "azuredevops_security_namespace_token" "test" {
  namespace_name = "Server"
}
`
}

// TestAccDataSecurityNamespaceToken_versionControlPrivileges tests token generation for VersionControlPrivileges namespace
func TestAccDataSecurityNamespaceToken_versionControlPrivileges(t *testing.T) {
	tfNode := "data.azuredevops_security_namespace_token.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_versionControlPrivileges(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "token", "Global"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_versionControlPrivileges() string {
	return `
data "azuredevops_security_namespace_token" "test" {
  namespace_name = "VersionControlPrivileges"
}
`
}
