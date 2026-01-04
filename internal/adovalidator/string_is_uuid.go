package adovalidator

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type stringIsUUID struct{}

func (v stringIsUUID) Description(ctx context.Context) string {
	return "validate this in UUID format"
}

func (v stringIsUUID) MarkdownDescription(ctx context.Context) string {
	return "validate this in UUID format"
}

func (_ stringIsUUID) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	str := req.ConfigValue

	if str.IsUnknown() || str.IsNull() {
		return
	}

	if _, errs := isUUID(str.ValueString(), req.Path.String()); len(errs) != 0 {
		for _, err := range errs {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid UUID string",
				err.Error())
		}
	}
}

func StringIsUUID() validator.String {
	return stringIsUUID{}
}

func isUUID(i any, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := uuid.ParseUUID(v); err != nil {
		errors = append(errors, fmt.Errorf("expected %q to be a valid UUID, got %v", k, v))
	}

	return warnings, errors
}
