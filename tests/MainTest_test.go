package tests

import (
	"fmt"
	"os"
	"testing"
	"youtubeAdsSkipper/tests/helperPackages/DB"
	commonstateacrosstest "youtubeAdsSkipper/tests/helperPackages/commonStateAcrossTest"

	helperfunc1_test "youtubeAdsSkipper/tests/helperFunc1"

	"github.com/joho/godotenv"
)

const (
	extensionID                    = "dpkehfkmkhmbmbofjaikccklcdpdmhpl"
	keyForStorageInChromeExtension = "key"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		panic(err)
	}
	err = DB.CreateDBForTest()
	if err != nil {
		panic(err)
	}
	cancelFunc, err := helperfunc1_test.StartTestServer("../", []string{}, "server.log", "8080")
	if err != nil {
		println("Error starting test server:", err.Error())
		panic(err)
	}
	defer cancelFunc()
	// println("Server started successfully, and sleeping for 30 sec")
	// time.Sleep(30 * time.Second)
	println("sleeping for 30 sec done \n\n")
	ctx, cancelFunc0, cancelFunc, err := helperfunc1_test.GetNewBrowserForChromeExtension(extensionID)
	if err != nil {
		println("there is a error in starting the browser  and the error is ", err.Error())
		panic(err)
	}
	commonstateacrosstest.BrowserContext = ctx
	commonstateacrosstest.CancelFunc = cancelFunc

	defer cancelFunc()
	defer cancelFunc0()

	code := m.Run()

	os.Exit(code)
}

// here load env and set the chrome context in the global state and let other runs
