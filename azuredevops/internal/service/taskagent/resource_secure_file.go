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
			"properties": {
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

// getSecureFileURL builds the secure file URL for a given HTTP method.
func getSecureFileURL(clients *client.AggregatedClient, projectID, secureFileID string, params url.Values) string {
	base := strings.TrimRight(clients.OrganizationURL, "/") + "/" + strings.TrimLeft(projectID, "/") + "/_apis/distributedtask/securefiles"
	baseParams := url.Values{}
	baseParams.Add("api-version", "6.0-preview.1")
	if params != nil {
		for k, v := range params {
			baseParams[k] = v
		}
	}
	if secureFileID != "" {
		base += "/" + secureFileID
	}
	return base + "?" + baseParams.Encode()
}

// getSecureFileProperties safely extracts the properties map from a secure file response.
func getSecureFileProperties(secureFile map[string]interface{}) map[string]interface{} {
	if remoteProps, ok := secureFile["properties"].(map[string]interface{}); ok {
		return remoteProps
	}
	return map[string]interface{}{}
}

// setSecureFileHashes updates the diff with hash values and forces new content if hashes change.
func setSecureFileHashes(d *schema.ResourceDiff, remoteProps map[string]interface{}) {
	oldSha1 := d.Get("file_hash_sha1").(string)
	newSha1, _ := remoteProps["file_hash_sha1"].(string)
	oldSha256 := d.Get("file_hash_sha256").(string)
	newSha256, _ := remoteProps["file_hash_sha256"].(string)
	d.SetNew("file_hash_sha1", newSha1)
	d.SetNew("file_hash_sha256", newSha256)
	if newSha1 != oldSha1 || newSha256 != oldSha256 {
		d.ForceNew("content")
	}
}

// buildPropertiesMap builds a string map from the resource data, always including hashes if present.
func buildPropertiesMap(d *schema.ResourceData) map[string]string {
	props := map[string]string{}
	if v, ok := d.GetOk("properties"); ok {
		for k, v2 := range v.(map[string]interface{}) {
			props[k] = v2.(string)
		}
	}
	if sha1Property, ok := d.Get("file_hash_sha1").(string); ok && sha1Property != "" {
		props["file_hash_sha1"] = sha1Property
	}
	if sha256Property, ok := d.Get("file_hash_sha256").(string); ok && sha256Property != "" {
		props["file_hash_sha256"] = sha256Property
	}
	return props
}

func calculateContentHashes(content string) (string, string) {
	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(content))
	sha1String := hex.EncodeToString(sha1Hash.Sum(nil))

	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(content))
	sha256String := hex.EncodeToString(sha256Hash.Sum(nil))

	return sha1String, sha256String
}

func resourceSecureFileCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	if d.Id() == "" {
		return nil
	}
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	secureFileID := d.Id()
	finalUrl := getSecureFileURL(clients, projectID, secureFileID, nil)
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
		return fmt.Errorf("error sending secure file read request: %v", err)
	}
	var secureFile map[string]interface{}
	err = clients.RawClient.UnmarshalBody(response, &secureFile)
	if err != nil {
		return err
	}
	remoteProps := getSecureFileProperties(secureFile)
	setSecureFileHashes(d, remoteProps)
	return nil
}

func resourceSecureFileCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	content := d.Get("content").(string)
	sha1String, sha256String := calculateContentHashes(content)

	d.Set("file_hash_sha1", sha1String)
	d.Set("file_hash_sha256", sha256String)
	// Build URL for secure file creation
	createURL := getSecureFileURL(clients, projectID, "", url.Values{"name": []string{name}})

	request, err := clients.RawClient.CreateRequestMessage(
		clients.Ctx,
		http.MethodPost,
		createURL,
		"",
		bytes.NewReader([]byte(content)),
		"application/octet-stream",
		"",
		map[string]string{},
	)
	if err != nil {
		return fmt.Errorf("error creating request message: %v", err)
	}
	response, err := clients.RawClient.SendRequest(request)
	if err != nil {
		return fmt.Errorf("error sending secure file create request: %v", err)
	}
	var secureFile map[string]interface{}
	err = clients.RawClient.UnmarshalBody(response, &secureFile)
	if err != nil {
		return fmt.Errorf("error parsing secure file response: %v", err)
	}
	secureFileID, ok := secureFile["id"].(string)
	if !ok {
		return fmt.Errorf("could not get secure file ID from response")
	}
	d.SetId(secureFileID)

	// PATCH to set properties (including hashes)
	properties := buildPropertiesMap(d)
	patchPayload := map[string]interface{}{
		"id":         secureFileID,
		"name":       name,
		"properties": properties,
	}
	patchBytes, err := json.Marshal(patchPayload)
	if err != nil {
		return fmt.Errorf("error marshaling patch payload: %v", err)
	}
	patchURL := getSecureFileURL(clients, projectID, secureFileID, nil)
	patchReq, err := clients.RawClient.CreateRequestMessage(
		clients.Ctx,
		http.MethodPatch,
		patchURL,
		"",
		bytes.NewReader(patchBytes),
		"application/json",
		"",
		map[string]string{},
	)
	if err != nil {
		return fmt.Errorf("error creating patch request: %v", err)
	}
	_, err = clients.RawClient.SendRequest(patchReq)
	if err != nil {
		return fmt.Errorf("error sending patch request: %v", err)
	}

	return resourceSecureFileRead(d, m)
}

func resourceSecureFileRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	secureFileID := d.Id()
	finalUrl := getSecureFileURL(clients, projectID, secureFileID, nil)
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
		return fmt.Errorf("error sending secure file read request: %v", err)
	}
	if response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}
	var secureFile map[string]interface{}
	err = clients.RawClient.UnmarshalBody(response, &secureFile)
	if err != nil {
		return fmt.Errorf("error parsing secure file response: %v", err)
	}
	if name, ok := secureFile["name"].(string); ok {
		d.Set("name", name)
	}
	props := getSecureFileProperties(secureFile)
	delete(props, "file_hash_sha1")
	delete(props, "file_hash_sha256")
	d.Set("properties", props)
	return nil
}

func resourceSecureFileUpdate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	secureFileID := d.Id()
	name := d.Get("name").(string)
	props := buildPropertiesMap(d)
	patchPayload := map[string]interface{}{
		"id":         secureFileID,
		"name":       name,
		"properties": props,
	}
	payloadBytes, err := json.Marshal(patchPayload)
	if err != nil {
		return fmt.Errorf("error marshaling update payload: %v", err)
	}
	patchURL := getSecureFileURL(clients, projectID, secureFileID, nil)
	request, err := clients.RawClient.CreateRequestMessage(
		clients.Ctx,
		http.MethodPatch,
		patchURL,
		"",
		bytes.NewReader(payloadBytes),
		"application/json",
		"",
		map[string]string{},
	)
	if err != nil {
		return fmt.Errorf("error creating update request: %v", err)
	}
	_, err = clients.RawClient.SendRequest(request)
	if err != nil {
		return fmt.Errorf("error sending secure file update request: %v", err)
	}
	return resourceSecureFileRead(d, m)
}

func resourceSecureFileDelete(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	secureFileID := d.Id()
	finalUrl := getSecureFileURL(clients, projectID, secureFileID, nil)
	request, err := clients.RawClient.CreateRequestMessage(
		clients.Ctx,
		http.MethodDelete,
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
	_, err = clients.RawClient.SendRequest(request)
	if err != nil {
		return fmt.Errorf("error sending secure file delete request: %v", err)
	}
	d.SetId("")
	return nil
}
