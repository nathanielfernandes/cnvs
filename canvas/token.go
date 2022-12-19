package canvas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const CANVAS_TOKEN_URL = "https://open.spotify.com/get_access_token?reason=transport&productType=web_player"

// no need for a mutex due to infrequent writes
var CANVAS_TOKEN = ""

func setCanvasToken() uint {
	// make request to get canvas token
	fmt.Println("Refreshing canvas token")

	resp, err := http.Get(CANVAS_TOKEN_URL)
	if err != nil {
		panic("Error getting canvas token: " + err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic("Error getting canvas token: " + resp.Status)
	}

	// decode response
	var token canvasTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		panic("Error decoding canvas token response: " + err.Error())
	}

	// set canvas token
	CANVAS_TOKEN = token.AccessToken

	return token.ExpiresOn
}

type canvasTokenResponse struct {
	ClientID    string `json:"clientId"`
	AccessToken string `json:"accessToken"`
	ExpiresOn   uint   `json:"accessTokenExpirationTimestampMs"`
	IsAnonymous bool   `json:"isAnonymous"`
}

func startCanvasTokenReferesher() {
	fmt.Println("Starting canvas token refresh")
	// set timer to refresh canvas token
	expiresTimeStampMs := setCanvasToken()

	go func() {
		for {
			// sleep until token expires, with a 30 second buffer
			time.Sleep(time.Duration(expiresTimeStampMs - 30000))

			// wait until token expires
			expiresTimeStampMs = setCanvasToken()
			fmt.Println("Canvas token refreshed")
		}
	}()
}
