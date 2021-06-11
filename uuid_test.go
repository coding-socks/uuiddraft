package uuiddraft

import (
	"bytes"
	"crypto/rand"
	"regexp"
	"testing"
	"time"
)

func Test_generator_NewV6(t *testing.T) {
	t.Run("structure", func(t *testing.T) {
		g := generator{
			rand: rand.Reader,
			now:  time.Now,
		}
		id, err := g.NewV6()
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := id.Version(), 0b0110; got != want {
			t.Errorf("Version() = %v, want %v", got, want)
		}
		if got, want := id.Variant(), 0b10; got != want {
			t.Errorf("Variant() = %b, want %b", got, want)
		}
		if got, want := id.String(), `(?i)[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`; !regexp.MustCompile(want).MatchString(got) {
			t.Errorf("String() = %v, want to match %v", got, want)
		}
	})
	t.Run("segment", func(t *testing.T) {
		times := []time.Time{
			time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(1970, 1, 1, 0, 0, 0, 50, time.UTC),
			time.Date(1970, 1, 1, 0, 0, 0, 100, time.UTC),
		}
		g := generator{
			rand: bytes.NewReader([]byte{
				255, 255, 255, 255, 255, 255,
				1, 1, 1, 1, 1, 1,
				128, 128, 128, 128, 128, 128,
			}),
			now: func() time.Time {
				var x time.Time
				x, times = times[0], times[1:]
				return x
			},
		}
		tests := []string{
			"1b21dd21-3814-6000-8000-ffffffffffff",
			"1b21dd21-3814-6000-8001-010101010101",
			"1b21dd21-3814-6001-8000-808080808080",
		}
		for i := range tests {
			id, err := g.NewV6()
			if err != nil {
				t.Error(err)
				return
			}
			if got, want := id.String(), tests[i]; got != want {
				t.Errorf("String() = %v, want %v", got, want)
			}
		}
	})
}
