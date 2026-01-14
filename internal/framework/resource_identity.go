package framework

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

type ResourceIdentity interface {
	// Convert from an import id to the identity
	FromId(id string)
	// Fields returns each identity field value, together with its path and its corresponding state path.
	Fields() []IdentityField
}

type IdentityField struct {
	PathIdentity path.Path
	PathState    path.Path
	Value        attr.Value
}
