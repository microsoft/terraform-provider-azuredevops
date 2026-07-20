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

// DataGitRefs schema and implementation for git refs data source
func DataGitRefs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGitRefsRead,
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
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_contains": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"refs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceGitRefsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	repoId := d.Get("repository_id").(string)

	args := git.GetRefsArgs{
		RepositoryId: &repoId,
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectID := v.(string)
		args.Project = &projectID
	}

	if v, ok := d.GetOk("filter"); ok {
		f := v.(string)
		args.Filter = &f
	}

	if v, ok := d.GetOk("filter_contains"); ok {
		fc := v.(string)
		args.FilterContains = &fc
	}

	var allRefs []git.GitRef
	for {
		resp, err := clients.GitReposClient.GetRefs(ctx, args)
		if err != nil {
			return diag.Errorf("Error reading git refs: %v", err)
		}
		if resp == nil {
			break
		}

		allRefs = append(allRefs, resp.Value...)

		if resp.ContinuationToken != "" {
			args.ContinuationToken = &resp.ContinuationToken
		} else {
			break
		}
	}

	var results []interface{}
	for _, ref := range allRefs {
		m := make(map[string]interface{})
		if ref.Name != nil {
			m["name"] = *ref.Name
		}
		if ref.ObjectId != nil {
			m["object_id"] = *ref.ObjectId
		}
		if ref.PeeledObjectId != nil {
			m["peeled_object_id"] = *ref.PeeledObjectId
		}
		if ref.Creator != nil && ref.Creator.Id != nil {
			m["creator"] = *ref.Creator.Id
		}
		if ref.Url != nil {
			m["url"] = *ref.Url
		}
		if ref.IsLocked != nil {
			m["is_locked"] = *ref.IsLocked
		}
		if ref.IsLockedBy != nil && ref.IsLockedBy.Id != nil {
			m["is_locked_by"] = *ref.IsLockedBy.Id
		}
		results = append(results, m)
	}

	d.SetId(repoId)
	if err := d.Set("refs", results); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting refs: %v", err))
	}

	return nil
}
