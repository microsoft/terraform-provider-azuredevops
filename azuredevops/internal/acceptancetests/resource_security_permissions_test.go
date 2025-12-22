package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// TestAccSecurityPermissions_ProjectNamespace tests setting permissions on the Project namespace
func TestAccSecurityPermissions_ProjectNamespace(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_security_permissions.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityPermissionsProjectNamespace(projectName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "namespace_id"),
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "3"),
					resource.TestCheckResourceAttr(tfNode, "permissions.GENERIC_READ", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.GENERIC_WRITE", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.DELETE", "deny"),
					resource.TestCheckResourceAttr(tfNode, "replace", "false"),
				),
			},
		},
	})
}

// TestAccSecurityPermissions_ProjectNamespaceUpdate tests updating permissions
func TestAccSecurityPermissions_ProjectNamespaceUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_security_permissions.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityPermissionsProjectNamespace(projectName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "3"),
					resource.TestCheckResourceAttr(tfNode, "permissions.GENERIC_READ", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.GENERIC_WRITE", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.DELETE", "deny"),
				),
			},
			{
				Config: hclSecurityPermissionsProjectNamespaceUpdated(projectName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "4"),
					resource.TestCheckResourceAttr(tfNode, "permissions.GENERIC_READ", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.GENERIC_WRITE", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.DELETE", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.MANAGE_PROPERTIES", "notset"),
				),
			},
		},
	})
}

// TestAccSecurityPermissions_GitRepositoryNamespace tests Git repository permissions
func TestAccSecurityPermissions_GitRepositoryNamespace(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()
	tfNode := "azuredevops_security_permissions.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityPermissionsGitRepoNamespace(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "namespace_id"),
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "4"),
					resource.TestCheckResourceAttr(tfNode, "permissions.GenericRead", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.GenericContribute", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ForcePush", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ManagePermissions", "deny"),
				),
			},
		},
	})
}

// TestAccSecurityPermissions_WithReplace tests the replace functionality
func TestAccSecurityPermissions_WithReplace(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_security_permissions.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityPermissionsWithReplace(projectName, false),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "replace", "false"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "2"),
				),
			},
			{
				Config: hclSecurityPermissionsWithReplace(projectName, true),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "replace", "true"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "2"),
				),
			},
		},
	})
}

// TestAccSecurityPermissions_BuildNamespace tests Build definition permissions
func TestAccSecurityPermissions_BuildNamespace(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_security_permissions.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityPermissionsBuildNamespace(projectName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "namespace_id"),
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "3"),
					resource.TestCheckResourceAttr(tfNode, "permissions.ViewBuilds", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.EditBuildQuality", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.DeleteBuilds", "deny"),
				),
			},
		},
	})
}

// TestAccSecurityPermissions_IdentityNamespace tests identity/group permissions
func TestAccSecurityPermissions_IdentityNamespace(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_security_permissions.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityPermissionsIdentityNamespace(projectName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "namespace_id"),
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "2"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Read", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.Write", "deny"),
				),
			},
		},
	})
}

// TestAccSecurityPermissions_MultiplePermissionSets tests managing multiple permission sets
func TestAccSecurityPermissions_MultiplePermissionSets(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode1 := "azuredevops_security_permissions.test1"
	tfNode2 := "azuredevops_security_permissions.test2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityPermissionsMultipleSets(projectName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					// First permission set (Readers)
					resource.TestCheckResourceAttrSet(tfNode1, "namespace_id"),
					resource.TestCheckResourceAttrSet(tfNode1, "token"),
					resource.TestCheckResourceAttrSet(tfNode1, "principal"),
					resource.TestCheckResourceAttr(tfNode1, "permissions.%", "2"),
					resource.TestCheckResourceAttr(tfNode1, "permissions.GENERIC_READ", "allow"),
					resource.TestCheckResourceAttr(tfNode1, "permissions.GENERIC_WRITE", "deny"),
					// Second permission set (Contributors)
					resource.TestCheckResourceAttrSet(tfNode2, "namespace_id"),
					resource.TestCheckResourceAttrSet(tfNode2, "token"),
					resource.TestCheckResourceAttrSet(tfNode2, "principal"),
					resource.TestCheckResourceAttr(tfNode2, "permissions.%", "2"),
					resource.TestCheckResourceAttr(tfNode2, "permissions.GENERIC_READ", "allow"),
					resource.TestCheckResourceAttr(tfNode2, "permissions.GENERIC_WRITE", "allow"),
				),
			},
		},
	})
}

