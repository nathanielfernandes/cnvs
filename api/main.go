package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nathanielfernandes/cnvs/token"
)

func main() {
	token.StartAccessTokenReferesher()

	http.HandleFunc("/canvas/", GetCanvas)
	http.HandleFunc("/r-canvas/", RedirectToCanvas)
	http.HandleFunc("/preview/", GetPreview)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}

	// after 2 hours kill the server to force a restart
	time.Sleep(2 * time.Hour)
	panic("Server killed")
}
