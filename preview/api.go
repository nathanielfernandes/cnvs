package preview

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/nathanielfernandes/cnvs/token"
	"github.com/tidwall/gjson"
)

const RAW_TRACK_URL = "https://api.spotify.com/v1/tracks/"

var TRACK_URL, _ = url.Parse(RAW_TRACK_URL)

func FetchTrack(trackId string) ([]byte, error) {
	req, err := http.NewRequest("GET", RAW_TRACK_URL+trackId, nil)
	req.Header.Add("Authorization", "Bearer "+token.ACCESS_TOKEN)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	return bytes, err
}

func FetchTrackPreview(trackId string) (PreviewResponse, error) {
	rtrack, err := FetchTrack(trackId)
	if err != nil {
		return PreviewResponse{}, err
	}

	track := string(rtrack)
	preview := gjson.Get(track, "preview_url").String()
	cover_art := gjson.Get(track, "album.images.0.url").String()

	return PreviewResponse{
		AudioURL:    preview,
		CoverArtURL: cover_art,
	}, nil
}
