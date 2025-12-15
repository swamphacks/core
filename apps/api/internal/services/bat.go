package services

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/swamphacks/core/apps/api/internal/bat"
	"github.com/swamphacks/core/apps/api/internal/db/repository"
	"github.com/swamphacks/core/apps/api/internal/db/sqlc"
	"github.com/swamphacks/core/apps/api/internal/tasks"
)

var (
	ErrListApplicationsFailure = errors.New("Failed to retrieve applications")
	ErrMissingRatings          = errors.New("Some applications are missing their review ratings")
	ErrAddResultFailure        = errors.New("Failed to add generated result")
)

type BatService struct {
	engine         *bat.BatEngine
	appRepo        *repository.ApplicationRepository
	eventRepo      *repository.EventRepository
	batResultRepo  *repository.BatResultsRepository
	batResultsRepo *repository.BatResultsRepository
	taskQueue      *asynq.Client
	logger         zerolog.Logger
}

func NewBatService(engine *bat.BatEngine, appRepo *repository.ApplicationRepository, eventRepo *repository.EventRepository, taskQueue *asynq.Client, logger zerolog.Logger) *BatService {
	return &BatService{
		engine:    engine,
		taskQueue: taskQueue,
		appRepo:   appRepo,
		eventRepo: eventRepo,
		logger:    logger.With().Str("service", "Bat Service").Str("component", "admissions").Logger(),
	}
}

func (s *BatService) AddResult(ctx context.Context, eventID uuid.UUID, acceptedApplicants []uuid.UUID, rejectedApplicants []uuid.UUID) (*sqlc.BatResult, error) {
	params := sqlc.AddResultParams{
		EventID:            eventID,
		AcceptedApplicants: acceptedApplicants,
		RejectedApplicants: rejectedApplicants,
	}

	result, err := s.batResultRepo.AddResult(ctx, params)
	if err != nil && errors.Is(err, repository.ErrDuplicateResult) {
		s.logger.Err(err).Msg("Could not insert result as it already exists.")
		return nil, ErrResultConflict
	} else if err != nil {
		s.logger.Err(err).Msg("An unknown error was caught!")
		return nil, ErrFailedToCreateResult
	}

	return result, nil
}

func (s *BatService) QueueCalculateAdmissionsTask(eventId uuid.UUID) (*asynq.TaskInfo, error) {
	task, err := tasks.NewTaskCalculateAdmissions(tasks.CalculateAdmissionsPayload{
		EventID: eventId,
	})
	if err != nil {
		s.logger.Err(err).Msg("Failed to create CalculateAdmissions task")
		return nil, err
	}

	info, err := s.taskQueue.Enqueue(task, asynq.Queue("bat"))
	if err != nil {
		s.logger.Err(err).Msg("Failed to queue CalculateAdmissions task")
		return nil, err
	}

	return info, nil
}

type ApplicantAdmissionsData struct {
	ID            uuid.UUID
	TeamID        uuid.NullUUID
	WeightedScore float64
	SortKey       float64
	IsUFStudent   bool // Is a University of Florida student?
	IsEarlyCareer bool // Is a freshman to sophomore?
}

// Unmarshal into this reduced struct since
// our application fields are fucked
// beyond all fuck man. TODO: FIX OUR APP FIELDS!
type ApplicationSchoolAndYear struct {
	School string `json:"school"`
	Year   string `json:"year"`
}

type QuotaState struct {
	TotalAccepted  int32
	TeamSlotsLeft  int32
	UFEarlyLeft    int32
	UFLateLeft     int32
	OtherEarlyLeft int32
	OtherLateLeft  int32
}

type TeamAdmissionData struct {
	TeamID               uuid.UUID
	MembersAdmissionData []ApplicantAdmissionsData
	AverageWeightedScore float64
	SortKey              float64
}

