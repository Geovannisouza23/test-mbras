package test

import (
	"testing"
	"backend-challenge-092025/internal/domain"
	"time"
)

func TestSentimentSpecificCases(t *testing.T) {
	processor := domain.NewProcessor()
	now := time.Now().UTC()

	cases := []struct {
		name     string
		userID   string
		content  string
		expected string
	}{
		{
			name:     "Não muito bom (usuário normal)",
			userID:   "user_abc",
			content:  "Não muito bom! #produto",
			expected: "negative",
		},
		{
			name:     "Super adorei! (user_mbras_123)",
			userID:   "user_mbras_123",
			content:  "Super adorei!",
			expected: "positive",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			msg := domain.Message{
				ID:        "msg1",
				UserID:    c.userID,
				Content:   c.content,
				Timestamp: now,
				Reactions:  1,
				Shares:     0,
				Views:      1,
				Hashtags:   []string{"#produto"},
			}
			analyzer := domain.NewSentimentAnalyzer(processor.Lexicon, processor.Intensifiers, processor.Negations)
			result := analyzer.Analyze(msg)
			if result.Classification != c.expected {
				t.Errorf("%s: esperado %s, obtido %s", c.name, c.expected, result.Classification)
			}
		})
	}
}
