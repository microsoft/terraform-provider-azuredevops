package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProviderResourcesMap(t *testing.T) {
	resources := Provider().ResourcesMap

	// Note: these tests are trivial and need to be changed quite a bit when a real implementation
	// exists... These are mostly here so that you write better tests during the hack... :)
	require.Equal(t, 1, len(resources), "Only one resource was expected to be defined by the provider.")
	require.Contains(t, resources, "azuredevops_foo", "Expected resource schema was not found in the provider.")
	require.NotNil(t, resources["azuredevops_foo"], "Resource schema cannot be nil.")
}
