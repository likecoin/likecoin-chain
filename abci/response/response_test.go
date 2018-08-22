package response

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMerge(t *testing.T) {
	r1 := R{
		Info: "Testing",
	}
	r2 := R{
		Code: 1,
		Data: []byte("{\"x\":1}"),
	}

	Convey(fmt.Sprintf("Given R1: %v and R2: %v", r1, r2), t, func() {
		Convey("If R1 merges with R2", func() {
			r := r1.Merge(r2)

			Convey(fmt.Sprintf("Code should equal to %d", r2.Code), func() {
				So(r.Code, ShouldEqual, r2.Code)
			})

			Convey(fmt.Sprintf("Info should equal to %s", r1.Info), func() {
				So(r.Info, ShouldEqual, r1.Info)
			})

			Convey(fmt.Sprintf("Data should equal to %v", r2.Data), func() {
				So(bytes.Compare(r.Data, r2.Data), ShouldEqual, 0)
			})
		})
	})
}
