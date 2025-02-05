package main

import (
	"fmt"
	"net/http"
	"os"
)

func ensureDir(dir string) {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		fmt.Printf("Directory '%s' not found.\n", dir)
		return
	}
}

func DefaultTo(value interface{}, defaults interface{}) interface{} {
	if value != nil {
		return value
	}

	return defaults
}

func main() {
	dir := "public"

	if len(os.Args) > 1 {
		dir = DefaultTo(os.Args[1], dir).(string)
	}

	ensureDir(dir)

	// Create a file server handler to serve the directory's contents
	fileServer := http.FileServer(http.Dir(dir))

	// Create a new HTTP server and handle normal requests
	http.Handle("/", fileServer)

	// Start the server on port 8080
	port := 8080
	fmt.Printf("Serving files from directory %s at http://localhost:%d\n", dir, port)
	fmt.Printf("Server started at http://localhost:%d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
