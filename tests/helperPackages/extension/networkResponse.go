package extension

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	commonchanneltype "youtubeAdsSkipper/tests/helperPackages/CommonChannelType" // Adjust import path if needed
	"youtubeAdsSkipper/tests/helperPackages/extension/types"                     // Adjust import path if needed

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/serviceworker"
	"github.com/chromedp/chromedp"
)

// GetResponseFromServerToChromeExtension listens specifically to the network requests
// made BY the chrome extension identified by ce.ExtensionId and sends the parsed
//
// Note: you should check to see wether this func has reuslted in the result, if not then we probally have soem error and you should proceed as error and return
func (ce ChromeExtension) GetResponseFromServerToChromeExtension(ctx context.Context, timeToKeepLookingForNetworkResponse time.Duration, resultChannel chan commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]) {
	// Create a timeout context for the entire operation, including finding the target
	chromePopulUrl := "chrome-extension://dpkehfkmkhmbmbofjaikccklcdpdmhpl/index.html"
	newTabCtx, cancelFuncTab, CancelFuncTimeout, err := ce.getNewTab(ctx, timeToKeepLookingForNetworkResponse, chromePopulUrl)
	if err != nil {
		println("there is a error in creatign new tab ->", err.Error())
		resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{Err: newTabCtx.Err()}
		return
	}
	defer cancelFuncTab()
	defer CancelFuncTimeout()
	responeFromYoutubeApiSuccessfullyReceivedChan := make(chan *types.YouTubeVideoResponse)
	println("listening to the network req form the service worker")

	// get the network response form it

	chromedp.ListenTarget(newTabCtx, func(ev any) {
		switch event := ev.(type) {
		case *network.EventResponseReceived:
			if strings.Contains(event.Response.URL, ce.ExtensionId) {
				println("event response received and url is -> ", event.Response.URL)
			}
			if event.Response.URL != "" {
				go func(reqID network.RequestID) {
					responseBody, err := network.GetResponseBody(reqID).Do(newTabCtx)
					if err != nil {
						// not exiting , reason being that the outer function will check and if this has not resulted in the response then we declare failure
						println("there is a error in getting the error body form the network request (accessing it and not parsing it, and we are not exicting out go func ) ->", err.Error())
						return
					}
					var jsonResponseOnYoutubePath types.YouTubeVideoResponse
					err = json.Unmarshal(responseBody, &jsonResponseOnYoutubePath)
					if err != nil {
						println("there is a error in marshalling the json received from the /youtube path in the youtube struct, and err->", err.Error())
						return
					}
					responeFromYoutubeApiSuccessfullyReceivedChan <- &jsonResponseOnYoutubePath
				}(event.RequestID)
			}
		}
	})

	err = chromedp.Run(newTabCtx, network.Enable(),
		network.SetBypassServiceWorker(false),
		network.SetCacheDisabled(true),
		serviceworker.Enable(),
	)
	if err != nil {
		resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{Err: err, Result: nil}
		return
	}

	time.Sleep(time.Minute * 2)
	select {
	case <-ctx.Done():
		println("Context already canceled:", ctx.Err().Error())
		resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{Err: ctx.Err()}
		return
	case <-newTabCtx.Done():
		println("Context already canceled:", newTabCtx.Err().Error())
		resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{Err: newTabCtx.Err()}
		return
	case ytBodyResponse := <-responeFromYoutubeApiSuccessfullyReceivedChan:
		resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{Result: ytBodyResponse}
		return
	}
}

func (ce ChromeExtension) getNewTab(ctx context.Context, timeToKeepLookingForNetworkResponse time.Duration, urlToOpen string) (newTabCtx context.Context, cancelTabFunc, cancelTimeoutFunc context.CancelFunc, err error) {
	newTabCtx, cancelFunc1 := chromedp.NewContext(ctx)
	newTabCtx, cancelFunc2 := context.WithTimeout(newTabCtx, timeToKeepLookingForNetworkResponse)
	err = chromedp.Run(newTabCtx, chromedp.Navigate(urlToOpen))
	return newTabCtx, cancelFunc1, cancelFunc2, err
}
