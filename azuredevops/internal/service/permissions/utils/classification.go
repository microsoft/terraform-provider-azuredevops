package utils

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/microsoft/azure-devops-go-api/azuredevops/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

const aclClassificationNodeTokenPrefix = "vstfs:///Classification/Node/"

func CreateClassificationNodeSecurityToken(context context.Context, workitemtrackingClient workitemtracking.Client, structureGroup workitemtracking.TreeStructureGroup, projectID string, path string) (string, error) {
	var aclToken string

	// you have to ommit the path property to get the
	// root ClassificationNode.
	rootClassificationNode, err := workitemtrackingClient.GetClassificationNode(context, workitemtracking.GetClassificationNodeArgs{
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
				node, err := workitemtrackingClient.GetClassificationNode(context, workitemtracking.GetClassificationNodeArgs{
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
