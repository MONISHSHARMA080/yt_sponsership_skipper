package askllm

import "fmt"

// we are not chekcing for that if the struct is empty or not
func (r ResponseForWhereToSkipVideo) GetStatusCode() int {
	return r.Status_code // Also return the actual status code, not just 1
}

func (r ResponseForWhereToSkipVideo) AreWeNotAllowedToReturnResponse() bool {
	return r.Message == "" && r.Status_code == 0 && r.StartTime == 0 && r.EndTime == 0
}

func (r *ResponseForWhereToSkipVideo) FillTheStructForError(errorMessageInTheResponse string, status_code int) {
	// this will fill the struct for the error
	r.Message = errorMessageInTheResponse
	r.Status_code = status_code
	r.StartTime = 0
	r.EndTime = 0
	r.ContainSponserSubtitle = false
	fmt.Printf("we have filled the struct for error  and it is -> %+v \n\n ", *r)
}

// the contains subtitle bool is assigned to be true cause duh !
func (r *ResponseForWhereToSkipVideo) FillTheStructForSuccess(MessageInTheResponse string, status_code int, startTime, endTime int64) {
	// this will fill the struct for the error
	r.Message = MessageInTheResponse
	r.Status_code = status_code
	r.StartTime = startTime
	r.EndTime = endTime
	r.ContainSponserSubtitle = true
	fmt.Printf("we have filled the struct for success and it is -> %+v  \n\n", *r)
}
