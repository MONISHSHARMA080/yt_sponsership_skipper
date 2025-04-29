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
	// var opts []chromedp.ExecAllocatorOption
	cwd, err := os.Getwd()
	if err != nil {
		return nil, nil, nil, err
	}
	extensionPath := filepath.Join(cwd, "../chromeExtension")
	println("the chrome extensionPath is ->", extensionPath)
	ublockPath := filepath.Join(cwd, "u-block")
	println("the ublockPath extensionPath is ->", ublockPath)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],

		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),          // Important: We don't want headless mode
		chromedp.Flag("enable-automation", false), // Disable the automation banner
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("extensions-on-chrome-urls", true),
		// chromedp.Flag("disable-extensions-except", extensionPath+","+ublockPath),

		chromedp.Flag("no-first-run", true),
		// Disable default apps and extensions
		chromedp.Flag("disable-default-apps", true),
		// chromedp.Flag("disable-extensions-except", extensionID),
		chromedp.Flag("load-extension", extensionPath),

		chromedp.Flag("use-gl", "desktop"),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("enable-features", "DisableAccountConsistency"),
		chromedp.Flag("disable-web-security", true),
		// chromedp.Flag("load-extension", extensionPath),
		chromedp.Flag("load-extension", extensionPath+","+ublockPath),            // Load both extensions
		chromedp.Flag("disable-extensions-except", extensionPath+","+ublockPath), // Optional: disable other extensions
		// chromedp.Flag("disable-extensions-except", extensionPath),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-sync", true),
		//---
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-signin-promo", true),
		chromedp.Flag("disable-sync-credential-backend", true),
		chromedp.Flag("enable-features", "PasswordImport"),
		chromedp.Flag("disable-features", "PasswordExport,Signin"),
		chromedp.Flag("safebrowsing-disable-extension-blacklist", true),

		chromedp.Flag("account-consistency", "disabled"),
		chromedp.Flag("disable-features", "PasswordExport,Signin,IdentityManager"),

		chromedp.Flag("unlimited-storage", true),
		chromedp.Flag("allow-unlimited-local-storage", true),
		chromedp.UserDataDir(filepath.Join(cwd, "../chrome-test-profile-doNotTouch/")),
		// for youtube video to work and not get detected/block by automation software

		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "IsolateOrigins,site-per-process"),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-infobars", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"),
		//

		chromedp.Flag("auto-open-devtools-for-tabs", true),
		// This ensures network events are captured
		chromedp.Flag("enable-network-service", true),
	)

	allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)

	// Create context
	ctx, _ := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(func(format string, args ...interface{}) {
			fmt.Printf(format, args...)
		}),
	)

	// Create a timeout for the entire operation
	ctx, cancel := context.WithTimeout(ctx, 40*time.Minute)

	// Navigate to the target URL
	err = chromedp.Run(ctx, network.Enable())
	if err != nil {
		cancel()
		cancel1()
		return nil, nil, nil, err
	}
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
