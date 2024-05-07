package feed

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
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
				ForceNew:     true,
			},
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Optional:     true,
				ForceNew:     true,
			},
			"features": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permanent_delete": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"restore": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"restored": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

var FeatureDefaults = map[string]interface{}{
	"permanent_delete": true,
	"restore":          true,
}

func resourceFeedCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	features := buildDefinitionFeatures(d)

	if v, ok := features["restore"]; ok {
		if restore := v.(bool); restore && isFeedRestorable(d, m) {
			err := restoreFeed(d, m)

			if err != nil {
				return fmt.Errorf("restoring feed. Name: %s, Error: %+v", name, err)
			}

			return resourceFeedRead(d, m)
		}
	}

	err := createFeed(d, m)

	if err != nil {
		return fmt.Errorf("creating new feed. Name: %s, Error: %+v", name, err)
	}

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
		return fmt.Errorf("error reading feed during read: %+v", err)
	}

	if getFeed != nil {
		d.SetId(getFeed.Id.String())
		d.Set("name", getFeed.Name)
		if getFeed.Project != nil {
			d.Set("project_id", getFeed.Project.Id.String())
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
	projectId := d.Get("project_id").(string)
	features := buildDefinitionFeatures(d)

	err := clients.FeedClient.DeleteFeed(clients.Ctx, feed.DeleteFeedArgs{
		FeedId:  &name,
		Project: &projectId,
	})

	if err != nil {
		return err
	}

	if v, ok := features["permanent_delete"]; ok {
		if permanentDelete := v.(bool); permanentDelete {
			err = clients.FeedClient.PermanentDeleteFeed(clients.Ctx, feed.PermanentDeleteFeedArgs{
				FeedId:  &name,
				Project: &projectId,
			})

			if err != nil {
				return err
			}
		}
	}

	d.SetId("")

	return nil
}

func isFeedRestorable(d *schema.ResourceData, m interface{}) bool {
	clients := m.(*client.AggregatedClient)
	name := d.Get("name").(string)
	projectId := d.Get("project_id").(string)

	change, err := clients.FeedClient.GetFeedChange(clients.Ctx, feed.GetFeedChangeArgs{
		FeedId:  &name,
		Project: &projectId,
	})

	return err == nil && *((*change).ChangeType) == feed.ChangeTypeValues.Delete
}

func createFeed(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	name := d.Get("name").(string)
	projectId := d.Get("project_id").(string)

	createFeed := feed.Feed{
		Name: &name,
	}

	_, err := clients.FeedClient.CreateFeed(clients.Ctx, feed.CreateFeedArgs{
		Feed:    &createFeed,
		Project: &projectId,
	})

	if err != nil {
		return err
	}

	d.Set("restored", false)

	return nil
}

func restoreFeed(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	name := d.Get("name").(string)
	projectId := d.Get("project_id").(string)

	path := "/isDeleted"

	patchJsons := []webapi.JsonPatchOperation{{
		From:  nil,
		Path:  &path,
		Op:    &webapi.OperationValues.Replace,
		Value: false,
	}}

	err := clients.FeedClient.RestoreDeletedFeed(clients.Ctx, feed.RestoreDeletedFeedArgs{
		FeedId:    &name,
		Project:   &projectId,
		PatchJson: &patchJsons,
	})

	if err != nil {
		return err
	}

	d.Set("restored", true)

	return nil
}

func buildDefinitionFeatures(d *schema.ResourceData) map[string]interface{} {
	features := d.Get("features").([]interface{})
	if len(features) != 0 {
		featureMap := features[0].(map[string]interface{})
		for k, v := range FeatureDefaults {
			if _, ok := featureMap[k]; !ok {
				featureMap[k] = v
			}
		}
		return featureMap
	}
	return FeatureDefaults
}
