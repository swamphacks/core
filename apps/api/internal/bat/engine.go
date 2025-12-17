package bat

import (
	"errors"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/google/uuid"
)

var (
	ErrImproperWeights  = errors.New("Passion and experience weights don't add to 1.0")
	ErrScoreOutOfBounds = errors.New("The score can only range from 1 to 5")
)

type BucketType int

const (
	BucketTypeUFEarly BucketType = iota
	BucketTypeUFLate
	BucketTypeOtherEarly
	BucketTypeOtherLate
)

type BucketConfig struct {
	Type        BucketType
	QuotaPtr    *int32
	RolloverPtr *int32
}

type AdmissionCandidate struct {
	UserID        uuid.UUID
	TeamID        uuid.NullUUID
	WeightedScore float64
	SortKey       float64
	IsUFStudent   bool
	IsEarlyCareer bool
}

// TeamEvaluationData holds the aggregated metrics for a group of applicants
// being considered as a single unit in the admissions process.
type TeamEvaluationData struct {
	// TeamID is the unique identifier for the team.
	TeamID uuid.UUID

	// Members is the list of all applicants belonging to this team.
	Members []AdmissionCandidate

	// AverageWeightedScore is the average WeightedScore of all team members.
	AverageWeightedScore float64

	// SortKey is the non-deterministic key used for ranking teams against each other.
	SortKey float64
}

type AdmissionContext struct {
	School string `json:"school"`
	Year   string `json:"year"`
}

type CategoryQuota struct {
	// EarlyLeft is the number of remaining slots for Early Career applicants (Freshman/Sophomore).
	EarlyLeft int32

	// LateLeft is the number of remaining slots for Upper Career applicants (Junior/Senior/Grad).
	LateLeft int32
}

// QuotaState tracks the remaining admission slots across all categories
// and for specific groups, ensuring the total capacity and category-specific
// limits are not exceeded during the admission process.
type QuotaState struct {
	// TotalAccepted is the running count of applicants already admitted.
	// This value is primarily informational and doesn't limit acceptance.
	TotalAccepted int32

	// TeamSlotsLeft is the maximum number of remaining *individual* applicants
	// that can be admitted as part of a team (before solo admissions).
	TeamSlotsLeft int32

	// UF holds the remaining quota slots for applicants who are University of Florida students,
	// categorized by career stage (Early vs. Late).
	UF CategoryQuota

	// Other holds the remaining quota slots for applicants who are NOT University of Florida students,
	// categorized by career stage (Early vs. Late).
	Other CategoryQuota
}

type BatEngine struct {
	passionWeight        float64
	experienceWeight     float64
	weightedBaseConstant float64
	Quota                QuotaState
}

func NewBatEngine(passionW, experienceW float64) (*BatEngine, error) {
	if !equalWithinTolerance(passionW+experienceW, 1.0, 1e-9) {
		return nil, ErrImproperWeights
	}

	// Allow this to be set outside at some point
	quota := QuotaState{
		TotalAccepted: 0,
		TeamSlotsLeft: 50,
		UF: CategoryQuota{
			EarlyLeft: 210,
			LateLeft:  140,
		},
		Other: CategoryQuota{
			EarlyLeft: 90,
			LateLeft:  60,
		},
	}

	return &BatEngine{
		passionWeight:        passionW,
		experienceWeight:     experienceW,
		weightedBaseConstant: 0.1,
		Quota:                quota,
	}, nil
}

func (b *BatEngine) CalculateWeightedScore(passionS, expS int32) (float64, error) {
	if 5 < passionS || 0 > passionS {
		return 0.0, ErrScoreOutOfBounds
	}

	if 5 < expS || 0 > expS {
		return 0.0, ErrScoreOutOfBounds
	}

	return (float64(passionS) * b.passionWeight) + (float64(expS) * b.experienceWeight) + b.weightedBaseConstant, nil
}

