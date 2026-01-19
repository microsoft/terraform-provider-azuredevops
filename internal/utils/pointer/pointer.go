package pointer

func To[T any](v *T) T {
	var zero T
	if v != nil {
		return *v
	}
	return zero
}

func From[T any](v T) *T {
	return &v
}
