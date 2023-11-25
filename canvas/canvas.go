package canvas

import (
	"fmt"
	"sync"

	"github.com/nathanielfernandes/cnvs/canvas/pb"
)

type CachedCanvasResponse struct {
	Canvas  *pb.CanvasResponse_Canvas
	Fetched bool
	Lock    sync.RWMutex
}

func (ccr *CachedCanvasResponse) Read() (*pb.CanvasResponse_Canvas, bool) {
	ccr.Lock.RLock()
	defer ccr.Lock.RUnlock()
	return ccr.Canvas, ccr.Fetched
}

var cacheLock = sync.RWMutex{}
var cached = make(map[string]*CachedCanvasResponse)

func GetCanvas(trackURI string) (*pb.CanvasResponse_Canvas, error) {
	cacheLock.RLock()

	if cachedCanvas, ok := cached[trackURI]; ok {
		if canvas, fetched := cachedCanvas.Read(); fetched {
			fmt.Println("Found cached canvas for", trackURI)
			cacheLock.RUnlock()
			return canvas, nil
		}
	}
	cacheLock.RUnlock()

	fmt.Println("Fetching canvas for", trackURI)

	cacheLock.Lock()
	canvas := CachedCanvasResponse{
		Canvas:  nil,
		Fetched: false,
		Lock:    sync.RWMutex{},
	}
	cached[trackURI] = &canvas
	cacheLock.Unlock()

	canvas.Lock.Lock()
	defer canvas.Lock.Unlock()

	canvases, err := FetchCanvases([]string{trackURI})

	if err != nil {
		return nil, err
	}

	if len(canvases) == 0 {
		return nil, nil
	}

	canvas.Canvas = canvases[0]
	canvas.Fetched = true

	return canvases[0], nil
}
