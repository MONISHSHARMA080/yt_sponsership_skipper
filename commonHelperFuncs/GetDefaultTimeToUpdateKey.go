package commonhelperfuncs

import (
	"os"
	"time"
)

// this func gets the time to update the key on the server from the env file and returns it, if the env file is empty it returns the
// hardcoded value(1 month and 1 day), also the env can only take value till s,m,h ; week etc will be  ignored
//
// use it to avoid using hardocded values and also will help in the integration testing, env key is
func GetTimeToExpireTheKey() int64 {
	timeForKeyExpiryStr := os.Getenv("TIMEFORKEYEXPIRY")

	defaultExpiryTime := time.Now().AddDate(0, 1, 1).Unix()
	if timeForKeyExpiryStr == "" {
		return defaultExpiryTime
	}

	timeForExpiry, err := time.ParseDuration(timeForKeyExpiryStr)
	if err != nil {
		return defaultExpiryTime
	}
	// make a assert here that if the time to skip the video is more than the (panic)
	return time.Now().Add(timeForExpiry).Unix()
}
