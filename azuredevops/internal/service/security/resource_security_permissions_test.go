package security

import (
	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityPermissions_BasicNameMapping(t *testing.T) {
	// Setup
	action1Name := "Read"
	action1DisplayName := "View"
	action1Bit := 1

	action2Name := "Write"
	action2DisplayName := "Edit"
	action2Bit := 2

	namespace := security.SecurityNamespaceDescription{
		Actions: &[]security.ActionDefinition{
			{
				Name:        &action1Name,
				DisplayName: &action1DisplayName,
				Bit:         &action1Bit,
			},
			{
				Name:        &action2Name,
				DisplayName: &action2DisplayName,
				Bit:         &action2Bit,
			},
		},
	}

	requestedPermissions := map[string]interface{}{
		"Read":  "allow",
		"Write": "deny",
	}

	// Execute
	actionMap, err := buildActionMap(namespace, requestedPermissions)

	// Verify
	require.NoError(t, err)
	assert.Equal(t, 1, actionMap["Read"])
	assert.Equal(t, 2, actionMap["Write"])
}

func TestSecurityPermissions_DisplayNameMapping(t *testing.T) {
	// Setup
	action1Name := "GenericRead"
	action1DisplayName := "View"
	action1Bit := 1

	action2Name := "GenericWrite"
	action2DisplayName := "Edit"
	action2Bit := 2

	namespace := security.SecurityNamespaceDescription{
		Actions: &[]security.ActionDefinition{
			{
				Name:        &action1Name,
				DisplayName: &action1DisplayName,
				Bit:         &action1Bit,
			},
			{
				Name:        &action2Name,
				DisplayName: &action2DisplayName,
				Bit:         &action2Bit,
			},
		},
	}

	requestedPermissions := map[string]interface{}{
		"View": "allow",
		"Edit": "deny",
	}

	// Execute
	actionMap, err := buildActionMap(namespace, requestedPermissions)

	// Verify
	require.NoError(t, err)
	assert.Equal(t, 1, actionMap["View"])
	assert.Equal(t, 2, actionMap["Edit"])
	assert.Equal(t, 1, actionMap["GenericRead"])
	assert.Equal(t, 2, actionMap["GenericWrite"])
}

func TestSecurityPermissions_MixedNameAndDisplayName(t *testing.T) {
	// Setup
	action1Name := "GenericRead"
	action1DisplayName := "View"
	action1Bit := 1

	action2Name := "GenericWrite"
	action2DisplayName := "Edit"
	action2Bit := 2

	namespace := security.SecurityNamespaceDescription{
		Actions: &[]security.ActionDefinition{
			{
				Name:        &action1Name,
				DisplayName: &action1DisplayName,
				Bit:         &action1Bit,
			},
			{
				Name:        &action2Name,
				DisplayName: &action2DisplayName,
				Bit:         &action2Bit,
			},
		},
	}

	requestedPermissions := map[string]interface{}{
		"View":         "allow",
		"GenericWrite": "deny",
	}

	// Execute
	actionMap, err := buildActionMap(namespace, requestedPermissions)

	// Verify
	require.NoError(t, err)
	assert.Equal(t, 1, actionMap["View"])
	assert.Equal(t, 2, actionMap["GenericWrite"])
}

func TestSecurityPermissions_DuplicatePermissionError(t *testing.T) {
	// Setup
	action1Name := "GenericRead"
	action1DisplayName := "View"
	action1Bit := 1

	namespace := security.SecurityNamespaceDescription{
		Actions: &[]security.ActionDefinition{
			{
				Name:        &action1Name,
				DisplayName: &action1DisplayName,
				Bit:         &action1Bit,
			},
		},
	}

	// Request the same permission using both Name and DisplayName
	requestedPermissions := map[string]interface{}{
		"GenericRead": "allow",
		"View":        "allow",
	}

	// Execute
	actionMap, err := buildActionMap(namespace, requestedPermissions)

	// Verify
	require.Error(t, err)
	assert.Nil(t, actionMap)
	assert.Contains(t, err.Error(), "permission specified multiple times")
	assert.Contains(t, err.Error(), "GenericRead")
	assert.Contains(t, err.Error(), "View")
}

func TestSecurityPermissions_AmbiguousDisplayName(t *testing.T) {
	// Setup - Two actions with the same DisplayName
	action1Name := "Project.Read"
	action1DisplayName := "View"
	action1Bit := 1

	action2Name := "Repo.Read"
	action2DisplayName := "View" // Same DisplayName as action1
	action2Bit := 2

	namespace := security.SecurityNamespaceDescription{
		Actions: &[]security.ActionDefinition{
			{
				Name:        &action1Name,
				DisplayName: &action1DisplayName,
				Bit:         &action1Bit,
			},
			{
				Name:        &action2Name,
				DisplayName: &action2DisplayName,
				Bit:         &action2Bit,
			},
		},
	}

	// Try to use the ambiguous DisplayName
	requestedPermissions := map[string]interface{}{
		"View": "allow",
	}

	// Execute
	actionMap, err := buildActionMap(namespace, requestedPermissions)

	// Verify
	require.Error(t, err)
	assert.Nil(t, actionMap)
	assert.Contains(t, err.Error(), "ambiguous")
	assert.Contains(t, err.Error(), "View")
}

func TestSecurityPermissions_AmbiguousDisplayNameNotRequested(t *testing.T) {
	// Setup - Two actions with the same DisplayName
	action1Name := "Project.Read"
	action1DisplayName := "View"
	action1Bit := 1

	action2Name := "Repo.Read"
	action2DisplayName := "View" // Same DisplayName as action1
	action2Bit := 2

	action3Name := "Write"
	action3DisplayName := "Edit"
	action3Bit := 4

	namespace := security.SecurityNamespaceDescription{
		Actions: &[]security.ActionDefinition{
			{
				Name:        &action1Name,
				DisplayName: &action1DisplayName,
				Bit:         &action1Bit,
			},
			{
				Name:        &action2Name,
				DisplayName: &action2DisplayName,
				Bit:         &action2Bit,
			},
			{
				Name:        &action3Name,
				DisplayName: &action3DisplayName,
				Bit:         &action3Bit,
			},
		},
	}

	// Use specific names and unambiguous DisplayName
	requestedPermissions := map[string]interface{}{
		"Project.Read": "allow",
		"Edit":         "deny",
	}

	// Execute
	actionMap, err := buildActionMap(namespace, requestedPermissions)

	// Verify - Should succeed because ambiguous DisplayName wasn't used
	require.NoError(t, err)
	assert.Equal(t, 1, actionMap["Project.Read"])
	assert.Equal(t, 4, actionMap["Edit"])
}

func TestSecurityPermissions_EmptyDisplayName(t *testing.T) {
	// Setup - Action with empty DisplayName
	action1Name := "Read"
	action1DisplayName := ""
	action1Bit := 1

	namespace := security.SecurityNamespaceDescription{
		Actions: &[]security.ActionDefinition{
			{
				Name:        &action1Name,
				DisplayName: &action1DisplayName,
				Bit:         &action1Bit,
			},
		},
	}

	requestedPermissions := map[string]interface{}{
		"Read": "allow",
	}

	// Execute
	actionMap, err := buildActionMap(namespace, requestedPermissions)

	// Verify
	require.NoError(t, err)
	assert.Equal(t, 1, actionMap["Read"])
	// Empty DisplayName should not be added to map
	_, exists := actionMap[""]
	assert.False(t, exists)
}

func TestSecurityPermissions_NoDisplayName(t *testing.T) {
	// Setup - Action without DisplayName field
	action1Name := "Read"
	action1Bit := 1

	namespace := security.SecurityNamespaceDescription{
		Actions: &[]security.ActionDefinition{
			{
				Name: &action1Name,
				Bit:  &action1Bit,
				// DisplayName is nil
			},
		},
	}

	requestedPermissions := map[string]interface{}{
		"Read": "allow",
	}

	// Execute
	actionMap, err := buildActionMap(namespace, requestedPermissions)

	// Verify
	require.NoError(t, err)
	assert.Equal(t, 1, actionMap["Read"])
}

func TestSecurityPermissions_InvalidPermissionName(t *testing.T) {
	// Setup
	action1Name := "Read"
	action1Bit := 1

	namespace := security.SecurityNamespaceDescription{
		Actions: &[]security.ActionDefinition{
			{
				Name: &action1Name,
				Bit:  &action1Bit,
			},
		},
	}

	requestedPermissions := map[string]interface{}{
		"Read":           "allow",
		"NonExistentPerm": "deny",
	}

	// Execute
	actionMap, err := buildActionMap(namespace, requestedPermissions)

	// Verify - Should not error during buildActionMap
	// The error will be caught later when checking if permission exists
	require.NoError(t, err)
	assert.Equal(t, 1, actionMap["Read"])
	_, exists := actionMap["NonExistentPerm"]
	assert.False(t, exists)
}
