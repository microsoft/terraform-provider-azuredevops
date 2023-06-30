//go:build (all || resource_check_branch_control) && !exclude_approvalsandchecks
// +build all resource_check_branch_control
// +build !exclude_approvalsandchecks

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

var checkName = "Extend a required template"

func TestAccCheckRequiredTemplate_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repositoryName := "test-repo"
	repositoryRef := "refs/heads/master"
	templatePath := "templ/path1.yaml"

	resourceType := "azuredevops_check_required_template"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclRequiredTemplateCheckResourceBasic(projectName, repositoryName, repositoryRef, templatePath),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, checkName),
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					// TODO: check nested properties
					// resource.TestCheckResourceAttr(tfCheckNode, "allowed_branches", branches),
					// resource.TestCheckResourceAttr(tfCheckNode, "display_name", checkName),
				),
			},
		},
	})
}

func TestAccCheckRequiredTemplate_complete(t *testing.T) {
	// TODO: rewrite function to required template check
	projectName := testutils.GenerateResourceName()
	repositoryType := "github"
	repositoryName := "proj/test-repo"
	repositoryRef := "refs/heads/master"
	templatePath := "templ/path1.yaml"

	resourceType := "azuredevops_check_required_template"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclRequiredTemplateCheckResourceComplete(projectName, repositoryType, repositoryName, repositoryRef, templatePath),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, checkName),
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					// TODO: check nested properties
					// resource.TestCheckResourceAttr(tfCheckNode, "allowed_branches", branches),
					// resource.TestCheckResourceAttr(tfCheckNode, "display_name", checkName),
					// resource.TestCheckResourceAttr(tfCheckNode, "verify_branch_protection", "true"),
					// resource.TestCheckResourceAttr(tfCheckNode, "ignore_unknown_protection_status", "false"),
				),
			},
		},
	})
}

func TestAccCheckRequiredTemplate_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repositoryNameFirst := "test-repo"
	repositoryRefFirst := "refs/heads/master"
	templatePathFirst := "templ/path1.yaml"

	repositoryTypeSecond := "github"
	repositoryNameSecond := "test-project/test-repo"
	repositoryRefSecond := "refs/heads/main"
	templatePathSecond := "templ/path2.yaml"

	resourceType := "azuredevops_check_required_template"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclRequiredTemplateCheckResourceBasic(projectName, repositoryNameFirst, repositoryRefFirst, templatePathFirst),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, checkName),
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					// TODO: how to check nested properties?
					// resource.TestCheckResourceAttr(tfCheckNode, "allowed_branches", ),
					// resource.TestCheckResourceAttr(tfCheckNode, "display_name", checkNameFirst),
				),
			},
			{
				Config: hclRequiredTemplateCheckResourceUpdate(projectName, repositoryTypeSecond, repositoryNameSecond, repositoryRefSecond, templatePathSecond),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, checkName),
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					// TODO: how to check nested properties?
					// resource.TestCheckResourceAttr(tfCheckNode, "allowed_branches", ),
					// resource.TestCheckResourceAttr(tfCheckNode, "display_name", checkNameFirst),
				),
			},
		},
	})
}

func hclRequiredTemplateCheckResourceBasic(projectName, repositoryName, repositoryRef, templatePath string) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_branch_control" "test" {
  project_id           = azuredevops_project.project.id
  target_resource_id   = azuredevops_serviceendpoint_generic.test.id
  target_resource_type = "endpoint"
  required_template {
	repository_name = "%s"
	repository_ref = "%s"
	template_path = "%s"
  }
}`, repositoryName, repositoryRef, templatePath)

	genericcheckResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericcheckResource, checkResource)
}

func hclRequiredTemplateCheckResourceComplete(projectName, repository_type, repositoryName, repositoryRef, templatePath string) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_branch_control" "test" {
  project_id                       = azuredevops_project.project.id
  target_resource_id               = azuredevops_serviceendpoint_generic.test.id
  target_resource_type             = "endpoint"
  required_template {
	repository_type = "%s"
	repository_name = "%s"
	repository_ref = "%s"
	template_path = "%s"
  }
}`, repository_type, repositoryName, repositoryRef, templatePath)

	genericcheckResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericcheckResource, checkResource)
}

func hclRequiredTemplateCheckResourceUpdate(projectName, repository_type, repositoryName, repositoryRef, templatePath string) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_branch_control" "test" {
  project_id                       = azuredevops_project.project.id
  target_resource_id               = azuredevops_serviceendpoint_generic.test.id
  target_resource_type             = "endpoint"
  required_template {
	repository_type = "%s"
	repository_name = "%s"
	repository_ref = "%s"
	template_path = "%s"
  }

}`, repository_type, repositoryName, repositoryRef, templatePath)

	genericcheckResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericcheckResource, checkResource)
}
