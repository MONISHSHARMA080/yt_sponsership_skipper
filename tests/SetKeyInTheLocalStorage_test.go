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
	"github.com/joho/godotenv"
)

const (
	extensionID                    = "dpkehfkmkhmbmbofjaikccklcdpdmhpl"
	keyForStorageInChromeExtension = "key"
)

func TestMain(t *testing.T) {
	// get the key

	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		t.Fatal(err)
	}
	err = DB.CreateDBForTest()
	if err != nil {
		t.Fatal(err)
	}
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
	err = chromeExtension.SetAndtestExtensionStorage(ctx, newKeyForNow)
	if err != nil {
		t.Fatal(err)
	}
	println("we were able to successfully set the value in the local storage and get the same value back")
	time.Sleep(time.Hour * 2)
}
