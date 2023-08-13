package build

import (
	"fmt"
	"log"
	"strings"

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
		Create:   resourceBuildFolderCreate,
		Read:     resourceBuildFolderRead,
		Update:   resourceBuildFolderUpdate,
		Delete:   resourceBuildFolderDelete,
		Importer: tfhelper.ImportProjectQualifiedResource(),
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
		return fmt.Errorf(" failed creating resource Build Folder, %+v", err)
	}

	flattenBuildFolder(d, createdBuildFolder, projectID)
	return resourceBuildFolderRead(d, m)
}

func resourceBuildFolderRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	path := d.Id()

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

	flattenBuildFolder(d, &buildFolder, projectID)
	return nil
}

func resourceBuildFolderUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	oldPath, _ := d.GetChange("path")
	buildFolder, projectID, err := expandBuildFolder(d)
	if err != nil {
		return fmt.Errorf(" failed to expand build folder configurations. Project ID: %s , Error: %+v", projectID, err)
	}

	updatedBuildFolder, err := clients.BuildClient.UpdateFolder(m.(*client.AggregatedClient).Ctx, build.UpdateFolderArgs{
		Project: &projectID,
		Path:    converter.String(oldPath.(string)),
		Folder:  buildFolder,
	})

	if err != nil {
		return fmt.Errorf("failed to update build folder.  Project ID: %s, Error: %+v ", projectID, err)
	}

	flattenBuildFolder(d, updatedBuildFolder, projectID)
	return resourceBuildFolderRead(d, m)
}

func resourceBuildFolderDelete(d *schema.ResourceData, m interface{}) error {
	if strings.EqualFold(d.Id(), "") {
		return nil
	}

	clients := m.(*client.AggregatedClient)

	projectID := d.Get("project_id").(string)
	path := d.Get("path").(string)

	err := clients.BuildClient.DeleteFolder(m.(*client.AggregatedClient).Ctx, build.DeleteFolderArgs{
		Project: &projectID,
		Path:    &path,
	})

	return err
}

func flattenBuildFolder(d *schema.ResourceData, buildFolder *build.Folder, projectID string) {
	d.SetId(*buildFolder.Path)
	d.Set("project_id", projectID)
	d.Set("path", buildFolder.Path)
	d.Set("description", buildFolder.Description)
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

// create a Folder object from the tf Resource Data
func expandBuildFolder(d *schema.ResourceData) (*build.Folder, string, error) {
	projectID := d.Get("project_id").(string)

	projectUuid, err := uuid.Parse(projectID)
	if err != nil {
		return nil, "", err
	}

	projectReference := core.TeamProjectReference{
		Id: &projectUuid,
	}

	buildFolder := build.Folder{
		Description: converter.String(d.Get("description").(string)),
		Path:        converter.String(d.Get("path").(string)),
		Project:     &projectReference,
	}

	return &buildFolder, projectID, nil
}
