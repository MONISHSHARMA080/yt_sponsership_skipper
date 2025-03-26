package helperfuncs

import (
	"encoding/json"
	"fmt"

	"github.com/razorpay/razorpay-go"
)

type RefundResult struct {
	Success *SuccessRefund `json:"success,omitempty"`
	Error   *RefundError   `json:"error,omitempty"`
}

// SuccessRefund represents a successful refund response
type SuccessRefund struct {
	ID             string                 `json:"id"`
	Entity         string                 `json:"entity"`
	Amount         int64                  `json:"amount"`
	Currency       string                 `json:"currency"`
	PaymentID      string                 `json:"payment_id"`
	Receipt        *string                `json:"receipt"`
	Notes          []interface{}          `json:"notes"`
	AcquirerData   map[string]interface{} `json:"acquirer_data"`
	CreatedAt      int64                  `json:"created_at"`
	BatchID        *string                `json:"batch_id"`
	Status         string                 `json:"status"`
	SpeedProcessed string                 `json:"speed_processed"`
	SpeedRequested string                 `json:"speed_requested"`
}

// RefundError represents an error response from Razorpay
type RefundError struct {
	Code        string                 `json:"code"`
	Description string                 `json:"description"`
	Source      string                 `json:"source"`
	Step        string                 `json:"step"`
	Reason      string                 `json:"reason"`
	Metadata    map[string]interface{} `json:"metadata"`
	Field       string                 `json:"field"`
}

// this function is designed to be used with the webhook, such that when the payment fails or we fail to add stuff to the
// DB, so we will just simply refund to the user
func RefundTheUser(ammountToRefund int64, paymentID string, razorPayClient *razorpay.Client) (*RefundResult, error) {
	data := map[string]interface{}{
		"speed": "normal",
		"notes": "",
	}

	resultFromRazorPay, err := razorPayClient.Payment.Refund(paymentID, int(ammountToRefund), data, nil)
	if err != nil {
		return nil, err
	}
	preetyPrintJson(resultFromRazorPay)
	var result RefundResult

	// Convert resultFromRazorPay to JSON
	jsonData, err := json.Marshal(resultFromRazorPay)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal refund response: %v", err)
	}

	// Unmarshal JSON into RefundResult
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refund response: %v", err)
	}

	// Check if the result has an error
	if result.Error != nil {
		return &result, fmt.Errorf("refund error: %s", result.Error.Description)
	}

	// If no error, ensure Success is populated
	if result.Success == nil {
		// Directly populate Success if it's not set
		var successRefund SuccessRefund
		err = json.Unmarshal(jsonData, &successRefund)
		if err != nil {
			return nil, fmt.Errorf("failed to parse success refund: %v", err)
		}
		result.Success = &successRefund
	}

	return &result, nil
}

// if err don't return it
func preetyPrintJson(data interface{}) {
	preetyJsonInByte, err := json.MarshalIndent(data, "", "       ")
	if err != nil {
		println("there is a error in printing the preety json in refund struct and it is ->", err.Error())
		return
	}
	println("\n\n\n\n the preety json resonse of refund form razorpay is ->: ", string(preetyJsonInByte), "\n\n\n\n\n")
}

// this function is a refund function wrapper, meant to be used such when we want to refund you call this function and if there is
// a error we are not implementing the refunding logic, thats why I do not want to deal with it
func AbstractRefundFunctionWrapper(ammountToRefund int, paymentID string, razorPayClient *razorpay.Client) {
	refundResult, err := RefundTheUser(int64(ammountToRefund), paymentID, razorPayClient)
	if err != nil {
		println("there was a error in calling the refunc function form the razorpay ->", err.Error())
		return
	}
	if refundResult.Error != nil {
		println("well there is a errorn in the repsonse of the refund func and the code is->", refundResult.Error.Code, "\n and the description is ->", refundResult.Error.Description)
		return
	} else {
		println("the result form the refund func is success and ammount is ->", refundResult.Success.Amount)
		return
	}
}
