package feed

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/feed"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
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
				Type:     schema.TypeString,
				Computed: true,
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
	role := feed.FeedRole(d.Get("role").(string))
	project_id := d.Get("project_id").(string)
	display_name := d.Get("display_name").(string)
	identity_descriptor := d.Get("identity_descriptor").(string)

	permission, identity_response, err := getFeedPermission(d, m)

	if err != nil && !utils.ResponseWasNotFound(err) {
		return fmt.Errorf("creating feed Permission for Feed : %s and Identity : %s, Error: %+v", feed_id, identity_descriptor, err)
	}

	if permission != nil {
		return fmt.Errorf("feed Permission for Feed : %s and Identity : %s already exists", feed_id, identity_descriptor)
	}

	_, err = clients.FeedClient.SetFeedPermissions(clients.Ctx, feed.SetFeedPermissionsArgs{
		FeedId:  &feed_id,
		Project: &project_id,
		FeedPermission: &[]feed.FeedPermission{
			{
				DisplayName:        &display_name,
				IdentityDescriptor: identity_response.Descriptor,
				IdentityId:         identity_response.Id,
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
	identity_descriptor := d.Get("identity_descriptor").(string)
	permission, identity_response, err := getFeedPermission(d, m)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading feed permission during read: %+v", err)
	}

	if permission != nil {
		if permission.DisplayName != nil {
			d.Set("display_name", *permission.DisplayName)
		}
		d.Set("role", *permission.Role)
		d.Set("identity_descriptor", identity_descriptor)
		d.Set("identity_id", identity_response.Id.String())
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

	_, identity_response, err := getFeedPermission(d, m)
	if err != nil {
		return fmt.Errorf("error reading feed permission during update: %+v", err)
	}

	_, err = clients.FeedClient.SetFeedPermissions(clients.Ctx, feed.SetFeedPermissionsArgs{
		FeedId:  &feed_id,
		Project: &project_id,
		FeedPermission: &[]feed.FeedPermission{
			{
				DisplayName:        &display_name,
				IdentityDescriptor: identity_response.Descriptor,
				IdentityId:         identity_response.Id,
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

	identity_response, err := getIdentity(d, m)

	if err != nil {
		return fmt.Errorf("deleting feed Permission for Feed : %s and Identity : %s, Error: %+v", feed_id, identity_descriptor, err)
	}

	_, err = clients.FeedClient.SetFeedPermissions(clients.Ctx, feed.SetFeedPermissionsArgs{
		FeedId:  &feed_id,
		Project: &project_id,
		FeedPermission: &[]feed.FeedPermission{
			{
				IdentityDescriptor: identity_response.Descriptor,
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

func getIdentity(d *schema.ResourceData, m interface{}) (*identity.Identity, error) {
	clients := m.(*client.AggregatedClient)
	identity_descriptor := d.Get("identity_descriptor").(string)

	storageKey, err := clients.GraphClient.GetStorageKey(clients.Ctx, graph.GetStorageKeyArgs{
		SubjectDescriptor: &identity_descriptor,
	})

	if err != nil {
		return nil, err
	}

	response, err := clients.IdentityClient.ReadIdentity(clients.Ctx, identity.ReadIdentityArgs{
		IdentityId: converter.String((*storageKey.Value).String()),
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func getFeedPermission(d *schema.ResourceData, m interface{}) (*feed.FeedPermission, *identity.Identity, error) {
	clients := m.(*client.AggregatedClient)

	feed_id := d.Get("feed_id").(string)
	identity_descriptor := d.Get("identity_descriptor").(string)
	project_id := d.Get("project_id").(string)

	identity_response, err := getIdentity(d, m)

	if err != nil {
		return nil, nil, err
	}

	permissions, err := clients.FeedClient.GetFeedPermissions(clients.Ctx, feed.GetFeedPermissionsArgs{
		FeedId:  &feed_id,
		Project: &project_id,
	})

	if err != nil {
		return nil, identity_response, err
	}

	for _, permission := range *permissions {
		if *permission.IdentityDescriptor == *identity_response.Descriptor {
			return &permission, identity_response, nil
		}
	}

	notFound := http.StatusNotFound
	message := fmt.Sprintf("error reading permission for Feed: %s and Identity: %s", feed_id, identity_descriptor)
	return nil, identity_response, azuredevops.WrappedError{
		StatusCode: &notFound,
		Message:    &message,
	}
}