func (b *BatEngine) GroupCandidates(admissionsData []AdmissionCandidate) ([]TeamEvaluationData, []AdmissionCandidate) {
	teamMap := make(map[uuid.UUID][]AdmissionCandidate)
	individualCandidates := make([]AdmissionCandidate, 0)

	// Sort admissions candidates on a valid TeamID
	for _, app := range admissionsData {
		if app.TeamID.Valid {
			teamMap[app.TeamID.UUID] = append(teamMap[app.TeamID.UUID], app)
		} else {
			individualCandidates = append(individualCandidates, app)
		}
	}

	teamsEvalData := make([]TeamEvaluationData, 0)
	for teamId, members := range teamMap {
		var totalScore float64
		for _, member := range members {
			totalScore += member.WeightedScore
		}

		teamsEvalData = append(teamsEvalData, TeamEvaluationData{
			TeamID:               teamId,
			Members:              members,
			AverageWeightedScore: totalScore / float64(len(members)),
		})
	}

	return teamsEvalData, individualCandidates
}

func (b *BatEngine) AcceptIndividuals(idvs []AdmissionCandidate) ([]AdmissionCandidate, []AdmissionCandidate) {
	accepted := make([]AdmissionCandidate, 0)

	pools := groupCandidateByType(idvs)

	buckets := []BucketConfig{
		{
			Type:        BucketTypeUFEarly,
			QuotaPtr:    &b.quota.UF.EarlyLeft,
			RolloverPtr: &b.quota.UF.LateLeft,
		},
		{
			Type:        BucketTypeUFLate,
			QuotaPtr:    &b.quota.UF.LateLeft,
			RolloverPtr: &b.quota.UF.EarlyLeft,
		},
		{
			Type:        BucketTypeOtherEarly,
			QuotaPtr:    &b.quota.Other.EarlyLeft,
			RolloverPtr: &b.quota.Other.LateLeft,
		},
		{
			Type:        BucketTypeOtherLate,
			QuotaPtr:    &b.quota.Other.LateLeft,
			RolloverPtr: &b.quota.Other.EarlyLeft,
		},
	}

	// Two pass is the minimum we need to ensure convergence
	// on the rollover quotas. As of right now we do *not* support
	// rollover between non-leaf nodes/conditions.
	for range 2 {
		for _, bucket := range buckets {
			candidates := pools[bucket.Type]
			if len(candidates) == 0 {
				if bucket.RolloverPtr != nil {
					*bucket.RolloverPtr += *bucket.QuotaPtr
					*bucket.QuotaPtr = 0
				}
				continue
			}

			b.ApplyIndividualSortKey(candidates)
			sort.Slice(candidates, func(i, j int) bool {
				return candidates[i].SortKey > candidates[j].SortKey
			})

			remainingQuota := *bucket.QuotaPtr
			countToAccept := min(int(remainingQuota), len(candidates))

			for i := range countToAccept {
				candidate := candidates[i]
				accepted = append(accepted, candidate)

				*bucket.QuotaPtr--
				b.quota.TotalAccepted++
			}

			// Rollover leftover slots to nearest neighbor condition
			if bucket.QuotaPtr != nil && *bucket.QuotaPtr > 0 {
				*bucket.RolloverPtr += *bucket.QuotaPtr
				*bucket.QuotaPtr = 0
			}

			// Removed accepted candidates (topK) from pool
			pools[bucket.Type] = candidates[countToAccept:]
		}
	}

	rejected := make([]AdmissionCandidate, 0)
	for _, remaining := range pools {
		rejected = append(rejected, remaining...)
	}

	return accepted, rejected
}

func groupCandidateByType(idvs []AdmissionCandidate) map[BucketType][]AdmissionCandidate {
	m := make(map[BucketType][]AdmissionCandidate)
	for _, idv := range idvs {
		t := determineCandidateBucketType(idv)

		m[t] = append(m[t], idv)
	}

	return m
}

func determineCandidateBucketType(idv AdmissionCandidate) BucketType {
	if idv.IsEarlyCareer && idv.IsUFStudent {
		return BucketTypeUFEarly
	} else if !idv.IsEarlyCareer && idv.IsUFStudent {
		return BucketTypeUFLate
	} else if idv.IsEarlyCareer && !idv.IsUFStudent {
		return BucketTypeOtherEarly
	} else {
		return BucketTypeOtherLate
	}
}

func (b *BatEngine) ScoreTeams(teams []TeamEvaluationData) {
	for i := range teams {
		team := &teams[i]

		var totalScore float64
		for _, member := range team.Members {
			totalScore += member.WeightedScore
		}

		team.AverageWeightedScore = totalScore / float64(len(team.Members))
	}
}

