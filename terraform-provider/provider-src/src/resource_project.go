package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/operations"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		//https://godoc.org/github.com/hashicorp/terraform/helper/schema#Schema
		Schema: map[string]*schema.Schema{
			"project_name": &schema.Schema{
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				DiffSuppressFunc: suppressCaseSensitivity,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"visibility": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "private",
				ValidateFunc: validation.StringInSlice([]string{
					"private",
					"public",
				}, false),
			},
			"version_control": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Git",
				ValidateFunc: validation.StringInSlice([]string{
					"Git",
					"Tfvc",
				}, true),
			},
			"work_item_template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Agile",
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"process_template_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type projectValues struct {
	projectName        string
	description        string
	visibility         string
	versionControl     string
	workItemTemplate   string
	workItemTemplateID string
	projectID          string
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	//instantiate client
	clients := m.(*aggregatedClient)

	values := projectValues{
		projectName:      d.Get("project_name").(string),
		description:      d.Get("description").(string),
		visibility:       d.Get("visibility").(string),
		versionControl:   d.Get("version_control").(string),
		workItemTemplate: d.Get("work_item_template").(string),
	}

	// lookup process template id
	processTemplateID, err := lookupProcessTemplateID(clients, values.workItemTemplate)

	if err != nil {
		return fmt.Errorf("Invalid work item template name: %+v", err)
	}

	values.workItemTemplateID = processTemplateID

	// create project
	err = projectCreate(clients, &values)
	if err != nil {
		return fmt.Errorf("Error creating project in Azure DevOps: %+v", err)
	}

	//lookup project id
	projectID, err := lookupProjectID(clients, values.projectName)

	if err != nil {
		return fmt.Errorf("Error looking up project ID for project: %v, %+v", values.projectName, err)
	}

	values.projectID = projectID

	//call set id
	d.Set("process_template_id", values.workItemTemplateID)
	d.Set("project_id", values.projectID)
	d.SetId(values.projectName)

	//read project and return
	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	// Setup client
	clients := m.(*aggregatedClient)
	projectName := d.Id()

	// Get the Project
	projectValues, err := projectRead(clients, projectName)

	if err != nil {
		return fmt.Errorf("Error looking up project given ID: %v %v", projectName, err)
	}

	// Assign project values
	d.Set("project_name", projectValues.projectName)
	d.Set("description", projectValues.description)
	d.Set("visibility", projectValues.visibility)
	d.Set("version_control", projectValues.versionControl)
	d.Set("work_item_template", projectValues.workItemTemplate)
	d.Set("project_id", projectValues.projectID)
	d.Set("process_template_id", projectValues.workItemTemplateID)

	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func projectCreate(clients *aggregatedClient, values *projectValues) error {

	operationRef, err := clients.CoreClient.QueueCreateProject(clients.ctx, core.QueueCreateProjectArgs{
		ProjectToCreate: &core.TeamProject{
			Name:        &values.projectName,
			Description: &values.description,
			Visibility:  convertVisibilty(values.visibility),
			Capabilities: &map[string]map[string]string{
				"versioncontrol": map[string]string{
					"sourceControlType": values.versionControl,
				},
				"processTemplate": map[string]string{
					"templateTypeId": values.workItemTemplateID,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	err = waitForAsyncOperationSuccess(clients, operationRef)
	if err != nil {
		return err
	}

	return nil
}

func projectRead(clients *aggregatedClient, projectName string) (projectValues, error) {
	t := true
	f := false
	var pv projectValues

	p, err := clients.CoreClient.GetProject(clients.ctx, core.GetProjectArgs{
		ProjectId:           &projectName,
		IncludeCapabilities: &t,
		IncludeHistory:      &f,
	})

	if err != nil {
		return pv, fmt.Errorf("Error getting project: %v", err)
	}

	pv = projectValues{
		projectName:        *p.Name,
		description:        *p.Description,
		visibility:         string(*p.Visibility),
		versionControl:     (*p.Capabilities)["versioncontrol"]["sourceControlType"],
		workItemTemplateID: (*p.Capabilities)["processTemplate"]["templateTypeId"],
		projectID:          p.Id.String(),
	}

	templateID, err := uuid.Parse(pv.workItemTemplateID)
	if err != nil {
		return pv, fmt.Errorf("Error parsing Work Item Template ID, got %v: %v", pv.workItemTemplateID, err)
	}

	process, err := clients.CoreClient.GetProcessById(clients.ctx, core.GetProcessByIdArgs{
		ProcessId: &templateID,
	})
	if err != nil {
		return pv, fmt.Errorf("Error looking up template by ID: %v", err)
	}
	pv.workItemTemplate = *process.Name

	return pv, nil
}

func lookupProjectID(clients *aggregatedClient, projectName string) (string, error) {
	projects, err := clients.CoreClient.GetProjects(clients.ctx, core.GetProjectsArgs{})
	if err != nil {
		return "", err
	}

	for _, project := range projects.Value {
		if *project.Name == projectName {
			return project.Id.String(), nil
		}
	}

	return "", fmt.Errorf("No project found")
}

func lookupProcessTemplateID(clients *aggregatedClient, processTemplateName string) (string, error) {
	processes, err := clients.CoreClient.GetProcesses(clients.ctx, core.GetProcessesArgs{})
	if err != nil {
		return "", err
	}

	for _, p := range *processes {
		if *p.Name == processTemplateName {
			return p.Id.String(), nil
		}
	}

	return "", fmt.Errorf("No process template found")
}

func waitForAsyncOperationSuccess(clients *aggregatedClient, operationRef *operations.OperationReference) error {
	maxAttempts := 30
	currentAttempt := 1

	for currentAttempt <= maxAttempts {
		//log.Printf("Checking status for operation with ID: %s", operationRef.Id)
		result, err := clients.OperationsClient.GetOperation(clients.ctx, operations.GetOperationArgs{
			OperationId: operationRef.Id,
			PluginId:    operationRef.PluginId,
		})

		if err != nil {
			return err
		}

		if *result.Status == operations.OperationStatusValues.Succeeded {
			// Sometimes without the sleep, the subsequent operations won't find the project...
			time.Sleep(2 * time.Second)
			return nil
		}

		currentAttempt++
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("Operation was not successful after %d attempts", maxAttempts)
}

func convertVisibilty(v string) *core.ProjectVisibility {
	if strings.ToLower(v) == "public" {
		return &core.ProjectVisibilityValues.Public
	}
	return &core.ProjectVisibilityValues.Private
}

func suppressCaseSensitivity(k, old, new string, d *schema.ResourceData) bool {
	if strings.ToLower(old) == strings.ToLower(new) {
		return true
	}
	return false
}
