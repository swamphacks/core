package bat

import (
	"errors"
	"math"
)

var (
	ErrImproperWeights  = errors.New("Passion and experience weights don't add to 1.0")
	ErrScoreOutOfBounds = errors.New("The score can only range from 1 to 5")
)

type BatEngine struct {
	passionWeight        float64
	experienceWeight     float64
	weightedBaseConstant float64
}

func NewBatEngine(passionW, experienceW float64) (*BatEngine, error) {
	if !equalWithinTolerance(passionW+experienceW, 1.0, 1e-9) {
		return nil, ErrImproperWeights
	}

	return &BatEngine{
		passionWeight:        passionW,
		experienceWeight:     experienceW,
		weightedBaseConstant: 0.1,
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
