package helpers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldConvertModifiedAtToTimestamp(t *testing.T) {
	Convey("", t, func() {
		strDate := "2018-07-12T13:16:28.485Z"
		timestamp, err := parseStringToTime(strDate)
		So(err, ShouldBeNil)
		So(timestamp, ShouldBeGreaterThan, 0)
	})
}
