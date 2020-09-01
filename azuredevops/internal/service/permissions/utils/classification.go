package utils

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

const aclClassificationNodeTokenPrefix = "vstfs:///Classification/Node/"

// CreateClassificationNodeSecurityToken create security namespace token for iterations and areas
func CreateClassificationNodeSecurityToken(context context.Context, workItemTrackingClient workitemtracking.Client, structureGroup workitemtracking.TreeStructureGroup, projectID string, path string) (string, error) {
	var aclToken string

	// you have to omit the path property to get the
	// root ClassificationNode.
	rootClassificationNode, err := workItemTrackingClient.GetClassificationNode(context, workitemtracking.GetClassificationNodeArgs{
		Project:        &projectID,
		StructureGroup: &structureGroup,
		Depth:          converter.Int(1),
	})
	if err != nil {
		return "", fmt.Errorf("Error getting root classification node: %w", err)
	}

	/*
	 * Token format
	 * Root Node: vstfs:///Classification/Node/<NodeIdentifier>"
	 * 1st child: vstfs:///Classification/Node/<NodeIdentifier>:vstfs:///Classification/Node/<NodeIdentifier>
	 */
	aclToken = aclClassificationNodeTokenPrefix + rootClassificationNode.Identifier.String()

	if path != "" {
		path = strings.TrimLeft(strings.TrimSpace(path), "/")
		if path != "" && (rootClassificationNode.HasChildren == nil || !*rootClassificationNode.HasChildren) {
			return "", fmt.Errorf("A path was specified but the root classification node has no children")
		} else if path != "" {
			// get the id for each classification in the provided path
			// we do this by querying each path element
			// 0: foo
			// 1: foo/bar
			// 3: foo/bar/baz
			var pathElem []string

			linq.From(strings.Split(path, "/")).
				Where(func(elem interface{}) bool {
					return len(elem.(string)) > 0
				}).
				ToSlice(&pathElem)

			for i := range pathElem {
				pathItem := strings.Join(pathElem[:i+1], "/")
				node, err := workItemTrackingClient.GetClassificationNode(context, workitemtracking.GetClassificationNodeArgs{
					Project:        &projectID,
					Path:           &pathItem,
					StructureGroup: &structureGroup,
					Depth:          converter.Int(1),
				})
				if err != nil {
					return "", fmt.Errorf("Error getting classification node: %w", err)
				}

				aclToken = aclToken + ":" + aclClassificationNodeTokenPrefix + node.Identifier.String()
			}
		}
	}

	log.Printf("[DEBUG] CreateClassificationNodeSecurityToken(): Discovered aclToken %q", aclToken)
	return aclToken, nil
}
