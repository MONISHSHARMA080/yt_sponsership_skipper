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
	print("time passes -->",time.Since(a).Seconds(), "\n\n\n")

}


func TestGetTheSubtitlesMediantime(t *testing.T) {
    httP_client_1 := http.Client{}
    youtubeUrl := "https://www.youtube.com/watch?v=sS6u5UU3t3c"
    
    // Create a channel with buffer size of 1 to prevent goroutine leaks
    channel_for_subtitles := make(chan string_and_error_channel_for_subtitles, 1)
    
    // Number of iterations
    iterations := 14
    
    // Slice to store individual run times
    runTimes := make([]time.Duration, iterations)
    
    var transcript *Transcripts
    
    // Run the test iterations times
    for i := 0; i < iterations; i++ {
        start := time.Now()
        
        go Get_the_subtitles(httP_client_1, youtubeUrl, channel_for_subtitles, transcript)
        
        result := <-channel_for_subtitles
        
        // Record the time taken for this iteration
        
        
        // Handle results
        if result.err != nil {
            t.Logf("Iteration %d error: %v", i+1, result.err)
			t.Fail()
            return
        }
        
        if result.transcript != nil {
			runTimes[i] = time.Since(start)
            t.Logf("Iteration %d success - Sample subtitle: Start=%v, Text=%v, Dur=%v",
                i+1,
                result.transcript.Subtitles[1].Start,
                result.transcript.Subtitles[1].Text,
                result.transcript.Subtitles[1].Dur)
        }
    }
    
    // Calculate average time
    var totalTime time.Duration
    for _, duration := range runTimes {
        totalTime += duration
    }
    avgTime := totalTime / time.Duration(iterations)
    
    // Print detailed timing information
    for i, duration := range runTimes {
        t.Logf("Run %d: %v", i+1, duration)
    }
    t.Logf("\nTiming Results:")
    t.Logf("Total time: %v", totalTime)
    t.Logf("Average time per iteration: %v", avgTime)
    t.Logf("Individual run times:")
}