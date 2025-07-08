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
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The file path",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"branch": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The branch name, no default",
				ExactlyOneOf: []string{"branch", "tag"},
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"tag": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Optional tag name, no default",
				ExactlyOneOf: []string{"branch", "tag"},
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The file's content",
			},
			"last_commit_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last commit message",
			},
		},
	}
}

func dataSourceGitRepositoryFileRead(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	clients := m.(*client.AggregatedClient)

	repoId := d.Get("repository_id").(string)
	file := d.Get("file").(string)

	var vDescriptor git.GitVersionDescriptor
	if v, ok := d.GetOk("branch"); ok {
		vDescriptor = git.GitVersionDescriptor{
			VersionType: &git.GitVersionTypeValues.Branch,
			Version:     converter.String(shortBranchName(v.(string))),
		}
	}
	if v, ok := d.GetOk("tag"); ok {
		vDescriptor = git.GitVersionDescriptor{
			VersionType: &git.GitVersionTypeValues.Tag,
			Version:     converter.String(shortTagName(v.(string))),
		}
	}

	repoItem, err := clients.GitReposClient.GetItem(ctx, git.GetItemArgs{
		RepositoryId:      &repoId,
		Path:              &file,
		IncludeContent:    converter.Bool(true),
		VersionDescriptor: &vDescriptor,
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return fmt.Errorf("Item not found, repositoryID: %s, %s: %s, file: %s. Error: %+v", repoId, string(*vDescriptor.VersionType), *vDescriptor.Version, file, err)
		}
		return fmt.Errorf("Get item failed, repositoryID: %s, %s: %s, file: %s. Error: %+v", repoId, string(*vDescriptor.VersionType), *vDescriptor.Version, file, err)
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
		return fmt.Errorf("Get commit failed, repositoryID: %s, commitID: %s. Error:  %+v", repoId, *repoItem.CommitId, err)
	}

	err = d.Set("last_commit_message", commit.Comment)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s:%s:%s", repoId, file, string(*vDescriptor.VersionType), *vDescriptor.Version))

	return nil
}

// shortTagName removes the tag prefix which some API endpoints require.
func shortTagName(tag string) string {
	return strings.TrimPrefix(tag, "refs/tags/")
}
