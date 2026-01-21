package security

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/security"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// TokenTemplate defines the structure for generating tokens for a namespace
type TokenTemplate struct {
	// RequiredIdentifiers are the identifiers that must be provided
	RequiredIdentifiers []string
	// OptionalIdentifiers are the identifiers that may be provided
	OptionalIdentifiers []string
	// BuildFunc is the function to build the token (receives identifiers and API clients)
	BuildFunc func(identifiers map[string]string, clients *client.AggregatedClient) (string, error)
}

// TokenTemplate defines the structure for generating tokens for a namespace
// This matches the logic from CreateClassificationNodeSecurityToken in the permissions utils
func createClassificationNodeToken(clients *client.AggregatedClient, projectID string, path string, structureGroup workitemtracking.TreeStructureGroup) (string, error) {
	const aclClassificationNodeTokenPrefix = "vstfs:///Classification/Node/"

	// Get the root ClassificationNode
	rootClassificationNode, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &structureGroup,
		Depth:          converter.Int(1),
	})
	if err != nil {
		return "", fmt.Errorf("error getting root classification node: %w", err)
	}

	/*
	 * Token format
	 * Root Node: vstfs:///Classification/Node/<NodeIdentifier>"
	 * 1st child: vstfs:///Classification/Node/<NodeIdentifier>:vstfs:///Classification/Node/<NodeIdentifier>
	 */
	aclToken := aclClassificationNodeTokenPrefix + rootClassificationNode.Identifier.String()

	if path != "" {
		path = strings.TrimLeft(strings.TrimSpace(path), "/")
		if path != "" && (rootClassificationNode.HasChildren == nil || !*rootClassificationNode.HasChildren) {
			return "", fmt.Errorf("a path was specified but the root classification node has no children")
		} else if path != "" {
			// Get the id for each classification in the provided path
			// We do this by querying each path element
			// 0: foo
			// 1: foo/bar
			// 2: foo/bar/baz
			pathSegments := strings.Split(path, "/")
			var pathElem []string

			// Filter out empty segments
			for _, elem := range pathSegments {
				if len(elem) > 0 {
					pathElem = append(pathElem, elem)
				}
			}

			for i := range pathElem {
				pathItem := strings.Join(pathElem[:i+1], "/")
				node, err := clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
					Project:        &projectID,
					Path:           &pathItem,
					StructureGroup: &structureGroup,
					Depth:          converter.Int(1),
				})
				if err != nil {
					return "", fmt.Errorf("error getting classification node: %w", err)
				}

				aclToken = aclToken + ":" + aclClassificationNodeTokenPrefix + node.Identifier.String()
			}
		}
	}

	return aclToken, nil
}

