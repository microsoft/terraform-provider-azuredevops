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
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/build"
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
			"allow_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
	for k, v := range params {
		baseParams[k] = v
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
func setSecureFileHashes(d *schema.ResourceDiff, remoteProps map[string]interface{}) error {
	oldSha1, oldSha256 := calculateContentHashes(d.Get("content").(string))
	newSha1, _ := remoteProps["file_hash_sha1"].(string)
	newSha256, _ := remoteProps["file_hash_sha256"].(string)
	if err := d.SetNew("file_hash_sha1", newSha1); err != nil {
		return err
	}
	if err := d.SetNew("file_hash_sha256", newSha256); err != nil {
		return err
	}
	if newSha1 != oldSha1 || newSha256 != oldSha256 {
		if err := d.ForceNew("content"); err != nil {
			return err
		}
	}
	return nil
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
	err = setSecureFileHashes(d, remoteProps)
	if err != nil {
		return err
	}
	return nil
}

// --- Allow Access helpers for Secure File ---

// updateSecureFileAllowAccess updates the allow_access property for a secure file
func updateSecureFileAllowAccess(clients *client.AggregatedClient, projectID, secureFileID, name string, allowAccess bool) error {
	resourceRefType := "securefile"
	defResourceRef := build.DefinitionResourceReference{
		Type:       &resourceRefType,
		Authorized: &allowAccess,
		Name:       &name,
		Id:         &secureFileID,
	}
	resources := []build.DefinitionResourceReference{defResourceRef}
	_, err := clients.BuildClient.AuthorizeProjectResources(
		clients.Ctx, build.AuthorizeProjectResourcesArgs{
			Resources: &resources,
			Project:   &projectID,
		},
	)
	return err
}

// readSecureFileAllowAccess reads the allow_access property for a secure file
func readSecureFileAllowAccess(clients *client.AggregatedClient, projectID, secureFileID string) (bool, error) {
	resourceRefType := "securefile"
	projectResources, err := clients.BuildClient.GetProjectResources(
		clients.Ctx,
		build.GetProjectResourcesArgs{
			Project: &projectID,
			Type:    &resourceRefType,
			Id:      &secureFileID,
		},
	)
	if err != nil {
		return false, err
	}
	for _, authResource := range *projectResources {
		if secureFileID == *authResource.Id {
			return *authResource.Authorized, nil
		}
	}
	return false, nil
}

func resourceSecureFileCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	content := d.Get("content").(string)
	sha1String, sha256String := calculateContentHashes(content)

	if err := d.Set("file_hash_sha1", sha1String); err != nil {
		return fmt.Errorf("error setting file_hash_sha1: %v", err)
	}
	if err := d.Set("file_hash_sha256", sha256String); err != nil {
		return fmt.Errorf("error setting file_hash_sha256: %v", err)
	}
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

	// Set allow_access if needed
	if err := updateSecureFileAllowAccess(clients, projectID, secureFileID, name, d.Get("allow_access").(bool)); err != nil {
		return fmt.Errorf("error setting allow_access for secure file: %v", err)
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

	// Read allow_access
	allowAccess, err := readSecureFileAllowAccess(clients, projectID, secureFileID)
	if err == nil {
		d.Set("allow_access", allowAccess)
	}
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

	// Update allow_access if changed
	if d.HasChange("allow_access") {
		if err := updateSecureFileAllowAccess(clients, projectID, secureFileID, name, d.Get("allow_access").(bool)); err != nil {
			return fmt.Errorf("error updating allow_access for secure file: %v", err)
		}
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
