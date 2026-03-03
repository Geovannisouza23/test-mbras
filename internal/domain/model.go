package domain

import (
	"time"
)

// Message represents a single feed message for analysis.
type Message struct {
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

// AnalyzeFeedResponse is the output structure for the endpoint.
type AnalyzeFeedResponse struct {
	SentimentDistribution map[string]float64 `json:"sentiment_distribution"`
	EngagementScore       float64            `json:"engagement_score"`
	TrendingTopics        []string           `json:"trending_topics"`
	AnomaliesDetected     []string           `json:"anomalies_detected"`
}
