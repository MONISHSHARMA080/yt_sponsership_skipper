package commonresultchannel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ResponseStatus defines interface for types that can provide a status code
type ResponseStatus interface {
	GetStatusCode() int
	// this method is there to tell me that if for eg the struct is empty then do not send the response
	AreWeNotAllowedToReturnResponse() bool
}

// ResultAndErrorChannel is a generic struct for handling API responses
// The generic type T must implement the ResponseStatus interface
type ResultAndErrorChannel[T ResponseStatus] struct {
	Result T
	Err    error
}

// SendResponse sends the HTTP response to the user
// will return a error if the generic struct is empty and not allowd to send a response
func (r *ResultAndErrorChannel[T]) SendResponse(w http.ResponseWriter) error {
	w.WriteHeader(r.Result.GetStatusCode())
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("the result that is being sent to the user is -> %+v \n\n\n", r.Result)
	if r.Err != nil {
		fmt.Printf("the error in the  -> %s \n", r.Err.Error())
	}
	if r.Result.AreWeNotAllowedToReturnResponse() {
		println("the result is empty and we are not allowed to send a response, we are returning")
		return fmt.Errorf("the result is empty and we are not allowed to send a response")
	}
	// since this is a error response we do not need to send the start time or the end time, as that will be 0
	err := json.NewEncoder(w).Encode(r.Result)
	if err != nil {
		println("error in encoding the json method is -->", err.Error())
	}
	return nil
}

