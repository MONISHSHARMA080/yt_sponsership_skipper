package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
	// "strings"
	// "text/template/parse"
	// "gotest.tools/assert"
)

func AskGroqabouttheSponsorship(httpClient *http.Client, channel_for_groq_response chan<- String_and_error_channel_for_groq_response, APIKEY_according_to_users_tier string, subtitlesInTheVideo *string) {
	err, http_req := factoryGroqPostReqCreator(APIKEY_according_to_users_tier, subtitlesInTheVideo)
	if err != nil {
		println("||1")
		channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: nil, SponsorshipContent: nil}
		return
	}

	http_response, err := httpClient.Do(http_req)
	if err != nil {
		println("||2")
		channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: http_response, SponsorshipContent: nil}
		return
	}
	defer http_response.Body.Close()

	println(http_response != nil, "http_response is not null")

	// Read and print the response body
	bodyBytes, err := io.ReadAll(http_response.Body)
	if err != nil {
		println("||3 - Error reading response body:", err.Error())
		channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: http_response, SponsorshipContent: nil}
		return
	}

	bodyBytes = []byte(string(bodyBytes))
	// println("\n=== START OF RESPONSE BODY ===")
	// println(string(bodyBytes))
	// println("=== END OF RESPONSE BODY ===\n")

	// Create a new reader from the body bytes for json.NewDecoder
	var groqApiResponse GroqApiResponse
	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&groqApiResponse)
	if err != nil {
		println("||4 - Error decoding response:", err.Error())
		println("Response body was:", string(bodyBytes))
		channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: http_response, SponsorshipContent: nil}
		return
	}

	if len(groqApiResponse.Choices) > 0 {
		// println("\nContent field:", groqApiResponse.Choices[0].Message.Content, "\n\n")
		// formatting json
		a := formatGroqJson(groqApiResponse.Choices[0].Message.Content)
		if a == "" {
			println(" can't format the groq json (exiting)-->", a)
			channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: http_response, SponsorshipContent: nil}
			return
		}
		println("formatted groq json is -->\n", a, "\n")
		groqApiResponse.Choices[0].Message.Content = a
		var sponsorshipContent SponsorshipContent
		err = json.Unmarshal([]byte(groqApiResponse.Choices[0].Message.Content), &sponsorshipContent)
		if err != nil {
			println("||5 - Error parsing content JSON:", err.Error())
			println("Content was:", groqApiResponse.Choices[0].Message.Content)
			channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: http_response, SponsorshipContent: nil}
			return
		}
		// converting it to utf-8 (string is not working)

		if !utf8.ValidString(sponsorshipContent.SponsorshipSubtitle) {
			println(" the strign is not valid utf-8")
			sponsorshipContent.SponsorshipSubtitle = string(sponsorshipContent.SponsorshipSubtitle)
			if !utf8.ValidString(sponsorshipContent.SponsorshipSubtitle) {
				println("even after the strign() still not valid utf-8")
			}
		}
		channel_for_groq_response <- String_and_error_channel_for_groq_response{err: nil, groqApiResponsePtr: &groqApiResponse, http_response_for_go_api_ptr: http_response, SponsorshipContent: &sponsorshipContent}
	} else {
		println("No choices in response")
		channel_for_groq_response <- String_and_error_channel_for_groq_response{err: fmt.Errorf("no choices presesnt in the gorq response"), groqApiResponsePtr: &groqApiResponse, http_response_for_go_api_ptr: http_response, SponsorshipContent: nil}
	}

	println("||6")
}

func factoryGroqPostReqCreator(GroqApiKey string, subtitlesInTheVideo *string) (error, *http.Request) {
	// stringify a json schema in the end

	url := "https://api.groq.com/openai/v1/chat/completions"
	println(os.Getenv("GROQ_MESSAGE_CONTENT"))
	payload := map[string]interface{}{
		"model": os.Getenv("GROQ_MODEL"),
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": "give me all the full text from sponsership section. If the sponser is mention multiple time give me only one form the sponsership segment. Only Valid JSON is allowed   -->" + *subtitlesInTheVideo,
			},
			{
				"role":    "system",
				"content": os.Getenv("GROQ_MESSAGE_CONTENT"),
			},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err, nil
	}

	createdHttpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err, nil
	}
	createdHttpReq.Header.Set("Content-Type", "application/json")
	createdHttpReq.Header.Set("Authorization", "Bearer "+GroqApiKey)

	return nil, createdHttpReq
}

