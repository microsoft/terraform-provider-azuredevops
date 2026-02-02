package framework

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ResourceTimeout specifies the default timeout value for each operation.
type ResourceTimeout struct {
	Create time.Duration
	Read   time.Duration
	Update time.Duration
	Delete time.Duration
}

type ResourceIdentity interface {
	// Convert from an import id to the identity
	FromId(id string) error
	// Fields returns each identity field value, together with its path and its corresponding state path.
	Fields() []IdentityField
}

type IdentityField struct {
	PathIdentity path.Path
	PathState    path.Path
	Value        attr.Value
}

// Resource interface defines the mandatory methods that a resource requires to implement.
// Some of the method can be implemented by embedding a utility struct (see the comments).
type Resource interface {
	BaseResource

	// Identity returns a ResourceIdentity.
	Identity() ResourceIdentity

	// Resource implements the framework Resource interface.
	//
	// For Create, the implement doesn't need to handle the protocol response if the req.Plan contains all the information
	//             for a follow-up Read(). Otherwise, the implement shall set those Read() related attributes to the state.
	// For Update/Delete, the implement doesn't need to handle the protocol response, which is done by wrapper.
	// For Read, the implement must handle the protocol response (e.g. set the state).
	//
	// NOTE: Since the Metadata() is implemented by the wrapper, the implement struct shall not implement it.
	// 		 Instead, it is supposed to embed the WithMetadata to meet the interface requirement.
	//
	resource.ResourceWithIdentity
}

// ResourceWithTimeout is an opt-in interface that can implement customized timeout.
type ResourceWithTimeout interface {
	resource.Resource

	// Timeout returns the timeout for each operation.
	Timeout() ResourceTimeout
}

type ResourceWithPostCreate interface {
	resource.Resource

	ShouldPostCreate(ctx context.Context, req resource.CreateRequest) bool
	PostCreate(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse)
}

type ResourceWithPostUpdate interface {
	resource.Resource

	ShouldPostUpdate(ctx context.Context, req resource.UpdateRequest) bool
	PostUpdate(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse)
}

// Additionally, a resource can opt-in any of the following interfaces.
//
// resource.ResourceWithConfigValidators
// resource.ResourceWithModifyPlan
// resource.ResourceWithMoveState
// resource.ResourceWithUpgradeState
// resource.ResourceWithValidateConfig
// resource.ResourceWithUpgradeIdentity
