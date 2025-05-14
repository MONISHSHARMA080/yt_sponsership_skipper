package askllm

import (
	"fmt"
	"net/http"
	commonresultchannel "youtubeAdsSkipper/pkg/askLLM/commonResultChannel"
	askllm "youtubeAdsSkipper/pkg/askLLM/groqHelper"
)

type HttpJsonResponder func(w http.ResponseWriter, message string, http_status_code int64)

// this fucntion will write to the htto response when we get the error , and also return a error just to verify that we have a error, and you can return , it the error in this
// is ==  nil then continue with yout work
func AskGroqAboutSponsorship(httpClient *http.Client, w http.ResponseWriter, method_to_write_http_and_json_to_respond HttpJsonResponder, apiKey string,
	result_for_subtitles askllm.String_and_error_channel_for_subtitles, ChanForResponseForGettingSubtitlesTiming chan askllm.ResponseForGettingSubtitlesTiming,
	resultChannel chan commonresultchannel.ResultAndErrorChannel[askllm.ResponseForWhereToSkipVideo],
) {
	response := commonresultchannel.ResultAndErrorChannel[askllm.ResponseForWhereToSkipVideo]{}

	channel_for_groqResponse := make(chan askllm.String_and_error_channel_for_groq_response)
	go askllm.AskGroqabouttheSponsorship(httpClient, channel_for_groqResponse, apiKey, &result_for_subtitles.String_value)
	groq_response := <-channel_for_groqResponse

	if groq_response.Err != nil && groq_response.GroqApiResponsePtr == nil || groq_response.SponsorshipContent == nil {
		if groq_response.Http_response_for_go_api_ptr != nil {
			println("the http response is not nil and the status code is ->", groq_response.Http_response_for_go_api_ptr.StatusCode)
			response.Result.FillTheStructForError("somethign went wrong on our side", http.StatusInternalServerError)
			response.Err = groq_response.Err
			resultChannel <- response
			return
		}
		if groq_response.Http_response_for_go_api_ptr.StatusCode == 429 {
			println("the response form the groq api is 429")
			response.Result.FillTheStructForError("the request time out on this tier", http.StatusTooManyRequests)
			response.Err = fmt.Errorf("the request time out on this tier (429)")
			resultChannel <- response
			return
		} else if groq_response.GroqApiResponsePtr == nil {
			println("groq error ", groq_response.Err.Error())
			response.Result.FillTheStructForError("somethign went wrong on our side", http.StatusInternalServerError)
			response.Err = groq_response.Err
			resultChannel <- response
			return
		}
		if groq_response.SponsorshipContent == nil {
			println(" sponsership content is not there ")
			response.Result.FillTheStructForError("somethign went wrong on our side", http.StatusInternalServerError)
			response.Err = groq_response.Err
			resultChannel <- response
			return
		}
	}
	println("the groq response is there ---")
	// getting error deciding the escaped json in the json response
	if groq_response.SponsorshipContent.DoesVideoHaveSponsorship && groq_response.SponsorshipContent.SponsorshipSubtitle != "" {
		// if the subtitles is found
		println("the sponsorship subtitles are found in the video")
		go askllm.GetTimeAndDurInTheSubtitles(result_for_subtitles.Transcript, &groq_response.SponsorshipContent.SponsorshipSubtitle, &result_for_subtitles.String_value, ChanForResponseForGettingSubtitlesTiming)

		SubtitlesTimming := <-ChanForResponseForGettingSubtitlesTiming
		println("we got the subtitle timing --")
		if SubtitlesTimming.Err != nil {
			if SubtitlesTimming.Err.Error() == "" {
				response.Result.FillTheStructForError("no subtitle found for the video", http.StatusNotFound)
				response.Err = fmt.Errorf("no subtitle found for the video")
				resultChannel <- response
				return
			}
			println("error in result_for_subtitles.err --> ", SubtitlesTimming.Err.Error())
			println(" we got the subtitles from the groq but the subtitles can't be found by out function")
			response.Result.FillTheStructForError("Something is wrong on our side, error getting subtitles timming ", http.StatusInternalServerError)
			response.Err = fmt.Errorf(" we got subtitle form groq but have problem finding it via our function and the error is :-> %s ", SubtitlesTimming.Err.Error())
			resultChannel <- response
			return
		} else if SubtitlesTimming.EndTime+SubtitlesTimming.StartTime <= 0 {
			// the subtitle is not being found via the function or somethign is wrong
			println("\n\n  ----- the subtitle is not gettign found despite of the llm telling us that it is there ------   \n\n")
			response.Result.FillTheStructForError("Something is wrong on our side, error getting subtitles timming is not found of the subtitle ", http.StatusInternalServerError)
			response.Err = fmt.Errorf("the subtitle is not gettign found despite of the llm telling us that it is there ")
			resultChannel <- response
			return
		}
		response.Result.FillTheStructForSuccess(" subtitle found for the video ", http.StatusOK, int64(SubtitlesTimming.StartTime), int64(SubtitlesTimming.EndTime), groq_response.SponsorshipContent.DoesVideoHaveSponsorship)
		response.Err = nil
		resultChannel <- response
		return
	} else {
		// return no to the user as either the video does not have subtitles or the groq did not return it
		println("the video does not have sponsorship subtitles, or it is not found")
		response.Result.FillTheStructForError("sponsership subtitles not found in the video", http.StatusOK)
		// nil as there is not a error as the video does not have error
		response.Err = nil
		resultChannel <- response
		return
	}
}
