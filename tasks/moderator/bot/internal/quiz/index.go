package quiz

import (
	_ "embed"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/quiz", IndexHandler)
}

//go:embed index.html
var index []byte

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write(index)
	if err != nil {
		log.Printf("failed send index page: %v", err)
	}
}
