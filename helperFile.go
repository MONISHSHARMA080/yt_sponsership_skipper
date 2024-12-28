package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	// _ "github.com/tursodatabase/go-libsql"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type UserInDb struct {
	accounID       int64
	email          string
	userName       string
	is_a_paid_user bool
}

// userForDB
// Signup_detail_of_user

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

func InsertUserInDB(db *sql.DB, userStructToEnter UserInDb) error {

	query := `
        INSERT OR IGNORE INTO UserAccount
        (accountid, email, userName, is_a_paid_user)
        VALUES (?, ?, ?, ?)
    `

	// should also check if the user already exists, if it does then do not insert it
	rows_returned, err := db.Query(query, userStructToEnter.accounID, userStructToEnter.email, userStructToEnter.userName, userStructToEnter.is_a_paid_user)
	println(" things to insert in db is --", query, userStructToEnter.accounID, userStructToEnter.email, userStructToEnter.userName, userStructToEnter.is_a_paid_user)

	if err != nil {
		println("error in the query --", err.Error())
		return err
	}
	error_in_row_returned := rows_returned.Err()
	if error_in_row_returned != nil {
		println("error_in_row_returned", error_in_row_returned.Error())
		return error_in_row_returned
	}
	err = rows_returned.Close()
	if err != nil {
		return err
	}
	return nil
}

func CheckIfTheUserInDb(db *sql.DB, userInDBStruct UserInDb) error {
	// the db has unique fields
	rows_returned, err := db.Query(" SELECT * FROM UserAccount WHERE accountid = ? AND email = ? AND UserToken = ? AND is_a_paid_user = ?  ", userInDBStruct.accounID, userInDBStruct.email, userInDBStruct.userName, userInDBStruct.is_a_paid_user)
	if err != nil {
		return err
	}
	var id int64
	var a string
	var b int64
	var c string
	var is_a_paid_user bool
	for rows_returned.Next() {
		error_returned := rows_returned.Scan(&id, &b, &a, &c, &is_a_paid_user)
		if error_returned != nil {
			return error_returned
		}
		println("id->", id, b, a, c)
	}
	erro := rows_returned.Close()
	if erro != nil {
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

func return_string_based_on_user_details_for_encryption_text(user_detail UserInDb, is_paid_user bool) string {
	paid_status := "false"
	if is_paid_user {
		paid_status = "true"
	}
	return fmt.Sprintf("%d-|-%s-|-%s-|-%s", user_detail.accounID, user_detail.email, user_detail.userName, paid_status)
}

func returnUserInDbFormEncryptedString(decypted_string_of_user_in_db string) (UserInDb, error) {

	parts := strings.Split(decypted_string_of_user_in_db, "-|-")
	if len(parts) < 4 {
		return UserInDb{}, fmt.Errorf("string has less than 4 parts")
	}
	accountID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		// Handle the error if parsing fails, e.g., return a zero-value struct
		return UserInDb{}, fmt.Errorf("can't parse the AccountID string in encreypted string form db to int")
	}
	paid_status_of_user, err := strconv.ParseBool(parts[3])
	if err != nil {
		// Handle the error if parsing fails, e.g., return a zero-value struct
		return UserInDb{}, fmt.Errorf("can't parse the paid_status string in encreypted string form db to bool")
	}
	return UserInDb{accounID: accountID, email: parts[1], userName: parts[2], is_a_paid_user: paid_status_of_user}, nil
}

func write_to_json_a_error_message() {
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

func method_to_write_http_and_json_to_respond(w http.ResponseWriter, message string, http_status_code int64) {
	println(" message in the-->", message)
	// http.Error(w, message, int(http_status_code))
	w.WriteHeader(int(http_status_code))
	err := json.NewEncoder(w).Encode(responseForWhereToSkipVideo{Message: message, Status_code: int(http_status_code)})
	if err != nil {
		println("error in the method____-->", err.Error())
	}
}

func (r *ResponseFromTheUserAuthStruct) writeJSONAndHttpForUserSignupFunc(w http.ResponseWriter) error {
	w.WriteHeader(int(r.Status_code))
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		println("error in the json encoding for the user signup -->", err.Error())
		return err
	}
	return nil
}

func (r *ResponseFromTheUserAuthStruct) handleJSONSentByUserError(err error, w http.ResponseWriter) {
	// do i need to send 400 or 500
	response := ResponseFromTheUserAuthStruct{
		Success: false,
	}

	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case err == io.EOF:
		response.Message = "Empty request body"
		response.Status_code = http.StatusBadRequest
	case errors.As(err, &syntaxError):
		response.Message = "Malformed JSON"
		response.Status_code = http.StatusBadRequest
	case errors.As(err, &unmarshalTypeError):
		println("incorrect json field error -->", err.Error())
		response.Message = "Incorrect JSON field type"
		response.Status_code = http.StatusBadRequest
	default:
		// Log the error for debugging
		println("Server error while decoding JSON:", err.Error())
		response.Message = "Internal server error"
		response.Status_code = http.StatusInternalServerError
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	println(response.Message, response.Status_code)
	w.WriteHeader(int(response.Status_code))
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func (usr *Signup_detail_of_user_temp) convertAccountIDToNumber() (Signup_detail_of_user, error) {
	// Remove any non-digit characters
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, usr.AccountID)
	// Take only first 19 digits
	if len(digits) > 19 {
		digits = digits[:19]
	}
	// Convert to int64
	num, err := strconv.ParseInt(digits, 10, 64)
	if err != nil {
		return Signup_detail_of_user{}, err
	}
	return Signup_detail_of_user{AccountID: num, UserToken: usr.UserToken}, nil
}
