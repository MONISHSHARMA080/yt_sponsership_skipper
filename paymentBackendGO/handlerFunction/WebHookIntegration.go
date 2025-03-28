package paymentbackendgo

import (
	"net/http"
	"youtubeAdsSkipper/paymentBackendGO/common"
	helperfuncs "youtubeAdsSkipper/paymentBackendGO/helperFuncs"
	"youtubeAdsSkipper/paymentBackendGO/structs"

	"github.com/razorpay/razorpay-go"
	"github.com/razorpay/razorpay-go/utils"
)

// function is used to take in the web hook message and  if everything is right then update the DB to insert the message
func WebHookIntegrationForPaymentCapture(razorpayKeyID, razorpaySecretID, webHookSecret string, envKeyAsByte []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		println("in the webHookEvent")

		// Decode the webHookEvent
		var webHookEvent structs.WebhookEvent
		var messageFormDbOnPaymentCapture structs.MessageForUserOnPaymentCapture
		resultChannelToGetToDB := make(chan common.ErrorAndResultStruct[bool])
		resultFromJSONDecode := make(chan common.ErrorAndResultStruct[string])

		razorPayClient := razorpay.NewClient(razorpayKeyID, razorpaySecretID)
		go webHookEvent.DecodeJSONResponseInStructAndGetRequestBodyOut(r, resultFromJSONDecode)

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

		// here if I get the error maybe I need to make a DB call  and check the order id against it to see the payment type, but for as
		// of now this is ok
		//
		//I am probally making the call to the DB for the getting the order ID form it, might as well check it
		ispaymentForOneTimeOnly, err := webHookEvent.IsThePaymentForOneTimePaymentTier()
		if err != nil {
			println("error in knowing if the event is free tier or not and  the error is ->", err.Error())
			// return
			// assuming it to be true, such that you will implement it later
			ispaymentForOneTimeOnly = true
			println("\n\n\n\n\n\n\n---------WE ARE NOT ABLE TO DETERMINE IF THE PAYMENT IS PAID OR NOT, SO WE ARE MAKING IT PAID ;;;FIX THIS NOW--------------\n\n\n\n\n\n\n\n\n\n")
			println("error in getting payment for one time in webhookEvent is ->", err.Error())
			helperfuncs.AbstractRefundFunctionWrapper(webHookEvent.Payload.Payment.Entity.Amount, webHookEvent.Payload.Payment.Entity.ID, razorPayClient)
		}
		println("\n\n\n--------webHook is from the razorpay indeed and now we are going to make the call to DB to set a message there for the user, the ammount paid by the user is ->",
			webHookEvent.Payload.Payment.Entity.Amount, " and the Currency is ", webHookEvent.Payload.Payment.Entity.Currency, "and the payment for one time only ", ispaymentForOneTimeOnly, "\n\n\n")

		// if  ispaymentForOneTimeOnly is true make a fucn on it to make it after the day or months form the env
		// should it be on my webHookEvent or on my InitializeStruct method and passing a bool on whether it is a free tier
		// or not should set it
		err = messageFormDbOnPaymentCapture.InitializeStruct(webHookEvent.Payload.Payment.Entity.ID, webHookEvent.Payload.Payment.Entity.Notes.IDPrimaryKey, ispaymentForOneTimeOnly)
		if err != nil {
			println("there is a error in Initializing the Struct and we can't send messgae to the Db ->", err.Error())
			helperfuncs.AbstractRefundFunctionWrapper(webHookEvent.Payload.Payment.Entity.Amount, webHookEvent.Payload.Payment.Entity.ID, razorPayClient)
			return
		}
		go messageFormDbOnPaymentCapture.AddMessageAfterUserPaymentReceived(helperfuncs.DbConnect(), resultChannelToGetToDB)
		resultFormTheDB := <-resultChannelToGetToDB
		if resultFormTheDB.Error != nil {
			println("there  is a performing  operations in the DB and it is -> ", resultFormTheDB.Error.Error())
			helperfuncs.AbstractRefundFunctionWrapper(webHookEvent.Payload.Payment.Entity.Amount, webHookEvent.Payload.Payment.Entity.ID, razorPayClient)
			return
		}
	}
}
