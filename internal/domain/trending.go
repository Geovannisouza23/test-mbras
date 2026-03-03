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
type TrendingTopics struct{}

func NewTrendingTopics() *TrendingTopics {
	return &TrendingTopics{}
}

func (tt *TrendingTopics) TopHashtags(messages []Message, sentiments []SentimentResult, now time.Time) []string {

	tagStats := map[string]*hashtagStats{}
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
			if l := len([]rune(tag)); l > 8 {
				w *= math.Log10(float64(l)) / math.Log10(8)
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
	// Sort
	var stats []*hashtagStats
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
	var result []string
	for i := 0; i < len(stats) && i < 5; i++ {
		result = append(result, stats[i].Tag)
	}
	return result
}
