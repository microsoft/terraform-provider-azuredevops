package azuredevops

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/securityhelper"
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
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
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
			},
			"branch_name": {
				Type:         schema.TypeString,
				ValidateFunc: validate.NoEmptyStrings,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"repository_id"},
			},
		}),
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
	repositoryID, repoOk := d.GetOkExists("repository_id")
	if repoOk {
		aclToken += "/" + repositoryID.(string)
	}
	branchName, branchOk := d.GetOkExists("branch_name")
	if branchOk {
		if !repoOk {
			return nil, fmt.Errorf("Unable to create ACL token for branch %s, because no respository is specified", branchName)
		}
		branch, err := getBranchByName(clients,
			converter.StringFromInterface(repositoryID),
			converter.StringFromInterface(branchName))
		if err != nil {
			return nil, err
		}
		branchPath := strings.Split(*branch.Name, "/")
		branchName, err = converter.EncodeUtf16HexString(branchPath[len(branchPath)-1])
		if err != nil {
			return nil, err
		}
		aclToken += "/refs/heads/" + branchName.(string)
	}
	return &aclToken, nil
}

func getBranchByName(clients *config.AggregatedClient, repositoryID *string, branchName *string) (*git.GitRef, error) {
	filter := "heads/" + *branchName
	res, err := clients.GitReposClient.GetRefs(clients.Ctx, git.GetRefsArgs{
		RepositoryId: repositoryID,
		Filter:       &filter,
	})
	if err != nil {
		return nil, err
	}
	item := linq.From(res.Value).FirstWith(func(elem interface{}) bool {
		return strings.HasSuffix(*(elem.(git.GitRef).Name), *branchName)
	})
	if item == nil {
		return nil, fmt.Errorf("No branch found with name [%s] in repository with id [%s]", *branchName, *repositoryID)
	}
	gitRef := item.(git.GitRef)
	return &gitRef, nil
}

func resourceGitPermissionsCreate(d *schema.ResourceData, m interface{}) error {
	debugWait()

	clients := m.(*config.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(securityhelper.SecurityNamespaceIDValues.GitRepositories,
		clients.Ctx,
		clients.SecurityClient,
		clients.IdentityClient)
	if err != nil {
		return err
	}

	aclToken, err := createGitToken(clients, d)
	if err != nil {
		return err
	}

	err = securityhelper.SetPrincipalPermissions(d, sn, aclToken, nil, false)
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

	principalPermissions, err := securityhelper.GetPrincipalPermissions(d, sn, aclToken)
	if err != nil {
		return err
	}

	d.Set("permissions", principalPermissions.Permissions)
	return nil
}

func resourceGitPermissionsUpdate(d *schema.ResourceData, m interface{}) error {
	debugWait()

	return resourceGitPermissionsCreate(d, m)
}

func resourceGitPermissionsDelete(d *schema.ResourceData, m interface{}) error {
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

	err = securityhelper.SetPrincipalPermissions(d, sn, aclToken, &securityhelper.PermissionTypeValues.NotSet, true)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceGitPermissionsImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	debugWait()

	// repoV2/#ProjectID#/#RepositoryID#/refs/heads/#BranchName#/#SubjectDescriptor#
	return nil, errors.New("resourceGitPermissionsImporter: Not implemented")
}