func getAPIKEYForGroqBasedOnUsersTeir(is_user_paid bool) (string, error) {
	// the env should have 0 to the numebr of keys in front of them
	// 2 api keys for the groq ; the errors are from parsing in to str
	var err error
	var number_ofKeys int64

	if is_user_paid == true {
		numberOFEnvKeyAccordignToUserTeir := os.Getenv("NO_OF_KEYS_FOR_PAID_USER")
		number_ofKeys, err = strconv.ParseInt(numberOFEnvKeyAccordignToUserTeir, 10, 32)
	} else {
		numberOFEnvKeyAccordignToUserTeir := os.Getenv("NO_OF_KEYS_FOR_UNPAID_USER")
		number_ofKeys, err = strconv.ParseInt(numberOFEnvKeyAccordignToUserTeir, 10, 32)
	}
	if err != nil {
		return "", err
	}
	random_number_for_apiKey := rand.Intn(int(number_ofKeys))
	println("the nummber of key in the env is --> ", number_ofKeys)
	if is_user_paid {
		a := strconv.Itoa(random_number_for_apiKey)
		println("random number generated is ->", a)
		return os.Getenv("API_KEY_PAID" + a), nil
	} else {
		a := strconv.Itoa(random_number_for_apiKey)
		println("random number generated is ->", a)
		return os.Getenv("API_KEY_UNPAID" + a), nil
	}
}

type ResponseForGettingSubtitlesTiming struct {
	startTime int
	endTime   int
	err       error
}

func GetTimeAndDurInTheSubtitles(transcripts *Transcripts, sponsership_subtitles_form_groq *string, full_captions *string, responseForTimmingChannel chan<- ResponseForGettingSubtitlesTiming) {
	// this function will return the timming of the sub form the args
	//
	// reach the first position and get the timeandDur and then the last and ...
	//
	//  just reach the end and  start from indexes and just calculate the dur, start
	// mean just go to that string location and get the dur and start
	//
	sponsershipSubtitlesStartIndex := strings.Index(strings.ToLower(*full_captions), strings.ToLower(*sponsership_subtitles_form_groq))
	if sponsershipSubtitlesStartIndex == -1 {
		println("subtitle is not there")
		responseForTimmingChannel <- ResponseForGettingSubtitlesTiming{0, 0, fmt.Errorf("can't find the subtitles")}
		return // string is not present
	}
	// make a new branch and try 2 pointer approach(splip on the " ")
	// or normilization in both the strings
	// compare both strings to see if they are eqaul or not(before any of above stradegy)
	sponsershipSubtitlesEndIndex := sponsershipSubtitlesStartIndex + len(*sponsership_subtitles_form_groq)
	sponsershipLengthTracker := 0
	sponsershipStartSubtitleIndex := 0
	for i, subtitle := range transcripts.Subtitles {
		sponsershipLengthTracker += len(subtitle.Text + " ") // cause this does not include space
		if sponsershipLengthTracker >= sponsershipSubtitlesStartIndex {
			// need a better way to reach the strign than this
			sponsershipStartSubtitleIndex = i - 3
			println("+++---", transcripts.Subtitles[i-1].Text)
			println("+++---", subtitle.Text)
			println("+++---", transcripts.Subtitles[i+1].Text)
			println(transcripts.Subtitles[sponsershipStartSubtitleIndex].Text, "  --at ", sponsershipStartSubtitleIndex, " and sponsershipSubtitlesStartIndex is ", sponsershipSubtitlesStartIndex, " sponsershipLengthTracker is ", sponsershipLengthTracker)
			println("ads from full_caption", (*full_captions)[sponsershipSubtitlesStartIndex:sponsershipSubtitlesEndIndex])
			break
		}
	}
	// after the for loop go for the last for loop (make it i-3 as well (see first))
	for i := 0; i < len(transcripts.Subtitles); i++ {
		// try utf8.RuneCountInString() as a last resort
	}
	println("sponsershipSubtitlesStartIndex, sponsershipSubtitlesEndIndex--", sponsershipSubtitlesStartIndex, sponsershipSubtitlesEndIndex)

	// a := strings.Compare(makeStringFromASubtitle(transcripts.Subtitles), *full_captions)
	// println("are both the strings equal -->", a)
	responseForTimmingChannel <- ResponseForGettingSubtitlesTiming{0, 0, fmt.Errorf("can't find the subtitles")}
	return
}

