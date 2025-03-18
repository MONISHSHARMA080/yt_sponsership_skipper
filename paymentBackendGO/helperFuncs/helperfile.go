package helperfuncs

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"youtubeAdsSkipper/paymentBackendGO/common"
)
func ExtractPriceFormEnv (price  string) (int64, error){
if price == "" {
			return 0, fmt.Errorf("the price name for the one time payment price is not there in the env")
		}
		priceInInt, err := strconv.ParseInt(price, 10, 64)
		if err != nil {
			return 0, err
		}
		return priceInInt, nil
}

// paymentPlanType should be "onetime" or "recurring"
func GetPaymentForThePlan(paymentPlanType string) (int64, error) {
	println("the payment type we got was ->", paymentPlanType)
	if paymentPlanType == "onetime" {
		price := os.Getenv("ONETIMEPAYMENTPRICE")
		if price == "" {
			return 0, fmt.Errorf("the price name for the one time payment price is not there in the env")
		}
		intVal, err := strconv.ParseInt(price, 10, 64)
		if err != nil {
			return 0, err
		}
		println("price will be ", intVal," for the payment plan ->", paymentPlanType)
		return intVal, nil
	} else if paymentPlanType == "recurringpayment" {
		price := os.Getenv("RECURRINGPAYMENTPRICE")
		if price == "" {
			return 0, fmt.Errorf("the price for the one time payment price is not there in the env")
		}
		intVal, err := strconv.ParseInt(price, 10, 64)
		if err != nil {
			return 0, err
		}
		return intVal, nil
	}
	return 0, fmt.Errorf("the payment type could only be of 2 types")
}

func DecryptAndWriteToChannel(ciphertextAsString string, EnvKey []byte, channErr chan<- common.ErrorAndResultStruct[string]) {
	// First, decode the base64 encoded string
	println("in the decrypt_and_write_to_channel ")
	ciphertextAsByte, err := base64.StdEncoding.DecodeString(ciphertextAsString)
	if err != nil {
		channErr <- common.ErrorAndResultStruct[string]{Error: fmt.Errorf("failed to decode base64: %v", err), Result: ""}
		return
	}
	// println("ciphertext as text is ->", ciphertextAsString)
	// println("decoded ciphertext length:", len(ciphertextAsByte))

	// Now decrypt the actual ciphertext
	stringAsByte, err := decrypt(ciphertextAsByte, EnvKey)
	if err != nil {
		channErr <- common.ErrorAndResultStruct[string]{Error: fmt.Errorf("failed to decrypt: %v", err), Result: ""}
		return
	}

	string_as_string := string(stringAsByte)
	channErr <- common.ErrorAndResultStruct[string]{Error: nil, Result: string_as_string}
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

func GetEmailAndNameFormKey(k string) (email, name string, isPaidUsers bool, err error) {
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
