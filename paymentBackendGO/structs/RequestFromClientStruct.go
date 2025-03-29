package structs

import (
	"encoding/json"
	"os"
	"youtubeAdsSkipper/paymentBackendGO/common"
	helperfuncs "youtubeAdsSkipper/paymentBackendGO/helperFuncs"
)

type RequestFromClientInPaymentStruct struct {
	UserKey string `json:"user_key"`
}

func (requestFromClient *RequestFromClientInPaymentStruct) ParseIntoJson(data []byte) error {
	err := json.Unmarshal(data, requestFromClient)
	if err != nil {
		return err
	}
	return nil
}

func (req *RequestFromClientInPaymentStruct) ValidateAndExtractInfo(envKey []byte, channelForRes chan common.ErrorAndResultStruct[string], ) (bool, *InfoHolder, error) {
	// go  helperfuncs.Decrypt_and_write_to_channel(req.UserKey, EnvKey byte, envenvKey , chan<- structs.ErrorAndResultStruct[string])(request.Key, os_env_key, channel_for_userDetails)

	// channelToDecryptUserKey := make(chan common.ErrorAndResultStruct[string])

	// go helperfuncs.DecryptAndWriteToChannel(req.UserKey, envKey, channelForRes)
	// this will also validate the plan type
	priceForRecurring, err := helperfuncs.ExtractPriceFormEnv(os.Getenv("RECURRINGPAYMENTPRICE"))
	if err != nil {
		return false, nil, err
	}
	priceForOneTime, err := helperfuncs.ExtractPriceFormEnv(os.Getenv("ONETIMEPAYMENTPRICE"))
	if err != nil {
		return false, nil, err
	}

	// price, err :=   helperfuncs.GetPaymentForThePlan(req.PlanType)
	// if err != nil{
	//   return false, nil , err
	// }

	// println("the price is ", price, "  for the plan type ", req.PlanType)
	// now decrypting the struct
	var InfoHolder InfoHolder
	InfoHolder.PriceForOneTime = priceForOneTime
	InfoHolder.PriceForRecurring = priceForRecurring
	// var userInDB  common

	// decryptedKey := <-channelForRes
	// if decryptedKey.Error != nil {
	// 	return false, nil, decryptedKey.Error
	// }

	// InfoHolder.DecryptedKey = decryptedKey.Result
	// //  InfoHolder.PlanType = req.PlanType
	// //  println("the plan type is ->", InfoHolder.PlanType)
	// println("decrypted key is ->", decryptedKey.Result, " and same in the infoHolder is ->", InfoHolder.DecryptedKey)


	return true, &InfoHolder, nil
}
