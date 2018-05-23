package beholder

import (
	"bufio"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func parse(s string) []Entity {
	parsed, err := parseRulesJSON(
		bufio.NewReader(
			strings.NewReader(s),
		),
	)
	So(err, ShouldBeNil)
	return parsed
}

func TestParseRulesJSON(t *testing.T) {
	Convey("parseRulesJSON", t, func() {
		Convey("should read sections in order", func() {
			es := parse(`
			{
				"First Section": {
					"content": "Single line"
				},
				"Second Section": {
					"content": []
				}
			}
			`)

			So(es, ShouldHaveLength, 2)
			So(es[0].GetName(), ShouldEqual, "First Section")
			So(es[1].GetName(), ShouldEqual, "Second Section")
		})

		Convey("should handle arrays in content", func() {
			es := parse(`
			{
				"First Section": {
					"content": [
						"Single line",
						[
							"Multi",
							"Line"
						]
					]
				}
			}
			`)

			So(es, ShouldHaveLength, 1)
			So(es[0].GetName(), ShouldEqual, "First Section")

			// TODO verify content
		})

		Convey("should handle tables in content arrays", func() {
			es := parse(`
			{
				"First Section": {
					"content": [
						{
							"table": {
							}
						},
						"Text"
					]
				}
			}
			`)

			So(es, ShouldHaveLength, 1)
			So(es[0].GetName(), ShouldEqual, "First Section")

			textual, ok := es[0].(Textual)
			So(ok, ShouldBeTrue)

			So(textual.GetText(), ShouldContain, "Text")
		})

		Convey("should not include ignored sections", func() {
			es := parse(`
			{
				"Feats": {
					"content": "Ignored"
				}
			}
			`)

			So(es, ShouldHaveLength, 0)
		})
	})
}
