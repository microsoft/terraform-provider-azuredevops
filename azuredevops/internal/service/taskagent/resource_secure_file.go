package taskagent

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
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
		Create: resourceSecureFileCreate,
		Read:   resourceSecureFileRead,
		Delete: resourceSecureFileDelete,
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
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "",
			},
			"file_hash_sha1": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"file_hash_sha256": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

// Metadata structure for file hash information
type FileHashMetadata struct {
	SHA1   string `json:"sha1"`
	SHA256 string `json:"sha256"`
}

// Read file content and calculate hashes
func readFileAndCalculateHashes(filePath string) ([]byte, string, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, "", "", fmt.Errorf("error reading file: %v", err)
	}

	// Reset file pointer to beginning to recalculate hashes
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, "", "", fmt.Errorf("error seeking file: %v", err)
	}

	// Calculate SHA1
	sha1Hash := sha1.New()
	if _, err := io.Copy(sha1Hash, file); err != nil {
		return nil, "", "", fmt.Errorf("error calculating SHA1: %v", err)
	}
	sha1String := hex.EncodeToString(sha1Hash.Sum(nil))

	// Reset file pointer to recalculate SHA256
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, "", "", fmt.Errorf("error seeking file: %v", err)
	}

	// Calculate SHA256
	sha256Hash := sha256.New()
	if _, err := io.Copy(sha256Hash, file); err != nil {
		return nil, "", "", fmt.Errorf("error calculating SHA256: %v", err)
	}
	sha256String := hex.EncodeToString(sha256Hash.Sum(nil))

	return fileContent, sha1String, sha256String, nil
}

// Create hash metadata description
func createMetadataDescription(baseDescription string, sha1String, sha256String string) string {
	hashMetadata := FileHashMetadata{
		SHA1:   sha1String,
		SHA256: sha256String,
	}

	metadataJSON, err := json.Marshal(hashMetadata)
	if err != nil {
		// If JSON marshaling fails, fall back to a simpler format
		return fmt.Sprintf("%s [SHA1:%s SHA256:%s]", baseDescription, sha1String, sha256String)
	}

	if baseDescription == "" {
		return string(metadataJSON)
	}
	return fmt.Sprintf("%s %s", baseDescription, string(metadataJSON))
}

// Extract hash metadata from description
func extractHashMetadataFromDescription(description string) (string, string, string) {
	baseDescription := description
	sha1String := ""
	sha256String := ""

	// Try to extract JSON metadata
	jsonStart := strings.LastIndex(description, "{")
	if jsonStart != -1 {
		possibleJSON := description[jsonStart:]
		var hashMetadata FileHashMetadata
		if err := json.Unmarshal([]byte(possibleJSON), &hashMetadata); err == nil {
			baseDescription = strings.TrimSpace(description[:jsonStart])
			sha1String = hashMetadata.SHA1
			sha256String = hashMetadata.SHA256
			return baseDescription, sha1String, sha256String
		}
	}

	// If JSON extraction fails, look for the simpler format
	sha1Start := strings.Index(description, "[SHA1:")
	if sha1Start != -1 {
		sha256End := strings.Index(description, "SHA256:")
		if sha256End != -1 {
			endBracket := strings.Index(description[sha256End:], "]")
			if endBracket != -1 {
				baseDescription = strings.TrimSpace(description[:sha1Start])
				sha1Part := description[sha1Start+6 : sha256End-1]
				sha256Part := description[sha256End+7 : sha256End+endBracket]
				sha1String = strings.TrimSpace(sha1Part)
				sha256String = strings.TrimSpace(sha256Part)
			}
		}
	}

	return baseDescription, sha1String, sha256String
}

func resourceSecureFileCreate(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)
	projectID := d.Get("project_id").(string)
	name := d.Get("name").(string)
	filePath := d.Get("path").(string)
	description := d.Get("description").(string)

	// Read file content and calculate hashes
	fileContent, sha1String, sha256String, err := readFileAndCalculateHashes(filePath)
	if err != nil {
		return fmt.Errorf("error reading file and calculating hashes: %v", err)
	}

	// Create a description that includes hash metadata
	descriptionWithMetadata := createMetadataDescription(description, sha1String, sha256String)

	// Build URL for secure file creation
	createURL := projectID + "/_apis/distributedtask/securefiles"
	queryParams := map[string]string{
		"name": name,
	}
	if descriptionWithMetadata != "" {
		queryParams["description"] = descriptionWithMetadata
	}
	// Convert query params to url.Values
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
		bytes.NewReader(fileContent),
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

	// Extract description and hash metadata
	description := ""
	if descriptionVal, ok := secureFile["description"]; ok && descriptionVal != nil {
		description = descriptionVal.(string)
	}

	baseDescription, sha1String, sha256String := extractHashMetadataFromDescription(description)
	d.Set("description", baseDescription)

	// Only set hash values if they exist in the description
	if sha1String != "" {
		d.Set("file_hash_sha1", sha1String)
	}
	if sha256String != "" {
		d.Set("file_hash_sha256", sha256String)
	}

	return nil
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
