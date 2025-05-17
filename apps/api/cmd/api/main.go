package main

import (
	"github.com/rs/zerolog"
)

// @title SwampHacks Core API
// @version 1.0
// @description This is the API for SwampHacks. It handles all of the automation and management for SwampHacks's events.
// @termsOfService https://app.swamphacks.com/terms

// @contact.name API Support
// @contact.url http://app.swamphacks.com/support
// @contact.email tech@swamphacks.com

// @host api.swamphacks.com
// @BasePath /v1
func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

}
