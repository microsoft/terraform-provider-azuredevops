//go:build (all || data_sources || data_agent_queue) && (!exclude_data_sources || !exclude_data_agent_queue)
// +build all data_sources data_agent_queue
// +build !exclude_data_sources !exclude_data_agent_queue

package acceptancetests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccPersonAccessToken_DataSource(t *testing.T) {
	authorization_id := uuid.New().String()
	personalAccessTokenData := testutils.HclPersonalAccessTokenDataSource(authorization_id)

	tfNode := "data.azuredevops_personal_access_token.pat"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: personalAccessTokenData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "id", authorization_id),
					resource.TestCheckResourceAttr(tfNode, "authorization_id", authorization_id),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "scope"),
					resource.TestCheckResourceAttrSet(tfNode, "target_accounts"),
					resource.TestCheckResourceAttrSet(tfNode, "token"),
					resource.TestCheckResourceAttrSet(tfNode, "valid_from"),
					resource.TestCheckResourceAttrSet(tfNode, "valid_to"),
					resource.TestCheckResourceAttrSet(tfNode, "valid_to"),
				),
			},
		},
	})
}
