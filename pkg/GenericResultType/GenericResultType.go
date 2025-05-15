package genericResultType

type ErrorAndResultType[T any] struct {
	Err    error
	Result T
}
