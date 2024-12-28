package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

type TokenResponseFromGoogleAuth struct {
	Email      string
	Name       string
	Error      error
	StatusCode int
}

func verifyGoogleAuthToken(authToken string, responseChannel chan<- TokenResponseFromGoogleAuth) {
	// Initialize the Google provider
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			"http://localhost:8080/callback",
			"profile", "email",
		),
	)

	// Get the Google provider
	provider, err := goth.GetProvider("google")
	if err != nil {
		responseChannel <- TokenResponseFromGoogleAuth{
			Error:      fmt.Errorf("failed to get provider: %v", err),
			StatusCode: http.StatusInternalServerError,
		}
		return
	}

	_, ok := provider.(*google.Provider)
	if !ok {
		responseChannel <- TokenResponseFromGoogleAuth{
			Error:      fmt.Errorf("failed to cast provider to google provider"),
			StatusCode: http.StatusInternalServerError,
		}
		return
	}

	// Create a new session
	session := google.Session{
		AccessToken: authToken,
	}

	// Fetch the user info using the session
	user, err := provider.FetchUser(&session)
	if err != nil {
		responseChannel <- TokenResponseFromGoogleAuth{
			Error:      fmt.Errorf("failed to fetch user: %v", err),
			StatusCode: http.StatusUnauthorized,
		}
		return
	}

	responseChannel <- TokenResponseFromGoogleAuth{
		Email:      user.Email,
		Name:       user.Name,
		Error:      nil,
		StatusCode: http.StatusOK,
	}
}