type TimeAndDurationFromSub struct {
	err                  error
	estimated_start_time float64
	start_time           float64
}

func getTimeAndDurFromSubtitles(subtitles *string, dur string, start string, prev_value_of_transcript_tracker int, index_of_substring int) TimeAndDurationFromSub {

	//  will be here if the subtiutle string is in the current index (or nearby i.e in the string), this function will go to the index of the substring and return
	//  the estimeated start time of the sponsership and

	duration, err := strconv.ParseFloat(dur, 10)

	if err != nil {
		return TimeAndDurationFromSub{err: err, estimated_start_time: 0, start_time: 0}
	}

	startTime, err := strconv.ParseFloat(start, 10)

	if err != nil {
		return TimeAndDurationFromSub{err: err, estimated_start_time: 0, start_time: 0}
	}

	println("asserting in getTimeAndDurFromSubtitles() prev_value_of_transcript_tracker is less than index_of_substring", index_of_substring >= prev_value_of_transcript_tracker)
	println("prev_value_of_transcript_tracker is ", prev_value_of_transcript_tracker, " and index_of_substring is ", index_of_substring)

	var length_of_subtitles int = len(*subtitles)
	var value_to_increment_index_by int = index_of_substring - prev_value_of_transcript_tracker

	if length_of_subtitles < value_to_increment_index_by {
		// handle error
		return TimeAndDurationFromSub{err: err, estimated_start_time: 0, start_time: 0}
	}

	prev_value_of_transcript_tracker = value_to_increment_index_by

	//  calc the time taken in speaking the whole by subtitle/caption[n] / by len(length_of_subtitles), now we got the w.p. unit and now will get the words
	//  per unit of the bytes in the beginning by wpU * value of bytes before and then we remove it form dur  and get the estimated time to skip

	estamited_time_to_skip := duration - (duration / float64(length_of_subtitles) * float64(value_to_increment_index_by))
	println("estimated time to skip is less than duration --> ", estamited_time_to_skip <= duration)

	return TimeAndDurationFromSub{err: nil, estimated_start_time: estamited_time_to_skip + startTime, start_time: startTime}
}

func parseXMLAndAdd(subtitles []Subtitle) string {
	var result strings.Builder
	for _, subtitle := range subtitles {
		result.WriteString(fmt.Sprintf("[%s] %s [%s] ",
			subtitle.Start,
			strings.TrimSpace(subtitle.Text),
			subtitle.Dur))
	}
	return result.String()
}

func makeStringFromASubtitle(a []Subtitle) string {
	var resultStr strings.Builder
	for _, subtitle := range a {
		i, err := resultStr.WriteString(subtitle.Text + " ")
		if err != nil {
			println("\nerror in building the string in iter..", i, " form subtitle-->", err.Error())
		}
	}
	return resultStr.String()
}

func formatGroqJson(JSONBYGroq string) string {
	// why cause ti hallucinates on json and sometimes write ```json{....}```
	startIndex := strings.Index(JSONBYGroq, "{")
	if startIndex == -1 {
		return ""
	}

	// Find the last occurrence of }
	endIndex := strings.LastIndex(JSONBYGroq, "}")
	if endIndex == -1 || endIndex <= startIndex {
		return ""
	}

	return strings.Replace(JSONBYGroq, `"true"`, "true", 1)[startIndex : endIndex+1]

	// // Extract the substring
	// return JSONBYGroq[startIndex : endIndex+1]
}

// func GetTimeAndDurInTheSubtitles2(transcripts *Transcripts, sponsership_subtitles_form_groq *string, full_captions *string, responseForTimmingChannel chan<- ResponseForGettingSubtitlesTiming) {
// 	// write it in a single for loop
// 	sponserShipSubtitlesFromGroqLowerCase := strings.ToLower(*sponsership_subtitles_form_groq)
// 	for i, subtitle := range transcripts.Subtitles {

// 		if strings.Index(sponserShipSubtitlesFromGroqLowerCase, strings.ToLower(subtitle.Text)) != -1 {
// 			// we found our first string
// 			getTimeAndDurFromSubtitles()
// 		} else {
// 			// subtitle is not in transcript
// 			continue
// 		}

// 	}
// }
