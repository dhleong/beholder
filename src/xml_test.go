package beholder

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestXmlParser(t *testing.T) {

	Convey("parseSpellUsers", t, func() {
		Convey("should handle a single, simple class", func() {
			spell := &Spell{
				Classes: "Captain",
			}

			parseSpellUsers(spell)

			So(spell.SpellUsers, ShouldHaveLength, 1)
			So(spell.SpellUsers[0].Name, ShouldEqual, "Captain")
			So(spell.SpellUsers[0].Variant, ShouldBeEmpty)
		})

		Convey("should handle multiple simple classes", func() {
			spell := &Spell{
				Classes: "Captain, Pilot, First Mate",
			}

			parseSpellUsers(spell)

			So(spell.SpellUsers, ShouldHaveLength, 3)
			So(spell.SpellUsers[0].Name, ShouldEqual, "Captain")
			So(spell.SpellUsers[1].Name, ShouldEqual, "Pilot")
			So(spell.SpellUsers[2].Name, ShouldEqual, "First Mate")
			So(spell.SpellUsers[0].Variant, ShouldBeEmpty)
			So(spell.SpellUsers[1].Variant, ShouldBeEmpty)
			So(spell.SpellUsers[2].Variant, ShouldBeEmpty)
		})

		Convey("should handle a single variant class", func() {
			spell := &Spell{
				Classes: "Captain (mreynolds)",
			}

			parseSpellUsers(spell)

			So(spell.SpellUsers, ShouldHaveLength, 1)
			So(spell.SpellUsers[0].Name, ShouldEqual, "Captain")
			So(spell.SpellUsers[0].Variant, ShouldEqual, "mreynolds")
		})

		Convey("should handle multiple variant classes", func() {
			spell := &Spell{
				Classes: "Captain (mreynolds), Engineer (itskaylee)",
			}

			parseSpellUsers(spell)

			So(spell.SpellUsers, ShouldHaveLength, 2)
			So(spell.SpellUsers[0].Name, ShouldEqual, "Captain")
			So(spell.SpellUsers[0].Variant, ShouldEqual, "mreynolds")
			So(spell.SpellUsers[1].Name, ShouldEqual, "Engineer")
			So(spell.SpellUsers[1].Variant, ShouldEqual, "itskaylee")
		})

		Convey("should handle a mix of classes", func() {
			spell := &Spell{
				Classes: "Captain, Engineer (itskaylee)",
			}

			parseSpellUsers(spell)

			So(spell.SpellUsers, ShouldHaveLength, 2)
			So(spell.SpellUsers[0].Name, ShouldEqual, "Captain")
			So(spell.SpellUsers[0].Variant, ShouldBeEmpty)
			So(spell.SpellUsers[1].Name, ShouldEqual, "Engineer")
			So(spell.SpellUsers[1].Variant, ShouldEqual, "itskaylee")
		})
	})
}
