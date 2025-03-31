package paymentbackendgo

import (
	"fmt"
	"net/http"
	commonstructs "youtubeAdsSkipper/commonStructs"
	"youtubeAdsSkipper/paymentBackendGO/common"
	helperfuncs "youtubeAdsSkipper/paymentBackendGO/helperFuncs"
	"youtubeAdsSkipper/paymentBackendGO/structs"
)

func VerifyPaymentSignature(razorpayKeyID, razorpaySecretID string, envKeyAsByte []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response structs.ResponseVerifyPaymentSignature
		var verifuPaymentLaterFromDB structs.TemporaryFieldToVerifyPaymentLater

		if r.Method != http.MethodPost {
			response.ReturnTheErrorInJsonResponse(w, r, "invalid method", "", http.StatusBadRequest, false)
			return
		}
		var request structs.RequestVerifyPaymentSignature
		// var userKey common.UserKey
		err := request.ParseIntoJson(r)
		if err != nil {
			response.ReturnTheErrorInJsonResponse(w, r, "error decoding JSON", "", http.StatusBadRequest, false)
			return
		}

		channelToDecryptUserKey := make(chan common.ErrorAndResultStruct[string])
		userFormKey := commonstructs.UserKey{}
		db := helperfuncs.DbConnect()
		resultFromGettingTokensFromDbChann := make(chan common.ErrorAndResultStruct[bool])

		go userFormKey.DecryptTheKey(request.UserKey, channelToDecryptUserKey)
		go verifuPaymentLaterFromDB.GetTokens(db, request.Email, resultFromGettingTokensFromDbChann)
		// get the email etc form the key

		// getting the order id form the db, use the request.email to get the order id we will check later to see if it is correct

		resultFromDecryptingKey := <-channelToDecryptUserKey
		if err := resultFromDecryptingKey.Error; err != nil {
			println("the error in decrypting the key is ->", err.Error())
			response.ReturnTheErrorInJsonResponse(w, r, "error is decrypting key", "", http.StatusBadRequest, false)
			return
		}
		if userFormKey.ShouldWeTellUserToGoGetANewKeyPanic() {
			println("\n\n ==the user should be upgraded as it's time ran out ===\n\n ")
			response.ReturnTheErrorInJsonResponse(w, r, "upgrade your key as it's time ran out", "", http.StatusUpgradeRequired, false)
			return
		}

		// checking if the email sent by user matches the email in the key, if not then here is a error
		// err = userKey.SetUserDetail()
		// if err != nil {
		// 	// 400 error as IDK what is causing it the bad key or smth
		// 	println("the error in decrypting the key is ->", err.Error())
		// 	response.ReturnTheErrorInJsonResponse(w, r, "error in getting information out of the key", http.StatusBadRequest, false)
		// 	return
		// }
		if userFormKey.Email == "" || userFormKey.Email != request.Email {
			println("the error is that the email in request is ->", request.Email, "<- and the one form the key is ->", userFormKey.Email, "<- and are they equal ->", userFormKey.Email != request.Email)
			response.ReturnTheErrorInJsonResponse(w, r, "the email sent in the request does not match the one in the key", "", http.StatusBadRequest, false)
			return
		}

		resultFromDBcall := <-resultFromGettingTokensFromDbChann
		if err := resultFromDBcall.Error; err != nil {
			println("the error in going to the db is ->", err.Error())
			response.ReturnTheErrorInJsonResponse(w, r, "error, probally email sent by you is incorrect", "", http.StatusBadRequest, false)
			return
		}

		// if there is a error in Decrypting the key then the user is not authorised and we should send them back
		// now I have gotten a email form the user we should check it against the one form the key too see if it is valid or not

		// println("++++++++++++++++++++++++++++++++++++++++++++")
		// fmt.Printf("\n\n the order id from the DB recurring is %s and the one time is  %s and the one from the request  is %s\n\n", verifuPaymentLaterFromDB.RecurringOrderID, verifuPaymentLaterFromDB.OnetimeOrderID, request.RazorpayOrderId)
		// tmp := request.RazorpayOrderId
		// fmt.Printf("\n ---> is the order id from the req same with recurring %t is it same as one time %t \n", tmp == verifuPaymentLaterFromDB.RecurringOrderID, tmp == verifuPaymentLaterFromDB.OnetimeOrderID)

		orderID := ""
		isUserTierOneTime := false
		if request.RazorpayOrderId == verifuPaymentLaterFromDB.OnetimeOrderID {
			orderID = verifuPaymentLaterFromDB.OnetimeOrderID
			isUserTierOneTime = true
			println("user selected one time payment ")
		} else {
			orderID = verifuPaymentLaterFromDB.RecurringOrderID
			isUserTierOneTime = false
			println("user selected recurring payment")
		}

		if orderID != request.RazorpayOrderId {
			response.ReturnTheErrorInJsonResponse(w, r, "the orderID does not match the one in the DB", "", http.StatusBadRequest, false)
			return
		}

		// if request.DidUserSelectedOneTimePayment {
		// 	println("user selected one time payment")
		// 	orderID = verifuPaymentLaterFromDB.OnetimeOrderID
		// } else {
		// 	println("user selected recurring payment")
		// 	orderID = verifuPaymentLaterFromDB.RecurringOrderID
		// }
		println("Order ID from DB:", orderID)
		println("Order ID from request:", request.RazorpayOrderId)

		signatureGeneratedFromOrderIdStoredInDb, err := helperfuncs.GetGeneratedSignature(orderID, request.RazorpayPaymentId, razorpaySecretID)
		if err != nil {
			println(" the error in generating the signature form the db is ->", err.Error(), " --++-- the signature generated is -> ", signatureGeneratedFromOrderIdStoredInDb)
			response.ReturnTheErrorInJsonResponse(w, r, "signature verification failed", "", http.StatusBadRequest, false)
			return
		}

		fmt.Printf("\n the signature generated form the   -- and the one form the Db/stored one is %s  \n", signatureGeneratedFromOrderIdStoredInDb)

		if signatureGeneratedFromOrderIdStoredInDb != request.RazorpaySignature {
			println("the generate signature is ->", signatureGeneratedFromOrderIdStoredInDb, "++---------- and form the razorpay is ->", request.RazorpaySignature)
			response.ReturnTheErrorInJsonResponse(w, r, "signature verification failed", "", http.StatusBadRequest, false)
			return
		}

		// now we can update the db
		//
		// just give the user generated key that resets after 2 days for the optimistic update state, we will update the key in the message form next day and till
		// then the payment situation will be sorted

		// here is the fake key  to the user

		// get this bool form the one time payment  form the money received
		resultChannel := make(chan common.ErrorAndResultStruct[string])
		helperfuncs.GetFakeKeyForAWhile(&userFormKey, isUserTierOneTime, resultChannel)
		resultForNewUser := <-resultChannel
		if resultForNewUser.Error != nil {
			response.ReturnTheErrorInJsonResponse(w, r, "can't generate a new key", "", http.StatusInternalServerError, false)
			println("error in generating the new key for the user ->", resultForNewUser.Error.Error())
			return
		}
		println("the new key is ->", resultForNewUser.Result)
		// as of now here is the right response
		response.ReturnTheErrorInJsonResponse(w, r, "success", resultForNewUser.Result, http.StatusOK, true)
	}
}
