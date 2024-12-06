package main

import (
	"strings"
	"testing"
)

func TestGetTimeAndDurInTheSubtitles(t *testing.T) {
	transcript := Transcripts{
		Subtitles: []Subtitle{
			{Text: "Hello, and welcome to the presentation.", Start: "1000", Dur: "3500"},
			{Text: "Today, we'll discuss the topic in depth.", Start: "4000", Dur: "3000"},
			{Text: "Feel free to ask any questions during the session.", Start: "7500", Dur: "4000"},
			{Text: "Now, let's dive into the first section.", Start: "12000", Dur: "3000"},
			{Text: "This section covers the basics of the subject.", Start: "15500", Dur: "4500"},
			{Text: "Make sure to take notes as we go along.", Start: "20500", Dur: "2500"},
			{Text: "Moving on, let's explore the advanced concepts.", Start: "24000", Dur: "3500"},
			{Text: "These concepts are essential for deeper understanding.", Start: "28000", Dur: "3000"},
			{Text: "We'll now summarize the key points covered.", Start: "32000", Dur: "2500"},
			{Text: "Thank you for your attention!", Start: "35000", Dur: "2000"},
		},
	}
	fullCaption := "Hello, and welcome to the presentation. Today, we'll discuss the topic in depth. Feel free to ask any questions during the session. Now, let's dive into the first section. This section covers the basics of the subject. Make sure to take notes as we go along. Moving on, let's explore the advanced concepts. These concepts are essential for deeper understanding. We'll now summarize the key points covered. Thank you for your attention!"
	sponserShipSubtitles := "the advanced concepts. These concepts are essential for deeper understanding. We'll now summarize the key points covered. Thank"
	pik := strings.Index(strings.ToLower(fullCaption), strings.ToLower(sponserShipSubtitles))
	println("and is -->", pik)
	startTime, endTime, err := GetTimeAndDurInTheSubtitles(&transcript, &sponserShipSubtitles, &fullCaption)
	if err != nil {
		println("error occurred")
		t.Fatal(err)
		panic(err)
	}
	println("StartTime is -->", startTime, "endtime is -->", endTime)
	if startTime != 24000 && endTime != 35000 {
		t.Fail()
		t.FailNow()
		t.Fatal("function did not provided the expected output")
	}
}

// 	// func get_the_subtitles(httpClient http.Client, youtubeUrl string, want_text_without_time bool, channel_for_subtitles chan<- string_and_error_channel) {
// 	httP_client_1 := http.Client{}
// 	youtubeUrl := "https://www.youtube.com/watch?v=xSBGYoS6z68"

// 	channel_for_subtitles := make(chan string_and_error_channel_for_subtitles)

// 	a := time.Now()
// 	go Get_the_subtitles(httP_client_1, youtubeUrl, channel_for_subtitles)
// 	result := <-channel_for_subtitles
// 	if result.err != nil {
// 		print("error ocurred -->" + result.err.Error() + "\n")
// 	}
// 	if result.transcript != nil {
// 		bb := result.transcript
// 		println("xml -->", bb.Subtitles[1].Start, bb.Subtitles[1].Text, bb.Subtitles[1].Dur)
// 	} else {
// 		println("")
// 	}
// 	print("time passes -->", time.Since(a).Seconds(), "\n\n\n")

// }

// func TestGetTheSubtitlesMediantime(t *testing.T) {
// 	httP_client_1 := http.Client{}
// 	youtubeUrl := "https://www.youtube.com/watch?v=sS6u5UU3t3c"

// 	// Create a channel with buffer size of 1 to prevent goroutine leaks
// 	channel_for_subtitles := make(chan string_and_error_channel_for_subtitles, 1)

// 	// Number of iterations
// 	iterations := 14

// 	// Slice to store individual run times
// 	runTimes := make([]time.Duration, iterations)

// 	// Run the test iterations times
// 	for i := 0; i < iterations; i++ {
// 		start := time.Now()

// 		go Get_the_subtitles(httP_client_1, youtubeUrl, channel_for_subtitles)

// 		result := <-channel_for_subtitles

// 		// Record the time taken for this iteration

// 		// Handle results
// 		if result.err != nil {
// 			t.Logf("Iteration %d error: %v", i+1, result.err)
// 			t.Fail()
// 			return
// 		}

// 		if result.transcript != nil {
// 			runTimes[i] = time.Since(start)
// 			t.Logf("Iteration %d success - Sample subtitle: Start=%v, Text=%v, Dur=%v",
// 				i+1,
// 				result.transcript.Subtitles[1].Start,
// 				result.transcript.Subtitles[1].Text,
// 				result.transcript.Subtitles[1].Dur)
// 		}
// 	}

// 	// Calculate average time
// 	var totalTime time.Duration
// 	for _, duration := range runTimes {
// 		totalTime += duration
// 	}
// 	avgTime := totalTime / time.Duration(iterations)

// 	// Print detailed timing information
// 	for i, duration := range runTimes {
// 		t.Logf("Run %d: %v", i+1, duration)
// 	}
// 	t.Logf("\nTiming Results:")
// 	t.Logf("Total time: %v", totalTime)
// 	t.Logf("Average time per iteration: %v", avgTime)
// 	t.Logf("Individual run times:")
// }
