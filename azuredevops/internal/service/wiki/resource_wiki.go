package wiki

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/wiki"
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(wiki.WikiTypeValues.ProjectWiki),
					string(wiki.WikiTypeValues.CodeWiki)},
					false),
			},
			"mapped_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"remote_url": {
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
	if mappedPath, ok := d.GetOk("mapped_path"); ok {
		wikiArgs.MappedPath = converter.String(mappedPath.(string))
	}
	if repositoryId, ok := d.GetOk("repository_id"); ok {
		repositoryUuid, err := uuid.Parse(repositoryId.(string))
		if err != nil {
			return err
		}
		wikiArgs.RepositoryId = &repositoryUuid
	}
	if version, ok := d.GetOk("versions"); ok {
		wikiArgs.Version = &git.GitVersionDescriptor{Version: converter.String(version.(string))}
	}

	resp, err := clients.WikiClient.CreateWiki(clients.Ctx, wiki.CreateWikiArgs{WikiCreateParams: wikiArgs})

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
	}
	if resp.Name != nil {
		d.Set("name", *resp.Name)
	}
	if resp.ProjectId != nil {
		d.Set("project_id", resp.ProjectId.String())
	}
	if resp.Type != nil {
		d.Set("type", *resp.Type)
	}
	if resp.MappedPath != nil {
		d.Set("mapped_path", *resp.MappedPath)
	}
	if resp.RemoteUrl != nil {
		d.Set("remote_url", *resp.RemoteUrl)
	}
	if resp.RepositoryId != nil {
		d.Set("repository_id", resp.RepositoryId.String())
	}
	if resp.Url != nil {
		d.Set("url", *resp.Url)
	}
	if resp.Versions != nil {
		d.Set("versions", (*resp.Versions)[0].Version)
	}

	return nil
}

func resourceWikiUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	if d.HasChange("name") {
		_, err := clients.WikiClient.UpdateWiki(clients.Ctx, wiki.UpdateWikiArgs{
			WikiIdentifier: converter.String(d.Id()),
			UpdateParameters: &wiki.WikiUpdateParameters{
				Name: converter.String(d.Get("name").(string)),
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceWikiDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	//  codewiki can be deleted normally, for project wiki the project needs to be deleted
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
		if resp.Id == nil {
			return errors.New("projectWiki can only be removed when attached project is removed.")
		}
	}
	d.SetId("")
	return nil
}
