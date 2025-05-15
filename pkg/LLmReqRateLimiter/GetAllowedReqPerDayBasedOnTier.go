package llmreqratelimiter

import (
	"fmt"
	"os"
	"strconv"
)

// if we have a erro(it will be form programmer mistake only, or the parse error) it will return error instead of panic
//
// returns the int64 req per day the user on a tier is allowed to make
func (rl *RateLimiterForUser) getAllowedReqPerDayBasedOnUserTier(userTier string) (int64, error) {
	arrayOfUserTier := []string{"recurring", "free tier", "one time"}
	found := false
	i := 0
	for j, tier := range arrayOfUserTier {
		if tier == userTier {
			found = true
			i = j
		}
	}
	if !found {
		return 0, fmt.Errorf("the userTier tier is not a valid user tier(didn't found it)")
	}

	allowedReqFromEnv := ""
	println(allowedReqFromEnv)
	switch i {
	case 0:
		// recurring
		reqPerDay := os.Getenv("ALLOWEDREQPERDAYFORRECURRINGTIER")
		if reqPerDay == "" {
			// hardcoded value is
			allowedReqFromEnv = "15"
		}
	case 1:
		// free tier
		reqPerDay := os.Getenv("ALLOWEDREQPERDAYFORFREETIER")
		if reqPerDay == "" {
			// hardcoded value is
			allowedReqFromEnv = "10"
		}
	case 2:
		// one time
		reqPerDay := os.Getenv("ALLOWEDREQPERDAYFORONETIER")
		if reqPerDay == "" {
			// hardcoded value is
			allowedReqFromEnv = "10"
		}
	default:
		return 0, fmt.Errorf("there are only 3 user tier and we got the value that is not in that , ranging over the array failed")
	}
	allowedReq, err := strconv.ParseInt(allowedReqFromEnv, 10, 64)
	if err != nil {
		println("there is a error in parsing the allowedReqFromEnv(which is either form the env or hardcoded) and it is ->", err.Error())
		return 0, err
	}
	println("the selected allowedReq for the user on tier: ", userTier, " is ->", allowedReq)

	return allowedReq, nil
}

// if we have a error (it will be form programmer mistake only, or the parse error) it will panic instead of error
//
// returns the int64 req per day the user on a tier is allowed to make
func (rl *RateLimiterForUser) getAllowedReqPerDayBasedOnUserTierPanic(userTier string) int64 {
	arrayOfUserTier := []string{"recurring", "free tier", "one time"}
	found := false
	i := 0
	for j, tier := range arrayOfUserTier {
		if tier == userTier {
			found = true
			i = j
		}
	}
	if !found {
		panic(fmt.Errorf(" --the userTier tier is not a valid user tier(didn't found it)--"))
	}

	allowedReqFromEnv := ""
	hardcodedReqPerDay := 3
	println(allowedReqFromEnv)
	switch i {
	case 0:
		// recurring
		reqPerDay := os.Getenv("ALLOWEDREQPERDAYFORRECURRINGTIER")
		if reqPerDay == "" {
			// hardcoded value is
			allowedReqFromEnv = "18"
			hardcodedReqPerDay = 18
		}
	case 1:
		// free tier
		reqPerDay := os.Getenv("ALLOWEDREQPERDAYFORFREETIER")
		if reqPerDay == "" {
			// hardcoded value is
			allowedReqFromEnv = "10"
			hardcodedReqPerDay = 10
		}
	case 2:
		// one time
		reqPerDay := os.Getenv("ALLOWEDREQPERDAYFORONETIER")
		if reqPerDay == "" {
			// hardcoded value is
			allowedReqFromEnv = "6"
			hardcodedReqPerDay = 6
		}
	default:
		panic(fmt.Errorf("there are only 3 user tier and we got the value that is not in that , ranging over the array failed"))
	}
	allowedReq, err := strconv.ParseInt(allowedReqFromEnv, 10, 64)
	if err != nil {
		println("there is a error in parsing the allowedReqFromEnv(which is either form the env or hardcoded) and it is ->", err.Error())
		println("we are returning the hardcoded req per day value and it is ->", hardcodedReqPerDay)
		return int64(hardcodedReqPerDay)
	}
	println("the selected allowedReq for the user on tier: ", userTier, " is ->", allowedReq)

	return allowedReq
}
