package stride

type Optional[T any] struct {
	Value T
	Valid bool
}

func NewOptional[T any](v T, valid bool) Optional[T] {
	return Optional[T]{Value: v, Valid: valid}
}
