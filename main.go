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

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
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

type userForDB struct {
	accountid      int64
	email          string
	UserToken      string
	is_a_paid_user bool
}

func main() {
	// httP_client_1 := http.Client{}
	// youtubeUrl := "https://www.youtube.com/watch?v=sS6u5UU3t3c"
	// want_text_without_time := true
	// channel_for_subtitles := make(chan string_and_error_channel)
	// println("sleeping")

	// time.Sleep(30000)

	err := godotenv.Load()
	if err != nil {
		println("Error loading .env file: %v", err)
		panic(err)
	}

	encryption_key := os.Getenv("encryption_key")
	encryption_key_as_byte := []byte(os.Getenv("encryption_key"))

	// key := []byte(os.Getenv("encryption_key"))
	// // --------- fill this
	// a := Signup_detail_of_user{Email: "abc@gamil.com", AccountID: 123, UserToken: "dwunewiunwioc"}
	// plaintext := []byte(return_string_based_on_user_details_for_encryption_text(a, false))
	// // -------fill this
	// // Encrypt
	// ciphertext, err := encrypt(plaintext, key)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Encrypted: %x\n", ciphertext)
	// // Decrypt
	// decryptedText, err := decrypt(ciphertext, key)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Decrypted: %s\n", decryptedText)
	//
	//
	//
	// dkj, _ := decrypt([]byte("NUntyC7TuShFIk9nGTVpWtta4lhqWKEAT3ej8ok9QHKu1rq3UY44AkqVKqgCbvxc9OTAUGUAhx50OxMS/SJF5D5ThdY="),key)

	println(encryption_key, "-----", len(encryption_key))

	httpClient := http.Client{}
	http.HandleFunc("/", getRoot)

	// if need new token use this
	http.HandleFunc("/signup", User_signup_handler(encryption_key))
	http.HandleFunc("/youtubeVideo", Return_to_client_where_to_skip_to_in_videos(encryption_key_as_byte, &httpClient))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err.Error())
	}

	// key := []byte(os.Getenv("encryption_key"))
	// // --------- fill this
	// plaintext := []byte("jeionew")
	// // -------fill this
	// // Encrypt
	// ciphertext, err := encrypt(plaintext, key)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Encrypted: %x\n", ciphertext)
	// // Decrypt
	// decryptedText, err := decrypt(ciphertext, key)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Decrypted: %s\n", decryptedText)
	// print(time.Since(startTime).Milliseconds())
	// ------for the new user
	// ------for the new user
	// defer db.Close()
	// user := userForDB{accountid: 138987, email: "mmmmmmhmewwed@heb.com", UserToken: "--------------++ebwuewbcifduewbdiewdiewiduewjhb"}
	// errorCh := make(chan error)
	// go func() {
	// 	errorCh <- InsertUserInDB(db, user)
	// }()
	// if err := <-errorCh; err != nil {
	// 	panic(err.Error())
	// }
	// error_ff :=  InsertUserInDB(db, userForDB{accountid: 3298,email: "iuewwed@heb.com", UserToken: "ebwuewbciuewbdiewdiewidu"})
	// error_ff :=  InsertUserInDB(db, userForDB{accountid: 3887,email: "iuewewisuewwed@heb.com", UserToken: "ebwuewbcifduewbdiewdiewiduewjhb"})
	// if error_ff!= nil {
	// 	panic(error_ff.Error())
	// }
	// httpClient := http.Client{}
	// youtubeUrl := "https://www.youtube.com/watch?v=X7LA_VnHoAg"
	// text_form_subtitile, err := get_the_subtitles(httpClient, youtubeUrl, true)
	// if err != nil {
	// 	// return the response but here I will panic
	// 	println("error I got -->", err)
	// 	panic(err.Error())
	// }
	// println("\n\n ------------", text_form_subtitile, "\n\n -----------")
	// --## good now it is done , just make a func to check the size of the string  which one is smaller just send that(wait if the plain text contains it
	// then how will I detect where it is , do some sort of loop on the text , just do that as it will be efficient )
	// --##  and get many llm keys
	// DbConnect()
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

	// fmt.Printf("\n\nstring in fetchAndReturnTheBodyAsString -->  %s --++\n\n", responseBodyString)
	return responseBodyString, nil
}

func fetchAndReturnTheBodyAsByte(youtubeVideoUrl string, httpClient *http.Client) ([]byte, error) {
	response, err := httpClient.Get(youtubeVideoUrl)
	defer response.Body.Close()

	if err != nil {
		return []byte{1}, err
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
	println("baseurl of the string is -->", baseUrl)
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

	httP_client_1 := http.Client{}
	htmlResponse, err := fetchAndReturnTheBodyAsString(youtubeUrl, &httP_client_1)
	if err != nil {
		channel_for_subtitles <- string_and_error_channel_for_subtitles{err: err, string_value: "", transcript: nil}
		return
	}

	var captionsDataInJson map[string]interface{}
	// probally take it as a htmlresponse *string
	err = convertHtmlToJsonAndWriteItToAMAp(htmlResponse, &captionsDataInJson)
	if err != nil {
		channel_for_subtitles <- string_and_error_channel_for_subtitles{err: err, string_value: "", transcript: nil}
		return
	}
	// printJson(captionsDataInJson)

	baseUrl, err := return_caption_url(captionsDataInJson)
	if err != nil {
		channel_for_subtitles <- string_and_error_channel_for_subtitles{err: err, string_value: "", transcript: nil}
		return
	}
	captionsInXML, errorF := fetchAndReturnTheBodyAsByte(baseUrl, &httpClient)
	if errorF != nil {
		channel_for_subtitles <- string_and_error_channel_for_subtitles{err: errorF, string_value: "", transcript: nil}
		return
	}

	transcripts := Transcripts{}
	errorInXMl := xml.Unmarshal(captionsInXML, &transcripts)
	if errorInXMl != nil {
		channel_for_subtitles <- string_and_error_channel_for_subtitles{err: errorInXMl, string_value: "", transcript: nil}
		return
	}
	// formatting the transcript to be in utf-8
	println("formatting the transctipt.subtitles.text to be utf-8")
	for i, text := range transcripts.Subtitles {
		transcripts.Subtitles[i].Text = html.UnescapeString(text.Text)
	}

	// for _, subtitle := range transcripts.Subtitles {
	// 	fmt.Printf("[start %s] %s [Duration: %s]\n", subtitle.Start, subtitle.Text, subtitle.Dur)
	// }

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
