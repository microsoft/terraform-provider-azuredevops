package azuredevops

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/operations"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/suppress"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/validate"
)

var projectCreateTimeoutDuration time.Duration = 60 * 3
var projectDeleteTimeoutDuration time.Duration = 60

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
				Type:     schema.TypeString,
				Optional: true,
				Default:  core.ProjectVisibilityValues.Private,
				ValidateFunc: validation.StringInSlice([]string{
					string(core.ProjectVisibilityValues.Private),
					string(core.ProjectVisibilityValues.Public),
				}, false),
			},
			"version_control": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  core.SourceControlTypesValues.Git,
				ValidateFunc: validation.StringInSlice([]string{
					string(core.SourceControlTypesValues.Git),
					string(core.SourceControlTypesValues.Tfvc),
				}, true),
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
	clients := m.(*config.AggregatedClient)
	project, err := expandProject(clients, d, true)
	if err != nil {
		return fmt.Errorf("Error converting terraform data model to Azure DevOps project reference: %+v", err)
	}

	err = createProject(clients, project, projectCreateTimeoutDuration)
	if err != nil {
		return fmt.Errorf("Error creating project: %v", err)
	}

	d.Set("project_name", *project.Name)
	return resourceProjectRead(d, m)
}

// Make API call to create the project and wait for an async success/fail response from the service
func createProject(clients *config.AggregatedClient, project *core.TeamProject, timeoutSeconds time.Duration) error {
	operationRef, err := clients.CoreClient.QueueCreateProject(clients.Ctx, core.QueueCreateProjectArgs{ProjectToCreate: project})
	if err != nil {
		return err
	}

	return waitForAsyncOperationSuccess(clients, operationRef, timeoutSeconds)
}

func waitForAsyncOperationSuccess(clients *config.AggregatedClient, operationRef *operations.OperationReference, timeoutSeconds time.Duration) error {
	timeout := time.After(timeoutSeconds * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			result, err := clients.OperationsClient.GetOperation(clients.Ctx, operations.GetOperationArgs{
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
	clients := m.(*config.AggregatedClient)

	id := d.Id()
	name := d.Get("project_name").(string)
	project, err := ProjectRead(clients, id, name)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error looking up project with ID %s and Name %s", id, name)
	}

	err = flattenProject(clients, d, project)
	if err != nil {
		return fmt.Errorf("Error flattening project: %v", err)
	}
	return nil
}

// ProjectRead Lookup a project using the ID, or name if the ID is not set. Note, usage of the name in place
// of the ID is an explicitly stated supported behavior:
//		https://docs.microsoft.com/en-us/rest/api/azure/devops/core/projects/get?view=azure-devops-rest-5.0
func ProjectRead(clients *config.AggregatedClient, projectID string, projectName string) (*core.TeamProject, error) {
	identifier := projectID
	if identifier == "" {
		identifier = projectName
	}

	return clients.CoreClient.GetProject(clients.Ctx, core.GetProjectArgs{
		ProjectId:           &identifier,
		IncludeCapabilities: converter.Bool(true),
		IncludeHistory:      converter.Bool(false),
	})
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*config.AggregatedClient)
	project, err := expandProject(clients, d, false)
	if err != nil {
		return fmt.Errorf("Error converting terraform data model to AzDO project reference: %+v", err)
	}

	err = updateProject(clients, project, projectCreateTimeoutDuration)
	if err != nil {
		return fmt.Errorf("Error updating project: %v", err)
	}
	return resourceProjectRead(d, m)
}

func updateProject(clients *config.AggregatedClient, project *core.TeamProject, timeoutSeconds time.Duration) error {
	operationRef, err := clients.CoreClient.UpdateProject(
		clients.Ctx,
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
	clients := m.(*config.AggregatedClient)
	id := d.Id()

	err := deleteProject(clients, id, projectDeleteTimeoutDuration)
	if err != nil {
		return fmt.Errorf("Error deleting project: %v", err)
	}

	return nil
}

func deleteProject(clients *config.AggregatedClient, id string, timeoutSeconds time.Duration) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("Invalid project UUID: %s", id)
	}

	operationRef, err := clients.CoreClient.QueueDeleteProject(clients.Ctx, core.QueueDeleteProjectArgs{
		ProjectId: &uuid,
	})

	if err != nil {
		return err
	}

	return waitForAsyncOperationSuccess(clients, operationRef, timeoutSeconds)
}

// Convert internal Terraform data structure to an AzDO data structure
func expandProject(clients *config.AggregatedClient, d *schema.ResourceData, forCreate bool) (*core.TeamProject, error) {
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

	visibility := core.ProjectVisibility(d.Get("visibility").(string))

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
		Visibility:   &visibility,
		Capabilities: capabilities,
	}

	return project, nil
}

func flattenProject(clients *config.AggregatedClient, d *schema.ResourceData, project *core.TeamProject) error {
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
func lookupProcessTemplateID(clients *config.AggregatedClient, templateName string) (string, error) {
	processes, err := clients.CoreClient.GetProcesses(clients.Ctx, core.GetProcessesArgs{})
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
func lookupProcessTemplateName(clients *config.AggregatedClient, templateID string) (string, error) {
	id, err := uuid.Parse(templateID)
	if err != nil {
		return "", fmt.Errorf("Error parsing Work Item Template ID, got %s: %v", templateID, err)
	}

	process, err := clients.CoreClient.GetProcessById(clients.Ctx, core.GetProcessByIdArgs{
		ProcessId: &id,
	})

	if err != nil {
		return "", fmt.Errorf("Error looking up template by ID: %v", err)
	}

	return *process.Name, nil
}

// ParseImportedProjectIDAndID : Parse the Id (projectId/int) or (projectName/int)
func ParseImportedProjectIDAndID(clients *config.AggregatedClient, id string) (string, int, error) {
	project, resourceID, err := tfhelper.ParseImportedID(id)
	if err != nil {
		return "", 0, err
	}

	// Get the project ID
	currentProject, err := ProjectRead(clients, project, project)
	if err != nil {
		return "", 0, err
	}

	return currentProject.Id.String(), resourceID, nil
}

// ParseImportedProjectIDAndUUID : Parse the Id (projectId/uuid) or (projectName/uuid)
func ParseImportedProjectIDAndUUID(clients *config.AggregatedClient, id string) (string, string, error) {
	project, resourceID, err := tfhelper.ParseImportedUUID(id)
	if err != nil {
		return "", "", err
	}

	// Get the project ID
	currentProject, err := ProjectRead(clients, project, project)
	if err != nil {
		return "", "", err
	}

	return currentProject.Id.String(), resourceID, nil
}
