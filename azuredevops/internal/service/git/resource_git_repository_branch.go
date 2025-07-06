package git

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

	branchName := d.Get("name").(string)
	if strings.HasPrefix(branchName, REF_BRANCH_PREFIX) {
		return diag.Errorf("Branch name must be in short format without refs/heads/ prefix, got: %q", branchName)
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
			return diag.FromErr(fmt.Errorf("Getting refs matching %q: %w", filter, err))
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

		switch {
		case gotRef.PeeledObjectId != nil:
			newObjectId = *gotRef.PeeledObjectId
		case gotRef.ObjectId != nil:
			newObjectId = *gotRef.ObjectId
		default:
			return diag.FromErr(fmt.Errorf("GetRefs response doesn't have a valid commit id."))
		}
	}

	if err := updateRefs(clients, git.UpdateRefsArgs{
		RefUpdates: &[]git.GitRefUpdate{{
			Name:        converter.String(REF_BRANCH_PREFIX + branchName),
			NewObjectId: &newObjectId,
			OldObjectId: converter.String("0000000000000000000000000000000000000000"),
		}},
		RepositoryId: converter.String(repoId),
	}); err != nil {
		return diag.FromErr(fmt.Errorf("Creating branch %q: %+v", branchName, err))
	}

	d.SetId(fmt.Sprintf("%s:%s", repoId, branchName))
	return resourceGitRepositoryBranchRead(ctx, d, m)
}

func resourceGitRepositoryBranchRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	repoId, branchName, err := tfhelper.ParseGitRepoBranchID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if strings.HasPrefix(branchName, REF_BRANCH_PREFIX) {
		return diag.Errorf("Branch name must be in short format without refs/heads/ prefix, got: %q", branchName)
	}

	gotBranch, err := clients.GitReposClient.GetBranch(clients.Ctx, git.GetBranchArgs{
		RepositoryId: &repoId,
		Name:         &branchName,
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		var branchErr azuredevops.WrappedError
		if errors.As(err, &branchErr) {
			regx := regexp.MustCompile(fmt.Sprintf("Branch \"%[1]s\" does not exist in the %[2]s repository.", branchName, repoId))
			if branchErr.Message != nil && regx.MatchString(*branchErr.Message) {
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(fmt.Errorf("Reading branch %q: %w", branchName, err))
	}

	d.Set("name", branchName)
	d.Set("repository_id", repoId)
	d.Set("last_commit_id", *gotBranch.Commit.CommitId)

	return nil
}

func resourceGitRepositoryBranchDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	repoId, branchName, err := tfhelper.ParseGitRepoBranchID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if strings.HasPrefix(branchName, REF_BRANCH_PREFIX) {
		return diag.Errorf("Branch name must be in short format without refs/heads/ prefix, got: %q", branchName)
	}

	gotBranch, err := clients.GitReposClient.GetBranch(clients.Ctx, git.GetBranchArgs{
		RepositoryId: &repoId,
		Name:         &branchName,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("Getting latest commit of %q: %w", branchName, err))
	}

	if err := updateRefs(clients, git.UpdateRefsArgs{
		RefUpdates: &[]git.GitRefUpdate{{
			Name:        converter.String(REF_BRANCH_PREFIX + branchName),
			OldObjectId: gotBranch.Commit.CommitId,
			NewObjectId: converter.String("0000000000000000000000000000000000000000"),
		}},
		RepositoryId: converter.String(repoId),
	}); err != nil {
		return diag.FromErr(fmt.Errorf("Deleting branch %q: %w", branchName, err))
	}

	return nil
}

func updateRefs(clients *client.AggregatedClient, args git.UpdateRefsArgs) error {
	updateRefResults, err := clients.GitReposClient.UpdateRefs(clients.Ctx, args)
	if err != nil {
		return fmt.Errorf("Updating refs: %w", err)
	}

	for _, refUpdate := range *updateRefResults {
		if !*refUpdate.Success {
			return fmt.Errorf("Update refs failed. Update Status: %s", *refUpdate.UpdateStatus)
		}
	}

	return nil
}

func withPrefix(prefix, name string) string {
	if strings.HasPrefix(name, prefix) {
		return name
	}
	return prefix + name
}
