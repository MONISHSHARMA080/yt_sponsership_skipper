package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"strings"
	"sync"
	"time"
	commonstructs "youtubeAdsSkipper/commonStructs"
	"youtubeAdsSkipper/paymentBackendGO/common"
	genericResulttype "youtubeAdsSkipper/pkg/GenericResultType"
	llmreqratelimiter "youtubeAdsSkipper/pkg/LLmReqRateLimiter"
	askllm "youtubeAdsSkipper/pkg/askLLM"
	commonresultchannel "youtubeAdsSkipper/pkg/askLLM/commonResultChannel"
	askllmHelper "youtubeAdsSkipper/pkg/askLLM/groqHelper"

	"go.uber.org/zap"
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
		fmt.Printf("the user request is %+v \n", signup_user_details_temp)
		if signup_user_details_temp.AccountID == "" || signup_user_details_temp.UserToken == "" {
			ResponseFromTheUserAuth.Status_code = http.StatusBadRequest
			ResponseFromTheUserAuth.Message = "Request body is invalid(the account id or the user token is invalid)"
			ResponseFromTheUserAuth.Success = false
			println("3.3")
			err := ResponseFromTheUserAuth.writeJSONAndHttpForUserSignupFunc(w)
			if err != nil {
				println("problem with json encoding in the struct method", err.Error())
			}
			return
		}

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

		userToInsert := commonstructs.UserInDb{
			AccountID:  signup_user_details_temp.AccountID,
			Email:      responseFormGoogleAuthToken.Email,
			UserName:   responseFormGoogleAuthToken.Name,
			IsUserPaid: false, //  default is false for the new user
		}
		fmt.Printf("made the db struct %+v \n\n", userToInsert)
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
	Transcript       string `json:"transcript"`
}

type responseForWhereToSkipVideo struct {
	Status_code            int    `json:"status"`
	Message                string `json:"message"`
	StartTime              int64  `json:"startTime"`
	EndTime                int64  `json:"endTime"`
	ContainSponserSubtitle bool   `json:"containSponserSubtitle"`
	Error                  string `json:"error,omitempty"`
}

