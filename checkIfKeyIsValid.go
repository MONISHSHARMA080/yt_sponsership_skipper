package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	commonstructs "youtubeAdsSkipper/commonStructs"
	"youtubeAdsSkipper/paymentBackendGO/common"
)

type responseStruct struct {
	Message      string `json:"message"`
	Status_code  int64  `json:"status_code"`
	Success      bool   `json:"success"`
	EncryptedKey string `json:"encrypted_key"`
	Email        string `json:"email"`
	Name         string `json:"name"`
}
type requestStruct struct {
	Key string `json:"key"`
}

func returnTheJsonResponseonError(message string, statusCode int64, success bool, w http.ResponseWriter) {
	http.Error(w, message, int(statusCode))
	response := responseStruct{Message: message, Status_code: statusCode, Success: success}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		println("error in encoding the json in checkIFKeyISValid and it is ->" + err.Error())
	}
}

func CheckIfKeyIsValid(os_env_key []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			returnTheJsonResponseonError("Invalid request method", http.StatusBadRequest, false, w)
			return
		}
		var request requestStruct
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			message, statusCode := whatCausedJSONDecodeError(err)
			returnTheJsonResponseonError(message, statusCode, false, w)
			return
		}
		if request.Key == "" {
			returnTheJsonResponseonError("key value can't be empty", http.StatusBadRequest, false, w)
			return
		}
		channelToDecryptUserKey := make(chan common.ErrorAndResultStruct[string])
		userFormKey := commonstructs.UserKey{}
		go userFormKey.DecryptTheKey(request.Key, channelToDecryptUserKey)
		// channel_for_userDetails := make(chan string_and_error_channel)
		// go decrypt_and_write_to_channel(request.Key, os_env_key, channel_for_userDetails)
		resultFormChannel := <-channelToDecryptUserKey
		if resultFormChannel.Error != nil {
			println("the error in decoding the key is ->", resultFormChannel.Error.Error())
			returnTheJsonResponseonError("key value can't be empty", http.StatusBadRequest, false, w)
			return
		}
		if userFormKey.ShouldWeTellUserToGoGetANewKeyPanic() {
			println("\n\n ==the user should be upgraded as it's time ran out ===\n\n ")
			method_to_write_http_and_json_to_respond(w, "upgrade your key as it's time ran out", http.StatusUpgradeRequired)
			return
		}
		// email, name := getEmailAndNameFormKey(resultFormChannel.string_value)
		email := userFormKey.Email
		name := userFormKey.UserName
		response := responseStruct{Message: "success", Status_code: http.StatusOK, Success: true, Email: email, Name: name, EncryptedKey: request.Key}
		fmt.Printf("\n\n the user struct is -> %v \n\n", userFormKey)
		println("the decoded json string is ->", userFormKey.GetDecryptedStringInTheStruct())
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			returnTheJsonResponseonError("error encoding json", http.StatusInternalServerError, false, w)
			println("the error in the json encoding is ->", err.Error()+"\n and the value returned is ->", resultFormChannel.Result)
			return
		}
	}
}

func whatCausedJSONDecodeError(e error) (string, int64) {
	switch e.(type) {
	case *json.SyntaxError:
		return "error is caused by bad json syntax", http.StatusBadRequest
	case *json.UnmarshalTypeError:
		return "error is caused by bad json type", http.StatusBadRequest
	case *json.InvalidUnmarshalError:
		return "error is caused by bad json field unmarshal", http.StatusBadRequest
	}
	println("the json decoding error is due to ->", e.Error())
	return "something went wrong on our side", http.StatusInternalServerError
}

func getEmailAndNameFormKey(k string) (email, name string) {
	strings := strings.Split(k, "-|-")
	println(strings)
	return strings[1], strings[2]
}
