package azuredevops

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/operations"
)

var projectCreateTimeoutSeconds time.Duration = 60
var projectDeleteTimeoutSeconds time.Duration = 60

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		//https://godoc.org/github.com/hashicorp/terraform/helper/schema#Schema
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validate.NoEmptyStrings,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"visibility": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "private",
				ValidateFunc: validation.StringInSlice([]string{"private", "public"}, false),
			},
			"version_control": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				Default:      "Git",
				ValidateFunc: validation.StringInSlice([]string{"Git", "Tfvc"}, true),
			},
			"work_item_template": {
				Type:             schema.TypeString,
				ForceNew:         true,
				Optional:         true,
				ValidateFunc:     validate.NoEmptyStrings,
				DiffSuppressFunc: suppress.CaseDifference,
				Default:          "Agile",
			},
			"process_template_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	project, err := expandProject(clients, d, true)
	if err != nil {
		return fmt.Errorf("Error converting terraform data model to AzDO project reference: %+v", err)
	}

	err = createProject(clients, project, projectCreateTimeoutSeconds)
	if err != nil {
		return fmt.Errorf("Error creating project in Azure DevOps: %+v", err)
	}

	d.Set("project_name", *project.Name)
	return resourceProjectRead(d, m)
}

// Make API call to create the project and wait for an async success/fail response from the service
func createProject(clients *aggregatedClient, project *core.TeamProject, timeoutSeconds time.Duration) error {
	operationRef, err := clients.CoreClient.QueueCreateProject(clients.ctx, core.QueueCreateProjectArgs{ProjectToCreate: project})
	if err != nil {
		return err
	}

	return waitForAsyncOperationSuccess(clients, operationRef, timeoutSeconds)
}

func waitForAsyncOperationSuccess(clients *aggregatedClient, operationRef *operations.OperationReference, timeoutSeconds time.Duration) error {
	timeout := time.After(timeoutSeconds * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			result, err := clients.OperationsClient.GetOperation(clients.ctx, operations.GetOperationArgs{
				OperationId: operationRef.Id,
				PluginId:    operationRef.PluginId,
			})

			if err != nil {
				return err
			}

			if *result.Status == operations.OperationStatusValues.Succeeded {
				// Sometimes without the sleep, the subsequent operations won't find the project...
				delay := os.Getenv("AZDO_PRJ_CREATE_DELAY")
				settleDelay := time.Duration(0)
				i, err := strconv.ParseInt(delay, 10, 64)
				if err == nil {
					settleDelay = time.Duration(i) * time.Second
				}
				log.Printf("Inserting artificial delay after project creation: %s\n", settleDelay.String())
				time.Sleep(settleDelay)
				return nil
			}
		case <-timeout:
			return fmt.Errorf("Operation was not successful after %d seconds", timeoutSeconds)
		}
	}
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)

	id := d.Id()
	name := d.Get("project_name").(string)
	project, err := projectRead(clients, id, name)
	if err != nil {
		return fmt.Errorf("Error looking up project with ID %s and Name %s", id, name)
	}

	err = flattenProject(clients, d, project)
	if err != nil {
		return fmt.Errorf("Error flattening project: %v", err)
	}
	return nil
}

// Lookup a project using the ID, or name if the ID is not set. Note, usage of the name in place
// of the ID is an explicitly stated supported behavior:
//		https://docs.microsoft.com/en-us/rest/api/azure/devops/core/projects/get?view=azure-devops-rest-5.0
func projectRead(clients *aggregatedClient, projectID string, projectName string) (*core.TeamProject, error) {
	identifier := projectID
	if identifier == "" {
		identifier = projectName
	}

	return clients.CoreClient.GetProject(clients.ctx, core.GetProjectArgs{
		ProjectId:           &identifier,
		IncludeCapabilities: converter.Bool(true),
		IncludeHistory:      converter.Bool(false),
	})
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	project, err := expandProject(clients, d, false)
	if err != nil {
		return fmt.Errorf("Error converting terraform data model to AzDO project reference: %+v", err)
	}

	err = updateProject(clients, project, projectCreateTimeoutSeconds)
	if err != nil {
		return fmt.Errorf("Error updating project in Azure DevOps: %+v", err)
	}
	return resourceProjectRead(d, m)
}

