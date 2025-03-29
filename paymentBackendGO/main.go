package paymentbackendgo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	commonstructs "youtubeAdsSkipper/commonStructs"
	"youtubeAdsSkipper/paymentBackendGO/common"
	helperfuncs "youtubeAdsSkipper/paymentBackendGO/helperFuncs"
	structs "youtubeAdsSkipper/paymentBackendGO/structs"

	"github.com/razorpay/razorpay-go"
)

// handler function that will called by the user and give them the oder id
func CreateAndReturnOrderId(razorpayKeyID, razorpaySecretID string, envKeyAsByte []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		println("in the payment func and envKeyAsByte is ->", envKeyAsByte, "\n\n")
		var responseFromTheServer structs.ResponseToTheUser
		if r.Method != http.MethodPost {
			responseFromTheServer.ReturnTheErrorInJsonResponse(w, r, "", "", "incorrect method", http.StatusBadRequest)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			responseFromTheServer.ReturnTheErrorInJsonResponse(w, r, "", "", "something wrong with the body", http.StatusBadRequest)
			return
		}
		var request structs.RequestFromClientInPaymentStruct
		var userFromTheRequest commonstructs.UserKey
		channelForRes := make(chan common.ErrorAndResultStruct[string])

		err = request.ParseIntoJson(bodyBytes)
		if err != nil {
			println("error reading request body:", err.Error())
			responseFromTheServer.ReturnTheErrorInJsonResponse(w, r, "", "", "can't decode json send  by you", http.StatusBadRequest)
			return
		}

		go userFromTheRequest.DecryptTheKey(request.UserKey, channelForRes)
		validated, infoHolder, err := request.ValidateAndExtractInfo(envKeyAsByte, channelForRes)

		if err != nil || !validated {
			println("error in validating ->", err.Error())
			responseFromTheServer.ReturnTheErrorInJsonResponse(w, r, "", "", "can't decode key send by you or ", http.StatusBadRequest)
			return
		}
		resutFromKeyDecryption := <- channelForRes
		if resutFromKeyDecryption.Error != nil{
			println("error in key decoidng ->", err.Error())
			responseFromTheServer.ReturnTheErrorInJsonResponse(w, r, "", "", "can't decode key send by you or ", http.StatusBadRequest)
			return
		}
		
	email, name, isPaidUser := userFromTheRequest.Email, userFromTheRequest.UserName, userFromTheRequest.IsUserPaid
	infoHolder.Email = email
	infoHolder.Name = name
	infoHolder.IsPaidUser = isPaidUser

		println("email ->", infoHolder.Email)
		// println("price ->", infoHolder.Price)
		println("name ->", infoHolder.Name)
		println("IsPaidUser ->", infoHolder.IsPaidUser)

		recurringChannel := make(chan common.ErrorAndResultStruct[string])
		oneTimeChannel := make(chan common.ErrorAndResultStruct[string])
		var RazorpayOrderForRecurring structs.RazorpayOrderResponse
		var RazorpayOrderForOneTime structs.RazorpayOrderResponse
		

		fmt.Printf("\n\n----------the useKey struct's decrypted key is -> %s ---------\n\n", userFromTheRequest.GetDecryptedStringInTheStruct())
		fmt.Printf("\n\n----------the useKey struct is -> %v ---------\n\n", userFromTheRequest)
		println(" it should be id primary key ->", userFromTheRequest.IDPrimaryKey)

		razorPayClient := razorpay.NewClient(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_SECRET_ID"))

		go RazorpayOrderForRecurring.AskRazorpayForTheOrderID(razorPayClient, infoHolder.PriceForRecurring, recurringChannel, userFromTheRequest.IDPrimaryKey)
		go RazorpayOrderForOneTime.AskRazorpayForTheOrderID(razorPayClient, infoHolder.PriceForOneTime, oneTimeChannel, userFromTheRequest.IDPrimaryKey)
		println("waiting for the func to finish")
		resFromOneTime := <-oneTimeChannel
		resFromRecurring := <-recurringChannel
		if resFromOneTime.Error != nil || resFromRecurring.Error != nil {
			println("error in validating ->")
			responseFromTheServer.ReturnTheErrorInJsonResponse(w, r, "", "", "trouble getting to the razorpay", http.StatusInternalServerError)
			return
		}
		// if err != nil {
		// 	println("error in validating ->", err.Error())
		// 	responseFromTheServer.ReturnTheErrorInJsonResponse(w, r, "", "", "trouble getting to the razorpay", http.StatusInternalServerError)
		// 	return
		// }
		println(" the oder id is (one time)->", RazorpayOrderForOneTime.ID, "  -- recurring ->", RazorpayOrderForRecurring.ID)

		// getting to the DB to store shit

		db := helperfuncs.DbConnect()
		resChannel := make(chan common.ErrorAndResultStruct[string])
		dbFieldTOVerifyPayment := structs.CreateDBFieldForStoringTempOrderId(RazorpayOrderForRecurring.ID, RazorpayOrderForOneTime.ID)
		go dbFieldTOVerifyPayment.SetTokensForTheUser(db, infoHolder.Email, resChannel)
		resultFromAddingTODb := <-resChannel
		if resultFromAddingTODb.Error != nil {
			println("there is a error in adding the tokens to the db ->", resultFromAddingTODb.Error.Error())
			responseFromTheServer.ReturnTheErrorInJsonResponse(w, r, "", "", "Internal server error, can't reach DB", http.StatusInternalServerError)
		} else {
			println("the result of adding the token to db is ->", resultFromAddingTODb.Result)
		}
		responseFromTheServer.ReturnTheErrorInJsonResponse(w, r, RazorpayOrderForRecurring.ID, RazorpayOrderForOneTime.ID, "success", http.StatusOK)
		timeTaken := time.Since(startTime)
		println("time taken is ->", timeTaken.Microseconds(), " Microseconds or ", timeTaken.Seconds(), " sec", " and", timeTaken.Milliseconds(), " ms")
	}
}
