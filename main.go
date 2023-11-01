package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func main() {
	const port uint16 = 8080

	http.HandleFunc("/", handler)
	slog.Info(fmt.Sprintf("Start listening on port %d", port))
	slog.Error(http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil).Error())
}
