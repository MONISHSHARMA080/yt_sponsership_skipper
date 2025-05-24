package askllm

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	commonresultchannel "youtubeAdsSkipper/pkg/askLLM/commonResultChannel"
	askllm "youtubeAdsSkipper/pkg/askLLM/groqHelper"

	"go.uber.org/zap"
	"google.golang.org/genai"
)

type geminiResponseType struct {
	DoesVideoHaveSponsorship bool   `json:"does_video_have_sponsorship"`
	SponsorshipSubtitle      string `json:"sponsorship_subtitle"`
}

func AskGeminiAboutSponsorShipAndGetTheSponsorTiming(videoScript string, result_for_subtitles askllm.String_and_error_channel_for_subtitles,
	ChanForResponseForGettingSubtitlesTiming chan askllm.ResponseForGettingSubtitlesTiming, resultChannel chan commonresultchannel.ResultAndErrorChannel[askllm.ResponseForWhereToSkipVideo],
	logger *zap.Logger,
) {
	response := commonresultchannel.ResultAndErrorChannel[askllm.ResponseForWhereToSkipVideo]{}
	ctx := context.Background()
	apiKey, err := getRandomApiKey()
	if err != nil {
		logger.Warn("there is a error in getting a random api key and it is->", zap.Error(err))
		response.Result.FillTheStructForError("somethign went wrong on our side", http.StatusInternalServerError)
		response.Err = err
		resultChannel <- response
		return
	}
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		logger.Warn("there is a getting the gemini client", zap.Error(err))
		response.Result.FillTheStructForError("somethign went wrong on our side", http.StatusInternalServerError)
		response.Err = err
		resultChannel <- response
		return
	}

	geminiSystemPrompt := os.Getenv("GEMINI_MESSAGE_CONTENT")
	logger.Info("gemini system prompt that we go from the env", zap.String("gemini system prompt", geminiSystemPrompt))

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(geminiSystemPrompt, genai.RoleUser),
		ResponseMIMEType:  "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"does_video_have_sponsorship": {Type: genai.TypeBoolean},
				"sponsorship_subtitle":        {Type: genai.TypeString},
			},
			// PropertyOrdering: []string{"recipeName", "ingredients"},
		},
	}
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(videoScript),
		config,
	)
	if err != nil {
		println("there is a getting the response form the gemini client  and it is ", err.Error())
		response.Result.FillTheStructForError("somethign went wrong on our side", http.StatusInternalServerError)
		response.Err = err
		resultChannel <- response
		return
	}
	logger.Info("gemini's response that we got", zap.Any("gemini's whole response", result))
	logger.Info("gemini's response.text() ", zap.String("gemini's result.text()", result.Text()))

	geminiResponse, err := getGemniResponseDecoded(result.Text())
	if err != nil {
		println("there is a error in decoding the gemini response and it is ", err.Error())
		response.Result.FillTheStructForError("somethign went wrong on our side", http.StatusInternalServerError)
		response.Err = err
		resultChannel <- response
		return
	}
	fmt.Printf("the gemini's decoded  response is -> %+v \n", geminiResponse)
	logger.Info("gemini's decoded response", zap.Any("decoded response from gemini", geminiResponse))

	if !geminiResponse.DoesVideoHaveSponsorship {
		println("the video does not contain subtitles")
		response.Result.FillTheStructForSuccess("Success, got the subtitles, and where to skip in the video", http.StatusOK, 0, 0, false)
		response.Err = nil
		resultChannel <- response
		return
	}
	go askllm.GetTimeAndDurInTheSubtitles(result_for_subtitles.Transcript, &geminiResponse.SponsorshipSubtitle, &result_for_subtitles.String_value, ChanForResponseForGettingSubtitlesTiming)
	subtitleTimingResponse := <-ChanForResponseForGettingSubtitlesTiming
	println("got the subtitleTimingResponse in the gemini")
	if subtitleTimingResponse.Err != nil {
		if subtitleTimingResponse.Err.Error() == "" {
			logger.Warn("there is a error in decoding the gemini's Response", zap.Error(subtitleTimingResponse.Err))
			response.Result.FillTheStructForError("no subtitle found for the video", http.StatusNotFound)
			response.Err = err
			resultChannel <- response
			return
		}
		logger.Warn("there is a error in decoding the gemini's Response", zap.Error(subtitleTimingResponse.Err))
		response.Result.FillTheStructForError("Something is wrong on our side, error getting subtitles timming ", http.StatusInternalServerError)
		response.Err = err
		resultChannel <- response
		return
	} else if subtitleTimingResponse.EndTime+subtitleTimingResponse.StartTime <= 0 {
		// the subtitle is not being found via the function or somethign is wrong
		println("\n\n  ----- the subtitle is not gettign found despite of the llm telling us that it is there ------   \n\n")
		logger.Info("the subtitle is not gettign found despite of the llm telling us that it is there", zap.Skip())
		logger.Warn("the subtitle is not gettign found despite of the llm telling us that it is there", zap.Skip())
		response.Result.FillTheStructForError("Something is wrong on our side, error getting subtitles timming is not found of the subtitle ", http.StatusInternalServerError)
		response.Err = fmt.Errorf("the subtitle is not gettign found despite of the llm telling us that it is there ")
		resultChannel <- response
		return
	}
	response.Result.FillTheStructForSuccess("Success, got the subtitles, and where to skip in the video", http.StatusOK, int64(subtitleTimingResponse.StartTime), int64(subtitleTimingResponse.EndTime), geminiResponse.DoesVideoHaveSponsorship)
	response.Err = nil
	resultChannel <- response
}

func getRandomApiKey() (string, error) {
	keysForPaidUser := os.Getenv("NO_OF_KEYS_FOR_PAID_USER")
	TotalKeysForPaidUser, err := strconv.ParseInt(keysForPaidUser, 10, 64)
	if err != nil {
		println("there is a error in converting the keys of paid user for gemini to int and it is ", err.Error())
		return "", err
	}
	randomKeyNo := rand.Intn(int(TotalKeysForPaidUser))
	keyStr := fmt.Sprintf("GEMINI_API_KEY%d", randomKeyNo)
	println("gemini key to get form the env is --> ", keyStr)
	envKeyForGemini := os.Getenv(keyStr)
	println("env key for gemini is ", envKeyForGemini[0:len(envKeyForGemini)-6])
	return envKeyForGemini, nil
}

func getGemniResponseDecoded(geminisResponse string) (geminiResponseType, error) {
	var geminiResponse geminiResponseType
	err := json.Unmarshal([]byte(geminisResponse), &geminiResponse)
	if err != nil {
		println("there is a error in unmarshalling the gemini response and it is ", err.Error())
		return geminiResponse, err
	}
	return geminiResponse, nil
}
