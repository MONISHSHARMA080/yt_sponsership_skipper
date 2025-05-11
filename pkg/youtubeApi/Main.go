package youtubeapi

import (
	// Keep context if you plan to use it within the function for other calls, though not strictly needed for this specific function body if service is pre-configured.
	"fmt"
	"io"
	"net/http"
	"strings"

	"google.golang.org/api/youtube/v3"
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
func GetSubtitlesForVideo(service *youtube.Service, videoID string) (string, error) {
	if videoID == "" {
		return "", fmt.Errorf("videoID cannot be empty")
	}
	if service == nil {
		return "", fmt.Errorf("youtube service cannot be nil")
	}

	// 1. List available caption tracks for the video.
	// We request the "snippet" part to get language and track kind information.
	// The videoID is the ID of the YouTube video (e.g., "dQw4w9WgXcQ").
	listCall := service.Captions.List([]string{"snippet"}, videoID)
	response, err := listCall.Do()
	if err != nil {
		// This can happen if the video doesn't exist, is private, or has captions disabled by the owner.
		return "", fmt.Errorf("failed to list captions for video %s: %w", videoID, err)
	}

	if len(response.Items) == 0 {
		return "", fmt.Errorf("no caption tracks found for video %s", videoID)
	}

	var manualEnCaptionID string
	var autoEnCaptionID string

	// 2. Find the desired English caption track.
	// We prioritize manual English captions ("standard") over auto-generated ones ("ASR").
	for _, item := range response.Items {
		if item.Snippet == nil {
			continue
		}
		// Normalize language code for comparison (e.g., "en-US" becomes "en")
		language := strings.ToLower(item.Snippet.Language)
		trackKind := item.Snippet.TrackKind // "standard" or "ASR"

		// Check for English language (e.g., "en", "en-us", "en-gb")
		if strings.HasPrefix(language, "en") {
			if trackKind == "standard" { // Manually created captions
				manualEnCaptionID = item.Id
				// Found the best option (manual English), no need to search further.
				break
			} else if trackKind == "ASR" { // Auto-generated captions
				// Store the first auto-generated English caption found, in case no manual one exists.
				if autoEnCaptionID == "" {
					autoEnCaptionID = item.Id
				}
			}
		}
	}

	var chosenCaptionID string
	if manualEnCaptionID != "" {
		chosenCaptionID = manualEnCaptionID
		// For debugging or logging:
		// fmt.Printf("Found manual English caption track ID: %s for video %s\n", chosenCaptionID, videoID)
	} else if autoEnCaptionID != "" {
		chosenCaptionID = autoEnCaptionID
		// For debugging or logging:
		// fmt.Printf("Found auto-generated English caption track ID: %s for video %s\n", chosenCaptionID, videoID)
	} else {
		return "", fmt.Errorf("no English captions (manual or auto-generated) found for video %s", videoID)
	}

	// 3. Download the chosen caption track.
	// We can specify the format using Tfmt (e.g., "srt", "vtt"). Default is "srt".
	downloadCall := service.Captions.Download(chosenCaptionID)
	downloadCall.Tfmt("srt") // Requesting SRT format

	// The Download() method on a Captions.Download call returns an *http.Response and an error.
	captionHTTPResp, err := downloadCall.Download()
	if err != nil {
		return "", fmt.Errorf("failed to initiate download for caption track %s: %w", chosenCaptionID, err)
	}
	defer captionHTTPResp.Body.Close()

	if captionHTTPResp.StatusCode != http.StatusOK {
		// Try to read the body for more error details from YouTube API
		bodyBytes, readErr := io.ReadAll(captionHTTPResp.Body)
		errorDetail := ""
		if readErr == nil {
			errorDetail = string(bodyBytes)
		}
		return "", fmt.Errorf("failed to download caption track %s, status: %s, details: %s",
			chosenCaptionID, captionHTTPResp.Status, errorDetail)
	}

	// Read the content of the caption file.
	captionBytes, err := io.ReadAll(captionHTTPResp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read caption content for track %s: %w", chosenCaptionID, err)
	}

	return string(captionBytes), nil
}
