package askllm

// this will make the import hard to decide where they are comming form
// import . "youtubeAdsSkipper/pkg/askLLM/groqHelper"

import (
	"encoding/json"
	"net/http"
	askllm "youtubeAdsSkipper/pkg/askLLM/groqHelper"
)

type HttpJsonResponder func(w http.ResponseWriter, message string, http_status_code int)

// this fucntion will write to the htto response when we get the error , and also return a error just to verify that we have a error, and you can return , it the error in this
// is ==  nil then continue with yout work
func AskGroqAboutSponsorship(httpClient *http.Client, w http.ResponseWriter, method_to_write_http_and_json_to_respond HttpJsonResponder, apiKey string,
	result_for_subtitles askllm.String_and_error_channel_for_subtitles, ChanForResponseForGettingSubtitlesTiming chan askllm.ResponseForGettingSubtitlesTiming,
) {
	channel_for_groqResponse := make(chan askllm.String_and_error_channel_for_groq_response)
	go askllm.AskGroqabouttheSponsorship(httpClient, channel_for_groqResponse, apiKey, &result_for_subtitles.String_value)
	groq_response := <-channel_for_groqResponse

	if groq_response.Err != nil && groq_response.GroqApiResponsePtr == nil || groq_response.SponsorshipContent == nil {
		if groq_response.Http_response_for_go_api_ptr != nil {
			println("the http response is not nil and the status code is ->", groq_response.Http_response_for_go_api_ptr.StatusCode)
			method_to_write_http_and_json_to_respond(w, "somethign went wrong on our side", http.StatusInternalServerError)
			return
		}
		if groq_response.Http_response_for_go_api_ptr.StatusCode == 429 {
			println("the response form the groq api is 429")
			method_to_write_http_and_json_to_respond(w, "the request time out on this tier", http.StatusTooManyRequests)
			return
		} else if groq_response.GroqApiResponsePtr == nil {
			println("groq error ", groq_response.Err.Error())
			method_to_write_http_and_json_to_respond(w, "somethign went wrong on our side", http.StatusInternalServerError)
			return
		}
		if groq_response.SponsorshipContent == nil {
			println(" sponsership content is not there ")
			method_to_write_http_and_json_to_respond(w, "somethign went wrong on our side", http.StatusInternalServerError)
			return
		}
	}
	// getting error deciding the escaped json in the json response
	if groq_response.SponsorshipContent.DoesVideoHaveSponsorship && groq_response.SponsorshipContent.SponsorshipSubtitle != "" {
		// if the subtitles is found
		go askllm.GetTimeAndDurInTheSubtitles(result_for_subtitles.Transcript, &groq_response.SponsorshipContent.SponsorshipSubtitle, &result_for_subtitles.String_value, ChanForResponseForGettingSubtitlesTiming)
	} else {
		// return no to the user as either the video does not have subtitles or the groq did not return it
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(askllm.ResponseForWhereToSkipVideo{Message: "sponsership subtitles not found in the video", Status_code: http.StatusOK, ContainSponserSubtitle: false})
		if err != nil {
			println("error in the method  encoding the json in the struct 123-->", err.Error())
		}
		return
	}
}
