package main

import "net/http"

// contains structs and interfaces
type String_and_error_channel_for_groq_response struct {
	err                error
	groqApiResponsePtr *GroqApiResponse
	http_response_for_go_api_ptr *http.Response
}

func (groqResp *String_and_error_channel_for_groq_response) checkIfSponsorshipSubtitlesExist (transcript *Transcripts) (bool,int64,error){
	//will check the groq response and if the ads is present (in the groq resp bool) then I will go through the transcript to find it and return the time to the user
	return false, 0, nil
}