package workitemtrackingprocess

import (
	"github.com/hashicorp/go-cty/cty/gocty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getOrder returns order if there is one defined, otherwise nil
func getOrder(d *schema.ResourceData) (*int, error) {
	rawPlan := d.GetRawPlan()
	if !rawPlan.IsKnown() || rawPlan.IsNull() {
		return nil, nil
	}

	order := rawPlan.GetAttr("order")
	if !order.IsKnown() || order.IsNull() {
		return nil, nil
	}

	var val int
	if err := gocty.FromCtyValue(order, &val); err != nil {
		return nil, err
	}
	return &val, nil
}
