package helpers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldCheckString(t *testing.T) {
	Convey("should check string convertions", t, func() {
		Convey("should check if a string is a integer", func() {
			So(IsInt("1"), ShouldBeTrue)
			So(IsInt("-1"), ShouldBeTrue)
			So(IsInt("1.5"), ShouldBeFalse)
			So(IsInt("1a"), ShouldBeFalse)
		})

		Convey("should check if a string is a float", func() {
			So(IsFloat("1.0"), ShouldBeTrue)
			So(IsFloat("10"), ShouldBeTrue)
			So(IsFloat("-1.0"), ShouldBeTrue)
			So(IsFloat("1a"), ShouldBeFalse)
		})
	})
}
