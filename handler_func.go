package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	commonstructs "youtubeAdsSkipper/commonStructs"
	"youtubeAdsSkipper/paymentBackendGO/common"
)

type JsonError_HTTPErrorCode_And_Message struct {
	Message     string `json:"message"`
	Status_code int64  `json:"status_code"`
}

type string_and_error_channel struct {
	err          error
	string_value string
}

type string_and_error_channel_for_subtitles struct {
	err          error
	string_value string
	transcript   *Transcripts
}

type Signup_detail_of_user struct {
	AccountID int64  `json:"account_id"`
	UserToken string `json:"user_token"`
}
type Signup_detail_of_user_temp struct {
	AccountID string `json:"account_id"`
	UserToken string `json:"user_token"`
}
type ResponseFromTheUserAuthStruct struct {
	Message      string `json:"message"`
	Status_code  int64  `json:"status_code"`
	Success      bool   `json:"success"`
	EncryptedKey string `json:"encrypted_key"`
}

func User_signup_handler(os_env_key string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ResponseFromTheUserAuth := ResponseFromTheUserAuthStruct{}
		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusBadRequest)
			a := ResponseFromTheUserAuthStruct{Message: "Invalid request method", Status_code: http.StatusBadRequest, Success: false}
			a.writeJSONAndHttpForUserSignupFunc(w)
			return
		}

		var signup_user_details_temp Signup_detail_of_user_temp // to refactor the function

		// Parsing JSON
		err := json.NewDecoder(r.Body).Decode(&signup_user_details_temp)
		if err != nil {
			// http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			ResponseFromTheUserAuth.handleJSONSentByUserError(err, w)
			return
		}
		// checking if the user has not provided the field
		println("22")

		db := DbConnect()
		resultAndErrChan := make(chan common.ErrorAndResultStruct[string])
		responseFromGoogleAuth := make(chan TokenResponseFromGoogleAuth)

		// --------------- check for the oauth and see wether the user detail is true or not and then take the email(from OAUTH2) etc, name etc and put that in the db

		go verifyGoogleAuthToken(signup_user_details_temp.UserToken, responseFromGoogleAuth)
		responseFormGoogleAuthToken := <-responseFromGoogleAuth
		if responseFormGoogleAuthToken.Error != nil {
			println("error occurred in the google Oauth", responseFormGoogleAuthToken.Error.Error())
			ResponseFromTheUserAuth.Status_code = int64(responseFormGoogleAuthToken.StatusCode)
			ResponseFromTheUserAuth.Message = "can't authenticate you"
			ResponseFromTheUserAuth.Success = false
			println("4")
			err := ResponseFromTheUserAuth.writeJSONAndHttpForUserSignupFunc(w)
			if err != nil {
				println("problem with json encoding in the struct method", err.Error())
			}
			return
		}
		// ---> now make the json response func <---
		// -------- this one should send the json response and not http error

		println("making the db struct")
		userToInsert := commonstructs.UserInDb{
			AccountID:  signup_user_details_temp.AccountID,
			Email:      responseFormGoogleAuthToken.Email,
			UserName:   responseFormGoogleAuthToken.Name,
			IsUserPaid: false, //  default is false for the new user
		}

		// println("is the Db struct valid (not added the free tier so it should be false->", userToInsert.IsUserValid())
		// userToInsert.AddUserToFreeTier()
		// println("is the Db struct valid added the free tier so it should be true->", userToInsert.IsUserValid())
		go userToInsert.InsertNewUserInDbAndGetNewKey(db, resultAndErrChan)
		// var encryptedKeyOrResult string
		// Wait for the goroutine to finish or timeout
		println("waiting for the channel to receive on  ")

		select {
		case result := <-resultAndErrChan:
			if result.Error != nil {
				println("there is a error in inserting teh user in Db ->", result.Error.Error())
				ResponseFromTheUserAuth.Message = "error inserting you in the DB"
				ResponseFromTheUserAuth.Success = false
				ResponseFromTheUserAuth.Status_code = http.StatusInternalServerError
				ResponseFromTheUserAuth.writeJSONAndHttpForUserSignupFunc(w)
				fmt.Printf("Error inserting user into DB: %v", result.Error)
				return
			} else {
				println("there is no problem in getting the key and we are returning the encrypted key to be ->", result.Result)
				ResponseFromTheUserAuth.Message = "successfully completed user signup"
				ResponseFromTheUserAuth.Success = true
				ResponseFromTheUserAuth.Status_code = http.StatusOK
				ResponseFromTheUserAuth.EncryptedKey = result.Result
				ResponseFromTheUserAuth.writeJSONAndHttpForUserSignupFunc(w)
				w.Header().Set("Content-Type", "application/json")
			}
		case <-time.After(8 * time.Second): // Timeout after 6 seconds
			println("timeout after 8 sec")
			ResponseFromTheUserAuth.Message = "timeout while inserting you in the DB"
			ResponseFromTheUserAuth.Success = false
			ResponseFromTheUserAuth.Status_code = http.StatusInternalServerError
			ResponseFromTheUserAuth.writeJSONAndHttpForUserSignupFunc(w)
			return
		}

		// new user the version will be 0
	}
}

