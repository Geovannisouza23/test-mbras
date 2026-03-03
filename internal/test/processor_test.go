package test

import (
	"testing"
	"time"

	domain "backend-challenge-092025/internal/domain"
)

func TestAnalyzeFeed_ValidRequest(t *testing.T) {
	now := time.Date(2026, 3, 2, 22, 0, 0, 0, time.UTC)
	p := domain.NewProcessor()
	req := domain.AnalyzeFeedRequest{
		Messages: []domain.Message{
			{
				UserID:    "user_mbras1",
				Content:   "ótimo trabalho!",
				Timestamp: now.Add(-2 * time.Minute),
				Reactions: 7,
				Shares:    0,
				Views:     10,
				Hashtags:  []string{"#golang"},
			},
			{
				UserID:    "user_test2",
				Content:   "não gostei",
				Timestamp: now.Add(-3 * time.Minute),
				Reactions: 0,
				Shares:    0,
				Views:     5,
				Hashtags:  []string{"#backend"},
			},
			{
				UserID:    "user_test3",
				Content:   "péssimo",
				Timestamp: now.Add(-4 * time.Minute),
				Reactions: 0,
				Shares:    0,
				Views:     5,
				Hashtags:  []string{"#fail"},
			},
		},
		TimeWindowMinutes: 10,
	}
	resp, code, msg := p.AnalyzeFeed(req, now)
	if code != 200 {
		t.Fatalf("expected 200, got %d, msg: %s", code, msg)
	}
	if resp.EngagementScore <= 0 {
		t.Errorf("expected positive engagement score, got %f", resp.EngagementScore)
	}
	if len(resp.TrendingTopics) == 0 {
		t.Errorf("expected trending topics, got none")
	}
	// Agora a distribuição deve ser em percentuais
	if got := resp.SentimentDistribution.Positive; int(got+0.5) != 33 {
		t.Errorf("expected ~33%% positive, got %f", got)
	}
	if got := resp.SentimentDistribution.Negative; int(got+0.5) != 67 {
		t.Errorf("expected ~67%% negative, got %f", got)
	}
	t.Logf("SentimentDistribution: %+v", resp.SentimentDistribution)
}
