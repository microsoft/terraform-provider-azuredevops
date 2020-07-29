package permissions

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceAreaPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceAreaPermissionsCreate,
		Read:   resourceAreaPermissionsRead,
		Update: resourceAreaPermissionsUpdate,
		Delete: resourceAREAPermissionsDelete,
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
			},
			"path": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				ForceNew:     true,
				Optional:     true,
			},
		}),
	}
}

func getAreaIDbyPath(context context.Context, workitemtrackingClient workitemtracking.Client, projectID string, path string) (*string, error) {
	area, err := workitemtrackingClient.GetClassificationNode(context, workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		Path:           &path,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Depth:          converter.Int(1),
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting Area: %w", err)
	}

	areaID := area.Identifier.String()
	return &areaID, nil
}

func createAreaToken(context context.Context, workitemtrackingClient workitemtracking.Client, d *schema.ResourceData) (*string, error) {
	const aclTokenPrefix = "vstfs:///Classification/Node/"

	var aclToken string
	projectID := d.Get("project_id").(string)

	// you have to ommit the path property to get the
	// root area.
	rootArea, err := workitemtrackingClient.GetClassificationNode(context, workitemtracking.GetClassificationNodeArgs{
		Project:        converter.String(projectID),
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Areas,
		Depth:          converter.Int(1),
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting Area: %+v", err)
	}

	/*
	 * Token format
	 * Root area: vstfs:///Classification/Node/<AreaIdentifier>:vstfs:///Classification/Node/f8c5b667-91dd-4fe7-bf23-3138c439d07e"
	 * 1st child: vstfs:///Classification/Node/<AreaIdentifier>:vstfs:///Classification/Node/<AreaIdentifier>
	 */
	aclToken = aclTokenPrefix + rootArea.Identifier.String()

	path, ok := d.GetOk("path")
	if ok {
		if !*rootArea.HasChildren {
			return nil, fmt.Errorf("A path was specified but the root area has no children")
		} else {
			// get the id for each area in the provided path
			// we do this by querying each path element
			// 0: foo
			// 1: foo/bar
			// 3: foo/bar/baz
			var pathElem []string

			linq.From(strings.Split(path.(string), "/")).
				Where(func(elem interface{}) bool {
					return len(elem.(string)) > 0
				}).
				ToSlice(&pathElem)

			for i := range pathElem {
				pathItem := strings.Join(pathElem[:i+1], "/")
				currID, err := getAreaIDbyPath(context, workitemtrackingClient, projectID, pathItem)
				if err != nil {
					return nil, fmt.Errorf("Failed to get ID for area %s, %w", pathItem, err)
				}
				aclToken = aclToken + ":" + aclTokenPrefix + *currID
			}
		}
	}

	log.Printf("[DEBUG] createAreaToken(): Discovered aclToken %q", aclToken)
	return &aclToken, nil
}

func resourceAreaPermissionsCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.CSS,
		clients.SecurityClient,
		clients.IdentityClient)

	if err != nil {
		return err
	}

	aclToken, err := createAreaToken(clients.Ctx, clients.WorkItemTrackingClient, d)
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
	clients := m.(*client.AggregatedClient)

	aclToken, err := createAreaToken(clients.Ctx, clients.WorkItemTrackingClient, d)
	if err != nil {
		return err
	}

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.CSS,
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
	return resourceAreaPermissionsCreate(d, m)
}

func resourceAREAPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	aclToken, err := createAreaToken(clients.Ctx, clients.WorkItemTrackingClient, d)
	if err != nil {
		return err
	}

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.CSS,
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
