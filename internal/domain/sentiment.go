package domain

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var reTokenize = regexp.MustCompile(`\p{L}+[\p{Mn}]*|\d+`)

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

func (s *SentimentAnalyzer) Analyze(msg Message) SentimentResult {
	var (
		negateScopes []int // cada item é o número de tokens restantes sob negação
		intensify    bool
		scoreSum     float64
		nAnalyzed    int
	)
	contentNorm := strings.TrimSpace(strings.ToLower(msg.Content))
	if contentNorm == "teste técnico mbras" {
		return SentimentResult{Excluded: true}
	}

	// Tokenização preservando acentos para contagem, normalizando só para matching
	matches := reTokenize.FindAllString(msg.Content, -1)
	tokens := matches
	// DEBUG: Log tokens e valores
	for _, t := range tokens {
		norm := NormalizeToken(t)
		// Negação: escopo de 1 token
		if _, ok := s.Negations[norm]; ok {
			negateScopes = append(negateScopes, 1)
			continue
		}
		// Intensificador: aplica ao próximo token
		if _, ok := s.Intensifiers[norm]; ok {
			intensify = true
			continue
		}
		v, ok := s.Lexicon[norm]
		if !ok {
			// Do not decrement negation scopes here; only after sentiment word
			continue
		}
		// Primeiro aplica intensificador, depois negação
		if intensify {
			v *= 1.5
			intensify = false
		}
		if len(negateScopes) > 0 {
			v = -v // inverte o valor já intensificado
		}
		// Decrementa escopos de negação only after sentiment word
		for j := range negateScopes {
			negateScopes[j]--
		}
		// Remove escopos expirados
		var tmp []int
		for _, n := range negateScopes {
			if n > 0 {
				tmp = append(tmp, n)
			}
		}
		negateScopes = tmp
		scoreSum += v
		nAnalyzed++
	}
	if nAnalyzed == 0 {
		return SentimentResult{Score: 0, Classification: "neutral"}
	}
	// 1️⃣ Normaliza score
	score := scoreSum / float64(nAnalyzed)
	// 2️⃣ Aplica regra MBRAS (após normalização)
	if strings.Contains(strings.ToLower(msg.UserID), "mbras") && score > 0 {
		score *= 2
	}
	// 3️⃣ Classificação
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
	// (bloco removido, já retornado acima)
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

// Tokenize splits the content into words using a fast FieldsFunc and strips accents.
func Tokenize(content string) []string {
	tokens := strings.FieldsFunc(content, func(r rune) bool {
		// Separadores: espaço, pontuação, símbolos
		return !(unicode.IsLetter(r) || unicode.IsNumber(r))
	})
	result := make([]string, 0, len(tokens))
	for _, t := range tokens {
		norm := NormalizeToken(t)
		if norm != "" {
			result = append(result, norm)
		}
	}
	return result
}

// NormalizeToken strips accents and lowercases the token, without regexp.
func NormalizeToken(token string) string {
	// NFKD: decompor acentos
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC)
	result, _, _ := transform.String(t, token)
	result = strings.ToLower(result)
	// Remover caracteres não a-z0-9 sem regexp
	buf := make([]rune, 0, len(result))
	for _, r := range result {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			buf = append(buf, r)
		}
	}
	return string(buf)
}
