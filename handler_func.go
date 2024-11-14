package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type JsonError_HTTPErrorCode_And_Message struct {
    Message string `json:"message"`
    Status_code int64 `json:"status_code"`
}

type string_and_error_channel struct {
    err error
    string_value string
}
type string_and_error_channel_for_subtitles struct {
    err error
    string_value string
	transcript *Transcripts
}



func User_signup_handler(os_env_key string) http.HandlerFunc {
  
	return func(w http.ResponseWriter, r *http.Request) {
		println("in signup ")
	
	//  time1:= time.Now()
  // when user signs up I want them to send me 
	if r.Method != http.MethodPost{
	  http.Error(w, "Invalid request method", http.StatusBadRequest)
	  return
	}
	var signup_user_details Signup_detail_of_user
  
	// Parsing JSON 
	err := json.NewDecoder(r.Body).Decode(&signup_user_details)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}
	// checking if the user has not provided the field
	if signup_user_details.AccountID == 0 || signup_user_details.Email == "" || signup_user_details.UserToken == "" {
	  http.Error(w, "Missing required fields", http.StatusBadRequest)
	  println("-->",signup_user_details.AccountID, "--", signup_user_details.Email, "--", signup_user_details.UserToken)
	  return
	}
	println( "User signed up successfully", signup_user_details.AccountID, "-", signup_user_details.Email, " - ", signup_user_details.UserToken )
    db := DbConnect()
  
	  errChan := make(chan error, 1)
  
	  // Run InsertUserInDB asynchronously
	  go func() {
		  userToInsert := userForDB{
			  accountid:      signup_user_details.AccountID,
			  email:          signup_user_details.Email,
			  UserToken:      signup_user_details.UserToken,
			  is_a_paid_user: false, // Assuming default is false
		  }
		  
		  err := InsertUserInDB(db, userToInsert) // Assuming 'db' is accessible here
		  errChan <- err // Send error (or nil) to channel
	  }()
  
	  // Wait for the goroutine to finish or timeout
	  select {
	  case err := <-errChan:
		  if err != nil {
			  http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
			  log.Printf("Error inserting user into DB: %v", err)
			  return
		  }
	  case <-time.After(5 * time.Second): // Timeout after 5 seconds
		  http.Error(w, "Database operation timed out", http.StatusInternalServerError)
		  return
	  }
  
	// now add the user to the db and give them the key
  
	key := []byte(os_env_key)
	  // --------- fill this
	  plaintext := []byte(return_string_based_on_user_details_for_encryption_text(signup_user_details, false) )
	  // -------fill this
	  // Encrypt	
	  
	  ciphertext, err := encrypt(plaintext, key)
	  if err != nil {
		  panic(err)
		}
		
		base64Ciphertext := base64.StdEncoding.EncodeToString(ciphertext)
		println("\n\nabout to encrypt -->", return_string_based_on_user_details_for_encryption_text(signup_user_details, false),"\n\n", "and ",base64Ciphertext ,"\n\n")
  
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(base64Ciphertext)
		// should have a strucutred json response on this path
  }
}


type request_for_youtubeVideo_struct struct{
	Youtube_Video_Id string `json:"youtube_Video_Id"`
	Encrypted_string string `json:"encrypted_string"`
}


