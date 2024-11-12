package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)



  
  func AskGroqabouttheSponsorship( httpClient *http.Client, channel_for_groq_response chan<- String_and_error_channel_for_groq_response )  {
	
	// make the structs for the response 
	// set up context it we want to run it async
	err ,http_req := factoryGroqPostReqCreator(os.Getenv("API_KEYS0"))
	if err!= nil {
		channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: nil} 
	}
	http_response , err := httpClient.Do(http_req)
	if err!= nil {
		channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: nil} 
	}
	var groqApiResponse GroqApiResponse
	err =  json.NewDecoder(http_response.Body).Decode(groqApiResponse)
	if err != nil{
		channel_for_groq_response <- String_and_error_channel_for_groq_response{err: err, groqApiResponsePtr: nil, http_response_for_go_api_ptr: http_response} 
	}
	println("\n groq api response--==>>", groqApiResponse.Choices[0].Message.Content, "<--====----\n")
	channel_for_groq_response <- String_and_error_channel_for_groq_response{err: nil, groqApiResponsePtr: &groqApiResponse, http_response_for_go_api_ptr: http_response } 
  }
  
  func factoryGroqPostReqCreator(GroqApiKey string) (error, *http.Request) {
  
	// GroqApiKey is asked as I need to decide wether the user is paid or free
  
	url := "https://api.groq.com/openai/v1/chat/completions"
	
	payload := map[string]interface{}{
		  "model": os.Getenv("MODEL"),
		  "messages": []map[string]string{
			  {
				  "role":    "user",
				  "content": os.Getenv("MESSAGE_CONTENT"),
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