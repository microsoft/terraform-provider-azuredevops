package permissions

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceWorkItemQueryPermissions schema and implementation for project permission resource
func ResourceWorkItemQueryPermissions() *schema.Resource {
	return &schema.Resource{
		Create: ResourceWorkItemQueryPermissionsCreateOrUpdate,
		Read:   ResourceWorkItemQueryPermissionsRead,
		Update: ResourceWorkItemQueryPermissionsCreateOrUpdate,
		Delete: ResourceWorkItemQueryPermissionsDelete,
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
				Optional:     true,
				ForceNew:     true,
			},
		}),
	}
}

// ResourceWorkItemQueryPermissionsCreateOrUpdate create or update project permissions
func ResourceWorkItemQueryPermissionsCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.WorkItemQueryFolders, createWorkItemQueryToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, nil, false); err != nil {
		return err
	}

	return ResourceWorkItemQueryPermissionsRead(d, m)
}

// ResourceWorkItemQueryPermissionsRead read project permissions
func ResourceWorkItemQueryPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.WorkItemQueryFolders, createWorkItemQueryToken)
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

// ResourceWorkItemQueryPermissionsDelete remove project permissions
func ResourceWorkItemQueryPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.WorkItemQueryFolders, createWorkItemQueryToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, &securityhelper.PermissionTypeValues.NotSet, true); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createWorkItemQueryToken(d *schema.ResourceData, clients *client.AggregatedClient) (string, error) {
	projectID, ok := d.GetOk("project_id")
	if !ok {
		return "", fmt.Errorf("Failed to get 'project_id' from schema")
	}
	aclToken := fmt.Sprintf("$/%s", projectID.(string))
	path, ok := d.GetOk("path")
	if ok {
		idList, err := getQueryIDsFromPath(clients.Ctx, clients.WorkItemTrackingClient, projectID.(string), path.(string))
		if err != nil {
			return "", err
		}

		aclToken = fmt.Sprintf("%s/%s", aclToken, strings.Join(*idList, "/"))
	}
	return aclToken, nil
}

func getQueryIDsFromPath(context context.Context, wiqClient workitemtracking.Client, projectID string, path string) (*[]string, error) {
	var pathItems []string
	var err error
	var qry *workitemtracking.QueryHierarchyItem
	ret := []string{}

	path = strings.TrimSpace(path)
	linq.From(strings.Split(path, "/")).
		Where(func(elem interface{}) bool {
			return len(elem.(string)) > 0
		}).
		ToSlice(&pathItems)

	qry, err = wiqClient.GetQuery(context, workitemtracking.GetQueryArgs{
		Project: &projectID,
		Query:   converter.String("Shared Queries"),
		Depth:   converter.Int(1),
	})
	if err != nil {
		return nil, err
	}
	ret = append(ret, qry.Id.String())
	if len(pathItems) > 0 {
		for _, v := range pathItems {
			if qry.Children == nil || len(*qry.Children) <= 0 {
				return nil, fmt.Errorf("Unable to find query [%s] in folder [%s] because it has no children", v, converter.ToString(qry.Name, qry.Id.String()))
			}

			uuid, ok := uuid.Parse(v)
			chldIdx := -1
			for idx, chldItem := range *qry.Children {
				if ok == nil && strings.EqualFold(uuid.String(), chldItem.Id.String()) {
					chldIdx = idx
				} else if chldItem.Name != nil && strings.EqualFold(*chldItem.Name, v) {
					chldIdx = idx
				}
			}

			if chldIdx < 0 {
				return nil, fmt.Errorf("Unable to find query [%s] in folder [%s]", v, converter.ToString(qry.Name, qry.Id.String()))
			}

			qry, err = wiqClient.GetQuery(context, workitemtracking.GetQueryArgs{
				Project: &projectID,
				Query:   converter.String((*qry.Children)[chldIdx].Id.String()),
				Depth:   converter.Int(1),
			})
			if err != nil {
				return nil, err
			}
			ret = append(ret, qry.Id.String())
		}
	}
	return &ret, nil
}
