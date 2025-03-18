package structs

import (
	"database/sql"
	"youtubeAdsSkipper/paymentBackendGO/common"
)

type TemporaryFieldToVerifyPaymentLater struct {
	UserAccountID    int
	RecurringOrderID string
	OnetimeOrderID   string
}

// factoryFunction
func a(DB *sql.DB, RecurringOrderID, OnetimeOrderID string) *TemporaryFieldToVerifyPaymentLater {
	var dbFeild TemporaryFieldToVerifyPaymentLater = TemporaryFieldToVerifyPaymentLater{UserAccountID: 0, RecurringOrderID: RecurringOrderID, OnetimeOrderID: OnetimeOrderID}
	return &dbFeild
}

func (DbField *TemporaryFieldToVerifyPaymentLater) getTheUserPrimaryKeyFormEmail(email string, resultChan chan common.ErrorAndResultStruct[int64]) {
}
