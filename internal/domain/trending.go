package domain

import (
	"math"
	"sort"
	"strings"
	"time"
)

type hashtagStats struct {
	Tag           string
	Weight        float64
	Frequency     int
	SentimentSum  float64
	SentimentType string
}

// TrendingTopics analyzes hashtags for trending topics.
// TrendingTopics analyzes hashtags for trending topics.
type TrendingTopics struct{}

// NewTrendingTopics creates a new TrendingTopics analyzer.
func NewTrendingTopics() *TrendingTopics {
	return &TrendingTopics{}
}

// TopHashtags returns the top trending hashtags based on messages and sentiment results.
func (tt *TrendingTopics) TopHashtags(messages []Message, sentiments []SentimentResult, now time.Time) []string {

	tagStats := make(map[string]*hashtagStats, len(messages)*2)
	log10Cache := make(map[int]float64)
	for i, msg := range messages {
		for _, tag := range msg.Hashtags {
			if !strings.HasPrefix(tag, "#") || len(tag) < 2 {
				continue
			}
			minutes := now.Sub(msg.Timestamp).Minutes()
			if minutes < 0.01 {
				minutes = 0.01
			}
			w := 1 + (1.0 / minutes)
			// Sentiment modifier
			mod := 1.0
			stype := sentiments[i].Classification
			switch stype {
			case "positive":
				mod = 1.2
			case "negative":
				mod = 0.8
			}
			w *= mod
			// Long hashtag adjustment
			l := len([]rune(tag))
			if l > 8 {
				logL, ok := log10Cache[l]
				if !ok {
					logL = math.Log10(float64(l)) / math.Log10(8)
					log10Cache[l] = logL
				}
				w *= logL
			}
			stat, ok := tagStats[tag]
			if !ok {
				stat = &hashtagStats{Tag: tag}
				tagStats[tag] = stat
			}
			stat.Weight += w
			stat.Frequency++
			stat.SentimentSum += mod
			stat.SentimentType = stype
		}
	}
	stats := make([]*hashtagStats, 0, len(tagStats))
	for _, s := range tagStats {
		stats = append(stats, s)
	}
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Weight != stats[j].Weight {
			return stats[i].Weight > stats[j].Weight
		}
		if stats[i].Frequency != stats[j].Frequency {
			return stats[i].Frequency > stats[j].Frequency
		}
		if stats[i].SentimentSum != stats[j].SentimentSum {
			return stats[i].SentimentSum > stats[j].SentimentSum
		}
		return stats[i].Tag < stats[j].Tag
	})
	result := make([]string, 0, 5)
	for i := 0; i < len(stats) && i < 5; i++ {
		result = append(result, stats[i].Tag)
	}
	return result
}
