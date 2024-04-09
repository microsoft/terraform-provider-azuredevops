package feed

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

func ResourceFeed() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeedCreate,
		Read:   resourceFeedRead,
		Update: resourceFeedUpdate,
		Delete: resourceFeedDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Required:     true,
			},
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,

				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceFeedCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	name := d.Get("name").(string)
	projectId := d.Get("project_id").(string)

	createFeed := feed.Feed{
		Name: &name,
	}

	createdFeed, err := clients.FeedClient.CreateFeed(clients.Ctx, feed.CreateFeedArgs{
		Feed:    &createFeed,
		Project: &projectId,
	})

	if err != nil {
		return err
	}

	d.SetId((*createdFeed).Id.String())

	return resourceFeedRead(d, m)
}

func resourceFeedRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	name := d.Get("name").(string)
	projectId := d.Get("project_id").(string)

	getFeed, err := clients.FeedClient.GetFeed(clients.Ctx, feed.GetFeedArgs{
		FeedId:  &name,
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
		d.SetId((*getFeed).Id.String())
		d.Set("name", (*getFeed).Name)
		project := (*getFeed).Project
		if project != nil {
			d.Set("project_id", (*project).Id.String())
		}
	}

	return nil
}

func resourceFeedUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	name := d.Get("name").(string)
	projectId := d.Get("project_id").(string)

	_, err := clients.FeedClient.UpdateFeed(clients.Ctx, feed.UpdateFeedArgs{
		Feed:    &feed.FeedUpdate{},
		FeedId:  &name,
		Project: &projectId,
	})

	if err != nil {
		return err
	}

	return resourceFeedRead(d, m)
}

func resourceFeedDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	name := d.Get("name").(string)
	project_id := d.Get("project_id").(string)

	err := clients.FeedClient.DeleteFeed(clients.Ctx, feed.DeleteFeedArgs{
		FeedId:  &name,
		Project: &project_id,
	})

	if err != nil {
		return err
	}

	err = clients.FeedClient.PermanentDeleteFeed(clients.Ctx, feed.PermanentDeleteFeedArgs{
		FeedId:  &name,
		Project: &project_id,
	})

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
