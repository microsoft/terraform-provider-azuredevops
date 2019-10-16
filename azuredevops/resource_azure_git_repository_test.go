package azuredevops

import (
	"testing"
)

/**
 * Begin unit tests
 */

// verifies that the create operation is considered failed if the initial API
// call fails.
func TestAzureGitRepo_Create_DoesNotSwallowErrorFromFailedCreateCall(t *testing.T) {
}

// verifies that a round-trip flatten/expand sequence will not result in data loss
func TestAzureGitRepo_FlattenExpand_RoundTrip(t *testing.T) {
}

// verifies that the resource ID is used for reads if the ID is set
func TestAzureGitRepo_Read_UsesIdIfSet(t *testing.T) {
}

// verifies that the name is used for reads if the ID is not set
func TestAzureGitRepo_Read_UsesNameIfIdNotSet(t *testing.T) {
}

/**
 * Begin acceptance tests
 */

// Verifies that the following sequence of events occurrs without error:
//	(1) TF apply creates resource
//	(2) TF state values are set
//	(3) resource can be queried by ID and has expected name
// 	(4) TF destroy deletes resource
//	(5) resource can no longer be queried by ID
func TestAccAzureGitRepo_Create(t *testing.T) {
}
