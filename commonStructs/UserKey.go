package commonstructs

type UserKey struct {
	AccountID  string `json:"account_id"`
	Email      string `json:"email"`
	UserName   string `json:"user_name"`
	IsUserPaid bool   `json:"is_user_paid"` // Assuming default is false
	//
	//
	//
	//the above feilds could be used as a foreign key and the new db row will have
	//
	UserTier            string `json:"user_tier"`               // could only be free tier | recurring | one time
	Version             int64  `json:"version"`                 // this is used to compare it to the db version and get the new version from the DB, default is 0
	CheckForKeyUpdateOn int64  `json:"check_for_key_update_on"` // this field is used when we will decrypt the key on the server and check if the time is more than the time on the server. If yes then ask user to update the key
}

// now the user message table will be
//
// for updating we will need to use the
