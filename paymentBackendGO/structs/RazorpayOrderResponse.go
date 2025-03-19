package structs

import (
	"encoding/json"
	"youtubeAdsSkipper/paymentBackendGO/common"

	"github.com/razorpay/razorpay-go"
)

type RazorpayOrderResponse struct {
	ID             string        `json:"id"`                        // The unique identifier of the order
	Entity         string        `json:"entity"`                    // Type of entity (always "order" for orders)
	Amount         int           `json:"amount"`                    // Payment amount in smallest currency sub-unit
	AmountPaid     int           `json:"amount_paid"`               // The amount paid against the order
	AmountDue      int           `json:"amount_due"`                // The amount pending against the order
	Currency       string        `json:"currency"`                  // ISO code for the currency (e.g., "INR")
	Receipt        string        `json:"receipt"`                   // Receipt number that corresponds to this order
	Status         string        `json:"status"`                    // Status of the order: created, attempted, or paid
	Attempts       int           `json:"attempts"`                  // Number of payment attempts made against this order
	Notes          []interface{} `json:"notes"`                     // Key-value pairs for additional information
	OfferID        interface{}   `json:"offer_id"`                  // ID of the offer applied, null if none
	CreatedAt      int64         `json:"created_at"`                // Unix timestamp of when the order was created
	PartialPayment bool          `json:"partial_payment,omitempty"` // Whether partial payment is allowed
}

func (rpResp *RazorpayOrderResponse) convertResponseToJSON(responseBody map[string]interface{}) error {
	jsonBytes, err := json.Marshal(responseBody)
	if err != nil {
		return err
	}

	// Then unmarshal into our struct
	if err := json.Unmarshal(jsonBytes, rpResp); err != nil {
		return err
	}
	return nil
}

func (respRPay *RazorpayOrderResponse) AskRazorpayForTheOrderID(client *razorpay.Client, ammount int64, resultChannel chan common.ErrorAndResultStruct[string])  {
	// client := razorpay.NewClient(razorpayKeyID, razorpaySecretID)

	data := map[string]interface{}{
		"amount":   ammount, // Amount is in currency subunits. Default currency is INR. Hence, 50000 refers to 50000 paise
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		println("error is ->.>", err.Error())
		resultChannel <- common.ErrorAndResultStruct[string]{Error: err, Result: ""}
		return 
	}

	err = respRPay.convertResponseToJSON(body)
	if err != nil {
		println("error is ->.>", err.Error())
		resultChannel <- common.ErrorAndResultStruct[string]{Error: err, Result: ""}
		return 
	}

		resultChannel <- common.ErrorAndResultStruct[string]{Error: nil, Result: respRPay.Status + string(respRPay.Attempts)}

	// return &RazorpayOrderResponse, nil
}
