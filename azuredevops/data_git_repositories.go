package azuredevops

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/datahelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
)

func dataGitRepositories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitRepositoriesRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validate.UUID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"name": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validate.NoEmptyStrings,
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
					},
				},
			},
		},
	}
}

func dataSourceGitRepositoriesRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)

	projectRepos, err := getGitRepositoriesByNameAndProject(d, clients)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error finding repositories. Error: %v", err)
	}
	log.Printf("[TRACE] plugin.terraform-provider-azuredevops: Read [%d] Git repositories", len(*projectRepos))

	results, err := flattenGitRepositories(projectRepos)
	if err != nil {
		return fmt.Errorf("Error flattening projects. Error: %v", err)
	}

	repoNames, err := datahelper.GetAttributeValues(results, "name")
	if err != nil {
		return fmt.Errorf("Failed to get list of repository names: %v", err)
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
	if nil == repoNames {
		names = []string{}
	} else {
		names = *repoNames
	}
	if len(names) <= 0 {
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

func flattenGitRepositories(repos *[]git.GitRepository) ([]interface{}, error) {
	if repos == nil {
		return []interface{}{}, nil
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

		results = append(results, output)
	}

	return results, nil
}

func getGitRepositoriesByNameAndProject(d *schema.ResourceData, clients *config.AggregatedClient) (*[]git.GitRepository, error) {
	var repos *[]git.GitRepository
	var err error
	name, projectID := d.Get("name").(string), d.Get("project_id").(string)
	includeHidden := d.Get("include_hidden").(bool)

	if name != "" && projectID != "" {
		repo, err := gitRepositoryRead(clients, "", name, projectID)
		if err != nil {
			return nil, err
		}
		repos = &[]git.GitRepository{*repo}
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
