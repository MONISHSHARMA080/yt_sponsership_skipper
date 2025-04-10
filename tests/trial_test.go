package tests

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	helperfunc1_test "youtubeAdsSkipper/tests/helperFunc1"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/cdproto/target"
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

	// Test setting and getting local storage in the service worker
	testExtensionStorage(ctx, extensionID)
	time.Sleep(120 * time.Second)

	// Navigate to YouTube to test content script
	if err := navigateToYouTube(ctx); err != nil {
		log.Fatal("Failed to navigate to YouTube:", err)
	}

	// Open extension popup
	if err := openExtensionPopup(ctx, extensionID); err != nil {
		log.Fatal("Failed to open extension popup:", err)
	}

	// Wait for user to see the results
	time.Sleep(1 * time.Second)
}

func testExtensionStorage(ctx context.Context, extensionID string) error {
	fmt.Println("Testing extension storage...")

	// Navigate to the service worker's devtools page
	serviceWorkerURL := fmt.Sprintf("chrome-extension://%s/service-worker.js", extensionID)
	serviceWorkerDebugURL := fmt.Sprintf("chrome://inspect/#extensions")

	if err := chromedp.Run(ctx, chromedp.Navigate(serviceWorkerDebugURL)); err != nil {
		return fmt.Errorf("failed to navigate to service worker debug page: %w", err)
	}

	// Wait a bit for the page to load
	time.Sleep(5 * time.Second)

	// Now directly access the storage via CDP protocol
	// First, get the target for our service worker
	targets, err := chromedp.Targets(ctx)
	if err != nil {
		return fmt.Errorf("failed to get targets: %w", err)
	}
	indexHTMLURL := fmt.Sprintf("chrome-extension://%s/index.html", extensionID)
	err = chromedp.Run(ctx, chromedp.Navigate(indexHTMLURL))
	if err != nil {
		println("can't navigate to the index html in the testing extension storage url and the error is ->", err.Error(), "\n\n")
		if err := chromedp.Run(ctx, chromedp.Navigate(serviceWorkerDebugURL)); err != nil {
			return fmt.Errorf("failed to navigate to service worker debug page: %w", err)
		}
	}
	println("the target length is ", len(targets), " and serviceWorkerURL is ", serviceWorkerURL)
	//
	//
	//
	// get the storage and set it
	//
	//
	err = accessStorageViaExtensionContext(ctx)
	if err != nil {
		println("the error in getting the storage in the chrome extension -->", err.Error())
		return err
	}
	return nil
}

func navigateToYouTube(ctx context.Context) error {
	fmt.Println("Navigating to YouTube...")
	return chromedp.Run(ctx, chromedp.Navigate("https://www.youtube.com"))
}

func openExtensionPopup(ctx context.Context, extensionID string) error {
	fmt.Println("Opening extension popup...")
	popupURL := fmt.Sprintf("chrome-extension://%s/index.html", extensionID)
	println("the poput url is ->", popupURL)
	return chromedp.Run(ctx, chromedp.Navigate(popupURL))
}

