package beholder

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func queryMatches(query, value string) bool {
	return NewQueryMatcher(query).Matches(value)
}

func TestQueryMatcher(t *testing.T) {

	Convey("QueryMatcher.Matches()", t, func() {
		Convey("should be case insensitive", func() {
			So(queryMatches("m", "Mal Reynolds"), ShouldBeTrue)
			So(queryMatches("r", "Mal Reynolds"), ShouldBeTrue)
		})
	})

}
