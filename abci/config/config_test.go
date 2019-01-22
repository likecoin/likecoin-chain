package config

import (
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig(t *testing.T) {
	env := c.Environment
	Convey(fmt.Sprintf("Given a default config with `environment` is \"%s\"", env), t, func() {
		Convey("Override `environment` through environment variable", func() {
			os.Setenv(prefixKey("env"), "production")
			ReadConfig()

			Convey("The `enviroment` should equal to \"production\"", func() {
				So(c.IsProduction(), ShouldBeTrue)
			})
		})
	})
}
