package feed

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

func DataFeed() *schema.Resource {
	return &schema.Resource{
		Read: dataFeedRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
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
				ValidateFunc: validation.IsUUID,
				ConflictsWith: []string{
					"name",
				},
			},
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
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
			return nil
		}
		return fmt.Errorf("reading feed during read: %+v", err)
	}

	if getFeed != nil {
		d.SetId(getFeed.Id.String())
		d.Set("name", *getFeed.Name)
		d.Set("feed_id", getFeed.Id.String())
		if getFeed.Project != nil {
			d.Set("project_id", getFeed.Project.Id.String())
		}
	}

	return nil
}
