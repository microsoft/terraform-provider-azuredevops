package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProviderResourcesMap(t *testing.T) {
	resources := Provider().ResourcesMap

	require.Equal(t, 3, len(resources), "Three resources were expected to be defined by the provider.")

	require.Contains(t, resources, "azuredevops_pipeline", "Expected resource schema was not found in the provider.")
	require.NotNil(t, resources["azuredevops_pipeline"], "Resource schema cannot be nil.")

	require.Contains(t, resources, "azuredevops_project", "Expected resource schema was not found in the provider.")
	require.NotNil(t, resources["azuredevops_project"], "Resource schema cannot be nil.")
}
