package youtubeapi

import (
	// Keep context if you plan to use it within the function for other calls, though not strictly needed for this specific function body if service is pre-configured.

	"context"
	"fmt"
	helperfuncs "youtubeAdsSkipper/pkg/youtubeApi/helperFuncs"

	"golang.org/x/oauth2"
)

const (
	clientSecretFile = "client_secret.json" // Make sure this is in your project root
	tokenFile        = "token.json"         // This file will store your access and refresh tokens
	// Scopes required:
	// youtube.YoutubeReadonlyScope for reading data like caption lists.
	// youtube.YoutubeForceSslScope is required by the API for caption downloads.
	youtubeReadonlyScope = "https://www.googleapis.com/auth/youtube.readonly"
	youtubeForceSslScope = "https://www.googleapis.com/auth/youtube.force-ssl"
)

// GetSubtitlesForVideo fetches subtitles for a given YouTube video ID.
// It prioritizes manually created English subtitles, then auto-generated English subtitles.
// The 'service' parameter is an initialized *youtube.Service.
// It returns the caption text in SRT format.
func GetSubtitlesForVideo(config *oauth2.Config, videoID string) (string, error) {
	tokenFile := "token.json"
	token, err := helperfuncs.TokenFromFile(tokenFile)
	if err != nil {
		return "", err
	}
	transcript, a, err := helperfuncs.FetchVideoCaptions(context.Background(), config, token, tokenFile, videoID)
	if err != nil {
		println("err:======8===DD", err.Error())
		return "", err
	}
	fmt.Printf("the oauth token is %+v \n", *a)
	return transcript, nil
}
