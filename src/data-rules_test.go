package beholder

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ShouldContainExactly(actual interface{}, expected ...interface{}) string {
	actualSlice := actual.([]string)
	var mismatch string
	if len(actualSlice) == len(expected) {
		for i, actualString := range actualSlice {
			if actualString != expected[i].(string) {
				mismatch = fmt.Sprintf(
					"\nAt [%d]:\n  expected: %s\n    actual: %s",
					i,
					expected[i],
					actualString,
				)
				break
			}
		}

		if mismatch == "" {
			return ""
		}
	}

	return fmt.Sprintf(
		"`%v` should contain exactly `%v`%s",
		actual,
		expected,
		mismatch,
	)
}

func TestGenerateEntities(t *testing.T) {
	Convey("generateEntities", t, func() {

		topRule, entities := generateEntities(
			RuleEntity,
			rule("Rule",
				"Rule introduction",
				section("Sub-section",
					"Subsection text",
				),
				section("Another subsection",
					"More subsection text",
				),
			),
			true,
		)

		Convey("should merge sections and create headers", func() {
			So(
				topRule.GetText(),

				ShouldContainExactly,
				"Rule introduction",
				"",
				"<h2>Sub-section</h2>",
				"Subsection text",
				"",
				"<h2>Another subsection</h2>",
				"More subsection text",
			)
		})

		Convey("should create sub-section entities", func() {
			So(
				entities[1].(Textual).GetText(),

				ShouldContainExactly,
				"Subsection text",
			)
		})

		Convey("Should create the expected number of Entities", func() {
			So(entities, ShouldHaveLength, 3)
		})
	})
}
