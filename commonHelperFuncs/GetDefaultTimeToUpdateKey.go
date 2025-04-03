package commonhelperfuncs

import (
	"os"
	"time"
)

// this func gets the time to update the key on the server from the env file and returns it, if the env file is empty it returns the
// hardcoded value(1 month and 1 day), also the env can only take value till s,m,h ; week etc will be  ignored
//
// use it to avoid using hardocded values and also will help in the integration testing, env key is "TIMEFORKEYEXPIRY"
func GetTimeToExpireTheKey(isThisCallForFakeKey bool) int64 {
	timeForKeyExpiryStr := os.Getenv("TIMEFORKEYEXPIRY")
	defaultExpiryTime := time.Now().AddDate(0, 1, 1).Unix()

	if isThisCallForFakeKey {
		timeForKeyExpiryStr = os.Getenv("TIMEFORFAKEKEYEXPIRY")
		defaultExpiryTime = time.Now().AddDate(0, 0, 1).Unix()
	}
	if timeForKeyExpiryStr == "" {
		return defaultExpiryTime
	}

	println("the time selected form env is ->", timeForKeyExpiryStr)
	timeForExpiry, err := time.ParseDuration(timeForKeyExpiryStr)
	if err != nil {
		return defaultExpiryTime
	}
	println(">>>>>>the time selected for the key is after ", timeForExpiry.Milliseconds(), "ms")
	// make a assert here that if the time to skip the video is more than the (panic)
	return time.Now().Add(timeForExpiry).Unix()
}
