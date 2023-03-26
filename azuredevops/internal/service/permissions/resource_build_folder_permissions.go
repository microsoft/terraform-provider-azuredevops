package permissions

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/build"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceBuildFolderPermissions schema and implementation for build permission resource
func ResourceBuildFolderPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceBuildFolderPermissionsCreateOrUpdate,
		Read:   resourceBuildFolderPermissionsRead,
		Update: resourceBuildFolderPermissionsCreateOrUpdate,
		Delete: resourceBuildFolderPermissionsDelete,
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		}),
	}
}

func resourceBuildFolderPermissionsCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Build, createBuildFolderToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, nil, false); err != nil {
		return err
	}

	return resourceBuildFolderPermissionsRead(d, m)
}

func resourceBuildFolderPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Build, createBuildFolderToken)
	if err != nil {
		return err
	}

	principalPermissions, err := securityhelper.GetPrincipalPermissions(d, sn)
	if err != nil {
		return err
	}
	if principalPermissions == nil {
		d.SetId("")
		log.Printf("[INFO] Permissions for ACL token %q not found. Removing from state", sn.GetToken())
		return nil
	}

	d.Set("permissions", principalPermissions.Permissions)
	return nil
}

func resourceBuildFolderPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Build, createBuildFolderToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, &securityhelper.PermissionTypeValues.NotSet, true); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createBuildFolderToken(d *schema.ResourceData, clients *client.AggregatedClient) (string, error) {
	projectID, ok := d.GetOk("project_id")
	if !ok {
		return "", fmt.Errorf("Failed to get 'project_id' from schema")
	}

	buildFolderPath, ok := d.GetOk("path")
	if !ok {
		return "", fmt.Errorf("Failed to get 'path' from schema")
	}

	buildFolders, err := clients.BuildClient.GetFolders(clients.Ctx, build.GetFoldersArgs{
		Project: converter.String(projectID.(string)),
		Path:    converter.String(buildFolderPath.(string)),
	})

	if err != nil {
		return "", fmt.Errorf(" failed to get the folder. Project ID: %s, Path: %s. %+v", projectID, buildFolderPath, err)
	}

	if buildFolders == nil || len(*buildFolders) == 0 {
		return "", fmt.Errorf(" folder not found. Project ID: %s, Path: %s.", projectID, buildFolderPath)
	}

	Folder := (*buildFolders)[0]

	var aclToken string

	// The token format is Project_ID/Path
	if *Folder.Path != "\\" {
		transformedPath := transformPath(*Folder.Path)

		aclToken = fmt.Sprintf("%s/%s", projectID.(string), transformedPath)
	} else {
		aclToken = fmt.Sprintf("%s/%s", projectID.(string), *Folder.Path)
	}

	return aclToken, nil
}
