package wiki

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/wiki"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceWiki() *schema.Resource {
	return &schema.Resource{
		Create: resourceWikiCreate,
		Read:   resourceWikiRead,
		Update: resourceWikiUpdate,
		Delete: resourceWikiDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"isdisabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"mappedpath": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"remoteurl": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"repository_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"versions": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceWikiCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	wikiType := wiki.WikiType(d.Get("type").(string))
	uuidProject, err := uuid.Parse(d.Get("project_id").(string))
	if err != nil {
		return err
	}
	wikiArgs := &wiki.WikiCreateParametersV2{Name: converter.String(d.Get("name").(string)), ProjectId: &uuidProject, Type: &wikiType}
	mappedPath, b := d.GetOk("mappedpath")

	if b {
		wikiArgs.MappedPath = converter.String(mappedPath.(string))
	}
	repositoryId, b := d.GetOk("repository_id")
	if b {
		repositoryUuid, err := uuid.Parse(repositoryId.(string))
		if err != nil {
			return err
		}
		wikiArgs.RepositoryId = &repositoryUuid
	}

	version, b := d.GetOk("versions")
	if b {
		wikiArgs.Version = &git.GitVersionDescriptor{Version: converter.String(version.(string))}
	}

	CreateWikiArgs := wiki.CreateWikiArgs{WikiCreateParams: wikiArgs}
	resp, err := clients.WikiClient.CreateWiki(clients.Ctx, CreateWikiArgs)

	if err != nil {
		return err
	}

	d.SetId(resp.Id.String())
	return resourceWikiRead(d, m)
}

func resourceWikiRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	resp, err := clients.WikiClient.GetWiki(clients.Ctx, wiki.GetWikiArgs{WikiIdentifier: converter.String(d.Id())})
	if err != nil {
		return err
	}

	if resp.Id != nil {
		d.SetId(resp.Id.String())
		d.Set("id", *resp.Id)
	}
	if resp.Name != nil {
		d.Set("name", *resp.Name)
	}
	if resp.ProjectId != nil {
		d.Set("project_id", *resp.ProjectId)
	}
	if resp.Type != nil {
		d.Set("type", *resp.Type)
	}
	if resp.MappedPath != nil {
		d.Set("mappedpath", *resp.MappedPath)
	}
	if resp.RemoteUrl != nil {
		d.Set("remoteurl", *resp.RemoteUrl)
	}
	if resp.RepositoryId != nil {
		d.Set("repository_id", *resp.RepositoryId)
	}
	if resp.Url != nil {
		d.Set("url", *resp.Url)
	}
	if resp.Versions != nil {
		d.Set("versions", *resp.Versions)
	}

	return nil
}

func resourceWikiUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	requiresUpdate := false
	var updateParameters wiki.WikiUpdateParameters

	if d.HasChange("name") {
		updateParameters.Name = converter.String(d.Get("name").(string))
		requiresUpdate = true
	}
	if requiresUpdate {
		_, err := clients.WikiClient.UpdateWiki(clients.Ctx, wiki.UpdateWikiArgs{
			WikiIdentifier:   converter.String(d.Id()),
			UpdateParameters: &updateParameters})
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceWikiDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	//  codewiki can be deleted normally, for project wiki the repo needs to be deleted
	wikiType := wiki.WikiType(d.Get("type").(string))
	if wikiType == "codeWiki" {

		_, err := clients.WikiClient.DeleteWiki(clients.Ctx, wiki.DeleteWikiArgs{
			WikiIdentifier: converter.String(d.Id()),
			Project:        converter.String(d.Get("project_id").(string))})
		if err != nil {
			return err
		}

	} else if wikiType == "projectWiki" {
		resp, err := clients.WikiClient.GetWiki(clients.Ctx, wiki.GetWikiArgs{WikiIdentifier: converter.String(d.Id())})
		if err != nil {
			return err
		}

		err = clients.GitReposClient.DeleteRepository(clients.Ctx, git.DeleteRepositoryArgs{
			RepositoryId: resp.RepositoryId,
			Project:      converter.String(d.Get("project_id").(string)),
		})
		if err != nil {
			return err
		}
	}
	d.SetId("")
	return nil
}
