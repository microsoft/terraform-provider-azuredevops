// +build all resource_authorization
// +build !exclude_resource_authorization

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccResourceAuthorization_CRUD(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	serviceEndpointHCL := testutils.HclServiceEndpointGitHubResource(projectName, serviceEndpointName)
	authedHCL := testutils.HclResourceAuthorization("azuredevops_serviceendpoint_github.serviceendpoint.id", true)
	unAuthedHCL := testutils.HclResourceAuthorization("azuredevops_serviceendpoint_github.serviceendpoint.id", false)

	tfAuthNode := "azuredevops_resource_authorization.auth"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf("%s\n%s", serviceEndpointHCL, authedHCL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfAuthNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfAuthNode, "resource_id"),
					resource.TestCheckResourceAttr(tfAuthNode, "authorized", "true"),
				),
			}, {
				Config: fmt.Sprintf("%s\n%s", serviceEndpointHCL, unAuthedHCL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfAuthNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfAuthNode, "resource_id"),
					resource.TestCheckResourceAttr(tfAuthNode, "authorized", "false"),
				),
			},
		},
	})
}

func TestAccResourceAuthorization_Definition_CRUD(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	variableGroupName := testutils.GenerateResourceName()
	repositoryName := testutils.GenerateResourceName()
	buildDefinitionName := testutils.GenerateResourceName()

	buildDefinitionDHCL := testutils.HclBuildDefinitionResourceTfsGit(projectName, repositoryName, buildDefinitionName, `\`)
	variableGroupHCL := testutils.HclVariableGroupResource(variableGroupName, true)
	authedHCL := testutils.HclDefinitionResourceAuthorization("azuredevops_variable_group.vg.id", "azuredevops_build_definition.build.id", "variablegroup", true)
	unAuthedHCL := testutils.HclDefinitionResourceAuthorization("azuredevops_variable_group.vg.id", "azuredevops_build_definition.build.id", "variablegroup", false)

	tfAuthNode := "azuredevops_resource_authorization.auth"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf("%s\n%s\n%s", buildDefinitionDHCL, variableGroupHCL, authedHCL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfAuthNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfAuthNode, "definition_id"),
					resource.TestCheckResourceAttrSet(tfAuthNode, "resource_id"),
					resource.TestCheckResourceAttr(tfAuthNode, "authorized", "true"),
				),
			}, {
				Config: fmt.Sprintf("%s\n%s\n%s", buildDefinitionDHCL, variableGroupHCL, unAuthedHCL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfAuthNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfAuthNode, "definition_id"),
					resource.TestCheckResourceAttrSet(tfAuthNode, "resource_id"),
					resource.TestCheckResourceAttr(tfAuthNode, "authorized", "false"),
				),
			},
		},
	})
}
