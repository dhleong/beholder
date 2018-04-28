package beholder

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func queryMatches(query, value string) bool {
	return NewQueryMatcher(query).Match(value).Matched
}

func ShouldMatch(actual interface{}, expected ...interface{}) string {
	if queryMatches(actual.(string), expected[0].(string)) {
		return ""
	}

	return fmt.Sprintf("`%v` should match `%v`", actual, expected[0])
}

func ShouldNotMatch(actual interface{}, expected ...interface{}) string {
	if !queryMatches(actual.(string), expected[0].(string)) {
		return ""
	}

	return fmt.Sprintf("`%v` should NOT match `%v`", actual, expected[0])
}

type comparison struct {
	label string
	fn    func(a, b float32) bool
}

func ShouldScore(actual interface{}, expected ...interface{}) string {
	qm := NewQueryMatcher(actual.(string))
	a := expected[0].(string)
	compare := expected[1].(*comparison)

	bs := expected[2:]
	bMatches := make([]MatchResult, 0, len(bs))

	matchA := qm.Match(a)

	if !matchA.Matched {
		return fmt.Sprintf("Expected Matcher('%s') to match: '%v'", actual, a)
	}

	for _, b := range bs {
		m := qm.Match(b.(string))
		if !m.Matched {
			return fmt.Sprintf("Expected Matcher('%s') to match: '%v'", actual, b)
		}

		bMatches = append(bMatches, m)
	}

	scoreA := matchA.Score

	for i, matchB := range bMatches {
		scoreB := matchB.Score
		b := bs[i]

		if compare.fn(scoreA, scoreB) {
			continue
		}

		return fmt.Sprintf(
			"Expected Matcher('%s') to score:\n  '%v' (%f)\n%v:\n  '%v' (%f)",
			actual,
			a, scoreA,
			compare.label,
			b, scoreB,
		)
	}

	return ""
}

var greaterThan = &comparison{
	label: "greater than",
	fn: func(a, b float32) bool {
		return a > b
	},
}

func TestQueryMatcher(t *testing.T) {

	Convey("QueryMatcher.Matches()", t, func() {
		Convey("should be case insensitive", func() {
			So("m", ShouldMatch, "Mal Reynolds")
			So("r", ShouldMatch, "Mal Reynolds")
		})

		Convey("should be fuzzy", func() {
			So("mr", ShouldMatch, "Mal Reynolds")
		})

		Convey("should score by longer subsequences", func() {
			So("mr", ShouldScore,
				"Mr. Reynolds",
				greaterThan,
				"Mary Jeynolds",
			)
		})

		Convey("should prefer matching at start", func() {
			So("wil", ShouldScore,
				"Will-o'-Wisp",
				greaterThan,
				"Wind Walk",
				"Psionic Weapon (Immortal)",
			)
		})

		Convey("should score word starts", func() {
			So("wow", ShouldScore,
				"Will-o'-wisp",
				greaterThan,
				"Row your wand",
			)
		})
	})

}