func (s *BatService) CalculateAdmissions(ctx context.Context, eventId uuid.UUID) error {
	s.logger.Info().Str("eventId", eventId.String()).Msg("")

	// First aggregate data necessary
	applications, err := s.appRepo.ListAdmissionCandidatesByEvent(ctx, eventId)
	if err != nil || len(applications) == 0 {
		return ErrListApplicationsFailure
	}

	var appAdmissionsData []ApplicantAdmissionsData
	for _, app := range applications {
		if app.ExperienceRating == nil || app.PassionRating == nil {
			return ErrMissingRatings
		}

		var applicationShoolAndYear ApplicationSchoolAndYear
		if err := json.Unmarshal(app.Application, &applicationShoolAndYear); err != nil {
			s.logger.Debug().Bytes("App", app.Application).Msg("Application data")
			return err
		}

		var teamId uuid.UUID
		if app.TeamID != nil {
			teamId = *app.TeamID
		}

		wScore, err := s.engine.CalculateWeightedScore(*app.PassionRating, *app.ExperienceRating)
		if err != nil {
			return err
		}
		appAdmissionsData = append(appAdmissionsData, ApplicantAdmissionsData{
			ID: app.UserID,
			TeamID: uuid.NullUUID{
				UUID:  teamId,
				Valid: app.TeamID != nil,
			},
			WeightedScore: wScore,
			SortKey:       0.0,
			IsUFStudent:   applicationShoolAndYear.School == "University of Florida",
			IsEarlyCareer: applicationShoolAndYear.Year == "first_year" || applicationShoolAndYear.Year == "second_year",
		})
	}

	quota := QuotaState{
		TotalAccepted:  0,
		TeamSlotsLeft:  50,
		UFEarlyLeft:    210,
		UFLateLeft:     140,
		OtherEarlyLeft: 90,
		OtherLateLeft:  60,
	}

	teams, solo := groupAndSortTeams(appAdmissionsData)
	applyTeamSortKey(teams)
	acceptedTeams, remaining, quota := admitTeams(teams, solo, quota)
	accepted, rejected, quota := admitSoloApplicants(remaining, quota)

	var acceptedUUIDs []uuid.UUID
	var rejectedUUIDs []uuid.UUID

	for _, applicant := range accepted {
		acceptedUUIDs = append(acceptedUUIDs, applicant.ID)
	}

	for _, applicant := range rejected {
		rejectedUUIDs = append(rejectedUUIDs, applicant.ID)
	}
	_, err = s.AddResult(ctx, eventId, acceptedUUIDs, rejectedUUIDs)

	if err != nil {
		return ErrAddResultFailure
	}

	s.logger.Info().Int("Accepted", len(accepted)+len(acceptedTeams)).Int("Rejected", len(rejected)).Msg("Finished Algo")

	return nil
}

func admitTeams(teams []TeamAdmissionData, solo []ApplicantAdmissionsData, initialQuota QuotaState) (
	[]ApplicantAdmissionsData, // Accepted Applicants
	[]ApplicantAdmissionsData, // remaining applicants (joined with solo)
	QuotaState, // Updated Quote
) {
	admittedApplicants := make([]ApplicantAdmissionsData, 0, initialQuota.TotalAccepted)
	remainingApplicants := make([]ApplicantAdmissionsData, 0)
	remainingApplicants = append(remainingApplicants, solo...)
	quota := initialQuota

	for _, team := range teams {
		// All remaining members get appended to solo/remaining selection
		if quota.TeamSlotsLeft <= int32(len(team.MembersAdmissionData)) {
			remainingApplicants = append(remainingApplicants, team.MembersAdmissionData...)
			continue
		}

		required := countTeamSlots(team.MembersAdmissionData)
		if canAdmitTeam(required, quota) {
			admittedApplicants = append(admittedApplicants, team.MembersAdmissionData...)

			quota.TotalAccepted += int32(len(team.MembersAdmissionData))
			quota.TeamSlotsLeft -= int32(len(team.MembersAdmissionData))

			quota.UFEarlyLeft -= required.UFEarlyLeft
			quota.UFLateLeft -= required.UFLateLeft
			quota.OtherEarlyLeft -= required.OtherEarlyLeft
			quota.OtherLateLeft -= required.OtherLateLeft
		} else {
			remainingApplicants = append(remainingApplicants, team.MembersAdmissionData...)
		}
	}

	return admittedApplicants, remainingApplicants, quota
}

type BucketConfig struct {
	Name        string
	QuotaPtr    *int32
	RolloverPtr *int32
}

