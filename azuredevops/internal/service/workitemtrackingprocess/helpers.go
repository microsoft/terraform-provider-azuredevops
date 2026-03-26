package workitemtrackingprocess

import (
	"github.com/hashicorp/go-cty/cty/gocty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getBoolAttributeFromConfig returns a *bool from the configuration for the given key.
// Returns nil if the key is not set in the configuration, allowing
// distinction between an unset value and an explicit false.
func getBoolAttributeFromConfig(d *schema.ResourceData, key string) (*bool, error) {
	rawPlan := d.GetRawPlan()
	if !rawPlan.IsKnown() || rawPlan.IsNull() {
		return nil, nil
	}

	value := rawPlan.GetAttr(key)
	if !value.IsKnown() || value.IsNull() {
		return nil, nil
	}

	var val bool
	if err := gocty.FromCtyValue(value, &val); err != nil {
		return nil, err
	}
	return &val, nil
}
