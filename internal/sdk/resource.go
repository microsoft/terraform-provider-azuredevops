package sdk

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ResourceTimeout specifies the default timeout value for each operation.
type ResourceTimeout struct {
	Create time.Duration
	Read   time.Duration
	Update time.Duration
	Delete time.Duration
}

type Resource interface {
	Timeout() ResourceTimeout

	// Enforced interfaces, which means every resource must implement them.
	resource.Resource
	resource.ResourceWithConfigure

	// The followings are interfaces that a resource can opt-in.
	// resource.ResourceWithConfigValidators
	// resource.ResourceWithModifyPlan
	// resource.ResourceWithImportState
	// resource.ResourceWithMoveState
	// resource.ResourceWithUpgradeState
	// resource.ResourceWithValidateConfig
}
