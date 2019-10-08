package azuredevops

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
)

func resourceBuildDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceBuildDefinitionCreate,
		Read:   resourceBuildDefinitionRead,
		Update: resourceBuildDefinitionUpdate,
		Delete: resourceBuildDefinitionDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"revision": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"agent_pool_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Hosted Ubuntu 1604",
			},
			"repository": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"yml_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"repo_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"repo_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"GitHub"}, false),
						},
						"branch_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "master",
						},
						"service_connection_id": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
		},
	}
}

func resourceBuildDefinitionCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	buildDefinition, projectID, err := expandBuildDefinition(d)
	if err != nil {
		return err
	}

	createdBuildDefinition, err := createBuildDefinition(clients, buildDefinition, projectID)
	if err != nil {
		return err
	}

	flattenBuildDefinition(d, createdBuildDefinition, projectID)
	return nil
}

func flattenBuildDefinition(d *schema.ResourceData, buildDefinition *build.BuildDefinition, projectID string) {
	d.SetId(strconv.Itoa(*buildDefinition.Id))

	d.Set("project_id", projectID)
	d.Set("name", *buildDefinition.Name)
	d.Set("repository", flattenRepository(buildDefinition))
	d.Set("agent_pool_name", *buildDefinition.Queue.Pool.Name)

	revision := 0
	if buildDefinition.Revision != nil {
		revision = *buildDefinition.Revision
	}

	d.Set("revision", revision)
}

func createBuildDefinition(clients *aggregatedClient, buildDefinition *build.BuildDefinition, project string) (*build.BuildDefinition, error) {
	createdBuild, err := clients.BuildClient.CreateDefinition(clients.ctx, build.CreateDefinitionArgs{
		Definition: buildDefinition,
		Project:    &project,
	})

	return createdBuild, err
}

func resourceBuildDefinitionRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	projectID, buildDefinitionID, err := parseIdentifiers(d)

	if err != nil {
		return err
	}

	buildDefinition, err := clients.BuildClient.GetDefinition(clients.ctx, build.GetDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &buildDefinitionID,
	})

	if err != nil {
		return err
	}

	flattenBuildDefinition(d, buildDefinition, projectID)
	return nil
}

func resourceBuildDefinitionDelete(d *schema.ResourceData, m interface{}) error {
	if d.Id() == "" {
		return nil
	}

	clients := m.(*aggregatedClient)
	projectID, buildDefinitionID, err := parseIdentifiers(d)
	if err != nil {
		return err
	}

	err = clients.BuildClient.DeleteDefinition(m.(*aggregatedClient).ctx, build.DeleteDefinitionArgs{
		Project:      &projectID,
		DefinitionId: &buildDefinitionID,
	})

	return err
}

func resourceBuildDefinitionUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	buildDefinition, projectID, err := expandBuildDefinition(d)
	if err != nil {
		return err
	}

	updatedBuildDefinition, err := clients.BuildClient.UpdateDefinition(m.(*aggregatedClient).ctx, build.UpdateDefinitionArgs{
		Definition:   buildDefinition,
		Project:      &projectID,
		DefinitionId: buildDefinition.Id,
	})

	if err != nil {
		return err
	}

	flattenBuildDefinition(d, updatedBuildDefinition, projectID)
	return nil
}

func parseIdentifiers(d *schema.ResourceData) (string, int, error) {
	projectID := d.Get("project_id").(string)
	buildDefinitionID, err := strconv.Atoi(d.Id())

	return projectID, buildDefinitionID, err
}

func flattenRepository(buildDefiniton *build.BuildDefinition) interface{} {
	process := buildDefiniton.Process.(map[string]interface{})
	return []map[string]interface{}{{
		"yml_path":              process["yamlFilename"].(string),
		"repo_name":             *buildDefiniton.Repository.Name,
		"repo_type":             *buildDefiniton.Repository.Type,
		"branch_name":           *buildDefiniton.Repository.DefaultBranch,
		"service_connection_id": (*buildDefiniton.Repository.Properties)["connectedServiceId"],
	}}
}

func expandBuildDefinition(d *schema.ResourceData) (*build.BuildDefinition, string, error) {
	projectID := d.Get("project_id").(string)
	repositories := d.Get("repository").(*schema.Set).List()
	repository := repositories[0].(map[string]interface{})

	repoName := repository["repo_name"].(string)
	repoType := repository["repo_type"].(string)
	repoURL := ""
	if strings.EqualFold(repoType, "github") {
		repoURL = fmt.Sprintf("https://github.com/%s.git", repoName)
	}

	// Look for the ID. This may not exist if we are within the context of a "create" operation,
	// so it is OK if it is missing.
	buildDefinitionID, err := strconv.Atoi(d.Id())
	var buildDefinitionReference *int
	if err == nil {
		buildDefinitionReference = &buildDefinitionID
	} else {
		buildDefinitionReference = nil
	}

	agentPoolName := d.Get("agent_pool_name").(string)
	buildDefinition := build.BuildDefinition{
		Id:       buildDefinitionReference,
		Name:     converter.String(d.Get("name").(string)),
		Revision: converter.Int(d.Get("revision").(int)),
		Repository: &build.BuildRepository{
			Url:           &repoURL,
			Id:            &repoName,
			Name:          &repoName,
			DefaultBranch: converter.String(repository["branch_name"].(string)),
			Type:          &repoType,
			Properties: &map[string]string{
				"connectedServiceId": repository["service_connection_id"].(string),
			},
		},
		Process: &build.YamlProcess{
			YamlFilename: converter.String(repository["yml_path"].(string)),
		},
		Queue: &build.AgentPoolQueue{
			Name: &agentPoolName,
			Pool: &build.TaskAgentPoolReference{
				Name: &agentPoolName,
			},
		},
		QueueStatus: &build.DefinitionQueueStatusValues.Enabled,
		Type:        &build.DefinitionTypeValues.Build,
		Quality:     &build.DefinitionQualityValues.Definition,
	}

	return &buildDefinition, projectID, nil
}
