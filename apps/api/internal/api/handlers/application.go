package handlers

import "github.com/swamphacks/core/apps/api/internal/services"

type ApplicationHandler struct {
	appService *services.ApplicationService
}

func NewApplicationHandler(appService *services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		appService: appService,
	}
}
