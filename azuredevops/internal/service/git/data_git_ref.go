package git

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// DataGitRef schema and implementation for git ref data source
func DataGitRef() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGitRefRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peeled_object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_locked_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGitRefRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	repoId := d.Get("repository_id").(string)
	name := d.Get("name").(string)

	args := git.GetRefsArgs{
		RepositoryId: &repoId,
		Filter:       &name,
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectID := v.(string)
		args.Project = &projectID
	}

	resp, err := clients.GitReposClient.GetRefs(ctx, args)
	if err != nil {
		return diag.Errorf("Error reading git ref: %v", err)
	}

	if resp == nil || len(resp.Value) == 0 {
		return diag.Errorf("Git ref not found for repository_id: %s, name: %s", repoId, name)
	}

	var match *git.GitRef
	for _, ref := range resp.Value {
		if ref.Name != nil && *ref.Name == name {
			match = &ref
			break
		}
	}

	if match == nil {
		// Just in case it wasn't an exact match
		match = &resp.Value[0]
	}

	d.SetId(fmt.Sprintf("%s:%s", repoId, name))
	if match.ObjectId != nil {
		d.Set("object_id", *match.ObjectId)
	}
	if match.PeeledObjectId != nil {
		d.Set("peeled_object_id", *match.PeeledObjectId)
	}
	if match.Creator != nil && match.Creator.Id != nil {
		d.Set("creator", *match.Creator.Id)
	}
	if match.Url != nil {
		d.Set("url", *match.Url)
	}
	if match.IsLocked != nil {
		d.Set("is_locked", *match.IsLocked)
	}
	if match.IsLockedBy != nil && match.IsLockedBy.Id != nil {
		d.Set("is_locked_by", *match.IsLockedBy.Id)
	}

	return nil
}
