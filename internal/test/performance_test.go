package test

import (
	"backend-challenge-092025/internal/domain"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func TestPerformance_AnalyzeFeed(t *testing.T) {
	now := time.Now().UTC()
	p := domain.NewProcessor()

	// Teste de tempo para 1000 mensagens
	var msgs1k []domain.Message
	for i := 0; i < 1000; i++ {
		ts := now.Add(-time.Duration(rand.Intn(60)) * time.Second).UTC()
		msgs1k = append(msgs1k, domain.Message{
			UserID:    "user_perf",
			Content:   "ótimo trabalho! #performance",
			Timestamp: ts,
			Reactions: rand.Intn(10),
			Shares:    rand.Intn(5),
			Views:     rand.Intn(100),
			Hashtags:  []string{"#performance"},
		})
	}
	start := time.Now()
	resp, code, _ := p.AnalyzeFeed(domain.AnalyzeFeedRequest{
		Messages:          msgs1k,
		TimeWindowMinutes: 60,
	}, now)
	dur := time.Since(start)
	if code != 200 {
		t.Fatalf("expected 200, got %d", code)
	}
	if len(resp.TrendingTopics) == 0 {
		t.Errorf("expected trending topics, got none")
	}
	if dur > 200*time.Millisecond {
		t.Errorf("performance for 1000 msgs: %v (alvo < 200ms)", dur)
	}

	// Teste de memória para 10k mensagens
	runtime.GC()
	var m2 runtime.MemStats
	var msgs10k []domain.Message
	for i := 0; i < 10000; i++ {
		ts := now.Add(-time.Duration(rand.Intn(60)) * time.Second).UTC()
		msgs10k = append(msgs10k, domain.Message{
			UserID:    "user_perf",
			Content:   "ótimo trabalho! #performance",
			Timestamp: ts,
			Reactions: rand.Intn(10),
			Shares:    rand.Intn(5),
			Views:     rand.Intn(100),
			Hashtags:  []string{"#performance"},
		})
	}
	p.AnalyzeFeed(domain.AnalyzeFeedRequest{
		Messages:          msgs10k,
		TimeWindowMinutes: 60,
	}, now)
	runtime.GC()
	runtime.ReadMemStats(&m2)
	usedMB := float64(m2.Alloc) / 1024.0 / 1024.0
	if usedMB > 20.0 {
		t.Errorf("memória usada para 10k mensagens: %.2fMB (alvo ≤ 20MB)", usedMB)
	}
}
