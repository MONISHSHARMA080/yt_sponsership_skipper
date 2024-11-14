package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	// _ "github.com/tursodatabase/go-libsql"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Signup_detail_of_user struct {
  AccountID  int64    `json:"account_id"`
  Email      string `json:"email"`
  UserToken  string `json:"user_token"`
}
type UserInDb struct {
  AccountID  int64    
  Email      string 
  UserToken  string 
  paid_status bool
}

// type GroqApiResponse struct {
//   ID                string             `json:"id"`
//   Object           string             `json:"object"`
//   Created          int64              `json:"created"`
//   Model            string             `json:"model"`
//   Choices          []GroqApiChoice    `json:"choices"`
//   Usage            GroqApiUsage       `json:"usage"`
//   SystemFingerprint string            `json:"system_fingerprint"`
//   XGroq            GroqApiXGroq       `json:"x_groq"`
// }

// type GroqApiChoice struct {
//   Index        int              `json:"index"`
//   Message      GroqApiMessage   `json:"message"`
//   LogProbs     any             `json:"logprobs"`
//   FinishReason string          `json:"finish_reason"`
// }

// type GroqApiMessage struct {
//   Role    string `json:"role"`
//   Content string `json:"content"`
// }

// type GroqApiUsage struct {
//   QueueTime        float64 `json:"queue_time"`
//   PromptTokens     int     `json:"prompt_tokens"`
//   PromptTime       float64 `json:"prompt_time"`
//   CompletionTokens int     `json:"completion_tokens"`
//   CompletionTime   float64 `json:"completion_time"`
//   TotalTokens      int     `json:"total_tokens"`
//   TotalTime        float64 `json:"total_time"`
// }

// type GroqApiXGroq struct {
//   ID string `json:"id"`
// }

// type Message struct {
// 	Role    string `json:"role"`
// 	Content string
// }

// type GroqApiResponse struct {
// 	// ID represents the unique identifier for the chat completion response.
// 	ID string `json:"id,omitempty"`
// 	// Object specifies the type of object returned in the response.
// 	Object string `json:"object,omitempty"`
// 	// Created indicates the timestamp when the response was created.
// 	Created int `json:"created,omitempty"`
// 	// Model specifies the model used for the chat completion.
// 	Model string `json:"model,omitempty"`
// 	// Choices represents a slice of choice structures containing information about each choice.
// 	Choices []struct {
// 		// Index specifies the index of the choice.
// 		Index int `json:"index,omitempty"`
// 		// Message contains the message content of the choice.
// 		Message Message `json:"message,omitempty"`
// 		// Logprobs represents the log probabilities of the choice.
// 		Logprobs interface{} `json:"logprobs,omitempty"`
// 		// FinishReason indicates the reason why the choice was finished.
// 		FinishReason string `json:"finish_reason,omitempty"`
// 	} `json:"choices,omitempty"`
// 	// Usage contains usage statistics for the chat completion.
// 	Usage struct {
// 		// QueueTime specifies the time spent in the queue.
// 		QueueTime float64 `json:"queue_time,omitempty"`
// 		// PromptTokens indicates the number of tokens in the prompt.
// 		PromptTokens int `json:"prompt_tokens,omitempty"`
// 		// PromptTime specifies the time spent processing the prompt.
// 		PromptTime float64 `json:"prompt_time,omitempty"`
// 		// CompletionTokens indicates the number of tokens in the completion.
// 		CompletionTokens int `json:"completion_tokens,omitempty"`
// 		// CompletionTime specifies the time spent generating the completion.
// 		CompletionTime float64 `json:"completion_time,omitempty"`
// 		// TotalTokens indicates the total number of tokens processed.
// 		TotalTokens int `json:"total_tokens,omitempty"`
// 		// TotalTime specifies the total time spent processing the request.
// 		TotalTime float64 `json:"total_time,omitempty"`
// 	} `json:"usage,omitempty"`
// 	// SystemFingerprint represents a unique identifier for the system.
// 	SystemFingerprint string `json:"system_fingerprint,omitempty"`
// 	// XGroq contains additional information about the Groq system.
// 	XGroq struct {
// 		// ID specifies the unique identifier for the Groq system.
// 		ID string `json:"id,omitempty"`
// 	} `json:"x_groq,omitempty"`
// }


func DbConnect() *sql.DB {
  
  
  // println(os.Getenv("TURSO_DATABASE_URL"),os.Getenv("TURSO_AUTH_TOKEN"))
  // url := "libsql://["+os.Getenv("TURSO_DATABASE_URL")+"].turso.io?authToken=["+os.Getenv("TURSO_AUTH_TOKEN")+"]"
  // url := os.Getenv("TURSO_DATABASE_URL")+".?authToken="+os.Getenv("TURSO_AUTH_TOKEN")
  dbURL := os.Getenv("TURSO_DATABASE_URL")
  authToken := os.Getenv("TURSO_AUTH_TOKEN")

  
  url := fmt.Sprintf("%s?authToken=%s", dbURL, authToken)
  // println(url,"\n\n")
  db, err := sql.Open("libsql", url)
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
    os.Exit(1)
  }
 
