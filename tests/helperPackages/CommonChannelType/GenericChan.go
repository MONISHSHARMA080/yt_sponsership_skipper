package commonchanneltype

type GenericResultChannel[T any] struct {
	Err    error
	Result T
}