// getQueryIDsFromPath resolves a path string to a list of query/folder IDs
// This matches the logic from getQueryIDsFromPath in the workitemquery permissions resource
func getQueryIDsFromPath(clients *client.AggregatedClient, projectID string, path string) ([]string, error) {
	path = strings.TrimSpace(path)

	// Parse path segments, filtering out empty strings
	var pathItems []string
	for _, segment := range strings.Split(path, "/") {
		if len(segment) > 0 {
			pathItems = append(pathItems, segment)
		}
	}

	// Start with "Shared Queries" folder
	qry, err := clients.WorkItemTrackingClient.GetQuery(clients.Ctx, workitemtracking.GetQueryArgs{
		Project: &projectID,
		Query:   converter.String("Shared Queries"),
		Depth:   converter.Int(1),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting Shared Queries folder: %w", err)
	}

	ret := []string{qry.Id.String()}

	// Traverse path segments
	for _, segment := range pathItems {
		if qry.Children == nil || len(*qry.Children) == 0 {
			return nil, fmt.Errorf("unable to find query [%s] in folder [%s] because it has no children", segment, getQueryName(qry))
		}

		// Try to parse segment as UUID first, otherwise match by name
		segmentUUID, parseErr := uuid.Parse(segment)
		childIdx := -1

		for idx, child := range *qry.Children {
			if parseErr == nil && strings.EqualFold(segmentUUID.String(), child.Id.String()) {
				childIdx = idx
				break
			} else if child.Name != nil && strings.EqualFold(*child.Name, segment) {
				childIdx = idx
				break
			}
		}

		if childIdx < 0 {
			return nil, fmt.Errorf("unable to find query [%s] in folder [%s]", segment, getQueryName(qry))
		}

		// Get the child query/folder with depth 1
		qry, err = clients.WorkItemTrackingClient.GetQuery(clients.Ctx, workitemtracking.GetQueryArgs{
			Project: &projectID,
			Query:   converter.String((*qry.Children)[childIdx].Id.String()),
			Depth:   converter.Int(1),
		})
		if err != nil {
			return nil, fmt.Errorf("error getting query: %w", err)
		}

		ret = append(ret, qry.Id.String())
	}

	return ret, nil
}

// getQueryName returns the name of a query, falling back to its ID if name is not available
func getQueryName(qry *workitemtracking.QueryHierarchyItem) string {
	if qry.Name != nil {
		return *qry.Name
	}
	return qry.Id.String()
}

// TokenTemplate defines the structure for generating tokens for a namespace
var namespaceTokenTemplates = map[utils.SecurityNamespaceID]TokenTemplate{
	// Git Repositories namespace
	// Token formats:
	// repoV2/project_id
	// repoV2/project_id/repository_id
	// repoV2/project_id/repository_id/ref_name (with ref_name segments after (refs/heads|tags) encoded in UTF-16LE hex)
	utils.SecurityNamespaceIDValues.GitRepositories: {
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{"repository_id", "ref_name"},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			if _, hasRef := identifiers["ref_name"]; hasRef {
				if _, hasRepo := identifiers["repository_id"]; !hasRepo {
					return "", fmt.Errorf("ref_name provided without repository_id; ref_name needs it's parent repository_id to build the token")
				}
			}
			if repoID, hasRepo := identifiers["repository_id"]; hasRepo {
				if refName, hasRef := identifiers["ref_name"]; hasRef {
					// split refName. The first 2 segments should stay as is, the rest need to be stringToUTF16LEHex
					segments := strings.Split(refName, "/")
					// skip the first two segments
					if len(segments) > 2 {
						for i := 2; i < len(segments); i++ {
							encoded, err := converter.EncodeUtf16HexString(segments[i])
							if err != nil {
								return "", fmt.Errorf("failed to encode segment '%s': %w", segments[i], err)
							}
							segments[i] = encoded
						}
					}
					// join back the segments
					refName = strings.Join(segments, "/")
					return fmt.Sprintf("repoV2/%s/%s/%s/", identifiers["project_id"], repoID, refName), nil
				}
				return fmt.Sprintf("repoV2/%s/%s", identifiers["project_id"], repoID), nil
			}
			return fmt.Sprintf("repoV2/%s", identifiers["project_id"]), nil
		},
	},
	// Project namespace
	// Token format: vstfs:///Classification/TeamProject/project_id
	utils.SecurityNamespaceIDValues.Project: {
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			return fmt.Sprintf("$PROJECT:vstfs:///Classification/TeamProject/%s", identifiers["project_id"]), nil
		},
	},
	// Build namespace
	// Token formats:
	// project_id
	// project_id/build_definition_id
	// project_id/path
	// project_id/path/build_definition_id
	// Note: In practice, getting the full path requires API calls, so path must be pre-transformed
	utils.SecurityNamespaceIDValues.Build: {
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{"path", "build_definition_id"},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			projectID := identifiers["project_id"]
			buildDefID, hasBuildDef := identifiers["build_definition_id"]
			path, hasPath := identifiers["path"]

			if hasBuildDef {
				if hasPath && path != "" && path != "\\" {
					// Remove leading/trailing slashes and convert backslashes to forward slashes
					transformedPath := strings.Trim(strings.ReplaceAll(path, "\\", "/"), "/")
					return fmt.Sprintf("%s/%s/%s", projectID, transformedPath, buildDefID), nil
				}
				return fmt.Sprintf("%s/%s", projectID, buildDefID), nil
			}
			if hasPath {
				// Remove leading/trailing slashes and convert backslashes to forward slashes
				transformedPath := strings.Trim(strings.ReplaceAll(path, "\\", "/"), "/")
				return fmt.Sprintf("%s/%s", projectID, transformedPath), nil
			}
			return projectID, nil
		},
	},
	// Service Endpoints namespace
	utils.SecurityNamespaceIDValues.ServiceEndpoints: {
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			return "", fmt.Errorf("service Endpoints namespace uses role assignments for permissions; token generation is not supported, Role assignment scope is distributedtask.project.serviceendpointrole")
		},
	},
	// CSS namespace (Areas)
	// Note: Requires API calls to resolve path to node identifiers
	utils.SecurityNamespaceIDValues.CSS: {
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{"path"},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			projectID := identifiers["project_id"]
			path := identifiers["path"]
			return createClassificationNodeToken(clients, projectID, path, workitemtracking.TreeStructureGroupValues.Areas)
		},
	},
	// Iteration namespace
	// Note: Requires API calls to resolve path to node identifiers
	utils.SecurityNamespaceIDValues.Iteration: {
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{"path"},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			projectID := identifiers["project_id"]
			path := identifiers["path"]
			return createClassificationNodeToken(clients, projectID, path, workitemtracking.TreeStructureGroupValues.Iterations)
		},
	},
	// Tagging namespace
	utils.SecurityNamespaceIDValues.Tagging: {
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{"project_id"},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			if projectID, hasProject := identifiers["project_id"]; hasProject {
				return fmt.Sprintf("/%s", projectID), nil
			}
			return "", nil
		},
	},
	// Service Hooks namespace
	utils.SecurityNamespaceIDValues.ServiceHooks: {
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{"project_id"},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			if projectID, hasProject := identifiers["project_id"]; hasProject {
				return fmt.Sprintf("PublisherSecurity/%s", projectID), nil
			}
			return "PublisherSecurity", nil
		},
	},
	// Work Item Query Folders namespace
	// Note: Requires API calls to resolve path to query IDs
	utils.SecurityNamespaceIDValues.WorkItemQueryFolders: {
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{"path"},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			projectID := identifiers["project_id"]
			aclToken := fmt.Sprintf("$/%s", projectID)

			if path, hasPath := identifiers["path"]; hasPath && path != "" {
				idList, err := getQueryIDsFromPath(clients, projectID, path)
				if err != nil {
					return "", err
				}
				aclToken = fmt.Sprintf("%s/%s", aclToken, strings.Join(idList, "/"))
			}

			return aclToken, nil
		},
	},
	// Analytics namespace
	utils.SecurityNamespaceIDValues.Analytics: {
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			return fmt.Sprintf("$/%s", identifiers["project_id"]), nil
		},
	},
	// AnalyticsViews namespace
	utils.SecurityNamespaceIDValues.AnalyticsViews: {
		RequiredIdentifiers: []string{"project_id"},
		OptionalIdentifiers: []string{},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			return fmt.Sprintf("$/Shared/%s", identifiers["project_id"]), nil
		},
	},
	// Collection namespace
	utils.SecurityNamespaceIDValues.Collection: {
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			return "NAMESPACE:", nil
		},
	},
	// Process namespace
	utils.SecurityNamespaceIDValues.Process: {
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{"workitem_template_id", "process_id"},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			// if process_id is given but not workitem_template_id, throw an error
			if _, hasProcess := identifiers["process_id"]; hasProcess {
				if _, hasTemplate := identifiers["workitem_template_id"]; !hasTemplate {
					return "", fmt.Errorf("process_id provided without workitem_template_id; process_id needs it's parent workitem_template_id to build the token")
				}
			}
			if templateID, hasTemplate := identifiers["workitem_template_id"]; hasTemplate {
				if processID, hasProcess := identifiers["process_id"]; hasProcess {
					return fmt.Sprintf("$PROCESS:%s:%s:", processID, templateID), nil
				}
				return fmt.Sprintf("$PROCESS:%s:", templateID), nil
			}
			return "$PROCESS:", nil
		},
	},
	// AuditLog namespace
	utils.SecurityNamespaceIDValues.AuditLog: {
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			return "AllPermissions", nil
		},
	},
	// BuildAdministration namespace
	utils.SecurityNamespaceIDValues.BuildAdministration: {
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			return "BuildPrivileges", nil
		},
	},
	// Server namespace
	utils.SecurityNamespaceIDValues.Server: {
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			return "FrameworkGlobalSecurity", nil
		},
	},
	// VersionControlPrivileges namespace
	utils.SecurityNamespaceIDValues.VersionControlPrivileges: {
		RequiredIdentifiers: []string{},
		OptionalIdentifiers: []string{},
		BuildFunc: func(identifiers map[string]string, clients *client.AggregatedClient) (string, error) {
			return "Global", nil
		},
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
		template, exists := namespaceTokenTemplates[utils.SecurityNamespaceID(namespaceID)]
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
	token, err := generateToken(d, namespaceID, clients)
	if err != nil {
		return fmt.Errorf("generating token: %v", err)
	}

	d.Set("token", token)
	d.SetId(fmt.Sprintf("ns-token-%s-%s", namespaceID.String(), token))

	return nil
}

func generateToken(d *schema.ResourceData, namespaceID uuid.UUID, clients *client.AggregatedClient) (string, error) {
	identifiers := make(map[string]string)

	// Get identifiers from the schema
	if ids, ok := d.GetOk("identifiers"); ok {
		for k, v := range ids.(map[string]interface{}) {
			identifiers[k] = v.(string)
		}
	}

	// Look up the template for this namespace
	template, exists := namespaceTokenTemplates[utils.SecurityNamespaceID(namespaceID)]
	if !exists {
		// For unknown namespaces, throw a not supported error
		return "", fmt.Errorf("unable to generate token for namespace %s", namespaceID.String())
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

	// Use the BuildFunc to generate the token
	if template.BuildFunc != nil {
		return template.BuildFunc(identifiers, clients)
	}

	// This should never happen since all templates now have BuildFunc
	return "", fmt.Errorf("no token generation function available for namespace %s", namespaceID.String())
}
