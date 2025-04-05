package routehandlerfunc

import (
	"time"
	commonhelperfuncs "youtubeAdsSkipper/commonHelperFuncs"
	commonstructs "youtubeAdsSkipper/commonStructs"
	"youtubeAdsSkipper/paymentBackendGO/common"
	"youtubeAdsSkipper/paymentBackendGO/structs"
)

// the new value is same regardless the tier and is form the commonHelperfunc(env or hardcoded)
func UpdateTheCheckForKeyUpdateToNewValue(DBStruct *structs.MessageForUserOnPaymentCapture, oldUser *commonstructs.UserKey) common.ErrorAndResultStruct[string] {
	// false cause we are hardcoding
	println("the old user key is ->", oldUser.CheckForKeyUpdateOn)
	oldTimeToCheck := oldUser.CheckForKeyUpdateOn
	oldUser.CheckForKeyUpdateOn = commonhelperfuncs.GetTimeToExpireTheKey(false)
	println("the old new key is ->", oldUser.CheckForKeyUpdateOn)
	println("adding time to new key and it is after :=>", time.Until(time.Unix(oldTimeToCheck, 0)).Round(time.Minute), " min")
	resultDBChannelForNewuser := make(chan common.ErrorAndResultStruct[string])
	go oldUser.EncryptTheUser(resultDBChannelForNewuser)
	resultForNewuser := <-resultDBChannelForNewuser
	return resultForNewuser
}

func DownGradeTheUserToFreeTierAndAlsoSetTheTimeAfterAMonth(DBStruct *structs.MessageForUserOnPaymentCapture, oldUser *commonstructs.UserKey) common.ErrorAndResultStruct[string] {
	println("the old User's Tier was -> ", oldUser.UserTier)
	println("asseritng the UserTier is not free over here ->", oldUser.UserTier != "free tier")
	oldUser.UserTier = "free tier"
	return UpdateTheCheckForKeyUpdateToNewValue(DBStruct, oldUser)
}

func UpdateTheUserToNewMessage(DBStruct *structs.MessageForUserOnPaymentCapture, oldUser *commonstructs.UserKey) common.ErrorAndResultStruct[string] {
	oldUser.UserTier = DBStruct.UserTier
	oldUser.Version = DBStruct.Version
	return UpdateTheCheckForKeyUpdateToNewValue(DBStruct, oldUser)
}
