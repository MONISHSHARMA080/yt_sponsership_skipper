package routehandlerfunc

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	UserKey string `json:"user_key"`
	// the email is only there to speed up the Db request and you should check it against the key to make sure user did not lied etc.
	EmailByRequestDoNotTrust string `json:"email"`
}

func (request *Request) DecodeJSONRequest(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return err
	}
	return nil
}
