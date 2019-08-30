package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
)

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourcePipelineCreate,
		Read:   resourcePipelineRead,
		Update: resourcePipelineUpdate,
		Delete: resourcePipelineDelete,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"pipeline_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"repository": &schema.Schema{
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
							Type:     schema.TypeString,
							Required: true,
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

type pipelineValues struct {
	projectID                     string
	pipelineName                  string
	repositoryName                string
	repositoryDefaultBranch       string
	repositoryPipelineYmlPath     string
	repositoryType                string
	repositoryURL                 string
	repositoryServiceConnectionID string
	agentPoolName                 string
	agentPoolID                   int
}

func resourcePipelineCreate(d *schema.ResourceData, m interface{}) error {
	repositories := d.Get("repository").(*schema.Set).List()
	repository := repositories[0].(map[string]interface{})
	repoName := repository["repo_name"].(string)
	pipelineName := d.Get("pipeline_name").(string)
	if pipelineName == "" {
		pipelineName = repoName + "_pipeline"
	}

	values := pipelineValues{
		projectID:                     d.Get("project_id").(string),
		pipelineName:                  pipelineName,
		repositoryName:                repoName,
		repositoryDefaultBranch:       repository["branch_name"].(string),
		repositoryPipelineYmlPath:     repository["yml_path"].(string),
		repositoryType:                repository["repo_type"].(string),
		repositoryURL:                 fmt.Sprintf("https://github.com/%s.git", repoName),
		repositoryServiceConnectionID: repository["service_connection_id"].(string),
		agentPoolName:                 "Hosted Ubuntu 1604",
		agentPoolID:                   224,
	}

	pipelineID, err := createPipeline(m.(*aggregatedClient), &values)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", pipelineID))
	return resourcePipelineRead(d, m)
}

func createPipeline(clients *aggregatedClient, values *pipelineValues) (int, error) {
	//get info from the client & create a build definition
	createRes, err := clients.BuildClient.CreateDefinition(clients.ctx, build.CreateDefinitionArgs{
		Definition: &build.BuildDefinition{
			Quality:     &build.DefinitionQualityValues.Definition,
			Name:        &values.pipelineName,
			Type:        &build.DefinitionTypeValues.Build,
			QueueStatus: &build.DefinitionQueueStatusValues.Enabled,
			Repository: &build.BuildRepository{
				Url:           &values.repositoryURL,
				Id:            &values.repositoryName,
				Name:          &values.repositoryName,
				DefaultBranch: &values.repositoryDefaultBranch,
				Type:          &values.repositoryType,
				Properties: &map[string]string{
					"connectedServiceId": values.repositoryServiceConnectionID,
				},
			},
			Process: &build.YamlProcess{
				YamlFilename: &values.repositoryPipelineYmlPath,
			},
			Queue: &build.AgentPoolQueue{
				Name: &values.agentPoolName,
				Pool: &build.TaskAgentPoolReference{
					Id:   &values.agentPoolID,
					Name: &values.agentPoolName,
				},
			},
		},
		Project: &values.projectID,
	})

	if err != nil {
		return 0, err
	}

	log.Printf("got response: %T", createRes)
	return *(createRes.Id), nil
}

func resourcePipelineRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePipelineDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePipelineUpdate(d *schema.ResourceData, m interface{}) error {
	return resourcePipelineRead(d, m)
}
