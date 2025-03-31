package structs

import (
	"database/sql"
	"fmt"
	"time"
	commonhelperfuncs "youtubeAdsSkipper/commonHelperFuncs"
	"youtubeAdsSkipper/paymentBackendGO/common"
)

type MessageForUserOnPaymentCapture struct {
	// AccountID  string `json:"account_id"`
	// Email      string `json:"email"`
	// UserName   string `json:"user_name"`
	// IsUserPaid bool   `json:"is_user_paid"` // Assuming default is false
	//
	//
	//
	//the above feilds could be used as a foreign key and the new db row will have
	//
	UserTier            string `json:"user_tier"`               // could only be free tier | recurring | one time
	Version             int64  `json:"version"`                 // this is used to compare it to the db version and get the new version from the DB, default is 0, and will be automatically set
	CheckForKeyUpdateOn int64  `json:"check_for_key_update_on"` // if it is recurring then make it after 1 month +1 day
	UserAccountID       int64  `json:"user_account_id"`         // this os for the primary key identification
	RazorpayPaymentID   string `json:"razorpay_payment_id"`
}

func (messageOnPaymentCapture *MessageForUserOnPaymentCapture) AddMessageAfterUserPaymentReceived(db *sql.DB, resultChannel chan common.ErrorAndResultStruct[bool]) {
	if messageOnPaymentCapture.IsMyStructUninitialized() {
		resultChannel <- common.ErrorAndResultStruct[bool]{Error: fmt.Errorf("the struct is not initialized"), Result: false}
		return
	}
	// take some sort of channel that will have error and return a  messageOnPaymentCapture and a error as a erturn type
	//
	//
	//
	// There is a error in the query such that It will fail, I want the primary key UserAccountID to autoIncrement and the
	// foreign key to be passed in
	//
	//
	//
	println("the user account ID is ->", messageOnPaymentCapture.UserAccountID)
	query := `
INSERT INTO messageForTheUserAfterPaymentCaptured (
        user_account_id, 
        user_tier, 
        razorpay_payment_id,
        check_for_key_update_on, 
        version
    )
    VALUES (
        ?, 
        ?, 
        ?,
        ?, 
        COALESCE(
            (SELECT MAX(version) FROM messageForTheUserAfterPaymentCaptured 
             WHERE user_account_id = ?), 
            0
        ) + 1
    )`

	_, err := db.Exec(query,
		messageOnPaymentCapture.UserAccountID,
		messageOnPaymentCapture.UserTier,
		messageOnPaymentCapture.RazorpayPaymentID,
		messageOnPaymentCapture.CheckForKeyUpdateOn,
		messageOnPaymentCapture.UserAccountID)
	if err != nil {
		resultChannel <- common.ErrorAndResultStruct[bool]{Error: err, Result: false}
		return
	}
	resultChannel <- common.ErrorAndResultStruct[bool]{Error: nil, Result: true}
}

// after creating the struct set the values in it, the version will be figured out by DB, note vlaid tiers
// are: recurring | one time | free tier
func (messageForUserOnPaymentCapture *MessageForUserOnPaymentCapture) InitializeStruct(RazorpayPaymentID string, UserAccountID int64, isOneTimePayment bool) error {
	UserTier := "recurring"
	if isOneTimePayment {
		UserTier = "one time"
	}

	validTiers := map[string]bool{
		"recurring": true,
		"free tier": true,
		"one time":  true,
	}
	// make some sort of check to see if the UserTier is a string that DB will take
	if !validTiers[UserTier] {
		println("\n\n\n\n\n\n\n\n\n\n\\n\n\n\n\n\n\n\n\n\n\n\n=====================================================================================================")
		println("we have a user tier that is not supposed to be , we should have crashed the program but we are setting the time to update still, tier is  ", UserTier)
		println("\n\n\n\n\n\n\n\n\n\n\\n\n\n\n\n\n\n\n\n\n\n\n=====================================================================================================")
		return fmt.Errorf("the UserTier string is not of a valid category")
	}
	// timeToCheckForKeyUpdate, err := messageForUserOnPaymentCapture.GetTimeToCheckForKeyUpdateOn(UserTier)
	// if err != nil {
	// 	return fmt.Errorf("can't get the time to check for key on update")
	// }
	// make some sort of check to see if the time here is more than the time on the server
	timeToCheckForKeyUpdate := commonhelperfuncs.GetTimeToExpireTheKey(false)

	if time.Now().Unix() >= timeToCheckForKeyUpdate {
		// if I am setting the time to be less than the currentTime then I should crash the program
		panic(fmt.Errorf("the time should be less than CheckForKeyUpdateOn"))
	}
	println("the user account is by razorpay is ->", UserAccountID)
	messageForUserOnPaymentCapture.UserTier = UserTier
	messageForUserOnPaymentCapture.CheckForKeyUpdateOn = timeToCheckForKeyUpdate
	messageForUserOnPaymentCapture.UserAccountID = UserAccountID
	messageForUserOnPaymentCapture.RazorpayPaymentID = RazorpayPaymentID
	return nil
}

