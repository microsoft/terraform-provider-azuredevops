package framework

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// DataSourceTimeout specifies the default timeout value for the data source
type DataSourceTimeout struct {
	Read time.Duration
}

// DataSource interface defines the mandatory methods that a data source requires to implement.
// Some of the method can be implemented by embedding a utility struct.
type DataSource interface {
	BaseResource
	datasource.DataSource
}

// DataSourceWithTimeout is an opt-in interface that can implement customized timeout.
type DataSourceWithTimeout interface {
	datasource.DataSource

	// Timeout returns the timeout for each operation.
	Timeout() DataSourceTimeout
}

// Additionally, a data source can opt-in any of the following interfaces.
//
// datasource.DataSourceWithConfigValidators
// datasource.DataSourceWithValidateConfig
