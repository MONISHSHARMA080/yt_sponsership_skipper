package tests

import (
	"fmt"
	"log"
	"testing"
	"time"
	"youtubeAdsSkipper/paymentBackendGO/common"
	"youtubeAdsSkipper/tests/helperPackages/DB"
	"youtubeAdsSkipper/tests/helperPackages/extension"
	userindb "youtubeAdsSkipper/tests/helperPackages/userInDb"
	userkey "youtubeAdsSkipper/tests/helperPackages/userKey"

	helperfunc1_test "youtubeAdsSkipper/tests/helperFunc1"

	"github.com/chromedp/chromedp"
)

const (
	extensionID                    = "dpkehfkmkhmbmbofjaikccklcdpdmhpl"
	keyForStorageInChromeExtension = "key"
)

func TestMain(t *testing.T) {
	// get the key

	ctx, cancelFunc, err := helperfunc1_test.GetNewBrowserForChromeExtension(extensionID)
	if err != nil {
		t.Fatal(err)
	}
	defer cancelFunc()
	// Start the browser
	if err := chromedp.Run(ctx); err != nil {
		log.Fatal("Failed to start browser:", err)
	}
	DB := DB.DbConnect()

	userInDb := userindb.Userindb{}
	getUserIdChann := make(chan common.ErrorAndResultStruct[int64])
	go userInDb.GenerateSpamUserAndSaveItInDB(DB, getUserIdChann)
	userInDBChannResult := <-getUserIdChann
	if userInDBChannResult.Error != nil {
		t.Fatal("there is a error in getting the user in the DB and it is ->" + userInDBChannResult.Error.Error())
	}
	println("the user's id is ->", userInDBChannResult.Result)

	userKey := userkey.UserKey{}
	println("browser started")
	// Get the extension ID
	// if err := getExtensionID(ctx, &extensionID); err != nil {
	// 	log.Fatal("Failed to get extension ID:", err)
	// }
	fmt.Printf("Extension ID: %s\n", extensionID)
	newKeyForNow := "IamAGod100"

	chromeExtension := extension.ChromeExtension{ExtensionId: extensionID}

	// Test setting and getting local storage in the service worker
	err = chromeExtension.SetAndtestExtensionStorage(ctx, newKeyForNow)
	if err != nil {
		t.Fatal(err)
	}
	println("we were able to successfully set the value in the local storage and get the same value back")
	time.Sleep(time.Minute * 2)
}
