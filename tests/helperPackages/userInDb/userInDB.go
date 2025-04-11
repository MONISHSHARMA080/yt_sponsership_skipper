package userindb

import (
	"database/sql"
	"math/rand"
	"strconv"
	"time"
	commonstructs "youtubeAdsSkipper/commonStructs"
	"youtubeAdsSkipper/paymentBackendGO/common"
)

type Userindb struct {
	user commonstructs.UserInDb
}

// generates spam user and  will set it on the struct,
// NOTE: the user will be in the free tier as defualt
func (userInDb *Userindb) generateSpamUser() {
	randomInt := rand.Intn(100)

	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Create a byte slice to store the result
	b := make([]byte, 88)

	// Fill the byte slice with random characters
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}

	userName := string(b[0:10]) + strconv.Itoa(randomInt)

	userInDb.user.AccountID = string(b)
	userInDb.user.Email = userName + "@gmail.com"
	userInDb.user.UserName = userName
	userInDb.user.AddUserToFreeTier()
}

// generates a spam user and saves it in the DB and on the struct, when we have saved the user in the DB we will return the primary ID of the use
func (usrk *Userindb) GenerateSpamUserAndSaveItInDB(db *sql.DB, getuserIdChann chan common.ErrorAndResultStruct[int64]) (int64, error) {
	usrk.generateSpamUser()
	go usrk.user.InsertNewUserInDb(db, getuserIdChann)
	userIdAndErrChann := <-getuserIdChann
	if userIdAndErrChann.Error != nil {
		return 0, userIdAndErrChann.Error
	}
	return userIdAndErrChann.Result, nil
}

func GenerateDB() {
}
