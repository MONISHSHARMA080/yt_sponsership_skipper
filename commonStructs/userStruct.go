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
	// return U.AccountID == "" || U.Email == "" || U.UserName == "" || U.UserTeir == ""
	return U.AccountID != "" && U.Email != "" && U.UserName != "" && U.UserTeir != ""
}

func (UserInDb *UserInDb) AddUserToFreeTier() {
	UserInDb.UserTeir = "free tier" // taken form the DB
}

// method to insert New User in Db
func (UserInDb *UserInDb) InsertNewUserInDb(db *sql.DB, resultChannel chan common.ErrorAndResultStruct[int64]) {
	UserInDb.AddUserToFreeTier()
	if !UserInDb.IsUserValid() {
		println("\n\n")
		println("the user in db is not valid lets see its fields tier ->", UserInDb.UserTeir, " account id ->", UserInDb.AccountID, " email ->", UserInDb.Email, " user name ->", UserInDb.UserName)
		println("is the user valid->", UserInDb.IsUserValid())
		println("\n\n")
		resultChannel <- common.ErrorAndResultStruct[int64]{Result: 0, Error: fmt.Errorf("the Db struct is not valid or it is empty")}
		return
	}
	var id int64
	query := `
        INSERT OR REPLACE INTO UserAccount
        (accountid, email, userName, is_a_paid_user)
        VALUES (?, ?, ?, ?)
        RETURNING id
    `
	err := db.QueryRow(query, UserInDb.AccountID, UserInDb.Email, UserInDb.UserName, UserInDb.IsUserPaid).Scan(&id)
	println("made the sql query")
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[int64]{Result: 0, Error: err}
		return
	}
	// the primary key is not returned here
	println("the db insert is successfull-- and the primary id is ->", id)
	resultChannel <- common.ErrorAndResultStruct[int64]{Result: id, Error: nil}
}

// method to insert New User in and get the encrypted key with it in the channel result string
func (UserInDb *UserInDb) InsertNewUserInDbAndGetNewKey(db *sql.DB, resultChannel chan common.ErrorAndResultStruct[string]) {
	resultChannelForDb := make(chan common.ErrorAndResultStruct[int64])
	resulChannellForEncryptionKey := make(chan common.ErrorAndResultStruct[string])
	go UserInDb.InsertNewUserInDb(db, resultChannelForDb)
	resultFromTheDB := <-resultChannelForDb
	if resultFromTheDB.Error != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{Result: "", Error: resultFromTheDB.Error}
		return
	}
	// take the user in the Db and make the key there
	newUserKey := UserKey{}
	err := newUserKey.InitializeTheStructForNewUser(*UserInDb, resultFromTheDB.Result)
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{Result: "", Error: err}
		return
	}
	println("about to encrypt the key ")
	go newUserKey.EncryptTheUser(resulChannellForEncryptionKey)
	resulFromEncryptedKey := <-resulChannellForEncryptionKey
	if resulFromEncryptedKey.Error != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{Result: resulFromEncryptedKey.Result, Error: resulFromEncryptedKey.Error}
		return
	}
	resultChannel <- common.ErrorAndResultStruct[string]{Result: resulFromEncryptedKey.Result, Error: nil}
}
