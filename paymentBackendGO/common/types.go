// paymentBackendGO/common/types.go
package common

type ErrorAndResultStruct[T any] struct {
    Error         error
    Result T
}