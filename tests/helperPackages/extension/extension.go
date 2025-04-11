package extension

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

type ChromeExtension struct {
	ExtensionId string
}

func (e *ChromeExtension) isTheStructEmpty() bool {
	return e.ExtensionId == ""
}

// this function is there to set the key in the local storage and verify that it is same; will return error as nil if the error is false
func (extension *ChromeExtension) SetAndtestExtensionStorage(ctx context.Context, newValueForKey string) error {
	if extension.isTheStructEmpty() {
		return fmt.Errorf("the extension is un initialized as the id is a empty string ")
	}
	//
	//
	// in this func we will first do is set the value in the storage and see if it is there, and if it is there we return, make it a bool to show that we are ready
	// also make a func that will get the user value either do it in the db (trial, make sqlite here and set the migration etc.) and create a fake one and set it in the db and make a key and give it to this one
	//
	//
	fmt.Println("Testing extension storage...")

	// serviceWorkerURL := fmt.Sprintf("chrome-extension://%s/service-worker.js", extensionID)
	// Navigate to the service worker's devtools page
	// serviceWorkerDebugURL := fmt.Sprintf("chrome://inspect/#extensions")
	err := extension.GoToExtensionPopupPage(ctx)
	if err != nil {
		return err
	}
	valueAtTheKey, err := extension.GetKeysValueFormStorage(ctx)
	if err != nil {
		println("the error in getting the storage in the chrome extension -->", err.Error())
		return err
	} else if valueAtTheKey == "" {
		return fmt.Errorf("the value of the key form the local storage is a empty string ")
	}
	println("the value we got form local storage is ->", valueAtTheKey)
	success, err := extension.SetKeysValueFormStorage(ctx, newValueForKey)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("the success of setting the value is false")
	}
	valueAtTheKeyAfterSettingNewValue, err := extension.GetKeysValueFormStorage(ctx)
	if err != nil {
		println("the error in getting the storage in the chrome extension -->", err.Error())
		return err
	} else if valueAtTheKeyAfterSettingNewValue == "" {
		return fmt.Errorf("the value of the key form the local storage is a empty string ")
	}
	println("the value we got form local storage is ->", valueAtTheKeyAfterSettingNewValue)
	//  didn't do as we might be running it on the same storage and result form previous test can still be stoed here -> valueAtTheKey == valueAtTheKeyAfterSettingNewValue
	if valueAtTheKeyAfterSettingNewValue != newValueForKey {
		return fmt.Errorf("there is a error as either the new value set by us is not eqault to the value returned ")
	}
	return nil
}

func (extension *ChromeExtension) GoToExtensionPopupPage(ctx context.Context) error {
	if extension.isTheStructEmpty() {
		return fmt.Errorf("the extension is un initialized as the id is a empty string ")
	}
	fmt.Println("Opening extension popup...")
	popupURL := fmt.Sprintf("chrome-extension://%s/index.html", extension.ExtensionId)
	println("the poput url is ->", popupURL)
	return chromedp.Run(ctx, chromedp.Navigate(popupURL))
}

// note we are not going to the index.html of the extension here, you have to do that
func (extension *ChromeExtension) GetKeysValueFormStorage(ctx context.Context) (string, error) {
	type StorageResult struct {
		Error string `json:"error"`
		Value string `json:"value"`
	}
	var result StorageResult
	// var rawJson []byte
	script := `
	      new Promise((resolve) => {
	          try {
	              chrome.storage.local.get(['key'], (items) => {
	                  if (chrome.runtime.lastError) {
	                      resolve({
	                          error: chrome.runtime.lastError.message,
	                          value: ""
	                      });
	                  } else {
	                      // If key exists, use its value, otherwise use a hardcoded string
	                      const value = items['key'];
	                      resolve({
	                          error: null,
	                          value: value
	                      });
	                  }
	              });
	          } catch (e) {
	              resolve({
	                  error: e.toString(),
	                  value: ""
	              });
	          }
	      })
	`

	err := chromedp.Run(ctx,
		chromedp.Evaluate(script, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams { return p.WithAwaitPromise(true) }),
	)
	fmt.Printf("\n the script execution is completed and ---the error in the result is %s  -- and the value is %s \n", result.Error, result.Value)
	if err != nil {
		println("the error in getting the value first time is ->", err.Error())
		return "", err
	} else if result.Error != "" {
		fmt.Printf("\n the error in the result is %s  -- and the value is %s \n", result.Error, result.Value)
		return "", fmt.Errorf("the error from the script to get value of the key is ->%s", result.Error)
	}
	return result.Value, nil
}

// note we are not going to the index.html of the extension here, you have to do that
func (extension *ChromeExtension) SetKeysValueFormStorage(ctx context.Context, ValueOfThekeyToSetInStorage string) (bool, error) {
	type StorageResult struct {
		Error   string `json:"error"`
		Success bool   `json:"success"`
	}
	var result StorageResult

	script := `
	      new Promise((resolve) => {
	          try {
	              chrome.storage.local.set({'key': '` + ValueOfThekeyToSetInStorage + `'}, () => {
	                  if (chrome.runtime.lastError) {
	                      resolve({
	                          error: chrome.runtime.lastError.message,
	                          success: false
	                      });
	                  } else {
	                      resolve({
	                          error: null,
	                          success: true
	                      });
	                  }
	              });
	          } catch (e) {
	              resolve({
	                  error: e.toString(),
	                  success: false
	              });
	          }
	      })
	`
	err := chromedp.Run(ctx,
		chromedp.Evaluate(script, &result, func(p *runtime.EvaluateParams) *runtime.EvaluateParams { return p.WithAwaitPromise(true) }),
	)
	fmt.Printf("\n the script execution is completed and ---the error in the result is %s -- and the success is %t \n", result.Error, result.Success)
	if err != nil {
		println("the error in setting the value is ->", err.Error())
		return false, err
	} else if result.Error != "" {
		fmt.Printf("\n the error in the result is %s -- and the success is %t \n", result.Error, result.Success)
		return false, fmt.Errorf("the error from the script to set value of the key is ->%s", result.Error)
	}
	return result.Success, nil
}
