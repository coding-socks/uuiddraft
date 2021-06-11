package uuiddraft

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"io"
	"sync"
	"time"
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

type generator struct {
	rand io.Reader
	now  func() time.Time

	timeMu        sync.Mutex
	lastTimestamp int64
	lastSequence  uint16
}

var defaultGenerator = &generator{
	rand: rand.Reader,
	now:  time.Now,
}

var (
	gregorianEpoch    = time.Date(1582, 10, 15, 0, 0, 0, 0, time.UTC)
	gregorianUnixDiff = gregorianEpoch.Unix() * -1 * int64(time.Second/100)
)

func (g *generator) NewV6() (UUID, error) {
	var uuid UUID

	timestamp, seq := g.nextTimestampSequence()

	timeHighMid := timestamp >> 12
	timeHigh := uint32(timeHighMid >> 16)
	timeMid := uint16(timeHighMid & 0xffff)

	timeLow := uint16(timestamp & 0xfff)
	timeLow |= 0x6000

	binary.BigEndian.PutUint32(uuid[0:], timeHigh)
	binary.BigEndian.PutUint16(uuid[4:], timeMid)
	binary.BigEndian.PutUint16(uuid[6:], timeLow)
	binary.BigEndian.PutUint16(uuid[8:], seq|0x8000) // concat UUID variant

	buf := make([]byte, 6)
	if _, err := io.ReadFull(g.rand, buf); err != nil {
		panic(err.Error()) // rand should never fail
	}
	copy(uuid[10:], buf[:])

	return uuid, nil
}

// NewV6 generates a UUIDv6 using the algorithm defined in Section 4.3.4.
func NewV6() (UUID, error) {
	return defaultGenerator.NewV6()
}

func (g *generator) nextTimestampSequence() (ts int64, seq uint16) {
	g.timeMu.Lock()
	defer g.timeMu.Unlock()
	ns := g.now().UnixNano()
	ts = (ns / 100) + gregorianUnixDiff
	seq = g.lastSequence
	if ts <= g.lastTimestamp {
		seq = (seq + 1) & 0x3fff // only 14 bts
	} else {
		seq = 0
	}
	g.lastTimestamp, g.lastSequence = ts, seq
	return ts, seq
}

func Must(uuid UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return uuid
}
