package uuiddraft

import (
	"crypto/rand"
	"regexp"
	"testing"
	"time"
)

func Test_generator_NewV6(t *testing.T) {
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
}
