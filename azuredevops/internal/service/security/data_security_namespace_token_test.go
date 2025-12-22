package security

import (
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken_GitRepositoriesWithRepository(t *testing.T) {
	d := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, map[string]interface{}{
		"identifiers": map[string]interface{}{
			"project_id":    "abc-123",
			"repository_id": "def-456",
		},
	})

	namespaceID := uuid.MustParse("2e9eb7ed-3c0a-47d4-87c1-0ffdd275fd87")
	token, err := generateToken(d, namespaceID)

	assert.NoError(t, err)
	assert.Equal(t, "repoV2/abc-123/def-456", token)
}

func TestGenerateToken_GitRepositoriesWithoutRepository(t *testing.T) {
	d := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, map[string]interface{}{
		"identifiers": map[string]interface{}{
			"project_id": "abc-123",
		},
	})

	namespaceID := uuid.MustParse("2e9eb7ed-3c0a-47d4-87c1-0ffdd275fd87")
	token, err := generateToken(d, namespaceID)

	assert.NoError(t, err)
	assert.Equal(t, "repoV2/abc-123", token)
}

func TestGenerateToken_ProjectNamespace(t *testing.T) {
	d := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, map[string]interface{}{
		"identifiers": map[string]interface{}{
			"project_id": "abc-123",
		},
	})

	namespaceID := uuid.MustParse("52d39943-cb85-4d7f-8fa8-c6baac873819")
	token, err := generateToken(d, namespaceID)

	assert.NoError(t, err)
	assert.Equal(t, "$PROJECT:vstfs:///Classification/TeamProject/abc-123", token)
}

func TestGenerateToken_MissingRequiredIdentifier(t *testing.T) {
	d := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, map[string]interface{}{
		"identifiers": map[string]interface{}{},
	})

	namespaceID := uuid.MustParse("52d39943-cb85-4d7f-8fa8-c6baac873819")
	_, err := generateToken(d, namespaceID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing required identifiers: project_id")
}

func TestGenerateToken_CollectionNamespace_NoIdentifiers(t *testing.T) {
	d := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, map[string]interface{}{
		"identifiers": map[string]interface{}{},
	})

	namespaceID := uuid.MustParse("3e65f728-f8bc-4ecd-8764-7e378b19bfa7")
	token, err := generateToken(d, namespaceID)

	assert.NoError(t, err)
	assert.Equal(t, "NAMESPACE:", token)
}

func TestGenerateToken_UnknownNamespace_WithProjectID(t *testing.T) {
	d := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, map[string]interface{}{
		"identifiers": map[string]interface{}{
			"project_id": "abc-123",
		},
	})

	namespaceID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	token, err := generateToken(d, namespaceID)

	assert.NoError(t, err)
	assert.Equal(t, "$/abc-123", token)
}

func TestGenerateToken_UnknownNamespace_NoIdentifiers(t *testing.T) {
	d := schema.TestResourceDataRaw(t, DataSecurityNamespaceToken().Schema, map[string]interface{}{
		"identifiers": map[string]interface{}{},
	})

	namespaceID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	_, err := generateToken(d, namespaceID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to generate token for namespace")
}

func TestGenerateToken_GitRepositoriesIdentifierInfo(t *testing.T) {
	// Test that we can retrieve identifier info for Git Repositories namespace
	namespaceID := uuid.MustParse("2e9eb7ed-3c0a-47d4-87c1-0ffdd275fd87")
	template, exists := namespaceTokenTemplates[namespaceID.String()]

	assert.True(t, exists, "Git Repositories namespace should have a template")
	assert.Equal(t, []string{"project_id"}, template.RequiredIdentifiers)
	assert.Equal(t, []string{"repository_id", "ref_name"}, template.OptionalIdentifiers)
}

func TestGenerateToken_ProjectIdentifierInfo(t *testing.T) {
	// Test that we can retrieve identifier info for Project namespace
	namespaceID := uuid.MustParse("52d39943-cb85-4d7f-8fa8-c6baac873819")
	template, exists := namespaceTokenTemplates[namespaceID.String()]

	assert.True(t, exists, "Project namespace should have a template")
	assert.Equal(t, []string{"project_id"}, template.RequiredIdentifiers)
	assert.Equal(t, []string{}, template.OptionalIdentifiers)
}

func TestGenerateToken_CollectionIdentifierInfo(t *testing.T) {
	// Test that we can retrieve identifier info for Collection namespace
	namespaceID := uuid.MustParse("3e65f728-f8bc-4ecd-8764-7e378b19bfa7")
	template, exists := namespaceTokenTemplates[namespaceID.String()]

	assert.True(t, exists, "Collection namespace should have a template")
	assert.Equal(t, []string{}, template.RequiredIdentifiers)
	assert.Equal(t, []string{}, template.OptionalIdentifiers)
}
