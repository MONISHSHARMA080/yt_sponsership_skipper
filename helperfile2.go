package main

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)



func AskGroqabouttheSponsorship(httpClient *http.Client, channel_for_groq_response chan<- String_and_error_channel_for_groq_response, APIKEY_according_to_users_tier string, subtitlesInTheVideo *string) {
    err, http_req := factoryGroqPostReqCreator(APIKEY_according_to_users_tier, subtitlesInTheVideo)
    if err != nil {
        println("||1")
        channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: nil}
        return
    }

    http_response, err := httpClient.Do(http_req)
    if err != nil {
        println("||2")
        channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: nil}
        return
    }
    defer http_response.Body.Close()

    println(http_response != nil, "http_response is not null")

    // Read and print the response body
    bodyBytes, err := io.ReadAll(http_response.Body)
    if err != nil {
        println("||3 - Error reading response body:", err.Error())
        channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: http_response}
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
        channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: http_response}
        return
    }

    if len(groqApiResponse.Choices) > 0 {
        println("\nContent field:", groqApiResponse.Choices[0].Message.Content, "\n\n")
        
        var sponsorshipContent SponsorshipContent
        err = json.Unmarshal([]byte(groqApiResponse.Choices[0].Message.Content), &sponsorshipContent)
        if err != nil {
            println("||5 - Error parsing content JSON:", err.Error())
            println("Content was:", groqApiResponse.Choices[0].Message.Content)
            channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: http_response}
            return
        }

        println("Parsed values:")
        println("DoesVideoHaveSponsorship:", sponsorshipContent.DoesVideoHaveSponsorship)
        println("SponsorshipSubtitle:", sponsorshipContent.SponsorshipSubtitle)
    } else {
        println("No choices in response")
    }

    println("||6")
    channel_for_groq_response <- String_and_error_channel_for_groq_response{err: nil, groqApiResponsePtr: &groqApiResponse, http_response_for_go_api_ptr: http_response}
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
