package tests

import (
	"context"
	"fmt"
	"testing"
	"time"
	"youtubeAdsSkipper/paymentBackendGO/common"
	"youtubeAdsSkipper/tests/helperPackages/DB"
	commonstateacrosstest "youtubeAdsSkipper/tests/helperPackages/commonStateAcrossTest"
	"youtubeAdsSkipper/tests/helperPackages/extension"
	userindb "youtubeAdsSkipper/tests/helperPackages/userInDb"
	userkey "youtubeAdsSkipper/tests/helperPackages/userKey"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func TestCheckForChromeExtensionToWebsiteKeyTransfer(t *testing.T) {
	// now I have 2 choices, 1) either I disable the network in the browser, as when It gets the key it will try too
	// see that if it is valid or not of not then get the new one
	// 2) or make the expiry time form the env large so while generating the key we know it willbe sufficiend and will not change

	// 1> get a new user and set the key in the local storage
	println("in the test to check if chromeExtension passes its key to the website")
	ctx := commonstateacrosstest.BrowserContext
	DB := DB.DbConnect()
	chromeExtension := extension.ChromeExtension{ExtensionId: extensionID}
	userKey := userkey.UserKey{}
	userInDb := userindb.Userindb{}
	// make the user
	getUserIdChann := make(chan common.ErrorAndResultStruct[int64])
	go userInDb.GenerateSpamUserAndSaveItInDB(DB, getUserIdChann)
	userInDBChannResult := <-getUserIdChann
	if userInDBChannResult.Error != nil {
		t.Fatal("there is a error in getting the user in the DB and it is ->" + userInDBChannResult.Error.Error())
	}
	println("the user's id primary key  is ->", userInDBChannResult.Result)
	// first change the env key's value to deprecate
	// println("about to change the env key's value during runtime, and it is ->", os.Getenv("TIMEFORFAKEKEYEXPIRY"))
	// // cause we do not want the other ones to be disturbed by our change here and cause unneccary bugs in the tests
	// defer os.Setenv("TIMEFORKEYEXPIRY", os.Getenv("TIMEFORKEYEXPIRY"))
	// err := os.Setenv("TIMEFORKEYEXPIRY", "1ms")
	// if err != nil {
	// 	t.Fatal("there is a error in setting the env key during runtime(which is imp) and it is ->", err.Error())
	// }
	// println("the changed env key at the runtime is ->", os.Getenv("TIMEFORKEYEXPIRY"))
	// value changed value will be 1 ms and a 1 sec pause
	// make the key
	getEncryptedKeyFromUser := make(chan common.ErrorAndResultStruct[string])
	go userKey.InitializeTheStructAndGetEncryptedKey(&userInDb, userInDBChannResult.Result, getEncryptedKeyFromUser)
	encryptedKeyFormChannel := <-getEncryptedKeyFromUser
	if encryptedKeyFormChannel.Error != nil {
		t.Fatal("there is a error in getting encrypted user key form  and it is ->" + encryptedKeyFormChannel.Error.Error())
	}
	userKeyInExtension := encryptedKeyFormChannel.Result
	err := chromeExtension.SetAndtestExtensionStorage(ctx, encryptedKeyFormChannel.Result)
	if err != nil {
		println("\n\n >>> the test for seeing  if we can write the user key in the strorage failled and the error is  " + err.Error() + "  <<< \n\n ")
		t.Fatal(err)
	}
	fmt.Printf("the user is %+v \n", userKey)
	// now go to the website and then see if we are able to
	// err = browserutil.DisableNetwork(ctx)
	// if err != nil {
	// 	println("there is a erorr in stopping the network in the chrome tab")
	// 	t.Fatal(err)
	// }
	err = chromedp.Run(ctx,
		network.Enable(), // enable network tracking :contentReference[oaicite:0]{index=0}
		chromedp.Navigate("http://localhost:5173"), // go to target page
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Put the browser into offline mode
			return network.EmulateNetworkConditions(true, 0, 0, 0).Do(ctx) // go fully offline :contentReference[oaicite:1]{index=1}
		}),
	)
	if err != nil {
		println("there is a error in either going to the website or in the network thingy ->", err.Error())
		t.Fatal(err)
	}
	// err = chromeExtension.NavigateToAWebPage(ctx, "http://localhost:5173")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	time.Sleep(time.Second * 5)
	result, err := chromeExtension.GetValueFormLocalStorageOfWebsite(ctx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf(" the result form getting the website's storage is %+v \n", *result)
	if result.Key == userKeyInExtension {
		println("the user key in the website is same as the one in the chrome extension ")
	} else {
		println("the user key by the chrome extension was not found to be the same as one in the website")
		t.Fatal("the userKey in website was not same as the one in with the chrome extension")
	}
}
