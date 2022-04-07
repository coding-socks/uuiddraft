package uuiddraft

import (
	"bytes"
	"crypto/rand"
	"regexp"
	"testing"
	"time"
)

func Test_generatorV6_New(t *testing.T) {
	t.Run("structure", func(t *testing.T) {
		g := generatorV6{
			rand: rand.Reader,
			now:  time.Now,
		}
		id, err := g.New()
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := id.Version(), 6; got != want {
			t.Errorf("Version() = %v, want %v", got, want)
		}
		if got, want := id.Variant(), 0b10; got != want {
			t.Errorf("Variant() = %b, want %b", got, want)
		}
		if got, want := id.String(), `(?i)[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}`; !regexp.MustCompile(want).MatchString(got) {
			t.Errorf("String() = %v, want to match %v", got, want)
		}
	})
	t.Run("vectors", func(t *testing.T) {
		g := generatorV6{
			rand: bytes.NewReader([]byte{
				0xf3, 0xc8, // seq
				0x9e, 0x6b, 0xde, 0xce, 0xd8, 0x46, // node
			}),
			now: func() time.Time {
				zone := time.FixedZone("GMT-5", -5*60*60)
				return time.Date(2022, 2, 22, 14, 22, 22, 0, zone)
			},
			lastSequence: -1,
		}
		id, err := g.New()
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := id.String(), "1ec9414c-232a-6b00-b3c8-9e6bdeced846"; got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
	t.Run("segment", func(t *testing.T) {
		times := []time.Time{
			time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(1970, 1, 1, 0, 0, 0, 50, time.UTC),
			time.Date(1970, 1, 1, 0, 0, 0, 100, time.UTC),
		}
		g := generatorV6{
			rand: bytes.NewReader([]byte{
				0x11, 0x11, // seq
				0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				0x02, 0x02, 0x02, 0x02, 0x02, 0x02,
				0x03, 0x03, 0x03, 0x03, 0x03, 0x03,
			}),
			now: func() time.Time {
				var x time.Time
				x, times = times[0], times[1:]
				return x
			},
			lastSequence: -1,
		}
		tests := []string{
			"1b21dd21-3814-6000-9111-010101010101",
			"1b21dd21-3814-6000-9112-020202020202",
			"1b21dd21-3814-6001-9112-030303030303",
		}
		for i := range tests {
			id, err := g.New()
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
