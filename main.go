package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strings"
	routehandlerfunc "youtubeAdsSkipper/RouteHandlerFunc/GetNewKey"
	paymentbackendgo "youtubeAdsSkipper/paymentBackendGO"
	handlerfunction "youtubeAdsSkipper/paymentBackendGO/handlerFunction"
	llmreqratelimiter "youtubeAdsSkipper/pkg/LLmReqRateLimiter"
	askllmHelper "youtubeAdsSkipper/pkg/askLLM/groqHelper"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// type Subtitle struct {
// 	Text  string `xml:",chardata"`
// 	Start string `xml:"start,attr"`
// 	Dur   string `xml:"dur,attr"`
// }

// // Transcripts holds an array of Subtitle elements
// type Transcripts struct {
// 	Subtitles []Subtitle `xml:"text"`
// }

// type alias so that I have common types
type (
	Transcripts = askllmHelper.Transcripts
	Subtitle    = askllmHelper.Subtitle
)

// ----------------- client redirect

var (
	// change this to the same URI you configured in Google Cloud Console:
	oauth2RedirectURL = "http://localhost:8080/oauth2callback"
	config            *oauth2.Config
)

func main() {
	err := godotenv.Load()
	if err != nil {
		println("Error loading .env file: %s", err.Error())
		railwayProjectName := os.Getenv("RAILWAY_PROJECT_NAME")
		if railwayProjectName == "" {
			panic(err)
		} else {
			println("the railway project name is -->", railwayProjectName)
		}
		isthisTestingEnv := os.Getenv("IS_THIS_TESTING_ENVIRONMENT")
		if isthisTestingEnv == "" {
			panic(err)
		} else {
			println("the is this testing env is -->", isthisTestingEnv)
			println("we are able to get the env")
		}
	}
	rateLimiterDb, err := llmreqratelimiter.GetRateLimiterDb("file:./rateLimitterForReq.db")
	if err != nil {
		println("there is a error in making the rate limiter DB and we are crashing, err is ->", err.Error())
		panic(err)
	}

	encryption_key := os.Getenv("encryption_key")
	encryption_key_as_byte := []byte(os.Getenv("encryption_key"))

	println(encryption_key, "-----", len(encryption_key))

	// -- structurred logger
	logger, err := zap.NewProduction()
	if err != nil {
		println("there is a error in initializing  the logger and it is ->", err.Error())
		panic(err)
	}
	defer logger.Sync()

	// config = loadOAuthConfig()
	httpClient := http.Client{}
	http.HandleFunc("/", getRoot)

	// http.HandleFunc("/login", handleLogin)
	// http.HandleFunc("/oauth2callback", handleCallback)
	// if need new token use this
	http.HandleFunc("/signup", User_signup_handler(encryption_key))
	http.HandleFunc("/youtubeVideo", Return_to_client_where_to_skip_to_in_videos(encryption_key_as_byte, &httpClient, rateLimiterDb, logger))
	http.HandleFunc("/checkIfKeyIsValid", CheckIfKeyIsValid(encryption_key_as_byte))
	http.HandleFunc("/makeAPayment", paymentbackendgo.CreateAndReturnOrderId(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_SECRET_ID"), encryption_key_as_byte))
	http.HandleFunc("/validatePayment", handlerfunction.VerifyPaymentSignature(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_SECRET_ID"), encryption_key_as_byte))
	http.HandleFunc("/webHookIntegrationForPaymentCapture", handlerfunction.WebHookIntegrationForPaymentCapture(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_SECRET_ID"), os.Getenv("RAZORPAYWEBHOOKSECRET"), encryption_key_as_byte))
	http.HandleFunc("/getNewKey", routehandlerfunc.GetNewKey())

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err.Error())
	}
}

func fetchAndReturnTheBodyAsString(youtubeVideoUrl string, httpClient *http.Client) (string, error) {
	response, err := httpClient.Get(youtubeVideoUrl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		return "", err
	}
	responseBodyString := buf.String()
	// println("the response body is ->", responseBodyString, "\n\n ++++")
	fmt.Printf("the youtube html status is %s, and the status code is %d \n\n ", response.Status, response.StatusCode)

	// fmt.Printf("\n\nstring in fetchAndReturnTheBodyAsString -->  %s --++\n\n", responseBodyString)
	return responseBodyString, nil
}

func fetchAndReturnTheBodyAsByte(youtubeVideoUrl string, httpClient *http.Client) ([]byte, error) {
	response, err := httpClient.Get(youtubeVideoUrl)
	if err != nil {
		return []byte{1}, err
	}
	defer response.Body.Close()
	responseBodyByte, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{1}, err
	}
	return responseBodyByte, nil
}

func printJson(data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}
	fmt.Println(string(jsonBytes))
}

func convertHtmlToJsonAndWriteItToAMAp(htmlResponse string, captionsData *map[string]interface{}) error {
	parts := strings.Split(htmlResponse, `"captions":`)
	if len(parts) <= 1 {
		println("no captions found for video ")
		return fmt.Errorf("no captions found for video")
	}
	jsonString := strings.Split(parts[1], `,"videoDetails"`)[0]

	err := json.Unmarshal([]byte(jsonString), captionsData)
	if err != nil {
		return err
	}
	return nil
}

func return_caption_url(captionsDataInJson map[string]interface{}) (string, error) {
	// why -->
	// error in the patrick video, https://www.youtube.com/watch?v=Wx51CffrBIg, it is returning the subtitles in arabic by default , so I should check
	// the lang to be english (auto etc) and  avoid the loss when converting bytes to str in fetchAndReturnTheBodyAsString

	playerCaptionsTracklistRenderer, ok := captionsDataInJson["playerCaptionsTracklistRenderer"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("playerCaptionsTracklistRenderer not found")
	}

	captionTracks, ok := playerCaptionsTracklistRenderer["captionTracks"].([]interface{})
	if !ok || len(captionTracks) == 0 {
		return "", fmt.Errorf("captionTracks not found or empty")
	}

	var englishTrack, englishAutoTrack, firstTrack map[string]interface{}

	for _, track := range captionTracks {
		captionTrack, ok := track.(map[string]interface{})
		if !ok {
			continue
		}
		languageCode, _ := captionTrack["languageCode"].(string)
		name, _ := captionTrack["name"].(map[string]interface{})
		simpleText, _ := name["simpleText"].(string)

		if languageCode == "en" {
			if simpleText == "English" {
				englishTrack = captionTrack
				println("in the english track")
				break
			} else if simpleText == "English (auto-generated)" {
				println("in the auto generated track")
				englishAutoTrack = captionTrack
			}
		}

		if firstTrack == nil {
			firstTrack = captionTrack
		}
	}

	var selectedTrack map[string]interface{}

	if englishTrack != nil {
		selectedTrack = englishTrack
	} else if englishAutoTrack != nil {
		selectedTrack = englishAutoTrack
	} else if firstTrack != nil {
		selectedTrack = firstTrack
	} else {
		return "", fmt.Errorf("no suitable caption track found")
	}

	baseUrl, ok := selectedTrack["baseUrl"].(string)
	// println("baseurl of the string is -->", baseUrl)
	if !ok {
		return "", fmt.Errorf("baseUrl not found in selected captionTrack")
	}

	return baseUrl, nil
}

// --- 2nd one for utf-8 strings

func generateSubtitleString(subtitles []Subtitle) string {
	var result strings.Builder
	// Concatenate the subtitle text with spaces in between
	for _, subtitle := range subtitles {
		result.WriteString(html.UnescapeString(strings.TrimSpace(subtitle.Text) + " "))
	}
	// Convert the result into a string and return
	return result.String()
}

func Get_the_subtitles(httpClient http.Client, youtubeUrl string, channel_for_subtitles chan<- string_and_error_channel_for_subtitles) {
	println(" in the get_the_subtitles func")
	httP_client_1 := &http.Client{}
	println("fetching the body and returning it as a string")

	htmlResponse, err := fetchAndReturnTheBodyAsString(youtubeUrl, httP_client_1)
	if err != nil {
		println("there is a error in fetching the youtube body and is ->", err.Error())
		channel_for_subtitles <- string_and_error_channel_for_subtitles{err: err, string_value: "", transcript: nil}
		return
	}
	println("does the html response contiand the ytInitialPlayerResponse", strings.Contains(htmlResponse, "ytInitialPlayerResponse"))

	indexa := strings.Index(htmlResponse, "ytInitialPlayerResponse")
	println("in the string after the ytInitialPlayerResponse do we have playerCaptionsTracklistRenderer there too", strings.Contains(htmlResponse[indexa:], "playerCaptionsTracklistRenderer"), " \n the text after that is(80 plus) ->")
	if len(htmlResponse) > indexa+5380 {
		a := htmlResponse[indexa : indexa+5380]
		println("+5380 text is --> ", a, "--- and the length of the indexa:indexa+5380 is ->", len(a))
	} else {
		println("we can't have the text  as +80 dosen't exist")
	}

	var captionsDataInJson map[string]interface{}
	// probally take it as a htmlresponse *string
	err = convertHtmlToJsonAndWriteItToAMAp(htmlResponse, &captionsDataInJson)
	if err != nil {
		println("there is a error in converting the html to json and writing it to a map->", err.Error())
		channel_for_subtitles <- string_and_error_channel_for_subtitles{err: err, string_value: "", transcript: nil}
		return
	}
	// printJson(captionsDataInJson)

	baseUrl, err := return_caption_url(captionsDataInJson)
	if err != nil {
		println("there is a error in returning the captiin url ->", err.Error())
		channel_for_subtitles <- string_and_error_channel_for_subtitles{err: err, string_value: "", transcript: nil}
		return
	}
	captionsInXML, errorF := fetchAndReturnTheBodyAsByte(baseUrl, &httpClient)
	if errorF != nil {
		println("there is a error in fetching the captionsInXML->", errorF.Error())
		channel_for_subtitles <- string_and_error_channel_for_subtitles{err: errorF, string_value: "", transcript: nil}
		return
	}

	transcripts := Transcripts{}
	errorInXMl := xml.Unmarshal(captionsInXML, &transcripts)
	if errorInXMl != nil {
		println("there is a error in Unmarshaling the xml in the transcript struct and it is ->", errorInXMl.Error())
		channel_for_subtitles <- string_and_error_channel_for_subtitles{err: errorInXMl, string_value: "", transcript: nil}
		return
	}
	// formatting the transcript to be in utf-8
	println("formatting the transctipt.subtitles.text to be utf-8")
	for i, text := range transcripts.Subtitles {
		transcripts.Subtitles[i].Text = html.UnescapeString(text.Text)
	}

	for _, subtitle := range transcripts.Subtitles {
		fmt.Printf("[start %s]- %s -[Duration: %s]\n", subtitle.Start, subtitle.Text, subtitle.Dur)
	}

	// 2. Second requirement: Generate single string with format "[start] text [dur]"

	// probally need an array of some sort , like encoding it in the string is not a good idead how will I decode it later; probally itegrate over the string or asyncly convert in a array
	// I think (not thought it through) either way I will itereate through the string so why not just do it once

	// now this block is useless as I will return the xml to myself

	channel_for_subtitles <- string_and_error_channel_for_subtitles{err: nil, string_value: generateSubtitleString(transcripts.Subtitles), transcript: &transcripts}
}

func GenerateSubtitleWithTime(Subtitles []Subtitle, channel_for_subtitles chan<- string) {
	// probally need an array of some sort , like encoding it in the string is not a good idead how will I decode it later

	var result strings.Builder
	for _, subtitle := range Subtitles {
		result.WriteString("[start")
		result.WriteString(subtitle.Start)
		result.WriteString("] ")
		result.WriteString(subtitle.Text)
		result.WriteString("[Duration: ")
		result.WriteString(subtitle.Dur)
		result.WriteString("]\n")
	}
	channel_for_subtitles <- result.String()
}

func GenerateSubtitleWithTimeWithoutChannels(Subtitles []Subtitle) string {
	// probally need an array of some sort , like encoding it in the string is not a good idead how will I decode it later

	var result strings.Builder
	for _, subtitle := range Subtitles {
		result.WriteString("[start")
		result.WriteString(subtitle.Start)
		result.WriteString("] ")
		result.WriteString(subtitle.Text)
		result.WriteString("[Duration: ")
		result.WriteString(subtitle.Dur)
		result.WriteString("]\n")
	}
	return result.String()
}
