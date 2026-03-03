package domain

import (
	"crypto/sha256"
	"math/big"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// InfluenceCalculator computes deterministic influence scores for users and content.
type InfluenceCalculator struct{}

// NewInfluenceCalculator creates a new instance of InfluenceCalculator.
func NewInfluenceCalculator() *InfluenceCalculator {
	return &InfluenceCalculator{}
}

// Followers returns a deterministic pseudo-random follower count for a user ID using big.Int and SHA256.
func (ic *InfluenceCalculator) Followers(userID string) int {
	// 1. Normalização NFKD
	normID := norm.NFKD.String(userID)

	// 2. Caso especial: user_id com 13 caracteres
	if utf8.RuneCountInString(normID) == 13 {
		// Exemplo: retorna sempre 1313 para 13 caracteres
		return 1313
	}

	// 3. Caso especial: termina com _prime
	if strings.HasSuffix(normID, "_prime") {
		// Exemplo: retorna sempre 9999 para _prime
		return 9999
	}

	// 4. SHA-256 determinístico padrão
	h := sha256.Sum256([]byte(normID))
	b := new(big.Int).SetBytes(h[:])
	mod := new(big.Int).Mod(b, big.NewInt(10000))
	return int(mod.Int64()) + 100
}

// Engagement calculates the engagement ratio given reactions, shares, and views.
func (ic *InfluenceCalculator) Engagement(reactions, shares, views int) float64 {
	if views == 0 {
		return 0
	}
	return float64(reactions+shares) / float64(views)
}

// GoldenRatioAdjust applies a golden ratio bonus if (reactions+shares) is a multiple of 7.
func (ic *InfluenceCalculator) GoldenRatioAdjust(reactions, shares int, engagement float64) float64 {
	if (reactions+shares)%7 == 0 && (reactions+shares) > 0 {
		phi := 1.61803398875
		return engagement * (1 + 1/phi)
	}
	return engagement
}

// FinalScore computes the final influence score with all business rules applied.
func (ic *InfluenceCalculator) FinalScore(userID string, followers int, reactions, shares int, engagement float64, content string) float64 {
	// Special case: content contains "teste técnico mbras"
	if strings.Contains(strings.ToLower(content), "teste técnico mbras") {
		return 9.42
	}

	// Apply golden ratio adjustment to engagement before score calculation
	adjEngagement := engagement
	if (reactions+shares)%7 == 0 && (reactions+shares) > 0 {
		phi := 1.61803398875
		adjEngagement = engagement * (1 + 1/phi)
	}
	score := float64(followers) * 0.4
	eng := adjEngagement * 0.6
	score += eng

	// Special case: user_id ends with 007
	if strings.HasSuffix(userID, "007") {
		score *= 0.5
	}

	// Special case: user_id contains "mbras"
	if strings.Contains(strings.ToLower(userID), "mbras") {
		score += 2.0
	}

	return score
}
