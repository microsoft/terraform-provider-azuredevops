package serviceendpoint

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// The unit tests in this package drive a mocked service endpoint client, which
	// reaches its target state on the first poll, so both intervals below only add
	// wall clock time to the test run.
	stateChangeDelay = 0
	deleteRetryInterval = 0
	os.Exit(m.Run())
}
