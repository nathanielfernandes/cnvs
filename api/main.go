package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nathanielfernandes/cnvs/canvas"
	"github.com/nathanielfernandes/cnvs/preview"
	"github.com/nathanielfernandes/cnvs/token"
)

func main() {
	token.StartAccessTokenReferesher()
	canvas.StartCanvasRunner()
	preview.StartPreviewRunner()

	http.HandleFunc("/canvas/", GetCanvas)
	http.HandleFunc("/r-canvas/", RedirectToCanvas)
	http.HandleFunc("/preview/", GetPreview)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}
}
