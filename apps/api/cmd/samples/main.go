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

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	earlyAppOpenTime := time.Date(2026, time.June, 20, 19, 13, 20, 0, loc)
	earlyAppCloseTime := time.Date(2026, time.August, 21, 23, 59, 59, 0, loc)
	appOpenTime := time.Date(2026, time.August, 22, 0, 0, 0, 0, loc)
	appCloseTime := time.Date(2026, time.September, 4, 23, 59, 59, 0, loc)
	startTime := time.Date(2026, time.October, 16, 19, 13, 20, 0, loc)
	endTime := time.Date(2026, time.October, 18, 19, 13, 20, 0, loc)

	_, err = hackathonRepo.CreateHackathon(context.TODO(), sqlc.CreateHackathonParams{
		ID:                      "xii",
		Name:                    "SwampHacks XII",
		ApplicationOpen:         appOpenTime,
		ApplicationClose:        appCloseTime,
		StartTime:               startTime,
		EndTime:                 endTime,
		Description:             "SwampHacks' 12th iteration",
		Location:                "Reitz Union",
		LocationUrl:             nil,
		MaxAttendees:            nil,
		RsvpDeadline:            nil,
		DecisionRelease:         nil,
		IsActive:                true,
		AcceptEarlyApplications: true,
		EarlyApplicationOpen:    earlyAppOpenTime,
		EarlyApplicationClose:   earlyAppCloseTime,
	})

	if err != nil {
		fmt.Println("something went wrong")
		fmt.Println(err)
	}
}
