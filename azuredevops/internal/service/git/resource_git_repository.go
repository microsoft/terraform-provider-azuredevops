package git

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// RepoInitType strategy for initializing the repo
type RepoInitType string

type repoInitTypeValuesType struct {
	Uninitialized RepoInitType
	Clean         RepoInitType
	Fork          RepoInitType
	Import        RepoInitType
}

// RepoInitTypeValues enum of strategy for initializing the repo
var RepoInitTypeValues = repoInitTypeValuesType{
	Uninitialized: "Uninitialized",
	Clean:         "Clean",
	Fork:          "Fork",
	Import:        "Import",
}

// A helper type that is used for transient info only used during repo creation
type repoInitializationMeta struct {
	initType            string
	sourceType          string
	sourceURL           string
	serviceConnectionID string
}

// ResourceGitRepository schema and implementation for git repo resource
func ResourceGitRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitRepositoryCreate,
		Read:   resourceGitRepositoryRead,
		Update: resourceGitRepositoryUpdate,
		Delete: resourceGitRepositoryDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResource(),
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true, // repositories cannot be moved
				ValidateFunc:     validation.NoZeroValues,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"name": {
				Type:             schema.TypeString,
				ForceNew:         false,
				Required:         true,
				ValidateFunc:     validation.NoZeroValues,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"initialization": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"init_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(RepoInitTypeValues.Clean),
								string(RepoInitTypeValues.Fork),
								string(RepoInitTypeValues.Import),
								string(RepoInitTypeValues.Uninitialized),
							}, false),
						},
						"source_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"Git"}, false),
							RequiredWith: []string{"initialization.0.source_url"},
						},
						"source_url": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "",
							RequiredWith: []string{"initialization.0.source_type"},
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
						},
						"service_connection_id": {
							Type:     schema.TypeString,
							Optional: true,
							RequiredWith: []string{
								"initialization.0.source_url",
								"initialization.0.source_type",
							},
							Default: "",
						},
					},
				},
			},
			"parent_repository_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.IsUUID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"default_branch": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"is_fork": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"remote_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ssh_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"web_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}
func resourceGitRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	repo, initialization, projectID, err := expandGitRepository(d)
	if err != nil {
		return fmt.Errorf(" failed expanding repository resource data (ProjectID:  %s, Repository: %s) Error: %+v",
			d.Get("project_id").(string), d.Get("name").(string), err)
	}

	if _, ok := d.GetOk("default_branch"); ok {
		if strings.EqualFold(initialization.initType, string(RepoInitTypeValues.Uninitialized)) {
			return fmt.Errorf(" Repository 'initialization.init_type = Uninitialized', there will be no branches, 'default_branch' cannot not be set.")
		}
	}

	var parentRepoRef *git.GitRepositoryRef = nil
	if parentRepoID, ok := d.GetOk("parent_repository_id"); ok {
		parentRepo, err := gitRepositoryRead(clients, parentRepoID.(string), "", "")
		if err != nil {
			return fmt.Errorf("Failed to locate parent repository [%s]: %+v", parentRepoID, err)
		}
		parentRepoRef = &git.GitRepositoryRef{
			Id:      parentRepo.Id,
			Name:    parentRepo.Name,
			Project: parentRepo.Project,
		}
	}

	createdRepo, err := createGitRepository(clients, repo.Name, projectID, parentRepoRef)
	if err != nil {
		return fmt.Errorf(" Creating repository in Azure DevOps: %+v", err)
	}

	// set the id immediately after successfully creating the repository, which will allow terraform to track the
	// repository in state and will taint the resource if further errors are encountered during initialization
	d.SetId(createdRepo.Id.String())

	if initialization != nil {
		if strings.EqualFold(initialization.initType, string(RepoInitTypeValues.Import)) && strings.EqualFold(initialization.sourceType, "Git") {
			importRequest := git.GitImportRequest{
				Parameters: &git.GitImportRequestParameters{
					GitSource: &git.GitImportGitSource{
						Url: &initialization.sourceURL,
					},
				},
				Repository: createdRepo,
			}

			if initialization.serviceConnectionID != "" {
				importRequest.Parameters.ServiceEndpointId = converter.UUID(initialization.serviceConnectionID)
				importRequest.Parameters.DeleteServiceEndpointAfterImportIsDone = converter.Bool(false)
			}

			_, importErr := createImportRequest(clients, importRequest, projectID.String(), *createdRepo.Name)
			if importErr != nil {
				return fmt.Errorf("Error import repository in Azure DevOps: %+v ", importErr)
			}
		}

		if strings.EqualFold(initialization.initType, string(RepoInitTypeValues.Clean)) ||
			strings.EqualFold(initialization.initType, string(RepoInitTypeValues.Fork)) {
			err = initializeGitRepository(clients, createdRepo, repo.DefaultBranch)
			if err != nil {
				return fmt.Errorf(" initializing repository in Azure DevOps: %+v ", err)
			}
		}
	}

	if !strings.EqualFold(initialization.initType, string(RepoInitTypeValues.Uninitialized)) || parentRepoRef != nil {
		err := waitForBranch(clients, repo.Name, projectID)
		if err != nil {
			return err
		}
	}

	// update default_branch
	if v := d.Get("default_branch").(string); v != "" {
		createdRepo.DefaultBranch = converter.String(v)
		_, err = updateGitRepository(clients, createdRepo, projectID)
		if err != nil {
			return fmt.Errorf(" updating repository : %+v", err)
		}
	}

	if v := d.Get("disabled").(bool); v {
		_, err = updateIsDisabledGitRepository(clients, createdRepo.Id.String(), projectID.String(), true)
		if err != nil {
			return fmt.Errorf("Error disabling created repository in Azure DevOps: %+v", err)
		}
	}

	return resourceGitRepositoryRead(d, m)
}