func admitSoloApplicants(solo []ApplicantAdmissionsData, quota QuotaState) (
	[]ApplicantAdmissionsData, // Accepted
	[]ApplicantAdmissionsData, // Rejected
	QuotaState,
) {
	var admittedSolo []ApplicantAdmissionsData

	pool := make(map[uuid.UUID]ApplicantAdmissionsData)
	for _, app := range solo {
		pool[app.ID] = app
	}

	buckets := []BucketConfig{
		{"UF_Early", &quota.UFEarlyLeft, &quota.UFLateLeft},
		{"UF_Late", &quota.UFLateLeft, &quota.UFEarlyLeft},
		{"Other_Early", &quota.OtherEarlyLeft, &quota.OtherLateLeft},
		{"Other_Late", &quota.OtherLateLeft, &quota.OtherEarlyLeft},
	}

	for {
		passAdmittedCount := int32(0)

		for _, bucket := range buckets {
			originalCount := *bucket.QuotaPtr
			if originalCount <= 0 {
				continue
			}

			bucketPool := filterApplicants(pool, bucket.Name)
			applyIndividualSortKey(bucketPool)
			sort.Slice(bucketPool, func(i, j int) bool {
				return bucketPool[i].SortKey > bucketPool[j].SortKey
			})

			numToSelect := min(originalCount, int32(len(bucketPool)))
			for i := range numToSelect {

				app := bucketPool[i]
				admittedSolo = append(admittedSolo, app)

				*bucket.QuotaPtr -= 1
				quota.TotalAccepted += 1
				delete(pool, app.ID)

				passAdmittedCount++
			}

			slotsUnfilled := *bucket.QuotaPtr
			if slotsUnfilled > 0 && bucket.RolloverPtr != nil {
				*bucket.RolloverPtr += slotsUnfilled

				*bucket.QuotaPtr = 0

			}
		}

		if passAdmittedCount == 0 {
			break // No admissions in a pass anymore!
		}
	}

	var rejected []ApplicantAdmissionsData
	for _, applicant := range pool {
		rejected = append(rejected, applicant)
	}

	return admittedSolo, rejected, quota

	// for bucket, targetCount := range buckets {
	// 	count := *targetCount

	// 	if count <= 0 {
	// 		continue // Skip if bucket is already full/completed
	// 	}

	// 	bucketPool := filterApplicants(pool, bucket)
	// 	applyIndividualSortKey(bucketPool)

	// 	// Sort
	// 	sort.Slice(bucketPool, func(i, j int) bool {
	// 		return bucketPool[i].SortKey > bucketPool[j].SortKey
	// 	})

	// 	numToSelect := int(min(count, int32(len(bucketPool))))

	// 	for i := range numToSelect {
	// 		applicant := bucketPool[i]
	// 		admittedSolo = append(admittedSolo, applicant)

	// 		*targetCount -= 1
	// 		quota.TotalAccepted += 1
	// 		delete(pool, applicant.ID)
	// 	}
	// }

	// var rejectedSolo []ApplicantAdmissionsData
	// for _, applicant := range pool {
	// 	rejectedSolo = append(rejectedSolo, applicant)
	// }

	// return admittedSolo, rejectedSolo, quota
}

func applyIndividualSortKey(apps []ApplicantAdmissionsData) {
	var Epsilon float64 = 0.001

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range apps {
		app := &apps[i]
		randVal := r.Float64()
		exp := 1.0 / (app.WeightedScore + Epsilon)

		app.SortKey = math.Pow(randVal, exp)
	}
}

func filterApplicants(
	pool map[uuid.UUID]ApplicantAdmissionsData,
	bucketName string,
) []ApplicantAdmissionsData {

	var filtered []ApplicantAdmissionsData

	for _, app := range pool {
		switch bucketName {
		case "UF_Early":
			if app.IsUFStudent && app.IsEarlyCareer {
				filtered = append(filtered, app)
			}
		case "UF_Late":
			if app.IsUFStudent && !app.IsEarlyCareer {
				filtered = append(filtered, app)
			}
		case "Other_Early":
			if !app.IsUFStudent && app.IsEarlyCareer {
				filtered = append(filtered, app)
			}
		case "Other_Late":
			if !app.IsUFStudent && !app.IsEarlyCareer {
				filtered = append(filtered, app)
			}
		default:
			// Should not happen if buckets map is defined correctly
		}
	}
	return filtered
}
func canAdmitTeam(req QuotaState, curr QuotaState) bool {
	return req.UFEarlyLeft <= curr.UFEarlyLeft &&
		req.UFLateLeft <= curr.UFLateLeft &&
		req.OtherEarlyLeft <= curr.OtherEarlyLeft &&
		req.OtherLateLeft <= curr.OtherLateLeft
}