func (DbField *MessageForUserOnPaymentCapture) IsMyStructUninitialized() bool {
	return DbField.UserTier == "" || DbField.RazorpayPaymentID == "" ||
		DbField.CheckForKeyUpdateOn == 0 && DbField.UserAccountID == 0
}

// will give us the hardcoded time to check for the key; will return error if the UserTier is not a valid one
func (msgForUser *MessageForUserOnPaymentCapture) GetTimeToCheckForKeyUpdateOn(UserTier string) (int64, error) {
	validTiers := map[string]bool{
		"recurring": true,
		"free tier": true,
		"one time":  true,
	}

	// Check if the UserTier is valid
	if !validTiers[UserTier] {
		return 0, fmt.Errorf("invalid user tier: %s", UserTier)
	}

	// Check if it's a tier that requires key update check (exclude free tier)
	if UserTier == "free tier" {
		return 0, nil
	}

	// Calculate time one month and one day from now
	// currentTime := time.Now()
	// timeAfter1monthAnd1Day := currentTime.AddDate(0, 1, 1)

	timeToUpdateKeyOn := commonhelperfuncs.GetTimeToExpireTheKey(false)
	// Convert to Unix timestamp
	return timeToUpdateKeyOn, nil
}

// type DBChannelResult struct {
// 	IsUserThere    bool
// 	IsMessageThere bool
// 	Error          error
// }

// this one gets the latest message for the user and fills the structs, it is intended to be used on the empty struct
// returns false if the message is not there in the DB(no rows)-- maybe update the time to check in Db and return it
//
// think about this, this will not give me message not found error as we are not checking for the message when the  user is on
// free tier, only when they make payment we gave them a fake key and they look after a day (eg) -> we receive a web hook and we set the message ->
// user check for the new key/message and they gets it -> they do it after the expiration time and then if the user is on free teir or on paid we need to check
//
// here we are assuming that the email is there in the DB if not then give them the error outside form it
func (DbField *MessageForUserOnPaymentCapture) GetLatestMessageForTheUser(db *sql.DB, email string, resultChannel chan common.ErrorAndResultStruct[bool],

// resultChannel2 chan common.ErrorAndResultStruct[DBChannelResult],
) {
	query := `
    SELECT 
        m.user_tier,
        m.version,
        m.check_for_key_update_on,
        m.user_account_id,
        m.razorpay_payment_id
    FROM messageForTheUserAfterPaymentCaptured m
    JOIN UserAccount ua ON m.user_account_id = ua.id
    WHERE ua.email = ?
    ORDER BY m.version DESC
    LIMIT 1`

	rowsReturned := db.QueryRow(query, email)

	err := rowsReturned.Scan(
		&DbField.UserTier,
		&DbField.Version,
		&DbField.CheckForKeyUpdateOn,
		&DbField.UserAccountID,
		&DbField.RazorpayPaymentID,
	)
	if err != nil {
		// this will be true if the email is not there or for the message is not there
		if err == sql.ErrNoRows {
			println("there are no rows in the DB(err)")
			resultChannel <- common.ErrorAndResultStruct[bool]{Result: false, Error: nil}
			return
		}
		println("\n\n the error in the getting latest message form the DB is ->", err.Error(), "\n\n")
		resultChannel <- common.ErrorAndResultStruct[bool]{Result: false, Error: err}
	} else {
		resultChannel <- common.ErrorAndResultStruct[bool]{Result: true, Error: nil}
	}
}