func resourceGitRepositoryRead(d *schema.ResourceData, m interface{}) error {
	repoID := d.Id()
	repoName := d.Get("name").(string)
	projectID := d.Get("project_id").(string)

	clients := m.(*client.AggregatedClient)
	repo, err := gitRepositoryRead(clients, repoID, repoName, projectID)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error looking up repository with ID %s and Name %s. Error: %v", repoID, repoName, err)
	}

	if repo == nil {
		d.SetId("")
		return nil
	}

	err = flattenGitRepository(d, repo)
	if err != nil {
		return fmt.Errorf("Failed to flatten Git repository: %w", err)
	}
	return nil
}

func resourceGitRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	repo, _, projectID, err := expandGitRepository(d)
	disabled := d.Get("disabled").(bool)
	if err != nil {
		return fmt.Errorf("Error converting terraform data model to AzDO project reference: %+v", err)
	}

	// you cannot update a disabled repo. We must temporarily enable to update
	repoActual, err := gitRepositoryRead(clients, repo.Id.String(), *repo.Name, projectID.String())
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error looking up repository with ID %s and Name %s. Error: %v", *repo.Id, *repo.Name, err)
	}

	// potentially enabling repo to match the config
	if *repoActual.IsDisabled && !disabled {
		repoActual, err = updateIsDisabledGitRepository(clients, repo.Id.String(), projectID.String(), disabled)
		if err != nil {
			return fmt.Errorf("Error enabling repository in Azure DevOps: %+v", err)
		}
	}

	if *repoActual.IsDisabled {
		return fmt.Errorf("A disabled repository cannot be updated, please enable the repository before attempting to update : %s", repo.Id.String())
	}

	_, err = updateGitRepository(clients, repo, projectID)
	if err != nil {
		return fmt.Errorf("Error updating repository in Azure DevOps: %+v", err)
	}

	// potentially disabling repo to match the config
	if !*repoActual.IsDisabled && disabled {
		_, err = updateIsDisabledGitRepository(clients, repo.Id.String(), projectID.String(), disabled)
		if err != nil {
			return fmt.Errorf("Error disabling repository in Azure DevOps: %+v", err)
		}
	}

	return resourceGitRepositoryRead(d, m)
}

func resourceGitRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	repoID := d.Id()
	repoName := d.Get("name").(string)
	projectID := d.Get("project_id").(string)
	clients := m.(*client.AggregatedClient)

	uuid, err := uuid.Parse(repoID)
	if err != nil {
		return fmt.Errorf("Invalid repositoryId UUID: %s", repoID)
	}

	// you cannot delete a disabled repo. We must enable in order to delete
	repoActual, err := gitRepositoryRead(clients, repoID, repoName, projectID)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error looking up repository with ID %s and Name %s. Error: %v", repoID, repoName, err)
	}

	if *repoActual.IsDisabled {
		return fmt.Errorf("A disabled repository cannot be deleted, please enable the repository before attempting to delete : %s", repoID)
	}

	err = deleteGitRepository(clients, uuid)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func waitForBranch(clients *client.AggregatedClient, repoName *string, projectID fmt.Stringer) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Waiting"},
		Target:  []string{"Synched"},
		Refresh: func() (interface{}, string, error) {
			state := "Waiting"
			gitRepo, err := gitRepositoryRead(clients, "", *repoName, projectID.String())
			if err != nil {
				return nil, "", fmt.Errorf("Error reading repository: %+v", err)
			}

			if converter.ToString(gitRepo.DefaultBranch, "") != "" {
				state = "Synched"
			}

			return state, state, nil
		},
		Timeout:                   60 * time.Second,
		MinTimeout:                2 * time.Second,
		Delay:                     1 * time.Second,
		ContinuousTargetOccurence: 1,
	}
	if _, err := stateConf.WaitForState(); err != nil { //nolint:staticcheck
		return fmt.Errorf("Error retrieving expected branch for repository [%s]: %+v", *repoName, err)
	}
	return nil
}

func createImportRequest(clients *client.AggregatedClient, gitImportRequest git.GitImportRequest, project string, repositoryID string) (*git.GitImportRequest, error) {
	args := git.CreateImportRequestArgs{
		ImportRequest: &gitImportRequest,
		Project:       &project,
		RepositoryId:  &repositoryID,
	}

	return clients.GitReposClient.CreateImportRequest(clients.Ctx, args)
}

func createGitRepository(clients *client.AggregatedClient, repoName *string, projectID *uuid.UUID, parentRepo *git.GitRepositoryRef) (*git.GitRepository, error) {
	args := git.CreateRepositoryArgs{
		GitRepositoryToCreate: &git.GitRepositoryCreateOptions{
			Name: repoName,
			Project: &core.TeamProjectReference{
				Id: projectID,
			},
			ParentRepository: parentRepo,
		},
	}
	createdRepository, err := clients.GitReposClient.CreateRepository(clients.Ctx, args)
	if err != nil {
		return nil, err
	}

	return createdRepository, nil
}

func initializeGitRepository(clients *client.AggregatedClient, repo *git.GitRepository, defaultBranch *string) error {
	branchName := converter.ToString(defaultBranch, "")
	if strings.EqualFold(branchName, "") {
		branchName = "refs/heads/master"
	}
	args := git.CreatePushArgs{
		RepositoryId: repo.Name,
		Project:      repo.Project.Name,
		Push: &git.GitPush{
			RefUpdates: &[]git.GitRefUpdate{
				{
					Name:        converter.String(branchName),
					OldObjectId: converter.String("0000000000000000000000000000000000000000"),
				},
			},
			Commits: &[]git.GitCommitRef{
				{
					Comment: converter.String("Initial commit."),
					Changes: &[]interface{}{
						git.Change{
							ChangeType: &git.VersionControlChangeTypeValues.Add,
							Item: git.GitItem{
								Path: converter.String("/README.md"),
							},
							NewContent: &git.ItemContent{
								ContentType: &git.ItemContentTypeValues.RawText,
								Content:     repo.Project.Name,
							},
						},
					},
				},
			},
		},
	}

	_, err := clients.GitReposClient.CreatePush(clients.Ctx, args)

	return err
}

func updateGitRepository(clients *client.AggregatedClient, repository *git.GitRepository, project fmt.Stringer) (*git.GitRepository, error) {
	if nil == project {
		return nil, fmt.Errorf("updateGitRepository: ID of project cannot be nil")
	}
	projectID := project.String()
	return clients.GitReposClient.UpdateRepository(
		clients.Ctx,
		git.UpdateRepositoryArgs{
			NewRepositoryInfo: repository,
			RepositoryId:      repository.Id,
			Project:           &projectID,
		})
}

func deleteGitRepository(clients *client.AggregatedClient, uuid uuid.UUID) error {
	return clients.GitReposClient.DeleteRepository(clients.Ctx, git.DeleteRepositoryArgs{
		RepositoryId: &uuid,
	})
}

