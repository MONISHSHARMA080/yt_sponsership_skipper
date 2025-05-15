package tests

import (
	"fmt"
	"testing"
	"time"
	"youtubeAdsSkipper/paymentBackendGO/common"
	"youtubeAdsSkipper/tests/helperPackages/DB"
	browserutil "youtubeAdsSkipper/tests/helperPackages/browserUtil"
	commonstateacrosstest "youtubeAdsSkipper/tests/helperPackages/commonStateAcrossTest"
	"youtubeAdsSkipper/tests/helperPackages/extension"
	userindb "youtubeAdsSkipper/tests/helperPackages/userInDb"
	userkey "youtubeAdsSkipper/tests/helperPackages/userKey"
	commonstructs "youtubeAdsSkipper/commonStructs"
)

func TestCheckForChromeExtensionToWebsiteKeyTransfer(t *testing.T) {
	// now I have 2 choices, 1) either I disable the network in the browser, as when It gets the key it will try too
	// see that if it is valid or not of not then get the new one
	// 2) or make the expiry time form the env large so while generating the key we know it willbe sufficiend and will not change
	// this is the already happening as we have not changed the key
	// helperfunc1_test.LogTestNameInTheServerLogFile(t)
	commonstateacrosstest.LogChan <- "\n\n\n\n\n\n -----------------" + t.Name() + "------------------\n\n\n\n\n\n"
	println("\n\n\n\n\n in the test to check if chromeExtension passes its key to the website")
	defer println("\n\n\n\n\n\n")
	ctx := commonstateacrosstest.BrowserContext
	DB := DB.DbConnect()
	chromeExtension := extension.ChromeExtension{ExtensionId: extensionID}
	userKey := userkey.UserKey{}
	userInDb := userindb.Userindb{}
	// make the user
	getUserIdChann := make(chan common.ErrorAndResultStruct[commonstructs.SignupResult])
	go userInDb.GenerateSpamUserAndSaveItInDB(DB, getUserIdChann)
	userInDBChannResult := <-getUserIdChann
	if userInDBChannResult.Error != nil {
		t.Fatal("there is a error in getting the user in the DB and it is ->" + userInDBChannResult.Error.Error())
	}
	println("the user's id primary key  is ->", userInDBChannResult.Result)
	getEncryptedKeyFromUser := make(chan common.ErrorAndResultStruct[string])
	go userKey.InitializeTheStructAndGetEncryptedKey(&userInDb, userInDBChannResult.Result.UserID, getEncryptedKeyFromUser)
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
	err = chromeExtension.NavigateToAWebPage(ctx, "http://localhost:5173")
	if err != nil {
		t.Fatal(err)
	}
	defer browserutil.EnableNetwork(ctx)
	time.Sleep(time.Second * 10)
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
