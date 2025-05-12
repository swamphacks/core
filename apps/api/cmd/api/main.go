package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("Request to [GET] / endpoint!")
	fmt.Fprintf(w, "Hello!\n")
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	http.HandleFunc("/", handler)

	fmt.Println("Server running on 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server.")
	}
}
