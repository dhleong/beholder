package beholder

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ShouldMatchMonster(actual any, expected ...any) string {
	query := actual.(string)
	monster := expected[0].(Monster)
	if NewQueryMatcher(query).Match(monster).Matched {
		return ""
	}
	return fmt.Sprintf("`%v` should match `%v`", query, monster)
}

func ShouldNotMatchMonster(actual any, expected ...any) string {
	query := actual.(string)
	monster := expected[0].(Monster)
	if !NewQueryMatcher(query).Match(monster).Matched {
		return ""
	}
	return fmt.Sprintf("`%v` should NOT match `%v`", query, monster)
}

func TestCRQueries(t *testing.T) {
	Convey("ChallengeRatingQuery", t, func() {
		Convey("should match a monster", func() {
			So("cr2", ShouldMatchMonster, Monster{
				Challenge: "2",
			})
		})

		Convey("should cooperate with name filter", func() {
			So("cr2 zom", ShouldMatchMonster, Monster{
				Named:     Named{Name: "Zombie"},
				Challenge: "2",
			})
			So("cr2 zom", ShouldNotMatchMonster, Monster{
				Named:     Named{Name: "Kelpie"},
				Challenge: "2",
			})
		})
	})
}
