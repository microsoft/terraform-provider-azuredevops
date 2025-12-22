package security

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"unicode/utf16"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/security"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// TokenTemplate defines the structure for generating tokens for a namespace
type TokenTemplate struct {
	// Template is the format string for the token (if BuildFunc is nil)
	Template string
	// RequiredIdentifiers are the identifiers that must be provided
	RequiredIdentifiers []string
	// OptionalIdentifiers are the identifiers that may be provided
	OptionalIdentifiers []string
	// BuildFunc is an optional custom function to build the token
	BuildFunc func(identifiers map[string]string) string
}

func stringToUTF16LEHex(s string) string {
	// Convert string to UTF-16 code points
	utf16Codes := utf16.Encode([]rune(s))

	// Convert to bytes in little-endian format
	bytes := make([]byte, len(utf16Codes)*2)
	for i, code := range utf16Codes {
		bytes[i*2] = byte(code)        // Low byte
		bytes[i*2+1] = byte(code >> 8) // High byte
	}

	// Convert bytes to hexadecimal string
	return hex.EncodeToString(bytes)
}

// namespaceTokenTemplates maps namespace IDs to their token generation templates
var namespaceTokenTemplates = map[string]TokenTemplate{
	// Git Repositories namespace
	"2e9eb7ed-3c0a-47d4-87c1-0ffdd275fd87": {
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{"repository_id", "ref_name"},
		BuildFunc: func(identifiers map[string]string) string {
			if repoID, hasRepo := identifiers["repository_id"]; hasRepo {
				if refName, hasRef := identifiers["ref_name"]; hasRef {
					// split refName. The first 2 segments should stay as is, the rest need to be stringToUTF16LEHex
					segments := strings.Split(refName, "/")
					// skip the first two segments
					if len(segments) > 2 {
						for i := 2; i < len(segments); i++ {
							segments[i] = stringToUTF16LEHex(segments[i])
						}
					}
					// join back the segments
					refName = strings.Join(segments, "/")
					return fmt.Sprintf("repoV2/%s/%s/%s/", identifiers["project_id"], repoID, refName)
				}
				return fmt.Sprintf("repoV2/%s/%s", identifiers["project_id"], repoID)
			}
			return fmt.Sprintf("repoV2/%s", identifiers["project_id"])
		},
	},
	// Project namespace
	"52d39943-cb85-4d7f-8fa8-c6baac873819": {
		Template:            "$PROJECT:vstfs:///Classification/TeamProject/%s",
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{},
	},
	// Build namespace
	// Token format: project_id or project_id/path/build_definition_id or project_id/build_definition_id
	// Note: In practice, getting the full path requires API calls, so path must be pre-transformed
	"33344d9c-fc72-4d6f-aba5-fa317101a7e9": {
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{"path", "build_definition_id"},
		BuildFunc: func(identifiers map[string]string) string {
			projectID := identifiers["project_id"]
			buildDefID, hasBuildDef := identifiers["build_definition_id"]
			path, hasPath := identifiers["path"]

			if hasBuildDef {
				if hasPath && path != "" && path != "\\" {
					// Remove leading/trailing slashes and convert backslashes to forward slashes
					transformedPath := strings.Trim(strings.ReplaceAll(path, "\\", "/"), "/")
					return fmt.Sprintf("%s/%s/%s", projectID, transformedPath, buildDefID)
				}
				return fmt.Sprintf("%s/%s", projectID, buildDefID)
			}
			return projectID
		},
	},
	// Service Endpoints namespace
	"49b48001-ca20-4adc-8111-5b60c903a50c": {
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{"serviceendpoint_id"},
		BuildFunc: func(identifiers map[string]string) string {
			projectID := identifiers["project_id"]
			if seID, hasSE := identifiers["serviceendpoint_id"]; hasSE {
				return fmt.Sprintf("endpoints/%s/%s", projectID, seID)
			}
			return fmt.Sprintf("endpoints/%s", projectID)
		},
	},
	// CSS namespace (Areas)
	// Note: Requires API calls to resolve path to node identifiers
	"83e28ad4-2d72-4ceb-97b0-c7726d5502c3": {
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{"node_id"},
		BuildFunc: func(identifiers map[string]string) string {
			// Format: vstfs:///Classification/Node/<NodeIdentifier>
			// For nested: vstfs:///Classification/Node/<RootID>:vstfs:///Classification/Node/<ChildID>
			if nodeID, hasNode := identifiers["node_id"]; hasNode {
				return fmt.Sprintf("vstfs:///Classification/Node/%s", nodeID)
			}
			return ""
		},
	},
	// Iteration namespace
	// Note: Requires API calls to resolve path to node identifiers
	"bf7bfa03-b2b7-47db-8113-fa2e002cc5b1": {
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{"node_id"},
		BuildFunc: func(identifiers map[string]string) string {
			// Format: vstfs:///Classification/Node/<NodeIdentifier>
			// For nested: vstfs:///Classification/Node/<RootID>:vstfs:///Classification/Node/<ChildID>
			if nodeID, hasNode := identifiers["node_id"]; hasNode {
				return fmt.Sprintf("vstfs:///Classification/Node/%s", nodeID)
			}
			return ""
		},
	},
	// Tagging namespace
	"bb50f182-8e5e-40b8-bc21-e8752a1e7ae2": {
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{"project_id"},
		BuildFunc: func(identifiers map[string]string) string {
			if projectID, hasProject := identifiers["project_id"]; hasProject {
				return fmt.Sprintf("/%s", projectID)
			}
			return ""
		},
	},
	// Service Hooks namespace
	"cb594ebe-87dd-4fc9-ac2c-6a10a4c92046": {
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{"project_id"},
		BuildFunc: func(identifiers map[string]string) string {
			if projectID, hasProject := identifiers["project_id"]; hasProject {
				return fmt.Sprintf("PublisherSecurity/%s", projectID)
			}
			return "PublisherSecurity"
		},
	},
	// Work Item Query Folders namespace
	// Note: Requires API calls to resolve path to query IDs
	"71356614-aad7-4757-8f2c-0fb3bff6f680": {
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{"query_id"},
		BuildFunc: func(identifiers map[string]string) string {
			projectID := identifiers["project_id"]
			if queryID, hasQuery := identifiers["query_id"]; hasQuery {
				return fmt.Sprintf("$/%s/%s", projectID, queryID)
			}
			return fmt.Sprintf("$/%s", projectID)
		},
	},
	// Analytics namespace
	"58450c49-b02d-465a-ab12-59ae512d6531": {
		Template:            "$/%s",
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{},
	},
	// AnalyticsViews namespace
	"d34d3680-dfe5-4cc6-a949-7d9c68f73cba": {
		Template:            "$/Shared/%s",
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{},
	},
	// Collection namespace
	"3e65f728-f8bc-4ecd-8764-7e378b19bfa7": {
		Template:            "NAMESPACE:",
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{},
	},
	// Process namespace
	"2dab47f9-bd70-49ed-9bd5-8eb051e59c02": {
		Template:            "$PROCESS:",
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{},
	},
	// AuditLog namespace
	"a6cc6381-a1ca-4b36-b3c1-4e65211e82b6": {
		Template:            "AllPermissions",
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{},
	},
	// BuildAdministration namespace
	"302acaca-b667-436d-a946-87133492041c": {
		Template:            "BuildPrivileges",
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{},
	},
	// Server namespace
	"1f4179b3-6bac-4d01-b421-71ea09171400": {
		Template:            "FrameworkGlobalSecurity",
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{},
	},
	// VersionControlPrivileges namespace
	"66312704-deb5-43f9-b51c-ab4ff5e351c3": {
		Template:            "Global",
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{},
	},
}

