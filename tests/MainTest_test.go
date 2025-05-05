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
	println("the api key form the env is ->", os.Getenv("API_KEYS"))
	envFileContent, err := os.ReadFile("../.env")
	if err != nil {
		println("Error opening .env file:", err.Error())
		panic(err)
	}
	fmt.Printf("\n\n --Content of .env file:\n %s \n\n\n\n\n", string(envFileContent))
	fmt.Printf("\n\n --Content of .env file:\n %s \n\n\n\n\n", string(envFileContent[:38]))

	err = DB.CreateDBForTest()
	if err != nil {
		panic(err)
	}
	commonstateacrosstest.LogChan = make(chan string, 100) // Increased buffer size
	cancelFunc, err := helperfunc1_test.StartTestServer("../", []string{}, "server.log", "8080", commonstateacrosstest.LogChan)
	if err != nil {
		println("Error starting test server:", err.Error())
		panic(err)
	}
	defer cancelFunc()
	commonstateacrosstest.LogChan <- "Test server started, proceeding to browser setup."
	defer close(commonstateacrosstest.LogChan)
	println("Test server started, proceeding to browser setup.")
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
