package feed

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"

	feedUtils "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/feed/utils"
)

func DataFeeds() *schema.Resource {
	feedResourceSchema := feedUtils.CommonFeedFields()

	feedResourceSchema["name"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	feedResourceSchema["feed_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	feedResourceSchema["project_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return &schema.Resource{
		Read: dataFeedsRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},
			"feeds": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: feedResourceSchema,
				},
			},
		},
	}
}

func dataFeedsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)

	args := feed.GetFeedsArgs{}
	if projectID != "" {
		args.Project = converter.String(projectID)
	}

	response, err := clients.FeedClient.GetFeeds(clients.Ctx, args)
	if err != nil {
		return fmt.Errorf("Error finding feeds: %+v", err)
	}

	var feedsData []any

	if response != nil {
		for _, feedObj := range *response {
			item := make(map[string]any)

			if feedObj.Id != nil {
				item["feed_id"] = feedObj.Id.String()
			}
			if feedObj.Name != nil {
				item["name"] = *feedObj.Name
			}
			if feedObj.Project != nil {
				item["project_id"] = feedObj.Project.Id.String()
			}
			if feedObj.Description != nil {
				item["description"] = *feedObj.Description
			}
			if feedObj.Url != nil {
				item["url"] = *feedObj.Url
			}
			if feedObj.BadgesEnabled != nil {
				item["badges_enabled"] = *feedObj.BadgesEnabled
			}
			if feedObj.HideDeletedPackageVersions != nil {
				item["hide_deleted_package_versions"] = *feedObj.HideDeletedPackageVersions
			}
			if feedObj.UpstreamEnabled != nil {
				item["upstream_enabled"] = *feedObj.UpstreamEnabled
			}

			if feedObj.UpstreamSources != nil {
				item["upstream_sources"] = feedUtils.FlattenUpstreamSources(*feedObj.UpstreamSources)
			}

			feedsData = append(feedsData, item)
		}
	}

	d.Set("feeds", feedsData)

	if projectID != "" {
		d.SetId(fmt.Sprintf("feeds-%s", projectID))
	} else {
		d.SetId("feeds-all")
	}

	return nil
}
