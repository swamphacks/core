package main

import (
	"context"
	"fmt"
	"time"

	"github.com/swamphacks/core/apps/api/internal/database"
	"github.com/swamphacks/core/apps/api/internal/database/repository"
	"github.com/swamphacks/core/apps/api/internal/database/sqlc"
)

func main() {
	db := database.NewDB("postgres://postgres:postgres@localhost:5432/coredb")
	defer db.Close()

	hackathonRepo := repository.NewHackathonRepository(db)

	appOpenTime := time.Date(2026, 3, 26, 19, 13, 20, 0, time.UTC)
	appCloseTime := time.Date(2026, 4, 26, 19, 13, 20, 0, time.UTC)
	startTime := time.Date(2026, 10, 26, 19, 13, 20, 0, time.UTC)
	endTime := time.Date(2026, 10, 29, 19, 13, 20, 0, time.UTC)

	_, err := hackathonRepo.CreateHackathon(context.TODO(), sqlc.CreateHackathonParams{
		ID:               "xii",
		Name:             "SwampHacks XII",
		ApplicationOpen:  appOpenTime,
		ApplicationClose: appCloseTime,
		StartTime:        startTime,
		EndTime:          endTime,
		Description:      "SwampHacks' 12th iteration",
		Location:         "Reitz Union",
		LocationUrl:      nil,
		MaxAttendees:     nil,
		RsvpDeadline:     nil,
		DecisionRelease:  nil,
		IsActive:         true,
	})

	if err != nil {
		fmt.Println("something went wrong")
		fmt.Println(err)
	}
}
