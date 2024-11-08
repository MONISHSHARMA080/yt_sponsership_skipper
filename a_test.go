package main

import (
	"net/http"
	"testing"
	"time"
)

func TestGetTheSubtitles(t *testing.T) {
	// func get_the_subtitles(httpClient http.Client, youtubeUrl string, want_text_without_time bool, channel_for_subtitles chan<- string_and_error_channel) {
	a:= time.Now()
	httP_client_1 := http.Client{}
	youtubeUrl := "https://www.youtube.com/watch?v=DjBrLbOJXHA"
	want_text_without_time := true
	channel_for_subtitles := make(chan string_and_error_channel)
	println("sleeping")
	time.Sleep(3000)
	println("finished sleeping")
	go Get_the_subtitles(httP_client_1, youtubeUrl, want_text_without_time, channel_for_subtitles)
	result := <- channel_for_subtitles
	if result.err != nil {
		print("error ocurred -->\n")
	}
	print("\n\n",result.string_value+"<<===,,,result value was this \n\n")
	print(time.Since(a).Seconds())

}
