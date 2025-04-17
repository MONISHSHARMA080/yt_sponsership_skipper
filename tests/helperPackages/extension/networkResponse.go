package extension

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	commonchanneltype "youtubeAdsSkipper/tests/helperPackages/CommonChannelType"
	"youtubeAdsSkipper/tests/helperPackages/extension/types"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func (ce *ChromeExtension) GetResponseFromServerToChromeExtension(ctx context.Context, timeToKeepLookingForNettworkResponse time.Duration, resultChannel chan commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]) {
	// var req network.RequestID
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeToKeepLookingForNettworkResponse)
	defer cancel()
	responseChan := make(chan *types.YouTubeVideoResponse, 1)

	err := chromedp.Run(ctxWithTimeout, network.Enable())
	if err != nil {
		resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{
			Result: nil,
			Err:    fmt.Errorf("error enabling network monitoring: %w", err),
		}
		return
	}

	chromedp.ListenTarget(ctxWithTimeout, func(ev interface{}) {
		switch resp := ev.(type) {
		case *network.EventResponseReceived:
			// Check if this is a response from our endpoint
			println("received a network response and it is on url ->", resp.Response.URL)
			if resp.Response.URL != "" && resp.Response.Status > 0 {
				// Look for the target endpoint in the URL
				if strings.Contains(resp.Response.URL, "youtubeVideo") {
					log.Printf("Found response from endpoint: %s, status: %d", resp.Response.URL, resp.Response.Status)

					// Get the response body
					rbp, err := network.GetResponseBody(resp.RequestID).Do(ctxWithTimeout)
					if err != nil {
						log.Printf("Error getting response body: %v", err)
						return
					}

					// Parse the response
					var response types.YouTubeVideoResponse
					if err := json.Unmarshal(rbp, &response); err != nil {
						log.Printf("Error parsing response JSON: %v", err)
						return
					}

					// Send the response to the channel
					responseChan <- &response
				}
			}
		}
	})

	// Wait for response or timeout
	select {
	case response := <-responseChan:
		resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{
			Result: response,
			Err:    nil,
		}
	case <-ctxWithTimeout.Done():
		resultChannel <- commonchanneltype.GenericResultChannel[*types.YouTubeVideoResponse]{
			Result: nil,
			Err:    fmt.Errorf("timeout waiting for response from /youtubeVideo endpoint after 30 seconds"),
		}
	}
}
