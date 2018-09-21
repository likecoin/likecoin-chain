package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsValidBigIntegerString(t *testing.T) {
	Convey("If input is greater than or equal to 0, it should pass", t, func() {
		So(IsValidBigIntegerString("0"), ShouldBeTrue)
		So(IsValidBigIntegerString("1"), ShouldBeTrue)
	})

	Convey("If input is less than 0, it should fail", t, func() {
		So(IsValidBigIntegerString("-1"), ShouldBeFalse)
	})
}