func Return_to_client_where_to_skip_to_in_videos(os_env_key []byte, httpClient *http.Client, rateLimiterDb *sql.DB, logger *zap.Logger) http.HandlerFunc {
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
			logger.Warn("invalid method for the path by the user", zap.String("method", r.Method))
			return
		}
		var request_for_youtubeVideo_struct request_for_youtubeVideo_struct
		err := json.NewDecoder(r.Body).Decode(&request_for_youtubeVideo_struct)
		if err != nil {
			println(err.Error())
			_, errorInUserResponse := err.(*json.UnmarshalTypeError)
			if errorInUserResponse {
				logger.Info("error decoding json as the user's request in json is not right", zap.Error(err))
				http.Error(w, "request json object does not match the expected result", http.StatusBadRequest)
				json.NewEncoder(w).Encode(responseForWhereToSkipVideo{Status_code: http.StatusBadRequest, Message: "request json object does not match the expected result", ContainSponserSubtitle: false})
				return
			}
			http.Error(w, "something went wrong on out side", http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responseForWhereToSkipVideo{Status_code: http.StatusInternalServerError, Message: "something went wrong on out side", ContainSponserSubtitle: false})
			logger.Error("json decode failed of th user request", zap.Error(err))
			return
		}

		if request_for_youtubeVideo_struct.Youtube_Video_Id == "" {
			logger.Info("the client left the youtube video id empty", zap.String("Youtube_Video_Id", ""))
			method_to_write_http_and_json_to_respond(w, "Parameter youtube_video_id  not provided", http.StatusBadRequest)
			return
		}
		if request_for_youtubeVideo_struct.Transcript == "" {
			logger.Info("the client left the youtube video Transcript empty", zap.String("Transcript", ""))
			method_to_write_http_and_json_to_respond(w, "Transcript is not present in the request", http.StatusBadRequest)
			return
		}

		channelToDecryptUserKey := make(chan common.ErrorAndResultStruct[string])
		userFormKey := commonstructs.UserKey{}

		// channel_for_userDetails := make(chan string_and_error_channel)
		channel_for_subtitles := make(chan string_and_error_channel_for_subtitles)
		ChanForResponseForGettingSubtitlesTiming := make(chan askllmHelper.ResponseForGettingSubtitlesTiming)

		go userFormKey.DecryptTheKey(request_for_youtubeVideo_struct.Encrypted_string, channelToDecryptUserKey)
		go request_for_youtubeVideo_struct.GetTheTranscript(channel_for_subtitles)

		// result_for_user_details := <-channel_for_userDetails
		resultForUserKeyChannel := <-channelToDecryptUserKey
		if resultForUserKeyChannel.Error != nil {
			// println("there is a erron in decoding the encrypted_key ->", resultForUserKeyChannel.Error.Error())
			logger.Info("there is a error in decoding the encrypted_key of the user ", zap.Error(resultForUserKeyChannel.Error))
			method_to_write_http_and_json_to_respond(w, "Something is wrong with your encrypted string", http.StatusBadRequest)
			return
		} else {
			logger.Info("the user is decrypted successfully from the user key channel and it is ", zap.Any("resultForUserKeyChannel", resultForUserKeyChannel.Result))
		}
		rateLimiterForUser := llmreqratelimiter.RateLimiterForUser{UserEmail: userFormKey.Email}
		chanForRateLimit := make(chan genericResulttype.ErrorAndResultType[bool])
		go rateLimiterForUser.ShouldWeRateLimitUser(rateLimiterDb, userFormKey.UserTier, chanForRateLimit)
		rateLimitTheUser := <-chanForRateLimit
		if rateLimitTheUser.Err != nil {
			// println("there is a error in getting that should we rate limit the user or not, so here is the error->", rateLimitTheUser.Err.Error())
			logger.Warn("there is a error in getting that should we rate limit the user or not", zap.Error(rateLimitTheUser.Err))
			method_to_write_http_and_json_to_respond(w, "Something went wrong on out side", http.StatusInternalServerError)
			return
		}

		if rateLimitTheUser.Result {
			// println("the user is getting rate limited as the result form the channel is true")
			logger.Info("the user is getting rate limited as the result form the channel is true", zap.Bool(" rateLimitTheUser.Result", rateLimitTheUser.Result))
			method_to_write_http_and_json_to_respond(w, "you have excceded your quota limit for today, see you tomorrow", http.StatusTooManyRequests)
			return
		} else {
			logger.Info("the user is not getting rate limited and we will try to add the user to the Db for using the service(increase the rate limiting no for the user ) ", zap.Bool(" rateLimitTheUser.Result", rateLimitTheUser.Result))
		}
		chanForUpdateTheDbForNewReq := make(chan genericResulttype.ErrorAndResultType[bool])
		go rateLimiterForUser.NewRequestMadeUpdateDb(rateLimiterDb, chanForUpdateTheDbForNewReq)

		if userFormKey.ShouldWeTellUserToGoGetANewKeyPanic() {
			// println("\n\n ==the user should be upgraded as it's time ran out ===\n\n ")
			logger.Warn("the user should be upgraded as it's time ran out ", zap.Bool("userFormKey.ShouldWeTellUserToGoGetANewKeyPanic", true))
			method_to_write_http_and_json_to_respond(w, "upgrade your key as it's time ran out", http.StatusUpgradeRequired)
			return
		}
		result_for_subtitles := <-channel_for_subtitles
		if result_for_subtitles.err != nil {
			method_to_write_http_and_json_to_respond(w, "Something is wrong on our side", http.StatusInternalServerError)
			// println("error in result_for_subtitles.err --> ", result_for_subtitles.err.Error())
			logger.Warn("there is a error in getting the result for decoding the transctipt given by the user ", zap.Error(result_for_subtitles.err))
			return
		}

		// print("\n string value is this --> ", result_for_subtitles.string_value, "<--string value was this ")
		logger.Info("the decoded transctipt are received(get a sense of how long the transcript array is ) ", zap.Int("the no of elements in transctipt('s subtitle) array is", len(result_for_subtitles.transcript.Subtitles))) // zap.String("full captions in the transcripts", result_for_subtitles.string_value),
		// zap.Any("result for subtitle channel's result", result_for_subtitles.transcript),

		resultFromSubtitiles := askllmHelper.String_and_error_channel_for_subtitles{Err: result_for_subtitles.err, String_value: result_for_subtitles.string_value, Transcript: result_for_subtitles.transcript}
		resultChannel := make(chan commonresultchannel.ResultAndErrorChannel[askllmHelper.ResponseForWhereToSkipVideo])
		if userFormKey.IsUserPaid {
			logger.Info("the user is paid and we are about to ask the gemini", zap.Skip())
			go askllm.AskGeminiAboutSponsorShipAndGetTheSponsorTiming(result_for_subtitles.string_value, resultFromSubtitiles, ChanForResponseForGettingSubtitlesTiming, resultChannel, logger)
		} else {
			// what about the free user and paid user channel/key_channel and prompt the groq
			apiKey, err := getAPIKEYForGroqBasedOnUsersTeir(userFormKey.IsUserPaid)
			if err != nil {
				method_to_write_http_and_json_to_respond(w, "Something is wrong on our side, error generating a random number", http.StatusInternalServerError)
				logger.Warn("there is somthing wrong with getting the api key for the groq ", zap.Error(err))
				return
			}
			logger.Info(" key picked info  ", zap.String("first 4 word of the key", apiKey[:len(apiKey)-4]), zap.Int("length of the key(groq)", len(apiKey)))
			logger.Info("the user is in free tier and we are about to ask the groq", zap.Skip())
			go askllm.AskGroqAboutSponsorship(httpClient, w, method_to_write_http_and_json_to_respond, apiKey, resultFromSubtitiles, ChanForResponseForGettingSubtitlesTiming, resultChannel, logger)
		}
		result := <-resultChannel
		logger.Info("got the result for the sponsorship", zap.Any("result returned form the channel ", result))
		// if we have a error then log it and either way send the response
		if result.Err != nil {
			fmt.Printf("\n the error in giving the sponsorship of the video is -> %s  \n ", result.Err.Error())
			logger.Warn("the result from the sponsorship is not ok(err) ", zap.Error(result.Err))
		}
		// seeing if we are able to inset the request in the DB
		updatedTheDbOnNewReq := <-chanForUpdateTheDbForNewReq
		if updatedTheDbOnNewReq.Err != nil {
			// well we have spend this respourse, we will let this go, se here is a free lunch for the user
			logger.Warn("there is a error in inserting the new req in the rate limiting DB", zap.Error(updatedTheDbOnNewReq.Err))
		} else {
			logger.Info(" successfully added the req by the user in the rate limiter Db ", zap.Skip())
		}
		err = result.SendResponse(w)
		if err != nil {
			logger.Warn(" we are going to panic as we should have filled the struct in a good way but clearly we did not do that", zap.Error(err))
			panic(err)
		}
	}
}

