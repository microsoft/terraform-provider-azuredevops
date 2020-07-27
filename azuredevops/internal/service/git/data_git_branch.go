package git

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

func DataGitBranch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitBranchRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsNotWhiteSpace,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"repo_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsNotWhiteSpace,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"project_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.IsUUID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGitBranchRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	name := d.Get("name").(string)
	repoName := d.Get("repo_name").(string)
	projectID := d.Get("project_id").(string)

	repoBranch, err := getGitBranchByNameAndRepo(clients, name, repoName, projectID)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return fmt.Errorf("Branch with name %s does not exist in repo %s project %s", name, repoName, projectID)
		}
		return fmt.Errorf("Error finding branches. Error: %v", err)
	}
	if repoBranch == nil {
		return fmt.Errorf("Branch with name %s does not exist in repo %s project %s", name, repoName, projectID)
	}

	err = flattenGitBranch(d, repoBranch)
	if err != nil {
		return fmt.Errorf("Error flattening Git branch: %w", err)
	}
	return nil
}

func getGitBranchByNameAndRepo(clients *client.AggregatedClient, name string, repoName string, projectID string) (*git.GitRef, error) {
	branch, err := gitBranchRead(clients, name, repoName, projectID)
	return branch, err
}
