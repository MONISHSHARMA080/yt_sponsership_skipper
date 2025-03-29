package routehandlerfunc

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	NewEncryptedKey string `json:"new_encrypted_key"`
	Message         string `json:"message"`
	StatusCode      int    `json:"status_code"`
}

func (resp *Response) ReturnJSONResponse(w http.ResponseWriter, newKey, message string, statusCode int) error {
	resp.Message = message
	resp.NewEncryptedKey = newKey
	resp.StatusCode = statusCode

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// resp.FillTheStruct(orderIdRecurr, orderIdOneTime, messageWeGot, statusCode)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		println("we have a problem in encoding json and the erro we got in encoding is ->",
			err.Error(), "\n\n --++-- and the error we were going to send to the user is ->", message)
		return err
	}
	return nil
}
