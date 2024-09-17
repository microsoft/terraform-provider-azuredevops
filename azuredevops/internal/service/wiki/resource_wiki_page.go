package wiki

import (
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/wiki"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

/*
To improve concurrent page api response:
"The wiki page has already been updated by another client, so you cannot update it. Please try again."
Add mutex lock to limit terraform providor concurrent create / update / delete request.
*/
var pageLock = sync.Mutex{}

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
			"etag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
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

	pageLock.Lock()
	defer pageLock.Unlock()

	_, err := clients.WikiClient.CreateOrUpdatePage(clients.Ctx, wiki.CreateOrUpdatePageArgs{
		Parameters:     &wiki_page_upsert_params,
		Project:        &project_id,
		WikiIdentifier: &wiki_id,
		Path:           &path,
	})

	if err != nil {
		return err
	}
	return resourceWikiPageRead(d, m)
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

	if err != nil {
		return err
	}
	etagValue := strings.Trim(strings.Join(*resp.ETag, " "), `\"`)
	d.Set("etag", etagValue)
	d.SetId(strconv.Itoa(*resp.Page.Id))
	return nil
}

func resourceWikiPageUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	project_id := d.Get("project_id").(string)
	wiki_id := d.Get("wiki_id").(string)

	etag := d.Get("etag").(string)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	content := d.Get("content").(string)
	wiki_page_upsert_params := wiki.WikiPageCreateOrUpdateParameters{
		Content: &content,
	}
	pageLock.Lock()
	defer pageLock.Unlock()

	_, err = clients.WikiClient.UpdatePageById(clients.Ctx, wiki.UpdatePageByIdArgs{
		Parameters:     &wiki_page_upsert_params,
		Project:        &project_id,
		WikiIdentifier: &wiki_id,
		Id:             &id,
		Version:        &etag,
	})

	if err != nil {
		return err
	}
	return resourceWikiPageRead(d, m)
}

func resourceWikiPageDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	project_id := d.Get("project_id").(string)
	wiki_id := d.Get("wiki_id").(string)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	pageLock.Lock()
	defer pageLock.Unlock()

	_, err = clients.WikiClient.DeletePageById(clients.Ctx, wiki.DeletePageByIdArgs{
		Project:        &project_id,
		WikiIdentifier: &wiki_id,
		Id:             &id,
	})
	if err != nil {
		return err
	}

	return nil
}
