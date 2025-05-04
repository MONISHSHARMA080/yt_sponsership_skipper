package extension

import (
	"context"
	"fmt"
	"time"
	commonchanneltype "youtubeAdsSkipper/tests/helperPackages/CommonChannelType"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

func (ce *ChromeExtension) DidWeSkippedTheAd(startTime, endTime float64, playbackTime []float64) (bool, error) {
	// check if the start time and end time is in the playback time
	if len(playbackTime) == 0 {
		return false, fmt.Errorf("the playback time array is empty")
	}
	// Look for a significant jump in playback time between startTime and endTime
	for i := 1; i < len(playbackTime); i++ {
		// timeDiff := playbackTime[i] - playbackTime[i-1]

		// Check if the jump occurred around our target segment
		if playbackTime[i-1] >= startTime-1 && playbackTime[i] <= endTime+1 {
			return true, nil
		}
	}

	// No significant skip found in the target range
	return false, nil
}

// will store the time in Array in the js side and when the vide is completed, then we will return it
// this will be a struct that will have result and error on it
func (ce *ChromeExtension) TrackVideoPlaybackTime(ctx context.Context, resultChannel chan commonchanneltype.GenericResultChannel[*[]float64]) {
	var result struct {
		Result []float64 `json:"result"`
		Error  string    `json:"error"`
	}
	script := `
async function getTheVideoPlaybackTimeForItsLife() {
  try {
    const videoPlayer = document.querySelector('video');
    if (!videoPlayer) {
      return { result: null, error: "The video player is not found" };
    }

  console.log("in the script form the chromedp")
    let resultArray = [];
    let i = 0;

    videoPlayer.addEventListener("timeupdate", () => {
      resultArray.push(videoPlayer.currentTime);
      console.log("Adding currentTime to array, index:" + i);
      i++
    });

    return await new Promise((resolve) => {
      videoPlayer.addEventListener("ended", () => {
        console.log("Video has ended, returning resultArray");
        resolve({ result: resultArray, error: "" });
      });
    });
  } catch (error) {
    console.log("Error:", error);
    return { result: null, error: "There was an error: " + error };
  }
}
getTheVideoPlaybackTimeForItsLife()
  `
	err := chromedp.Run(ctx,
		chromedp.Evaluate(script, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}),
	)
	println("in the track video playback time func")
	if err != nil {
		resultChannel <- commonchanneltype.GenericResultChannel[*[]float64]{Result: nil, Err: fmt.Errorf("there is a error in running the script to get the videoPlayback time (array) out of the js exec and it is :->%s", err.Error())}
		return
	}
	println("the results are here for the time array")
	resultChannel <- commonchanneltype.GenericResultChannel[*[]float64]{Result: &result.Result, Err: nil}
}

