package commonstructs

import (
	"database/sql"
	"fmt"
	"youtubeAdsSkipper/paymentBackendGO/common"
)

type UserInDb struct {
	AccountID  string
	Email      string
	UserName   string
	IsUserPaid bool // Assuming default is false
	UserTeir   string
}

// returns true if the user is valid, the UserTeir should be set before using it
func (U *UserInDb) IsUserValid() bool {
	// if len(U.AccountID) >= 1 {
	// 	return false
	// } else if len(U.Email) >= 1 {
	// 	return false
	// } else if len(U.UserName) >= 1 {
	// 	return false
	// } else if !U.IsUserPaid {
	// 	return false
	// } // will always br false
	// // else if U.UserTeir != "RecurringPayment" && U.UserTeir != "OneTime" {
	// // 	return false
	// // }
	// return true
	return U.AccountID == "" || U.Email == "" || U.UserName == "" || U.UserTeir == ""
}

func (UserInDb *UserInDb) AddUserToFreeTier() {
	UserInDb.UserTeir = "free tier" // taken form the DB
}

// we are not checking if the user is free tier, this method just copies it and calls it a day, and will error if the user in db is empty,
//
//hard code values/assumptions: the user tier will be free and the version to be 0, and gthe time to check for update on is after 1 month and 1 day
// func   InitializeTheStructForNewUser( userKey *UserKey, userInDb UserInDb , primaryKeyOfTheUserReturnedFromTheDB int)error{
// 	userInDb.AddUserToFreeTier()
// 	if !userInDb.IsUserValid(){
// 		return fmt.Errorf("the user struct is not valid ")
// 	}
// 	userKey.AccountID = userInDb.AccountID
// 	userKey.UserName = userInDb.UserName
// 	userKey.Email = userInDb.Email
// 	userKey.IsUserPaid = userInDb.IsUserPaid
// 	userKey.UserTier = userInDb.UserTeir
// 	userKey.Version = 0
// 	userKey.IDPrimaryKey = int64(primaryKeyOfTheUserReturnedFromTheDB)
// 	userKey.CheckForKeyUpdateOn = time.Now().AddDate(0,1,1).Unix()
// 	return nil
// }

// method to insert New User in Db
func (UserInDb *UserInDb) InsertNewUserInDb(db *sql.DB, resultChannel chan common.ErrorAndResultStruct[string]) {
	UserInDb.AddUserToFreeTier()
	if !UserInDb.IsUserValid() {
		resultChannel <- common.ErrorAndResultStruct[string]{Result: "", Error: fmt.Errorf("the Db struct is not valid or it is empty")}
		return
	}
	query := `
        INSERT OR IGNORE INTO UserAccount
        (accountid, email, userName, is_a_paid_user)
        VALUES (?, ?, ?, ?)
    `
	_, err := db.Exec(query, UserInDb.AccountID, UserInDb.Email, UserInDb.UserName, UserInDb.IsUserPaid)
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{Result: "", Error: err}
		return
	}
	resultChannel <- common.ErrorAndResultStruct[string]{Result: "", Error: nil}
}

// method to insert New User in and get the encrypted key with it in the channel result string
func (UserInDb *UserInDb) InsertNewUserInDbAndGetNewKey(db *sql.DB, resultChannel chan common.ErrorAndResultStruct[string]) {
	resultChannelForDb := make(chan common.ErrorAndResultStruct[string])
	go UserInDb.InsertNewUserInDb(db, resultChannelForDb)
	resultFromTheDB := <-resultChannelForDb
	if resultFromTheDB.Error != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{Result: "", Error: resultFromTheDB.Error}
		return
	}
	// take the user in the Db and make the key there
	newUserKey := UserKey{}
	newUserKey.InitializeTheStructForNewUser(*UserInDb, 23)
}

// method to insert  User in Db,
// not implemented yet
// func (UserInDb *UserInDb) InsertUserInDb(db *sql.DB, resultChannel chan common.ErrorAndResultStruct[string]) {
// }
