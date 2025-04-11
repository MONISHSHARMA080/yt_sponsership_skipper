package tests

import (
	"fmt"
	"log"
	"testing"
	"time"
	"youtubeAdsSkipper/tests/helperPackages/extension"

	helperfunc1_test "youtubeAdsSkipper/tests/helperFunc1"

	"github.com/chromedp/chromedp"
)

const (
	extensionID                    = "dpkehfkmkhmbmbofjaikccklcdpdmhpl"
	keyForStorageInChromeExtension = "key"
)

func TestTesticale(t *testing.T) {
	ctx, cancelFunc, err := helperfunc1_test.GetNewBrowserForChromeExtension(extensionID)
	if err != nil {
		t.Fatal(err)
	}
	defer cancelFunc()
	// Start the browser
	if err := chromedp.Run(ctx); err != nil {
		log.Fatal("Failed to start browser:", err)
	}

	println("browser started")
	// Get the extension ID
	// if err := getExtensionID(ctx, &extensionID); err != nil {
	// 	log.Fatal("Failed to get extension ID:", err)
	// }
	fmt.Printf("Extension ID: %s\n", extensionID)
	newKeyForNow := "IamAGod"

	chromeExtension := extension.ChromeExtension{ExtensionId: extensionID}

	// Test setting and getting local storage in the service worker
	err = chromeExtension.SetAndtestExtensionStorage(ctx, newKeyForNow)
	if err != nil {
		t.Fatal(err)
	}
	println("we were able to successfully set the value in the local storage and get the same value back")

	time.Sleep(4 * time.Minute)
}
