package commonstructs

import "youtubeAdsSkipper/paymentBackendGO/common"

type UserKey struct {
	AccountID  string `json:"account_id"`
	Email      string `json:"email"`
	UserName   string `json:"user_name"`
	IsUserPaid bool   `json:"is_user_paid"` // Assuming default is false
	//
	//
	//
	//the above feilds could be used as a foreign key and the new db row will have
	//
	UserTier            string `json:"user_tier"`               // could only be free tier | recurring | one time
	Version             int64  `json:"version"`                 // this is used to compare it to the db version and get the new version from the DB, default is 0
	CheckForKeyUpdateOn int64  `json:"check_for_key_update_on"` // this field is used when we will decrypt the key on the server and check if the time is more than the time on the server. If yes then ask user to update the key
	IDPrimaryKey        int64  `json:"id_primary_key"`          // this is the id (primary key) in the DB, would give it to the razorpay during orderID
	// creation  and will take it out during the webHook-- will save us a DB call -- or in the webHook for payment capture we cam make an
	// additional DB call that will search the order id in both of the columns  and give us the orderID but then we need to put some sort of
	// wait there as It can be updated by the client , thats why it is easier to just put it here

	encryptedUserKey string // also known as cipherText, that will be decoded into the key eg. -> =mkdkcccno/ubuinewc889nxkn==
	decryptedUserKey string // this is the string representation of the user struct 29012093-|-name-|-email-|-kskjdc...
}

// now the user message table will be
//
// for updating we will need to use the

// method will encrypt The User in the struct and will give you the Encrypted Key out, also setting the key on the private feild
func (userKey *UserKey) EncryptTheUserAndGetTheKey(result <-chan common.ErrorAndResultStruct[string]) {
	// to be implemented
}

// method will return the encrypted key form the struct if we go else "", encryptedUserKey is a private property indeed
func (userKey *UserKey) GetEncryptedKey() string {
	return userKey.encryptedUserKey
}
