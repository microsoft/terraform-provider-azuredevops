package git

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

const (
	REF_BRANCH_PREFIX = "refs/heads/"
	REF_TAG_PREFIX    = "refs/tags/"
)

// ResourceGitRepositoryBranch schema to manage the lifecycle of a git repository branch
func ResourceGitRepositoryBranch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitRepositoryBranchCreate,
		ReadContext:   resourceGitRepositoryBranchRead,
		DeleteContext: resourceGitRepositoryBranchDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"repository_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"ref_branch": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"ref_tag", "ref_commit_id"},
			},
			"ref_tag": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"ref_branch", "ref_commit_id"},
			},
			"ref_commit_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"ref_branch", "ref_tag"},
			},
			"last_commit_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGitRepositoryBranchCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	repoId := d.Get("repository_id").(string)

	name := d.Get("name").(string)
	shortBranchName := withoutPrefix(REF_BRANCH_PREFIX, name)
	longBranchName := withPrefix(REF_BRANCH_PREFIX, name)
	if name != shortBranchName {
		return diag.Errorf("Branch name must be in short format without refs/heads/ prefix, got: %q", name)
	}

	var newObjectId string
	if v, ok := d.GetOk("ref_commit_id"); ok {
		newObjectId = v.(string)
	} else {
		var rs string
		if v, ok := d.GetOk("ref_branch"); ok {
			rs = withPrefix(REF_BRANCH_PREFIX, v.(string))
		}
		if v, ok := d.GetOk("ref_tag"); ok {
			rs = withPrefix(REF_TAG_PREFIX, v.(string))
		}

		filter := strings.TrimPrefix(rs, "refs/")
		gotRefs, err := clients.GitReposClient.GetRefs(clients.Ctx, git.GetRefsArgs{
			RepositoryId: converter.String(repoId),
			Filter:       converter.String(filter),
			Top:          converter.Int(1),
			PeelTags:     converter.Bool(true),
		})
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error getting refs matching %q: %w", filter, err))
		}

		if len(gotRefs.Value) == 0 {
			return diag.FromErr(fmt.Errorf("No refs found that match ref %q.", rs))
		}

		gotRef := gotRefs.Value[0]
		if gotRef.Name == nil {
			return diag.FromErr(fmt.Errorf("Got unexpected GetRefs response, a ref without a name was returned."))
		}

		// Check for complete match. Sometimes refs exist that match prefix with Ref, but do not match completely.
		if *gotRef.Name != rs {
			return diag.FromErr(fmt.Errorf("Ref %q not found, closest match is %q.", filter, *gotRef.Name))
		}

		if gotRef.PeeledObjectId != nil {
			newObjectId = *gotRef.PeeledObjectId
		} else if gotRef.ObjectId != nil {
			newObjectId = *gotRef.ObjectId
		} else {
			return diag.FromErr(fmt.Errorf("GetRefs response doesn't have a valid commit id."))
		}
	}

	_, err := updateRefs(clients, git.UpdateRefsArgs{
		RefUpdates: &[]git.GitRefUpdate{{
			Name:        &longBranchName,
			NewObjectId: &newObjectId,
			OldObjectId: converter.String("0000000000000000000000000000000000000000"),
		}},
		RepositoryId: converter.String(repoId),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error creating branch %q: %w", shortBranchName, err))
	}

	d.SetId(fmt.Sprintf("%s:%s", repoId, shortBranchName))

	return resourceGitRepositoryBranchRead(ctx, d, m)
}

func resourceGitRepositoryBranchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	repoId, name, err := tfhelper.ParseGitRepoBranchID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	shortBranchName := withoutPrefix(REF_BRANCH_PREFIX, name)
	if name != shortBranchName {
		return diag.Errorf("Branch name must be in short format without refs/heads/ prefix, got: %q", name)
	}

	gotBranch, err := clients.GitReposClient.GetBranch(clients.Ctx, git.GetBranchArgs{
		RepositoryId: converter.String(repoId),
		Name:         converter.String(shortBranchName),
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("Error reading branch %q: %w", name, err))
	}

	d.SetId(fmt.Sprintf("%s:%s", repoId, shortBranchName))
	d.Set("name", shortBranchName)
	d.Set("repository_id", repoId)
	d.Set("last_commit_id", *gotBranch.Commit.CommitId)

	return nil
}

func resourceGitRepositoryBranchDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	repoId, name, err := tfhelper.ParseGitRepoBranchID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	shortBranchName := withoutPrefix(REF_BRANCH_PREFIX, name)
	longBranchName := withPrefix(REF_BRANCH_PREFIX, name)
	if name != shortBranchName {
		return diag.Errorf("Branch name must be in short format without refs/heads/ prefix, got: %q", name)
	}

	gotBranch, err := clients.GitReposClient.GetBranch(clients.Ctx, git.GetBranchArgs{
		RepositoryId: converter.String(repoId),
		Name:         converter.String(shortBranchName),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error getting latest commit of %q: %w", name, err))
	}

	_, err = updateRefs(clients, git.UpdateRefsArgs{
		RefUpdates: &[]git.GitRefUpdate{{
			Name:        converter.String(longBranchName),
			OldObjectId: gotBranch.Commit.CommitId,
			NewObjectId: converter.String("0000000000000000000000000000000000000000"),
		}},
		RepositoryId: converter.String(repoId),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error deleting branch %q: %w", name, err))
	}

	return nil
}

func updateRefs(clients *client.AggregatedClient, args git.UpdateRefsArgs) (*[]git.GitRefUpdateResult, error) {
	updateRefResults, err := clients.GitReposClient.UpdateRefs(clients.Ctx, args)
	if err != nil {
		return nil, err
	}

	for _, refUpdate := range *updateRefResults {
		if !*refUpdate.Success {
			return nil, fmt.Errorf("Error got invalid GitRefUpdate.UpdateStatus: %s", *refUpdate.UpdateStatus)
		}
	}

	return updateRefResults, nil
}

func withPrefix(prefix, name string) string {
	if strings.HasPrefix(name, prefix) {
		return name
	}
	return prefix + name
}

func withoutPrefix(prefix, name string) string {
	if strings.HasPrefix(name, prefix) {
		return name[len(prefix):]
	}
	return name
}
