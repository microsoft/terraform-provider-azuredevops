package azuredevops

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/commonhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/securityhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
)

func resourceAreaPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceAreaPermissionsCreate,
		Read:   resourceAreaPermissionsRead,
		Update: resourceAreaPermissionsUpdate,
		Delete: resourceAREAPermissionsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceAreaPermissionsImporter,
		},
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validate.UUID,
				Required:     true,
				ForceNew:     true,
			},
			"path": {
				Type:         schema.TypeString,
				ValidateFunc: validate.NoEmptyStrings,
				ForceNew:     true,
				Optional:     true,
			},
		}),
	}
}

func getAreaIDbyPath(clients *config.AggregatedClient, d *schema.ResourceData, path string) (*string, error) {
	var areaID string = ""
	projectID := d.Get("project_id").(string)

	args := workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		Path:           &path,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Depth:          converter.Int(999),
	}

	area, err := clients.WitClient.GetClassificationNode(clients.Ctx, args)
	if err != nil {
		return &areaID, fmt.Errorf("Error getting Area: %+v", err)
	}

	areaID = area.Identifier.String()
	return &areaID, nil
}

func createAreaToken(clients *config.AggregatedClient, d *schema.ResourceData) (*string, error) {
	var aclToken string
	var aclTokenPrefix string = "vstfs:///Classification/Node/"
	projectID := d.Get("project_id").(string)

	// you have to ommit the path property to get the
	// root area.
	args := workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Depth:          converter.Int(999),
	}

	rootArea, err := clients.WitClient.GetClassificationNode(clients.Ctx, args)
	if err != nil {
		return nil, fmt.Errorf("Error getting Area: %+v", err)
	}

	/*
	 * Token format
	 * Root area: vstfs:///Classification/Node/<AreaIdentifier>:vstfs:///Classification/Node/f8c5b667-91dd-4fe7-bf23-3138c439d07e"
	 * 1st child: vstfs:///Classification/Node/<AreaIdentifier>:vstfs:///Classification/Node/<AreaIdentifier>
	 */
	path, ok := d.GetOk("path")

	if !ok {
		// no path was specified we use the root area
		aclToken = "vstfs:///Classification/Node/" + rootArea.Identifier.String()
	} else {
		if !*rootArea.HasChildren {
			return &aclToken, fmt.Errorf("A path was specified but the root area has no children")
		} else {
			// get the id for each area in the provided path
			// we do this by querying each path element
			// 0: foo
			// 1: foo/bar
			// 3: foo/bar/baz
			aclToken = aclTokenPrefix + rootArea.Identifier.String()
			pathElem := strings.Split(path.(string), "/")
			for i, v := range pathElem {
				var pathQuery string
				if i == 0 {
					pathQuery = v
				} else {
					pathQuery = strings.Join(commonhelper.SelectArrayRange(pathElem, 0, i), "/")
				}
				currID, _ := getAreaIDbyPath(clients, d, pathQuery)
				aclToken = aclToken + ":" + aclTokenPrefix + *currID
			}
		}
	}

	log.Printf("[DEBUG] createAreaToken(): Discovered aclToken %q", aclToken)
	return &aclToken, nil
}

func resourceAreaPermissionsCreate(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.CSS,
		clients.Ctx,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	aclToken, err := createAreaToken(clients, d)
	if err != nil {
		return err
	}

	err = securityhelper.SetPrincipalPermissions(d, sn, aclToken, nil, false)
	if err != nil {
		return err
	}

	return resourceAreaPermissionsRead(d, m)
}

func resourceAreaPermissionsRead(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)

	aclToken, err := createAreaToken(clients, d)
	if err != nil {
		return err
	}

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.CSS,
		clients.Ctx,
		clients.SecurityClient,
		clients.IdentityClient)
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

func resourceAreaPermissionsUpdate(d *schema.ResourceData, m interface{}) error {
	debugWait()

	return resourceAreaPermissionsCreate(d, m)
}

func resourceAREAPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)

	aclToken, err := createAreaToken(clients, d)
	if err != nil {
		return err
	}

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.CSS,
		clients.Ctx,
		clients.SecurityClient,
		clients.IdentityClient)
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

func resourceAreaPermissionsImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	debugWait()

	// repoV2/#ProjectID#/#RepositoryID#/refs/heads/#BranchName#/#SubjectDescriptor#
	return nil, errors.New("resourceAreaPermissionsImporter: Not implemented")
}
