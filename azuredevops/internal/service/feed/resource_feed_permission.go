package feed

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

func ResourceFeedPermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeedPermissionCreate,
		Read:   resourceFeedPermissionRead,
		Update: resourceFeedPermissionUpdate,
		Delete: resourceFeedPermissionDelete,
		Schema: map[string]*schema.Schema{
			"feed_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
			},
			"identity_descriptor": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Required:     true,
				ForceNew:     true,
			},
			"identity_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Computed:     true,
			},
			"role": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					string(feed.FeedRoleValues.Reader),
					string(feed.FeedRoleValues.Contributor),
					string(feed.FeedRoleValues.Administrator),
					string(feed.FeedRoleValues.Collaborator),
				}, false),
				Required: true,
			},
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Optional:     true,
				ForceNew:     true,
			},
			"display_name": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Optional:     true,
			},
		},
	}
}

func resourceFeedPermissionCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	feed_id := d.Get("feed_id").(string)
	identity_descriptor := d.Get("identity_descriptor").(string)
	role := feed.FeedRole(d.Get("role").(string))
	project_id := d.Get("project_id").(string)
	display_name := d.Get("display_name").(string)

	permission, getFeedPermissionErr := getFeedPermission(d, m)

	if getFeedPermissionErr != nil && !utils.ResponseWasNotFound(getFeedPermissionErr) {
		return fmt.Errorf("creating feed Permission for Feed : %s and Identity : %s, Error: %+v", feed_id, identity_descriptor, getFeedPermissionErr)
	}

	if permission != nil {
		return fmt.Errorf("feed Permission for Feed : %s and Identity : %s already exists", feed_id, identity_descriptor)
	}

	_, err := clients.FeedClient.SetFeedPermissions(clients.Ctx, feed.SetFeedPermissionsArgs{
		FeedId:  &feed_id,
		Project: &project_id,
		FeedPermission: &[]feed.FeedPermission{
			{
				DisplayName:        &display_name,
				IdentityDescriptor: &identity_descriptor,
				Role:               &role,
			},
		},
	})

	if err != nil {
		return fmt.Errorf("creating feed Permission for Feed : %s and Identity : %s, Error: %+v", feed_id, identity_descriptor, err)
	}

	id, _ := uuid.NewUUID()
	d.SetId(fmt.Sprintf("fp-%s", id.String()))

	return resourceFeedPermissionRead(d, m)
}

func resourceFeedPermissionRead(d *schema.ResourceData, m interface{}) error {
	permission, err := getFeedPermission(d, m)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading feed permission during read: %+v", err)
	}

	if permission != nil {
		d.Set("display_name", *permission.DisplayName)
		d.Set("role", *permission.Role)
		d.Set("identity_descriptor", *permission.IdentityDescriptor)
		d.Set("identity_id", (*permission.IdentityId).String())
	}

	return nil
}

func resourceFeedPermissionUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	feed_id := d.Get("feed_id").(string)
	identity_descriptor := d.Get("identity_descriptor").(string)
	role := feed.FeedRole(d.Get("role").(string))
	project_id := d.Get("project_id").(string)
	display_name := d.Get("display_name").(string)

	_, getFeedPermissionErr := getFeedPermission(d, m)
	if getFeedPermissionErr != nil {
		return fmt.Errorf("error reading feed permission during update: %+v", getFeedPermissionErr)
	}

	_, err := clients.FeedClient.SetFeedPermissions(clients.Ctx, feed.SetFeedPermissionsArgs{
		FeedId:  &feed_id,
		Project: &project_id,
		FeedPermission: &[]feed.FeedPermission{
			{
				DisplayName:        &display_name,
				IdentityDescriptor: &identity_descriptor,
				Role:               &role,
			},
		},
	})

	if err != nil {
		return fmt.Errorf("updating feed Permission for Feed : %s and Identity : %s, Error: %+v", feed_id, identity_descriptor, err)
	}

	return resourceFeedPermissionRead(d, m)
}

func resourceFeedPermissionDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	feed_id := d.Get("feed_id").(string)
	identity_descriptor := d.Get("identity_descriptor").(string)
	role := feed.FeedRoleValues.None
	project_id := d.Get("project_id").(string)

	_, err := clients.FeedClient.SetFeedPermissions(clients.Ctx, feed.SetFeedPermissionsArgs{
		FeedId:  &feed_id,
		Project: &project_id,
		FeedPermission: &[]feed.FeedPermission{
			{
				IdentityDescriptor: &identity_descriptor,
				Role:               &role,
			},
		},
	})

	if err != nil {
		return fmt.Errorf("deleting feed Permission for Feed : %s and Identity : %s, Error: %+v", feed_id, identity_descriptor, err)
	}

	d.SetId("")
	return nil
}

func getFeedPermission(d *schema.ResourceData, m interface{}) (*feed.FeedPermission, error) {
	clients := m.(*client.AggregatedClient)

	feed_id := d.Get("feed_id").(string)
	identity_descriptor := d.Get("identity_descriptor").(string)
	project_id := d.Get("project_id").(string)

	permissions, err := clients.FeedClient.GetFeedPermissions(clients.Ctx, feed.GetFeedPermissionsArgs{
		FeedId:             &feed_id,
		Project:            &project_id,
		IdentityDescriptor: &identity_descriptor,
	})

	if err != nil {
		return nil, err
	}

	for _, permission := range *permissions {
		return &permission, nil
	}

	notFound := http.StatusNotFound
	message := fmt.Sprintf("error reading permission for Feed: %s and Identity: %s", feed_id, identity_descriptor)
	return nil, azuredevops.WrappedError{
		StatusCode: &notFound,
		Message:    &message,
	}
}
