package acceptance

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/microsoft/terraform-provider-azuredevops/internal/acceptance/planchecks"
)

// ImportStep returns a Test Step which Imports the Resource, optionally
// ignoring any fields which may not be imported (for example, as they're
// not returned from the API)
func (d TestData) ImportStep(ignore ...string) resource.TestStep {
	resourceAddr := d.ResourceAddr()
	if strings.HasPrefix(resourceAddr, "data.") {
		return resource.TestStep{
			ResourceName: resourceAddr,
			SkipFunc: func() (bool, error) {
				return false, fmt.Errorf("data sources (%q) do not support import - remove the ImportStep / ImportStepFor`", resourceAddr)
			},
		}
	}

	step := resource.TestStep{
		ResourceName:      resourceAddr,
		ImportState:       true,
		ImportStateVerify: true,
	}

	if len(ignore) > 0 {
		step.ImportStateVerifyIgnore = ignore
	}

	return step
}

func (d TestData) MigratePlanStep(config string) resource.TestStep {
	step := resource.TestStep{
		Config: config,
		ConfigPlanChecks: resource.ConfigPlanChecks{
			PreApply: []plancheck.PlanCheck{
				planchecks.IsResourceAction(d.ResourceAddr(), plancheck.ResourceActionNoop),
			},
		},
	}
	return step
}
