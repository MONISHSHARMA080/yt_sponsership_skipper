package structs

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ResponseToTheUser struct{
	OrderId  	string  `json:"order_id"`
	Message  	string  `json:"message"`
	StatusCode  int	    `json:"status_code"`
	PlanType    string  `json:"plan_type"`
}

func (resp * ResponseToTheUser) FillTheStruct ( orderId string, messageWeGot string, statusCode int, planType string)  {
	resp.Message =  messageWeGot
	resp.OrderId = orderId
	resp.StatusCode = statusCode
	resp.PlanType = strings.ToLower(planType)
}

func (resp * ResponseToTheUser) ReturnTheErrorInJsonResponse (w http.ResponseWriter, r *http.Request, orderId string, messageWeGot string, statusCode int, planType string)  error{
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)

	resp.FillTheStruct(orderId, messageWeGot, statusCode, planType)

	err:=json.NewEncoder(w).Encode(resp)

	if err != nil {
		println("we have a problem in encoding json and the erro we got in encoding is ->", 
		err.Error(), "\n\n --++-- and the error we were going to send to the user is ->", messageWeGot)
		return err
	}
	println("the response form the server is ->")
	return nil
	
}