// TestAccSecurityPermissions_OrganizationGroup tests with organization-level groups
func TestAccSecurityPermissions_OrganizationGroup(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_security_permissions.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityPermissionsOrganizationGroup(projectName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "namespace_id"),
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckResourceAttrSet(tfNode, "principal"),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "2"),
				),
			},
		},
	})
}

// TestAccSecurityPermissions_CaseSensitivity tests case-insensitive permission values
func TestAccSecurityPermissions_CaseSensitivity(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_security_permissions.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclSecurityPermissionsCaseSensitivity(projectName),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "permissions.%", "3"),
					resource.TestCheckResourceAttr(tfNode, "permissions.GENERIC_READ", "allow"),
					resource.TestCheckResourceAttr(tfNode, "permissions.GENERIC_WRITE", "deny"),
					resource.TestCheckResourceAttr(tfNode, "permissions.DELETE", "notset"),
				),
			},
		},
	})
}

// HCL configurations for tests

func hclSecurityPermissionsProjectNamespace(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_security_namespaces" "all" {
}

locals {
  project_namespace = [
    for ns in data.azuredevops_security_namespaces.all.namespaces :
    ns if ns.name == "Project"
  ][0]
}

data "azuredevops_security_namespace_token" "project" {
  namespace_name = "Project"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

data "azuredevops_group" "test" {
  project_id = azuredevops_project.test.id
  name       = "Readers"
}

resource "azuredevops_security_permissions" "test" {
  namespace_id = local.project_namespace.namespace_id
  token        = data.azuredevops_security_namespace_token.project.token
  principal    = data.azuredevops_group.test.descriptor
  permissions = {
    GENERIC_READ  = "allow"
    GENERIC_WRITE = "deny"
    DELETE        = "deny"
  }
  replace = false
}
`, projectName)
}

func hclSecurityPermissionsProjectNamespaceUpdated(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_security_namespaces" "all" {
}

locals {
  project_namespace = [
    for ns in data.azuredevops_security_namespaces.all.namespaces :
    ns if ns.name == "Project"
  ][0]
}

data "azuredevops_security_namespace_token" "project" {
  namespace_name = "Project"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

data "azuredevops_group" "test" {
  project_id = azuredevops_project.test.id
  name       = "Readers"
}

resource "azuredevops_security_permissions" "test" {
  namespace_id = local.project_namespace.namespace_id
  token        = data.azuredevops_security_namespace_token.project.token
  principal    = data.azuredevops_group.test.descriptor
  permissions = {
    GENERIC_READ      = "allow"
    GENERIC_WRITE     = "allow"
    DELETE            = "deny"
    MANAGE_PROPERTIES = "notset"
  }
  replace = false
}
`, projectName)
}

func hclSecurityPermissionsGitRepoNamespace(projectName, repoName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
  initialization {
    init_type = "Clean"
  }
}

data "azuredevops_security_namespaces" "all" {
}

locals {
  git_namespace = [
    for ns in data.azuredevops_security_namespaces.all.namespaces :
    ns if ns.name == "Git Repositories"
  ][0]
}

data "azuredevops_security_namespace_token" "git_repo" {
  namespace_name = "Git Repositories"
  identifiers = {
    project_id    = azuredevops_project.test.id
    repository_id = azuredevops_git_repository.test.id
  }
}

data "azuredevops_group" "test" {
  project_id = azuredevops_project.test.id
  name       = "Contributors"
}

resource "azuredevops_security_permissions" "test" {
  namespace_id = local.git_namespace.namespace_id
  token        = data.azuredevops_security_namespace_token.git_repo.token
  principal    = data.azuredevops_group.test.descriptor
  permissions = {
    GenericRead       = "allow"
    GenericContribute = "allow"
    ForcePush         = "deny"
    ManagePermissions = "deny"
  }
  replace = false
}
`, projectName, repoName)
}

func hclSecurityPermissionsWithReplace(projectName string, replace bool) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_security_namespaces" "all" {
}

locals {
  project_namespace = [
    for ns in data.azuredevops_security_namespaces.all.namespaces :
    ns if ns.name == "Project"
  ][0]
}

data "azuredevops_security_namespace_token" "project" {
  namespace_name = "Project"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

data "azuredevops_group" "test" {
  project_id = azuredevops_project.test.id
  name       = "Readers"
}

resource "azuredevops_security_permissions" "test" {
  namespace_id = local.project_namespace.namespace_id
  token        = data.azuredevops_security_namespace_token.project.token
  principal    = data.azuredevops_group.test.descriptor
  permissions = {
    GENERIC_READ  = "allow"
    GENERIC_WRITE = "deny"
  }
  replace = %t
}
`, projectName, replace)
}

func hclSecurityPermissionsBuildNamespace(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_security_namespaces" "all" {
}

