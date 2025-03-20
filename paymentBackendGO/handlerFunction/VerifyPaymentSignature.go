package paymentbackendgo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"youtubeAdsSkipper/paymentBackendGO/common"
	helperfuncs "youtubeAdsSkipper/paymentBackendGO/helperFuncs"
	"youtubeAdsSkipper/paymentBackendGO/structs"
)

func VerifyPaymentSignature(razorpayKeyID, razorpaySecretID string, envKeyAsByte []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response structs.ResponseVerifyPaymentSignature
		var verifuPaymentLaterFromDB structs.TemporaryFieldToVerifyPaymentLater

		if r.Method != http.MethodPost {
			response.ReturnTheErrorInJsonResponse(w, r, "invalid method", http.StatusBadRequest, false)
			return
		}
		var request structs.RequestVerifyPaymentSignature
		var userKey common.UserKey
		err := request.ParseIntoJson(r)
		if err != nil {
			response.ReturnTheErrorInJsonResponse(w, r, "error decoding JSON", http.StatusBadRequest, false)
			return
		}

		db := helperfuncs.DbConnect()
		resultFromGettingTokensFromDbChann := make(chan common.ErrorAndResultStruct[bool])

		// get the email etc form the key
		userKey.EncryptedUserKey = request.UserKey
		channForKeyResult := make(chan common.ErrorAndResultStruct[string])

		go userKey.DecryptKey(envKeyAsByte, channForKeyResult)
		// getting the order id form the db, use the request.email to get the order id we will check later to see if it is correct

		go verifuPaymentLaterFromDB.GetTokens(db, request.Email, resultFromGettingTokensFromDbChann)

		resultFromDecryptingKey := <-channForKeyResult
		if err := resultFromDecryptingKey.Error; err != nil {
			println("the error in decrypting the key is ->", err.Error())
			response.ReturnTheErrorInJsonResponse(w, r, "error is decrypting key", http.StatusBadRequest, false)
			return
		}

		// checking if the email sent by user matches the email in the key, if not then here is a error
		err = userKey.SetUserDetail()
		if err != nil {
			// 400 error as IDK what is causing it the bad key or smth
			println("the error in decrypting the key is ->", err.Error())
			response.ReturnTheErrorInJsonResponse(w, r, "error in getting information out of the key", http.StatusBadRequest, false)
			return
		}
		if userKey.UserInTheDb.Email != "" && userKey.UserInTheDb.Email != request.Email {
			println("the error is that the email in request is ->",request.Email,"<- and the one form the key is ->", userKey.UserInTheDb.Email,"<- and are they equal ->",userKey.UserInTheDb.Email != request.Email)
			response.ReturnTheErrorInJsonResponse(w, r, "the email does not match the one in the key", http.StatusBadRequest, false)
			return
		}

		resultFromDBcall := <-resultFromGettingTokensFromDbChann
		if err := resultFromDBcall.Error; err != nil {
			println("the error in going to the db is ->", err.Error())
			response.ReturnTheErrorInJsonResponse(w, r, "error, probally email sent by you is incorrect", http.StatusBadRequest, false)
			return
		}

		// if there is a error in Decrypting the key then the user is not authorised and we should send them back
		// now I have gotten a email form the user we should check it against the one form the key too see if it is valid or not
		//

		// do not
		var orderID string
		if request.DidUserSelectedOneTimePayment {
			println("user selected one time payment")
			orderID = verifuPaymentLaterFromDB.OnetimeOrderID
		} else {
			println("user selected recurring payment")
			orderID = verifuPaymentLaterFromDB.RecurringOrderID
		}
		println("Order ID from DB:", orderID)
		println("Order ID from request:", request.RazorpayOrderId)
		data := request.RazorpayOrderId + "|" + request.RazorpayPaymentId

		h := hmac.New(sha256.New, []byte(razorpaySecretID))

		// Write the data to the HMAC
		intReturned, err := h.Write([]byte(data))
		if err != nil {
			response.ReturnTheErrorInJsonResponse(w, r, "signature verification failed", http.StatusBadRequest, false)
			return
		}
		println("the int returned is ->", intReturned)

		generatedSignature := hex.EncodeToString(h.Sum(nil))
		if generatedSignature != request.RazorpaySignature {
			println("the generate signature is ->", generatedSignature, "++---------- and form the razorpay is ->", request.RazorpaySignature)
			response.ReturnTheErrorInJsonResponse(w, r, "signature verification failed", http.StatusBadRequest, false)
			return
		}

			println("the generate signature is ->", generatedSignature, "++---------- and form the razorpay is ->", request.RazorpaySignature)

		// now we can update the db
		// see what the webhook returns and if it is  same shit(based on that response) make the db table and update it here
		//
		//
		//
		//
		// as of now here is the right response
		response.ReturnTheErrorInJsonResponse(w, r, "success", http.StatusOK, true)
	}
}
