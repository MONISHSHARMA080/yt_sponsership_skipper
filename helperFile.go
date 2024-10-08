package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
  "encoding/json"


	_ "github.com/tursodatabase/go-libsql"
)

type Signup_detail_of_user struct {
  AccountID  int64    `json:"account_id"`
  Email      string `json:"email"`
  UserToken  string `json:"user_token"`
}

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
  return fmt.Sprintf("%d-%s-%s-%s", user_detail.AccountID, user_detail.Email, user_detail.UserToken, paid_status)
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