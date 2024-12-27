package feed

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func ResourceFeedRetentionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFeedRetentionPolicyCreate,
		ReadContext:   resourceFeedRetentionPolicyRead,
		UpdateContext: resourceFeedRetentionPolicyUpdate,
		DeleteContext: resourceFeedRetentionPolicyDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				ids := strings.Split(d.Id(), "/")
				if len(ids) == 1 {
					d.SetId(ids[0])
				} else {
					projectNameOrID, resourceID, err := tfhelper.ParseImportedName(d.Id())
					if err != nil {
						return nil, fmt.Errorf(" Parsing the resource ID. Expect in format `projectID/feedID`. Error: %v", err)
					}

					if projectNameOrID, err = tfhelper.GetRealProjectId(projectNameOrID, meta); err == nil {
						d.Set("project_id", projectNameOrID)
						d.SetId(resourceID)
					}

					if err != nil {
						return nil, err
					}
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"feed_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"count_limit": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 5000),
			},
			"days_to_keep_recently_downloaded_packages": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 4000),
			},
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceFeedRetentionPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	feedId := d.Get("feed_id").(string)
	projectId := d.Get("project_id").(string)
	_, err := clients.FeedClient.SetFeedRetentionPolicies(clients.Ctx, feed.SetFeedRetentionPoliciesArgs{
		Policy: &feed.FeedRetentionPolicy{
			CountLimit:                           converter.Int(d.Get("count_limit").(int)),
			DaysToKeepRecentlyDownloadedPackages: converter.Int(d.Get("days_to_keep_recently_downloaded_packages").(int)),
		},
		FeedId:  &feedId,
		Project: &projectId,
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf(" Creating Feed Retention Policy. FeedId: %s, Error: %+v", feedId, err))
	}

	d.SetId(feedId)
	return resourceFeedRetentionPolicyRead(clients.Ctx, d, m)
}

func resourceFeedRetentionPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	feedID := d.Id()
	projectId := d.Get("project_id").(string)
	policy, err := clients.FeedClient.GetFeedRetentionPolicies(clients.Ctx, feed.GetFeedRetentionPoliciesArgs{
		FeedId:  &feedID,
		Project: &projectId,
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf(" Failed get Feed Retention Policy. Projecct ID: %s , Feed ID: %s. Error: %+v", projectId, feedID, err))
	}

	if policy != nil {
		if policy.CountLimit != nil {
			d.Set("count_limit", policy.CountLimit)
		}

		if policy.DaysToKeepRecentlyDownloadedPackages != nil {
			d.Set("days_to_keep_recently_downloaded_packages", policy.DaysToKeepRecentlyDownloadedPackages)
		}
	}

	return nil
}

func resourceFeedRetentionPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	feedId := d.Get("feed_id").(string)
	projectId := d.Get("project_id").(string)
	_, err := clients.FeedClient.SetFeedRetentionPolicies(clients.Ctx, feed.SetFeedRetentionPoliciesArgs{
		Policy: &feed.FeedRetentionPolicy{
			CountLimit:                           converter.Int(d.Get("count_limit").(int)),
			DaysToKeepRecentlyDownloadedPackages: converter.Int(d.Get("days_to_keep_recently_downloaded_packages").(int)),
		},
		FeedId:  &feedId,
		Project: &projectId,
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf(" Updating Feed Retention Policy. ProjectID: %s, FeedId: %s, Error: %+v", projectId, feedId, err))
	}
	return resourceFeedRetentionPolicyRead(clients.Ctx, d, m)
}

func resourceFeedRetentionPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	feedId := d.Get("feed_id").(string)
	projectId := d.Get("project_id").(string)
	err := clients.FeedClient.DeleteFeedRetentionPolicies(clients.Ctx, feed.DeleteFeedRetentionPoliciesArgs{
		FeedId:  &feedId,
		Project: &projectId,
	})

	if err != nil {
		return diag.FromErr(fmt.Errorf(" Deleting Feed Retention Policy. ProjectID: %s, FeedId: %s, Error: %+v", projectId, feedId, err))
	}
	return nil
}
