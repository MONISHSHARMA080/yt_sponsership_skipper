package tests

import (
	"os"
	"testing"
	helperfunc1_test "youtubeAdsSkipper/tests/helperFunc1"
	commonabstractedfunctions "youtubeAdsSkipper/tests/helperPackages/commonAbstractedFunctions"
	commonstateacrosstest "youtubeAdsSkipper/tests/helperPackages/commonStateAcrossTest"
	"youtubeAdsSkipper/tests/helperPackages/extension"
	userindb "youtubeAdsSkipper/tests/helperPackages/userInDb"
	userkey "youtubeAdsSkipper/tests/helperPackages/userKey"
)

func TestAAAAPaymentToUpgradeTheUserPlan(t *testing.T) {
	// go to the website and  then see if we are able to make a payment and  is it working
	//
	//also run the localtunnel and create the razorpay webhook form the lib, and delete when the test is over
	// we are not creating a new webhook as we don't need it, we will just create the tunnel at the already existing one
	ctx := commonstateacrosstest.BrowserContext

	tunnelCMD, err := helperfunc1_test.StartLocalTunnel()
	if err != nil {
		t.Fatal("there is an error in starting the local tunnel and it is ->", err.Error())
		// println("there is an error in starting the local tunnel and it is ->", err.Error())
	}

	defer helperfunc1_test.StopLocalTunnel(tunnelCMD)
	println("the env for the sqlite file is ->", os.Getenv("TURSO_DATABASE_URL"), " and we in the testing environment ->", os.Getenv("IS_THIS_TESTING_ENVIRONMENT"))
	userKey := userkey.UserKey{}
	userInDb := userindb.Userindb{}
	chromeExtension := extension.ChromeExtension{ExtensionId: extensionID}
	err = commonabstractedfunctions.SaveUserInLocalStorageOfChromeExtension(ctx, extensionID, &userKey, &userInDb, &chromeExtension)
	if err != nil {
		t.Fatal("there is an error in saving the user in the local storage of the chrome extension and it is ->", err.Error())
	}
	err = chromeExtension.MakeThePaymentAndGetOnPaidTier(ctx, false)
	if err != nil {
		println("there is a error in making the payment and it is ->", err.Error())
		t.Fatal(err)
	}
	println("we have clicked the payment button and now will do the one for razorpay")
}
