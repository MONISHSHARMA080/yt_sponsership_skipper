package tests

import (
	"fmt"
	"log"
	"testing"
	"youtubeAdsSkipper/paymentBackendGO/common"
	"youtubeAdsSkipper/tests/helperPackages/DB"
	commonstateacrosstest "youtubeAdsSkipper/tests/helperPackages/commonStateAcrossTest"
	"youtubeAdsSkipper/tests/helperPackages/extension"
	userindb "youtubeAdsSkipper/tests/helperPackages/userInDb"
	userkey "youtubeAdsSkipper/tests/helperPackages/userKey"

	"github.com/chromedp/chromedp"
)

func TestSetKeyInTheLocalStorage(t *testing.T) {
	// get the key
	commonstateacrosstest.LogChan <- "\n\n\n\n\n\n -----------------" + t.Name() + "------------------\n\n\n\n\n\n"

	ctx := commonstateacrosstest.BrowserContext
	// Start the browser
	if err := chromedp.Run(ctx); err != nil {
		log.Fatal("Failed to start browser:", err)
	}
	DB := DB.DbConnect()
	userKey := userkey.UserKey{}
	userInDb := userindb.Userindb{}

	getUserIdChann := make(chan common.ErrorAndResultStruct[int64])
	getEncryptedKeyFromUser := make(chan common.ErrorAndResultStruct[string])

	go userInDb.GenerateSpamUserAndSaveItInDB(DB, getUserIdChann)

	userInDBChannResult := <-getUserIdChann
	if userInDBChannResult.Error != nil {
		t.Fatal("there is a error in getting the user in the DB and it is ->" + userInDBChannResult.Error.Error())
	}
	println("the user's id primary key  is ->", userInDBChannResult.Result)

	go userKey.InitializeTheStructAndGetEncryptedKey(&userInDb, userInDBChannResult.Result, getEncryptedKeyFromUser)
	encryptedKeyFormChannel := <-getEncryptedKeyFromUser
	if encryptedKeyFormChannel.Error != nil {
		t.Fatal("there is a error in getting encrypted user key form  and it is ->" + encryptedKeyFormChannel.Error.Error())
	}
	fmt.Printf("the user is %v", userKey)
	println("browser started")
	// Get the extension ID
	// if err := getExtensionID(ctx, &extensionID); err != nil {
	// 	log.Fatal("Failed to get extension ID:", err)
	// }
	fmt.Printf("Extension ID: %s\n", extensionID)
	newKeyForNow := encryptedKeyFormChannel.Result
	println("the new key that we are setting in the local storage is ->", newKeyForNow, " --")

	chromeExtension := extension.ChromeExtension{ExtensionId: extensionID}

	// Test setting and getting local storage in the service worker
	err := chromeExtension.SetAndtestExtensionStorage(ctx, newKeyForNow)
	if err != nil {
		println("\n\n >>> the test for seeing  if we can write the user key in the strorage failled and the error is  " + err.Error() + "  <<< \n\n ")
		t.Fatal(err)
	}
	println("<<--we were able to successfully set the value in the local storage and get the same value back <<--\n\n")
}
