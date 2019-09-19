package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProviderResourcesMap(t *testing.T) {
	resources := Provider().ResourcesMap

	// Note: these tests are trivial and need to be changed quite a bit when a real implementation
	// exists... These are mostly here so that you write better tests during the hack... :)
	require.Equal(t, 2, len(resources), "Two resources were expected to be defined by the provider.")
	require.Contains(t, resources, "azuredevops_foo", "Expected resource schema was not found in the provider.")
	require.Contains(t, resources, "azuredevops_project", "Expected resource schema was not found in the provider.")
	require.NotNil(t, resources["azuredevops_foo"], "Resource schema cannot be nil.")
	require.NotNil(t, resources["azuredevops_project"], "Resource schema cannot be nil.")
}
