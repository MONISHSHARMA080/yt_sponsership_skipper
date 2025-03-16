package paymentbackendgo

import (
	"io"
	"net/http"
	"os"
	"youtubeAdsSkipper/paymentBackendGO/common"
	structs "youtubeAdsSkipper/paymentBackendGO/structs"

	"github.com/razorpay/razorpay-go"
)

// handler function that will called by the user and give them the oder id
func CreateAndReturnOrderId(razorpayKeyID, razorpaySecretID string, envKeyAsByte []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		println("in the payment func and envKeyAsByte is ->", envKeyAsByte, "\n\n")
		var responseFromTheServer structs.ResponseToTheUser
		if r.Method != http.MethodPost {
			responseFromTheServer.ReturnTheErrorInJsonResponse(w,r,"", "incorrect method", http.StatusBadRequest,"" )
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			responseFromTheServer.ReturnTheErrorInJsonResponse(w,r,"", "something wrong with the body", http.StatusBadRequest,"" )
			return
		}
		var request structs.RequestFromClientInPaymentStruct
		channelForRes := make(chan common.ErrorAndResultStruct[string])

		err = request.ParseIntoJson(bodyBytes)
		if err != nil {
			println("error reading request body:", err.Error())
			responseFromTheServer.ReturnTheErrorInJsonResponse(w,r,"", "can't decode json send  by you", http.StatusBadRequest,"" )
			return
		}

		validated, infoHolder, err := request.ValidateAndExtractInfo(envKeyAsByte, channelForRes)

		if err != nil || !validated  {
			println("error in validating ->", err.Error())
			responseFromTheServer.ReturnTheErrorInJsonResponse(w,r,"", "can't decode key send by you or ", http.StatusBadRequest, infoHolder.PlanType )
			return
		}
		println("decrypted key ->", infoHolder.DecryptedKey)
		println("email ->", infoHolder.Email)
		println("price ->", infoHolder.Price)
		println("name ->", infoHolder.Name)
		println("IsPaidUser ->", infoHolder.IsPaidUser)
     	println("the plan type is ->", infoHolder.PlanType)
     	println("the plan type is ->", request.PlanType)
		//  now let's make the call to thr razopay
		var RazorpayOrder structs.RazorpayOrderResponse

		razorPayClient := razorpay.NewClient(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_SECRET_ID"))
		responsePtr, err := RazorpayOrder.AskRazorpayForTheOrderID(razorPayClient, infoHolder.Price)


		if err != nil {
			println("error in validating ->", err.Error())
			responseFromTheServer.ReturnTheErrorInJsonResponse(w,r,"", "trouble getting to the razorpay", http.StatusInternalServerError, infoHolder.PlanType )
			return
		}
		// as ID is the order ID
		println(" the oder id is ->", responsePtr.ID)
		responseFromTheServer.ReturnTheErrorInJsonResponse(w,r, responsePtr.ID, "success", http.StatusOK, infoHolder.PlanType )
	}
}








// func TrialReq(razorpayKeyID, razorpaySecretID string) error {
// 	client := razorpay.NewClient(razorpayKeyID, razorpaySecretID)
//
// 	data := map[string]interface{}{
// 		"amount":   50000, // Amount is in currency subunits. Default currency is INR. Hence, 50000 refers to 50000 paise
// 		"currency": "INR",
// 		"receipt":  "some_receipt_id",
// 	}
// 	body, err := client.Order.Create(data, nil)
// 	if err != nil {
// 		return err
// 	}
//
// 	println(body)
//
// 	var RazorpayOrderResponse structs.RazorpayOrderResponse
// 	err = RazorpayOrderResponse.ConvertResponseToJSON(body)
// 	if err != nil {
// 		return err
// 	}
// 	println(RazorpayOrderResponse.Amount, RazorpayOrderResponse.Attempts)
//
// 	return nil
// }
