package userindb

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	commonstructs "youtubeAdsSkipper/commonStructs"
	"youtubeAdsSkipper/paymentBackendGO/common"
)

type Userindb struct {
	user commonstructs.UserInDb
}

type GetUserDetailsOutInterface interface {
	// this will get the user details out and if it is empty etc will return the error
	GetUserDetailsOut() (userAccountID, UserName, UserEmail, userTier string, isUserPaid bool, err error)
}

// generates spam user and  will set it on the struct,
// NOTE: the user will be in the free tier as defualt
func (userInDb *Userindb) generateSpamUser() {
	randomInt := rand.Intn(100000)

	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create a byte slice to store the result
	b := make([]byte, 40)

	// Fill the byte slice with random characters
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}

	userName := string(b[0:9]) + strconv.Itoa(randomInt)

	userInDb.user.AccountID = string(b)
	userInDb.user.Email = userName + "@gmail.com"
	userInDb.user.UserName = "TestUser:" + userName
	userInDb.user.AddUserToFreeTier()
}

// generates a spam user and saves it in the DB and on the struct, when we have saved the user in the DB we will return the primary ID of the use
func (usrk *Userindb) GenerateSpamUserAndSaveItInDB(db *sql.DB, getuserIdChann chan common.ErrorAndResultStruct[commonstructs.SignupResult]) {
	usrk.generateSpamUser()
	go usrk.user.InsertNewUserInDb(db, getuserIdChann)
}

func (u *Userindb) IsMyStructEmpty() bool {
	return u.user.Email == "" || u.user.AccountID == "" || u.user.UserName == "" || u.user.UserTeir == ""
}

func (u *Userindb) GetUserDetailsOut() (userAccountID, UserName, UserEmail, userTier string, isUserPaid bool, err error) {
	if u.IsMyStructEmpty() {
		return "", "", "", "", false, fmt.Errorf("the struct is empty, fill it first")
	}
	return u.user.AccountID, u.user.UserName, u.user.Email, u.user.UserTeir, u.user.IsUserPaid, nil
}
