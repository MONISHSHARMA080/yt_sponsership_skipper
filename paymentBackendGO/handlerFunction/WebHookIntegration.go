package paymentbackendgo

import (
	"net/http"
	"youtubeAdsSkipper/paymentBackendGO/common"
	"youtubeAdsSkipper/paymentBackendGO/structs"

	"github.com/razorpay/razorpay-go/utils"
)

// function is used to take in the web hook message and  if everything is right then update the DB to insert the message
func WebHookIntegrationForPaymentCapture(razorpayKeyID, razorpaySecretID, webHookSecret string, envKeyAsByte []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		println("in the webHookEvent")

		// Decode the webHookEvent
		var webhookEvent structs.WebhookEvent
		resultFromJSONDecode := make(chan common.ErrorAndResultStruct[string])
		go webhookEvent.DecodeJSONResponseInStructAndGetRequestBodyOut(r, resultFromJSONDecode)

		// getting the webhook
		webhookSignature := r.Header.Get("X-Razorpay-Signature")
		if webhookSignature == "" {
			http.Error(w, "Missing webhook signature", http.StatusBadRequest)
			return
		}

		webHookEventReceiveChannel := <-resultFromJSONDecode
		if webHookEventReceiveChannel.Error != nil {
			println("there is a error in the webHookEvent Decoding(in json) ->", webHookEventReceiveChannel.Error.Error())
			return
		}
		webHookBodyFromReqAsAString := webHookEventReceiveChannel.Result
		// Process the webhook event...
		// validate the webHook event and if it is correct then go into the db
		// razorPayClient := razorpay.NewClient(razorpayKeyID, razorpaySecretID, webHookSecret)
		isWebHookCorrect := utils.VerifyWebhookSignature(webHookBodyFromReqAsAString, webhookSignature, webHookSecret)
		if !isWebHookCorrect {
			println("---- webhook is not form razorpay ----")
			return
		}
		// if webhookEvent.Payload.Payment.Entity.Currency !=
		// how do I get the tier of the request

		ispaymentForOneTimeOnly, err := webhookEvent.IsThePaymentForOneTimePaymentTier()
		if err != nil {
			println("error in knowing if the event is free tier or not and  the error is ->", err.Error())
			return
		}
		println("\n\n\n--------webHook is from the razorpay indeed and now we are going to make the call to DB to set a message there for the user, the ammount paid by the user is ->",
			webhookEvent.Payload.Payment.Entity.Amount, " and the Currency is ", webhookEvent.Payload.Payment.Entity.Currency, "and the payment for one time only ", ispaymentForOneTimeOnly, "\n\n\n")

		// now make sure the ammount paid is correct and if everything is alright if it is then set the message in the DB
		//
		//
		//
		//for inserting in the message table we would need to use the onconflict update the, or SELECT COALESCE(MAX(version), 0-- depending
		// upon if I want many rows for the user or just a single one(many)
	}
}
