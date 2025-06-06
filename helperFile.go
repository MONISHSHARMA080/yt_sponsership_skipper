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

	_ "github.com/tursodatabase/go-libsql"
	// _ "github.com/tursodatabase/libsql-client-go/libsql"
)

type UserInDb struct {
	accounID       int64
	email          string
	userName       string
	is_a_paid_user bool
}

func DbConnect() *sql.DB {
	url := ""
	dbURL := ""
	authToken := ""
	isThisTestingEnv := os.Getenv("IS_THIS_TESTING_ENVIRONMENT")
	if isThisTestingEnv == "true" {
		dbURL = os.Getenv("TURSO_DATABASE_URL")
		url = dbURL
	} else {
		// in any case we are in prod
		dbURL = os.Getenv("TURSO_DATABASE_URL")
		authToken = os.Getenv("TURSO_AUTH_TOKEN")
		url = fmt.Sprintf("%s?authToken=%s", dbURL, authToken)
	}

	// println(url,"\n\n")
	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}

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

func decrypt_and_write_to_channel(ciphertextAsString string, EnvKey []byte, channErr chan<- string_and_error_channel) {
	// First, decode the base64 encoded string
	println("in the decrypt_and_write_to_channel ")
	ciphertextAsByte, err := base64.StdEncoding.DecodeString(ciphertextAsString)
	if err != nil {
		channErr <- string_and_error_channel{err: fmt.Errorf("failed to decode base64: %v", err), string_value: ""}
		return
	}

	// Now decrypt the actual ciphertext
	stringAsByte, err := decrypt(ciphertextAsByte, EnvKey)
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
	println("---- writting to the json  ----")
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
		fmt.Printf("Server error while decoding JSON for the user --- %+v --- is ->%s \n\n", err.Error())
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
	println("converting the acc ID to number, id is  -->", usr.AccountID)
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
