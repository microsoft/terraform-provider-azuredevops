package git

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

func ResourceGitRepositoryFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitRepositoryFileCreate,
		Read:   resourceGitRepositoryFileRead,
		Update: resourceGitRepositoryFileUpdate,
		Delete: resourceGitRepositoryFileDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), ":")
				branch := "main"

				if len(parts) > 2 {
					return nil, fmt.Errorf("Invalid ID specified. Supplied ID must be written as <repository>/<file path> (when branch is \"master\") or <repository>/<file path>:<branch>")
				}

				if len(parts) == 2 {
					branch = parts[1]
				}

				clients := m.(*client.AggregatedClient)
				repo, file := splitRepoFilePath(parts[0])
				if err := checkRepositoryFileExists(clients, repo, file, branch); err != nil {
					return nil, err
				}

				d.SetId(fmt.Sprintf("%s/%s", repo, file))
				d.Set("branch", branch)
				d.Set("overwrite_on_create", false)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"repository_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The repository name",
			},
			"file": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The file path to manage",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file's content",
			},
			"branch": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The branch name, defaults to \"master\"",
				Default:     "refs/heads/master",
			},
			"commit_message": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The commit message when creating or updating the file",
			},
			"overwrite_on_create": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable overwriting existing files, defaults to \"false\"",
				Default:     false,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Second),
		},
	}
}

func resourceGitRepositoryPushArgs(d *schema.ResourceData, objectID string, changeType git.VersionControlChangeType) (*git.CreatePushArgs, error) {
	var message *string
	if commitMessage, hasCommitMessage := d.GetOk("commit_message"); hasCommitMessage {
		cm := commitMessage.(string)
		message = &cm
	}

	repo := d.Get("repository_id").(string)
	content := d.Get("content").(string)
	file := d.Get("file").(string)
	branch := d.Get("branch").(string)

	change := git.GitChange{
		ChangeType: &changeType,
		Item: git.GitItem{
			Path: &file,
		},
		NewContent: &git.ItemContent{
			Content:     &content,
			ContentType: &git.ItemContentTypeValues.RawText,
		},
	}

	args := &git.CreatePushArgs{
		RepositoryId: &repo,
		Push: &git.GitPush{
			RefUpdates: &[]git.GitRefUpdate{
				{
					Name:        &branch,
					OldObjectId: &objectID,
				},
			},
			Commits: &[]git.GitCommitRef{
				{
					Comment: message,
					Changes: &[]interface{}{change},
				},
			},
		},
	}

	return args, nil
}

