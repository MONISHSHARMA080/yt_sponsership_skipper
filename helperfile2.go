package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)



  
  func AskGroqabouttheSponsorship( httpClient *http.Client, channel_for_groq_response chan<- String_and_error_channel_for_groq_response, APIKEY_according_to_users_tier string, subtitlesInTheVideo *string )  {
	
	// make the structs for the response 
	// set up context it we want to run it async
	err ,http_req := factoryGroqPostReqCreator(APIKEY_according_to_users_tier, subtitlesInTheVideo )
	if err!= nil {
		channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: nil} 
	}
	http_response , err := httpClient.Do(http_req)
	if err!= nil {
		channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: nil} 
	}
	var groqApiResponse GroqApiResponse
	err =  json.NewDecoder(http_response.Body).Decode(&groqApiResponse)
	if err != nil{
		channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: http_response} 
	}
	println("\n groq api response--==>>", groqApiResponse.Choices[0].Message.Content, "<--====----\n")
	channel_for_groq_response <- String_and_error_channel_for_groq_response{err: nil, groqApiResponsePtr: &groqApiResponse, http_response_for_go_api_ptr: http_response } 
  }
  
  func factoryGroqPostReqCreator(GroqApiKey string, subtitlesInTheVideo *string) (error, *http.Request) {
  
	// stringify a json schema in the end 
  
	url := "https://api.groq.com/openai/v1/chat/completions"
	
	payload := map[string]interface{}{
		  "model": os.Getenv("GROQ_MODEL"),
		  "messages": []map[string]string{
			  {
				  "role":    "user",
				  "content": "don't forget I only need json form you nothing else; sutitles-->"+*subtitlesInTheVideo,
			  },
			  {
				"role": "model",
				"content": os.Getenv("GROQ_MESSAGE_CONTENT"),
			  },
		  },
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err,nil
	}

	createdHttpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err,nil
	}
	createdHttpReq.Header.Set("Content-Type", "application/json")
	createdHttpReq.Header.Set("Authorization", "Bearer "+GroqApiKey)
  
  return nil, createdHttpReq
  
  }

  func getAPIKEYForGroqBasedOnUsersTeir(is_user_paid bool)(string, error){
	// the env should have 0 to the numebr of keys in front of them 
	// 2 api keys for the groq ; the errors are from parsing in to str
	var err error
	var number_ofKeys int64	

	if is_user_paid == true{
		numberOFEnvKeyAccordignToUserTeir := os.Getenv("NO_OF_KEYS_FOR_PAID_USER")
		number_ofKeys, err = strconv.ParseInt(numberOFEnvKeyAccordignToUserTeir, 10,32)
	}else{
		numberOFEnvKeyAccordignToUserTeir := os.Getenv("NO_OF_KEYS_FOR_UNPAID_USER")
		number_ofKeys, err = strconv.ParseInt(numberOFEnvKeyAccordignToUserTeir, 10,32)
	}
	if err!= nil{
		return "", err
	}
	random_number_for_apiKey := rand.Intn(int(number_ofKeys))
	
	if is_user_paid == true {
		return os.Getenv("API_KEY_PAID"+strconv.Itoa(random_number_for_apiKey)),nil
	}else{
		return os.Getenv("API_KEY_UNPAID"+strconv.Itoa(random_number_for_apiKey)),nil
	}
  }