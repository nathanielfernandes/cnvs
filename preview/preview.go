package preview

import (
	"sync"

	"github.com/nathanielfernandes/cnvs/token"
)

var previewRequests = make(chan string)
var pendingPreviews = make(map[string]chan PreviewResponse)

var cacheLock = sync.RWMutex{}
var cachedPreviews = make(map[string]PreviewResponse)

func StartPreviewRunner() {
	if token.ACCESS_TOKEN == "" {
		panic("No access token, cannot start preview runner")
	}

	// Listen for new requests
	go func() {
		for trackURI := range previewRequests {
			// Make a new request
			preview, err := FetchTrackPreview(trackURI)

			if err != nil {
				println(err.Error())
			}

			pendingPreviews[trackURI] <- preview
		}
	}()
}

func StartScrapeRunner() {
	// Listen for new requests
	go func() {
		for trackURI := range previewRequests {
			// Make a new request
			preview, err := ScrapeTrackPreview(trackURI)

			if err != nil {
				println(err.Error())
			}

			pendingPreviews[trackURI] <- preview
		}
	}()
}

func CacheGet(trackURI string) (PreviewResponse, bool) {
	cacheLock.RLock()
	preview, ok := cachedPreviews[trackURI]
	cacheLock.RUnlock()
	return preview, ok
}

func GetPreview(trackURI string) (PreviewResponse, error) {
	// Check the cache
	if preview, ok := CacheGet(trackURI); ok {
		return preview, nil
	}

	// Make a new request
	previewChan := make(chan PreviewResponse)
	pendingPreviews[trackURI] = previewChan
	previewRequests <- trackURI

	return <-previewChan, nil
}
