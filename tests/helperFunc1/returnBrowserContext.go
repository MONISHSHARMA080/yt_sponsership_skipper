package helperfunc1_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

func GetNewBrowserForChromeExtension(extensionID string) (context.Context, context.CancelFunc, error) {
	// var opts []chromedp.ExecAllocatorOption
	cwd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}
	extensionPath := filepath.Join(cwd, "../chromeExtension")
	println("the chrome extensionPath is ->", extensionPath)

	// // Common options for both modes
	// baseOpts := []chromedp.ExecAllocatorOption{
	// 	chromedp.NoFirstRun,
	// 	chromedp.NoDefaultBrowserCheck,
	// 	// Disable security features that might block extension access
	// 	chromedp.Flag("disable-web-security", true),
	// 	// Disable content blocking features to prevent ERR_BLOCKED_BY_CLIENT
	// 	chromedp.Flag("disable-client-side-phishing-detection", true),
	// 	chromedp.Flag("disable-features", "BlockInsecurePrivateNetworkRequests,IsolateOrigins,site-per-process,SafeBrowsing"),
	// 	// Allow extension access
	// 	chromedp.Flag("silent-debugger-extension-api", false),
	// 	chromedp.Flag("allow-file-access-from-files", true),
	// 	// Create a user data directory to persist extension state
	// 	chromedp.UserDataDir(filepath.Join(cwd, "chrome-test-profile")),
	// 	chromedp.Flag("headless", false),          // Important: We don't want headless mode
	// 	chromedp.Flag("enable-automation", false), // Disable the automation banner
	// }
	//
	// if useExistingBrowser {
	// 	// Connect to existing Chrome browser with remote debugging
	// 	opts = append(baseOpts,
	// 		chromedp.Flag("remote-debugging-port", "9222"),
	// 		// Important: Make sure extensions are still enabled even in existing browser
	// 		chromedp.Flag("load-extension", extensionPath),
	// 		chromedp.Flag("disable-extensions-except", extensionPath),
	// 		chromedp.Flag("disable-popup-blocking", true),
	// 		chromedp.Flag("disable-sync", true),
	// 		chromedp.Flag("safebrowsing-disable-extension-blacklist", true),
	// 	)
	// } else {
	// 	// Create a new Chrome instance with the extension loaded
	// 	opts = append(baseOpts,
	// 		chromedp.DisableGPU,
	// 		// Load the extension
	// 		chromedp.Flag("load-extension", extensionPath),
	// chromedp.Flag("disable-extensions-except", extensionPath),
	// 		// Disable all browser filtering
	// 		chromedp.Flag("disable-popup-blocking", true),
	// 		chromedp.Flag("disable-sync", true),
	// 		// Disable safe browsing blocking
	// 		chromedp.Flag("safebrowsing-disable-extension-blacklist", true),
	// 	)
	// }
	//
	// // Append default options to our custom ones
	// opts = append(chromedp.DefaultExecAllocatorOptions[:], opts...)
	//
	// // Create a new allocator context with the custom options
	// allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)
	//
	// // Create a browser context with standard logging (not debug)
	// ctx, cancel2 := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	//
	// // Create a context with timeout (2 minutes)
	// ctx, cancel3 := context.WithTimeout(ctx, 120*time.Second)
	//
	// return ctx, cancel1, cancel2, cancel3, nil

	opts := append(chromedp.DefaultExecAllocatorOptions[:],

		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", false),          // Important: We don't want headless mode
		chromedp.Flag("enable-automation", false), // Disable the automation banner

		// chromedp.Flag("extensions-on-chrome-urls", true),
		chromedp.Flag("disable-extensions-except", extensionID),
		chromedp.Flag("load-extension", extensionPath),

		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("load-extension", extensionPath),
		chromedp.Flag("disable-extensions-except", extensionPath),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("safebrowsing-disable-extension-blacklist", true),
		//---
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-signin-promo", true),
		chromedp.Flag("disable-sync-credential-backend", true),
		chromedp.Flag("enable-features", "PasswordImport"),
		chromedp.Flag("disable-features", "PasswordExport,Signin"),

		chromedp.Flag("account-consistency", "disabled"),
		chromedp.Flag("disable-features", "PasswordExport,Signin,IdentityManager"),

		chromedp.Flag("unlimited-storage", true),
		chromedp.Flag("allow-unlimited-local-storage", true),
		// chromedp.UserDataDir(filepath.Join(cwd, "../chrome-test-profile-doNotTouch/")),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	// Create context
	ctx, _ := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(func(format string, args ...interface{}) {
			fmt.Printf(format, args...)
		}),
	)

	// Create a timeout for the entire operation
	ctx, cancel = context.WithTimeout(ctx, 130*time.Second)

	// Navigate to the target URL

	return ctx, cancel, nil
}
