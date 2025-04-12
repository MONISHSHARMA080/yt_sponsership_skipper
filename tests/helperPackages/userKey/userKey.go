package userkey

import (
	commonhelperfuncs "youtubeAdsSkipper/commonHelperFuncs"
	commonstructs "youtubeAdsSkipper/commonStructs"
	"youtubeAdsSkipper/paymentBackendGO/common"
	userindb "youtubeAdsSkipper/tests/helperPackages/userInDb"
)

type UserKey struct {
	User commonstructs.UserKey
}

// initialize the struct and set the value on it if the error is not there(form empty user in DB struct); note we are harcoding it to 0 as the new use it will always be one
func (uk *UserKey) InitializeTheStruct(userInDb userindb.GetUserDetailsOutInterface, primaryKey int64) error {
	userAccountID, UserName, UserEmail, userTier, isUserPaid, err := userInDb.GetUserDetailsOut()
	if err != nil {
		return err
	}
	uk.User.AccountID = userAccountID
	uk.User.UserName = UserName
	uk.User.Email = UserEmail
	uk.User.UserTier = userTier
	uk.User.IsUserPaid = isUserPaid
	uk.User.Version = 0 //

	uk.User.IDPrimaryKey = primaryKey
	uk.User.CheckForKeyUpdateOn = commonhelperfuncs.GetTimeToExpireTheKey(false)
	return nil
}

// make a func that will take the userInDB fields and will set it on the usr In Key and will give you the Encrypted key for the user
// NOTE: the assumtion is the user is new and thats why the Version is hardcoded to 0
func (uk *UserKey) InitializeTheStructAndGetEncryptedKey(userInDb userindb.GetUserDetailsOutInterface, IDPrimaryKey int64, resultChannel chan common.ErrorAndResultStruct[string]) {
	err := uk.InitializeTheStruct(userInDb, IDPrimaryKey)
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{Result: "", Error: err}
		return
	}
	go uk.User.EncryptTheUser(resultChannel)
}