func countTeamSlots(members []ApplicantAdmissionsData) QuotaState {
	reqQuota := QuotaState{
		TotalAccepted:  0,
		TeamSlotsLeft:  0,
		UFEarlyLeft:    0,
		UFLateLeft:     0,
		OtherEarlyLeft: 0,
		OtherLateLeft:  0,
	}

	for _, member := range members {
		if member.IsUFStudent {
			if member.IsEarlyCareer {
				reqQuota.UFEarlyLeft += 1
			} else {
				reqQuota.UFLateLeft += 1
			}
		} else {
			if member.IsEarlyCareer {
				reqQuota.OtherEarlyLeft += 1
			} else {
				reqQuota.OtherLateLeft += 1
			}
		}

		reqQuota.TotalAccepted += 1
		reqQuota.TeamSlotsLeft += 1
	}

	return reqQuota
}

func applyTeamSortKey(teams []TeamAdmissionData) {
	const Epsilon float64 = 0.001

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range teams {
		// Get reference to apply in place
		team := &teams[i]

		randVal := r.Float64()

		exponent := 1.0 / (team.AverageWeightedScore + Epsilon)
		team.SortKey = math.Pow(randVal, exponent)
	}

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].SortKey > teams[j].SortKey
	})
}

func groupAndSortTeams(data []ApplicantAdmissionsData) ([]TeamAdmissionData, []ApplicantAdmissionsData) {
	teamMap := make(map[uuid.UUID][]ApplicantAdmissionsData)
	var soloApplicants []ApplicantAdmissionsData
	for _, app := range data {
		if app.TeamID.Valid {
			teamMap[app.TeamID.UUID] = append(teamMap[app.TeamID.UUID], app)
		} else {
			soloApplicants = append(soloApplicants, app)
		}
	}

	var teams []TeamAdmissionData
	for teamID, members := range teamMap {
		var totalScore float64
		for _, member := range members {
			totalScore += member.WeightedScore
		}

		teams = append(teams, TeamAdmissionData{
			TeamID:               teamID,
			MembersAdmissionData: members,
			AverageWeightedScore: totalScore / float64(len(members)),
		})
	}

	return teams, soloApplicants
}

func (s *BatService) LogAdmissionsStats(data []ApplicantAdmissionsData) {
	type scoreStats struct {
		count int
		sum   float64
		min   float64
		max   float64
	}

	newStats := func() scoreStats {
		return scoreStats{
			min: math.Inf(1),
			max: math.Inf(-1),
		}
	}

	update := func(s *scoreStats, score float64) {
		s.count++
		s.sum += score
		if score < s.min {
			s.min = score
		}
		if score > s.max {
			s.max = score
		}
	}

	finalize := func(s scoreStats) map[string]any {
		if s.count == 0 {
			return map[string]any{
				"count": 0,
			}
		}
		return map[string]any{
			"count": s.count,
			"avg":   s.sum / float64(s.count),
			"min":   s.min,
			"max":   s.max,
		}
	}

	var (
		total          = newStats()
		uf             = newStats()
		nonUF          = newStats()
		earlyCareer    = newStats()
		upperCareer    = newStats()
		teamApplicants = newStats()
		soloApplicants = newStats()
	)

	for _, a := range data {
		update(&total, a.WeightedScore)

		if a.IsUFStudent {
			update(&uf, a.WeightedScore)
		} else {
			update(&nonUF, a.WeightedScore)
		}

		if a.IsEarlyCareer {
			update(&earlyCareer, a.WeightedScore)
		} else {
			update(&upperCareer, a.WeightedScore)
		}

		if a.TeamID.Valid {
			update(&teamApplicants, a.WeightedScore)
		} else {
			update(&soloApplicants, a.WeightedScore)
		}
	}

	s.logger.Info().
		Int("total_applicants", len(data)).
		Fields(map[string]any{
			"overall":         finalize(total),
			"uf_students":     finalize(uf),
			"non_uf_students": finalize(nonUF),
			"early_career":    finalize(earlyCareer),
			"upper_career":    finalize(upperCareer),
			"team_applicants": finalize(teamApplicants),
			"solo_applicants": finalize(soloApplicants),
		}).
		Msg("Admissions statistics snapshot")
}
