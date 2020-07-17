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

func resourceIterationPermissions() *schema.Resource {
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

func getIterationIDbyPath(clients *config.AggregatedClient, d *schema.ResourceData, path string) (*string, error) {
	var IterationID string = ""
	projectID := d.Get("project_id").(string)

	args := workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		Path:           &path,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
		Depth:          converter.Int(999),
	}

	Iteration, err := clients.WitClient.GetClassificationNode(clients.Ctx, args)
	if err != nil {
		return &IterationID, fmt.Errorf("Error getting Iteration: %+v", err)
	}

	IterationID = Iteration.Identifier.String()
	return &IterationID, nil
}

func createIterationToken(clients *config.AggregatedClient, d *schema.ResourceData) (*string, error) {
	var aclToken string
	var aclTokenPrefix string = "vstfs:///Classification/Node/"
	projectID := d.Get("project_id").(string)

	// you have to ommit the path property to get the
	// root Iteration.
	args := workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
		Depth:          converter.Int(999),
	}

	rootIteration, err := clients.WitClient.GetClassificationNode(clients.Ctx, args)
	if err != nil {
		return nil, fmt.Errorf("Error getting Iteration: %+v", err)
	}

	/*
	 * Token format
	 * Root Iteration: vstfs:///Classification/Node/<IterationIdentifier>:vstfs:///Classification/Node/f8c5b667-91dd-4fe7-bf23-3138c439d07e"
	 * 1st child: vstfs:///Classification/Node/<IterationIdentifier>:vstfs:///Classification/Node/<IterationIdentifier>
	 */
	path, ok := d.GetOk("path")

	if !ok {
		// no path was specified we use the root Iteration
		aclToken = "vstfs:///Classification/Node/" + rootIteration.Identifier.String()
	} else {
		if !*rootIteration.HasChildren {
			return &aclToken, fmt.Errorf("A path was specified but the root Iteration has no children")
		} else {
			// get the id for each Iteration in the provided path
			// we do this by querying each path element
			// 0: foo
			// 1: foo/bar
			// 3: foo/bar/baz
			aclToken = aclTokenPrefix + rootIteration.Identifier.String()
			pathElem := strings.Split(path.(string), "/")
			for i, v := range pathElem {
				var pathQuery string
				if i == 0 {
					pathQuery = v
				} else {
					pathQuery = strings.Join(commonhelper.SelectArrayRange(pathElem, 0, i), "/")
				}
				currID, _ := getIterationIDbyPath(clients, d, pathQuery)
				aclToken = aclToken + ":" + aclTokenPrefix + *currID
			}
		}
	}

	log.Printf("[DEBUG] createIterationToken(): Discovered aclToken %q", aclToken)
	return &aclToken, nil
}

func resourceIterationPermissionsCreate(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.Iteration,
		clients.Ctx,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	aclToken, err := createIterationToken(clients, d)
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
	debugWait()

	clients := m.(*config.AggregatedClient)

	aclToken, err := createIterationToken(clients, d)
	if err != nil {
		return err
	}

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.Iteration,
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

func resourceIterationPermissionsUpdate(d *schema.ResourceData, m interface{}) error {
	debugWait()

	return resourceIterationPermissionsCreate(d, m)
}

func resourceIterationPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)

	aclToken, err := createIterationToken(clients, d)
	if err != nil {
		return err
	}

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.Iteration,
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

func resourceIterationPermissionsImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	debugWait()

	// repoV2/#ProjectID#/#RepositoryID#/refs/heads/#BranchName#/#SubjectDescriptor#
	return nil, errors.New("resourceIterationPermissionsImporter: Not implemented")
}
