package feed

import (
	"fmt"

	"github.com/google/uuid"
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

func resourceFeedCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.FeedClient)
	name := d.Get("name").(string)
	project := d.Get("project").(string)

	createFeed := feed.Feed{
		Name: &name,
	}

	err := clients.FeedClient.CreateFeed(clients.Ctx, &feed.CreateFeedArgs{
		Feed:    &createFeed,
		Project: &project,
	})

	if err != nil {
		return err
	}

	d.SetId("feed-" + uuid.New().String())

	return resourceFeedRead(d, m)
}

func resourceFeedRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.FeedClient)

	name := d.Get("name").(string)
	project := d.Get("project").(string)

	getFeed, err := clients.FeedClient.CreateFeed(clients.Ctx, &feed.GetFeedArgs{
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
		d.Set("name", *getFeed.Name)
		d.Set("project", *getFeed.Project.Name)
	}

	return nil
}

func resourceFeedUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.FeedClient)
	name := d.Get("name").(string)
	project := d.Get("project").(string)

	updateFeed := &feed.FeedUpdate{}

	err := clients.FeedClient.UpdateFeed(clients.Ctx, &feed.UpdateFeedArgs{
		Feed:    updateFeed,
		FeedId:  &name,
		Project: &project,
	})

	if err != nil {
		return err
	}

	d.SetId(d.Id())

	return resourceFeedRead(d, m)
}

func resourceFeedDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.FeedClient)
	name := d.Get("name").(string)
	project := d.Get("project").(string)

	err := clients.FeedClient.DeleteFeed(clients.Ctx, &feed.UpdateFeedArgs{
		FeedId:  &name,
		Project: &project,
	})

	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}
