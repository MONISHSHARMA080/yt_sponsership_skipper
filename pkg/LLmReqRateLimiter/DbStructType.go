package llmreqratelimiter

import (
	"database/sql"
	"fmt"
	"time"
	genericresulttype "youtubeAdsSkipper/pkg/GenericResultType"
)

type RateLimiterForUser struct {
	// ID               int64     `db:"id"`
	UserEmail        string    `db:"user_email"`
	RequestTimestamp time.Time `db:"request_timestamp"`
}

// this only checks the email
func (rt *RateLimiterForUser) IsMyStructEmpty() bool {
	return rt.UserEmail == ""
}

// this function is to insert the user in the Db , ie. they are making a req and we want to count this
func (rt *RateLimiterForUser) NewRequestMadeUpdateDb(db *sql.DB, errAndResChannel chan genericresulttype.ErrorAndResultType[bool]) {
	if rt.IsMyStructEmpty() {
		errAndResChannel <- genericresulttype.ErrorAndResultType[bool]{Result: false, Err: fmt.Errorf("the struct is not initialized")}
		return
	}
	// currentTime := time.Now().Format("2006-01-02 15:04:05") // Standard Go time formatting
	sqlStatement := `
		INSERT INTO rate_limit_user(user_email, request_timestamp)
		VALUES (?, date('now'));
	`
	_, err := db.Exec(sqlStatement, rt.UserEmail)
	if err != nil {
		// Handle the error appropriately (log it, return it, etc.)
		errAndResChannel <- genericresulttype.ErrorAndResultType[bool]{Result: false, Err: fmt.Errorf("failed to log request for user %s: %w", rt.UserEmail, err)}
		return
	}
	errAndResChannel <- genericresulttype.ErrorAndResultType[bool]{Result: false, Err: nil}
}

func (rt *RateLimiterForUser) GetAllReqFromUserToday(db *sql.DB, errAndResChannel chan genericresulttype.ErrorAndResultType[int64]) {
	if rt.IsMyStructEmpty() {
		errAndResChannel <- genericresulttype.ErrorAndResultType[int64]{Result: 0, Err: fmt.Errorf("the struct is not initialized")}
		return
	}
	howManyReqUserHaveMade := `
    SELECT COUNT(*) 
		FROM rate_limit_user
		WHERE user_email = ?
		AND date(request_timestamp) = date('now')
  `
	var count int
	err := db.QueryRow(howManyReqUserHaveMade, rt.UserEmail).Scan(&count)
	if err != nil {
		errAndResChannel <- genericresulttype.ErrorAndResultType[int64]{Result: 0, Err: fmt.Errorf("failed to get request count for user %s: %w", rt.UserEmail, err)}
		return
	}
	fmt.Printf("all the user(email:%s) request in a day is -> %d \n\n", rt.UserEmail, count)
	errAndResChannel <- genericresulttype.ErrorAndResultType[int64]{Result: int64(count), Err: nil}
}

func (rt *RateLimiterForUser) ShouldWeRateLimitUser(db *sql.DB, userTier string, errAndResChannel chan genericresulttype.ErrorAndResultType[bool]) {
	if rt.IsMyStructEmpty() {
		errAndResChannel <- genericresulttype.ErrorAndResultType[bool]{Result: false, Err: fmt.Errorf("the struct is not initialized")}
		return
	}
	chanForAllUserReq := make(chan genericresulttype.ErrorAndResultType[int64])
	go rt.GetAllReqFromUserToday(db, chanForAllUserReq)
	allUserReqToday := <-chanForAllUserReq
	if allUserReqToday.Err != nil {
		errAndResChannel <- genericresulttype.ErrorAndResultType[bool]{Result: false, Err: allUserReqToday.Err}
		return
	}
	userReqAllowedPerDay := rt.getAllowedReqPerDayBasedOnUserTierPanic(userTier)
	println("the response came form the Db and the func and allowed req in a day (for the tier ", userTier, " ) and it ", userReqAllowedPerDay, " and the req made today by the user is ", allUserReqToday.Result)
	if userReqAllowedPerDay <= allUserReqToday.Result {
		errAndResChannel <- genericresulttype.ErrorAndResultType[bool]{Result: true, Err: nil}
		return
	}
	errAndResChannel <- genericresulttype.ErrorAndResultType[bool]{Result: false, Err: nil}
}
