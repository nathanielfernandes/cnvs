package preview

import (
	"sync"
)

var cacheLock = sync.RWMutex{}
var cachedPreviews = make(map[string]PreviewResponse)

func GetPreview(trackURI string) (PreviewResponse, error) {
	cacheLock.RLock()

	if preview, ok := cachedPreviews[trackURI]; ok {
		cacheLock.RUnlock()
		return preview, nil
	}
	cacheLock.RUnlock()

	preview, err := ScrapeTrackPreview(trackURI)

	if err != nil {
		return PreviewResponse{}, err
	}

	cacheLock.Lock()
	cachedPreviews[trackURI] = preview
	cacheLock.Unlock()

	return preview, nil
}
