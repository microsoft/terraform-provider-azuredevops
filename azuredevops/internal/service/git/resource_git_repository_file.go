package git

import (
	"context"
	"fmt"
	"strings"

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
			"repository": {
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
				Default:     "main",
			},
			"commit_message": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The commit message when creating or updating the file",
			},
			"commit_author": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The commit author name, defaults to the authenticated user's name",
			},
			"commit_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The commit author email address, defaults to the authenticated user's email address",
			},
			"object_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The blob object id of the file",
			},
			"overwrite_on_create": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable overwriting existing files, defaults to \"false\"",
				Default:     false,
			},
		},
	}
}

func resourceGitRepositoryPushArgs(d *schema.ResourceData, changeType git.VersionControlChangeType) (*git.CreatePushArgs, error) {
	var message string
	if commitMessage, hasCommitMessage := d.GetOk("commit_message"); hasCommitMessage {
		message = commitMessage.(string)
	}

	var objectID string
	if fileObjectID, hasObjectID := d.GetOk("object_id"); hasObjectID {
		objectID = fileObjectID.(string)
	}

	commitAuthor, hasCommitAuthor := d.GetOk("commit_author")
	commitEmail, hasCommitEmail := d.GetOk("commit_email")

	if hasCommitAuthor && !hasCommitEmail {
		return nil, fmt.Errorf("Cannot set commit_author without setting commit_email")
	}

	if hasCommitEmail && !hasCommitAuthor {
		return nil, fmt.Errorf("Cannot set commit_email without setting commit_author")
	}

	var author *git.GitUserDate
	if hasCommitAuthor && hasCommitEmail {
		name := commitAuthor.(string)
		email := commitEmail.(string)
		author = &git.GitUserDate{Name: &name, Email: &email}
	}

	repo := d.Get("repo").(string)
	content := d.Get("content").(string)
	file := d.Get("file").(string)
	branch := d.Get("branch").(string)

	change := &git.GitChange{
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
					Author:  author,
					Comment: &message,
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

	repo := d.Get("repository").(string)
	file := d.Get("file").(string)
	branch := d.Get("branch").(string)

	if err := checkRepositoryBranchExists(clients, repo, branch); err != nil {
		return err
	}

	args, err := resourceGitRepositoryPushArgs(d, git.VersionControlChangeTypeValues.Add)
	if err != nil {
		return err
	}

	if (*args.Push.Commits)[0].Comment == nil {
		m := fmt.Sprintf("Add %s", file)
		(*args.Push.Commits)[0].Comment = &m
	}

	item, err := clients.GitReposClient.GetItem(ctx, git.GetItemArgs{
		RepositoryId: &repo,
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) == false {
			return err
		}
	}

	if item.Content != nil {
		if d.Get("overwrite_on_create").(bool) {
			(*args.Push.RefUpdates)[0].OldObjectId = item.ObjectId
		} else {
			return fmt.Errorf("Refusing to overwrite existing file. Configure `overwrite_on_create` to `true` to override.")
		}
	}

	// Create a new or overwritten file
	_, err = clients.GitReposClient.CreatePush(ctx, *args)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", repo, file))

	return resourceGitRepositoryFileRead(d, m)
}

func boolPointer(val bool) *bool {
	return &val
}

func splitRepoFilePath(path string) (string, string) {
	parts := strings.Split(path, "/")
	return parts[0], strings.Join(parts[1:], "/")
}

func resourceGitRepositoryFileRead(d *schema.ResourceData, m interface{}) error {
	ctx := context.Background()
	clients := m.(*client.AggregatedClient)

	repo, file := splitRepoFilePath(d.Id())
	branch := d.Get("branch").(string)

	if err := checkRepositoryBranchExists(clients, repo, branch); err != nil {
		return err
	}

	item, err := clients.GitReposClient.GetItem(ctx, git.GetItemArgs{
		RepositoryId:      &repo,
		Path:              &file,
		VersionDescriptor: &git.GitVersionDescriptor{Version: &branch},
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("content", item.Content)
	d.Set("repository", repo)
	d.Set("file", file)
	d.Set("object_id", item.ObjectId)

	commit, err := clients.GitReposClient.GetCommit(ctx, git.GetCommitArgs{
		RepositoryId: &repo,
		CommitId:     item.CommitId,
	})
	if err != nil {
		return err
	}

	d.Set("commit_author", commit.Author.Name)
	d.Set("commit_email", commit.Author.Email)
	d.Set("commit_message", commit.Comment)

	return nil
}

func resourceGitRepositoryFileUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	ctx := context.Background()

	repo := d.Get("repository").(string)
	file := d.Get("file").(string)
	branch := d.Get("branch").(string)

	if err := checkRepositoryBranchExists(clients, repo, branch); err != nil {
		return err
	}

	args, err := resourceGitRepositoryPushArgs(d, git.VersionControlChangeTypeValues.Edit)
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

	repo := d.Get("repository").(string)
	file := d.Get("file").(string)
	branch := d.Get("branch").(string)
	objectID := d.Get("object_id").(string)
	message := fmt.Sprintf("Delete %s", file)

	change := &git.GitChange{
		ChangeType: &git.VersionControlChangeTypeValues.Delete,
		Item: git.GitItem{
			Path: &file,
		},
	}

	_, err := clients.GitReposClient.CreatePush(ctx, git.CreatePushArgs{
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
	ctx := context.Background()
	c.GitReposClient.GetBranch(ctx, git.GetBranchArgs{
		RepositoryId: &repo,
		Name:         &branch,
	})

	return nil
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
