package commonstructs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	// "time"
	// commonstructs "youtubeAdsSkipper/commonStructs"
	"youtubeAdsSkipper/paymentBackendGO/common"
)

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
	Version             int64  `json:"version"`                 // this is used to compare it to the db version and get the new version from the DB, default is 1
	CheckForKeyUpdateOn int64  `json:"check_for_key_update_on"` // this field is used when we will decrypt the key on the server and check if the time is more than the time on the server. If yes then ask user to update the key
	IDPrimaryKey        int64  `json:"id_primary_key"`          // this is the id (primary key) in the DB, would give it to the razorpay during orderID
	// creation  and will take it out during the webHook-- will save us a DB call -- or in the webHook for payment capture we cam make an
	// additional DB call that will search the order id in both of the columns  and give us the orderID but then we need to put some sort of
	// wait there as It can be updated by the client , thats why it is easier to just put it here

	encryptedUserKey string // also known as cipherText, that will be decoded into the key eg. -> =mkdkcccno/ubuinewc889nxkn==
	decryptedUserKey string // this is the string representation of the user struct 29012093-|-name-|-email-|-kskjdc...
}

// should return false if the struct is initialized jsut now or things smth is not set
// A better approach: if we int all the && feilds and add them and it is equal to 0 then may be there is a error
func (u *UserKey) IsMyStructEmpty() bool {
	return u.AccountID == "" || u.Email == "" || u.UserName == "" || u.UserTier == "" && !u.IsUserPaid &&
		u.Version == 0 && u.CheckForKeyUpdateOn == 0 && u.IDPrimaryKey == 0
}

// method will return the encrypted key form the struct if we go else "", encryptedUserKey is a private property indeed
func (userKey *UserKey) GetEncryptedKey() string {
	return userKey.encryptedUserKey
}

// returns the decrypted json string stored in the struct; note it does not decode the encrypted key and returns the json output
// this is just to get the value form the struct
func (userKey *UserKey) GetDecryptedStringInTheStruct() string {
	return userKey.decryptedUserKey
}

// decrypts the encryptedUserKey on the struct  (probally set it first using the method); and sets it on the struct and returns it too
//
// working: will take the encrypted key and get the decoded string out(JSON) and then try to decode the JSON string into the struct
func (userKey *UserKey) DecryptTheKey(encryptedUserKeyToDecrypt string, resultChannel chan common.ErrorAndResultStruct[string]) {
	startTime := time.Now()
	userKey.encryptedUserKey = encryptedUserKeyToDecrypt

	if userKey.encryptedUserKey == "" {
		resultChannel <- common.ErrorAndResultStruct[string]{Result: "", Error: fmt.Errorf("the encrypted key is empty and we can't decode it")}
		return
	}

	// Get encryption key from environment variable
	encryptionKey := os.Getenv("encryption_key")
	if encryptionKey == "" {
		resultChannel <- common.ErrorAndResultStruct[string]{
			Result: "",
			Error:  fmt.Errorf("encryption key not found in environment variables"),
		}
		return
	}

	// Decode the base64 encrypted key
	ciphertext, err := base64.StdEncoding.DecodeString(userKey.encryptedUserKey)
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{
			Result: "",
			Error:  fmt.Errorf("failed to decode base64 encrypted key: %v", err),
		}
		return
	}

	// Prepare encryption key (ensure 32 bytes for AES-256)
	keyBytes := []byte(encryptionKey)
	if len(keyBytes) > 32 {
		keyBytes = keyBytes[:32]
	} else if len(keyBytes) < 32 {
		// Pad the key if it's too short
		keyBytes = append(keyBytes, make([]byte, 32-len(keyBytes))...)
	}

	// Create AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{
			Result: "",
			Error:  fmt.Errorf("failed to create cipher block: %v", err),
		}
		return
	}

	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{
			Result: "",
			Error:  fmt.Errorf("failed to create GCM: %v", err),
		}
		return
	}

	// Extract nonce (first 12 bytes)
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		resultChannel <- common.ErrorAndResultStruct[string]{
			Result: "",
			Error:  fmt.Errorf("ciphertext too short"),
		}
		return
	}
	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data
	decryptedData, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{
			Result: "",
			Error:  fmt.Errorf("decryption failed: %v", err),
		}
		return
	}

	// Store decrypted JSON string in the struct
	userKey.decryptedUserKey = string(decryptedData)

	err = json.Unmarshal(decryptedData, userKey)
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{Result: "", Error: err}
		return
	}

	println("time taken in userKey decoding is ->", time.Since(startTime).Microseconds(), " Microseconds")
	// Send result through the channel
	resultChannel <- common.ErrorAndResultStruct[string]{
		Result: userKey.decryptedUserKey,
		Error:  nil,
	}
}

