package sdk

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var _ diag.Diagnostic = diagResourceNotFound{}

// diagResourceNotFound implies the resource is not found.
type diagResourceNotFound struct {
	resourceType string
	identity     string
}

func NewDiagResourceNotFound(resourceType, identity string) diagResourceNotFound {
	return diagResourceNotFound{resourceType: resourceType, identity: identity}
}

func IsDiagResourceNotFound(o diag.Diagnostic) bool {
	_, ok := o.(diagResourceNotFound)
	return ok
}

// Summary implements diag.Diagnostic.
func (d diagResourceNotFound) Summary() string {
	return "Resource not found at the service side."
}

// Detail implements diag.Diagnostic.
func (d diagResourceNotFound) Detail() string {
	return fmt.Sprintf("resource_type=%s, identity=%s", d.resourceType, d.identity)
}

// Equal implements diag.Diagnostic.
func (d diagResourceNotFound) Equal(o diag.Diagnostic) bool {
	if do, ok := o.(diagResourceNotFound); ok {
		return d == do
	}
	return false
}

// Severity implements diag.Diagnostic.
func (d diagResourceNotFound) Severity() diag.Severity {
	return diag.SeverityError
}
