package routehandlerfunc

import (
	"fmt"
	"net/http"
	"time"
	routehandlerfunc "youtubeAdsSkipper/RouteHandlerFunc/GetNewKey/HelperFunc"
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
		timeNow := time.Now().Unix()

		resultFromTheDB := <-resultDBChannel
		// here the email is valid, now the row is 0
		if resultFromTheDB.Error != nil {
			println("there is a error form the DB request ->", resultFromTheDB.Error.Error())
			response.ReturnJSONResponse(w, "", "something went wrong on our side", http.StatusInternalServerError)
			return
		}
		println("the user's tier in the key provided by them is ->", oldUser.UserTier)

		// meaning if the time to update the key is not now aka key(user's tier) is vaild
		// if the current time is less than the time to upgrade on  then do not give them a new key back
		if timeNow < oldUser.CheckForKeyUpdateOn {
			println("the user tier is valid and time has not came yet to update the key")
			response.ReturnJSONResponse(w, "", "you don't need to update the key as you tier is still valid", http.StatusBadRequest)
			return
		}
		//  this means the message is not there in the DB
		if !resultFromTheDB.Result {
			// here change the key of the user and
			println("the resultFromTheDB is false with no errr for the ", Request.EmailByRequestDoNotTrust, " meaning there is no message")
			println("here we would give the user key and adjust it for the future time")
			// set time on the key to be a month in the future, don't do DBStruct.UserTier as when the db does not return the field will be empty
			resultForNewUser := routehandlerfunc.UpdateTheCheckForKeyUpdateToNewValue(&DBStruct, &oldUser)
			if resultForNewUser.Error != nil {
				println("error in gettting time to CheckForKeyUpdateOn ->", resultForNewUser.Error.Error())
				response.ReturnJSONResponse(w, "", "something went wrong on our side in giving you your new key", http.StatusInternalServerError)
				return
			}
			response.ReturnJSONResponse(w, resultForNewUser.Result, "success", http.StatusOK)
			return
		}

		println("actual case")
		// now before we return the new struct we also need to see if the version in the key is same as the one in the Db
		// why: cause if the user  cancells the payment to the paid teirs we need to make them the free tier

		println("time to update the user form the key is on ->", oldUser.CheckForKeyUpdateOn, " and the time on the server is ->", timeNow, " is the timeNow >= time to update the user on ->", timeNow >= oldUser.CheckForKeyUpdateOn)
		fmt.Printf("\n the version in DB struct is %d and the one in the user struct is %d \n\n ", DBStruct.Version, oldUser.Version)
		// here  meaning that the message for the user is not there, or user has not made a new payment
		// so lets downgrade the user to the free tier, here also check if the time has expired for the user in that tier
		// this is just a safeguard, meanign  if the version is same && the time is >= time to update user on, checked earlier ->&& oldUser.CheckForKeyUpdateOn >= timeNow
		if DBStruct.Version == oldUser.Version {
			println("the user has downgrade to the new tier as the version form the DB is same as the version from the key")
			resultForNewUser := routehandlerfunc.DownGradeTheUserToFreeTierAndAlsoSetTheTimeAfterAMonth(&DBStruct, &oldUser)
			if resultForNewUser.Error != nil {
				println("error in gettting time to CheckForKeyUpdateOn ->", resultForNewUser.Error.Error())
				response.ReturnJSONResponse(w, "", "something went wrong on our side in giving you your new key", http.StatusInternalServerError)
				return
			}
			response.ReturnJSONResponse(w, resultForNewUser.Result, "success", http.StatusOK)
			return

		} else {
			// there is a new message for the user
			println("updating you to new message")

			resultForNewUser := routehandlerfunc.UpdateTheUserToNewMessage(&DBStruct, &oldUser)
			if resultForNewUser.Error != nil {
				println("error in gettting time to CheckForKeyUpdateOn ->", resultForNewUser.Error.Error())
				response.ReturnJSONResponse(w, "", "something went wrong on our side in giving you your new key", http.StatusInternalServerError)
				return
			}
			response.ReturnJSONResponse(w, resultForNewUser.Result, "success", http.StatusOK)
			return
		}
	}
}
