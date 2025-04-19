package tests

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
	commonchanneltype "youtubeAdsSkipper/tests/helperPackages/CommonChannelType"
	commonstateacrosstest "youtubeAdsSkipper/tests/helperPackages/commonStateAcrossTest"
	"youtubeAdsSkipper/tests/helperPackages/extension"
	"youtubeAdsSkipper/tests/helperPackages/extension/types"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
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
	youtubeUrl := []string{"https://www.youtube.com/watch?v=NOfUCMzBNVg", "https://www.youtube.com/watch?v=korOpibkm6g", "https://www.youtube.com/watch?v=D3cjV3tNd88", "https://www.youtube.com/watch?v=WVn4FPULFWA"}
	chromeExtension := extension.ChromeExtension{ExtensionId: extensionID}
	getAPIResponseFromNetworkChann := make(chan commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse])
	success := false
	for i, pageUrl := range youtubeUrl {
		// go to the url
		// go chromeExtension.GetResponseFromServerToChromeExtension(ctx, time.Minute*20, getAPIResponseFromNetworkChann)

		var infos []*target.Info
		if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			infos, err = target.GetTargets().Do(ctx)
			return err
		})); err != nil {
			println(" we have error is ->", err.Error())
			log.Fatal(err)
		}

		// 5. Find your extension’s service worker.
		var swTargetID target.ID
		prefix := "chrome-extension://" + extensionID
		for i, ti := range infos {
			println("index:", i, " and the ti title is ->", ti.Title, "ti's url is ->", ti.URL, "--and type is ", ti.Type)

			if ti.Type == "service_worker" && strings.HasPrefix(ti.URL, prefix) {
				swTargetID = ti.TargetID
				break
			}
		}
		if swTargetID == "" {
			t.Fatal("the taerget ID of the service worker is empty")
		}
		println("the target ID of the service worker is ->", swTargetID)

		swCtx, swCancel := chromedp.NewContext(ctx, chromedp.WithTargetID(swTargetID))
		defer swCancel()

		// Create a channel to receive network events
		networkEventsChan := make(chan *network.EventResponseReceived, 100)

		// Listen for network response events
		chromedp.ListenTarget(swCtx, func(ev interface{}) {
			switch e := ev.(type) {
			case *network.EventResponseReceived:
				networkEventsChan <- e
			}
		})

		// Enable network events
		if err := chromedp.Run(swCtx, network.Enable()); err != nil {
			t.Fatalf("Failed to enable network monitoring: %v", err)
		}

		// Navigate to the YouTube URL in the main browser context
		if err := chromedp.Run(ctx, chromedp.Navigate(pageUrl)); err != nil {
			t.Fatalf("Failed to navigate to URL: %v", err)
		}

		// Process network events with a timeout
		go func() {
			for {
				select {
				case resp := <-networkEventsChan:
					// Filter for API calls related to your extension
					fmt.Printf("response from the networkEventsChan is %v \n\n", resp)
					println("the response url is ->", resp.Response.URL)
					if strings.Contains(resp.Response.URL, "sponsorblock") ||
						strings.Contains(resp.Response.URL, "your-api-endpoint") {

						fmt.Printf("Detected relevant API call: %s\n", resp.Response.URL)

						// Get the response body
						var responseBody string
						err := chromedp.Run(swCtx, chromedp.ActionFunc(func(ctx context.Context) error {
							body, err := network.GetResponseBody(resp.RequestID).Do(ctx)
							if err != nil {
								return err
							}
							responseBody = string(body)
							return nil
						}))
						if err != nil {
							fmt.Printf("Error getting response body: %v\n", err)
							continue
						}

						fmt.Printf("Response body: %s\n", responseBody)

						// Parse the response and check for sponsor segments
						// ... your code to handle the response ...
					}
				case <-time.After(30 * time.Second):
					// Timeout if no relevant responses are detected
					fmt.Printf("No relevant API responses detected for URL %s\n", pageUrl)
					return
				}
			}
		}()

		// 7. Enable Network on that session.

		// 8. Listen for network requests from the service worker.
		// chromedp.ListenTarget(ctx, func(ev interface{}) {
		// 	if req, ok := ev.(*network.EventRequestWillBeSent); ok && req.SessionID == sessionID {
		// 		fmt.Printf("→ SW request: %s\n", req.Request.URL)
		// 	}
		// })
		println("the service workrer target id is ->", swTargetID)
		println("sleeping for 4 sec to ensure that the we are able to intercept the message form the service worker")
		time.Sleep(time.Second * 4)
		err := chromeExtension.NavigateToAWebPage(ctx, pageUrl)
		if err != nil {
			println("there is a error in navigating to the youtubeUrl at index ", i)
			continue
		}
		println("we are in the youtube video ->", pageUrl)
		stopChannelToStopChekingIfTheVideoIsPlaying := make(chan struct{})
		defer close(stopChannelToStopChekingIfTheVideoIsPlaying) // this is a send only channel so only we can close it
		resultChanForTrackingPlayBackTime := make(chan commonchanneltype.GenericResultChannel[*[]float64])
		go chromeExtension.EnsureVideoIsPlayingPeriodically(ctx, time.Millisecond*700, stopChannelToStopChekingIfTheVideoIsPlaying, false)
		go chromeExtension.TrackVideoPlaybackTime(ctx, resultChanForTrackingPlayBackTime)
		playBackTimeChan := <-resultChanForTrackingPlayBackTime
		println("the playBakc result is here")
		if playBackTimeChan.Err != nil {
			t.Fatal("there is a error in getting the playBackTime[] and it is ->" + playBackTimeChan.Err.Error())
		}
		println("is the playBackTime array not nil->", playBackTimeChan.Result != nil, " and the array lenght is:", len(*playBackTimeChan.Result))
		APIResponseFormNetwork := <-getAPIResponseFromNetworkChann
		if APIResponseFormNetwork.Err != nil {
			println("there is a error in getting the api respons form the NEtwork and it is ->", APIResponseFormNetwork.Err.Error())
			t.Fatal(APIResponseFormNetwork.Err)
		}
		fmt.Printf("we got the API resp form the network and it is -> %v", APIResponseFormNetwork.Result)

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
