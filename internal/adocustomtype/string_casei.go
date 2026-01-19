package adocustomtype

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// The Schema Type
var _ basetypes.StringTypable = StringCaseInsensitiveType{}

type StringCaseInsensitiveType struct {
	basetypes.StringType
}

func (t StringCaseInsensitiveType) Equal(o attr.Type) bool {
	other, ok := o.(StringCaseInsensitiveType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t StringCaseInsensitiveType) String() string {
	return "StringCaseInsensitive"
}

func (t StringCaseInsensitiveType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := StringCaseInsensitiveValue{
		StringValue: in,
	}

	return value, nil
}

func (t StringCaseInsensitiveType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func (t StringCaseInsensitiveType) ValueType(ctx context.Context) attr.Value {
	return StringCaseInsensitiveValue{}
}

// The Value Type
var _ basetypes.StringValuable = StringCaseInsensitiveValue{}
var _ basetypes.StringValuableWithSemanticEquals = StringCaseInsensitiveValue{}

type StringCaseInsensitiveValue struct {
	basetypes.StringValue
}

func (v StringCaseInsensitiveValue) StringSemanticEquals(ctx context.Context, o basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	oval, ok := o.(StringCaseInsensitiveValue)

	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", o),
		)
		return false, diags
	}

	return strings.EqualFold(oval.ValueString(), v.ValueString()), diags
}

func (v StringCaseInsensitiveValue) Equal(o attr.Value) bool {
	other, ok := o.(StringCaseInsensitiveValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v StringCaseInsensitiveValue) Type(ctx context.Context) attr.Type {
	return StringCaseInsensitiveType{}
}