func updateProject(clients *aggregatedClient, project *core.TeamProject, timeoutSeconds time.Duration) error {

	operationRef, err := clients.CoreClient.UpdateProject(
		clients.ctx,
		core.UpdateProjectArgs{
			ProjectUpdate: project,
			ProjectId:     project.Id,
		})

	if err != nil {
		return err
	}

	return waitForAsyncOperationSuccess(clients, operationRef, timeoutSeconds)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*aggregatedClient)
	id := d.Id()

	return deleteProject(clients, id, projectDeleteTimeoutSeconds)
}

func deleteProject(clients *aggregatedClient, id string, timeoutSeconds time.Duration) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("Invalid project UUID: %s", id)
	}

	operationRef, err := clients.CoreClient.QueueDeleteProject(clients.ctx, core.QueueDeleteProjectArgs{
		ProjectId: &uuid,
	})

	if err != nil {
		return err
	}

	return waitForAsyncOperationSuccess(clients, operationRef, timeoutSeconds)
}

// Convert internal Terraform data structure to an AzDO data structure
func expandProject(clients *aggregatedClient, d *schema.ResourceData, forCreate bool) (*core.TeamProject, error) {
	workItemTemplate := d.Get("work_item_template").(string)
	processTemplateID, err := lookupProcessTemplateID(clients, workItemTemplate)
	if err != nil {
		return nil, err
	}

	// an "error" is OK here as it is expected in the case that the ID is not set in the resource data
	var projectID *uuid.UUID
	parsedID, err := uuid.Parse(d.Id())
	if err == nil {
		projectID = &parsedID
	}

	visibility := d.Get("visibility").(string)

	var capabilities *map[string]map[string]string
	if forCreate {
		capabilities = &map[string]map[string]string{
			"versioncontrol": {
				"sourceControlType": d.Get("version_control").(string),
			},
			"processTemplate": {
				"templateTypeId": processTemplateID,
			},
		}
	}

	project := &core.TeamProject{
		Id:           projectID,
		Name:         converter.String(d.Get("project_name").(string)),
		Description:  converter.String(d.Get("description").(string)),
		Visibility:   convertVisibilty(visibility),
		Capabilities: capabilities,
	}

	return project, nil
}

func convertVisibilty(v string) *core.ProjectVisibility {
	if strings.ToLower(v) == "public" {
		return &core.ProjectVisibilityValues.Public
	}
	return &core.ProjectVisibilityValues.Private
}

func flattenProject(clients *aggregatedClient, d *schema.ResourceData, project *core.TeamProject) error {
	description := converter.ToString(project.Description, "")
	processTemplateID := (*project.Capabilities)["processTemplate"]["templateTypeId"]
	processTemplateName, err := lookupProcessTemplateName(clients, processTemplateID)

	if err != nil {
		return err
	}

	d.SetId(project.Id.String())
	d.Set("project_name", *project.Name)
	d.Set("visibility", *project.Visibility)
	d.Set("description", description)
	d.Set("version_control", (*project.Capabilities)["versioncontrol"]["sourceControlType"])
	d.Set("process_template_id", processTemplateID)
	d.Set("work_item_template", processTemplateName)

	return nil
}

// given a process template name, get the process template ID
func lookupProcessTemplateID(clients *aggregatedClient, templateName string) (string, error) {
	processes, err := clients.CoreClient.GetProcesses(clients.ctx, core.GetProcessesArgs{})
	if err != nil {
		return "", err
	}

	for _, p := range *processes {
		// Process names are case insensitive
		if strings.EqualFold(*p.Name, templateName) {
			return p.Id.String(), nil
		}
	}

	return "", fmt.Errorf("No process template found")
}

// given a process template ID, get the process template name
func lookupProcessTemplateName(clients *aggregatedClient, templateID string) (string, error) {
	id, err := uuid.Parse(templateID)
	if err != nil {
		return "", fmt.Errorf("Error parsing Work Item Template ID, got %s: %v", templateID, err)
	}

	process, err := clients.CoreClient.GetProcessById(clients.ctx, core.GetProcessByIdArgs{
		ProcessId: &id,
	})

	if err != nil {
		return "", fmt.Errorf("Error looking up template by ID: %v", err)
	}

	return *process.Name, nil
}