// this func encrypts the user struct into a key and returns and sets encrypted key in the channel and on the struct
//
// working: it will take the struct and json encode it, and take that json string and encypt it
func (userKey *UserKey) EncryptTheUser(resultChannel chan common.ErrorAndResultStruct[string]) {
	// check if the struct is empty, if it is return error
	if userKey.IsMyStructEmpty() {
		resultChannel <- common.ErrorAndResultStruct[string]{
			Result: "",
			Error:  fmt.Errorf("the struct is empty or not initialized"),
		}
		return
	}

	encryptionKey := os.Getenv("encryption_key")
	if encryptionKey == "" {
		resultChannel <- common.ErrorAndResultStruct[string]{
			Result: "",
			Error:  fmt.Errorf("encryption key not found in environment variables"),
		}
		return
	}

	// Convert user struct to JSON
	jsonData, err := json.Marshal(userKey)
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{
			Result: "",
			Error:  fmt.Errorf("failed to marshal user data: %v", err),
		}
		return
	}

	// Prepare encryption key (ensure 32 bytes for AES-256)
	keyBytes := []byte(encryptionKey)
	if len(keyBytes) > 32 {
		keyBytes = keyBytes[:32]
	} else if len(keyBytes) < 32 {
		// Pad the key if it's too short
		keyBytes = append(keyBytes, make([]byte, 32-len(keyBytes))...)
	}

	// Create AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{
			Result: "",
			Error:  fmt.Errorf("failed to create cipher block: %v", err),
		}
		return
	}

	println("creating the GCM mode")
	// Create GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{
			Result: "",
			Error:  fmt.Errorf("failed to create GCM: %v", err),
		}
		return
	}

	// Generate a secure random nonce
	nonce, err := generateSecureNonce(gcm.NonceSize())
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[string]{
			Result: "",
			Error:  fmt.Errorf("failed to generate nonce: %v", err),
		}
		return
	}

	// Encrypt the data
	encryptedData := gcm.Seal(nonce, nonce, jsonData, nil)

	// Base64 encode the encrypted data
	encryptedBase64 := base64.StdEncoding.EncodeToString(encryptedData)

	// Store the encrypted key in the struct
	userKey.encryptedUserKey = encryptedBase64

	println("about to write to the result channel in the encryptedUserKey ->", encryptedBase64)
	// Send result through the channel
	resultChannel <- common.ErrorAndResultStruct[string]{
		Result: encryptedBase64,
		Error:  nil,
	}
}

// returns true if the user should update the key, compares the value to the current time
// will not panic if the struct is  not initialized
func (userKey *UserKey) ShouldWeTellUserToGoGetANewKey() (bool, error) {
	// note if the user is in free tier return false as in DBfeild for it we are returning 0, if user in free tier
	// we won't tell them to go get a new key ever and if not then we might, is the reason
	if userKey.IsMyStructEmpty() {
		return false, fmt.Errorf("the struct is not initialized")
	}
	if userKey.UserTier == "free tier" {
		return false, nil
	}
	return time.Now().Unix() >= userKey.CheckForKeyUpdateOn, nil
}

// returns true if the user should update the key, compares the value to the current time
// will panic if the struct is  not initialized
func (userKey *UserKey) ShouldWeTellUserToGoGetANewKeyPanic() bool {
	// note if the user is in free tier return false as in DBfeild for it we are returning 0, if user in free tier
	// we won't tell them to go get a new key ever and if not then we might, is the reason
	println("----++-- in the func to see if we should tell user to update the key")
	if userKey.IsMyStructEmpty() {
		panic("the struct is not initialized (userKey )")
	}
	if userKey.UserTier == "free tier" {
		println("the user is on free tier")
		return false
	}
	fmt.Printf("Time remaining until key update: %f\n", time.Until(time.Unix(userKey.CheckForKeyUpdateOn, 0)).Seconds())
	return time.Now().Unix() >= userKey.CheckForKeyUpdateOn
}

//
//
//
// ------ helper functions --------
//
//
//

// cryptoRandRead generates cryptographically secure random bytes
func cryptoRandRead(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, fmt.Errorf("cannot read into empty byte slice")
	}

	// Use crypto/rand for cryptographically secure random number generation
	n, err = rand.Read(b)
	if err != nil {
		return 0, fmt.Errorf("failed to generate secure random bytes: %v", err)
	}

	// Verify that all bytes were filled
	if n != len(b) {
		return n, fmt.Errorf("incomplete random byte generation: expected %d, got %d", len(b), n)
	}

	return n, nil
}

// generateSecureNonce creates a cryptographically secure nonce of specified size
func generateSecureNonce(nonceSize int) ([]byte, error) {
	nonce := make([]byte, nonceSize)
	_, err := cryptoRandRead(nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to generate secure nonce: %v", err)
	}
	return nonce, nil
}
