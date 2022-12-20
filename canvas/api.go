package canvas

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/nathanielfernandes/cnvs/canvas/pb"
	"github.com/nathanielfernandes/cnvs/token"
	"google.golang.org/protobuf/proto"
)

const RAW_CANVASES_URL = "https://spclient.wg.spotify.com/canvaz-cache/v0/canvases"

var CANVASES_URL, _ = url.Parse(RAW_CANVASES_URL)

func FetchCanvases(trackURIs []string) ([]*pb.CanvasResponse_Canvas, error) {
	cr := pb.CanvasRequest{}

	for _, uri := range trackURIs {
		track := &pb.CanvasRequest_Track{
			TrackUri: uri,
		}
		cr.Tracks = append(cr.Tracks, track)
	}

	buff, err := proto.Marshal(&cr)
	if err != nil {
		return nil, err
	}

	r := io.NopCloser(bytes.NewReader(buff))

	req := http.Request{
		Method: "POST",
		URL:    CANVASES_URL,
		Body:   r,
		Header: http.Header{
			"content-type":    []string{"application/protobuf"},
			"accept":          []string{"application/protobuf"},
			"accept-language": []string{"en"},
			"user-agent":      []string{"Spotify/8.6.85 iOS/14.4.2 (iPhone12,1)"},
			"accept-encoding": []string{"gzip, deflate, br"},
			"authorization":   []string{"Bearer " + token.ACCESS_TOKEN},
		},
	}

	resp, err := http.DefaultClient.Do(&req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buff, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var canvasResp pb.CanvasResponse
	if err := proto.Unmarshal(buff, &canvasResp); err != nil {
		return nil, err
	}

	return canvasResp.GetCanvases(), err
}

func FetchCanvasesMapped(trackURIs []string) (map[string]*pb.CanvasResponse_Canvas, error) {
	canvases, err := FetchCanvases(trackURIs)

	if err != nil {
		return nil, err
	}

	mapped := make(map[string]*pb.CanvasResponse_Canvas)

	for _, canvas := range canvases {
		mapped[canvas.GetTrackUri()] = canvas
	}

	for _, uri := range trackURIs {
		if _, ok := mapped[uri]; !ok {
			mapped[uri] = nil
		}
	}

	return mapped, nil
}
