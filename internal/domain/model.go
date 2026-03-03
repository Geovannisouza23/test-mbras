package domain

import (
	"errors"
	"time"
)

// ErrNotFound é retornado quando um recurso não é encontrado (mock infra)
var ErrNotFound = errors.New("not found")

// Message represents a single feed message for analysis.
type Message struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Reactions int       `json:"reactions"`
	Shares    int       `json:"shares"`
	Views     int       `json:"views"`
	Hashtags  []string  `json:"hashtags"`
}

// AnalyzeFeedRequest is the input for the /analyze-feed endpoint.
type AnalyzeFeedRequest struct {
	Messages          []Message `json:"messages"`
	TimeWindowMinutes int       `json:"time_window_minutes"`
}

// SentimentResult holds the sentiment analysis for a message.
type SentimentResult struct {
	Score          float64
	Classification string // "positive", "neutral", "negative"
	Excluded       bool   // true if excluded from distribution
}

// FeedSentimentFlags indica casos especiais detectados
type FeedSentimentFlags struct {
	MbrasEmployee      bool `json:"mbras_employee"`
	SpecialPattern     bool `json:"special_pattern"`
	CandidateAwareness bool `json:"candidate_awareness"`
}

// FeedSentimentDistribution percentuais dos sentimentos
type FeedSentimentDistribution struct {
	Positive float64 `json:"positive"`
	Negative float64 `json:"negative"`
	Neutral  float64 `json:"neutral"`
}

// AnalyzeFeedResponse é a resposta do endpoint, alinhada ao modelo test-mbras
type AnalyzeFeedResponse struct {
	SentimentDistribution FeedSentimentDistribution `json:"sentiment_distribution"`
	EngagementScore       float64                   `json:"engagement_score"`
	TrendingTopics        []string                  `json:"trending_topics"`
	InfluenceRanking      map[string]float64        `json:"influence_ranking,omitempty"`
	AnomalyDetected       bool                      `json:"anomaly_detected"`
	AnomalyType           string                    `json:"anomaly_type,omitempty"`
	Flags                 FeedSentimentFlags        `json:"flags"`
	ProcessingTimeMs      float64                   `json:"processing_time_ms"`
}
