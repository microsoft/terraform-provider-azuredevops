package wiki

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/wiki"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func ResourceWikiPage() *schema.Resource {
	return &schema.Resource{
		Create: resourceWikiPageCreate,
		Read:   resourceWikiPageRead,
		Update: resourceWikiPageUpdate,
		Delete: resourceWikiPageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"wiki_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			/*
				"comment" Comment to be associated with the page operation.
				"version" (Optional in case of ProjectWiki).
			*/
		},
	}
}

func resourceWikiPageCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	project_id := d.Get("project_id").(string)
	wiki_id := d.Get("wiki_id").(string)
	path := d.Get("path").(string)
	content := d.Get("content").(string)
	wiki_page_upsert_params := wiki.WikiPageCreateOrUpdateParameters{
		Content: &content,
	}
	resp, err := clients.WikiClient.CreateOrUpdatePage(clients.Ctx, wiki.CreateOrUpdatePageArgs{
		Parameters:     &wiki_page_upsert_params,
		Project:        &project_id,
		WikiIdentifier: &wiki_id,
		Path:           &path,
	})
	d.Set("e_tag", *resp.ETag)
	d.SetId(strconv.Itoa(*resp.Page.Id))
	if err != nil {
		return err
	}

	return resourceWikiPageRead(d, m) // need this?
}

func resourceWikiPageRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	project_id := d.Get("project_id").(string)
	wiki_id := d.Get("wiki_id").(string)
	path := d.Get("path").(string)

	resp, err := clients.WikiClient.GetPage(clients.Ctx, wiki.GetPageArgs{
		Project:        &project_id,
		WikiIdentifier: &wiki_id,
		Path:           &path,
	})
	d.SetId(strconv.Itoa(*resp.Page.Id))
	d.Set("e_tag", *resp.ETag)
	if err != nil {
		return err
	}

	return nil
}

func resourceWikiPageUpdate(d *schema.ResourceData, m interface{}) error {

	resourceWikiPageDelete(d, m)
	resourceWikiPageCreate(d, m)
	// TODO: try update API

	return resourceWikiPageRead(d, m)
}

func resourceWikiPageDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	project_id := d.Get("project_id").(string)
	wiki_id := d.Get("wiki_id").(string)
	id, _ := strconv.Atoi(d.Id())

	_, err := clients.WikiClient.DeletePageById(clients.Ctx, wiki.DeletePageByIdArgs{
		Project:        &project_id,
		WikiIdentifier: &wiki_id,
		Id:             &id,
		Comment:        converter.String("Delete path"),
	})
	if err != nil {
		return err
	}
	d.SetId("")

	return nil
}
