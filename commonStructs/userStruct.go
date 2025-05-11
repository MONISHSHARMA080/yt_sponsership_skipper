package commonstructs

import (
	"database/sql"
	"fmt"
	"time"
	commonhelperfuncs "youtubeAdsSkipper/commonHelperFuncs"
	"youtubeAdsSkipper/paymentBackendGO/common"
)

type UserInDb struct {
	AccountID  string
	Email      string
	UserName   string
	IsUserPaid bool // Assuming default is false
	UserTeir   string
}

// UserSignupQueryResult holds the combined result of user upsertion and their latest message.
type SignupResult struct {
	UserID        int64
	LastUserTier  sql.NullString
	LastPaymentID sql.NullString
	LastCheckTime sql.NullInt64
	LastVersion   sql.NullInt64
}

// if the user is in the DB before then all of them should return true else if any one of them  is not there then we will set the user in the free tier
func (res *SignupResult) WasTheUserAlreadyInTheDB() bool {
	return res.LastUserTier.Valid && res.LastPaymentID.Valid && res.LastCheckTime.Valid && res.LastVersion.Valid
}

// returns true if the user is valid, the UserTeir should be set before using it
func (U *UserInDb) IsUserValid() bool {
	// return U.AccountID == "" || U.Email == "" || U.UserName == "" || U.UserTeir == ""
	return U.AccountID != "" && U.Email != "" && U.UserName != "" && U.UserTeir != ""
}

func (UserInDb *UserInDb) AddUserToFreeTier() {
	UserInDb.UserTeir = "free tier" // taken form the DB
}

// method to insert New User in Db, now the tier is not decided here
func (UserInDb *UserInDb) InsertNewUserInDb(db *sql.DB, resultChannel chan common.ErrorAndResultStruct[SignupResult]) {
	UserInDb.AddUserToFreeTier()
	var res SignupResult
	if !UserInDb.IsUserValid() {
		println("\n\n")
		println("the user in db is not valid lets see its fields tier ->", UserInDb.UserTeir, " account id ->", UserInDb.AccountID, " email ->", UserInDb.Email, " user name ->", UserInDb.UserName)
		println("is the user valid->", UserInDb.IsUserValid())
		println("\n\n")
		resultChannel <- common.ErrorAndResultStruct[SignupResult]{Result: res, Error: fmt.Errorf("the Db struct is not valid or it is empty")}
		return
	}
	resultChannelForCreatingNewUser := make(chan common.ErrorAndResultStruct[SignupResult])
	resultChannelForSeeIfWeHaveReturningUser := make(chan common.ErrorAndResultStruct[SignupResult])
	go UserInDb.dbcallToCreateNewUser(db, resultChannelForCreatingNewUser)
	go UserInDb.dbcallToSeeIfWeHaveReturningUser(db, resultChannelForSeeIfWeHaveReturningUser)
	returningUser := <-resultChannelForSeeIfWeHaveReturningUser
	newUser := <-resultChannelForCreatingNewUser
	if returningUser.Error != nil {
		println("there is a error in the returningUser DB call and it is ->", returningUser.Error.Error())
		resultChannel <- common.ErrorAndResultStruct[SignupResult]{Result: res, Error: returningUser.Error}
		return
	}
	if newUser.Error != nil {
		println("there is a error in the new user db call and it is ->", newUser.Error.Error())
		resultChannel <- common.ErrorAndResultStruct[SignupResult]{Result: res, Error: newUser.Error}
		return
	}

	res.LastCheckTime = returningUser.Result.LastCheckTime
	res.LastPaymentID = returningUser.Result.LastPaymentID
	res.LastUserTier = returningUser.Result.LastUserTier
	res.LastVersion = returningUser.Result.LastVersion
	res.UserID = newUser.Result.UserID
	fmt.Printf("the  new result user is -> %+v \n", res)
	resultChannel <- common.ErrorAndResultStruct[SignupResult]{Result: res, Error: nil}
}

// DB call to see if it is a returning user
func (UserInDb *UserInDb) dbcallToSeeIfWeHaveReturningUser(db *sql.DB, resultChannel chan common.ErrorAndResultStruct[SignupResult]) {
	// if the emil is null then how will we know
	var res SignupResult
	checkIfReturningUser := `
SELECT
  m.user_tier,
  m.razorpay_payment_id,
  m.check_for_key_update_on,
  m.version
FROM messageForTheUserAfterPaymentCaptured AS m
INNER JOIN UserAccount AS u
  ON m.user_account_id = u.id
WHERE u.email = ?
ORDER BY m.id DESC
LIMIT 1;
  `
	row := db.QueryRow(checkIfReturningUser, UserInDb.Email)
	err := row.Scan(&res.LastUserTier, &res.LastPaymentID, &res.LastCheckTime, &res.LastVersion)
	if err != nil {
		if err == sql.ErrNoRows {
			println("there is a no rows err, meanign that the user is new and does not have the column filled ")
			fmt.Printf("the result struct is -> %+v \n", res)
			resultChannel <- common.ErrorAndResultStruct[SignupResult]{Result: res, Error: nil}
			return
		}
		println("there is a error in checkIfReturningUser db query and it is ->", err.Error())
		resultChannel <- common.ErrorAndResultStruct[SignupResult]{Result: res, Error: nil}
		return
	}
	fmt.Printf("the result of check is existing user db call is -> %+v \n\n", res)
	resultChannel <- common.ErrorAndResultStruct[SignupResult]{Result: res, Error: nil}
}

