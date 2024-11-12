package main

import (
	"net/http"
	"testing"
	"time"
)

func TestGetTheSubtitles(t *testing.T) {
	// func get_the_subtitles(httpClient http.Client, youtubeUrl string, want_text_without_time bool, channel_for_subtitles chan<- string_and_error_channel) {

	httP_client_1 := http.Client{}
	youtubeUrl := "https://www.youtube.com/watch?v=sS6u5UU3t3c"

	channel_for_subtitles := make(chan string_and_error_channel_for_subtitles)
	println("sleeping")

	time.Sleep(30000)

	println("finished sleeping")
	a := time.Now()
	var transcript *Transcripts
	go Get_the_subtitles(httP_client_1, youtubeUrl, channel_for_subtitles, transcript)
	result := <-channel_for_subtitles
	if result.err != nil {
		print("error ocurred -->" + result.err.Error() + "\n")
	}
	if result.transcript != nil {
		bb := result.transcript
		println("xml -->", bb.Subtitles[1].Start, bb.Subtitles[1].Text, bb.Subtitles[1].Dur )
	}else{
		println()
	}
	print("\n\n", result.string_value+"<<===,,,result value was this \n\n")
	print("time passes -->",time.Since(a).Seconds())

}
