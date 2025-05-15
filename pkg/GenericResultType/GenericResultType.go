package genericresulttype

type ErrorAndResultType[T any] struct {
	Err    error
	Result T
}
