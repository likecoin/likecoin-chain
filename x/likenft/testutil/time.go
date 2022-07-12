package testutil

import (
	"time"
)

func MustParseTime(layout string, value string) *time.Time {
	time, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}

	return &time
}
