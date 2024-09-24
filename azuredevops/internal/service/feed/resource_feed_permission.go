package feed

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

const (
	syncing = "Syncing"
	failed  = "Failed"
	succeed = "Succeeded"
)

func ResourceFeedPermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeedPermissionCreate,
		Read:   resourceFeedPermissionRead,
		Update: resourceFeedPermissionUpdate,
		Delete: resourceFeedPermissionDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
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

	feedId := d.Get("feed_id").(string)
	role := feed.FeedRole(d.Get("role").(string))
	projectId := d.Get("project_id").(string)
	displayName := d.Get("display_name").(string)
	identityDescriptor := d.Get("identity_descriptor").(string)

	permission, identityResponse, err := getFeedPermission(d, m)

	if err != nil && !utils.ResponseWasNotFound(err) {
		return fmt.Errorf("creating feed Permission for Feed : %s and Identity : %s, Error: %+v", feedId, identityDescriptor, err)
	}

	if permission != nil {
		return fmt.Errorf("feed Permission for Feed : %s and Identity : %s already exists", feedId, identityDescriptor)
	}

	_, err = clients.FeedClient.SetFeedPermissions(clients.Ctx, feed.SetFeedPermissionsArgs{
		FeedId:  &feedId,
		Project: &projectId,
		FeedPermission: &[]feed.FeedPermission{
			{
				DisplayName:        &displayName,
				IdentityDescriptor: identityResponse.Descriptor,
				IdentityId:         identityResponse.Id,
				Role:               &role,
			},
		},
	})

	if err != nil {
		return fmt.Errorf("creating feed Permission for Feed : %s and Identity : %s, Error: %+v", feedId, identityDescriptor, err)
	}

	err = checkPermissions(d, m)
	if err != nil {
		return fmt.Errorf(" Sync Feed Permission for Feed failed: %+v", err)
	}

	id, _ := uuid.NewUUID()
	d.SetId(fmt.Sprintf("fp-%s", id.String()))

	return resourceFeedPermissionRead(d, m)
}

func resourceFeedPermissionRead(d *schema.ResourceData, m interface{}) error {
	identityDescriptor := d.Get("identity_descriptor").(string)
	permission, identityResponse, err := getFeedPermission(d, m)
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
		d.Set("identity_descriptor", identityDescriptor)
		d.Set("identity_id", identityResponse.Id.String())
	}

	return nil
}

func resourceFeedPermissionUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	feedId := d.Get("feed_id").(string)
	identityDescriptor := d.Get("identity_descriptor").(string)
	role := feed.FeedRole(d.Get("role").(string))
	projectId := d.Get("project_id").(string)
	displayName := d.Get("display_name").(string)

	_, identityResponse, err := getFeedPermission(d, m)
	if err != nil {
		return fmt.Errorf("error reading feed permission during update: %+v", err)
	}

	_, err = clients.FeedClient.SetFeedPermissions(clients.Ctx, feed.SetFeedPermissionsArgs{
		FeedId:  &feedId,
		Project: &projectId,
		FeedPermission: &[]feed.FeedPermission{
			{
				DisplayName:        &displayName,
				IdentityDescriptor: identityResponse.Descriptor,
				IdentityId:         identityResponse.Id,
				Role:               &role,
			},
		},
	})

	if err != nil {
		return fmt.Errorf("updating feed Permission for Feed : %s and Identity : %s, Error: %+v", feedId, identityDescriptor, err)
	}

	err = checkPermissions(d, m)
	if err != nil {
		return fmt.Errorf(" Sync Feed Permission for Feed failed: %+v", err)
	}

	return resourceFeedPermissionRead(d, m)
}

func resourceFeedPermissionDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	feedId := d.Get("feed_id").(string)
	identityDescriptor := d.Get("identity_descriptor").(string)
	role := feed.FeedRoleValues.None
	projectId := d.Get("project_id").(string)

	identityResponse, err := getIdentity(d, m)

	if err != nil {
		return fmt.Errorf("deleting feed Permission for Feed : %s and Identity : %s, Error: %+v", feedId, identityDescriptor, err)
	}

	_, err = clients.FeedClient.SetFeedPermissions(clients.Ctx, feed.SetFeedPermissionsArgs{
		FeedId:  &feedId,
		Project: &projectId,
		FeedPermission: &[]feed.FeedPermission{
			{
				IdentityDescriptor: identityResponse.Descriptor,
				Role:               &role,
			},
		},
	})

	if err != nil {
		return fmt.Errorf("deleting feed Permission for Feed : %s and Identity : %s, Error: %+v", feedId, identityDescriptor, err)
	}

	return nil
}

func getIdentity(d *schema.ResourceData, m interface{}) (*identity.Identity, error) {
	clients := m.(*client.AggregatedClient)
	identityDescriptor := d.Get("identity_descriptor").(string)

	storageKey, err := clients.GraphClient.GetStorageKey(clients.Ctx, graph.GetStorageKeyArgs{
		SubjectDescriptor: &identityDescriptor,
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

	feedId := d.Get("feed_id").(string)
	identityDescriptor := d.Get("identity_descriptor").(string)
	projectId := d.Get("project_id").(string)

	identityResponse, err := getIdentity(d, m)

	if err != nil {
		return nil, nil, err
	}

	permissions, err := clients.FeedClient.GetFeedPermissions(clients.Ctx, feed.GetFeedPermissionsArgs{
		FeedId:  &feedId,
		Project: &projectId,
	})

	if err != nil {
		return nil, identityResponse, err
	}

	for _, permission := range *permissions {
		if *permission.IdentityDescriptor == *identityResponse.Descriptor {
			return &permission, identityResponse, nil
		}
	}

	message := fmt.Sprintf("error reading permission for Feed: %s and Identity: %s", feedId, identityDescriptor)
	return nil, identityResponse, azuredevops.WrappedError{
		StatusCode: converter.Int(http.StatusNotFound),
		Message:    &message,
	}
}

func checkPermissions(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	stateConf := &resource.StateChangeConf{
		ContinuousTargetOccurence: 2,
		Delay:                     5 * time.Second,
		MinTimeout:                10 * time.Second,
		Timeout:                   d.Timeout(schema.TimeoutCreate),
		Pending:                   []string{syncing},
		Target:                    []string{succeed, failed},
		Refresh:                   pollPermissions(d, m),
	}
	if _, err := stateConf.WaitForStateContext(clients.Ctx); err != nil {
		return fmt.Errorf(" Failed waiting for Feed Permission create. %v ", err)
	}
	return nil
}

func pollPermissions(d *schema.ResourceData, m interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		_, _, err := getFeedPermission(d, m)
		if err != nil {
			if utils.ResponseWasNotFound(err) {
				return nil, syncing, nil
			}
			return nil, failed, nil
		}
		return "", succeed, nil
	}
}
