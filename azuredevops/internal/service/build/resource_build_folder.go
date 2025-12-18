package build

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/validate"
)

// ResourceBuildFolder schema and implementation for build folder resource
func ResourceBuildFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceBuildFolderCreate,
		Read:   resourceBuildFolderRead,
		Update: resourceBuildFolderUpdate,
		Delete: resourceBuildFolderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				projectNameOrID, path, err := tfhelper.ParseImportedName(d.Id(), "projectid/resourceName")
				if err != nil {
					return nil, fmt.Errorf("parsing the resource ID from the Terraform resource data: %v", err)
				}

				if projectID, err := tfhelper.GetRealProjectId(projectNameOrID, m); err == nil {
					d.SetId(projectID)
					d.Set("project_id", projectID)
					d.Set("path", path)
					return []*schema.ResourceData{d}, nil
				}
				return nil, err
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.Path,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  ``,
			},
		},
	}
}

func resourceBuildFolderCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	description := d.Get("description").(string)
	path := d.Get("path").(string)

	createdBuildFolder, err := createBuildFolder(clients, path, projectID, description)
	if err != nil {
		return fmt.Errorf("failed creating resource Build Folder, %+v", err)
	}

	d.SetId(createdBuildFolder.Project.Id.String())
	return resourceBuildFolderRead(d, m)
}

func resourceBuildFolderRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	path := d.Get("path").(string)

	buildFolders, err := clients.BuildClient.GetFolders(clients.Ctx, build.GetFoldersArgs{
		Project: &projectID,
		Path:    &path,
	})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	if len(*buildFolders) == 0 {
		d.SetId("")
		log.Printf("[TRACE] plugin.terraform-provider-azuredevops: Folder [%s] not found. Removing from state.", path)
		return nil
	}

	buildFolder := (*buildFolders)[0]

	d.Set("project_id", projectID)

	if buildFolder.Path != nil {
		d.Set("path", buildFolder.Path)
	}

	if buildFolder.Description != nil {
		d.Set("description", buildFolder.Description)
	}
	return nil
}

func resourceBuildFolderUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	oldPath, path := d.GetChange("path")
	projectID := d.Get("project_id").(string)
	projectUuid, err := uuid.Parse(projectID)
	if err != nil {
		return fmt.Errorf("failed to parse Project ID. Project ID: %s , Error: %+v", projectID, err)
	}

	_, err = clients.BuildClient.UpdateFolder(m.(*client.AggregatedClient).Ctx, build.UpdateFolderArgs{
		Project: &projectID,
		Path:    converter.String(oldPath.(string)),
		Folder: &build.Folder{
			Description: converter.String(d.Get("description").(string)),
			Path:        converter.String(path.(string)),
			Project: &core.TeamProjectReference{
				Id: &projectUuid,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update build folder.  Project ID: %s, Error: %+v ", projectID, err)
	}

	return resourceBuildFolderRead(d, m)
}

func resourceBuildFolderDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	err := clients.BuildClient.DeleteFolder(m.(*client.AggregatedClient).Ctx, build.DeleteFolderArgs{
		Project: converter.ToPtr(d.Get("project_id").(string)),
		Path:    converter.ToPtr(d.Get("path").(string)),
	})

	return err
}

// create a Folder object to pass to the API
func createBuildFolder(clients *client.AggregatedClient, path string, project string, description string) (*build.Folder, error) {
	projectUuid, err := uuid.Parse(project)
	if err != nil {
		return nil, err
	}

	createdBuild, err := clients.BuildClient.CreateFolder(clients.Ctx, build.CreateFolderArgs{
		Folder: &build.Folder{
			Description: &description,
			Path:        &path,
			Project: &core.TeamProjectReference{
				Id: &projectUuid,
			},
		},
		Project: &project,
		Path:    &path,
	})

	return createdBuild, err
}
