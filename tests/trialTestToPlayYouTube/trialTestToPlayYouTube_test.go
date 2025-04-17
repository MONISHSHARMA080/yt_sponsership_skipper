package trialtesttoplayyoutube

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func TestToWatchYoutbeVideo(t *testing.T) {
	// Define browser options with stealth settings
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "IsolateOrigins,site-per-process"),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-infobars", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// Add common cookies and headers
	err := chromedp.Run(ctx, network.Enable())
	if err != nil {
		log.Fatal(err)
	}

	// Run tasks
	err = chromedp.Run(ctx,
		// Set realistic viewport
		emulation.SetDeviceMetricsOverride(1920, 1080, 1.0, false),
		// Clear JavaScript properties that might reveal automation
		chromedp.Evaluate(`
			// Overwrite navigator properties
			Object.defineProperty(navigator, 'webdriver', {
				get: () => false,
			});
			// Clear automation-related properties
			delete window.cdc_adoQpoasnfa76pfcZLmcfl_Array;
			delete window.cdc_adoQpoasnfa76pfcZLmcfl_Promise;
			delete window.cdc_adoQpoasnfa76pfcZLmcfl_Symbol;
		`, nil),
		// Load YouTube
		chromedp.Navigate("https://www.youtube.com/watch?v=dQw4w9WgXcQ"),
		// Wait for video player
		chromedp.WaitVisible("video", chromedp.ByQuery),
		// Give video time to start
		chromedp.Sleep(5*time.Second),
		// Check if video is playing
		chromedp.Evaluate(`
			(() => {
				const video = document.querySelector('video');
				return video && !video.paused && video.currentTime > 0;
			})()
		`, nil),
		// Keep browser open for visual confirmation
		chromedp.Sleep(30*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Minute * 10)

	// Navigate to the target URL
}
