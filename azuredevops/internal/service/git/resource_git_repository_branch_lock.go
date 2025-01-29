package git

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceGitRepositoryBranchLock() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitRepositoryBranchLockCreate,
		ReadContext:   resourceGitRepositoryBranchLockRead,
		DeleteContext: resourceGitRepositoryBranchLockDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"branch": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"is_locked": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
		},
	}
}

func resourceGitRepositoryBranchLockCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	repoID := d.Get("repository_id").(string)
	branch := d.Get("branch").(string)
	isLocked := d.Get("is_locked").(bool)

	if strings.HasPrefix(branch, REF_BRANCH_PREFIX) {
		return diag.Errorf("Branch name must be in short format without refs/heads/ prefix, got: %q", branch)
	}

	branchRef := REF_BRANCH_PREFIX + branch

	// Get current ref to get the objectId
	refs, err := clients.GitReposClient.GetRefs(ctx, git.GetRefsArgs{
		RepositoryId: &repoID,
		Filter:       converter.String(strings.TrimPrefix(branchRef, "refs/")),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error getting branch ref: %v", err))
	}
	if len(refs.Value) == 0 {
		return diag.Errorf("Branch %s not found", branch)
	}

	// Update the ref with lock status
	refUpdate := git.GitRefUpdate{
		IsLocked: converter.Bool(isLocked),
	}

	_, err = clients.GitReposClient.UpdateRef(ctx, git.UpdateRefArgs{
		NewRefInfo:   &refUpdate,
		RepositoryId: &repoID,
		Filter:       converter.String(strings.TrimPrefix(branchRef, "refs/")),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error updating branch lock: %v", err))
	}
	d.SetId(fmt.Sprintf("%s:%s", repoID, branch))
	return resourceGitRepositoryBranchLockRead(ctx, d, m)
}

func resourceGitRepositoryBranchLockRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	repoID, branch, err := parseBranchLockID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	branchRef := REF_BRANCH_PREFIX + branch
	refs, err := clients.GitReposClient.GetRefs(ctx, git.GetRefsArgs{
		RepositoryId: &repoID,
		Filter:       converter.String(strings.TrimPrefix(branchRef, "refs/")),
	})

	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("Error reading branch ref: %v", err))
	}

	if len(refs.Value) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("is_locked", refs.Value[0].IsLocked)
	return nil
}

func resourceGitRepositoryBranchLockDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// For delete, we'll just set isLocked to false
	clients := m.(*client.AggregatedClient)
	repoID, branch, err := parseBranchLockID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	branchRef := REF_BRANCH_PREFIX + branch
	refs, err := clients.GitReposClient.GetRefs(ctx, git.GetRefsArgs{
		RepositoryId: &repoID,
		Filter:       converter.String(strings.TrimPrefix(branchRef, "refs/")),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error getting branch ref: %v", err))
	}
	if len(refs.Value) == 0 {
		return nil
	}

	// Update the ref with lock status of False
	refUpdate := git.GitRefUpdate{
		IsLocked: converter.Bool(false),
	}

	_, err = clients.GitReposClient.UpdateRef(ctx, git.UpdateRefArgs{
		NewRefInfo:   &refUpdate,
		RepositoryId: &repoID,
		Filter:       converter.String(strings.TrimPrefix(branchRef, "refs/")),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error unlocking branch: %v", err))
	}

	return nil
}

func parseBranchLockID(id string) (repoID string, branch string, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = fmt.Errorf("Invalid branch lock ID: %s", id)
		return
	}
	repoID = parts[0]
	branch = parts[1]
	return
}
