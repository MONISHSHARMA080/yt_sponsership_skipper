package structs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"youtubeAdsSkipper/paymentBackendGO/common"
)

// WebhookEvent represents the top-level structure of Razorpay webhook events
type WebhookEvent struct {
	Entity    string         `json:"entity"`
	AccountID string         `json:"account_id"`
	Event     string         `json:"event"`
	Contains  []string       `json:"contains"`
	Payload   PaymentPayload `json:"payload"`
	CreatedAt int64          `json:"created_at"`
}

// PaymentPayload represents the payment payload in the webhook
type PaymentPayload struct {
	Payment PaymentEntity `json:"payment"`
}

// PaymentEntity wraps the payment entity details
type PaymentEntity struct {
	Entity PaymentDetails `json:"entity"`
}

// PaymentDetails contains all possible fields from various payment methods
type PaymentDetails struct {
	ID                string  `json:"id"`
	Entity            string  `json:"entity"`
	Amount            int     `json:"amount"`
	Currency          string  `json:"currency"`
	BaseAmount        int     `json:"base_amount"`
	Status            string  `json:"status"`
	OrderID           string  `json:"order_id"`
	InvoiceID         *string `json:"invoice_id"`
	International     bool    `json:"international"`
	Method            string  `json:"method"` // "netbanking", "card", "wallet", "upi", etc.
	AmountRefunded    int     `json:"amount_refunded"`
	AmountTransferred int     `json:"amount_transferred"`
	RefundStatus      *string `json:"refund_status"`
	Captured          bool    `json:"captured"`
	Description       *string `json:"description"`
	CardID            *string `json:"card_id"`
	Bank              *string `json:"bank"` // For netbanking
	Wallet            *string `json:"wallet"`
	VPA               *string `json:"vpa"` // For UPI
	Email             string  `json:"email"`
	Contact           string  `json:"contact"`
	Notes             struct {
		IDPrimaryKey int64 `json:"id_primary_key"`
	} `json:"notes"`
	Fee              *int         `json:"fee"`
	Tax              *int         `json:"tax"`
	ErrorCode        *string      `json:"error_code"`
	ErrorDescription *string      `json:"error_description"`
	ErrorSource      *string      `json:"error_source"`
	ErrorStep        *string      `json:"error_step"`
	ErrorReason      *string      `json:"error_reason"`
	AcquirerData     AcquirerData `json:"acquirer_data"`
	CreatedAt        int64        `json:"created_at"`
	Card             *CardDetails `json:"card,omitempty"`     // For card payments
	UPI              *UPIDetails  `json:"upi,omitempty"`      // For UPI payments
	TokenID          *string      `json:"token_id,omitempty"` // For tokenized card payments
}

// AcquirerData contains payment processor details that vary by payment method
type AcquirerData struct {
	BankTransactionID *string `json:"bank_transaction_id,omitempty"` // For netbanking
	TransactionID     *string `json:"transaction_id,omitempty"`      // For wallets
	RRN               *string `json:"rrn,omitempty"`                 // For UPI and some card payments
	AuthCode          *string `json:"auth_code,omitempty"`           // For card payments
}

// CardDetails contains details specific to card payments
type CardDetails struct {
	EMI           bool    `json:"emi"`
	Entity        string  `json:"entity"`
	ID            string  `json:"id"`
	IIN           string  `json:"iin"`
	International bool    `json:"international"`
	Issuer        *string `json:"issuer"`
	Last4         string  `json:"last4"`
	Name          string  `json:"name"`
	Network       string  `json:"network"`
	SubType       string  `json:"sub_type"`
	Type          string  `json:"type"`
}

// UPIDetails contains details specific to UPI payments
type UPIDetails struct {
	PayerAccountType string `json:"payer_account_type"`
	VPA              string `json:"vpa"`
	Flow             string `json:"flow"`
}

// Helper functions for parsing webhook data

func (W *WebhookEvent) DecodeJSONResponseInStructAndGetRequestBodyOut(req *http.Request, resultChannelForJsonDecodeErrAndRequestBody chan common.ErrorAndResultStruct[string]) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		resultChannelForJsonDecodeErrAndRequestBody <- common.ErrorAndResultStruct[string]{Result: "", Error: err}
		return
	}

	// Store the body as string for webhook verification
	bodyString := string(bodyBytes)

	// Create a new reader for the JSON decoder since we've read the body
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	err = json.NewDecoder(req.Body).Decode(W)
	if err != nil {
		resultChannelForJsonDecodeErrAndRequestBody <- common.ErrorAndResultStruct[string]{Result: "", Error: err}
		return
	}
	resultChannelForJsonDecodeErrAndRequestBody <- common.ErrorAndResultStruct[string]{Result: bodyString, Error: nil}
}

// GetCreatedTime returns the created time of the payment as a time.Time
func (w *WebhookEvent) GetCreatedTime() time.Time {
	return time.Unix(w.Payload.Payment.Entity.CreatedAt, 0)
}

// GetPaymentDetails returns the payment details from the webhook
func (w *WebhookEvent) GetPaymentDetails() PaymentDetails {
	return w.Payload.Payment.Entity
}

// IsCardPayment returns true if the payment was made via card
func (w *WebhookEvent) IsCardPayment() bool {
	return w.Payload.Payment.Entity.Method == "card"
}

// IsNetbankingPayment returns true if the payment was made via netbanking
func (w *WebhookEvent) IsNetbankingPayment() bool {
	return w.Payload.Payment.Entity.Method == "netbanking"
}

// IsWalletPayment returns true if the payment was made via wallet
func (w *WebhookEvent) IsWalletPayment() bool {
	return w.Payload.Payment.Entity.Method == "wallet"
}

// IsUPIPayment returns true if the payment was made via UPI
func (w *WebhookEvent) IsUPIPayment() bool {
	return w.Payload.Payment.Entity.Method == "upi"
}

// GetBank returns the bank name for netbanking payments
func (w *WebhookEvent) GetBank() string {
	if bank := w.Payload.Payment.Entity.Bank; bank != nil {
		return *bank
	}
	return ""
}

// GetWallet returns the wallet provider for wallet payments
func (w *WebhookEvent) GetWallet() string {
	if wallet := w.Payload.Payment.Entity.Wallet; wallet != nil {
		return *wallet
	}
	return ""
}

// returns true if the payment made by the user is regarding free teir
func (w *WebhookEvent) IsThePaymentForOneTimePaymentTier() (bool, error) {
	recurringPaymentPrice, err := strconv.ParseInt(os.Getenv("RECURRINGPAYMENTPRICE"), 10, 64)
	if err != nil {
		return false, err
	}
	oneTimePayment, err := strconv.ParseInt(os.Getenv("ONETIMEPAYMENTPRICE"), 10, 64)
	if err != nil {
		return false, err
	}

	println("price paid by the user acc to web hook is ->", w.Payload.Payment.Entity.Amount, " and the recurring payment price is:", recurringPaymentPrice, " and the one time payement price is:", oneTimePayment)

	if w.Payload.Payment.Entity.Amount == int(recurringPaymentPrice) {
		return false, nil
	} else if w.Payload.Payment.Entity.Amount == int(oneTimePayment) {
		return true, nil
	}
	return false, fmt.Errorf("can't find the error in the ")
}

// GetVPA returns the VPA for UPI payments
func (w *WebhookEvent) GetVPA() string {
	if vpa := w.Payload.Payment.Entity.VPA; vpa != nil {
		return *vpa
	}
	return ""
}
