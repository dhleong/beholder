package tui

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTextFormatter(t *testing.T) {
	Convey("formatText", t, func() {
		Convey("should handle <b> tags anywhere", func() {
			So(
				formatText([]string{
					"<h1>Serenity</h1>",
					"A <b>Firefly</b>-class ship",
				}), ShouldEqual,
				"[::b]Serenity[::-]\nA [::b]Firefly[::-]-class ship",
			)
		})
	})
}
