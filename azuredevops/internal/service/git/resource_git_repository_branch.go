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

// ResourceGitRepositoryBranch schema to manage the lifecycle of a git repository branch
func ResourceGitRepositoryBranch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitRepositoryBranchCreate,
		ReadContext:   resourceGitRepositoryBranchRead,
		DeleteContext: resourceGitRepositoryBranchDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitRepositoryBranchImport,
		},
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
				Computed:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"ref_branch", "ref_tag"},
			},
			"branch_reference": {
				Type:     schema.TypeString,
				Computed: true,
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
	branchName := d.Get("name").(string)
	branchRefHeadName := withPrefix("refs/heads/", branchName)
	ref, hasRef := d.GetOk("ref_branch")
	tag, hasTag := d.GetOk("ref_tag")
	_, hasCommitId := d.GetOk("ref_commit_id")

	if !hasRef && !hasTag && !hasCommitId {
		return diag.Errorf("One of 'ref' or 'tag' or 'commit_id' must be provided.")
	}

	// Get a commitId from a head or tag if it is not provided in the resource
	if !hasCommitId {
		var rs string
		if hasRef {
			rs = withPrefix("refs/heads/", ref.(string))
		}
		if hasTag {
			rs = withPrefix("refs/tags/", tag.(string))
		}

		// Azuredevops GetRefs api returns refs whose "prefix" matches Filter sorted from shortest to longest
		// Top1 should return best match
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

		// Check if ref was a tag and we need to use PeeledObjectId to get the commit id of the tag
		var refObjectIdSha *string
		if gotRef.PeeledObjectId != nil {
			refObjectIdSha = gotRef.PeeledObjectId
		} else if gotRef.ObjectId != nil {
			refObjectIdSha = gotRef.ObjectId
		} else {
			return diag.FromErr(fmt.Errorf("GetRefs response doesn't have a valid commit id."))
		}
		d.Set("ref_commit_id", *refObjectIdSha)
	}
	newObjectId := d.Get("ref_commit_id").(string)

	_, err := updateRefs(clients, git.UpdateRefsArgs{
		RefUpdates: &[]git.GitRefUpdate{{
			Name:        &branchRefHeadName,
			NewObjectId: &newObjectId,
			OldObjectId: converter.String("0000000000000000000000000000000000000000"),
		}},
		RepositoryId: converter.String(repoId),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error creating branch %q: %w", branchName, err))
	}

	d.SetId(fmt.Sprintf("%s:%s", repoId, branchName))

	return resourceGitRepositoryBranchRead(ctx, d, m)
}

func resourceGitRepositoryBranchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	repoId, name, err := tfhelper.ParseGitRepoBranchID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	gotBranch, err := clients.GitReposClient.GetBranch(clients.Ctx, git.GetBranchArgs{
		RepositoryId: converter.String(repoId),
		Name:         converter.String(name),
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("Error reading branch %q: %w", name, err))
	}

	d.SetId(fmt.Sprintf("%s:%s", repoId, name))
	d.Set("name", name)
	d.Set("repository_id", repoId)
	d.Set("branch_reference", converter.String(withPrefix("refs/heads/", name)))
	d.Set("last_commit_id", *gotBranch.Commit.CommitId)

	return nil
}

func resourceGitRepositoryBranchDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	repoId, name, err := tfhelper.ParseGitRepoBranchID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	gotBranch, err := clients.GitReposClient.GetBranch(clients.Ctx, git.GetBranchArgs{
		RepositoryId: converter.String(repoId),
		Name:         converter.String(name),
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error getting latest commit of %q: %w", name, err))
	}

	_, err = updateRefs(clients, git.UpdateRefsArgs{
		RefUpdates: &[]git.GitRefUpdate{{
			Name:        converter.String(withPrefix("refs/heads/", name)),
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

func resourceGitRepositoryBranchImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	_, _, err := tfhelper.ParseGitRepoBranchID(d.Id())
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
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
