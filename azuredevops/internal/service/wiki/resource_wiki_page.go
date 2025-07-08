package wiki

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/wiki"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

/*
To improve concurrent page api response:
"The wiki page has already been updated by another client, so you cannot update it. Please try again."
Add mutex lock to limit terraform provider concurrent create / update / delete request.
*/
var pageLock = sync.Mutex{}

func ResourceWikiPage() *schema.Resource {
	return &schema.Resource{
		Create: resourceWikiPageCreate,
		Read:   resourceWikiPageRead,
		Update: resourceWikiPageUpdate,
		Delete: resourceWikiPageDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"wiki_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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

	pageLock.Lock()
	defer pageLock.Unlock()

	resp, err := clients.WikiClient.CreateOrUpdatePage(clients.Ctx, wiki.CreateOrUpdatePageArgs{
		Parameters: &wiki.WikiPageCreateOrUpdateParameters{
			Content: converter.String(d.Get("content").(string)),
		},
		Project:        converter.String(d.Get("project_id").(string)),
		WikiIdentifier: converter.String(d.Get("wiki_id").(string)),
		Path:           converter.String(d.Get("path").(string)),
	})
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(*resp.Page.Id))
	return resourceWikiPageRead(d, m)
}

func resourceWikiPageRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	resp, err := clients.WikiClient.GetPage(clients.Ctx, wiki.GetPageArgs{
		Project:        converter.String(d.Get("project_id").(string)),
		WikiIdentifier: converter.String(d.Get("wiki_id").(string)),
		Path:           converter.String(d.Get("path").(string)),
		IncludeContent: converter.Bool(true),
	})
	if err != nil {
		return err
	}

	if resp.ETag != nil && len(*resp.ETag) > 0 {
		etagValue := strings.Trim(strings.Join(*resp.ETag, " "), `\"`)
		d.Set("etag", etagValue)
	}

	if resp.Page != nil {
		d.Set("content", resp.Page.Content)
		d.Set("path", resp.Page.Path)
	}
	return nil
}

func resourceWikiPageUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Parse wiki page ID: %s. Error: %+v", d.Id(), err)
	}

	pageLock.Lock()
	defer pageLock.Unlock()

	_, err = clients.WikiClient.UpdatePageById(clients.Ctx, wiki.UpdatePageByIdArgs{
		Parameters: &wiki.WikiPageCreateOrUpdateParameters{
			Content: converter.String(d.Get("content").(string)),
		},
		Id:             &id,
		Project:        converter.String(d.Get("project_id").(string)),
		WikiIdentifier: converter.String(d.Get("wiki_id").(string)),
		Version:        converter.String(d.Get("etag").(string)),
	})
	if err != nil {
		return err
	}

	return resourceWikiPageRead(d, m)
}

func resourceWikiPageDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Parse wiki page ID: %s. Error: %+v", d.Id(), err)
	}

	pageLock.Lock()
	defer pageLock.Unlock()

	_, err = clients.WikiClient.DeletePageById(clients.Ctx, wiki.DeletePageByIdArgs{
		Project:        converter.String(d.Get("project_id").(string)),
		WikiIdentifier: converter.String(d.Get("wiki_id").(string)),
		Id:             &id,
	})
	if err != nil {
		return err
	}

	return nil
}
