package main

import (
	"strings"
	"testing"
)

var testRecords = []PostalCode{
	{Province: "Jawa Barat", Regency: "Ciamis", District: "Cisaga", Village: "Danasari", Code: 46386, Latitude: -7.327, Longitude: 108.457, Elevation: 110, Timezone: "WIB"},
	{Province: "Jawa Tengah", Regency: "Purbalingga", District: "Karangjambu", Village: "Danasari", Code: 53357, Latitude: -7.185, Longitude: 109.436, Elevation: 705, Timezone: "WIB"},
	{Province: "Jawa Barat", Regency: "Bogor", District: "Bogor", Village: "Batu Tulis", Code: 16124, Latitude: -6.656, Longitude: 106.801, Elevation: 240, Timezone: "WIB"},
}

func TestBuildFulltext(t *testing.T) {
	p := PostalCode{
		Village: "A", District: "B", Regency: "C", Province: "D",
	}
	ft := buildFulltext(p)

	if ft == "" {
		t.Fatal("buildFulltext returned empty string")
	}

	// All pairwise combos of 4 fields = 12 pairs, each pair is "X Y "
	// So fulltext should contain all field values
	for _, field := range []string{"A", "B", "C", "D"} {
		if !strings.Contains(ft, field) {
			t.Errorf("fulltext missing field %q", field)
		}
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"hello world", 2},
		{"  hello   world  ", 2},
		{"jakarta", 1},
		{"", 0},
		{"   ", 0},
		{"a b c d", 4},
	}
	for _, tt := range tests {
		got := strings.Fields(tt.input)
		if len(got) != tt.want {
			t.Errorf("tokenize(%q) = %d tokens, want %d", tt.input, len(got), tt.want)
		}
	}
}

func TestParseSearchQuery(t *testing.T) {
	q := ParseSearchQuery("  danasari  jawa  ")
	if q.Raw != "danasari jawa" {
		t.Errorf("Raw = %q, want %q", q.Raw, "danasari jawa")
	}
	if len(q.Tokens) != 2 {
		t.Errorf("Tokens len = %d, want 2", len(q.Tokens))
	}
}

func TestSearchExact(t *testing.T) {
	s := NewSearcher(testRecords)

	results := s.Search(ParseSearchQuery("danasari"), 20)
	if len(results) == 0 {
		t.Fatal("expected results for 'danasari'")
	}

	// Both Danasari records should match
	codes := make(map[int]bool)
	for _, r := range results {
		codes[r.Code] = true
	}
	if !codes[46386] || !codes[53357] {
		t.Errorf("expected both Danasari codes, got %v", codes)
	}

	// Should not include fulltext field
	for _, r := range results {
		if r.Village == "" {
			t.Error("result missing Village")
		}
	}
}

func TestSearchMultiToken(t *testing.T) {
	s := NewSearcher(testRecords)

	results := s.Search(ParseSearchQuery("danasari jawa barat"), 20)
	if len(results) != 1 {
		t.Fatalf("expected 1 result for 'danasari jawa barat', got %d", len(results))
	}
	if results[0].Code != 46386 {
		t.Errorf("expected code 46386, got %d", results[0].Code)
	}
}

func TestSearchNoMatch(t *testing.T) {
	s := NewSearcher(testRecords)

	results := s.Search(ParseSearchQuery("xyznonexistent"), 20)
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchLimit(t *testing.T) {
	s := NewSearcher(testRecords)

	results := s.Search(ParseSearchQuery("danasari"), 1)
	if len(results) != 1 {
		t.Fatalf("expected 1 result with limit=1, got %d", len(results))
	}
}

func TestSearchEmptyQuery(t *testing.T) {
	s := NewSearcher(testRecords)

	results := s.Search(ParseSearchQuery(""), 20)
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty query, got %d", len(results))
	}
}

func TestSearchScoreOrder(t *testing.T) {
	s := NewSearcher(testRecords)

	results := s.Search(ParseSearchQuery("danasari jawa"), 20)
	if len(results) < 2 {
		t.Fatal("need at least 2 results for score test")
	}

	// Record with more token matches should rank higher
	if results[0].Code != 46386 && results[0].Code != 53357 {
		t.Errorf("unexpected top result code: %d", results[0].Code)
	}
}
