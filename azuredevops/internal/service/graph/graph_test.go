package graph

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// The unit tests in this package drive a mocked graph client, which reports the
	// membership in sync on the first check, so both intervals below only add wall
	// clock time to the test run.
	membershipSyncDelay = 0
	membershipSyncMinTimeout = 0
	os.Exit(m.Run())
}
