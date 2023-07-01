//go:build (all || resource_check_required_template) && !exclude_approvalsandchecks
// +build all resource_check_required_template
// +build !exclude_approvalsandchecks

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

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
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "required_template.0.repository_name", repositoryName),
					resource.TestCheckResourceAttr(tfCheckNode, "required_template.0.repository_ref", repositoryRef),
					resource.TestCheckResourceAttr(tfCheckNode, "required_template.0.template_path", templatePath),
				),
			},
		},
	})
}

func TestAccCheckRequiredTemplate_complete(t *testing.T) {
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
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "required_template.0.repository_type", repositoryType),
					resource.TestCheckResourceAttr(tfCheckNode, "required_template.0.repository_name", repositoryName),
					resource.TestCheckResourceAttr(tfCheckNode, "required_template.0.repository_ref", repositoryRef),
					resource.TestCheckResourceAttr(tfCheckNode, "required_template.0.template_path", templatePath),
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
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "required_template.0.repository_name", repositoryNameFirst),
					resource.TestCheckResourceAttr(tfCheckNode, "required_template.0.repository_ref", repositoryRefFirst),
					resource.TestCheckResourceAttr(tfCheckNode, "required_template.0.template_path", templatePathFirst),
				),
			},
			{
				Config: hclRequiredTemplateCheckResourceUpdate(projectName, repositoryTypeSecond, repositoryNameSecond, repositoryRefSecond, templatePathSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "required_template.0.repository_type", repositoryTypeSecond),
					resource.TestCheckResourceAttr(tfCheckNode, "required_template.0.repository_name", repositoryNameSecond),
					resource.TestCheckResourceAttr(tfCheckNode, "required_template.0.repository_ref", repositoryRefSecond),
					resource.TestCheckResourceAttr(tfCheckNode, "required_template.0.template_path", templatePathSecond),
				),
			},
		},
	})
}

func hclRequiredTemplateCheckResourceBasic(projectName, repositoryName, repositoryRef, templatePath string) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_required_template" "test" {
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
resource "azuredevops_check_required_template" "test" {
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
resource "azuredevops_check_required_template" "test" {
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
