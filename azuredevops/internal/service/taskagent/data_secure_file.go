package taskagent

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func DataSecureFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecureFileRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"properties": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"file_hash_sha1": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_hash_sha256": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSecureFileRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)

	// Build URL for listing secure files
	listURL := "_apis/distributedtask/securefiles"
	queryParams := url.Values{}
	queryParams.Add("api-version", "6.0-preview.1")

	finalUrl := strings.TrimRight(clients.OrganizationURL, "/") + "/" +
		strings.TrimLeft(projectID, "/") + "/" +
		strings.TrimLeft(listURL, "/") + "?" +
		queryParams.Encode()

	request, err := clients.RawClient.CreateRequestMessage(
		clients.Ctx,
		http.MethodGet,
		finalUrl,
		"",
		nil,
		"",
		"application/json",
		map[string]string{},
	)
	if err != nil {
		return fmt.Errorf("error creating request message: %v", err)
	}

	response, err := clients.RawClient.SendRequest(request)
	if err != nil {
		return fmt.Errorf("error sending secure file list request: %v", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return clients.RawClient.UnwrapError(response)
	}

	var result struct {
		Value []map[string]interface{} `json:"value"`
	}
	if err := clients.RawClient.UnmarshalBody(response, &result); err != nil {
		return fmt.Errorf("error parsing secure file list response: %v", err)
	}

	var found map[string]interface{}
	for _, sf := range result.Value {
		if sfName, ok := sf["name"].(string); ok && sfName == name {
			found = sf
			break
		}
	}
	if found == nil {
		return fmt.Errorf("secure file with name '%s' not found in project '%s'", name, projectID)
	}

	// Set ID and attributes
	id, ok := found["id"].(string)
	if !ok {
		return fmt.Errorf("could not get secure file ID from response")
	}
	d.SetId(id)
	d.Set("name", name)

	// Set properties if present
	if props, ok := found["properties"].(map[string]interface{}); ok {
		d.Set("properties", props)
		if sha1, ok := props["file_hash_sha1"].(string); ok {
			d.Set("file_hash_sha1", sha1)
		}
		if sha256, ok := props["file_hash_sha256"].(string); ok {
			d.Set("file_hash_sha256", sha256)
		}
	}

	return nil
}
