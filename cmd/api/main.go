package main

import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Rihal's API!"))
	})
	//Starting the HTTP server
	http.ListenAndServe(":8080", nil)
}