// db call to create a new user (if there is one)
func (UserInDb *UserInDb) dbcallToCreateNewUser(db *sql.DB, resultChannel chan common.ErrorAndResultStruct[SignupResult]) {
	var res SignupResult
	query := `

INSERT OR REPLACE INTO UserAccount
(accountid, email, userName, is_a_paid_user)
VALUES (?, ?, ?, ?)
ON CONFLICT(email) DO UPDATE SET
accountid = excluded.accountid,
userName = excluded.userName,
is_a_paid_user = excluded.is_a_paid_user
RETURNING id

  `
	row := db.QueryRow(query, UserInDb.AccountID, UserInDb.Email, UserInDb.UserName, UserInDb.IsUserPaid)
	println("made the sql query")
	err := row.Scan(&res.UserID)
	if err != nil {
		println("there is a error in getting user form the DB and it is ->", err.Error())
		resultChannel <- common.ErrorAndResultStruct[SignupResult]{Result: res, Error: nil}
		return
	}

	fmt.Printf("the result of db call for new user is -> %+v \n\n", res)
	resultChannel <- common.ErrorAndResultStruct[SignupResult]{Result: res, Error: nil}
}

// method to insert New User in and get the encrypted key with it in the channel result string
func (UserInDb *UserInDb) InsertNewUserInDbAndGetNewKey(db *sql.DB, resultChannel chan common.ErrorAndResultStruct[string]) {
	resultChannelForDb := make(chan common.ErrorAndResultStruct[SignupResult])
	resulChannellForEncryptionKey := make(chan common.ErrorAndResultStruct[string])
	go UserInDb.InsertNewUserInDb(db, resultChannelForDb)
	resultFromTheDB := <-resultChannelForDb
	if resultFromTheDB.Error != nil {
		println("there is a error in inserting the user ")
		resultChannel <- common.ErrorAndResultStruct[string]{Result: "", Error: resultFromTheDB.Error}
		return
	}
	// here we are going to see if the user is in the free tier or
	// take the user in the Db and make the key there
	//
	//
	//---------- also check the key deprecate time here if the user is deperecated then do it and move the user to free tier
	//
	//
	//
	userInfo := resultFromTheDB.Result
	// settign the user tier (in user in DB)
	if resultFromTheDB.Result.WasTheUserAlreadyInTheDB() {
		println("the user was in the DB before ")
		// fill the already user detail in the userINDb
		if time.Now().Unix() < userInfo.LastCheckTime.Int64 {
			println("the user tier is valid and we are setting it the same")
			UserInDb.UserTeir = userInfo.LastUserTier.String
		} else {
			println("the user tier expired, adding them to the free tier")
			UserInDb.AddUserToFreeTier()
		}
	} else {
		println("the user was not in the DB beofre, we are adding them to the free tier")
		UserInDb.AddUserToFreeTier()
	}
	// if user is on free tier we need to add it here
	newUserKey := UserKey{}
	err := newUserKey.InitializeTheStructForTheUser(*UserInDb, resultFromTheDB.Result.UserID, resultFromTheDB.Result.LastVersion.Int64)
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{Result: "", Error: err}
		return
	}

	// now adding the info to the key
	if resultFromTheDB.Result.WasTheUserAlreadyInTheDB() {
		println("the user was in the DB before ")
		// fill the already user detail in the userINDb
		// the  user tier is already filled
		newUserKey.IDPrimaryKey = userInfo.UserID
		newUserKey.Version = userInfo.LastVersion.Int64
		// now if the user tier is still valid
		if time.Now().Unix() < userInfo.LastCheckTime.Int64 {
			newUserKey.CheckForKeyUpdateOn = userInfo.LastCheckTime.Int64
			newUserKey.IsUserPaid = true
		} else {
			// if the user time ran out then the user
			newUserKey.CheckForKeyUpdateOn = commonhelperfuncs.GetTimeToExpireTheKey(false)
		}
	} else {
		println("the user was not in the DB beofre, we are adding them to the free tier")
		UserInDb.AddUserToFreeTier()
	}
	fmt.Printf("the user key for the user is -> %+v \n\n", newUserKey)
	println("about to encrypt the key ")
	go newUserKey.EncryptTheUser(resulChannellForEncryptionKey)
	resulFromEncryptedKey := <-resulChannellForEncryptionKey
	if resulFromEncryptedKey.Error != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{Result: resulFromEncryptedKey.Result, Error: resulFromEncryptedKey.Error}
		return
	}
	resultChannel <- common.ErrorAndResultStruct[string]{Result: resulFromEncryptedKey.Result, Error: nil}
}
