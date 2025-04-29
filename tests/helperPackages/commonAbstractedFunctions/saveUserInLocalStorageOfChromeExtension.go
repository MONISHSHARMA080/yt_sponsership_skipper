package commonabstractedfunctions

import (
	"context"
	"fmt"
	"youtubeAdsSkipper/paymentBackendGO/common"
	"youtubeAdsSkipper/tests/helperPackages/DB"
	"youtubeAdsSkipper/tests/helperPackages/extension"
	userindb "youtubeAdsSkipper/tests/helperPackages/userInDb"
	userkey "youtubeAdsSkipper/tests/helperPackages/userKey"
)

func SaveUserInLocalStorageOfChromeExtension(ctx context.Context, extensionID string, userKey *userkey.UserKey, userInDb *userindb.Userindb, chromeExtension *extension.ChromeExtension) error {
	DB := DB.DbConnect()
	// userKey  userkey.UserKey{}
	// userInDb userindb.Userindb{}

	getUserIdChann := make(chan common.ErrorAndResultStruct[int64])
	getEncryptedKeyFromUser := make(chan common.ErrorAndResultStruct[string])

	go userInDb.GenerateSpamUserAndSaveItInDB(DB, getUserIdChann)

	userInDBChannResult := <-getUserIdChann
	if userInDBChannResult.Error != nil {
		println("there is a error in getting the user in the DB and it is ->" + userInDBChannResult.Error.Error())
		return userInDBChannResult.Error
	}
	println("the user's id primary key  is ->", userInDBChannResult.Result)

	go userKey.InitializeTheStructAndGetEncryptedKey(userInDb, userInDBChannResult.Result, getEncryptedKeyFromUser)
	encryptedKeyFormChannel := <-getEncryptedKeyFromUser
	if encryptedKeyFormChannel.Error != nil {
		println("there is a error in getting encrypted user key form  and it is ->" + encryptedKeyFormChannel.Error.Error())
		return encryptedKeyFormChannel.Error
	}
	fmt.Printf("the user is %v", *userKey)
	println("browser started")
	// Get the extension ID
	// if err := getExtensionID(ctx, &extensionID); err != nil {
	// 	log.Fatal("Failed to get extension ID:", err)
	// }
	newKeyForNow := encryptedKeyFormChannel.Result

	// Test setting and getting local storage in the service worker
	err := chromeExtension.SetAndtestExtensionStorage(ctx, newKeyForNow)
	if err != nil {
		println("\n\n >>> the test for seeing  if we can write the user key in the strorage failled and the error is  " + err.Error() + "  <<< \n\n ")
		return err
	}
	println("<<--we were able to successfully set the value in the local storage and get the same value back <<--\n\n")
	return nil
}
