package beholder

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpellListsDataSource(t *testing.T) {

	Convey("extractClasses", t, func() {

		destMap := map[string][]*ClassSpell{}

		Convey("merges same-named classes", func() {
			spell := Spell{
				Classes: "Crew (frommybunk), Crew (justashepherd)",
			}
			parseSpellUsers(&spell)
			extractClassSpells(destMap, spell)

			So(destMap, ShouldHaveLength, 1)
			So(destMap["Crew"], ShouldHaveLength, 1)

			classSpell := destMap["Crew"][0]
			So(classSpell.Variants, ShouldHaveLength, 2)
			So(classSpell.VariantsOnly, ShouldBeTrue)
		})

		Convey("merges 'variant' and 'any' classes", func() {
			spell := Spell{
				Classes: "Crew (drsimon), Crew",
			}
			parseSpellUsers(&spell)
			extractClassSpells(destMap, spell)

			So(destMap, ShouldHaveLength, 1)
			So(destMap["Crew"], ShouldHaveLength, 1)

			classSpell := destMap["Crew"][0]
			So(classSpell.Variants, ShouldHaveLength, 1)
			So(classSpell.VariantsOnly, ShouldBeFalse)
		})
	})
}
