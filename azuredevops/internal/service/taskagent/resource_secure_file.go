package taskagent

import (
	"bytes"
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

// ResourceSecureFile schema and implementation for secure file resource
func ResourceSecureFile() *schema.Resource {
	return &schema.Resource{
		Create:        resourceSecureFileCreate,
		Read:          resourceSecureFileRead,
		Update:        resourceSecureFileUpdate,
		Delete:        resourceSecureFileDelete,
		CustomizeDiff: resourceSecureFileCustomizeDiff,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Importer: tfhelper.ImportProjectQualifiedResource(),
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
				Description:  "The ID of the Azure DevOps project.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "The name of the secure file. Must be unique within the project.",
			},
			"content": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				Description:  "The content of the secure file. This is the actual file data that will be stored securely.",
			},
			"properties": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Key-value map of properties. The provider automatically adds file_hash_sha1 and file_hash_sha256.",
			},
			"file_hash_sha1": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "SHA1 hash of the file content. Computed from the content during creation.",
			},
			"file_hash_sha256": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "SHA256 hash of the file content. Computed from the content during creation.",
			},
		},
	}
}

func resourceSecureFileCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	// Skip on new resources
	if d.Id() == "" {
		return nil
	}

	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	secureFileID := d.Id()

	// Build URL for getting secure file
	readURL := fmt.Sprintf("_apis/distributedtask/securefiles/%s", secureFileID)
	queryParams := url.Values{}
	queryParams.Add("api-version", "6.0-preview.1")

	finalUrl := strings.TrimRight(clients.OrganizationURL, "/") + "/" +
		strings.TrimLeft(projectID, "/") + "/" +
		strings.TrimLeft(readURL, "/") + "?" +
		queryParams.Encode()

	// Create request message
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

	// Send the request
	response, err := clients.RawClient.SendRequest(request)
	if err != nil {
		return fmt.Errorf("error sending secure file read request: %v", err)
	}

	// Check response status
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return clients.RawClient.UnwrapError(response)
	}

	// Parse the response body
	var secureFile map[string]interface{}
	err = clients.RawClient.UnmarshalBody(response, &secureFile)
	if err != nil {
		return err
	}
	// Extract properties
	remoteProps := map[string]interface{}{}
	if props, ok := secureFile["properties"].(map[string]interface{}); ok {
		for k, v := range props {
			remoteProps[k] = v
		}
	}

	oldSha1 := d.Get("file_hash_sha1").(string)
	newSha1 := remoteProps["file_hash_sha1"]
	oldSha256 := d.Get("file_hash_sha256").(string)
	newSha256 := remoteProps["file_hash_sha256"]
	if newSha1 != oldSha1 || newSha256 != oldSha256 {
		// File changed remotely, schedule replacement
		d.ForceNew("content")
	}
	return nil
}

func resourceSecureFileCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	content := d.Get("content").(string)
	properties := map[string]string{}
	if v, ok := d.GetOk("properties"); ok {
		for k, v2 := range v.(map[string]interface{}) {
			properties[k] = v2.(string)
		}
	}

	// Calculate hashes for the content
	contentBytes := []byte(content)
	sha1Hash := sha1.New()
	sha1Hash.Write(contentBytes)
	sha1String := hex.EncodeToString(sha1Hash.Sum(nil))

	sha256Hash := sha256.New()
	sha256Hash.Write(contentBytes)
	sha256String := hex.EncodeToString(sha256Hash.Sum(nil))

	// Optionally add hashes to properties
	properties["file_hash_sha1"] = sha1String
	properties["file_hash_sha256"] = sha256String
	if err := d.Set("properties", properties); err != nil {
		return fmt.Errorf("unable to set properties with hash fields: %w", err)
	}
	// Build URL for secure file creation (no properties, only name)
	createURL := projectID + "/_apis/distributedtask/securefiles"
	queryParams := map[string]string{
		"name": name,
	}
	urlValues := url.Values{}
	for key, value := range queryParams {
		urlValues.Add(key, value)
	}
	finalUrl := strings.TrimRight(clients.OrganizationURL, "/") + "/" + strings.TrimLeft(createURL, "/") + "?" + urlValues.Encode()

	// Create request message for secure file creation
	request, err := clients.RawClient.CreateRequestMessage(
		clients.Ctx,
		http.MethodPost,
		finalUrl,
		"6.0-preview.1",
		bytes.NewReader(contentBytes),
		"application/octet-stream",
		"",
		map[string]string{},
	)
	if err != nil {
		return fmt.Errorf("error creating request message: %v", err)
	}

	// Send the request
	response, err := clients.RawClient.SendRequest(request)
	if err != nil {
		return fmt.Errorf("error sending secure file create request: %v", err)
	}

	// Check response status
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return clients.RawClient.UnwrapError(response)
	}

	// Parse the response body
	var secureFile map[string]interface{}
	err = clients.RawClient.UnmarshalBody(response, &secureFile)
	if err != nil {
		return fmt.Errorf("error parsing secure file response: %v", err)
	}

	// Store ID in resource state
	secureFileID, ok := secureFile["id"].(string)
	if !ok {
		return fmt.Errorf("could not get secure file ID from response")
	}
	d.SetId(secureFileID)

	// Store hash values in the state
	d.Set("file_hash_sha1", sha1String)
	d.Set("file_hash_sha256", sha256String)

	// PATCH to set properties
	patchPayload := map[string]interface{}{
		"id":         secureFileID,
		"name":       name,
		"properties": properties,
	}
	patchBytes, err := json.Marshal(patchPayload)
	if err != nil {
		return fmt.Errorf("error marshaling patch payload: %v", err)
	}
	patchURL := strings.TrimRight(clients.OrganizationURL, "/") + "/" + strings.TrimLeft(projectID, "/") + "/_apis/distributedtask/securefiles/" + secureFileID + "?api-version=6.0-preview.1"
	patchReq, err := clients.RawClient.CreateRequestMessage(
		clients.Ctx,
		http.MethodPatch,
		patchURL,
		"6.0-preview.1",
		bytes.NewReader(patchBytes),
		"application/json",
		"",
		map[string]string{},
	)
	if err != nil {
		return fmt.Errorf("error creating patch request: %v", err)
	}
	patchResp, err := clients.RawClient.SendRequest(patchReq)
	if err != nil {
		return fmt.Errorf("error sending patch request: %v", err)
	}
	if patchResp.StatusCode < 200 || patchResp.StatusCode >= 300 {
		return clients.RawClient.UnwrapError(patchResp)
	}

	return resourceSecureFileRead(d, m)
}

func resourceSecureFileRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	secureFileID := d.Id()

	// Build URL for getting secure file
	readURL := fmt.Sprintf("_apis/distributedtask/securefiles/%s", secureFileID)
	queryParams := url.Values{}
	queryParams.Add("api-version", "6.0-preview.1")

	finalUrl := strings.TrimRight(clients.OrganizationURL, "/") + "/" +
		strings.TrimLeft(projectID, "/") + "/" +
		strings.TrimLeft(readURL, "/") + "?" +
		queryParams.Encode()

	// Create request message
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

	// Send the request
	response, err := clients.RawClient.SendRequest(request)
	if err != nil {
		return fmt.Errorf("error sending secure file read request: %v", err)
	}

	// Handle 404 (not found)
	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	// Check response status
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return clients.RawClient.UnwrapError(response)
	}

	// Parse the response body
	var secureFile map[string]interface{}
	err = clients.RawClient.UnmarshalBody(response, &secureFile)
	if err != nil {
		return fmt.Errorf("error parsing secure file response: %v", err)
	}

	// Set values from response
	d.Set("name", secureFile["name"].(string))

	// Extract properties
	props := map[string]interface{}{}
	if remoteProps, ok := secureFile["properties"].(map[string]interface{}); ok {
		for k, v := range remoteProps {
			props[k] = v
		}
	}

	if _, ok := props["file_hash_sha1"].(string); ok {
		// d.Set("file_hash_sha1", sha1)
		delete(props, "file_hash_sha1")
	}
	if _, ok := props["file_hash_sha256"].(string); ok {
		// d.Set("file_hash_sha256", sha256)
		delete(props, "file_hash_sha256")
	}
	d.Set("properties", props)

	return nil
}

func resourceSecureFileUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	secureFileID := d.Id()
	name := d.Get("name").(string)

	// Always preserve hash values
	hashSha1, _ := d.Get("file_hash_sha1").(string)
	hashSha256, _ := d.Get("file_hash_sha256").(string)

	// Build properties map, always including hashes
	props := map[string]string{}
	if v, ok := d.GetOk("properties"); ok {
		for k, v2 := range v.(map[string]interface{}) {
			props[k] = v2.(string)
		}
	}
	// Always set hashes, even if user tries to remove them
	props["file_hash_sha1"] = hashSha1
	props["file_hash_sha256"] = hashSha256
	if err := d.Set("properties", props); err != nil {
		return fmt.Errorf("error updating properties with hash fields: %w", err)
	}
	patchPayload := map[string]interface{}{
		"id":         secureFileID,
		"name":       name,
		"properties": props,
	}

	payloadBytes, err := json.Marshal(patchPayload)
	if err != nil {
		return fmt.Errorf("error marshaling update payload: %v", err)
	}

	patchURL := strings.TrimRight(clients.OrganizationURL, "/") + "/" +
		strings.TrimLeft(projectID, "/") + "/_apis/distributedtask/securefiles/" + secureFileID + "?api-version=6.0-preview.1"

	request, err := clients.RawClient.CreateRequestMessage(
		clients.Ctx,
		http.MethodPatch,
		patchURL,
		"6.0-preview.1",
		bytes.NewReader(payloadBytes),
		"application/json",
		"",
		map[string]string{},
	)
	if err != nil {
		return fmt.Errorf("error creating update request: %v", err)
	}

	response, err := clients.RawClient.SendRequest(request)
	if err != nil {
		return fmt.Errorf("error sending secure file update request: %v", err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return clients.RawClient.UnwrapError(response)
	}

	return resourceSecureFileRead(d, m)
}

func resourceSecureFileDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	secureFileID := d.Id()

	// Build URL for deleting secure file
	deleteURL := fmt.Sprintf("_apis/distributedtask/securefiles/%s", secureFileID)
	queryParams := url.Values{}

	finalUrl := strings.TrimRight(clients.OrganizationURL, "/") + "/" +
		strings.TrimLeft(projectID, "/") + "/" +
		strings.TrimLeft(deleteURL, "/") + "?" +
		queryParams.Encode()

	// Create request message
	request, err := clients.RawClient.CreateRequestMessage(
		clients.Ctx,
		http.MethodDelete,
		finalUrl,
		"6.0-preview.1",
		nil,
		"",
		"application/json",
		map[string]string{},
	)
	if err != nil {
		return fmt.Errorf("error creating request message: %v", err)
	}

	// Send the request
	response, err := clients.RawClient.SendRequest(request)
	if err != nil {
		return fmt.Errorf("error sending secure file delete request: %v", err)
	}

	// Check response status
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return clients.RawClient.UnwrapError(response)
	}

	d.SetId("")
	return nil
}
