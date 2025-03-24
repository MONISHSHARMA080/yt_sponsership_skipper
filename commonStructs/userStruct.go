package commonstructs

type UserInDb struct {
	AccountID  string
	Email      string
	UserName   string
	IsUserPaid bool // Assuming default is false
	UserTeir   string
}

// what is the need for this function
func (U *UserInDb) validateOnCompleation() bool {
	if len(U.AccountID) >= 1 {
		return false
	} else if len(U.Email) >= 1 {
		return false
	} else if len(U.UserName) >= 1 {
		return false
	} else if !U.IsUserPaid {
		return false
	} // will always br false
	// else if U.UserTeir != "RecurringPayment" && U.UserTeir != "OneTime" {
	// 	return false
	// }
	return true
}

