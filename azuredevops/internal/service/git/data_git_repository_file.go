package git

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func DataGitRepositoryFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitRepositoryFileRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The repository ID",
				ValidateFunc: validation.IsUUID,
			},
			"file": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file path",
			},
			"branch": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The branch name, no default",
				ConflictsWith: []string{"tag"},
			},
			"tag": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Optional tag name, no default",
				ConflictsWith: []string{"branch"},
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The file's content",
			},
			"commit_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The commit message",
			},
		},
	}
}

func dataSourceGitRepositoryFileRead(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	clients := m.(*client.AggregatedClient)

	repoId := d.Get("repository_id").(string)
	file := d.Get("file").(string)
	branch := d.Get("branch").(string)
	tag := d.Get("tag").(string)

	var vdVersionType *git.GitVersionType
	var vdVersion *string

	if len(branch) > 0 {
		vdVersionType = &git.GitVersionTypeValues.Branch
		vdVersion = converter.String(shortBranchName(branch))
	} else if len(tag) > 0 {
		vdVersionType = &git.GitVersionTypeValues.Tag
		vdVersion = converter.String(shortTagName(tag))
	} else {
		return fmt.Errorf("One of 'branch' or 'tag' must be specified")
	}

	repoItem, err := clients.GitReposClient.GetItem(ctx, git.GetItemArgs{
		RepositoryId:   &repoId,
		Path:           &file,
		IncludeContent: converter.Bool(true),
		VersionDescriptor: &git.GitVersionDescriptor{
			Version:     vdVersion,
			VersionType: vdVersionType,
		},
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return fmt.Errorf("Item not found, repositoryID: %s, %s: %s, file: %s. Error: %+v", repoId, string(*vdVersionType), *vdVersion, file, err)
		}
		return fmt.Errorf("Get item failed, repositoryID: %s, %s: %s, file: %s. Error: %+v", repoId, string(*vdVersionType), *vdVersion, file, err)
	}
	err = d.Set("content", repoItem.Content)
	if err != nil {
		return err
	}
	commit, err := clients.GitReposClient.GetCommit(ctx, git.GetCommitArgs{
		RepositoryId: &repoId,
		CommitId:     repoItem.CommitId,
	})

	if err != nil {
		return fmt.Errorf("Get commit failed, repositoryID: %s, branch: %s, file: %s. Error:  %+v", repoId, branch, file, err)
	}

	err = d.Set("commit_message", commit.Comment)
	if err != nil {
		return err
	}

	d.SetId(repoId + "/" + file + ":" + string(*vdVersionType) + ":" + *vdVersion)

	return nil
}

// shortTagName removes the tag prefix which some API endpoints require.
func shortTagName(tag string) string {
	return strings.TrimPrefix(tag, "refs/tags/")
}
