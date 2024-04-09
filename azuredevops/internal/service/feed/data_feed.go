package feed

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

func DataFeed() *schema.Resource {
	return &schema.Resource{
		Read: dataFeedRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Optional:     true,
				AtLeastOneOf: []string{
					"name", "feed_id",
				},
			},
			"feed_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				ConflictsWith: []string{
					"name",
				},
			},
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Optional:     true,
			},
		},
	}
}

func dataFeedRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	name := d.Get("name").(string)
	id := d.Get("feed_id").(string)
	projectId := d.Get("project_id").(string)

	identifier := id
	if identifier == "" {
		identifier = name
	}

	getFeed, err := clients.FeedClient.GetFeed(clients.Ctx, feed.GetFeedArgs{
		FeedId:  &identifier,
		Project: &projectId,
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading feed during read: %+v", err)
	}

	if getFeed != nil {
		d.SetId((*getFeed.Id).String())
		d.Set("name", *getFeed.Name)
		d.Set("feed_id", (*getFeed.Id).String())
		if getFeed.Project != nil {
			d.Set("project_id", (*getFeed.Project.Id).String())
		}
	}

	return nil
}
