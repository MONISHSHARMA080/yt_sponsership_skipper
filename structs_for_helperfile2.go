package main

import (

	"net/http"
)

// contains structs and interfaces
type String_and_error_channel_for_groq_response struct {
	err                error
	groqApiResponsePtr *GroqApiResponse
	http_response_for_go_api_ptr *http.Response
    SponsorshipContent *SponsorshipContent
}

func (groqResp *String_and_error_channel_for_groq_response) checkIfSponsorshipSubtitlesExist (transcript *Transcripts) (bool,int64,error){
	//will check the groq response and if the ads is present (in the groq resp bool) then I will go through the transcript to find it and return the time to the user
	
	
	return false, 0, nil
}









// ---------------

// type GroqApiResponse struct {
//     ID        string `json:"id,omitempty"`
//     Object    string `json:"object,omitempty"`
//     Created   int    `json:"created,omitempty"`
//     Model     string `json:"model,omitempty"`
//     Choices   []struct {
//         Index        int    `json:"index,omitempty"`
//         Message     struct {
//             Role    string `json:"role"`
//             Content string `json:"content"` // This is a JSON string that needs parsing
//         } `json:"message,omitempty"`
//         Logprobs     interface{} `json:"logprobs,omitempty"`
//         FinishReason string      `json:"finish_reason,omitempty"`
//     } `json:"choices,omitempty"`
//     Usage struct {
//         QueueTime        float64 `json:"queue_time,omitempty"`
//         PromptTokens     int     `json:"prompt_tokens,omitempty"`
//         PromptTime       float64 `json:"prompt_time,omitempty"`
//         CompletionTokens int     `json:"completion_tokens,omitempty"`
//         CompletionTime   float64 `json:"completion_time,omitempty"`
//         TotalTokens      int     `json:"total_tokens,omitempty"`
//         TotalTime       float64 `json:"total_time,omitempty"`
//     } `json:"usage,omitempty"`
//     SystemFingerprint string `json:"system_fingerprint,omitempty"`
//     XGroq            struct {
//         ID string `json:"id,omitempty"`
//     } `json:"x_groq,omitempty"`
// }

// // SponsorshipContent represents the parsed content structure
// type SponsorshipContent struct {
//     DoesVideoHaveSponsorship bool   `json:"does_video_have_sponsorship"`
//     SponsorshipSubtitle     string `json:"sponsorship_subtitle"`
// }

// // ParseContent parses the content string from the Groq response
// func ParseGroqContent(contentStr string) (bool, string, error) {
//     var content SponsorshipContent
//     err := json.Unmarshal([]byte(contentStr), &content)
//     if err != nil {
//         return false, "", fmt.Errorf("failed to parse content: %v", err)
//     }
//     return content.DoesVideoHaveSponsorship, content.SponsorshipSubtitle, nil
// }
type GroqApiResponse struct {
    ID              string `json:"id"`
    Object          string `json:"object"`
    Created         int64  `json:"created"`
    Model           string `json:"model"`
    Choices         []struct {
        Index        int `json:"index"`
        Message      struct {
            Role    string `json:"role"`
            Content string `json:"content"` // Changed to string to handle raw JSON
        } `json:"message"`
        Logprobs     interface{} `json:"logprobs"`
        FinishReason string      `json:"finish_reason"`
    } `json:"choices"`
    Usage struct {
        QueueTime        float64 `json:"queue_time"`
        PromptTokens     int     `json:"prompt_tokens"`
        PromptTime       float64 `json:"prompt_time"`
        CompletionTokens int     `json:"completion_tokens"`
        CompletionTime   float64 `json:"completion_time"`
        TotalTokens      int     `json:"total_tokens"`
        TotalTime        float64 `json:"total_time"`
    } `json:"usage"`
    SystemFingerprint string `json:"system_fingerprint"`
    XGroq             struct {
        ID string `json:"id"`
    } `json:"x_groq"`
}

type SponsorshipContent struct {
    DoesVideoHaveSponsorship bool `json:"does_video_have_sponsorship"`
    SponsorshipSubtitle      string `json:"sponsorship_subtitle"`
}