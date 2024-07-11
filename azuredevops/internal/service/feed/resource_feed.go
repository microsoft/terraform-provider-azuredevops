package feed

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/webapi"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
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
							Default:  false,
						},
						"restore": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func resourceFeedCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	name := d.Get("name").(string)
	projectId := d.Get("project_id").(string)
	features := expandFeedFeatures(d.Get("features").([]interface{}))

	if v, ok := features["restore"]; ok && v.(bool) {
		if isFeedRestorable(d, m) {
			err := restoreFeed(d, m)
			if err != nil {
				return fmt.Errorf("restoring feed. Name: %s, Error: %+v", name, err)
			}
			return resourceFeedRead(d, m)
		}
	}

	feedDetail, err := clients.FeedClient.CreateFeed(clients.Ctx, feed.CreateFeedArgs{
		Feed: &feed.Feed{
			Name: &name,
		},
		Project: &projectId,
	})

	if err != nil {
		return fmt.Errorf("creating new feed. Name: %s, Error: %+v", name, err)
	}

	d.SetId(feedDetail.Id.String())
	return resourceFeedRead(d, m)
}

func resourceFeedRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	name := d.Get("name").(string)
	projectId := d.Get("project_id").(string)

	feedDetail, err := clients.FeedClient.GetFeed(clients.Ctx, feed.GetFeedArgs{
		FeedId:  &name,
		Project: &projectId,
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(" failed get feed. Projecct ID: %s , Feed Name %s : . Error: %+v", projectId, name, err)
	}

	if feedDetail != nil {
		d.SetId(feedDetail.Id.String())
		d.Set("name", feedDetail.Name)
		if feedDetail.Project != nil {
			d.Set("project_id", feedDetail.Project.Id.String())
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
	features := expandFeedFeatures(d.Get("features").([]interface{}))

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
	return nil
}

func isFeedRestorable(d *schema.ResourceData, m interface{}) bool {
	clients := m.(*client.AggregatedClient)

	change, err := clients.FeedClient.GetFeedChange(clients.Ctx, feed.GetFeedChangeArgs{
		FeedId:  converter.String(d.Get("name").(string)),
		Project: converter.String(d.Get("project_id").(string)),
	})

	return err == nil && *((*change).ChangeType) == feed.ChangeTypeValues.Delete
}

func restoreFeed(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	err := clients.FeedClient.RestoreDeletedFeed(clients.Ctx, feed.RestoreDeletedFeedArgs{
		FeedId:  converter.String(d.Get("name").(string)),
		Project: converter.String(d.Get("project_id").(string)),
		PatchJson: &[]webapi.JsonPatchOperation{{
			From:  nil,
			Path:  converter.String("/isDeleted"),
			Op:    &webapi.OperationValues.Replace,
			Value: false,
		}},
	})

	if err != nil {
		return err
	}

	return nil
}

func expandFeedFeatures(input []interface{}) map[string]interface{} {
	if len(input) == 0 || input[0] == nil {
		return map[string]interface{}{}
	}
	return input[0].(map[string]interface{})
}
