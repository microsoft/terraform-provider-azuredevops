package azuredevops

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/securityhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
)

func resourceProjectPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectPermissionsCreate,
		Read:   resourceProjectPermissionsRead,
		Update: resourceProjectPermissionsUpdate,
		Delete: resourceProjectPermissionsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceProjectPermissionsImporter,
		},
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validate.UUID,
				Required:     true,
				ForceNew:     true,
			},
		}),
	}
}

func createProjectToken(d *schema.ResourceData) (*string, error) {
	projectID, ok := d.GetOk("project_id")
	if !ok {
		return nil, fmt.Errorf("Failed to get 'project_id' from schema")
	}
	aclToken := fmt.Sprintf("$PROJECT:vstfs:///Classification/TeamProject/%s", projectID.(string))
	return &aclToken, nil
}

func resourceProjectPermissionsCreate(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.Project,
		clients.Ctx,
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
	debugWait()

	clients := m.(*config.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.Project, clients.Ctx, clients.SecurityClient, clients.IdentityClient)
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
	debugWait()

	return resourceProjectPermissionsCreate(d, m)
}

func resourceProjectPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.Project, clients.Ctx, clients.SecurityClient, clients.IdentityClient)
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

func resourceProjectPermissionsImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	debugWait()

	// $PROJECT:vstfs:///Classification/TeamProject/<ProjectId>/<SubjectDescriptor>
	return nil, errors.New("resourceProjectPermissionsImporter: Not implemented")
}
