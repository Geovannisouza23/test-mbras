package domain

import (
	"crypto/sha256"
	"encoding/binary"
	"strings"
)

// InfluenceCalculator computes deterministic influence scores.
type InfluenceCalculator struct{}

func NewInfluenceCalculator() *InfluenceCalculator {
	return &InfluenceCalculator{}
}

func (ic *InfluenceCalculator) Followers(userID string) int {
	h := sha256.Sum256([]byte(userID))
	v := binary.BigEndian.Uint32(h[:4])
	return int(v%10000) + 100
}

func (ic *InfluenceCalculator) Engagement(reactions, shares, views int) float64 {
	if views == 0 {
		return 0
	}
	return float64(reactions+shares) / float64(views)
}

func (ic *InfluenceCalculator) GoldenRatioAdjust(reactions, shares int, engagement float64) float64 {
	if (reactions+shares)%7 == 0 && (reactions+shares) > 0 {
		phi := 1.61803398875
		return engagement * (1 + 1/phi)
	}
	return engagement
}

func (ic *InfluenceCalculator) FinalScore(userID string, followers int, engagement float64, content string) float64 {
	if strings.Contains(strings.ToLower(content), "teste técnico mbras") {
		return 9.42
	}
	score := float64(followers)*0.4 + engagement*0.6
	if strings.HasSuffix(userID, "007") {
		score *= 0.5
	}
	if strings.Contains(strings.ToLower(userID), "mbras") {
		score += 2.0
	}
	return score
}
