package serviceendpoint

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// The unit tests in this package drive a mocked service endpoint client, so the
	// back-off that `deleteServiceEndpoint` applies between two delete attempts only
	// adds wall clock time to the test run.
	deleteRetryInterval = 0
	os.Exit(m.Run())
}
