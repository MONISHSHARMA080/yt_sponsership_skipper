package commonstructs

import (
	"fmt"
	commonhelperfuncs "youtubeAdsSkipper/commonHelperFuncs"
)

// we are not checking if the user is free tier, this method just copies it and calls it a day, and will error if the user in db is empty,
//
// hard code values/assumptions: the user tier will be free and the version to be 0, and gthe time to check for update on is after 1 month and 1 day
func (userKey *UserKey) InitializeTheStructForNewUser(userInDb UserInDb, primaryKeyOfTheUserReturnedFromTheDB int64) error {
	userInDb.AddUserToFreeTier()
	if !userInDb.IsUserValid() {
		return fmt.Errorf("the user struct is not valid ")
	}
	userKey.AccountID = userInDb.AccountID
	userKey.UserName = userInDb.UserName
	userKey.Email = userInDb.Email
	userKey.IsUserPaid = userInDb.IsUserPaid
	userKey.UserTier = userInDb.UserTeir
	userKey.Version = 0
	userKey.IDPrimaryKey = primaryKeyOfTheUserReturnedFromTheDB
	userKey.CheckForKeyUpdateOn = commonhelperfuncs.GetTimeToExpireTheKey()
	return nil
}

