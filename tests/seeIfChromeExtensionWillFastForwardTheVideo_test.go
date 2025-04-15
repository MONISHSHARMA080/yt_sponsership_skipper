package tests

import (
	"testing"
	"time"
	commonstateacrosstest "youtubeAdsSkipper/tests/helperPackages/commonStateAcrossTest"
	"youtubeAdsSkipper/tests/helperPackages/extension"
)

func TestSeeIfChromeExtensionSkipsTheVideo(t *testing.T) {
	// mf We only got 2 days (sunday and monday) today implement 2 ;
	// sunday: this one and other one
	//
	//
	//
	// Implementation: I will give you 4 url form tldr and big box swe (both)
	// then go to the first one's yt page,
	// listen to the yt video's current time and get it in the go (by saving it in the array etc, get the current time event listener)
	// [DO THIS BEFORE] when in the go land, get the network req form the chrome extension, if the sponser is not there re run it or go
	// to the next one, and re do it,
	// if we do have a sponser then get the start and the end time form the  req and see in the current time array that are we skipping it
	// is there a time skip in the array (make it some sort of error prone like there can be 1 to 2 sec lag that is ok but check for
	// the time skipped is in the range of the network req)

	ctx := commonstateacrosstest.BrowserContext
	youtubeUrl := []string{"https://www.youtube.com/watch?v=korOpibkm6g", "https://www.youtube.com/watch?v=Fndo8vzTPoM", "https://www.youtube.com/watch?v=NOfUCMzBNVg", "https://www.youtube.com/watch?v=D3cjV3tNd88", "https://www.youtube.com/watch?v=WVn4FPULFWA"}
	chromeExtension := extension.ChromeExtension{ExtensionId: extensionID}
	success := false
	for i, pageUrl := range youtubeUrl {
		// go to the url
		err := chromeExtension.NavigateToAWebPage(ctx, pageUrl)
		if err != nil {
			println("there is a error in navigating to the youtubeUrl at index ", i)
			continue
		}
		time.Sleep(time.Second * 1)
		println("we are in the youtube video ->", pageUrl)
		stopChannelToStopChekingIfTheVideoIsPlaying := make(chan struct{})
		defer close(stopChannelToStopChekingIfTheVideoIsPlaying) // this is a send only channel so only we can close it
		go chromeExtension.EnsureVideoIsPlayingPeriodically(ctx, time.Millisecond*700, stopChannelToStopChekingIfTheVideoIsPlaying, false)
		chromeExtension.IfThereIsAAdThenFinishIt(ctx, time.Millisecond*800, false)

		//
		// or) may be just get the U-block lite and let it skip the ads instead
		//
		//
		// or in the js we can take the duration of the video if it keeps changing then we will know that we have encountoured an ad and we will take the last one
		// we will not need this level of autism
		//
		//or
		//
		//run it for 2 min and then in a

		// now get the current time array and see if we have skipped the video

		println("the first video in the array has ended and we are about to go to the new one")
		time.Sleep(time.Minute * 6)
	}

	// if after all we are not able to success fully predict one of them then fail: all the videos have ads so we should be able to predict one of them
	if !success {
		t.Fatal("we are not able to skip the yt video in the chrome extension(in all 4 of the url)")
	}
}
