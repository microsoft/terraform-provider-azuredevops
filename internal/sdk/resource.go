package sdk

import (
	"context"
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

type ResourceWithTimeout interface {
	resource.Resource

	// Timeout returns the timeout for each operation.
	Timeout() ResourceTimeout
}

type Resource interface {
	// ResourceType returns the resource type
	ResourceType() string

	// Resource implements the framework Resource interface.
	//
	// For Create/Update/Delete, the implement doesn't need to handle the protocol response, which is done by wrapper.
	// For Read, the implement must handle the protocol response (e.g. set the state).
	//
	// NOTE: Since the Metadata() is implemented by the wrapper, the implement struct shall not implement it.
	// 		 Instead, it is supposed to embed the WithMetadata to meet the interface requirement.
	//
	resource.Resource

	// SetMeta sets the provider meta to the resource.
	// This is implemented by the ImplMeta, the implement struct shall simply embed it.
	SetMeta(meta.Meta)

	// Log logs a message
	// This is implemented by the ImplLog, the implement struct shall simply embed it.
	Log(ctx context.Context, msg string, additionalFields ...map[string]any)

	// The followings are interfaces that a resource can opt-in by implementing them.

	// resource.ResourceWithConfigValidators
	// resource.ResourceWithModifyPlan
	// resource.ResourceWithImportState
	// resource.ResourceWithMoveState
	// resource.ResourceWithUpgradeState
	// resource.ResourceWithValidateConfig
	// ResourceWithTimeout
}
