package helperfuncs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// FetchVideoCaptions fetches captions for the given videoID,
// automatically refreshing the token if it’s expired.
func FetchVideoCaptions(ctx context.Context, config *oauth2.Config, tok *oauth2.Token, tokenFile, videoID string) (string, *oauth2.Token, error) {
	// 0. Ensure our token is valid; refresh & save if expired.
	newTok, err := ensureTokenValid(ctx, config, tok, tokenFile)
	if err != nil {
		return "", tok, fmt.Errorf("token refresh failed: %v", err)
	}
	// Use the (possibly refreshed) token from now on
	tok = newTok

	// 1. Build an authenticated HTTP client
	httpClient := config.Client(ctx, tok)

	// 2. Create the YouTube service
	service, err := youtube.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return "", tok, fmt.Errorf("unable to create YouTube service: %v", err)
	}

	// 3. List caption tracks for the video
	listResp, err := service.Captions.
		List([]string{"id", "snippet"}, videoID).
		Do()
	println("did the youtube captions call")
	if err != nil {
		return "", tok, fmt.Errorf("captions.list API error: %v", err)
	}
	fmt.Printf("\n the captions list is ->%+v   \n\n", *listResp)
	if len(listResp.Items) == 0 {
		return "", tok, fmt.Errorf("no captions found for video(in the api call) %q", videoID)
	}

	// 4. Pick the first caption track (or filter by language)
	captionID := listResp.Items[0].Id

	// 5. Download that caption track in SRT format
	downloadCall := service.Captions.Download(captionID)
	resp, err := downloadCall.Download()
	if err != nil {
		println("there is a error in donwloading the captions")
		return "", tok, fmt.Errorf("captions.download API error: %v", err)
	}
	defer resp.Body.Close()
	fmt.Printf("the captions form the download call is (request, status code %d and status is %s bodyis )-> %+v \n\n", resp.StatusCode, resp.Status, resp.Body)

	// 6. Read the response body
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return "", tok, fmt.Errorf("reading caption response body: %v", err)
	}

	return buf.String(), tok, nil
}

// ensureTokenValid takes your existing token (loaded from file or just exchanged),
// and uses oauth2.ReuseTokenSource to automatically refresh it if needed.
// It returns a valid token, and will save an updated token.json if a refresh occurred.
func ensureTokenValid(ctx context.Context, config *oauth2.Config, tok *oauth2.Token, tokenFile string) (*oauth2.Token, error) {
	src := config.TokenSource(ctx, tok)
	newTok, err := src.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token from source: %v", err)
	}
	if newTok.AccessToken != tok.AccessToken {
		log.Println("Access token refreshed, saving new token to disk")
		saveToken(tokenFile, newTok)
	}
	return newTok, nil
}

func saveToken(path string, token *oauth2.Token) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to create token file %q: %v", path, err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(token); err != nil {
		log.Fatalf("Unable to encode token to %q: %v", path, err)
	}
	fmt.Printf("Token saved to %s\n", path)
}

// tokenFromFile retrieves a token from a file, or returns an error.
func TokenFromFile(path string) (*oauth2.Token, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read token file: %w", err)
	}
	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("unmarshal token JSON: %w", err)
	}
	return &token, nil
}
