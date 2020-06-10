// SPDX-License-Identifier: MIT
// SPDX-License-Identifier: Unlicense

package processor

import "testing"

func TestRankResultsTFIDF(t *testing.T) {
	ml1 := map[string][][]int{}
	ml1["this"] = [][]int{{1}}
	ml1["is"] = [][]int{{1}}
	ml1["a"] = [][]int{{1}, {2}}
	ml1["sample"] = [][]int{{1}}

	ml2 := map[string][][]int{}
	ml2["this"] = [][]int{{1}}
	ml2["is"] = [][]int{{1}}
	ml2["another"] = [][]int{{1}, {2}}
	ml2["example"] = [][]int{{1}, {2}, {3}}

	s := []*fileJob{
		{
			MatchLocations: ml1,
			Location:       "/test/other.go",
			Bytes:          12,
		},
		{
			MatchLocations: ml2,
			Location:       "/test/test.go",
			Bytes:          12,
		},
	}

	s = rankResultsTFIDF(2, s)

	if s[0].Score > s[1].Score {
		t.Error("index 0 should have lower score than 1")
	}
}

func TestRankResultsLocation(t *testing.T) {
	ml := map[string][][]int{}
	ml["test"] = [][]int{{1}, {2}, {3}}

	s := []*fileJob{
		{
			MatchLocations: ml,
			Location:       "/test/other.go",
		},
		{
			MatchLocations: ml,
			Location:       "/test/test.go",
		},
	}

	s = rankResultsLocation(s)

	if s[0].Score > s[1].Score {
		t.Error("index 0 should have lower score than 1")
	}
}

func TestCalculateDocumentFrequency(t *testing.T) {
	ml := map[string][][]int{}
	ml["test"] = [][]int{{1}, {2}, {3}}

	s := []*fileJob{
		{
			MatchLocations: ml,
		},
		{
			MatchLocations: ml,
		},
	}

	freq := calculateDocumentFrequency(s)

	if len(freq) != 1 || freq["test"] != 6 {
		t.Error("did not work as expected")
	}
}

func TestSortResults(t *testing.T) {
	s := []*fileJob{
		{
			Filename: "1",
			Location: "",
			Score:    10,
		},
		{
			Filename: "2",
			Location: "",
			Score:    20,
		},
	}
	sortResults(s)

	if s[0].Filename != "2" {
		t.Error("expected 2 first")
	}
}

func TestSortResultsEqualScore(t *testing.T) {
	s := []*fileJob{
		{
			Filename: "1",
			Location: "2",
			Score:    10,
		},
		{
			Filename: "2",
			Location: "1",
			Score:    10,
		},
	}
	sortResults(s)

	if s[0].Filename != "2" {
		t.Error("expected 2 first")
	}
}