func Return_to_client_where_to_skip_to_in_videos(os_env_key []byte, httpClient *http.Client) http.HandlerFunc {
// take the video id out and hash  ,  and api will return (on success)

//  ads : boolean, if true then starts at _ _ _ and ends at _ _ _ 
// also has to search for the strings in the transcript by myself (probally should be using a on user device llm(like apple on device and cloud too) from gemini 
// but there accuraccy is meh! I think )

return func(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost{
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		// this cpuld return in error , if there is a error decoding json then I am sending the same error in the request code 
		err:=json.NewEncoder(w).Encode(JsonError_HTTPErrorCode_And_Message{Message:"Invalid request method", Status_code:http.StatusBadRequest})
		if err!= nil{
		// idk
		}
		return
	}

	var request_for_youtubeVideo_struct request_for_youtubeVideo_struct
	err:= json.NewDecoder(r.Body).Decode(&request_for_youtubeVideo_struct)

	if err!= nil{
		println(err.Error())
		http.Error(w, "something went wrong on out side", http.StatusInternalServerError)
	  json.NewEncoder(w).Encode(JsonError_HTTPErrorCode_And_Message{Status_code: http.StatusInternalServerError, Message: "something went wrong on out side" })
	  return
	}

	if request_for_youtubeVideo_struct.Youtube_Video_Id == ""{
		method_to_write_http_and_json_to_respond(w, "Parameter youtube_video_id  not provided",http.StatusBadRequest)
	  return
	}

	// now got the  video Id and the  encrypted string now we need to make a go routing to see whether the encrypted string is valid and fetch subtitles of the
	// yt video, and if the id is valid then we can just send it to groq based on the account of the user \

	channel_for_userDetails := make(chan string_and_error_channel )
	channel_for_subtitles := make(chan string_and_error_channel_for_subtitles )
	// channel_for_subtitles := make(chan string_and_error_channel )

	var transcript *Transcripts
	go decrypt_and_write_to_channel(request_for_youtubeVideo_struct.Encrypted_string, os_env_key, channel_for_userDetails)
	go Get_the_subtitles(*httpClient , request_for_youtubeVideo_struct.Youtube_Video_Id,  channel_for_subtitles, transcript ) 

	result_for_user_details := <- channel_for_userDetails

	if result_for_user_details.err != nil {
		method_to_write_http_and_json_to_respond(w,"Something is wrong with your encrypted string", http.StatusBadRequest)
		return 
	}
	// if the paid user has the 
	if len(result_for_user_details.string_value) >=3000{

	}
	userInDb , err := returnUserInDbFormEncryptedString(result_for_user_details.string_value)
	if err!= nil{
		// this could be a bad request too 
		println(" in the returnUserInDbFormEncryptedString's error -->", err.Error(), "\n")
		method_to_write_http_and_json_to_respond(w,"Something went wron on out side , error recognizing you form the auth token", http.StatusInternalServerError)
	}
	
	println("result_for_user_details--++",result_for_user_details.string_value,"\n user in db is -> ",userInDb.UserToken,userInDb.AccountID,userInDb.Email,userInDb.paid_status)

	result_for_subtitles := <- channel_for_subtitles
	if result_for_subtitles.err != nil {
		method_to_write_http_and_json_to_respond(w,"Something is wrong on our side", http.StatusInternalServerError)
		println("error in result_for_subtitles.err --> ", result_for_subtitles.err.Error())
	}
	print("\n string value is this --> ",result_for_subtitles.string_value, "<--string value was this ")
	
	// what about the free user and paid user channel/key_channel and prompt the groq 
	channel_for_groqResponse := make(chan String_and_error_channel_for_groq_response)
	apiKey, err := getAPIKEYForGroqBasedOnUsersTeir(userInDb.paid_status)
	if err != nil {
		method_to_write_http_and_json_to_respond(w,"Something is wrong on our side, error generating a random number", http.StatusInternalServerError)
		println("error in result_for_subtitles.err --> ", result_for_subtitles.err.Error())
	}
	println("and the random key picked by the logic is --> ",apiKey)

	go AskGroqabouttheSponsorship(httpClient, channel_for_groqResponse, apiKey, &result_for_subtitles.string_value)
	groq_response := <- channel_for_groqResponse

	if groq_response.err != nil && groq_response.groqApiResponsePtr == nil {
		if groq_response.http_response_for_go_api_ptr.StatusCode == 429 {
			method_to_write_http_and_json_to_respond(w,"the request time out on this tier", http.StatusTooManyRequests)
			return
		}else if groq_response.groqApiResponsePtr == nil {
			// erro decoding json 
			println(groq_response.err.Error())
			method_to_write_http_and_json_to_respond(w,"somethign went wrong on our side", http.StatusInternalServerError)
			return 
		}else{
			// should only do explicit comparison
			// println("error in the groq request -->", groq_response.err.Error(), "   status code -->", groq_response.http_response_for_go_api_ptr.StatusCode)
			// method_to_write_http_and_json_to_respond(w,"Something is wrong with your encrypted string", http.StatusBadRequest)
			// return
		}
	}
	// printJson(groq_response)
	println("=------ending---=-=--")
	if groq_response.groqApiResponsePtr.Choices[0].Message.Content != ""{
		var sponsorshipContent SponsorshipContent
		err := json.Unmarshal([]byte(groq_response.groqApiResponsePtr.Choices[0].Message.Content), &sponsorshipContent)
		if err != nil {
			println("error in the json unmarshall in the sponsrsgip content --", err.Error())
		}else{
			println("\n ++++",sponsorshipContent.DoesVideoHaveSponsorship, "--",sponsorshipContent.SponsorshipSubtitle)
		}

	
	}else{
		println("content is empty")
	}

	// getting error deciding the escaped json in the json response




	// For now, we'll just send it back as a response
	w.Header().Set("Content-Type", "application/json")
	method_to_write_http_and_json_to_respond(w,"in the end", http.StatusOK)

}
}

