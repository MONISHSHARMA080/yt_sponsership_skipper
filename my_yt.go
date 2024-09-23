package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"
)

type Subtitle struct {
	Text  string `xml:",chardata"`
	Start string `xml:"start,attr"`
	Dur   string `xml:"dur,attr"`
}

// Transcripts holds an array of Subtitle elements
type Transcripts struct {
	Subtitles []Subtitle `xml:"text"`
}


func main() {

	startTime := time.Now()

	httpClient := http.Client{}
	htmlResponse, err := fetchAndReturnTheBodyAsString("https://www.youtube.com/watch?v=X7LA_VnHoAg", &httpClient)
	if err != nil {
		panic(err.Error())
	}

	var captionsDataInJson map[string]interface{}

	err = convertHtmlToJsonAndWriteItToAMAp(htmlResponse, &captionsDataInJson)
	if err != nil {
		panic(err.Error())
	}

	// printJson(captionsDataInJson)

	baseUrl, err := return_caption_url(captionsDataInJson)
	if err != nil {

		panic(err.Error())
	}

	captionsInXML, errorF := fetchAndReturnTheBodyAsByte(baseUrl, &httpClient)
	if errorF != nil {
		panic(errorF.Error())
	}



	transcripts := Transcripts{}
	errorInXMl := xml.Unmarshal(captionsInXML, &transcripts)
	if errorInXMl != nil {
		panic(errorInXMl.Error())
	}

	// for _, subtitle := range transcripts.Subtitles {
	// 	fmt.Printf("[start %s] %s [Duration: %s]\n", subtitle.Start,  subtitle.Text, subtitle.Dur)
	// }

	// 2. Second requirement: Generate single string with format "[start] text [dur]"
	resultString := generateSubtitleString(transcripts.Subtitles)
	println(resultString )

	println("baseUrl -->", baseUrl)
	fmt.Printf("Time taken: %d ms\n", time.Since(startTime).Milliseconds())

	// --## good now it is done , just make a func to check the size of the string  which one is smaller just send that(wait if the plain text contains it 
	// then how will I detect where it is , do some sort of loop on the text , just do that as it will be efficient )
	// --##  and get many llm keys

}


func fetchAndReturnTheBodyAsString(youtubeVideoUrl string, httpClient *http.Client) (string, error) {
	response, err := httpClient.Get(youtubeVideoUrl)
	defer response.Body.Close()

	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		return "", err
	}
	responseBodyString := buf.String()

	// fmt.Printf("\n\nstring in fetchAndReturnTheBodyAsString -->  %s --++\n\n", responseBodyString)
	return responseBodyString, nil
}
func fetchAndReturnTheBodyAsByte(youtubeVideoUrl string, httpClient *http.Client) ( []byte, error) {
	response, err := httpClient.Get(youtubeVideoUrl)
	defer response.Body.Close()

	if err != nil {
		return []byte{1} , err
	}
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
	//the lang to be english (auto etc) and  avoid the loss when converting bytes to str in fetchAndReturnTheBodyAsString

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
				break
			} else if simpleText == "English (auto-generated)" {
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
	if !ok {
		return "", fmt.Errorf("baseUrl not found in selected captionTrack")
	}

	return baseUrl, nil
}

// fetchAndReturnTheBodyAsString returns the body of the HTTP response as a string
// generateSubtitleString formats subtitles into a single string "[start] text [dur]"
// func generateSubtitleString(subtitles []Subtitle) string {
// 	var result string
// 	for _, subtitle := range subtitles {
// 		result += fmt.Sprintf(" %s ",subtitle.Text)
// 	}
// 	result += "\n"
// 	return result
// }


// --- 2nd one for utf-8 strings 

func generateSubtitleString(subtitles []Subtitle) string {
	var result string
	for _, subtitle := range subtitles {
		// Concatenate the subtitle text with spaces in between
		result += html.UnescapeString(strings.TrimSpace(subtitle.Text))
	}
	// Convert the result into a string and return
	return result
}