type request_for_youtubeVideo_struct struct {
	Youtube_Video_Id string `json:"youtube_Video_Id"`
	Encrypted_string string `json:"encrypted_string"`
}

type responseForWhereToSkipVideo struct {
	Status_code            int    `json:"status"`
	Message                string `json:"message"`
	StartTime              int64  `json:"startTime"`
	EndTime                int64  `json:"endTime"`
	ContainSponserSubtitle bool   `json:"containSponserSubtitle"`
	Error                  string `json:"error,omitempty"`
}

func Return_to_client_where_to_skip_to_in_videos(os_env_key []byte, httpClient *http.Client) http.HandlerFunc {
	// take the video id out and hash  ,  and api will return (on success)

	//  ads : boolean, if true then starts at _ _ _ and ends at _ _ _
	// also has to search for the strings in the transcript by myself (probally should be using a on user device llm(like apple on device and cloud too) from gemini
	// but there accuraccy is meh! I think )

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusBadRequest)
			// this cpuld return in error , if there is a error decoding json then I am sending the same error in the request code
			_ = json.NewEncoder(w).Encode(responseForWhereToSkipVideo{Message: "Invalid request method", Status_code: http.StatusBadRequest, ContainSponserSubtitle: false})
			// if err != nil {
			// 	// idk
			// }
			return
		}
		var request_for_youtubeVideo_struct request_for_youtubeVideo_struct
		err := json.NewDecoder(r.Body).Decode(&request_for_youtubeVideo_struct)
		if err != nil {
			println(err.Error())
			_, errorInUserResponse := err.(*json.UnmarshalTypeError)
			if errorInUserResponse {
				http.Error(w, "request json object does not match the expected result", http.StatusBadRequest)
				json.NewEncoder(w).Encode(responseForWhereToSkipVideo{Status_code: http.StatusBadRequest, Message: "request json object does not match the expected result", ContainSponserSubtitle: false})
				return
			}
			http.Error(w, "something went wrong on out side", http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responseForWhereToSkipVideo{Status_code: http.StatusInternalServerError, Message: "something went wrong on out side", ContainSponserSubtitle: false})
			return
		}

		if request_for_youtubeVideo_struct.Youtube_Video_Id == "" {
			method_to_write_http_and_json_to_respond(w, "Parameter youtube_video_id  not provided", http.StatusBadRequest)
			return
		}

		channelToDecryptUserKey := make(chan common.ErrorAndResultStruct[string])
		userFormKey := commonstructs.UserKey{}

		// channel_for_userDetails := make(chan string_and_error_channel)
		channel_for_subtitles := make(chan string_and_error_channel_for_subtitles)
		ChanForResponseForGettingSubtitlesTiming := make(chan ResponseForGettingSubtitlesTiming)

		go userFormKey.DecryptTheKey(request_for_youtubeVideo_struct.Encrypted_string, channelToDecryptUserKey)
		// go decrypt_and_write_to_channel(request_for_youtubeVideo_struct.Encrypted_string, os_env_key, channel_for_userDetails)
		go Get_the_subtitles(*httpClient, request_for_youtubeVideo_struct.Youtube_Video_Id, channel_for_subtitles)
		// getting new user encrypted key decrypted

		// result_for_user_details := <-channel_for_userDetails
		resultForUserKeyChannel := <-channelToDecryptUserKey
		if resultForUserKeyChannel.Error != nil {
			println("there is a erron in decoding the encrypted_key ->", resultForUserKeyChannel.Error.Error())
			method_to_write_http_and_json_to_respond(w, "Something is wrong with your encrypted string", http.StatusBadRequest)
			return
		} else {
			fmt.Printf("resultForUserKeyChannel: %v\n", resultForUserKeyChannel.Result)
		}

		if userFormKey.ShouldWeTellUserToGoGetANewKeyPanic() {
			println("\n\n ==the user should be upgraded as it's time ran out ===\n\n ")
			method_to_write_http_and_json_to_respond(w, "upgrade your key as it's time ran out", http.StatusUpgradeRequired)
			return
		}

		// if result_for_user_details.err != nil {
		//
		// 	method_to_write_http_and_json_to_respond(w, "Something is wrong with your encrypted string", http.StatusBadRequest)
		// 	return
		// }
		// if the paid user has the

		// userInDb, err := returnUserInDbFormEncryptedString(result_for_user_details.string_value)
		// if err != nil {
		// 	// this could be a bad request too
		// 	println(" in the returnUserInDbFormEncryptedString's error -->", err.Error(), "\n")
		// 	method_to_write_http_and_json_to_respond(w, "Something went wron on out side , error recognizing you form the auth token", http.StatusInternalServerError)
		// }

		// println("result_for_user_details--++", result_for_user_details.string_value, "\n user in db is -> ", userInDb.userName, userInDb.accounID, userInDb.email, userInDb.is_a_paid_user)
		// if len(result_for_user_details.string_value) >= 4700 && !userInDb.is_a_paid_user {
		// 	// block it I guess
		// 	println(" +++++++++++string is longer than 3000")
		// 	method_to_write_http_and_json_to_respond(w, "Something is wrong with your encrypted string", http.StatusBadRequest)
		// 	return
		// }
		result_for_subtitles := <-channel_for_subtitles
		if result_for_subtitles.err != nil {
			method_to_write_http_and_json_to_respond(w, "Something is wrong on our side", http.StatusInternalServerError)
			println("error in result_for_subtitles.err --> ", result_for_subtitles.err.Error())
			return
		}
		print("\n string value is this --> ", result_for_subtitles.string_value, "<--string value was this ")

		// what about the free user and paid user channel/key_channel and prompt the groq
		channel_for_groqResponse := make(chan String_and_error_channel_for_groq_response)
		apiKey, err := getAPIKEYForGroqBasedOnUsersTeir(userFormKey.IsUserPaid)
		if err != nil {
			method_to_write_http_and_json_to_respond(w, "Something is wrong on our side, error generating a random number", http.StatusInternalServerError)
			println("error in result_for_subtitles.err --> ", result_for_subtitles.err.Error())
		}
		println("and the random key picked by the logic is --> ", apiKey)

		go AskGroqabouttheSponsorship(httpClient, channel_for_groqResponse, apiKey, &result_for_subtitles.string_value)
		groq_response := <-channel_for_groqResponse

		if groq_response.err != nil && groq_response.groqApiResponsePtr == nil || groq_response.SponsorshipContent == nil {
			if groq_response.http_response_for_go_api_ptr != nil {
				method_to_write_http_and_json_to_respond(w, "somethign went wrong on our side", http.StatusInternalServerError)
				return
			}
			if groq_response.http_response_for_go_api_ptr.StatusCode == 429 {
				method_to_write_http_and_json_to_respond(w, "the request time out on this tier", http.StatusTooManyRequests)
				return
			} else if groq_response.groqApiResponsePtr == nil {
				// erro decoding json
				println("groq error ", groq_response.err.Error())
				method_to_write_http_and_json_to_respond(w, "somethign went wrong on our side", http.StatusInternalServerError)
				return
			}
			if groq_response.SponsorshipContent == nil {
				println(" sponsership content is not there ")
				method_to_write_http_and_json_to_respond(w, "somethign went wrong on our side", http.StatusInternalServerError)
				return
			}
		}
		println("555")
		// getting error deciding the escaped json in the json response
		if groq_response.SponsorshipContent.DoesVideoHaveSponsorship && groq_response.SponsorshipContent.SponsorshipSubtitle != "" {
			// if the subtitles is found
			go GetTimeAndDurInTheSubtitles(result_for_subtitles.transcript, &groq_response.SponsorshipContent.SponsorshipSubtitle, &result_for_subtitles.string_value, ChanForResponseForGettingSubtitlesTiming)
		} else {
			// return no to the user as either the video does not have subtitles or the groq did not return it
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(responseForWhereToSkipVideo{Message: "sponsership subtitles not found in the video", Status_code: http.StatusOK, ContainSponserSubtitle: false})
			if err != nil {
				println("error in the method  encoding the json in the struct 123-->", err.Error())
			}
			return
		}
		SubtitlesTimming := <-ChanForResponseForGettingSubtitlesTiming
		if SubtitlesTimming.err != nil {
			if SubtitlesTimming.err.Error() == "" {
				method_to_write_http_and_json_to_respond(w, "no subtitle found for the video", http.StatusNotFound)
			}
			println("--==", SubtitlesTimming.err.Error())
			method_to_write_http_and_json_to_respond(w, "Something is wrong on our side, error getting subtitles timming ", http.StatusInternalServerError)
			println("error in result_for_subtitles.err --> ", SubtitlesTimming.err.Error())
			return
		}
		err = json.NewEncoder(w).Encode(responseForWhereToSkipVideo{Status_code: http.StatusOK, Error: "", Message: "subtitles found", StartTime: int64(SubtitlesTimming.startTime), EndTime: int64(SubtitlesTimming.endTime), ContainSponserSubtitle: true})
		if err != nil {
			println("error decoding json", SubtitlesTimming.err.Error(), SubtitlesTimming.endTime, SubtitlesTimming.startTime)
			method_to_write_http_and_json_to_respond(w, "Subtitles found but not able to provide json response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
	}
}