// checks if we are still in the ad using ticker, if the StopCheckingOnError is true them we will return form the func on single eror
// Prefer keeping stopCheckingOnError false
// if we get the value of areWeInAAd to  be true we will still check after 5 sec to see that the ad is over(confirm)
func (e *ChromeExtension) IfThereIsAAdThenFinishIt(ctx context.Context, intervalToKeepCheckingIfWeAreStillInAAD time.Duration, stopCheckingOnError bool) []error {
	ticker := time.NewTicker(intervalToKeepCheckingIfWeAreStillInAAD)
	defer ticker.Stop()
	errorsCollectedInCheckingIFWeAreInAAD := []error{}
	sleepAndCheckForADAgain := time.Second * 3
	// this is here so that if there are 2 ads then we need to skip both of them
	ifThisTheFirstTimeADIsOver := true
	done := make(chan struct{})

	go func() {
		// defer close(done) // Signal when goroutine completes
		defer println("the go func to keep checking for checking if we in a AD has closed")
		for range ticker.C {
			println("checking if the ad is playing   ")
			areWeInAAd, err := e.IsTheYoutubeVideoPlayingAnAd(ctx)
			println("got the value form id the AD is over func and it is ->", areWeInAAd, " and the err is nil->", err == nil)
			if err != nil {
				println("there is a error in checking if the AD is palying and it is ->", err.Error())
				errorsCollectedInCheckingIFWeAreInAAD = append(errorsCollectedInCheckingIFWeAreInAAD, err)
				if stopCheckingOnError {
					close(done)
					return
				}
				// else keep going
				continue
			}
			println("are we still in a AD ->", areWeInAAd)
			if !areWeInAAd {
				println("this was the first time when the ad was over and we are resetting the ticker ")
				// now wait for 5 sec to check if there is not additional add
				if ifThisTheFirstTimeADIsOver {
					// we will wait and check for the Second AD to finish/not appear by waiting for 4 sec
					println("this is the first time the AD has finish we will sleep for some time and then re check it again to be for sure")
					ifThisTheFirstTimeADIsOver = false
					time.Sleep(sleepAndCheckForADAgain)
					continue
				} else {
					// if not then there is no ad and we will return
					println("there is not ad after the second time waiting for the ad  ")
					close(done)
					return
				}
			} else {
				println("--AD still running --")
			}
		}
	}()
	<-done
	println("the for loop over the tiker finished")

	return nil
}

// note: run this func after the page visible one as this will not make sure that the page is loaded
func (e *ChromeExtension) IsTheYoutubeVideoPlayingAnAd(ctx context.Context) (bool, error) {
	script := `
      (function(){
          const player = document.querySelector('.html5-video-player');
          return player ? player.classList.contains('ad-showing') : false;
      })()
    `
	var isAd bool
	err := chromedp.Run(ctx,
		chromedp.Evaluate(script, &isAd),
	)
	return isAd, err
}

// note makse sure to close the channel form the outside
func (e *ChromeExtension) EnsureVideoIsPlayingPeriodically(ctx context.Context, intervalToCheckTheVideoPlayingAt time.Duration, stopChan <-chan struct{}, shouldWeStopOnError bool) {
	ticker := time.NewTicker(intervalToCheckTheVideoPlayingAt)
	println("in the ensure video is playing func")
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
				let vp = document.querySelector("#ytd-player");
				if (vp !== null && vp.player_ && vp.player_.isReady()) {
					// Get current time and duration to calculate percentage
					const currentTime = vp.getPlayer().getCurrentTime();
					const duration = vp.getPlayer().getDuration();
					const percentagePlayed = (currentTime / duration * 100).toFixed(2);
					
					
					// Set up ended event listener if not already set
					const video = document.querySelector('video');
					if (video) {
						// Only add the listener if it hasn't been added before
						if (!video.hasAttribute('data-end-listener-added')) {
							video.addEventListener('ended', function() {
								console.log("Video ended event triggered");
								vp.getPlayer().pauseVideo();
								video.pause();
							});
							// Mark that we've added the listener to avoid duplicates
							video.setAttribute('data-end-listener-added', 'true');
							console.log("Added ended event listener to video");
						}
						
						video.muted = true;
						video.play();
					} else {
						console.log("No video element found");
					}
					
					// Only play if not completed
					vp.getPlayer().playVideo();
					vp.getPlayer().setPlaybackRate(2);
					
					console.log("is my video paused", video ? video.paused : "no video element");
					
					return "Video playing successfully - " + percentagePlayed + "% complete ------current time of video is -> "+ currentTime + " and the duration is -> " + duration;
				} else {
					console.log("Video player not ready or not found:", vp);
					return "Video player not ready";
				}
			} catch (err) {
				console.error("Error ensuring video is playing:", err);
				return "Error: " + err.toString();
			}
		})()

	 `
	err := chromedp.Run(ctx,
		chromedp.Evaluate(script, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}),
	)
	println("the result of the script to keep the video running is  ->", result)
	return err
}
