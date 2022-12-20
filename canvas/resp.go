package canvas

type CanvasTokenResponse struct {
	ClientID    string `json:"clientId"`
	AccessToken string `json:"accessToken"`
	ExpiresOn   uint   `json:"accessTokenExpirationTimestampMs"`
	IsAnonymous bool   `json:"isAnonymous"`
}
