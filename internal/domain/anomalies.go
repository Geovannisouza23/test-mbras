package domain

import (
	"sort"
	"time"
)

// AnomalyDetector detects bursts, alternations, and syncs.
type AnomalyDetector struct{}

func NewAnomalyDetector() *AnomalyDetector {
	return &AnomalyDetector{}
}

func (ad *AnomalyDetector) Detect(messages []Message, sentiments []SentimentResult) []string {
	var anomalies []string
	if burst := detectBurst(messages); burst {
		anomalies = append(anomalies, "burst")
	}
	if alt := detectAlternating(sentiments); alt {
		anomalies = append(anomalies, "alternating_sentiment")
	}
	if sync := detectSynchronized(messages); sync {
		anomalies = append(anomalies, "synchronized_posting")
	}
	return anomalies
}

func detectBurst(messages []Message) bool {
	userTimes := map[string][]time.Time{}
	for _, m := range messages {
		userTimes[m.UserID] = append(userTimes[m.UserID], m.Timestamp)
	}
	for _, times := range userTimes {
		if len(times) <= 10 {
			continue
		}
		sort.Slice(times, func(i, j int) bool { return times[i].Before(times[j]) })
		for i := 0; i <= len(times)-11; i++ {
			if times[i+10].Sub(times[i]) <= 5*time.Minute {
				return true
			}
		}
	}
	return false
}

func detectAlternating(sentiments []SentimentResult) bool {
	if len(sentiments) < 10 {
		return false
	}
	pattern := []string{"positive", "negative"}
	for i := 0; i <= len(sentiments)-10; i++ {
		ok := true
		for j := 0; j < 10; j++ {
			if sentiments[i+j].Classification != pattern[j%2] {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}

func detectSynchronized(messages []Message) bool {
	if len(messages) < 3 {
		return false
	}
	times := make([]time.Time, len(messages))
	for i, m := range messages {
		times[i] = m.Timestamp
	}
	sort.Slice(times, func(i, j int) bool { return times[i].Before(times[j]) })
	for i := 0; i <= len(times)-3; i++ {
		if times[i+2].Sub(times[i]) <= 2*time.Second {
			return true
		}
	}
	return false
}
