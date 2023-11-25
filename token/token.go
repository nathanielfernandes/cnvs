package token

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const ACCESS_TOKEN_URL = "https://open.spotify.com/get_access_token?reason=transport&productType=web_player"

// no need for a mutex due to infrequent writes
var ACCESS_TOKEN = ""

func setAccessToken() uint {
	// make request to get canvas token
	fmt.Println("Refreshing access token")

	resp, err := http.Get(ACCESS_TOKEN_URL)
	if err != nil {
		panic("Error getting access token: " + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic("Error getting access token: " + resp.Status)
	}

	// decode response
	var token accessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		panic("Error decoding access token response: " + err.Error())
	}

	// set canvas token
	ACCESS_TOKEN = token.AccessToken

	return token.ExpiresOn
}

type accessTokenResponse struct {
	ClientID    string `json:"clientId"`
	AccessToken string `json:"accessToken"`
	ExpiresOn   uint   `json:"accessTokenExpirationTimestampMs"`
	IsAnonymous bool   `json:"isAnonymous"`
}

func StartAccessTokenReferesher() {
	fmt.Println("Starting access token refresh")
	// set timer to refresh canvas token
	expiresTimeStampMs := setAccessToken()

	// set timer to refresh canvas token
	go func() {
		for {
			// wait until token expires
			// convert to seconds
			expiresOn := expiresTimeStampMs / 1000
			// get current time
			now := uint(time.Now().UnixMilli() / 1000)
			// wait until token expires
			waitTime := expiresOn - now
			// wait until token expires
			time.Sleep(time.Duration(waitTime) * time.Second)
			// refresh token
			expiresTimeStampMs = setAccessToken()
		}
	}()
}
