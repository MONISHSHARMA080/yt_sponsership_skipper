package structs

import (
	"database/sql"
	"youtubeAdsSkipper/paymentBackendGO/common"
)

type TemporaryFieldToVerifyPaymentLater struct {
	// UserAccountID    int // commenting as we will get it through email in sql query
	RecurringOrderID string
	OnetimeOrderID   string
}

// factoryFunction
func CreateDBFieldForStoringTempOrderId(RecurringOrderID, OnetimeOrderID string) *TemporaryFieldToVerifyPaymentLater {
	var dbFeild TemporaryFieldToVerifyPaymentLater = TemporaryFieldToVerifyPaymentLater{RecurringOrderID: RecurringOrderID, OnetimeOrderID: OnetimeOrderID}
	return &dbFeild
}

// this struct is small so that I can use it in my struct method and is not ment ot be used outside
func (DbField *TemporaryFieldToVerifyPaymentLater) getTheUserPrimaryKeyFormEmail(db *sql.DB, email string, resultChan chan common.ErrorAndResultStruct[int64]) {
	// this one is designed to run using a go func and
	//
	// wait we can make it a single query
}

// well query here is designed to make sure that every user only gets 1 row in the table, this way prevoius ones are invalidated
// by default. Suppose you have a user token for the user and user requests ne one then we will only update the field
func (DbField *TemporaryFieldToVerifyPaymentLater) SetTokensForTheUser(db *sql.DB, email string, resultChan chan common.ErrorAndResultStruct[string]) {
	query := `
INSERT INTO temporaryFieldToVerifyParymentLater (user_account_id, recurring_order_id, onetime_order_id)
	SELECT id, ?, ? FROM UserAccount WHERE email = ?
	ON CONFLICT(user_account_id) 
	DO UPDATE SET 
	    recurring_order_id = excluded.recurring_order_id,
	    onetime_order_id = excluded.onetime_order_id;`
	_, err := db.Query(query, DbField.RecurringOrderID, DbField.OnetimeOrderID, email)
	if err != nil {
		resultChan <- common.ErrorAndResultStruct[string]{Error: err, Result: ""}
	}
	resultChan <- common.ErrorAndResultStruct[string]{Error: nil, Result: "worked hopefully"}
}

// gets token  based on  the user Email
func (DbField *TemporaryFieldToVerifyPaymentLater) GetTokens(DB *sql.DB, emailOFTheUser string, resultChan chan common.ErrorAndResultStruct[bool]) {
	query := `
	SELECT recurring_order_id, onetime_order_id 
	FROM temporaryFieldToVerifyParymentLater 
	JOIN UserAccount ON temporaryFieldToVerifyParymentLater.user_account_id = UserAccount.id 
	WHERE UserAccount.email = ?`

	err := DB.QueryRow(query, emailOFTheUser).Scan(&DbField.RecurringOrderID, &DbField.OnetimeOrderID)
	if err != nil {
		resultChan <- common.ErrorAndResultStruct[bool]{Result: false, Error: err}
		return
	}

	resultChan <- common.ErrorAndResultStruct[bool]{Result: true, Error: nil}
}
