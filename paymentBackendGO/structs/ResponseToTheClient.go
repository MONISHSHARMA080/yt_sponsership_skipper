package structs

import (
	"encoding/json"
	"net/http"
)

type ResponseToTheUser struct{
	OrderIdForRecurring  	string  `json:"order_id_for_recurring"`
	OrderIdForOneTime	  	string  `json:"order_id_for_onetime"`
	Message  				string  `json:"message"`
	StatusCode  			int	    `json:"status_code"`
}

func (resp * ResponseToTheUser) FillTheStruct ( orderIdRecurr, orderIdOneTime string, messageWeGot string, statusCode int, )  {
	resp.Message =  messageWeGot
	resp.OrderIdForOneTime = orderIdOneTime
	resp.OrderIdForRecurring = orderIdRecurr
	resp.StatusCode = statusCode
	// resp.PlanType = strings.ToLower(planType)
}

func (resp * ResponseToTheUser) ReturnTheErrorInJsonResponse (w http.ResponseWriter, r *http.Request, orderIdRecurr, orderIdOneTime string, messageWeGot string, statusCode int, )  error{
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)

	resp.FillTheStruct(orderIdRecurr, orderIdOneTime, messageWeGot, statusCode, )
		println(" the oder id is (one time)->", orderIdOneTime,"  -- recurring ->", orderIdRecurr )

	err:=json.NewEncoder(w).Encode(resp)

	if err != nil {
		println("we have a problem in encoding json and the erro we got in encoding is ->", 
		err.Error(), "\n\n --++-- and the error we were going to send to the user is ->", messageWeGot)
		return err
	}
	println("the response form the server is, recurring id  ->",resp.OrderIdForRecurring, "--one time ->",resp.OrderIdForOneTime)
	return nil
	
}