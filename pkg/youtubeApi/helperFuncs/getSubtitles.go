package helperfuncs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

func GetYoutubeService() error {
	ctx := context.Background()

	clientSecretFile, err := os.Open("client_secret.json")
	if err != nil {
		fmt.Printf("Unable to read client secret file: %s \n", err.Error())
		return err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(clientSecretFile).Decode(t)
	defer clientSecretFile.Close()
	// return t, err

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		fmt.Printf("Unable to parse client secret file to config: %s", err.Error())
		return err
	}
	client := config.Client(ctx)
	// client := getClient(ctx, config)
	service, err := youtube.New(client)
}