locals {
  build_namespace = [
    for ns in data.azuredevops_security_namespaces.all.namespaces :
    ns if ns.name == "Build"
  ][0]
}

data "azuredevops_security_namespace_token" "build" {
  namespace_name = "Build"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

data "azuredevops_group" "test" {
  project_id = azuredevops_project.test.id
  name       = "Readers"
}

resource "azuredevops_security_permissions" "test" {
  namespace_id = local.build_namespace.namespace_id
  token        = data.azuredevops_security_namespace_token.build.token
  principal    = data.azuredevops_group.test.descriptor
  permissions = {
    ViewBuilds        = "allow"
    EditBuildQuality  = "deny"
    DeleteBuilds      = "deny"
  }
  replace = false
}
`, projectName)
}

func hclSecurityPermissionsIdentityNamespace(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_security_namespaces" "all" {
}

locals {
  identity_namespace = [
    for ns in data.azuredevops_security_namespaces.all.namespaces :
    ns if ns.name == "Identity"
  ][0]
}

data "azuredevops_group" "readers" {
  project_id = azuredevops_project.test.id
  name       = "Readers"
}

data "azuredevops_group" "contributors" {
  project_id = azuredevops_project.test.id
  name       = "Contributors"
}

resource "azuredevops_security_permissions" "test" {
  namespace_id = local.identity_namespace.namespace_id
  token        = data.azuredevops_group.readers.descriptor
  principal    = data.azuredevops_group.contributors.descriptor
  permissions = {
    Read  = "allow"
    Write = "deny"
  }
  replace = false
}
`, projectName)
}

func hclSecurityPermissionsMultipleSets(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_security_namespaces" "all" {
}

locals {
  project_namespace = [
    for ns in data.azuredevops_security_namespaces.all.namespaces :
    ns if ns.name == "Project"
  ][0]
}

data "azuredevops_security_namespace_token" "project" {
  namespace_name = "Project"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

data "azuredevops_group" "readers" {
  project_id = azuredevops_project.test.id
  name       = "Readers"
}

data "azuredevops_group" "contributors" {
  project_id = azuredevops_project.test.id
  name       = "Contributors"
}

resource "azuredevops_security_permissions" "test1" {
  namespace_id = local.project_namespace.namespace_id
  token        = data.azuredevops_security_namespace_token.project.token
  principal    = data.azuredevops_group.readers.descriptor
  permissions = {
    GENERIC_READ  = "allow"
    GENERIC_WRITE = "deny"
  }
  replace = false
}

resource "azuredevops_security_permissions" "test2" {
  namespace_id = local.project_namespace.namespace_id
  token        = data.azuredevops_security_namespace_token.project.token
  principal    = data.azuredevops_group.contributors.descriptor
  permissions = {
    GENERIC_READ  = "allow"
    GENERIC_WRITE = "allow"
  }
  replace = false
}
`, projectName)
}

func hclSecurityPermissionsOrganizationGroup(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_security_namespaces" "all" {
}

locals {
  project_namespace = [
    for ns in data.azuredevops_security_namespaces.all.namespaces :
    ns if ns.name == "Project"
  ][0]
}

data "azuredevops_security_namespace_token" "project" {
  namespace_name = "Project"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

data "azuredevops_group" "test" {
  name = "Project Collection Build Service Accounts"
}

resource "azuredevops_security_permissions" "test" {
  namespace_id = local.project_namespace.namespace_id
  token        = data.azuredevops_security_namespace_token.project.token
  principal    = data.azuredevops_group.test.descriptor
  permissions = {
    GENERIC_READ  = "allow"
    GENERIC_WRITE = "allow"
  }
  replace = false
}
`, projectName)
}

func hclSecurityPermissionsCaseSensitivity(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_security_namespaces" "all" {
}

locals {
  project_namespace = [
    for ns in data.azuredevops_security_namespaces.all.namespaces :
    ns if ns.name == "Project"
  ][0]
}

data "azuredevops_security_namespace_token" "project" {
  namespace_name = "Project"
  identifiers = {
    project_id = azuredevops_project.test.id
  }
}

data "azuredevops_group" "test" {
  project_id = azuredevops_project.test.id
  name       = "Readers"
}

resource "azuredevops_security_permissions" "test" {
  namespace_id = local.project_namespace.namespace_id
  token        = data.azuredevops_security_namespace_token.project.token
  principal    = data.azuredevops_group.test.descriptor
  permissions = {
    GENERIC_READ  = "Allow"
    GENERIC_WRITE = "Deny"
    DELETE        = "NotSet"
  }
  replace = false
}
`, projectName)
}
