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
				ValidateFunc: validation.NoZeroValues,
				Required:     true,
			},
			"project": {
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
	project := d.Get("project").(string)

	getFeed, err := clients.FeedClient.GetFeed(clients.Ctx, feed.GetFeedArgs{
		FeedId:  &name,
		Project: &project,
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading feed during read: %+v", err)
	}

	if getFeed != nil {
		d.SetId((*getFeed).Id.String())
		d.Set("name", (*getFeed).Name)
		project := (*getFeed).Project
		if project != nil {
			d.Set("project", (*project).Name)
		}
	}

	return nil
}
