package preview

import (
	"sync"

	"github.com/nathanielfernandes/cnvs/lib/token"
)

var previewRequests = make(chan string)
var pendingPreviews = make(map[string]chan string)

var cacheLock = sync.RWMutex{}
var cachedPreviews = make(map[string]string)

func StartPreviewRunner() {
	if token.ACCESS_TOKEN == "" {
		panic("No access token, cannot start preview runner")
	}

	// Listen for new requests
	go func() {
		for trackURI := range previewRequests {
			// Make a new request
			preview, err := FetchTrackPreviewUrl(trackURI)

			if err != nil {
				println(err.Error())
			}

			pendingPreviews[trackURI] <- preview
		}
	}()
}

func CacheGet(trackURI string) (string, bool) {
	cacheLock.RLock()
	preview, ok := cachedPreviews[trackURI]
	cacheLock.RUnlock()
	return preview, ok
}

func GetPreview(trackURI string) (string, error) {
	// Check the cache
	if preview, ok := CacheGet(trackURI); ok {
		return preview, nil
	}

	// Make a new request
	previewChan := make(chan string)
	pendingPreviews[trackURI] = previewChan
	previewRequests <- trackURI

	return <-previewChan, nil
}
