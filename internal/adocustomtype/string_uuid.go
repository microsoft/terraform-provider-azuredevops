package adocustomtype

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// The Schema Type
var _ basetypes.StringTypable = StringUUIDType{}

type StringUUIDType struct {
	basetypes.StringType
}

func (t StringUUIDType) Equal(o attr.Type) bool {
	other, ok := o.(StringUUIDType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t StringUUIDType) String() string {
	return "StringUUID"
}

func (t StringUUIDType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := StringUUIDValue{
		StringValue: in,
	}
	if in.ValueString() != "" {
		var err error
		value.UUID, err = uuid.Parse(in.ValueString())
		if err != nil {
			return value, diag.Diagnostics{diag.NewErrorDiagnostic("Parsing UUID", err.Error())}
		}
	}
	return value, nil
}

func (t StringUUIDType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
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

func (t StringUUIDType) ValueType(ctx context.Context) attr.Value {
	return StringUUIDValue{}
}

// The Value Type
var _ basetypes.StringValuable = StringUUIDValue{}

type StringUUIDValue struct {
	UUID uuid.UUID
	basetypes.StringValue
}

func (v StringUUIDValue) Equal(o attr.Value) bool {
	other, ok := o.(StringUUIDValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v StringUUIDValue) Type(ctx context.Context) attr.Type {
	return StringUUIDType{}
}
