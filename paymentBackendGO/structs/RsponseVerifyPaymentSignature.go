package structs

import (
	"encoding/json"
	"net/http"
)

type ResponseVerifyPaymentSignature struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
	// NewKey     string  `json:"ney_key"` // maybe we should return the key here
}

func (resp *ResponseVerifyPaymentSignature) FillTheStruct(Success bool, messageWeGot string, statusCode int) {
	resp.Message = messageWeGot
	resp.StatusCode = statusCode
	resp.Success = Success
}

func (resp *ResponseVerifyPaymentSignature) ReturnTheErrorInJsonResponse(w http.ResponseWriter, r *http.Request, messageWeGot string, statusCode int, success bool) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp.FillTheStruct(success, messageWeGot, statusCode)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		println("we have a problem in encoding json and the erro we got in encoding is ->",
			err.Error(), "\n\n --++-- and the error we were going to send to the user is ->", messageWeGot)
		return err
	}
	println("the response form the server is,->", resp.StatusCode, "--message to the user ->", resp.Message, " is success ->", resp.Success)
	return nil
}
