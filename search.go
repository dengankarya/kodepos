package main

import (
	"sort"
	"strconv"
	"strings"
)

// Searcher provides token-based fulltext search over postal code records.
type Searcher struct {
	data          []PostalCode
	fulltextLower []string // lowercase, for scoring
	index         map[string][]int
}

// NewSearcher builds the fulltext index from the given records.
func NewSearcher(data []PostalCode) *Searcher {
	s := &Searcher{
		data:          data,
		fulltextLower: make([]string, len(data)),
		index:         make(map[string][]int, len(data)*4),
	}

	for i, record := range data {
		ft := buildFulltext(record)
		ftLower := strings.ToLower(ft)
		s.fulltextLower[i] = ftLower

		for _, tok := range strings.Fields(ftLower) {
			s.index[tok] = append(s.index[tok], i)
		}
	}

	return s
}

// buildFulltext generates all pairwise combinations of the record's fields.
func buildFulltext(p PostalCode) string {
	fields := []string{
		p.Village,
		p.District,
		p.Regency,
		p.Province,
		strconv.Itoa(p.Code),
	}

	var b strings.Builder
	for i, a := range fields {
		for j, b_ := range fields {
			if i != j {
				b.WriteString(a)
				b.WriteByte(' ')
				b.WriteString(b_)
				b.WriteByte(' ')
			}
		}
	}
	return b.String()
}

// SearchQuery holds the parsed search parameters.
type SearchQuery struct {
	Tokens []string // lowercase
	Raw    string
}

// ParseSearchQuery normalizes and tokenizes a raw search string.
func ParseSearchQuery(raw string) SearchQuery {
	cleaned := strings.ReplaceAll(raw, "  ", " ")
	cleaned = strings.TrimSpace(cleaned)

	tokens := strings.Fields(cleaned)
	lower := make([]string, len(tokens))
	for i, t := range tokens {
		lower[i] = strings.ToLower(t)
	}

	return SearchQuery{
		Tokens: lower,
		Raw:    cleaned,
	}
}

// Search finds postal codes matching all query tokens. Returns up to limit results.
func (s *Searcher) Search(query SearchQuery, limit int) []PostalCodeResult {
	if len(query.Tokens) == 0 {
		return nil
	}

	n := len(s.data)

	// Track which indices match ALL tokens so far.
	present := make([]bool, n)

	// Seed with first token's matches.
	for _, idx := range s.index[query.Tokens[0]] {
		present[idx] = true
	}

	// Intersect with remaining tokens.
	for _, tok := range query.Tokens[1:] {
		// Build set of matches for this token.
		matches := make([]bool, n)
		for _, idx := range s.index[tok] {
			matches[idx] = true
		}
		// Intersect.
		for i := range present {
			present[i] = present[i] && matches[i]
		}
	}

	// Score candidates.
	type scored struct {
		idx   int
		score float64
	}
	results := make([]scored, 0, 64)
	for i, ok := range present {
		if !ok {
			continue
		}
		ft := s.fulltextLower[i]
		matchCount := 0
		for _, tok := range query.Tokens {
			if strings.Contains(ft, tok) {
				matchCount++
			}
		}
		results = append(results, scored{i, float64(matchCount) / float64(len(query.Tokens))})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	if limit > len(results) {
		limit = len(results)
	}

	out := make([]PostalCodeResult, limit)
	for i, r := range results[:limit] {
		p := s.data[r.idx]
		out[i] = PostalCodeResult{
			Province:  p.Province,
			Regency:   p.Regency,
			District:  p.District,
			Village:   p.Village,
			Code:      p.Code,
			Latitude:  p.Latitude,
			Longitude: p.Longitude,
			Elevation: p.Elevation,
			Timezone:  p.Timezone,
		}
	}
	return out
}
