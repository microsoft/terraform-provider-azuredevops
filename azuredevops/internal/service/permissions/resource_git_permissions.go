package permissions

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceGitPermissions schema and implementation for Git repository permission resource
func ResourceGitPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitPermissionsCreateOrUpdate,
		Read:   resourceGitPermissionsRead,
		Update: resourceGitPermissionsCreateOrUpdate,
		Delete: resourceGitPermissionsDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: securityhelper.CreatePermissionResourceSchema(map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Required:     true,
				ForceNew:     true,
			},
			"repository_id": {
				Type:         schema.TypeString,
				ValidateFunc: validation.IsUUID,
				Optional:     true,
				ForceNew:     true,
			},
			"branch_name": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"repository_id"},
			},
		}),
	}
}

func resourceGitPermissionsCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.GitRepositories, createGitToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, nil, false); err != nil {
		return err
	}

	return resourceGitPermissionsRead(d, m)
}

func resourceGitPermissionsRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.GitRepositories, createGitToken)
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

func resourceGitPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	sn, err := securityhelper.NewSecurityNamespace(d, clients, securityhelper.SecurityNamespaceIDValues.GitRepositories, createGitToken)
	if err != nil {
		return err
	}

	if err := securityhelper.SetPrincipalPermissions(d, sn, &securityhelper.PermissionTypeValues.NotSet, true); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func createGitToken(d *schema.ResourceData, clients *client.AggregatedClient) (string, error) {
	projectID, ok := d.GetOk("project_id")
	if !ok {
		return "", fmt.Errorf("Failed to get 'project_id' from schema")
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
			return "", fmt.Errorf("Unable to create ACL token for branch %s, because no repository is specified", branchName)
		}

		re := regexp.MustCompile(`(/?refs/heads/)?(.*)+`)
		branchPath := re.FindStringSubmatch(branchName.(string))

		paths := strings.Split(branchPath[len(branchPath)-1], "/")
		encodedPaths := make([]string, len(paths))
		for i, subBranchPath := range paths {
			encoded, err := converter.EncodeUtf16HexString(subBranchPath)
			if err != nil {
				return "", err
			}
			encodedPaths[i] = encoded
		}

		aclToken += "/refs/heads/" + strings.Join(encodedPaths, "/")
	}
	return aclToken, nil
}
