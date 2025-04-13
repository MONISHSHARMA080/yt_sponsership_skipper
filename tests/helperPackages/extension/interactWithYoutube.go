package extension

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func (e *ChromeExtension) IfThereIsAAdThenFinishIt(ctx context.Context) error {
	areWeInAAd, err := e.IsTheYoutubeVideoPlayingAnAd(ctx)
	if err != nil {
		return err
	}
	if !areWeInAAd {
		println("there is not add in the youtube video")
		return nil
	}
	println("there is a add in the youtube video")
	// if the add go more than 4 min then we will fail as a ad will not be more than 4 min
	err = e.WaitForAdToFinish(ctx, time.Minute*4)
	if err != nil {
		return err
	}
	return nil
}

// note: run this func after the page visible one as this will not make sure that the page is loaded
func (e *ChromeExtension) IsTheYoutubeVideoPlayingAnAd(ctx context.Context) (bool, error) {
	adPlayerClass := "ytp-ad-player-overlay-layout"
	// assume the page is loaded , now we are going to see if the add one is visible
	return e.IsTheElementVisible(adPlayerClass, ctx)
}

// WaitForAdToFinish waits until a YouTube ad finishes playing and returns
// It takes a context and a timeout duration as parameters
// Returns nil if the ad finishes successfully, or an error if timeout occurs or other issues
func (e *ChromeExtension) WaitForAdToFinish(ctx context.Context, timeout time.Duration) error {
	// Create a context with the specified timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Create a ticker to check for ad status periodically
	ticker := time.NewTicker(600 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check if the ad is still playing
			isAdPlaying, err := e.IsTheYoutubeVideoPlayingAnAd(ctx)
			if err != nil {
				return err
			}

			// If ad is no longer playing, return success
			if !isAdPlaying {
				return nil
			}

		case <-timeoutCtx.Done():
			// Context timeout or cancellation
			return fmt.Errorf("timeout waiting for ad to finish")
		}
	}
}

// note makse sure to close the channel form the outside
func (e *ChromeExtension) EnsureVideoIsPlayingPeriodically(ctx context.Context, intervalToCheckTheVideoPlayingAt time.Duration, stopChan <-chan struct{}, shouldWeStopOnError bool) {
	ticker := time.NewTicker(intervalToCheckTheVideoPlayingAt)
	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := e.MakeSureTheVideoIsPlaying(ctx)
				if err != nil {
					fmt.Printf("Error in ensuring video periodically is : %v", err)
					if shouldWeStopOnError {
						println("returning as the making sure the video is playing func returend a error")
						return
					}
				}
			case <-stopChan:
				fmt.Println("Stopping periodic video check")
				return
			}
		}
	}()
}

// if the video is not playing then we will play it, if it is then return
func (e *ChromeExtension) MakeSureTheVideoIsPlaying(ctx context.Context) error {
	var result string
	script := `
  (function() {
    try {
      const video = document.querySelector('video');
      if (!video) {
        return "No video element found";
      }
      
      // Check if video is paused
      if (video.paused) {
        // Try to click the play button
        const playButton = document.querySelector('.ytp-play-button');
        if (playButton) {
          playButton.click();
        } else {
          // Direct play attempt if button not found
          video.play();
        }
        
        // Wait briefly and check if it's playing now
        return new Promise(resolve => {
          setTimeout(() => {
            if (video.paused) {
              resolve("Video still paused after play attempt");
            } else {
              resolve(null);
            }
          }, 1000);
        });
      } else {
        return null; // Already playing
      }
    } catch(e) {
      return "Error: " + e.toString();
    }
  })()
  `
	err := chromedp.Run(ctx,
		chromedp.Evaluate(script, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}),
	)
	return err
}
