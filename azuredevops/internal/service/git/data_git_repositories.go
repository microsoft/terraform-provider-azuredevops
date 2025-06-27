package git

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/datahelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/suppress"
)

// DataGitRepositories schema and implementation for git repo data source
func DataGitRepositories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitRepositoriesRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.IsUUID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsNotWhiteSpace,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"include_hidden": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"repositories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ssh_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"web_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remote_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_branch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGitRepositoriesRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	name := d.Get("name").(string)
	projectID := d.Get("project_id").(string)
	includeHidden := d.Get("include_hidden").(bool)

	projectRepos, err := getGitRepositoriesByNameAndProject(clients, name, projectID, includeHidden)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			return nil
		}
		return fmt.Errorf("finding repositories. Error: %v", err)
	}

	results := flattenGitRepositories(projectRepos)
	repoNames, err := datahelper.GetAttributeValues(results, "name")
	if err != nil {
		return fmt.Errorf("failed to get list of repository names: %v", err)
	}

	id, err := createGitRepositoryDataSourceID(d, &repoNames)
	if err != nil {
		return err
	}

	d.SetId(id)
	err = d.Set("repositories", results)
	if err != nil {
		d.SetId("")
		return err
	}
	return nil
}

func createGitRepositoryDataSourceID(d *schema.ResourceData, repoNames *[]string) (string, error) {
	h := sha1.New()
	var names []string
	if repoNames != nil {
		names = *repoNames
	}
	if len(names) == 0 {
		names = append(names, "empty")
	}
	projectID := d.Get("project_id").(string)
	if projectID != "" {
		names = append([]string{projectID}, names...)
	}
	if _, err := h.Write([]byte(strings.Join(names, "-"))); err != nil {
		return "", fmt.Errorf("Unable to compute hash for Git repository names: %v", err)
	}
	return "gitRepos#" + base64.URLEncoding.EncodeToString(h.Sum(nil)), nil
}

func flattenGitRepositories(repos *[]git.GitRepository) []interface{} {
	if repos == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)

	for _, element := range *repos {
		output := make(map[string]interface{})
		if element.Name != nil {
			output["name"] = *element.Name
		}

		if element.Id != nil {
			output["id"] = element.Id.String()
		}

		if element.Url != nil {
			output["url"] = *element.Url
		}

		if element.RemoteUrl != nil {
			output["remote_url"] = *element.RemoteUrl
		}

		if element.SshUrl != nil {
			output["ssh_url"] = *element.SshUrl
		}

		if element.WebUrl != nil {
			output["web_url"] = *element.WebUrl
		}

		if element.Project != nil && element.Project.Id != nil {
			output["project_id"] = element.Project.Id.String()
		}

		if element.Size != nil {
			output["size"] = *element.Size
		}

		if element.DefaultBranch != nil {
			output["default_branch"] = *element.DefaultBranch
		}

		if element.IsDisabled != nil {
			output["disabled"] = *element.IsDisabled
		}

		results = append(results, output)
	}

	return results
}

func getGitRepositoriesByNameAndProject(clients *client.AggregatedClient, name string, projectID string, includeHidden bool) (*[]git.GitRepository, error) {
	var repos *[]git.GitRepository
	var err error

	if name != "" && projectID != "" {
		repo, err := gitRepositoryRead(clients, "", name, projectID)
		if err != nil {
			return nil, err
		}

		if repo != nil {
			repos = &[]git.GitRepository{*repo}
		}
	} else {
		repos, err = clients.GitReposClient.GetRepositories(clients.Ctx, git.GetRepositoriesArgs{
			Project:       converter.String(projectID),
			IncludeHidden: converter.Bool(includeHidden),
		})
		if err != nil {
			return nil, err
		}
		if name != "" {
			for _, repo := range *repos {
				if strings.EqualFold(*repo.Name, name) {
					repos = &[]git.GitRepository{repo}
					break
				}
			}
		}
	}
	return repos, nil
}
