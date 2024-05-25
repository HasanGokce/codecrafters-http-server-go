package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "HTTP/1.1 200 OK\r\n\r\n")

	})

	err := http.ListenAndServe(":4221", nil)
	if err != nil {
		fmt.Println("Failed to start server on port 4221:", err)
	}

}
