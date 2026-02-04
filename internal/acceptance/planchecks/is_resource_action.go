package planchecks

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var _ plancheck.PlanCheck = isResourceAction{}

type isResourceAction struct {
	resourceAddress string
	actionType      plancheck.ResourceActionType
}

// CheckPlan implements the plan check logic.
func (e isResourceAction) CheckPlan(ctx context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	for _, rc := range req.Plan.ResourceChanges {
		if e.resourceAddress == rc.Address {
			switch e.actionType {
			case plancheck.ResourceActionNoop:
				if !rc.Change.Actions.NoOp() {
					resp.Error = fmt.Errorf("'%s' - expected action to be %s", rc.Address, e.actionType)
					return
				}
			case plancheck.ResourceActionCreate:
				if !rc.Change.Actions.Create() {
					resp.Error = fmt.Errorf("'%s' - expected action to be %s", rc.Address, e.actionType)
					return
				}
			case plancheck.ResourceActionRead:
				if !rc.Change.Actions.Read() {
					resp.Error = fmt.Errorf("'%s' - expected action to be %s", rc.Address, e.actionType)
					return
				}
			case plancheck.ResourceActionUpdate:
				if !rc.Change.Actions.Update() {
					resp.Error = fmt.Errorf("'%s' - expected action to be %s", rc.Address, e.actionType)
					return
				}
			case plancheck.ResourceActionDestroy:
				if !rc.Change.Actions.Delete() {
					resp.Error = fmt.Errorf("'%s' - expected action to be %s", rc.Address, e.actionType)
					return
				}
			case plancheck.ResourceActionDestroyBeforeCreate:
				if !rc.Change.Actions.DestroyBeforeCreate() {
					resp.Error = fmt.Errorf("'%s' - expected action to be %s", rc.Address, e.actionType)
					return
				}
			case plancheck.ResourceActionCreateBeforeDestroy:
				if !rc.Change.Actions.CreateBeforeDestroy() {
					resp.Error = fmt.Errorf("'%s' - expected action to be %s", rc.Address, e.actionType)
					return
				}
			case plancheck.ResourceActionReplace:
				if !rc.Change.Actions.Replace() {
					resp.Error = fmt.Errorf("'%s' - expected action to be %s", rc.Address, e.actionType)
					return
				}
			default:
				resp.Error = fmt.Errorf("%s - unexpected ResourceActionType: %s", rc.Address, e.actionType)
				return
			}
		}
	}

	// If the resource we're looking for in the plan doesn't exist we don't return an error since we have tests in the provider
	// that rely on a setup configuration before the pertinent test step is applied
}

// IsResourceAction returns a plan check that asserts that a given resource will have a specific resource change type in the plan.
// Valid actionType are an enum of type plancheck.ResourceActionType, examples: NoOp, DestroyBeforeCreate, Update (in-place), etc.
func IsResourceAction(resourceAddress string, actionType plancheck.ResourceActionType) plancheck.PlanCheck {
	return isResourceAction{
		resourceAddress: resourceAddress,
		actionType:      actionType,
	}
}
