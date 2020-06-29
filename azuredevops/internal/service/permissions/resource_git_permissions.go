package permissions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/debug"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/validate"
)

// ResourceGitPermissions schema and implementation for Git repository permission resource
func ResourceGitPermissions() *schema.Resource {
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

func resourceGitPermissionsCreate(d *schema.ResourceData, m interface{}) error {
	debug.Wait()

	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.GitRepositories,
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
	debug.Wait()

	clients := m.(*client.AggregatedClient)

	aclToken, err := createGitToken(clients, d)
	if err != nil {
		return err
	}

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.GitRepositories,
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
	debug.Wait()

	return resourceGitPermissionsCreate(d, m)
}

func resourceGitPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	debug.Wait()

	clients := m.(*client.AggregatedClient)

	aclToken, err := createGitToken(clients, d)
	if err != nil {
		return err
	}

	sn, err := securityhelper.NewSecurityNamespace(clients.Ctx,
		securityhelper.SecurityNamespaceIDValues.GitRepositories,
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
	debug.Wait()

	// repoV2/#ProjectID#/#RepositoryID#/refs/heads/#BranchName#/#SubjectDescriptor#
	return nil, errors.New("resourceGitPermissionsImporter: Not implemented")
}

func createGitToken(clients *client.AggregatedClient, d *schema.ResourceData) (*string, error) {
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
	repositoryID, repoOk := d.GetOk("repository_id")
	if repoOk {
		aclToken += "/" + repositoryID.(string)
	}
	branchName, branchOk := d.GetOk("branch_name")
	if branchOk {
		if !repoOk {
			return nil, fmt.Errorf("Unable to create ACL token for branch %s, because no repository is specified", branchName)
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

func getBranchByName(clients *client.AggregatedClient, repositoryID *string, branchName *string) (*git.GitRef, error) {
	filter := "heads/" + *branchName
	currentToken := ""
	args := git.GetRefsArgs{
		RepositoryId: repositoryID,
		Filter:       &filter,
	}
	for hasMore := true; hasMore; {
		if currentToken != "" {
			args.ContinuationToken = &currentToken
		}
		res, err := clients.GitReposClient.GetRefs(clients.Ctx, args)
		if err != nil {
			return nil, err
		}
		currentToken = res.ContinuationToken
		hasMore = currentToken != ""
		item := linq.From(res.Value).FirstWith(func(elem interface{}) bool {
			return strings.HasSuffix(*(elem.(git.GitRef).Name), *branchName)
		})
		if item != nil {
			gitRef := item.(git.GitRef)
			return &gitRef, nil
		}
	}
	return nil, fmt.Errorf("No branch found with name [%s] in repository with id [%s]", *branchName, *repositoryID)
}
