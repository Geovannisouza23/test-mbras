package domain

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

type SentimentAnalyzer struct {
	Lexicon      map[string]float64
	Intensifiers map[string]struct{}
	Negations    map[string]struct{}
}

func NewSentimentAnalyzer(
	lexicon map[string]float64,
	intensifiers map[string]struct{},
	negations map[string]struct{},
) *SentimentAnalyzer {
	return &SentimentAnalyzer{
		Lexicon:      lexicon,
		Intensifiers: intensifiers,
		Negations:    negations,
	}
}

type SentimentResult struct {
	Score          float64
	Classification string
	Excluded       bool
}

func (s *SentimentAnalyzer) Analyze(msg Message) SentimentResult {
	contentNorm := strings.TrimSpace(strings.ToLower(msg.Content))

	// META case
	if contentNorm == "teste técnico mbras" {
		return SentimentResult{Excluded: true}
	}

	tokens := Tokenize(msg.Content)
	normTokens := make([]string, len(tokens))
	for i, t := range tokens {
		normTokens[i] = NormalizeToken(t)
	}

	var (
		scoreSum  float64
		nAnalyzed int
		negScopes []int
	)

	for i := 0; i < len(tokens); i++ {
		orig := tokens[i]
		norm := normTokens[i]

		// Ignora hashtags
		if strings.HasPrefix(orig, "#") {
			negScopes = decreaseScopes(negScopes)
			continue
		}

		// Intensificador
		mult := 1.0
		if _, ok := s.Intensifiers[norm]; ok {
			if i+1 < len(normTokens) {
				next := normTokens[i+1]
				if _, ok2 := s.Lexicon[next]; ok2 &&
					!strings.HasPrefix(tokens[i+1], "#") {
					mult = 1.5
				}
			}
			negScopes = decreaseScopes(negScopes)
			continue
		}

		// Negação
		if _, ok := s.Negations[norm]; ok {
			negScopes = append(negScopes, 3)
			continue
		}

		val, ok := s.Lexicon[norm]
		if !ok {
			negScopes = decreaseScopes(negScopes)
			continue
		}

		v := val * mult

		// Aplica negação se houver escopo ativo
		if len(negScopes)%2 == 1 {
			v = -v
		}

		scoreSum += v
		nAnalyzed++

		negScopes = decreaseScopes(negScopes)
	}

	if nAnalyzed == 0 {
		return SentimentResult{
			Score:          0,
			Classification: "neutral",
		}
	}

	score := scoreSum / float64(nAnalyzed)

	// 🔥 Regra MBRAS (antes da classificação)
	if strings.Contains(strings.ToLower(msg.UserID), "mbras") && score > 0 {
		score *= 2
	}

	class := "neutral"
	if score > 0.1 {
		class = "positive"
	} else if score < -0.1 {
		class = "negative"
	}

	return SentimentResult{
		Score:          score,
		Classification: class,
	}
}

func decreaseScopes(scopes []int) []int {
	var updated []int
	for _, v := range scopes {
		if v-1 > 0 {
			updated = append(updated, v-1)
		}
	}
	return updated
}

func Tokenize(content string) []string {
	var tokens []string
	var token strings.Builder

	for _, r := range content {
		if unicode.IsSpace(r) {
			if token.Len() > 0 {
				tokens = append(tokens, token.String())
				token.Reset()
			}
			continue
		}
		token.WriteRune(r)
	}

	if token.Len() > 0 {
		tokens = append(tokens, token.String())
	}

	return tokens
}

func NormalizeToken(token string) string {
	normed := norm.NFKD.String(token)

	var b strings.Builder
	for _, r := range normed {
		// Remove marcas de combinação (acentos)
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}

	return strings.ToLower(b.String())
}
