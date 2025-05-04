package helperfunc1_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func GetNewBrowserForChromeExtension(extensionID string) (context.Context, context.CancelFunc, context.CancelFunc, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, nil, nil, err
	}
	
	extensionPath := filepath.Join(cwd, "../chromeExtension")
	println("the chrome extensionPath is ->", extensionPath)
	
	ublockPath := filepath.Join(cwd, "u-block")
	println("the ublockPath extensionPath is ->", ublockPath)

	// Configure options for non-headless Chrome in Docker
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),            // Explicitly set non-headless mode
		chromedp.Flag("no-sandbox", true),           // Required in Docker
		chromedp.Flag("disable-dev-shm-usage", true), // Avoid using /dev/shm which is limited in Docker
		chromedp.Flag("remote-debugging-address", "0.0.0.0"), // Allow external connections
		chromedp.Flag("remote-debugging-port", "9222"),
		chromedp.Flag("disable-gpu", false),         // Enable GPU for better rendering
		chromedp.Flag("window-size", "1280,720"),
		chromedp.Flag("enable-automation", false),   // Disable the automation banner
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("extensions-on-chrome-urls", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("load-extension", extensionPath+","+ublockPath),
		chromedp.Flag("disable-extensions-except", extensionPath+","+ublockPath),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-signin-promo", true),
		chromedp.Flag("disable-sync-credential-backend", true),
		chromedp.Flag("enable-features", "PasswordImport"),
		chromedp.Flag("disable-features", "PasswordExport,Signin,IdentityManager"),
		chromedp.Flag("safebrowsing-disable-extension-blacklist", true),
		chromedp.Flag("account-consistency", "disabled"),
		chromedp.Flag("unlimited-storage", true),
		chromedp.Flag("allow-unlimited-local-storage", true),
		// Use a persistent user data directory in Docker
		chromedp.UserDataDir("/tmp/chrome-profile"),
		// For better automation detection evasion
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("disable-infobars", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"),
		// Network settings
		chromedp.Flag("enable-network-service", true),
		chromedp.WSURLReadTimeout(60*time.Second),
	)

	// Create execution allocator context
	allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)
	
	// Create Chrome context
	ctx, _ := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(func(format string, args ...interface{}) {
			fmt.Printf(format, args...)
		}),
		// Increase browser start timeout for Docker environment
		// chromedp.WithBrowserOption(
		// 	chromedp.WithBrowserStartTimeout(2 * time.Minute),
		// ),
	)

	// Create a timeout for the entire operation
	ctx, cancel := context.WithTimeout(ctx, 40*time.Minute)

	// Enable network events
	err = chromedp.Run(ctx, network.Enable())
	if err != nil {
		cancel()
		cancel1()
		return nil, nil, nil, err
	}

	// Set viewport and anti-automation measures
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
	)
	if err != nil {
		cancel()
		cancel1()
		return nil, nil, nil, err
	}

	return ctx, cancel, cancel1, nil
}