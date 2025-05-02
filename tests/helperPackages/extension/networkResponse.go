package extension

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	commonchanneltype "youtubeAdsSkipper/tests/helperPackages/CommonChannelType" // Adjust import path if needed
	"youtubeAdsSkipper/tests/helperPackages/extension/types"                     // Adjust import path if needed

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/serviceworker"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

// GetResponseFromServerToChromeExtension listens specifically to the network requests
// made BY the chrome extension identified by ce.ExtensionId and sends the parsed
//
// Note: you should check to see wether this func has reuslted in the result, if not then we probally have soem error and you should proceed as error and return
func (ce ChromeExtension) GetResponseFromServerToChromeExtension(ctx context.Context, timeToKeepLookingForNetworkResponse time.Duration, resultChannel chan commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]) {
	println("getting context without timepout")
	ctx, cancelFunc := context.WithTimeout(ctx, timeToKeepLookingForNetworkResponse)
	// if err != nil {
	//    println("the error in getting the ")
	// 	resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{Result: nil, Err: nil}
	// 	return
	// }
	defer cancelFunc()
	var infos []*target.Info
	if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		infos, err = target.GetTargets().Do(ctx)
		return err
	})); err != nil {
		println(" we have error in getting the targets ->", err.Error())
		resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{Result: nil, Err: nil}
		return
	}

	// 5. Find your extension’s service worker.
	var swTargetID target.ID
	prefix := "chrome-extension://" + ce.ExtensionId
	for i, ti := range infos {
		println("index:", i, " and the ti title is ->", ti.Title, "ti's url is ->", ti.URL, "--and type is ", ti.Type)

		if ti.Type == "service_worker" && strings.HasPrefix(ti.URL, prefix) {
			println("found the service worker with the id ->", ti.TargetID)
			swTargetID = ti.TargetID
			break
		}
	}
	if swTargetID == "" {
		resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{Result: nil, Err: fmt.Errorf("the taerget ID of the service worker is empty")}
		return
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
		resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{Result: nil, Err: fmt.Errorf("failed to enable network monitoring: %v", err)}
		return
	}
	type responseForJsonBodyChannel struct {
		ytResp  *types.YouTubeVideoResponse
		err     error
		success bool
	}
	resultChanForJsonBody := make(chan responseForJsonBodyChannel)
	var youtubePathResponseInJson *types.YouTubeVideoResponse
	// Process network events with a timeout
	go func() {
		for {
			resp, ok := <-networkEventsChan
			if !ok {
				err := fmt.Errorf("the result is not ok in receiving form the channel 0 ")
				resultChanForJsonBody <- responseForJsonBodyChannel{err: err, success: false, ytResp: nil}
				return
			}
			// Filter for API calls related to your extension
			fmt.Printf("response from the networkEventsChan is %+v \n\n", resp)
			fmt.Printf("response is -> %+v \n\n", *resp.Response)
			println("the response url is ->", resp.Response.URL)
			if strings.Contains(resp.Response.URL, "localhost") ||
				strings.Contains(resp.Response.URL, "/youtubeVideo") {

				fmt.Printf("Detected relevant API call: %s\n", resp.Response.URL)

				// Get the response body
				var responseBody string
				err := chromedp.Run(swCtx, chromedp.ActionFunc(func(ctx context.Context) error {
					body, err := network.GetResponseBody(resp.RequestID).Do(ctx)
					if err != nil {
						println("error Unmarshalling the json in body")
						resultChanForJsonBody <- responseForJsonBodyChannel{err: err, success: false, ytResp: nil}
						return err
					}
					responseBody = string(body)
					println("the response body is ->\n ", responseBody, " \n\n---")
					err = json.Unmarshal(body, &youtubePathResponseInJson)
					if err != nil {
						println("there is a error in parsing the json into youtube path response ->", err.Error())
						resultChanForJsonBody <- responseForJsonBodyChannel{err: err, success: false, ytResp: nil}
						return err
					}
					println("got the json decoded in the struct and returnign it")
					resultChanForJsonBody <- responseForJsonBodyChannel{err: nil, success: true, ytResp: youtubePathResponseInJson}
					return nil
				}))
				if err != nil {
					fmt.Printf("Error getting response body: %v\n", err)
					continue
				}

				fmt.Printf("Response body: %s\n", responseBody)

			}
		}
	}()

	// make the resultChanForJsonBody to be a struct that has err and bool that way only one palce will write to the resultChannel
	// resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{Result: nil, Err: nil}
	select {
	case res := <-resultChanForJsonBody:
		// like what am I supose to do here except for return
		// that's why make this channle into a struct with error , ptr to body struct and bool
		if res.success {
			resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{Result: res.ytResp, Err: nil}
		} else {
			resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{Result: nil, Err: res.err}
		}
		return
	case <-ctx.Done():
		println("context deadline reached so we will just quit")
		resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{Result: nil, Err: fmt.Errorf("ctx deadline reached")}
	}
}

func (ce ChromeExtension) getNewTab(ctx context.Context, timeToKeepLookingForNetworkResponse time.Duration, urlToOpen string) (newTabCtx context.Context, cancelTabFunc, cancelTimeoutFunc context.CancelFunc, err error) {
	newTabCtx, cancelFunc1 := chromedp.NewContext(ctx)
	newTabCtx, cancelFunc2 := context.WithTimeout(newTabCtx, timeToKeepLookingForNetworkResponse)
	err = chromedp.Run(newTabCtx, chromedp.Navigate(urlToOpen),

		network.Enable(),
		network.SetBypassServiceWorker(false),
		network.SetCacheDisabled(true),
		serviceworker.Enable(),
		network.SetExtraHTTPHeaders(map[string]interface{}{
			"X-DevTools-Emulate-Network-Conditions-Client-Id": "Chrome.Headless",
		}),
		target.SetDiscoverTargets(true),
		target.SetAutoAttach(true, false),
	)

	return newTabCtx, cancelFunc1, cancelFunc2, err
}
