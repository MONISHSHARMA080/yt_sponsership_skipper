package askllm

// we are not chekcing for that if the struct is empty or not
func (r ResponseForWhereToSkipVideo) GetStatusCode() int {
    return r.Status_code  // Also return the actual status code, not just 1
}

func (r ResponseForWhereToSkipVideo) AreWeNotAllowedToReturnResponse() bool {
   return r.Message == "" && r.Status_code == 0 && r.StartTime == 0 && r.EndTime == 0
}

func (r *ResponseForWhereToSkipVideo) FillTheStructForError(message string, status_code int )  {
    // this will fill the struct for the error
    r.Message = message
    r.Status_code = status_code
    r.StartTime = 0
    r.EndTime = 0
    r.ContainSponserSubtitle = false
}

