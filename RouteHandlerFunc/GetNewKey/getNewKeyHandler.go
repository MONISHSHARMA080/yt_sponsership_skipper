package routehandlerfunc

import (
	"fmt"
	"net/http"
	commonstructs "youtubeAdsSkipper/commonStructs"
	"youtubeAdsSkipper/paymentBackendGO/common"
	helperfuncs "youtubeAdsSkipper/paymentBackendGO/helperFuncs"
	"youtubeAdsSkipper/paymentBackendGO/structs"
)

func GetNewKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response
		var Request Request

		if r.Method != http.MethodPost {
			// return json response
			response.ReturnJSONResponse(w, "", "method not allowed", http.StatusBadRequest)
			return
		}
		err := Request.DecodeJSONRequest(r)
		if err != nil {
			println("there is a error in json decoding the response", err.Error())
			response.ReturnJSONResponse(w, "", "can't decode the json", http.StatusBadRequest)
			return
		}
		// get the key string into the user struct and also get the
		oldUser := commonstructs.UserKey{}
		resultUserChannel := make(chan common.ErrorAndResultStruct[string])
		resultDBChannel := make(chan common.ErrorAndResultStruct[bool])
		var DBStruct structs.MessageForUserOnPaymentCapture
		db := helperfuncs.DbConnect()

		go oldUser.DecryptTheKey(Request.UserKey, resultUserChannel)
		go DBStruct.GetLatestMessageForTheUser(db, Request.EmailByRequestDoNotTrust, resultDBChannel)

		userDecryptionResult := <-resultUserChannel
		if userDecryptionResult.Error != nil {
			println("there is a error in decrypting the userKey ->", userDecryptionResult.Error.Error())
			response.ReturnJSONResponse(w, "", "can't decode your key", http.StatusBadRequest)
			return
		}
		fmt.Printf("the user struct is -> %s \n", oldUser.GetDecryptedStringInTheStruct())

		if Request.EmailByRequestDoNotTrust != oldUser.Email {
			fmt.Printf("the email enterd by the user in req (%s) is not same as the one in the key (%s) we are returning the error \n\n", Request.EmailByRequestDoNotTrust, oldUser.Email)
			response.ReturnJSONResponse(w, "", "you email in the request does not match the one in the key", http.StatusBadRequest)
			return
		}
		resultFromTheDB := <-resultDBChannel
		// here the email is valid, now the row is 0
		if resultFromTheDB.Error != nil {
			println("there is a error form the DB request ->", resultFromTheDB.Error.Error())
			// return json response
			return
		}
		// now before we return the new struct we also need to see if the version in the key is same as the one in the Db
		// why: cause if the user  cancells the payment to the paid teirs we need to make them the free tier
	}
}
