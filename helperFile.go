package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

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
  println("--++")
}




func User_signup_handler() http.HandlerFunc {
  
  return func(w http.ResponseWriter, r *http.Request) {
    db := DbConnect()
   time1:= time.Now()
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

  key := []byte(os.Getenv("encryption_key"))
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
 print( base64Ciphertext)
  fmt.Printf("time taken in user_signup_func is ", time.Since(time1).Milliseconds())
  }
}


func return_string_based_on_user_details_for_encryption_text(user_detail Signup_detail_of_user, is_paid_user bool) string{
  if is_paid_user == true{
    return string(user_detail.AccountID)+"-"+user_detail.Email+"-"+user_detail.UserToken+"-"+"true"
  }else{
    return string(user_detail.AccountID)+"-"+user_detail.Email+"-"+user_detail.UserToken+"-"+"false"
  }
}