func (b *BatEngine) ApplyIndividualSortKey(idv []AdmissionCandidate) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range idv {
		indie := &idv[i]

		indie.SortKey = generateSortKey(r, indie.WeightedScore)
	}
}

func (b *BatEngine) ApplyTeamSortKey(teams []TeamEvaluationData) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range teams {
		team := &teams[i]

		team.SortKey = generateSortKey(r, team.AverageWeightedScore)
	}

	// Sort my descending for top-k results
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].SortKey > teams[j].SortKey
	})
}

func (b *BatEngine) AcceptTeams(teams []TeamEvaluationData) ([]AdmissionCandidate, []AdmissionCandidate) {
	accepted := make([]AdmissionCandidate, 0)
	remaining := make([]AdmissionCandidate, 0)

	b.ScoreTeams(teams)
	b.ApplyTeamSortKey(teams)

	for _, team := range teams {
		if b.quota.TeamSlotsLeft <= int32(len(team.Members)) {
			remaining = append(remaining, team.Members...)
			continue
		}

		requiredQuota := getTeamQuotaRequirement(&team)

		if b.canAcceptTeam(requiredQuota) {
			accepted = append(accepted, team.Members...)
			size := int32(len(team.Members))

			b.quota.TotalAccepted += size
			b.quota.TeamSlotsLeft -= size

			// Adjust primary quota buckets
			b.quota.UF.EarlyLeft -= requiredQuota.UF.EarlyLeft
			b.quota.UF.LateLeft -= requiredQuota.UF.LateLeft
			b.quota.Other.EarlyLeft -= requiredQuota.Other.EarlyLeft
			b.quota.Other.LateLeft -= requiredQuota.Other.LateLeft
		} else {
			remaining = append(remaining, team.Members...)
		}

	}

	return accepted, remaining
}

func getTeamQuotaRequirement(team *TeamEvaluationData) QuotaState {
	reqQuota := QuotaState{
		TotalAccepted: 0,
		TeamSlotsLeft: 0,
		UF: CategoryQuota{
			EarlyLeft: 0,
			LateLeft:  0,
		},
		Other: CategoryQuota{
			EarlyLeft: 0,
			LateLeft:  0,
		},
	}

	for _, member := range team.Members {
		if member.IsUFStudent {
			if member.IsEarlyCareer {
				reqQuota.UF.EarlyLeft += 1
			} else {
				reqQuota.UF.LateLeft += 1
			}
		} else {
			if member.IsEarlyCareer {
				reqQuota.Other.EarlyLeft += 1
			} else {
				reqQuota.Other.LateLeft += 1
			}
		}

		reqQuota.TotalAccepted += 1
		reqQuota.TeamSlotsLeft += 1
	}

	return reqQuota
}

// Checks whether the current quota can fullfill the required
// quota calculated in getTeamQuotaRequirement
func (b BatEngine) canAcceptTeam(req QuotaState) bool {
	return req.UF.EarlyLeft <= b.quota.UF.EarlyLeft &&
		req.UF.LateLeft <= b.quota.UF.LateLeft &&
		req.Other.EarlyLeft <= b.quota.Other.EarlyLeft &&
		req.Other.LateLeft <= b.quota.Other.LateLeft
}

// generateSortKey generates a priority key for weighted random sampling.
//
// This follows the Efraimidis–Spirakis algorithm for weighted sampling without
// replacement. For an item with weight w, a uniform random value v ∈ (0,1) is
// drawn and transformed as v^(1/w). Sorting items by this key and selecting
// the top-k yields a sample where selection probability is proportional to w.
//
// Source: https://www.sciencedirect.com/science/article/abs/pii/S002001900500298X
func generateSortKey(rand *rand.Rand, weight float64) float64 {
	v := rand.Float64()
	exp := 1.0 / (weight)

	return math.Pow(v, exp)
}

// func (b *BatEngine) GroupAndScoreAppli

// equalWithinTolerance checks if two float64 values are equal within a given tolerance.
// It handles exact equality, zero values, and relative differences. We recommend you
// set the tolerance to 1e-9.
//
// TL;DR Compares a and b with a tolerance of e.
func equalWithinTolerance(a, b, e float64) bool {
	if a == b {
		return true
	}

	d := math.Abs(a - b)

	if b == 0 {
		return d < e
	}

	return (d / math.Abs(b)) < e
}
