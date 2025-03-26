package structs

import (
	"database/sql"
	"fmt"
	"time"
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
		return fmt.Errorf("the UserTier string is not of a valid category")
	}
	timeToCheckForKeyUpdate, err := messageForUserOnPaymentCapture.getTimeToCheckForKeyUpdateOn(UserTier)
	if err != nil {
		return fmt.Errorf("can't get the time to check for key on update")
	}
	// make some sort of check to see if the time here is more than the time on the server
	if time.Now().Unix() >= timeToCheckForKeyUpdate {
		return fmt.Errorf("the time should be less than CheckForKeyUpdateOn")
	}
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
func (msgForUser *MessageForUserOnPaymentCapture) getTimeToCheckForKeyUpdateOn(UserTier string) (int64, error) {
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
	currentTime := time.Now()
	timeAfter1monthAnd1Day := currentTime.AddDate(0, 1, 1)

	// Convert to Unix timestamp
	return timeAfter1monthAnd1Day.Unix(), nil
}
