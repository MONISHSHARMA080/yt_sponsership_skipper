package youtubeapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	// change this to the same URI you configured in Google Cloud Console:
	oauth2RedirectURL = "http://localhost:8080/oauth2callback"
	config            *oauth2.Config
)

// loadOAuthConfig reads client_secret.json and builds an oauth2.Config.
func loadOAuthConfig() *oauth2.Config {
	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Error reading client_secret.json: %v", err)
	}
	cfg, err := google.ConfigFromJSON(b,
		youtube.YoutubeReadonlyScope,
		youtube.YoutubeUploadScope,
		// add extra scopes here if needed
	)
	if err != nil {
		log.Fatalf("Error parsing client_secret.json: %v", err)
	}
	cfg.RedirectURL = oauth2RedirectURL
	return cfg
}

func main() {
	config = loadOAuthConfig()

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/oauth2callback", handleCallback)

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// handleLogin redirects user to Google’s OAuth consent page.
func handleLogin(w http.ResponseWriter, r *http.Request) {
	// "state" can be used to verify callback integrity.
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

// handleCallback handles the OAuth callback, exchanges code for token,
// and creates a YouTube service client.
func handleCallback(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse query", http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "Code not found in query", http.StatusBadRequest)
		return
	}

	// Exchange the code for a token
	ctx := context.Background()
	tok, err := config.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Persist tok (tok.AccessToken, tok.RefreshToken) as you see fit
	// e.g. saveToken("token.json", tok)

	// Create YouTube client
	client := config.Client(ctx, tok)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		http.Error(w, "YouTube client creation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Example API call: fetch your channel info
	resp, err := service.Channels.
		List([]string{"snippet", "statistics"}).
		Mine(true).
		Do()
	if err != nil {
		http.Error(w, "API call error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Hello, %s! You have %d subscribers.",
		resp.Items[0].Snippet.Title,
		resp.Items[0].Statistics.SubscriberCount)
}