// Lookup an Azure Git Repository using the ID, or name if the ID is not set.
func gitRepositoryRead(clients *client.AggregatedClient, repoID string, repoName string, projectID string) (*git.GitRepository, error) {
	identifier := repoID
	if strings.EqualFold(identifier, "") {
		identifier = repoName
	}

	repo, err := clients.GitReposClient.GetRepository(clients.Ctx, git.GetRepositoryArgs{
		RepositoryId: converter.String(identifier),
		Project:      converter.String(projectID),
	})

	// If the repository is disabled, the repository cannot be obtained through the GET API
	if utils.ResponseWasNotFound(err) {
		var allRepo *[]git.GitRepository
		allRepo, err = clients.GitReposClient.GetRepositories(clients.Ctx, git.GetRepositoriesArgs{
			Project: converter.String(projectID),
			// This flag is used to include disabled repos
			IncludeHidden: converter.Bool(true),
		})
		if err != nil {
			return nil, err
		}
		for _, gitRepo := range *allRepo {
			if strings.EqualFold((*gitRepo.Id).String(), identifier) ||
				strings.EqualFold(*gitRepo.Name, identifier) {
				repo = &gitRepo
				break
			}
		}
	}
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func flattenGitRepository(d *schema.ResourceData, repository *git.GitRepository) error {
	d.SetId(repository.Id.String())
	d.Set("name", repository.Name)
	if repository.Project == nil || repository.Project.Id == nil {
		return fmt.Errorf("Unable to flatten Git repository without a valid projectID")
	}
	d.Set("project_id", repository.Project.Id.String())
	d.Set("default_branch", repository.DefaultBranch)
	d.Set("is_fork", repository.IsFork)
	d.Set("remote_url", repository.RemoteUrl)
	d.Set("size", repository.Size)
	d.Set("ssh_url", repository.SshUrl)
	d.Set("url", repository.Url)
	d.Set("web_url", repository.WebUrl)
	d.Set("disabled", repository.IsDisabled)

	return nil
}

// Convert internal Terraform data structure to an AzDO data structure. Note: only the params that are
// not generated by the service are expanded here
func expandGitRepository(d *schema.ResourceData) (*git.GitRepository, *repoInitializationMeta, *uuid.UUID, error) {
	// an "error" is OK here as it is expected in the case that the ID is not set in the resource data
	var repoID *uuid.UUID
	id := d.Id()
	if strings.EqualFold(id, "") {
		log.Print("[DEBUG] expandGitRepository: ID is empty (not set)")
	} else {
		parsedID, err := uuid.Parse(id)
		if err == nil {
			repoID = &parsedID
		}
	}

	projectID, err := uuid.Parse(d.Get("project_id").(string))
	if err != nil {
		return nil, nil, nil, err
	}

	repo := &git.GitRepository{
		Id:            repoID,
		Name:          converter.String(d.Get("name").(string)),
		DefaultBranch: converter.String(d.Get("default_branch").(string)),
	}

	var initialization *repoInitializationMeta = nil
	initData := d.Get("initialization").([]interface{})
	// Note: If configured, this will be of length 1 based on the schema definition above.
	if len(initData) == 1 {
		initValues := initData[0].(map[string]interface{})

		initialization = &repoInitializationMeta{
			initType:            initValues["init_type"].(string),
			sourceType:          initValues["source_type"].(string),
			sourceURL:           initValues["source_url"].(string),
			serviceConnectionID: initValues["service_connection_id"].(string),
		}

		if strings.EqualFold(initialization.initType, "clean") {
			initialization.sourceType = ""
			initialization.sourceURL = ""
			initialization.serviceConnectionID = ""
		}
	} else if len(initData) > 1 {
		return nil, nil, nil, fmt.Errorf("Multiple initialization blocks")
	}
	return repo, initialization, &projectID, nil
}

// When enabling or disabling a repo, isDisabled must be the only data in the post body
func updateIsDisabledGitRepository(clients *client.AggregatedClient, repoID string, projectID string, isDisabled bool) (*git.GitRepository, error) {
	uuid, err := uuid.Parse(repoID)
	if err != nil {
		return nil, fmt.Errorf("Invalid repositoryId UUID: %s", repoID)
	}
	repo, err := clients.GitReposClient.UpdateRepository(
		clients.Ctx,
		git.UpdateRepositoryArgs{
			NewRepositoryInfo: &git.GitRepository{IsDisabled: converter.Bool(isDisabled)},
			RepositoryId:      &uuid,
			Project:           converter.String(projectID),
		})
	if err != nil {
		return nil, fmt.Errorf(" updating isDisabled on repository : %+v", err)
	}

	return repo, nil
}
