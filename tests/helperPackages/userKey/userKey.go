package userkey

import (
	commonstructs "youtubeAdsSkipper/commonStructs"
)

type UserKey struct {
	User commonstructs.UserKey
}

// make a func that will take the userInDB fields and will set it on the usr In Key and will give you the Encrypted key for the user
func (userKey *UserKey) InitializeTheStructAndGetEncryptedKey() {
}
