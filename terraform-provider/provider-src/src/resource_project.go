package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
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
			"organization": &schema.Schema{ //default func to pull from env?
				Type:     schema.TypeString,
				Required: true,
				//DiffSupressFunc: suppressCaseSensitivity, TODO
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				//DiffSupressFunc: suppressCaseSensitivity, TODO
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"visibility": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "private",
				ValidateFunc: validateVisibility,
			},
			"version_control": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Git",
				ValidateFunc: validateVersionControl,
			},
			"work_item_template": &schema.Schema{
				Type:     schema.TypeString, //need validation function  (TODO investigate if you have a custom template but don't have access, is this even a thing)
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
	ctx := context.Background()
	clients, err := getClients(ctx)

	if err != nil {
		return fmt.Errorf("Error creating client: %+v", err)
	}
	//extract project values from TF file
	//organization := d.Get("organization").(string)

	values := projectValues{
		projectName:      d.Get("project_name").(string),
		description:      d.Get("description").(string),
		visibility:       d.Get("visibility").(string),
		versionControl:   d.Get("version_control").(string),
		workItemTemplate: d.Get("work_item_template").(string),
	}
	// lookup process template id
	processTemplateID, err := lookupProcessTemplateID(ctx, clients, values.workItemTemplate)

	if err != nil {
		return fmt.Errorf("Invalid work item template name: %+v", err)
	}

	values.workItemTemplateID = processTemplateID //clean up

	// create project
	projectCreate(ctx, clients, &values)

	//lookup project id
	projectID, err := lookupProjectID(ctx, clients, values.projectName)

	if err != nil {
<<<<<<< HEAD
		return fmt.Errorf("Error looking up project ID for project: %+v", values.projectName)
=======
		return fmt.Errorf("Error looking up project ID for project: %v, %+v", values.projectName, err)
>>>>>>> d9a514c1a9553a4cb1feccc79cd9709261cd6f6f
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
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceProjectRead(d, m)
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func projectCreate(ctx context.Context, clients *AggregatedClient, values *projectValues) (map[string]string, error) {

	operationRef, err := clients.CoreClient.QueueCreateProject(ctx, core.QueueCreateProjectArgs{
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
		return nil, err
	}

	err = waitForAsyncOperationSuccess(ctx, clients, operationRef)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func lookupProjectID(ctx context.Context, clients *AggregatedClient, projectName string) (string, error) {
	projects, err := clients.CoreClient.GetProjects(ctx, core.GetProjectsArgs{})
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

func lookupProcessTemplateID(ctx context.Context, clients *AggregatedClient, processTemplateName string) (string, error) {
	processes, err := clients.CoreClient.GetProcesses(ctx, core.GetProcessesArgs{})
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

func waitForAsyncOperationSuccess(ctx context.Context, clients *AggregatedClient, operationRef *operations.OperationReference) error {
	maxAttempts := 30
	currentAttempt := 1

	for currentAttempt <= maxAttempts {
		//log.Printf("Checking status for operation with ID: %s", operationRef.Id)
		result, err := clients.OperationsClient.GetOperation(ctx, operations.GetOperationArgs{
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

//diff suppress function to compare strings as case-insensitive
func suppressCaseSensitivity(k, old, new string, d *schema.ResourceData) bool {
	if strings.ToLower(old) == strings.ToLower(new) {
		return true
	}
	return false
}

func validateVisibility(val interface{}, key string) (warns []string, errs []error) {

	valInsensitive := strings.ToLower(val.(string))

	switch valInsensitive {
	case "public":
		if val != "public" {
			errs = append(errs, fmt.Errorf("visibility must be lower case. got: %v", val))
		}
	case "private":
		if val != "private" {
			errs = append(errs, fmt.Errorf("visibility must be lower case. got: %v", val))
		}
	default:
		errs = append(errs, fmt.Errorf("Invalid visiblity value.  Valid values are public/private, got: %v", val))
	}
	return
}

func validateVersionControl(val interface{}, key string) (warns []string, errs []error) {

	valInsensitive := strings.ToLower(val.(string))

	switch valInsensitive {
	case "git":
		return
	case "Tfvc":
		return
	default:
		errs = append(errs, fmt.Errorf("Invalid version control value.  Valid values are Git/Tfvc , got: %v", val))
	}
	return
}
