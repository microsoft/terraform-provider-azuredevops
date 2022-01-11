package permissions

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/build"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceBuildDefinitionPermissions schema and implementation for build permission resource
func ResourceBuildDefinitionPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceBuildDefinitionPermissionsCreateOrUpdate,
		Read:   resourceBuildDefinitionPermissionsRead,
		Update: resourceBuildDefinitionPermissionsCreateOrUpdate,
		Delete: resourceBuildDefinitionPermissionsDelete,
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
			},
			"build_definition_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		}),
	}
}

func resourceBuildDefinitionPermissionsCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Build, createBuildToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, nil, false); err != nil {
		return err
	}

	return resourceBuildDefinitionPermissionsRead(d, m)
}

func resourceBuildDefinitionPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Build, createBuildToken)
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

func resourceBuildDefinitionPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.Build, createBuildToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, &securityhelper.PermissionTypeValues.NotSet, true); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createBuildToken(d *schema.ResourceData, clients *client.AggregatedClient) (string, error) {
	projectID, ok := d.GetOk("project_id")
	if !ok {
		return "", fmt.Errorf("Failed to get 'project_id' from schema")
	}

	buildDefinitionID, err := getBuildDefinitionID(d)
	if err != nil {
		return "", err
	}

	definition, err := clients.BuildClient.GetDefinition(clients.Ctx, build.GetDefinitionArgs{
		Project:      converter.String(projectID.(string)),
		DefinitionId: converter.Int(buildDefinitionID),
	})

	if err != nil {
		return "", err
	}

	var aclToken string

	// The token format is Project_ID/Build_Definition_ID
	// or Project_ID/Path/Build_Definition_ID

	if *definition.Path != "\\" {
		transformedPath := transformPath(*definition.Path)

		aclToken = fmt.Sprintf("%s/%s/%d", projectID.(string), transformedPath, buildDefinitionID)
	} else {
		aclToken = fmt.Sprintf("%s/%d", projectID.(string), buildDefinitionID)
	}

	return aclToken, nil
}

// transformPath must return a path with forward slashes
func transformPath(path string) string {
	paths := strings.Split(path, "\\")
	transformedPath := strings.Join(paths, "/")

	// remove slash at front of string
	transformedPath = strings.TrimPrefix(transformedPath, "/")

	// remove slash at end of string
	transformedPath = strings.TrimPrefix(transformedPath, "")

	return transformedPath
}

func getBuildDefinitionID(d *schema.ResourceData) (int, error) {
	buildID, ok := d.GetOk("build_definition_id")
	if !ok {
		return -1, fmt.Errorf("Failed to get 'build_definition_id' from schema")
	}

	id, err := strconv.Atoi(buildID.(string))
	if err != nil {
		return -1, err
	}

	return id, nil
}
