package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"
)



func User_signup_handler(os_env_key string) http.HandlerFunc {
  
	return func(w http.ResponseWriter, r *http.Request) {
	  db := DbConnect()
	//  time1:= time.Now()
  // when user signs up I want them to send me 
	if r.Method != http.MethodPost{
	  http.Error(w, "Invalid request method", http.StatusBadRequest)
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