package canvas

import (
	"fmt"
	"sync"
	"time"

	"github.com/bep/debounce"
	"github.com/nathanielfernandes/cnvs/canvas/pb"
	"github.com/nathanielfernandes/cnvs/token"
)

var canvasRequests = make(chan string)
var pendingCanvases = make(map[string]chan *pb.CanvasResponse_Canvas)
var debounced = debounce.New(200 * time.Millisecond)

var cacheLock = sync.RWMutex{}
var cachedCanvases = make(map[string]*pb.CanvasResponse_Canvas)

func StartCanvasRunner() {
	if token.ACCESS_TOKEN == "" {
		panic("No access token, cannot start canvas runner")
	}

	pending := []string{}

	// Listen for new requests
	go func() {
		for trackURI := range canvasRequests {
			pending = append(pending, trackURI)
			// Debounce to bactch requests
			debounced(func() {
				// Make a new request
				go func() {
					canvases, err := FetchCanvasesMapped(pending)

					if err != nil {
						fmt.Println(err.Error())
					}

					// Send the response to all the pending requests
					for trackURI, canvas := range canvases {
						pendingCanvases[trackURI] <- canvas

						// Cache the canvas
						cacheLock.Lock()
						cachedCanvases[trackURI] = canvas
						cacheLock.Unlock()
					}

					// Clear the pending requests
					pending = []string{}
				}()
			})
		}
	}()
}

func CacheGet(trackURI string) (*pb.CanvasResponse_Canvas, bool) {
	cacheLock.RLock()
	canvas, ok := cachedCanvases[trackURI]
	cacheLock.RUnlock()
	return canvas, ok
}

func GetCanvas(trackURI string) (*pb.CanvasResponse_Canvas, error) {
	// If the canvas is cached, just return it
	if canvas, ok := CacheGet(trackURI); ok {
		fmt.Println("Using cached canvas for", trackURI)
		return canvas, nil
	}

	// If there's already a pending request for this track, just wait for it
	if ch, ok := pendingCanvases[trackURI]; ok {
		fmt.Println("Waiting for pending request for", trackURI)
		return <-ch, nil
	}

	fmt.Println("Making new request for", trackURI)
	// Otherwise, make a new request
	ch := make(chan *pb.CanvasResponse_Canvas, 1)
	pendingCanvases[trackURI] = ch

	// Make a new request
	canvasRequests <- trackURI

	return <-ch, nil
}
