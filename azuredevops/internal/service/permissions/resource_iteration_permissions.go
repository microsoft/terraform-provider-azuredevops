package permissions

import (
	"context"
	"errors"
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

func ResourceIterationPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceIterationPermissionsCreate,
		Read:   resourceIterationPermissionsRead,
		Update: resourceIterationPermissionsUpdate,
		Delete: resourceIterationPermissionsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceIterationPermissionsImporter,
		},
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

func getIterationIDbyPath(context context.Context, workitemtrackingClient workitemtracking.Client, projectID string, path string) (*string, error) {
	var IterationID string = ""

	Iteration, err := workitemtrackingClient.GetClassificationNode(context, workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		Path:           &path,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
		Depth:          converter.Int(1),
	})
	if err != nil {
		return &IterationID, fmt.Errorf("Error getting Iteration: %w", err)
	}

	IterationID = Iteration.Identifier.String()
	return &IterationID, nil
}

func createIterationToken(context context.Context, workitemtrackingClient workitemtracking.Client, d *schema.ResourceData) (*string, error) {
	const aclTokenPrefix = "vstfs:///Classification/Node/"
	var aclToken string
	projectID := d.Get("project_id").(string)

	// you have to ommit the path property to get the
	// root Iteration.
	rootIteration, err := workitemtrackingClient.GetClassificationNode(context, workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
		Depth:          converter.Int(1),
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting Iteration: %w", err)
	}

	/*
	 * Token format
	 * Root Iteration: vstfs:///Classification/Node/<IterationIdentifier>:vstfs:///Classification/Node/f8c5b667-91dd-4fe7-bf23-3138c439d07e"
	 * 1st child: vstfs:///Classification/Node/<IterationIdentifier>:vstfs:///Classification/Node/<IterationIdentifier>
	 */
	aclToken = aclTokenPrefix + rootIteration.Identifier.String()
	path, ok := d.GetOk("path")

	if ok {
		if !*rootIteration.HasChildren {
			return nil, fmt.Errorf("A path was specified but the root Iteration has no children")
		} else {
			// get the id for each Iteration in the provided path
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
				currID, err := getIterationIDbyPath(context, workitemtrackingClient, projectID, pathItem)
				if err != nil {
					return nil, fmt.Errorf("Failed to get ID for iteration %s, %w", pathItem, err)
				}
				aclToken = aclToken + ":" + aclTokenPrefix + *currID
			}
		}
	}

	log.Printf("[DEBUG] createIterationToken(): Discovered aclToken %q", aclToken)
	return &aclToken, nil
}

func resourceIterationPermissionsCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.Iteration,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	aclToken, err := createIterationToken(clients.Ctx, clients.WorkItemTrackingClient, d)
	if err != nil {
		return err
	}

	err = securityhelper.SetPrincipalPermissions(d, sn, aclToken, nil, false)
	if err != nil {
		return err
	}

	return resourceIterationPermissionsRead(d, m)
}

func resourceIterationPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.Iteration,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	aclToken, err := createIterationToken(clients.Ctx, clients.WorkItemTrackingClient, d)
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

func resourceIterationPermissionsUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceIterationPermissionsCreate(d, m)
}

func resourceIterationPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.Iteration,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	aclToken, err := createIterationToken(clients.Ctx, clients.WorkItemTrackingClient, d)
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

func resourceIterationPermissionsImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// repoV2/#ProjectID#/#RepositoryID#/refs/heads/#BranchName#/#SubjectDescriptor#
	return nil, errors.New("resourceIterationPermissionsImporter: Not implemented")
}
