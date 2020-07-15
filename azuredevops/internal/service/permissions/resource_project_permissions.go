package permissions

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
)

// ResourceProjectPermissions schema and implementation for project permission resource
func ResourceProjectPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectPermissionsCreate,
		Read:   resourceProjectPermissionsRead,
		Update: resourceProjectPermissionsUpdate,
		Delete: resourceProjectPermissionsDelete,
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
			},
		}),
	}
}

func resourceProjectPermissionsCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.Project,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	aclToken, err := createProjectToken(d)
	if err != nil {
		return err
	}

	err = securityhelper.SetPrincipalPermissions(d, sn, aclToken, nil, false)
	if err != nil {
		return err
	}

	return resourceProjectPermissionsRead(d, m)
}

func resourceProjectPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.Project,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	aclToken, err := createProjectToken(d)
	if err != nil {
		return err
	}

	principalPermissions, err := securityhelper.GetPrincipalPermissions(d, sn, aclToken)
	if err != nil {
		return err
	}

	d.Set("permissions", principalPermissions.Permissions)
	return nil
}

func resourceProjectPermissionsUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceProjectPermissionsCreate(d, m)
}

func resourceProjectPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.Project,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	aclToken, err := createProjectToken(d)
	if err != nil {
		return err
	}

	err = securityhelper.SetPrincipalPermissions(d, sn, aclToken, &securityhelper.PermissionTypeValues.NotSet, true)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createProjectToken(d *schema.ResourceData) (*string, error) {
	projectID, ok := d.GetOk("project_id")
	if !ok {
		return nil, fmt.Errorf("Failed to get 'project_id' from schema")
	}
	aclToken := fmt.Sprintf("$PROJECT:vstfs:///Classification/TeamProject/%s", projectID.(string))
	return &aclToken, nil
}
