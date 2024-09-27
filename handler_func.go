package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"fmt"
)

type JsonError_HTTPErrorCode_And_Message struct {
    Message string `json:"message"`
    Status_code int64 `json:"status_code"`
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
  
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(base64Ciphertext)
	// json.NewEncoder(w).Encode(string(time.Since(time1).Milliseconds()))
	// fmt.Printf("\n\n\n time taken in user_signup_func is %d", time.Since(time1).Milliseconds())
	}
}


type request_for_youtubeVideo_struct struct{
	Youtube_Video_Id string `json:"youtube_Video_Id"`
	Encrypted_string string `json:"encrypted_string"`
}


func Return_to_client_where_to_skip_to_in_videos(os_env_key string) http.HandlerFunc {
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
	
	// all_url_params := r.URL.Query()
	// youtube_video_id := all_url_params.Get("youtube_video_id")
	
	// if youtube_video_id == ""{
	//   http.Error(w, "Parameter youtube_video_id  not provided", http.StatusBadRequest)
	//   json.NewEncoder(w).Encode(JsonError_HTTPErrorCode_And_Message{Message:"Parameter youtube_video_id  not provided", Status_code:http.StatusBadRequest  })
	//   return
	// }

	var request_for_youtubeVideo_struct request_for_youtubeVideo_struct
	err:= json.NewDecoder(r.Body).Decode(&request_for_youtubeVideo_struct)

	if err!= nil{
		println(err.Error())
		http.Error(w, "something went wrong on out side", http.StatusInternalServerError)
	  json.NewEncoder(w).Encode(JsonError_HTTPErrorCode_And_Message{Status_code: http.StatusInternalServerError, Message: "something went wrong on out side" })
	  return
	}

	if request_for_youtubeVideo_struct.Youtube_Video_Id == ""{
	  http.Error(w, "Parameter youtube_video_id  not provided", http.StatusBadRequest)
	  json.NewEncoder(w).Encode(JsonError_HTTPErrorCode_And_Message{Message:"Parameter youtube_video_id  not provided", Status_code:http.StatusBadRequest  })
	  return
	}
	// now got the  video Id and the  encrypted string now we need to make a go routing to see whether the encrypted string is valid and fetch subtitles of the
	// yt video, and if the id is valid then we can just send it to groq based on the account of the user 

	// For now, we'll just send it back as a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}

}


func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}
