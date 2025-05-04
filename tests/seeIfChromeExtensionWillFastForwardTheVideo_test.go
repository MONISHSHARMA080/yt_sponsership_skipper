package tests

import (
	"fmt"
	"testing"
	"time"
	commonchanneltype "youtubeAdsSkipper/tests/helperPackages/CommonChannelType"
	commonstateacrosstest "youtubeAdsSkipper/tests/helperPackages/commonStateAcrossTest"
	"youtubeAdsSkipper/tests/helperPackages/extension"
	"youtubeAdsSkipper/tests/helperPackages/extension/types"
)

func TestSeeIfChromeExtensionSkipsTheVideo(t *testing.T) {
	// Implementation: I will give you 4 url form tldr and big box swe (both)
	// then go to the first one's yt page,
	// listen to the yt video's current time and get it in the go (by saving it in the array etc, get the current time event listener)
	// [DO THIS BEFORE] when in the go land, get the network req form the chrome extension, if the sponser is not there re run it or go
	// to the next one, and re do it,
	// if we do have a sponser then get the start and the end time form the  req and see in the current time array that are we skipping it
	// is there a time skip in the array (make it some sort of error prone like there can be 1 to 2 sec lag that is ok but check for
	// the time skipped is in the range of the network req)

	ctx := commonstateacrosstest.BrowserContext
	youtubeUrl := []string{"https://www.youtube.com/watch?v=korOpibkm6g", "https://www.youtube.com/watch?v=D3cjV3tNd88", "https://www.youtube.com/watch?v=NOfUCMzBNVg", "https://www.youtube.com/watch?v=WVn4FPULFWA"}
	chromeExtension := extension.ChromeExtension{ExtensionId: extensionID}
	getAPIResponseFromNetworkChann := make(chan commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse])
	success := false
	for i, pageUrl := range youtubeUrl {

		// there are problems with this test, for eg why is the  function after the video completion does not return
		// the problem for that is that the channel getAPIResponseFromNetworkChann does not return

		go chromeExtension.GetResponseFromServerToChromeExtension(ctx, time.Minute*8, getAPIResponseFromNetworkChann)
		println("sleeping for 2 sec to ensure that the we are able to intercept the message form the service worker")
		time.Sleep(time.Second * 2)
		err := chromeExtension.NavigateToAWebPage(ctx, pageUrl)
		if err != nil {
			println("there is a error in navigating to the youtubeUrl at index ", i)
			continue
		}
		println("we are in the youtube video ->", pageUrl)
		// make sure to close this channel as is depends on me to close it
		stopChannelToStopChekingIfTheVideoIsPlaying := make(chan struct{})
		// defer close(stopChannelToStopChekingIfTheVideoIsPlaying) // this is a send only channel so only we can close it
		resultChanForTrackingPlayBackTime := make(chan commonchanneltype.GenericResultChannel[*[]float64])
		go chromeExtension.EnsureVideoIsPlayingPeriodically(ctx, time.Second*2, stopChannelToStopChekingIfTheVideoIsPlaying, false)
		go chromeExtension.TrackVideoPlaybackTime(ctx, resultChanForTrackingPlayBackTime)
		playBackTimeChan := <-resultChanForTrackingPlayBackTime
		println("the playBakc result is here")
		if playBackTimeChan.Err != nil {
			t.Fatal("there is a error in getting the playBackTime[] and it is ->" + playBackTimeChan.Err.Error())
		}

		println("closing the channel to see that the video is still playing as we have gotten the result")
		close(stopChannelToStopChekingIfTheVideoIsPlaying) // this is a send only channel so only we can close it

		println("is the playBackTime array not nil->", playBackTimeChan.Result != nil, " and the array lenght is:", len(*playBackTimeChan.Result))
		APIResponseFormNetwork := <-getAPIResponseFromNetworkChann
		if APIResponseFormNetwork.Err != nil {
			println("there is a error in getting the api respons form the NEtwork and it is ->", APIResponseFormNetwork.Err.Error())
			t.Fatal(APIResponseFormNetwork.Err)
		}
		fmt.Printf("we got the API resp form the network and it is -> %v", APIResponseFormNetwork.Result)
		fmt.Printf("\n Api response does contain subtitle %t -- message:%s, -- status:%d -- and start time: %f -- and end time:%f  \n", APIResponseFormNetwork.Result.ContainSponserSubtitle, APIResponseFormNetwork.Result.Message, APIResponseFormNetwork.Result.Status, APIResponseFormNetwork.Result.StartTime, APIResponseFormNetwork.Result.EndTime)
		println("the video time in the array is ->")
		for i, timeInArray := range *playBackTimeChan.Result {
			fmt.Printf("index:%d and time:%.6f \n", i, timeInArray)
		}
		didWeSkippedTheSponsorSegment, err := chromeExtension.DidWeSkippedTheAd(APIResponseFormNetwork.Result.StartTime, APIResponseFormNetwork.Result.EndTime, *playBackTimeChan.Result)
		if err != nil {
			println("there is a error in checking if we skipped the ad and it is ->", err.Error())
			t.Fatal(err)
		}
		println("did we skipped the ad ->", didWeSkippedTheSponsorSegment)
		if didWeSkippedTheSponsorSegment {
			println("we have successfully skipped the ad in the video and are exiting form the loop")
			success = true
			break
		}

		println("the first video in the array has ended and we are about to go to the new one")
	}

	// if after all we are not able to success fully predict one of them then fail: all the videos have ads so we should be able to predict one of them
	if !success {
		t.Fatal("we are not able to skip the yt video in the chrome extension(in all 4 of the url)")
	} else {
		println("we are able to skip the ad in the video and we are exiting the test")
	}
}