func (req *request_for_youtubeVideo_struct) GetTheTranscript(channelToReturnSubtitles chan<- string_and_error_channel_for_subtitles) {
	if req.Transcript == "" {
		channelToReturnSubtitles <- string_and_error_channel_for_subtitles{err: fmt.Errorf("the transcript by the user is empty"), string_value: "", transcript: nil}
		return
	}

	lenOfTranscript := len(req.Transcript)
	println("the transcripts sent by user's lenght is ", lenOfTranscript)
	transcripts := Transcripts{}
	errorInXMl := xml.Unmarshal([]byte(req.Transcript), &transcripts)
	if errorInXMl != nil {
		println("there is a error in Unmarshaling the xml in the transcript struct and it is ->", errorInXMl.Error())
		channelToReturnSubtitles <- string_and_error_channel_for_subtitles{err: errorInXMl, string_value: "", transcript: nil}
		return
	}
	lenOfSubtitles := len(transcripts.Subtitles)

	var wg sync.WaitGroup
	noOfWorkers := 12
	noOfWorkers = min(noOfWorkers, lenOfSubtitles)
	// here even if the thing can result in 3.4 (float) but since we are operating on a int the decimal value is discarded
	chunkSize := (lenOfSubtitles + noOfWorkers - 1) / noOfWorkers
	println(" the chunksize is ->", chunkSize, " and the lenOfSubtitles is :", lenOfSubtitles)

	wg.Add(noOfWorkers)
	endArrayLen := 0

	for w := 0; w < noOfWorkers; w++ {
		arrStart := w * chunkSize
		arrEnd := arrStart + chunkSize
		arrEnd = min(arrEnd, lenOfSubtitles)
		endArrayLen += arrEnd
		go func(subtitleArray []Subtitle) {
			defer wg.Done()
			for i, subtitleArrayElement := range subtitleArray {
				subtitleArray[i].Text = SanatizeStrignsForSearchingWithoutPtr(html.UnescapeString(subtitleArrayElement.Text))
				// println("\n -- at index:", i, "++", transcripts.Subtitles[i].Text, "++")
				// fmt.Printf(" -- at index:%d ---|||%s|||---  \n ", i, transcripts.Subtitles[i].Text)
			}
		}(transcripts.Subtitles[arrStart:arrEnd])
	}

	// formatting the transcript to be in utf-8
	wg.Wait()
	// see that the subtitles are perfectly arranged and well sanitized
	// for i, subtitle := range transcripts.Subtitles {
	// 	fmt.Printf(" -- at index:%d ---|||%s|||---  \n ", i, subtitle.Text)
	// }
	println("formatting the transctipt.subtitles.text to be utf-8")

	channelToReturnSubtitles <- string_and_error_channel_for_subtitles{err: nil, string_value: generateSubtitleString(transcripts.Subtitles), transcript: &transcripts}
}

// this one will not use a ptr
// takes the string and cleans it by removign any more than one " "(spaces) and "\n" with " " and also make it in the lower cap
func SanatizeStrignsForSearchingWithoutPtr(s string) string {
	newStr := strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(s), "  ", " "), "\n", " ")
	return strings.Join(strings.Fields(newStr), " ")
}
