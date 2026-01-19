package fwtype

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func BoolValue[T ~bool](v *T) types.Bool {
	if v == nil {
		return types.BoolNull()
	}
	return types.BoolValue(bool(*v))
}

func BoolValueOrZero[T ~bool](v *T) types.Bool {
	if v == nil {
		return types.BoolValue(false)
	}
	return types.BoolValue(bool(*v))
}

func Int32Value[T ~int32](v *T) types.Int32 {
	if v == nil {
		return types.Int32Null()
	}
	return types.Int32Value(int32(*v))
}

func Int32ValueOrZero[T ~int32](v *T) types.Int32 {
	if v == nil {
		return types.Int32Value(0)
	}
	return types.Int32Value(int32(*v))
}

func Int64Value[T ~int64](v *T) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*v))
}

func Int64ValueOrZero[T ~int64](v *T) types.Int64 {
	if v == nil {
		return types.Int64Value(0)
	}
	return types.Int64Value(int64(*v))
}

func StringValue[T ~string](v *T) types.String {
	if v == nil {
		return types.StringNull()
	}
	return types.StringValue(string(*v))
}

func StringValueOrZero[T ~string](v *T) types.String {
	if v == nil {
		return types.StringValue("")
	}
	return types.StringValue(string(*v))
}
