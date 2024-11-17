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
	"text/template/parse"

	"gotest.tools/assert"
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

	println("\n=== START OF RESPONSE BODY ===")
	println(string(bodyBytes))
	println("=== END OF RESPONSE BODY ===\n")

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
		println("\nContent field:", groqApiResponse.Choices[0].Message.Content, "\n\n")

		var sponsorshipContent SponsorshipContent
		err = json.Unmarshal([]byte(groqApiResponse.Choices[0].Message.Content), &sponsorshipContent)
		if err != nil {
			println("||5 - Error parsing content JSON:", err.Error())
			println("Content was:", groqApiResponse.Choices[0].Message.Content)
			channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: http_response, SponsorshipContent: nil}
			return
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

	payload := map[string]interface{}{
		"model": os.Getenv("GROQ_MODEL"),
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": "don't forget I only need json form you nothing else; sutitles-->" + *subtitlesInTheVideo,
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

func GetTimeAndDurInTheSubtitles(transcripts *Transcripts, sponsership_subtitles_form_groq *string, full_captions *string) (int, int, error) {
	// should probally use/communicate with goroutine

	// do some sort of two pointer  to pruse through
	//  take the fist string of the comparision text and then search for the string in the transcript and at a match check the next word if the next-to-next is
	//  not there then moove on to the words until you find the same patter again

	// -------- or ----------

	// make a string with string builder and then then comapre it

	// -------- or ----------

	// get the long builder string as a ptr to check of the string is there if it si then we will get where thhe sponser ship subtitle is in the
	// total caption text, (assume it in in the last ) we will itereate through the last subtitles[] and get to the text before it (comparing it to the ) (subtitle[] and stringFromGroq as a two ptr technioque)

	// -------- or -------

	// search for the subtitles in the caption string and at what number of char it is there , and then itereate overt the transcript until we get the desired
	// length/index of the string



	// the len is 1 indexed and the index is  0 so take it out
	length_of_full_captions := len(*full_captions)
	length_of_subtitles := len(transcripts.Subtitles)
	//cause if the index is there it contian if too and no need to compare the too
	sponsership_subtitles_index := strings.Index(strings.ToLower(*full_captions), strings.ToLower(*sponsership_subtitles_form_groq))
	if sponsership_subtitles_index == -1 {
		return 0,0,fmt.Errorf("error getting the substring position in the string")
	}
	tracker_for_len_of_sub_in_transcript := 0
	tracker_for_len_of_sub_in_transcript_prev_value := 0
	var timeAndDurationFromStartSub TimeAndDurationFromSub
	var whereTheIndexForSubtitlesWas int
	// make this one a seperate function to test
	for i := 0; i < length_of_subtitles; i++ {
		// here iterate over the full subtitles, in individual subtitle get the length -1 of that and compare it to index of subtitles, if it is less than add it 
		// to a tracker var and keep going until it is    
		tracker_for_len_of_sub_in_transcript_prev_value = tracker_for_len_of_sub_in_transcript	
		 tracker_for_len_of_sub_in_transcript = len(transcripts.Subtitles[i].Text) -1
		 if tracker_for_len_of_sub_in_transcript >=sponsership_subtitles_index{
			// subtitles in sponsership starts form here
			whereTheIndexForSubtitlesWas = i
			// what to do 1) get where it starts  and the duration and divide duration by the words there and return the startTime + dur and 
			// call that field an estimated field and a raw start field
			// 2) return the same (above func) but for the end and also how do we determine where it ends --> keep some sort of tracker that has the 
			// len of subtitle  string there and itereate over the subtitles until you react that length 
			timeAndDurationFromStartSub = getTimeAndDurFromSubtitles(&transcripts.Subtitles[i].Text, transcripts.Subtitles[i].Dur, transcripts.Subtitles[i].Start, tracker_for_len_of_sub_in_transcript_prev_value, sponsership_subtitles_index )
			if timeAndDurationFromStartSub.err != nil {
				// handel error 
				return 0,0, timeAndDurationFromStartSub.err
			}
			

			println("the line in subtiltes is -->", transcripts.Subtitles[i].Text)
		 }
	}
	// now for the end one's for loop
	// this one is till the end of the 
	index_of_end_substring := sponsership_subtitles_index + len(*sponsership_subtitles_form_groq)

	// -1 as len is 1 indexed  
	for i := whereTheIndexForSubtitlesWas - length_of_subtitles -1 ; i < length_of_subtitles; i++ {
		
	}
	
	return 0, 0, nil
}

type TimeAndDurationFromSub struct{
	err error
	estimated_start_time int64
	start_time int64
}

func getTimeAndDurFromSubtitles(subtitles *string, dur string, start string, prev_value_of_transcript_tracker int, index_of_substring int)TimeAndDurationFromSub{
	
	//  will be here if the subtiutle string is in the current index (or nearby i.e in the string), this function will go to the index of the substring and return
	//  the estimeated start time of the sponsership and
	
	duration, err :=strconv.ParseInt(dur, 10, 64)

	if  err!= nil {
		return TimeAndDurationFromSub{err: err, estimated_start_time:0, start_time: 0}
	}
	
	startTime, err :=strconv.ParseInt(start, 10, 64)
	
	if  err!= nil {
		return TimeAndDurationFromSub{err: err, estimated_start_time:0, start_time: 0}
	}
	
	println("asserting in getTimeAndDurFromSubtitles() prev_value_of_transcript_tracker is less than index_of_substring",index_of_substring >= prev_value_of_transcript_tracker)
	println("prev_value_of_transcript_tracker is ",prev_value_of_transcript_tracker ," and index_of_substring is ", index_of_substring)

	var length_of_subtitles int = len(*subtitles)
	var value_to_increment_index_by int = index_of_substring - prev_value_of_transcript_tracker
	
	if length_of_subtitles < value_to_increment_index_by{
		// handle error 
		return TimeAndDurationFromSub{err: err, estimated_start_time:0, start_time: 0}
	}
	
	prev_value_of_transcript_tracker = value_to_increment_index_by
	
	//  calc the time taken in speaking the whole by subtitle/caption[n] / by len(length_of_subtitles), now we got the w.p. unit and now will get the words 
	//  per unit of the bytes in the beginning by wpU * value of bytes before and then we remove it form dur  and get the estimated time to skip 
	
	estamited_time_to_skip:=  duration - ( duration/int64(length_of_subtitles) * int64(value_to_increment_index_by) )
	println("estimated time to skip is less than duration --> ", estamited_time_to_skip <= duration)

	return TimeAndDurationFromSub{err: nil, estimated_start_time:estamited_time_to_skip + startTime, start_time: startTime}
}