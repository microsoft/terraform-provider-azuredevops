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
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// ResourceGitRef schema to manage the lifecycle of a git ref
func ResourceGitRef() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitRefCreate,
		ReadContext:   resourceGitRefRead,
		DeleteContext: resourceGitRefDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"object_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGitRefCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)
	repoId := d.Get("repository_id").(string)
	refName := d.Get("name").(string)

	var newObjectId string
	if v, ok := d.GetOk("ref_commit_id"); ok {
		newObjectId = v.(string)
	} else {
		var rs string
		if v, ok := d.GetOk("ref_branch"); ok {
			rs = withPrefix(REF_BRANCH_PREFIX, v.(string))
		} else if v, ok := d.GetOk("ref_tag"); ok {
			rs = withPrefix(REF_TAG_PREFIX, v.(string))
		} else {
			return diag.Errorf("One of ref_branch, ref_tag or ref_commit_id must be specified")
		}

		filter := strings.TrimPrefix(rs, "refs/")
		args := git.GetRefsArgs{
			RepositoryId: &repoId,
			Filter:       &filter,
			Top:          converter.Int(1),
			PeelTags:     converter.Bool(true),
		}
		if v, ok := d.GetOk("project_id"); ok {
			projectID := v.(string)
			args.Project = &projectID
		}
		gotRefs, err := clients.GitReposClient.GetRefs(ctx, args)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Getting refs matching %q: %w", filter, err))
		}

		if gotRefs == nil || len(gotRefs.Value) == 0 {
			return diag.FromErr(fmt.Errorf("No refs found that match ref %q.", rs))
		}

		gotRef := gotRefs.Value[0]
		if gotRef.Name == nil {
			return diag.FromErr(fmt.Errorf("Got unexpected GetRefs response, a ref without a name was returned."))
		}

		if *gotRef.Name != rs {
			return diag.FromErr(fmt.Errorf("Ref %q not found, closest match is %q.", rs, *gotRef.Name))
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

	args := git.UpdateRefsArgs{
		RefUpdates: &[]git.GitRefUpdate{{
			Name:        &refName,
			NewObjectId: &newObjectId,
			OldObjectId: converter.String("0000000000000000000000000000000000000000"),
		}},
		RepositoryId: &repoId,
	}
	if v, ok := d.GetOk("project_id"); ok {
		projectID := v.(string)
		args.Project = &projectID
	}

	if err := updateRefs(clients, args); err != nil {
		return diag.FromErr(fmt.Errorf("Creating ref %q: %+v", refName, err))
	}

	d.SetId(fmt.Sprintf("%s:%s", repoId, refName))
	return resourceGitRefRead(ctx, d, m)
}

func resourceGitRefRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return diag.Errorf("Invalid resource ID: %s", d.Id())
	}
	repoId := parts[0]
	refName := parts[1]

	args := git.GetRefsArgs{
		RepositoryId: &repoId,
		Filter:       &refName,
	}
	if v, ok := d.GetOk("project_id"); ok {
		projectID := v.(string)
		args.Project = &projectID
	}

	gotRefs, err := clients.GitReposClient.GetRefs(ctx, args)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Reading ref %q: %w", refName, err))
	}

	var match *git.GitRef
	if gotRefs != nil {
		for _, ref := range gotRefs.Value {
			if ref.Name != nil && *ref.Name == refName {
				match = &ref
				break
			}
		}
	}

	if match == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", refName)
	d.Set("repository_id", repoId)
	if match.ObjectId != nil {
		d.Set("object_id", *match.ObjectId)
	}

	return nil
}

func resourceGitRefDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clients := m.(*client.AggregatedClient)

	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return diag.Errorf("Invalid resource ID: %s", d.Id())
	}
	repoId := parts[0]
	refName := parts[1]
	objectId := d.Get("object_id").(string)

	args := git.UpdateRefsArgs{
		RefUpdates: &[]git.GitRefUpdate{{
			Name:        &refName,
			OldObjectId: &objectId,
			NewObjectId: converter.String("0000000000000000000000000000000000000000"),
		}},
		RepositoryId: &repoId,
	}
	if v, ok := d.GetOk("project_id"); ok {
		projectID := v.(string)
		args.Project = &projectID
	}

	if err := updateRefs(clients, args); err != nil {
		return diag.FromErr(fmt.Errorf("Deleting ref %q: %w", refName, err))
	}

	return nil
}
