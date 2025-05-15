package tests

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
	commonstructs "youtubeAdsSkipper/commonStructs"
	"youtubeAdsSkipper/paymentBackendGO/common"
	"youtubeAdsSkipper/tests/helperPackages/DB"
	commonstateacrosstest "youtubeAdsSkipper/tests/helperPackages/commonStateAcrossTest"
	"youtubeAdsSkipper/tests/helperPackages/extension"
	userindb "youtubeAdsSkipper/tests/helperPackages/userInDb"
	userkey "youtubeAdsSkipper/tests/helperPackages/userKey"
	// "github.com/chromedp/chromedp"
)

func TestSeeIfWeReplaceDeprecatedKey(t *testing.T) {
	// the goal of this test is to see that are we able to replace the deprecated key in the extension
	// and fetch the new key when the server tells us to do so
	//
	// working: (X)1)create a new user and set the user key in the storage
	// (X) 2) then run the extension on the video(AD) ,
	// (X) 3) now when the video is done we expect the key to be different than the one we set before
	// 4) if we really want to be pedantic, get the new key and decrypt it and also see that the time of expiry of the new key will be

	// power off the phone and then do it , it is a 44-28 min test wth cursor

	commonstateacrosstest.LogChan <- "\n\n\n\n\n\n -----------------" + t.Name() + "------------------\n\n\n\n\n\n"
	ctx := commonstateacrosstest.BrowserContext
	chromeExtension := extension.ChromeExtension{ExtensionId: extensionID}
	print("\n\n in the replcae deprecated key \n\n")
	DB := DB.DbConnect()
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
	// first change the env key's value to deprecate
	println("about to change the env key's value during runtime, and it is ->", os.Getenv("TIMEFORFAKEKEYEXPIRY"))
	// cause we do not want the other ones to be disturbed by our change here and cause unneccary bugs in the tests
	defer os.Setenv("TIMEFORKEYEXPIRY", os.Getenv("TIMEFORKEYEXPIRY"))
	err := os.Setenv("TIMEFORKEYEXPIRY", "1ms")
	if err != nil {
		t.Fatal("there is a error in setting the env key during runtime(which is imp) and it is ->", err.Error())
	}
	println("the changed env key at the runtime is ->", os.Getenv("TIMEFORKEYEXPIRY"))
	// value changed value will be 1 ms and a 1 sec pause
	// make the key
	getEncryptedKeyFromUser := make(chan common.ErrorAndResultStruct[string])
	go userKey.InitializeTheStructAndGetEncryptedKey(&userInDb, userInDBChannResult.Result.UserID, getEncryptedKeyFromUser)
	encryptedKeyFormChannel := <-getEncryptedKeyFromUser
	if encryptedKeyFormChannel.Error != nil {
		t.Fatal("there is a error in getting encrypted user key form  and it is ->" + encryptedKeyFormChannel.Error.Error())
	}
	oldUserKey := userKey.User.GetEncryptedKey()
	err = chromeExtension.SetAndtestExtensionStorage(ctx, encryptedKeyFormChannel.Result)
	if err != nil {
		println("\n\n >>> the test for seeing  if we can write the user key in the strorage failled and the error is  " + err.Error() + "  <<< \n\n ")
		t.Fatal(err)
	}
	fmt.Printf("the user is %+v \n", userKey)
	// now lets play the video, and for every few second check it to see if we are still working
	youtubeUrls := []string{"https://www.youtube.com/watch?v=korOpibkm6g", "https://www.youtube.com/watch?v=D3cjV3tNd88", "https://www.youtube.com/watch?v=NOfUCMzBNVg", "https://www.youtube.com/watch?v=WVn4FPULFWA"}
	youTubeVideo := youtubeUrls[rand.Intn(len(youtubeUrls))]
	println("the youtube video selected at random is ->", youTubeVideo)
	err = chromeExtension.NavigateToAWebPage(ctx, youTubeVideo)
	if err != nil {
		println("there is a error in going to the youTubeVideo's url ")
		t.Fatal(err)
	}
	stopChannelToStopChekingIfTheVideoIsPlaying := make(chan struct{})
	defer close(stopChannelToStopChekingIfTheVideoIsPlaying) // this is a send only channel so only we can close it
	go chromeExtension.EnsureVideoIsPlayingPeriodically(ctx, time.Millisecond*700, stopChannelToStopChekingIfTheVideoIsPlaying, false)
	// now probally watch for the network req etc or just sleep for 5 sec and check again
	time.Sleep(time.Second * 5)
	newKeyFromStorage, err := chromeExtension.GetKeysValueFormStorageByGoingToTheIndexPage(ctx)
	if err != nil {
		t.Fatal(err)
	}
	println("the new key is ->", newKeyFromStorage, " \n\n ==== and the old key is -->", oldUserKey, " \n\n are they not equal (ans shoukd be yes here)->", newKeyFromStorage != oldUserKey)
	if newKeyFromStorage == oldUserKey {
		t.Fatal("the old userKey is same as the new one, meaning either we did not change or we might need more time")
	}
}
