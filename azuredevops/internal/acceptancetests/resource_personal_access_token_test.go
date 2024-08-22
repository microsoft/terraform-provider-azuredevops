//go:build (all || resource_agent_queue) && !exclude_resource_agent_queue
// +build all resource_agent_queue
// +build !exclude_resource_agent_queue

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccResourcePersonalAccessToken_CreateAndUpdate(t *testing.T) {
	tokenName := testutils.GenerateResourceName()
	tfNode := "azuredevops_personal_access_token.pat"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclPersonalAccessTokenResource(tokenName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "authorization_id"),
					resource.TestCheckResourceAttrSet(tfNode, "scope"),
					resource.TestCheckResourceAttrSet(tfNode, "target_accounts"),
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckResourceAttrSet(tfNode, "valid_from"),
					resource.TestCheckResourceAttrSet(tfNode, "valid_to"),
				),
			}, {
				ResourceName:      tfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
