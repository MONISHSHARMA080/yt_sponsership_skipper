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
	sponsershipSubtitlesEndIndex := sponsershipSubtitlesStartIndex + len(*sponsership_subtitles_form_groq)
	sponsershipLengthTracker := 0
	sponsershipStartSubtitleIndex := 0
	sponsershipEndSubtitleIndex := 0
	lengthOfTranscriptSubtitle := len(transcripts.Subtitles)
	startingSubtitleFound := false
	endingSubtitleFound := false
	// subtitleFoundInForstString:= false // not using as just check if the sponsershipStartSubtitleIndex!=0
	for i, subtitle := range transcripts.Subtitles {
		sponsershipLengthTracker += len(subtitle.Text + " ") // cause this does not include space
		if sponsershipLengthTracker >= sponsershipSubtitlesStartIndex {
			startingSubtitleFound = true
			if sponsershipLengthTracker > sponsershipStartSubtitleIndex {
				// we can say that we have came overBoard(cause sponsershipLengthTracker > sponsershipStartSubtitleIndex ) so if we go back it should not be that far off
				// make sure i is not 0(starting case, it can)
				// what I am going to do it get the first 2 words of the sponserSub and search it in this-1 it !found then  2nd if not !found still then search the first in the
				// or calculate the len of 1st 2 words of caption sub it it is in the index of
				if i > 0 {
					println("got overboard, the text in prev one is(I think correct)-->", transcripts.Subtitles[i-1].Text)
					println("this text is -->", subtitle.Text)
					// now the subtitle is either here or in prev subtitle
					sponsershipStartSubtitleIndex = getIndexOfSponserSubtitleForEndFromAdjacentIndex(*transcripts, i, sponsership_subtitles_form_groq, true, lengthOfTranscriptSubtitle)
					if sponsershipStartSubtitleIndex == i-1 {
						sponsershipLengthTracker = sponsershipLengthTracker - len(subtitle.Text+" ") // over counted
					} // else we are just on track
					println("the correct text is -->", transcripts.Subtitles[sponsershipStartSubtitleIndex].Text)
				} else {
					sponsershipStartSubtitleIndex = i
				}
			} else {
				sponsershipStartSubtitleIndex = i
			}
			break
		}
	}
	for i := sponsershipStartSubtitleIndex + 1; i < lengthOfTranscriptSubtitle; i++ {
		sponsershipLengthTracker += len(transcripts.Subtitles[i].Text + " ")
		if sponsershipLengthTracker >= sponsershipSubtitlesEndIndex {
			endingSubtitleFound = true
			if sponsershipLengthTracker > sponsershipSubtitlesEndIndex {
				// if i>0 will result in true
				// i+1 cause this checks on the prev and current one(yes it is a hack)
				sponsershipEndSubtitleIndex = getIndexOfSponserSubtitleForEndFromAdjacentIndex(*transcripts, i, sponsership_subtitles_form_groq, false, lengthOfTranscriptSubtitle)
				// if sponsershipStartSubtitleIndex == i-1 { // do I need to do it as it is unnecessary (in the end)
				// 	sponsershipLengthTracker = sponsershipLengthTracker - len(transcripts.Subtitles[i].Text+" ") // over counted
				// }
			}
		}
	}
	println("correct ending text--", transcripts.Subtitles[sponsershipEndSubtitleIndex].Text)

	println("sponsershipSubtitlesStartIndex, sponsershipSubtitlesEndIndex, sponsershipSubtitlesEndIndex--", sponsershipSubtitlesStartIndex, sponsershipSubtitlesEndIndex, sponsershipEndSubtitleIndex)

	if startingSubtitleFound && endingSubtitleFound {
		// return the time(implement the fucntion)
		startTimeOfSubtitle, err := getTimeAndDurFromSubtitles(transcripts, sponsershipStartSubtitleIndex)
		if err != nil {
			responseForTimmingChannel <- ResponseForGettingSubtitlesTiming{0, 0, err}
			return
		}
		endTimeOfSubtitle, err := getTimeAndDurFromSubtitles(transcripts, sponsershipEndSubtitleIndex)
		if err != nil {
			responseForTimmingChannel <- ResponseForGettingSubtitlesTiming{0, 0, err}
			return
		}
		responseForTimmingChannel <- ResponseForGettingSubtitlesTiming{int(startTimeOfSubtitle), int(endTimeOfSubtitle), nil}
		return

	} else {
		// return o
		responseForTimmingChannel <- ResponseForGettingSubtitlesTiming{0, 0, fmt.Errorf("Starting or ending subtitle was not found")}
		return
	}
}

func getTimeAndDurFromSubtitles(t *Transcripts, subtitleIndex int) (float64, error) {
	startTime, err := strconv.ParseFloat(t.Subtitles[subtitleIndex].Start, 64)
	if err != nil {
		return 0, err
	}
	return startTime, nil
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
	// Trim any leading whitespace
	JSONBYGroq = strings.TrimSpace(JSONBYGroq)

	// Find the first '{' and last '}'
	firstBrace := strings.Index(JSONBYGroq, "{")
	lastBrace := strings.LastIndex(JSONBYGroq, "}")

	// Extract the JSON substring
	if firstBrace == -1 || lastBrace == -1 {
		return JSONBYGroq
	}

	trimmedJSON := JSONBYGroq[firstBrace : lastBrace+1]

	// Replace "true" with true, but only the first occurrence
	formattedJSON := strings.Replace(trimmedJSON, `"true"`, "true", 1)

	return formattedJSON
}

