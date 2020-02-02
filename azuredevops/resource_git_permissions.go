package azuredevops

import (
	"errors"
	"fmt"

	"github.com/ahmetb/go-linq"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/securityhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
)

func resourceGitPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitPermissionsCreate,
		Read:   resourceGitPermissionsRead,
		Update: resourceGitPermissionsUpdate,
		Delete: resourceGitPermissionsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGitPermissionsImporter,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validate.UUID,
				Required:     true,
				ForceNew:     true,
			},
			"repository_id": {
				Type:         schema.TypeString,
				ValidateFunc: validate.UUID,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"project_id"},
			},
			"branch_name": {
				Type:         schema.TypeString,
				ValidateFunc: validate.NoEmptyStrings,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"project_id", "repository_id"},
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

func createGitToken(clients *config.AggregatedClient, d *schema.ResourceData) (*string, error) {
	projectID, ok := d.GetOk("project_id")
	if !ok {
		return nil, fmt.Errorf("Failed to get 'project_id' from schema")
	}

	/*
	 * Token format
	 * ACL for ALL Git repositories in a project:                 repoV2/#ProjectID#
	 * ACL for a Git repository in a project:                     repoV2/#ProjectID#/#RepositoryID#
	 * ACL for all branches inside a Git repository in a project: repoV2/#ProjectID#/#RepositoryID#/refs/heads
	 * ACL for a branch inside a Git repository in a project:     repoV2/#ProjectID#/#RepositoryID#/refs/heads/#BranchID#
	 */
	aclToken := "repoV2/" + projectID.(string)
	repositoryID, ok := d.GetOkExists("repository_id")
	if ok {
		aclToken += "/" + repositoryID.(string)
	}
	branchName, ok := d.GetOkExists("branch_name")
	if ok {
		branch, err := getBranchByName(clients,
			converter.StringFromInterface(repositoryID),
			converter.StringFromInterface(branchName))
		if err != nil {
			return nil, err
		}
		aclToken += "/" + *branch.ObjectId
	}
	return &aclToken, nil
}

func getBranchByName(clients *config.AggregatedClient, repositoryID *string, branchName *string) (*git.GitRef, error) {
	filter := "heads/" + *branchName
	res, err := clients.GitClient.GetRefs(clients.Ctx, git.GetRefsArgs{
		RepositoryId: repositoryID,
		Filter:       &filter,
	})
	if err != nil {
		return nil, err
	}
	item := linq.From(res.Value).FirstWith(func(elem interface{}) bool {
		return *(elem.(git.GitRef).Name) == *branchName
	})
	if item == nil {
		return nil, fmt.Errorf("No branch found with name [%s] in repository with id [%s]", *repositoryID, *branchName)
	}
	gitRef := item.(git.GitRef)
	return &gitRef, nil
}

func resourceGitPermissionsCreate(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)

	aclToken, err := createGitToken(clients, d)
	if err != nil {
		return err
	}

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.GitRepositories,
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

	err = sn.SetPrincipalPermissions(&setPermissions, aclToken)
	if err != nil {
		return err
	}

	return resourceGitPermissionsRead(d, m)
}

func resourceGitPermissionsRead(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)

	aclToken, err := createGitToken(clients, d)
	if err != nil {
		return err
	}

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.GitRepositories,
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

	principalList := []string{*converter.StringFromInterface(principal)}
	principalPermissions, err := sn.GetPrincipalPermissions(aclToken, &principalList)
	if err != nil {
		return err
	}
	if principalPermissions == nil || len(*principalPermissions) != 1 {
		return fmt.Errorf("Failed to retrive current permissions for principal [%s]", principalList[0])
	}
	d.SetId(fmt.Sprintf("%s/%s", *aclToken, principal.(string)))
	for key := range ((*principalPermissions)[0]).Permissions {
		if _, ok := permissions.(map[string]interface{})[string(key)]; !ok {
			delete(((*principalPermissions)[0]).Permissions, key)
		}
	}
	d.Set("permissions", ((*principalPermissions)[0]).Permissions)
	return nil
}

func resourceGitPermissionsUpdate(d *schema.ResourceData, m interface{}) error {
	debugWait()

	return resourceGitPermissionsCreate(d, m)
}

func resourceGitPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	debugWait()

	// force all specified permissions to 'NotSet'
	return errors.New("resourceGitPermissionsDelete: Not implemented")
}

func resourceGitPermissionsImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	debugWait()

	// repoV2/#ProjectID#/#RepositoryID#/refs/heads/#BranchID#/#SubjectDescriptor#
	return nil, errors.New("resourceGitPermissionsImporter: Not implemented")
}
