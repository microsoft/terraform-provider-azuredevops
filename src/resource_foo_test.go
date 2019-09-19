package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestResourceFooCRUDOperations(t *testing.T) {
	mockResourceData := &schema.ResourceData{}

	// Note: these tests are trivial and need to be changed quite a bit when a real implementation
	// exists... These are mostly here so that you write better tests during the hack... :)
	require.Nil(t, resourceFooRead(mockResourceData, nil))
	require.Nil(t, resourceFooUpdate(mockResourceData, nil))
	require.Nil(t, resourceFooDelete(mockResourceData, nil))
}

func TestResourceFooSchema(t *testing.T) {
	resource := resourceFoo()

	require.NotNil(t, resource.Create)
	require.NotNil(t, resource.Read)
	require.NotNil(t, resource.Update)
	require.NotNil(t, resource.Delete)

	require.NotNil(t, resource.Schema)
	require.Equal(t, 1, len(resource.Schema))
	require.Contains(t, resource.Schema, "fookey")
	require.Equal(t, schema.TypeString, resource.Schema["fookey"].Type)
	require.True(t, resource.Schema["fookey"].Required)
}
