package common

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

type UserKey struct {
	EncryptedUserKey string // also known as cipherText, that will be decoded into the key eg. -> =mkdkcccno/ubuinewc889nxkn==
	DecryptedUserKey string // this is the string representation of the user struct 29012093-|-name-|-email-|-kskjdc...
	UserInTheDb      UserInDB
}

func (key *UserKey) Encrypt() error {
	return fmt.Errorf("not implemented")
}

// this function demands that the  EncryptedUserKey is set in the struct
// it will convert the encrypted key into decrypted and also set the value in the struct
func (userKey *UserKey) DecryptKey(EnvKey []byte, channErr chan<- ErrorAndResultStruct[string]) {
	// First, decode the base64 encoded string
	if userKey.EncryptedUserKey == "" {
	}
	ciphertextAsByte, err := base64.StdEncoding.DecodeString(userKey.EncryptedUserKey)
	if err != nil {
		channErr <- ErrorAndResultStruct[string]{Error: fmt.Errorf("failed to decode base64: %v", err), Result: ""}
		return
	}
	// println("ciphertext as text is ->", ciphertextAsString)
	// println("decoded ciphertext length:", len(ciphertextAsByte))

	// Now decrypt the actual ciphertext
	stringAsByte, err := decrypt(ciphertextAsByte, EnvKey)
	if err != nil {
		channErr <- ErrorAndResultStruct[string]{Error: fmt.Errorf("failed to decrypt: %v", err), Result: ""}
		return
	}

	string_as_string := string(stringAsByte)
	userKey.DecryptedUserKey = string_as_string
	channErr <- ErrorAndResultStruct[string]{Error: nil, Result: string_as_string}
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	println("decoded ciphertext length:", len(ciphertext), " and it is ", string(ciphertext))

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

func (userKey *UserKey) GetEmailAndNameFormKey() (email, name string, isPaidUsers bool, err error) {
	k := userKey.DecryptedUserKey
	if k == "" {
		return "", "", false, fmt.Errorf("the DecryptedUserKey is empty ")
	}

	strings := strings.Split(k, "-|-")
	println("length of the decrypted key is ", len(strings))
	if len(strings) < 4 {
		return "", "", false, fmt.Errorf("the array is not of right length")
	}
	var isPaidUser bool
	isPaidUser, err = strconv.ParseBool(strings[3])
	if err != nil {
		return "", "", false, fmt.Errorf("can't parse bool in 3rd position of the array ")
	}

	return strings[1], strings[2], isPaidUser, nil
}

// this function is used to make the user IN db struct attached here get ready to be used by taking the value out of the
// DecryptedUserKey and put it in the user struct
func (usr *UserKey) SetUserDetail() error {
	//
	if usr.DecryptedUserKey == "" {
		return fmt.Errorf("the decrypted key is not there")
	}
	strings := strings.Split(usr.DecryptedUserKey, "-|-")

	email, name, isPaidUser, err := usr.GetEmailAndNameFormKey()
	if err != nil {
		return err
	}

	// no array len checking as already done in the above func

	// Parse accountID from string
	accountID, err := strconv.ParseInt(strings[0], 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse accountID: %v", err)
	}
	usr.UserInTheDb = UserInDB{AccountID: accountID, Email: email, UserName: name, IsPaidUser: isPaidUser}

	return nil
}
