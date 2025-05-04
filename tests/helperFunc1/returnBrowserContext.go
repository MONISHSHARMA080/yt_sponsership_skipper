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

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// Core settings for non-headless mode
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),
		chromedp.Flag("no-sandbox", true),             // Required for Docker
		chromedp.Flag("disable-dev-shm-usage", true),  // Good for Docker stability
		chromedp.Flag("disable-gpu", true),            // Better for Docker
		chromedp.Flag("enable-automation", false),     // Disable the automation banner
		
		// Extension-related flags
		chromedp.Flag("load-extension", extensionPath+","+ublockPath),
		chromedp.Flag("disable-extensions-except", extensionPath+","+ublockPath),
		chromedp.Flag("extensions-on-chrome-urls", true),
		chromedp.Flag("disable-popup-blocking", true),
		
		// Security and permission settings
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("disable-infobars", true),
		
		// Storage settings
		chromedp.Flag("unlimited-storage", true),
		chromedp.Flag("allow-unlimited-local-storage", true),
		
		// Account/sync settings
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-signin-promo", true),
		chromedp.Flag("disable-sync-credential-backend", true),
		chromedp.Flag("account-consistency", "disabled"),
		chromedp.Flag("disable-features", "PasswordExport,Signin,IdentityManager"),
		chromedp.Flag("enable-features", "PasswordImport,DisableAccountConsistency"),
		
		// Viewport settings
		chromedp.Flag("window-size", "1920,1080"),
		
		// Debug settings
		chromedp.Flag("remote-debugging-port", "9222"),
		chromedp.Flag("auto-open-devtools-for-tabs", true),
		chromedp.Flag("enable-network-service", true),
		
		// Anti-detection settings
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"),
		
		// Networking timeout
		chromedp.WSURLReadTimeout(60*time.Second),
	)

	allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)

	// Create context with more verbose logging
	ctx, _ := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(func(format string, args ...interface{}) {
			fmt.Printf("ChromeDP: "+format+"\n", args...)
		}),
		chromedp.WithDebugf(func(format string, args ...interface{}) {
			fmt.Printf("ChromeDP Debug: "+format+"\n", args...)
		}),
	)

	// Create a timeout for the entire operation
	ctx, cancel := context.WithTimeout(ctx, 40*time.Minute)

	// Set up network
	err = chromedp.Run(ctx, network.Enable())
	if err != nil {
		cancel()
		cancel1()
		return nil, nil, nil, fmt.Errorf("failed to enable network: %w", err)
	}
	
	// Set up anti-detection measures
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
		return nil, nil, nil, fmt.Errorf("failed to set up anti-detection: %w", err)
	}

	return ctx, cancel, cancel1, nil
}