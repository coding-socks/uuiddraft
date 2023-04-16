package uuiddraft

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

var (
	DefaultV6Generator = NewV6Generator()
	DefaultV7Generator = NewV7Generator()
	DefaultV8Generator = NewV8Generator()
)

// A UUID is a 128 bit (16 byte) Universal Unique Identifier as defined in RFC
// 4122.
type UUID [16]byte

// Variant returns the value of the UUID's variant segment.
func (u UUID) Variant() int {
	return int(u[8] >> 6)
}

// Version returns the value of the UUID's version segment.
func (u UUID) Version() int {
	return int(u[6] >> 4)
}

func (u UUID) String() string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])

	return string(buf)
}

func Equal(a, b UUID) bool {
	return bytes.Equal(a[:], b[:])
}

var (
	nilUUID = UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	maxUUID = UUID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

func IsNil(uuid UUID) bool {
	return uuid == nilUUID
}

func IsMax(uuid UUID) bool {
	return uuid == maxUUID
}

// clockSequence is a 14 bits counter
type clockSequence uint16

func randomClockSequence() clockSequence {
	b := make([]byte, 2)
	if _, err := rand.Read(b); err != nil {
		panic(err) // theoretically should never happen
	}
	return ((clockSequence(b[0]) << 8) | clockSequence(b[1])) & clockSequence(0xdfff)
}

func (cs clockSequence) Incr() clockSequence {
	return (cs + 1) & clockSequence(0xdfff)
}

var (
	gregEpoch = time.Date(1582, time.October, 15, 0, 0, 0, 0, time.UTC)
)

type V6Generator struct {
	now  func() time.Time
	rand io.Reader

	node     []byte
	cs       clockSequence
	mu       sync.Mutex
	prevTime time.Time
}

func NewV6Generator() *V6Generator {
	return &V6Generator{
		now:  time.Now,
		cs:   randomClockSequence(),
		rand: rand.Reader,
	}
}

// gregFormat returns a 60-bit timestamp represented by UTC as
// a count of 100- nanosecond intervals since 00:00:00.00, 15 October 1582.
func gregFormat(t time.Time) int64 {
	return (t.Unix()-gregEpoch.Unix())*1e7 + int64(t.Nanosecond()-gregEpoch.Nanosecond())/1e2
}

// Generate generates a UUID Version 6 based on
// https://www.ietf.org/archive/id/draft-ietf-uuidrev-rfc4122bis-00.html#name-uuid-version-6
func (g *V6Generator) Read(id *UUID) error {
	g.mu.Lock()
	if len(g.node) == 0 {
		// init arbitrary node ID when it is not set.
		r := make([]byte, 6)
		if _, err := g.rand.Read(r); err != nil {
			g.mu.Unlock()
			return fmt.Errorf("could not initialise node ID: %w", err) // fail fast
		}
		g.node = r
	}
	n := g.now()
	if n.Before(g.prevTime) {
		g.cs = g.cs.Incr()
	}
	g.prevTime = n
	cs := g.cs
	g.mu.Unlock()

	t := gregFormat(n)
	binary.BigEndian.PutUint32(id[:4], uint32(t>>28))  // time_high
	binary.BigEndian.PutUint16(id[4:6], uint16(t>>12)) // time_mid
	binary.BigEndian.PutUint16(id[6:8], uint16(t))     // time_low_and_version
	binary.BigEndian.PutUint16(id[8:10], uint16(cs))   // clk_seq_hi_res + clk_seq_low
	copy(id[10:], g.node)                              // node 0-5
	id[6] = (id[6] & 0x0f) | 0x60                      // ver
	id[8] = (id[8] & 0x3f) | 0x80                      // var
	return nil
}

// V6 reads a UUID from DefaultV6Generator.
func V6() (UUID, error) {
	var id UUID
	if err := DefaultV6Generator.Read(&id); err != nil {
		return UUID{}, err
	}
	return id, nil
}

type V7Generator struct {
	now  func() time.Time
	rand io.Reader
}

func NewV7Generator() *V7Generator {
	return &V7Generator{
		now:  time.Now,
		rand: rand.Reader,
	}
}

// Generate generates a UUID Version 7 based on
// https://www.ietf.org/archive/id/draft-ietf-uuidrev-rfc4122bis-00.html#name-uuid-version-7
func (g V7Generator) Read(id *UUID) error {
	r := make([]byte, 10)
	if _, err := g.rand.Read(r); err != nil {
		return err // fail fast
	}
	um := g.now().UnixMilli()
	binary.BigEndian.PutUint32(id[:4], uint32(um>>16)) // unix_ts_ms
	binary.BigEndian.PutUint16(id[4:6], uint16(um))    // unix_ts_ms
	copy(id[6:], r)                                    // rand
	id[6] = (id[6] & 0x0f) | 0x70                      // ver
	id[8] = (id[8] & 0x3f) | 0x80                      // var
	return nil
}

// V7 reads a UUID from DefaultV7Generator.
func V7() (UUID, error) {
	var id UUID
	if err := DefaultV7Generator.Read(&id); err != nil {
		return UUID{}, err
	}
	return id, nil
}

type V8Generator struct {
	r io.Reader
}

func NewV8Generator() *V8Generator {
	return &V8Generator{r: rand.Reader}
}

// Generate generates a UUID Version 8 based on
// https://www.ietf.org/archive/id/draft-ietf-uuidrev-rfc4122bis-00.html#name-uuid-version-8.
func (g V8Generator) Read(id *UUID) error {
	b := make([]byte, 16)
	if _, err := g.r.Read(b); err != nil {
		return err
	}
	copy(id[:], b)
	id[6] = (id[6] & 0x0f) | 0x80 // ver
	id[8] = (id[8] & 0x3f) | 0x80 // var
	return nil
}

// V8 reads a UUID from DefaultV8Generator.
func V8() (UUID, error) {
	var id UUID
	if err := DefaultV8Generator.Read(&id); err != nil {
		return UUID{}, err
	}
	return id, nil
}

func Must(uuid UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return uuid
}

var ErrInvalidUUID = errors.New("invalid UUID")

// Parse parses the "hex-and-dash" string representation of a UUID.
//
// Format: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
func Parse(raw string) (UUID, error) {
	if len(raw) != 36 {
		return UUID{}, ErrInvalidUUID
	}
	if raw[8] != '-' && raw[13] != '-' && raw[18] != '-' && raw[23] != '-' {
		return UUID{}, ErrInvalidUUID
	}
	src := raw[:8] + raw[9:13] + raw[14:18] + raw[19:23] + raw[24:]
	id := UUID{}
	_, err := hex.Decode(id[:], []byte(src))
	return id, err
}