// DataSecurityNamespaceToken schema and implementation for security namespace token data source
func DataSecurityNamespaceToken() *schema.Resource {
	return &schema.Resource{
		Read: dataSecurityNamespaceTokenRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"namespace_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				ExactlyOneOf: []string{"namespace_id", "namespace_name"},
				Description:  "The ID of the security namespace",
			},
			"namespace_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"namespace_id", "namespace_name"},
				Description:  "The name of the security namespace (e.g., 'Git Repositories', 'Project')",
			},
			"identifiers": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Map of identifiers required for token generation (e.g., project_id, repository_id)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"return_identifier_info": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, returns the required and optional identifiers for the namespace instead of generating a token",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The generated security token for the namespace",
			},
			"required_identifiers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of required identifiers for this namespace (only populated when return_identifier_info is true)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"optional_identifiers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of optional identifiers for this namespace (only populated when return_identifier_info is true)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSecurityNamespaceTokenRead(d *schema.ResourceData, m interface{}) error {
	clients := m.(*client.AggregatedClient)

	var namespaceID uuid.UUID
	var err error

	// Get namespace ID either directly or by name
	if nsID, ok := d.GetOk("namespace_id"); ok {
		namespaceID, err = uuid.Parse(nsID.(string))
		if err != nil {
			return fmt.Errorf("invalid namespace_id: %v", err)
		}
	} else if nsName, ok := d.GetOk("namespace_name"); ok {
		// Query namespaces to find by name
		namespaces, err := clients.SecurityClient.QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{})
		if err != nil {
			return fmt.Errorf("querying security namespaces: %v", err)
		}

		found := false
		for _, ns := range *namespaces {
			if ns.Name != nil && *ns.Name == nsName.(string) {
				if ns.NamespaceId != nil {
					namespaceID = *ns.NamespaceId
					found = true
					break
				}
			}
		}

		if !found {
			return fmt.Errorf("namespace with name '%s' not found", nsName.(string))
		}
	}

	// Check if we should return identifier info instead of generating a token
	returnIdentifierInfo := d.Get("return_identifier_info").(bool)

	if returnIdentifierInfo {
		// Look up the template for this namespace
		template, exists := namespaceTokenTemplates[namespaceID.String()]
		if !exists {
			return fmt.Errorf("no template information available for namespace %s", namespaceID.String())
		}

		// Set the required and optional identifiers
		d.Set("required_identifiers", template.RequiredIdentifiers)
		d.Set("optional_identifiers", template.OptionalIdentifiers)
		d.SetId(fmt.Sprintf("ns-info-%s", namespaceID.String()))

		return nil
	}

	// Generate token based on namespace and provided parameters
	token, err := generateToken(d, namespaceID)
	if err != nil {
		return fmt.Errorf("generating token: %v", err)
	}

	d.Set("token", token)
	d.SetId(fmt.Sprintf("ns-token-%s-%s", namespaceID.String(), token))

	return nil
}

