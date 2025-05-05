// This is a copy of github.com/microsoft/azure-devops-go-api/azuredevops/pipelineschecks/models.go
// The existing version does not contain the "Timeout" property on the CheckConfiguration struct

// This file cannot be under "internal", because azdosdkmocks/pipelines_checks_v5_extras_mock.go depends on it.

package dashboardextras

import (
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/dashboard"
)

type UpdateDashboardArgs struct {
	// (required) The initial state of the dashboard
	Dashboard *dashboard.Dashboard
	// (required) Project ID or project name
	Project *string
	// (Optional) Team I
	Team *string
}
