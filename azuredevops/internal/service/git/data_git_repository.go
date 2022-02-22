package git

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

// DataGitRepository schema and implementation for Git repository data source
func DataGitRepository() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitRepositoryRead,
		Schema: map[string]*schema.Schema{
			"name": {
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
			"default_branch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_fork": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"remote_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ssh_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"web_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGitRepositoryRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	name := d.Get("name").(string)
	projectID := d.Get("project_id").(string)

	projectRepos, err := getGitRepositoriesByNameAndProject(clients, name, projectID, true)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return fmt.Errorf("Repository with name %s does not exist in project %s", name, projectID)
		}
		return fmt.Errorf("Error finding repositories. Error: %v", err)
	}
	if projectRepos == nil || 0 >= len(*projectRepos) {
		return fmt.Errorf("Repository with name %s does not exist in project %s", name, projectID)
	}
	if 1 < len(*projectRepos) {
		return fmt.Errorf("Multiple Repositories with name %s found in project %s", name, projectID)
	}

	err = flattenGitRepository(d, &(*projectRepos)[0])
	if err != nil {
		return fmt.Errorf("Error flattening Git repository: %w", err)
	}
	return nil
}
