package azuredevops

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/securityhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/suppress"
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
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validate.UUID,
				Required:     true,
				ForceNew:     true,
			},
			"principal": {
				Type:         schema.TypeString,
				ValidateFunc: validate.NoEmptyStrings,
				Required:     true,
				ForceNew:     true,
			},
			"replace": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"permissions": {
				// Unable to define a validation function, because the
				// keys and values can only be validated with an initialized
				// security client as we must load the security namespace
				// definition and the available permission settings, and a validation
				// function in Terraform only receives the parameter name and the
				// current value as argument
				Type:     schema.TypeMap,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: suppress.CaseDifference,
			},
		},
	}
}

func resourceProjectPermissionsCreate(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)
	projectID, ok := d.GetOk("project_id")
	if !ok {
		return fmt.Errorf("Failed to get 'project_id' from schema")
	}

	aclToken := fmt.Sprintf("$PROJECT:vstfs:///Classification/TeamProject/%s", projectID.(string))
	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.Project,
		clients.Ctx,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	principal, ok := d.GetOk("principal")
	if !ok {
		return fmt.Errorf("Failed to get 'principal' from schema")
	}

	permissions, ok := d.GetOk("permissions")
	if !ok {
		return fmt.Errorf("Failed to get 'permissions' from schema")
	}

	bReplace := d.Get("replace")
	permissionMap := make(map[securityhelper.ActionName]securityhelper.PermissionType, len(permissions.(map[string]interface{})))
	for key, elem := range permissions.(map[string]interface{}) {
		permissionMap[securityhelper.ActionName(key)] = securityhelper.PermissionType(elem.(string))
	}
	setPermissions := []securityhelper.SetPrincipalPermission{
		securityhelper.SetPrincipalPermission{
			Replace: bReplace.(bool),
			PrincipalPermission: securityhelper.PrincipalPermission{
				SubjectDescriptor: principal.(string),
				Permissions:       permissionMap,
			},
		}}

	err = sn.SetPrincipalPermissions(&setPermissions, &aclToken)
	if err != nil {
		return err
	}

	return resourceProjectPermissionsRead(d, m)
}

func resourceProjectPermissionsRead(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)
	projectID, ok := d.GetOk("project_id")
	if !ok {
		return fmt.Errorf("Failed to get 'project_id' from schema")
	}

	aclToken := fmt.Sprintf("$PROJECT:vstfs:///Classification/TeamProject/%s", projectID.(string))
	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.Project, clients.Ctx, clients.SecurityClient, clients.IdentityClient)
	if err != nil {
		return err
	}

	principal, ok := d.GetOk("principal")
	if !ok {
		return fmt.Errorf("Failed to get 'principal' from schema")
	}

	permissions, ok := d.GetOk("permissions")
	if !ok {
		return fmt.Errorf("Failed to get 'permissions' from schema")
	}

	principalList := []string{*converter.StringFromInterface(principal)}
	principalPermissions, err := sn.GetPrincipalPermissions(&aclToken, &principalList)
	if err != nil {
		return err
	}
	if principalPermissions == nil || len(*principalPermissions) != 1 {
		return fmt.Errorf("Failed to retrive current permissions for principal [%s]", principalList[0])
	}
	d.SetId(fmt.Sprintf("%s/%s", aclToken, principal.(string)))
	for key := range ((*principalPermissions)[0]).Permissions {
		if _, ok := permissions.(map[string]interface{})[string(key)]; !ok {
			delete(((*principalPermissions)[0]).Permissions, key)
		}
	}
	d.Set("permissions", ((*principalPermissions)[0]).Permissions)
	return nil
}

func resourceProjectPermissionsUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceProjectPermissionsCreate(d, m)
}

func resourceProjectPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	debugWait()
	// force all specified permissions to 'NotSet'
	return errors.New("resourceProjectPermissionsDelete: Not implemented")
}

func resourceProjectPermissionsImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// $PROJECT:vstfs:///Classification/TeamProject/<ProjectId>/<SubjectDescriptor>
	return nil, errors.New("resourceProjectPermissionsImporter: Not implemented")
}
