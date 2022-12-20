package main

import (
	"encoding/json"
	"net/http"

	"github.com/nathanielfernandes/cnvs/lib/canvas"
	"github.com/nathanielfernandes/cnvs/lib/preview"
)

func addCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func GetCanvas(w http.ResponseWriter, r *http.Request) {
	addCors(w)
	// trim the /canvas/ part of the path
	track := GetSpotifyUriFromAny(r.URL.Path[8:])

	if track == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid track id"))
		return
	}

	cv, err := canvas.GetCanvas(track)

	if cv == nil || err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Canvas not found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	b, err := json.MarshalIndent(cv, "", "  ")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(b)
}

func RedirectToCanvas(w http.ResponseWriter, r *http.Request) {
	addCors(w)

	// trim the /r-canvas/ part of the path
	track := r.URL.Path[10:]

	if track == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid track id"))
		return
	}

	canvas, err := canvas.GetCanvas(track)

	if canvas == nil || err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Canvas not found"))
		return
	}

	http.Redirect(w, r, canvas.CanvasUrl, http.StatusFound)
}

func GetPreview(w http.ResponseWriter, r *http.Request) {
	addCors(w)

	// trim the /preview/ part of the path
	track := GetSpotifyIdFromAny(r.URL.Path[9:])

	if track == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid track id"))
		return
	}

	url, err := preview.GetPreview(track)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Preview not found"))
		return
	}

	w.Write([]byte(url))
}