func getIndexOfSponserSubtitleFromAdjacentIndex2(transcript Transcripts, currentIndex int, subtitleCaptionByGroq *string, returnFirst bool) int {
	firstWord, secondWord, err := getFirstTwoWords(subtitleCaptionByGroq)
	if err != nil {
		// If only one character, return the current index
		return currentIndex
	}

	// Search indices to check (previous, current, next)
	searchIndices := []int{currentIndex - 1, currentIndex, currentIndex + 1}

	// To store matching indices
	matchingIndices := []int{}

	for _, idx := range searchIndices {
		// Ensure index is within slice bounds
		if idx >= 0 && idx < len(transcript.Subtitles) {
			subtitleText := strings.ToLower(transcript.Subtitles[idx].Text)
			// Check if both first and second words are present
			if strings.Contains(subtitleText, strings.ToLower(firstWord)) &&
				strings.Contains(subtitleText, strings.ToLower(secondWord)) {
				matchingIndices = append(matchingIndices, idx)
			}
		}
	}

	// If no matching indices found, return current index
	if len(matchingIndices) == 0 {
		println("can't find so returning the basic")
		return currentIndex
	}

	// Determine which index to return based on the flag
	if returnFirst {
		// Return the first (lowest) matching index
		return matchingIndices[0]
	} else {
		// Return the last (highest) matching index
		return matchingIndices[len(matchingIndices)-1]
	}
}
func getIndexOfSponserSubtitleForEndFromAdjacentIndex(transcript Transcripts, currentIndex int, subtitleCaptionByGroq *string, returnBasedOnFirstWord bool, lengthOfSubtitle int) int {
	firstWord, secondWord, err := getFirstTwoWords(subtitleCaptionByGroq)
	if err != nil {
		// probally only 1 character  <<- jsut gonna return the current index
		return currentIndex
	}
	firstWordPresentInCurrentIndex := strings.Contains(strings.ToLower(transcript.Subtitles[currentIndex].Text), strings.ToLower(firstWord))
	secondWordPresentInCurrentIndex := strings.Contains(strings.ToLower(transcript.Subtitles[currentIndex].Text), strings.ToLower(secondWord))
	if returnBasedOnFirstWord {
		if currentIndex-1 > 0 {
			return currentIndex
		}
		firstWordPresentInPrevIndex := strings.Contains(strings.ToLower(transcript.Subtitles[currentIndex-1].Text), strings.ToLower(firstWord))
		secondWordPresentInPrevIndex := strings.Contains(strings.ToLower(transcript.Subtitles[currentIndex-1].Text), strings.ToLower(secondWord))
		if firstWordPresentInCurrentIndex && !secondWordPresentInCurrentIndex || firstWordPresentInCurrentIndex && secondWordPresentInCurrentIndex {
			return currentIndex
		} else {
			println("assesting -->firstWordPresentInPrevIndex && secondWordPresentInPrevIndex", firstWordPresentInPrevIndex && secondWordPresentInPrevIndex == true)
			return currentIndex - 1
		}
	} else {
		if currentIndex+1 >= lengthOfSubtitle {
			return currentIndex // safe bet
		}
		firstWordPresentInNextIndex := strings.Contains(strings.ToLower(transcript.Subtitles[currentIndex+1].Text), strings.ToLower(firstWord))
		secondWordPresentInNextIndex := strings.Contains(strings.ToLower(transcript.Subtitles[currentIndex+1].Text), strings.ToLower(secondWord))
		if secondWordPresentInNextIndex && !firstWordPresentInNextIndex || firstWordPresentInNextIndex && secondWordPresentInNextIndex {
			return currentIndex + 1
		} else {
			println("assesting -->firstWordPresentInNextIndex && secondWordPresentInNextIndex", firstWordPresentInNextIndex && secondWordPresentInNextIndex == true)
			return currentIndex
		}
	}

	// for _, idx := range []int{currentIndex - 1, currentIndex, currentIndex + 1} {
	// 	// Ensure index is within slice bounds
	// 	if idx >= 0 && idx < len(transcript.Subtitles) {
	// 		subtitleText := strings.ToLower(transcript.Subtitles[idx].Text)
	// 		// Check if both first and second words are present
	// 		if strings.Contains(subtitleText, strings.ToLower(firstWord)) &&
	// 			strings.Contains(subtitleText, strings.ToLower(secondWord)) {
	// 			return idx
	// 		}
	// 	}
	// }
	// println("can't find so returning the basic")
	// return currentIndex // as a safe bet
}

func getFirstTwoWords(stringToPerformOperationOn *string) (string, string, error) {
	strToReturn := strings.Fields(*stringToPerformOperationOn)
	if len(strToReturn) < 2 {
		println("+++", strToReturn)
		return "", "", fmt.Errorf("the string has less than 2 words")
	}
	println("==", strToReturn[0])
	return strToReturn[0], strToReturn[1], nil
}
