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
		  println(" it is done")
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
	channel_for_subtitles := make(chan string_and_error_channel )

	go decrypt_and_write_to_channel(request_for_youtubeVideo_struct.Encrypted_string, os_env_key, channel_for_userDetails)
	go Get_the_subtitles(*httpClient , request_for_youtubeVideo_struct.Youtube_Video_Id, true, channel_for_subtitles ) 

	result_for_user_details := <- channel_for_userDetails

	if result_for_user_details.err != nil {
		method_to_write_http_and_json_to_respond(w,"Something is wrong with your encrypted string", http.StatusBadRequest)
		return 
	}
	
	println("result_for_user_details--++",result_for_user_details.string_value)
	// println("http client is nil -->", http_client == nil)
	
	if err != nil {
		method_to_write_http_and_json_to_respond(w,"Something is wrong on our side", http.StatusInternalServerError)
		println(http.StatusInternalServerError, "-----:::-----" ,err.Error(), "|||||||||\n")
		return 
	}
	result_for_subtitles := <- channel_for_subtitles
	if result_for_subtitles.err != nil {
		method_to_write_http_and_json_to_respond(w,"Something is wrong on our side", http.StatusInternalServerError)
		println("error in result_for_subtitles.err --> ", result_for_subtitles.err.Error())
	}
	print(result_for_subtitles.string_value, "<--string value was this ")
	// ciphertext, err := base64.StdEncoding.DecodeString(request_for_youtubeVideo_struct.Encrypted_string)


	// if err != nil {
	// 	http.Error(w, "Invalid encrypted string", http.StatusBadRequest)
	// 	json.NewEncoder(w).Encode(JsonError_HTTPErrorCode_And_Message{Message: "Invalid encrypted string", Status_code: http.StatusBadRequest})
	// 	return
	// }

	// // Decrypt the ciphertext
	// plaintext, err := decrypt(ciphertext, os_env_key)
	// if err != nil {
	// 	http.Error(w, "Error decrypting string", http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(JsonError_HTTPErrorCode_And_Message{Message: "Error decrypting string", Status_code: http.StatusInternalServerError})
	// 	return
	// }

	// // Convert plaintext bytes to string
	// decryptedString := string(plaintext)

	// // Log the decrypted string
	// log.Printf("Decrypted string: %s", decryptedString)


	// println("value form the channel  is -->", result.string_value, channel_err)
	// For now, we'll just send it back as a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
}

