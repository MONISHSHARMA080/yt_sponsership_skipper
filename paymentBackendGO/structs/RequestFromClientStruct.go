package structs

import (
	"encoding/json"
	"youtubeAdsSkipper/paymentBackendGO/common"
	helperfuncs "youtubeAdsSkipper/paymentBackendGO/helperFuncs"
)

type RequestFromClientInPaymentStruct struct{
    PlanType    string  `json:"plan_type"` // Can be "recurring" or "onetime"
    UserKey     string    `json:"user_key"`
}

func (requestFromClient *RequestFromClientInPaymentStruct) ParseIntoJson(data []byte)error{
    err:=  json.Unmarshal(data, requestFromClient)
    if err != nil {
      return err
    }
    return nil
}

func (req *RequestFromClientInPaymentStruct) ValidateAndExtractInfo(envKey []byte, channelForRes chan common.ErrorAndResultStruct[string]) (bool, *InfoHolder ,error){
		// go  helperfuncs.Decrypt_and_write_to_channel(req.UserKey, EnvKey byte, envenvKey , chan<- structs.ErrorAndResultStruct[string])(request.Key, os_env_key, channel_for_userDetails)

    go helperfuncs.DecryptAndWriteToChannel( req.UserKey, envKey, channelForRes )
    // this will also validate the plan type
    price, err :=   helperfuncs.GetPaymentForThePlan(req.PlanType)
    if err != nil{
      return false, nil , err
    }
    println("the price is ", price, "  for the plan type ", req.PlanType)
    // now decrypting the struct
    var InfoHolder InfoHolder
    InfoHolder.Price = price
    // var userInDB  common
     decryptedKey := <- channelForRes
     if decryptedKey.Error != nil{
       return false, nil, decryptedKey.Error
     }
     InfoHolder.DecryptedKey =  decryptedKey.Result
     InfoHolder.PlanType = req.PlanType
     println("the plan type is ->", InfoHolder.PlanType)
     println("decrypted key is ->", decryptedKey.Result, " and same in the infoHolder is ->", InfoHolder.DecryptedKey)
    email, name, isPaidUser, err :=  helperfuncs.GetEmailAndNameFormKey(InfoHolder.DecryptedKey)
    if err != nil {
      return false, nil, err
    }
    InfoHolder.Email = email
    InfoHolder.Name = name
    InfoHolder.IsPaidUser = isPaidUser

    return true, &InfoHolder, nil
	  // resultFormChannel := <-channel_for_userDetails
	  // if resultFormChannel.err != nil {
		// 	println("the error in decoding the key is ->", resultFormChannel.err.Error())
		// 	returnTheJsonResponseonError("key value can't be empty", http.StatusBadRequest, false, w)
		// 	return
		// }
		// email, name := getEmailAndNameFormKey(resultFormChannel.string_value)
}
// func (req *RequestFromClientInPaymentStruct) ValidateUserPlanType(planType string) bool{

//   return false
// }




