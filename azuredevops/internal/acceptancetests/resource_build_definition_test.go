// +build all resource_build_definition
// +build !exclude_resource_build_definition

package acceptancetests

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
)

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccBuildDefinition_Create_Update_Import(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	buildDefinitionPathEmpty := `\`
	buildDefinitionNameFirst := testutils.GenerateResourceName()
	buildDefinitionNameSecond := testutils.GenerateResourceName()
	buildDefinitionNameThird := testutils.GenerateResourceName()

	buildDefinitionPathFirst := `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	buildDefinitionPathSecond := `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	buildDefinitionPathThird := `\` + buildDefinitionNameFirst + `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	buildDefinitionPathFourth := `\` + buildDefinitionNameSecond + `\` + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfBuildDefNode := "azuredevops_build_definition.build"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclBuildDefinitionResourceGitHub(projectName, buildDefinitionNameFirst, buildDefinitionPathEmpty),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(buildDefinitionNameFirst),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", buildDefinitionNameFirst),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", buildDefinitionPathEmpty),
				),
			}, {
				Config: testutils.HclBuildDefinitionResourceGitHub(projectName, buildDefinitionNameSecond, buildDefinitionPathEmpty),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(buildDefinitionNameSecond),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", buildDefinitionNameSecond),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", buildDefinitionPathEmpty),
				),
			}, {
				Config: testutils.HclBuildDefinitionResourceGitHub(projectName, buildDefinitionNameFirst, buildDefinitionPathFirst),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(buildDefinitionNameFirst),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", buildDefinitionNameFirst),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", buildDefinitionPathFirst),
				),
			}, {
				Config: testutils.HclBuildDefinitionResourceGitHub(projectName, buildDefinitionNameFirst,
					buildDefinitionPathSecond),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(buildDefinitionNameFirst),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", buildDefinitionNameFirst),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", buildDefinitionPathSecond),
				),
			}, {
				Config: testutils.HclBuildDefinitionResourceGitHub(projectName, buildDefinitionNameFirst, buildDefinitionPathThird),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(buildDefinitionNameFirst),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", buildDefinitionNameFirst),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", buildDefinitionPathThird),
				),
			}, {
				Config: testutils.HclBuildDefinitionResourceGitHub(projectName, buildDefinitionNameFirst, buildDefinitionPathFourth),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(buildDefinitionNameFirst),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", buildDefinitionNameFirst),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", buildDefinitionPathFourth),
				),
			}, {
				Config: testutils.HclBuildDefinitionResourceTfsGit(projectName, gitRepoName, buildDefinitionNameThird, buildDefinitionPathEmpty),
				Check: resource.ComposeTestCheckFunc(
					checkBuildDefinitionExists(buildDefinitionNameThird),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "revision"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "repository.0.repo_id"),
					resource.TestCheckResourceAttr(tfBuildDefNode, "name", buildDefinitionNameThird),
					resource.TestCheckResourceAttr(tfBuildDefNode, "path", buildDefinitionPathEmpty),
				),
			}, {
				// Resource Acceptance Testing https://www.terraform.io/docs/extend/resources/import.html#resource-acceptance-testing-implementation
				ResourceName:      tfBuildDefNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfBuildDefNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Verifies a build for Bitbucket can happen. Note: the update/import logic is tested in other tests
func TestAccBuildDefinitionBitbucket_Create(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      testutils.HclBuildDefinitionResourceBitbucket(projectName, "build-def-name", "\\", ""),
				ExpectError: regexp.MustCompile("bitbucket repositories need a referenced service connection ID"),
			}, {
				Config: testutils.HclBuildDefinitionResourceBitbucket(projectName, "build-def-name", "\\", "some-service-connection"),
				Check:  checkBuildDefinitionExists("build-def-name"),
			},
		},
	})
}

// Verifies a build for with variables can create and update, including secret variables
func TestAccBuildDefinition_WithVariables_CreateAndUpdate(t *testing.T) {
	name := testutils.GenerateResourceName()
	tfNode := "azuredevops_build_definition.b"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkBuildDefinitionDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclBuildDefinitionWithVariables("foo1", "bar1", name),
				Check:  checkForVariableValues(tfNode, "foo1", "bar1"),
			}, {
				Config: testutils.HclBuildDefinitionWithVariables("foo2", "bar2", name),
				Check:  checkForVariableValues(tfNode, "foo2", "bar2"),
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

// Given the name of an AzDO build definition, this will return a function that will check whether
// or not the definition (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func checkBuildDefinitionExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		buildDef, ok := s.RootModule().Resources["azuredevops_build_definition.build"]
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
// *after* terrafform destroys the resource but *before* the state is wiped clean.
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
	clients := testutils.GetProvider().Meta().(*config.AggregatedClient)
	return clients.BuildClient.GetDefinition(clients.Ctx, build.GetDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &buildDefID,
	})
}
