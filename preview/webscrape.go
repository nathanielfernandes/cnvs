package preview

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

const RAW_SCRAPE_URL = "https://open.spotify.com/embed/track/"

var SCRAPE_URL, _ = url.Parse(RAW_SCRAPE_URL)

var DATA_REGEX = regexp.MustCompile(`(?s)<script id="__NEXT_DATA__" type="application/json">(.+)</script>`)

func ScrapeTrack(trackId string) (string, error) {
	req, err := http.NewRequest("GET", RAW_SCRAPE_URL+trackId, nil)

	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	str := string(bytes)

	matches := DATA_REGEX.FindStringSubmatch(str)

	if len(matches) < 2 {
		return "", errors.New("no match")
	}

	return matches[1], nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ScrapeTrackPreview(trackId string) (PreviewResponse, error) {
	rtrack, err := ScrapeTrack(trackId)

	if err != nil {
		return PreviewResponse{}, err
	}

	data := gjson.Get(rtrack, "props.pageProps.state.data").String()

	background_color := gjson.Get(data, "backgroundColor").String()
	track_name := gjson.Get(data, "entity.name").String()

	rartists := gjson.Get(data, "entity.artists").Array()
	artists := make([]Artist, len(rartists))
	for i, artist := range rartists {
		m := artist.Map()

		name := m["name"].String()
		url := strings.Split(m["uri"].String(), ":")[2]
		url = "https://open.spotify.com/artist/" + url
		artists[i] = Artist{
			Name: name,
			URL:  url,
		}
	}

	audio_preview_url := gjson.Get(data, "entity.audioPreview.url").String()

	rcover_art := gjson.Get(data, "entity.coverArt.sources").Array()

	cover_art := CoverArt{}

	min_idx := len(rcover_art) - 1

	if len(rcover_art) > 0 {
		cover_art = CoverArt{
			Small:  rcover_art[min_idx].Map()["url"].String(),
			Medium: rcover_art[min(1, min_idx)].Map()["url"].String(),
			Large:  rcover_art[min(2, min_idx)].Map()["url"].String(),
		}
	}

	return PreviewResponse{
		AudioURL:        audio_preview_url,
		CoverArt:        cover_art,
		TrackName:       track_name,
		Artists:         artists,
		BackgroundColor: background_color,
	}, nil

}
