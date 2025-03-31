package routehandlerfunc

import (
	commonhelperfuncs "youtubeAdsSkipper/commonHelperFuncs"
	commonstructs "youtubeAdsSkipper/commonStructs"
	"youtubeAdsSkipper/paymentBackendGO/common"
	"youtubeAdsSkipper/paymentBackendGO/structs"
)

// the new value is same regardless the tier and is form the commonHelperfunc(env or hardcoded)
func UpdateTheCheckForKeyUpdateToNewValue(DBStruct *structs.MessageForUserOnPaymentCapture, oldUser *commonstructs.UserKey) common.ErrorAndResultStruct[string] {
	// newTimeToCheckForUpdateOn, err := DBStruct.GetTimeToCheckForKeyUpdateOn(oldUser.UserTier)
	// if err != nil {
	// 	println("error in gettting time to CheckForKeyUpdateOn ->", err.Error())
	// 	// response.ReturnJSONResponse(w, "", "something went wrong on our side in giving you your new key", http.StatusInternalServerError)
	// 	return common.ErrorAndResultStruct[string]{Result: "", Error: err}
	// }
	oldUser.CheckForKeyUpdateOn = commonhelperfuncs.GetTimeToExpireTheKey()

	resultDBChannelForNewuser := make(chan common.ErrorAndResultStruct[string])
	go oldUser.EncryptTheUser(resultDBChannelForNewuser)
	resultForNewuser := <-resultDBChannelForNewuser
	return resultForNewuser
}

func DownGradeTheUserToFreeTierAndAlsoSetTheTimeAfterAMonth(DBStruct *structs.MessageForUserOnPaymentCapture, oldUser *commonstructs.UserKey) common.ErrorAndResultStruct[string] {
	println("the UserTier is -> ", oldUser.UserTier)
	println("asseritng the UserTier is not free over here ->", oldUser.UserTier != "free tier")
	oldUser.UserTier = "free tier"
	return UpdateTheCheckForKeyUpdateToNewValue(DBStruct, oldUser)
}

func UpdateTheUserToNewMessage(DBStruct *structs.MessageForUserOnPaymentCapture, oldUser *commonstructs.UserKey) common.ErrorAndResultStruct[string] {
	oldUser.UserTier = DBStruct.UserTier
	oldUser.Version = DBStruct.Version
	return UpdateTheCheckForKeyUpdateToNewValue(DBStruct, oldUser)
}