// this function will navigate to the
// Find the extension's background page target
func accessStorageViaExtensionContext(ctx context.Context) error {
	targets, err := chromedp.Targets(ctx)
	if err != nil {
		return err
	}

	var extensionTarget *target.Info
	for i, t := range targets {
		fmt.Printf("the targets at the index %d  is -> %v", i, t)
		if t.Type == "service_worker" && t.URL == fmt.Sprintf("chrome-extension://%s/service-worker.js", extensionID) {
			extensionTarget = t
			break
		}
	}

	if extensionTarget == nil {
		return fmt.Errorf("extension service worker not found")
	}

	// Create a new context for the extension's service worker
	extCtx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(extensionTarget.TargetID))
	defer cancel()
	errorFormtrialJS := GetKeysValueFormStorage(ctx)
	if errorFormtrialJS != nil {
		println("the error seeing if my promise returs is true and it is ->", errorFormtrialJS.Error())
	}
	type StorageResult struct {
		Error string `json:"error"`
		Value string `json:"value"`
	}
	var result StorageResult

	script := `
    new Promise((resolve) => {
        chrome.storage.local.get(["` + keyForStorageInChromeExtension + `"], (data)=> {
            resolve(data["` + keyForStorageInChromeExtension + `"]);
        });
    });
    `
	println("\n\n the script to get key is ->", script, "\n\n")
	valueOfKeyInStorage := ""
	var rawJson []byte
	err = chromedp.Run(extCtx, chromedp.Evaluate(`
        new Promise((resolve) => {
            try {
                chrome.storage.local.get(['key'], (items) => {
                    if (chrome.runtime.lastError) {
                        resolve({
                            error: chrome.runtime.lastError.message,
                            value: "ERROR_OCCURRED"
                        });
                    } else {
                        // If key exists, use its value, otherwise use a hardcoded string
                        const value = items['key'] || "NO_KEY_FOUND_HARDCODED_STRING";
                        resolve({
                            error: null,
                            value: value
                        });
                    }
                });
            } catch (e) {
                resolve({
                    error: e.toString(),
                    value: "EXCEPTION_OCCURRED"
                });
            }
        })
    `, &rawJson, func(p *runtime.EvaluateParams) *runtime.EvaluateParams { return p.WithAwaitPromise(true) }),
	)
	fmt.Printf("\n the script execution is completed and ---the error in the result is %s  -- and the value is %s and the value form the rawJSON is ->%s \n", result.Error, result.Value, rawJson)
	if err != nil {
		println("the error in getting the value first time is ->", err.Error())
		return err
	} else if result.Error != "" {
		fmt.Printf("\n the error in the result is %s  -- and the value is %s \n", result.Error, result.Value)
		// return fmt.Errorf("the error from the script to get value of the key is ->%s", result.Error)
	}
	println("the result from the getting the value form the storage is ->", result.Value)

	// time.Sleep(3 * time.Second)
	// Now you can run commands in the extension's context
	if err := chromedp.Run(extCtx, chromedp.ActionFunc(func(ctx context.Context) error {
		script := `chrome.storage.local.set({key: "--------this is a new value to see id we are takign the value form the sctript the while time"});`
		_, _, err := runtime.Evaluate(script).Do(ctx)
		return err
	})); err != nil {
		return err
	}
	println("the value of the key is ->", valueOfKeyInStorage)
	err = chromedp.Run(extCtx, chromedp.Evaluate(`
        new Promise((resolve) => {
            try {
                chrome.storage.local.get(['key'], (items) => {
                    if (chrome.runtime.lastError) {
                        resolve({
                            error: chrome.runtime.lastError.message,
                            value: "ERROR_OCCURRED"
                        });
                    } else {
                        // If key exists, use its value, otherwise use a hardcoded string
                        const value = items['key'] || "NO_KEY_FOUND_HARDCODED_STRING";
                        resolve({
                            error: null,
                            value: value
                        });
                    }
                });
            } catch (e) {
                resolve({
                    error: e.toString(),
                    value: "EXCEPTION_OCCURRED"
                });
            }

        })
    `, &rawJson))
	fmt.Printf("\n the script execution is completed and ---the error in the result is %s  -- and the value is %s and the value form the rawJSON is ->%s \n", result.Error, result.Value, string(rawJson))
	if err != nil {
		println("the error in getting the value first time is ->", err.Error())
		return err
	} else if result.Error != "" {
		fmt.Printf("\n the error in the result is %s  -- and the value is %s \n", result.Error, result.Value)
		return fmt.Errorf("the error from the script to get value of the key is ->%s", result.Error)
	}
	println("the result from the getting the value form the storage is ->", result.Value)
	return nil
}

func GetKeysValueFormStorage(ctx context.Context) error {
	type StorageResult struct {
		Error string `json:"error"`
		Value string `json:"value"`
	}
	var result StorageResult
	var rawJson []byte
	script := `
	      new Promise((resolve) => {
	          try {
	              chrome.storage.local.get(['key'], (items) => {
	                  if (chrome.runtime.lastError) {
	                      resolve({
	                          error: chrome.runtime.lastError.message,
	                          value: "ERROR_OCCURRED"
	                      });
	                  } else {
	                      // If key exists, use its value, otherwise use a hardcoded string
	                      const value = items['key']||"----"+JSON.stringify(items)+"NO_ORIGINAL_VALUE_FOUND" ;
	                      resolve({
	                          error: null,
	                          value: value
	                      });
	                  }
	              });
	          } catch (e) {
	              resolve({
	                  error: e.toString(),
	                  value: "EXCEPTION_OCCURRED"
	              });
	          }
	      })

	`
	println("\n +the raw json is before running the script is ->", string(rawJson))

	// ------- new script for trial
	// script := `
	//    new Promise(resolve => {
	//        // Your async operation here
	//  setTimeout(() => resolve({value: "Promise result", error:"there is no error "}), 1000);
	//    })
	//  `
	err := chromedp.Run(ctx,
		chromedp.Evaluate(script, &rawJson, func(p *runtime.EvaluateParams) *runtime.EvaluateParams { return p.WithAwaitPromise(true) }),
	)
	fmt.Printf("\n the script execution is completed and ---the error in the result is %s  -- and the value is %s and the value form the rawJSON is ->%s \n", result.Error, result.Value, string(rawJson))
	if err != nil {
		println("the error in getting the value first time is ->", err.Error())
		return err
	} else if result.Error != "" {
		fmt.Printf("\n the error in the result is %s  -- and the value is %s \n", result.Error, result.Value)
		// return fmt.Errorf("the error from the script to get value of the key is ->%s", result.Error)
	}
	return nil
}
