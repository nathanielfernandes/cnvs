package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nathanielfernandes/cnvs/canvas"
)

func main() {
	canvas.StartCanvasRunner()

	http.HandleFunc("/canvas/", GetCanvas)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}
}