func generateToken(d *schema.ResourceData, namespaceID uuid.UUID) (string, error) {
	identifiers := make(map[string]string)

	// Get identifiers from the schema
	if ids, ok := d.GetOk("identifiers"); ok {
		for k, v := range ids.(map[string]interface{}) {
			identifiers[k] = v.(string)
		}
	}

	// Look up the template for this namespace
	template, exists := namespaceTokenTemplates[namespaceID.String()]
	if !exists {
		// For unknown namespaces, provide a basic fallback
		if projectID, hasProject := identifiers["project_id"]; hasProject {
			return fmt.Sprintf("$/%s", projectID), nil
		}
		return "", fmt.Errorf("unable to generate token for namespace %s with provided identifiers. Please check documentation for required identifiers", namespaceID.String())
	}

	// Validate required identifiers
	var missing []string
	for _, key := range template.RequiredIdentifiers {
		if _, exists := identifiers[key]; !exists {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return "", fmt.Errorf("missing required identifiers: %s", strings.Join(missing, ", "))
	}

	// If there's a custom build function, use it
	if template.BuildFunc != nil {
		return template.BuildFunc(identifiers), nil
	}

	// Build the list of values to substitute into the template
	// Order: required identifiers first, then optional identifiers
	var values []interface{}

	// Add required identifiers in order
	for _, key := range template.RequiredIdentifiers {
		values = append(values, identifiers[key])
	}

	// Add optional identifiers in order (if provided)
	for _, key := range template.OptionalIdentifiers {
		if val, exists := identifiers[key]; exists {
			values = append(values, val)
		}
	}

	// Generate the token using the template
	// If there are no placeholders, return the template as-is
	if len(values) == 0 {
		return template.Template, nil
	}

	return fmt.Sprintf(template.Template, values...), nil
}