func resourceGitRepositoryFileCreate(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	clients := m.(*client.AggregatedClient)

	repo := d.Get("repository_id").(string)
	file := d.Get("file").(string)
	branch := d.Get("branch").(string)
	overwriteOnCreate := d.Get("overwrite_on_create").(bool)

	if err := checkRepositoryBranchExists(clients, repo, branch); err != nil {
		return err
	}

	changeType := git.VersionControlChangeTypeValues.Add
	item, err := clients.GitReposClient.GetItem(ctx, git.GetItemArgs{
		RepositoryId: &repo,
		Path:         &file,
	})
	if err != nil && !utils.ResponseWasNotFound(err) {
		return err
	}

	if item != nil {
		if !overwriteOnCreate {
			return fmt.Errorf("Refusing to overwrite existing file. Configure `overwrite_on_create` to `true` to override.")
		} else {
			changeType = git.VersionControlChangeTypeValues.Edit
		}
	}

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		objectID, err := getLastCommitId(clients, repo, branch)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		args, err := resourceGitRepositoryPushArgs(d, objectID, changeType)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		if (*args.Push.Commits)[0].Comment == nil {
			m := fmt.Sprintf("Add %s", file)
			(*args.Push.Commits)[0].Comment = &m
		}

		_, err = clients.GitReposClient.CreatePush(ctx, *args)
		if err != nil {
			if utils.ResponseContainsStatusMessage(err, "has already been updated by another client") {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", repo, file))
	return resourceGitRepositoryFileRead(d, m)
}

func resourceGitRepositoryFileRead(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	clients := m.(*client.AggregatedClient)

	repo, file := splitRepoFilePath(d.Id())
	branch := d.Get("branch").(string)

	if err := checkRepositoryBranchExists(clients, repo, branch); err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		branch = strings.TrimPrefix(branch, "refs/heads/")
		item, err := clients.GitReposClient.GetItem(ctx, git.GetItemArgs{
			RepositoryId:   &repo,
			Path:           &file,
			IncludeContent: boolPointer(true),
			VersionDescriptor: &git.GitVersionDescriptor{
				Version:     &branch,
				VersionType: &git.GitVersionTypeValues.Branch,
			},
		})
		if err != nil {
			if utils.ResponseWasNotFound(err) {
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		d.Set("content", item.Content)
		d.Set("repository_id", repo)
		d.Set("file", file)

		commit, err := clients.GitReposClient.GetCommit(ctx, git.GetCommitArgs{
			RepositoryId: &repo,
			CommitId:     item.CommitId,
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}

		d.Set("commit_message", commit.Comment)

		return nil
	})
}

func resourceGitRepositoryFileUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	ctx := context.Background()

	repo := d.Get("repository_id").(string)
	file := d.Get("file").(string)
	branch := d.Get("branch").(string)

	if err := checkRepositoryBranchExists(clients, repo, branch); err != nil {
		return err
	}

	objectID, err := getLastCommitId(clients, repo, branch)
	if err != nil {
		return err
	}

	args, err := resourceGitRepositoryPushArgs(d, objectID, git.VersionControlChangeTypeValues.Edit)
	if err != nil {
		return err
	}

	if *(*args.Push.Commits)[0].Comment == fmt.Sprintf("Add %s", file) {
		m := fmt.Sprintf("Update %s", file)
		(*args.Push.Commits)[0].Comment = &m
	}

	_, err = clients.GitReposClient.CreatePush(ctx, *args)
	if err != nil {
		return err
	}

	return resourceGitRepositoryFileRead(d, m)
}

func resourceGitRepositoryFileDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	ctx := context.Background()

	repo := d.Get("repository_id").(string)
	file := d.Get("file").(string)
	branch := d.Get("branch").(string)
	message := fmt.Sprintf("Delete %s", file)

	objectID, err := getLastCommitId(clients, repo, branch)
	if err != nil {
		return err
	}

	change := &git.GitChange{
		ChangeType: &git.VersionControlChangeTypeValues.Delete,
		Item: git.GitItem{
			Path: &file,
		},
	}

	_, err = clients.GitReposClient.CreatePush(ctx, git.CreatePushArgs{
		RepositoryId: &repo,
		Push: &git.GitPush{
			RefUpdates: &[]git.GitRefUpdate{
				{
					Name:        &branch,
					OldObjectId: &objectID,
				},
			},
			Commits: &[]git.GitCommitRef{
				{
					Comment: &message,
					Changes: &[]interface{}{change},
				},
			},
		},
	})
	if err != nil {
		return nil
	}

	return nil
}

// checkRepositoryBranchExists tests if a branch exists in a repository.
func checkRepositoryBranchExists(c *client.AggregatedClient, repo, branch string) error {
	branch = strings.TrimPrefix(branch, "refs/heads/")
	ctx := context.Background()
	_, err := c.GitReposClient.GetBranch(ctx, git.GetBranchArgs{
		RepositoryId: &repo,
		Name:         &branch,
	})
	return err
}

// checkRepositoryFileExists tests if a file exists in a repository.
func checkRepositoryFileExists(c *client.AggregatedClient, repo, file, branch string) error {
	ctx := context.Background()
	_, err := c.GitReposClient.GetItem(ctx, git.GetItemArgs{
		RepositoryId: &repo,
		Path:         &file,
	})
	if err != nil {
		return nil
	}

	return nil
}

func getLastCommitId(c *client.AggregatedClient, repo, branch string) (string, error) {
	ctx := context.Background()
	commits, err := c.GitReposClient.GetCommits(ctx, git.GetCommitsArgs{
		RepositoryId:   &repo,
		Top:            intPointer(1),
		SearchCriteria: &git.GitQueryCommitsCriteria{},
	})
	if err != nil {
		return "", err
	}
	return *(*commits)[0].CommitId, nil
}

func boolPointer(val bool) *bool {
	return &val
}

func intPointer(val int) *int {
	return &val
}

func splitRepoFilePath(path string) (string, string) {
	parts := strings.Split(path, "/")
	return parts[0], strings.Join(parts[1:], "/")
}
