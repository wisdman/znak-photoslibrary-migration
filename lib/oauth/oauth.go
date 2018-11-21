package oauth

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var clientID string = os.Getenv("GOOGLE_CLIENT_ID")
var clientSecret string = os.Getenv("GOOGLE_CLIENT_SECRET")

func init() {
	if clientID == "" || clientSecret == "" {
		log.Fatal(`Error: GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables must be set.`)
	}
}

func generateOAuthState() (string, error) {
	var n uint64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &n); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", n), nil
}

func requestOAuthCode(url string) (string, error) {
	log.Printf("Open %s", url)
	fmt.Print("Enter code: ")

	var code string
	_, err := fmt.Scanln(&code)

	if err != nil {
		return "", err
	}

	return code, nil
}

func New(ctx context.Context, scopes []string) (*http.Client, error) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       scopes,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
	}

	authState, err := generateOAuthState()
	if err != nil {
		return nil, err
	}

	authCodeURL := config.AuthCodeURL(authState)

	authCode, err := requestOAuthCode(authCodeURL)
	if err != nil {
		return nil, err
	}

	accessToken, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, err
	}

	return config.Client(ctx, accessToken), nil
}
