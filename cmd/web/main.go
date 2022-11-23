package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

var (
	path = flag.String("path", ".", "path to the folder to serve. Defaults to the current folder")
	port = flag.String("port", "8080", "port to serve on. Defaults to 8080")
)

func main() {

	mux := Routes()

	fmt.Println("Server listening on port :8080")

	flag.Parse()

	dirname, err := filepath.Abs(*path)
	if err != nil {
		log.Fatalf("Could not get absolute path to directory: %s: %s", dirname, err.Error())
	}

	log.Printf("Serving %s on port %s", dirname, *port)

	err = Serve(dirname, *port, mux)
	if err != nil {
		log.Fatalf("Could not serve directory: %s: %s", dirname, err.Error())
	}

}

func Serve(dirname string, port string, mux http.Handler) error {
	fs := http.FileServer(http.Dir(dirname))
	http.Handle("/", fs)

	return http.ListenAndServe(":"+port, mux)
}