//   _, b := db.Query("CREATE TABLE UserAccount (    accountid INT,email TEXT,    strUserToken TEXT);")
//   if b!=nil{
//     panic(b.Error()) 
//   }
  return db
}

func InsertUserInDB(db *sql.DB, userStructToEnter userForDB)error{

  query :=`
        INSERT INTO UserAccount (accountid, email, UserToken, is_a_paid_user )
        SELECT ?, ?, ?, ?
        
    `
  // should also check if the user already exists, if it does then do not insert it
rows_returned, err :=  db.Query(query, userStructToEnter.accountid, userStructToEnter.email, userStructToEnter.UserToken, userStructToEnter.is_a_paid_user )

  if err!= nil{
    return err
  }
error_in_row_returned :=rows_returned.Err()
if error_in_row_returned!= nil{
  println(error_in_row_returned.Error())  
  return error_in_row_returned
}
err = rows_returned.Close()
  if err!=  nil {
    return err
  }
  return nil
  

}

func CheckIfTheUserInDb(db *sql.DB, userInDBStruct userForDB)error{
// the db has unique fields 
  rows_returned, err :=  db.Query(" SELECT * FROM UserAccount WHERE accountid = ? AND email = ? AND UserToken = ? AND is_a_paid_user = ?  ", userInDBStruct.accountid, userInDBStruct.email, userInDBStruct.UserToken, userInDBStruct.is_a_paid_user )
if err!= nil{
  return err
}
var id int64
var a string
var b int64
var c string
var is_a_paid_user bool
for rows_returned.Next(){
  error_returned := rows_returned.Scan(&id, &b, &a, &c, &is_a_paid_user)
  if error_returned!= nil {
    return error_returned
  }
  println("id->", id, b, a, c)
}
erro:= rows_returned.Close()
if erro!= nil {
  return erro
}

return nil
}


func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	
  block, err := aes.NewCipher(key)
	
  if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

  iv := ciphertext[:aes.BlockSize]

  if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}


func getRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is my website!\n")
  println("in root")
}




func return_string_based_on_user_details_for_encryption_text(user_detail Signup_detail_of_user, is_paid_user bool) string {
  paid_status := "false"
  if is_paid_user {
      paid_status = "true"
  }
  return fmt.Sprintf("%d-|-%s-|-%s-|-%s", user_detail.AccountID, user_detail.Email, user_detail.UserToken, paid_status)
}

func returnUserInDbFormEncryptedString(decypted_string_of_user_in_db string) (UserInDb, error) {
  
  parts := strings.Split(decypted_string_of_user_in_db, "-|-")
  if len(parts) < 4 {
    return UserInDb{} , fmt.Errorf("string has less than 4 parts")
  }
  accountID, err := strconv.ParseInt(parts[0], 10, 64)
  if err != nil {
    // Handle the error if parsing fails, e.g., return a zero-value struct
    return UserInDb{}, fmt.Errorf("can't parse the AccountID string in encreypted string form db to int")
  }
  paid_status_of_user , err := strconv.ParseBool(parts[3])
  if err != nil {
    // Handle the error if parsing fails, e.g., return a zero-value struct
    return UserInDb{}, fmt.Errorf("can't parse the paid_status string in encreypted string form db to bool")
  }
  return UserInDb{AccountID: accountID, Email: parts[1], UserToken: parts[2], paid_status: paid_status_of_user}, nil
}

func write_to_json_a_error_message(){
  // if encoding json for  the message
}

func decrypt_and_write_to_channel(ciphertextAsString string, key []byte, channErr chan<- string_and_error_channel) {
  // First, decode the base64 encoded string
  println("in the decrypt_and_write_to_channel ")
  ciphertextAsByte, err := base64.StdEncoding.DecodeString(ciphertextAsString)
  if err != nil {
      channErr <- string_and_error_channel{err: fmt.Errorf("failed to decode base64: %v", err), string_value: ""}
      return
  }

  // Now decrypt the actual ciphertext
  stringAsByte, err := decrypt(ciphertextAsByte, key)
  if err != nil {
      channErr <- string_and_error_channel{err: fmt.Errorf("failed to decrypt: %v", err), string_value: ""}
      return
  }

  string_as_string := string(stringAsByte)
  channErr <- string_and_error_channel{err: nil, string_value: string_as_string}
}


func method_to_write_http_and_json_to_respond( w http.ResponseWriter, message string, http_status_code int64){

  http.Error(w, message, int(http_status_code))
  json.NewEncoder(w).Encode(JsonError_HTTPErrorCode_And_Message{Message:message, Status_code:http_status_code  })

}
