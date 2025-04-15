package extension

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// checks if we are still in the ad using ticker, if the StopCheckingOnError is true them we will return form the func on single eror
// Prefer keeping stopCheckingOnError false
// if we get the value of areWeInAAd to  be true we will still check after 5 sec to see that the ad is over(confirm)
func (e *ChromeExtension) IfThereIsAAdThenFinishIt(ctx context.Context, intervalToKeepCheckingIfWeAreStillInAAD time.Duration, stopCheckingOnError bool) []error {
	ticker := time.NewTicker(intervalToKeepCheckingIfWeAreStillInAAD)
	defer ticker.Stop()
	errorsCollectedInCheckingIFWeAreInAAD := []error{}
	firstTimeTheADIsOver := true
	// recheckForAdAfterFirstTimeADFinish := time.Second * 3
	done := make(chan struct{})

	go func() {
		defer close(done) // Signal when goroutine completes
		for range ticker.C {
			println("checking if the ad is playing and is this the first time the AD is over -> ", firstTimeTheADIsOver)
			areWeInAAd, err := e.IsTheYoutubeVideoPlayingAnAd(ctx)
			println("got the value form id the AD is over func and it is ->", areWeInAAd, " and the err is nil->", err == nil)
			if err != nil {
				println("there is a error in checking if the AD is palying and it is ->", err.Error())
				errorsCollectedInCheckingIFWeAreInAAD = append(errorsCollectedInCheckingIFWeAreInAAD, err)
				if stopCheckingOnError {
					return
				}
				// else keep going
				continue
			}
			println("are we still in a AD ->", areWeInAAd)
			if !areWeInAAd {
				if firstTimeTheADIsOver {
					// of this is the first thme ad is over then we will set it to false
					firstTimeTheADIsOver = false
					println("this was the first time when the ad was over and we are resetting the ticker ")
					// now wait for 5 sec to check if there is not additional add
					// ticker.Reset(recheckForAdAfterFirstTimeADFinish)
					// now sleep for 5 sec and then run the func  again
					time.Sleep(time.Second * 4)

				} else {
					// if  firstTimeTheADIsOver is false then we have spend our 5 sec to see if there is no additional AD, and now we exit
					println(" this was the 2nd time that the ads skipped value was true ")
					return
				}
				// IDK why is this for, this is causing the ticker to return
				// return
			} else {
				println("we are still in a AD")
			}
		}
	}()
	<-done
	println("the for loop over the tiker finished")

	// areWeInAAd, err := e.IsTheYoutubeVideoPlayingAnAd(ctx)
	// println("checking if the ad is playing ")
	// if err != nil {
	// 	return err
	// }
	// println("is the video playing an ad ->", areWeInAAd)
	// if !areWeInAAd {
	// 	println("there is not add in the youtube video")
	// 	return nil
	// }
	// println("there is a add in the youtube video")
	// if the add go more than 4 min then we will fail as a ad will not be more than 4 min
	// err := e.WaitForAdToFinish(ctx, time.Minute*4)
	// if err != nil {
	// 	return err
	// }
	return nil
}

// note: run this func after the page visible one as this will not make sure that the page is loaded
func (e *ChromeExtension) IsTheYoutubeVideoPlayingAnAd(ctx context.Context) (bool, error) {
	adPlayerClass := ".ytp-ad-player-overlay-layout"
	// assume the page is loaded , now we are going to see if the add one is visible
	return e.IsTheElementVisible(adPlayerClass, ctx)
}

// --------removed it as we do not need it and the above(3rd ish) method is bettor and not AI generated and suits my need

// WaitForAdToFinish waits until a YouTube ad finishes playing and returns
// It takes a context and a timeout duration as parameters
// Returns nil if the ad finishes successfully, or an error if timeout occurs or other issues
// func (e *ChromeExtension) WaitForAdToFinish(ctx context.Context, timeout time.Duration) error {
// 	// Create a context with the specified timeout
// 	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
// 	defer cancel()
//
// 	// Create a ticker to check for ad status periodically
// 	ticker := time.NewTicker(600 * time.Millisecond)
// 	defer ticker.Stop()
//
// 	for {
// 		select {
// 		case <-ticker.C:
// 			// Check if the ad is still playing
// 			isAdPlaying, err := e.IsTheYoutubeVideoPlayingAnAd(ctx)
// 			if err != nil {
// 				return err
// 			}
//
// 			// If ad is no longer playing, return success
// 			if !isAdPlaying {
// 				return nil
// 			}
//
// 		case <-timeoutCtx.Done():
// 			// Context timeout or cancellation
// 			return fmt.Errorf("timeout waiting for ad to finish")
// 		}
// 	}
// }

// note makse sure to close the channel form the outside
func (e *ChromeExtension) EnsureVideoIsPlayingPeriodically(ctx context.Context, intervalToCheckTheVideoPlayingAt time.Duration, stopChan <-chan struct{}, shouldWeStopOnError bool) {
	ticker := time.NewTicker(intervalToCheckTheVideoPlayingAt)
	defer ticker.Stop()
	go func() {
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
