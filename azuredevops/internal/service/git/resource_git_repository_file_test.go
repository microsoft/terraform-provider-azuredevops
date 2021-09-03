package git

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShortBranchName(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{name: "basic", input: "refs/heads/master", expectedOutput: "master"},
		{name: "none", input: "master", expectedOutput: "master"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := shortBranchName(tt.input)
			require.Equal(t, tt.expectedOutput, output)
		})
	}
}

func TestSplitRepoFilePath(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedRepoId   string
		expectedFilePath string
	}{
		{name: "basic", input: "foo/bar", expectedRepoId: "foo", expectedFilePath: "bar"},
		{name: "nested", input: "foo/bar/baz.txt", expectedRepoId: "foo", expectedFilePath: "bar/baz.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoId, filePath := splitRepoFilePath(tt.input)
			require.Equal(t, tt.expectedRepoId, repoId)
			require.Equal(t, tt.expectedFilePath, filePath)
		})
	}
}
