package sdk

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/microsoft/terraform-provider-azuredevops/internal/meta"
)

// ResourceTimeout specifies the default timeout value for each operation.
type ResourceTimeout struct {
	Create time.Duration
	Read   time.Duration
	Update time.Duration
	Delete time.Duration
}

type Resource interface {
	// Type returns the resource type
	Type() string

	// SetMeta sets the provider meta to the resource.
	// This is implemented by the WithMeta, the implementer struct shall simply embed it.
	SetMeta(meta.Meta)

	// Timeout returns the default timeout for each operation.
	Timeout() ResourceTimeout

	// Resource implements the framework Resource interface.
	// Since the Metadata() is implemented by the wrapper, the implementer struct shall not implement it.
	// Instead, it is supposed to embed the WithMetadata to meet the interface requirement.
	resource.Resource

	// The followings are interfaces that a resource can opt-in by implementing them.
	// resource.ResourceWithConfigValidators
	// resource.ResourceWithModifyPlan
	// resource.ResourceWithImportState
	// resource.ResourceWithMoveState
	// resource.ResourceWithUpgradeState
	// resource.ResourceWithValidateConfig
}
