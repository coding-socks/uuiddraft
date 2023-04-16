package uuiddraft

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

func TestIsNil(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			uuid: UUID{},
			want: true,
		},
		{
			uuid: UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: true,
		},
		{
			uuid: UUID{0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f},
			want: false,
		},
		{
			uuid: UUID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNil(tt.uuid); got != tt.want {
				t.Errorf("IsNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsMax(t *testing.T) {
	tests := []struct {
		name string
		uuid UUID
		want bool
	}{
		{
			uuid: UUID{},
			want: false,
		},
		{
			uuid: UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			want: false,
		},
		{
			uuid: UUID{0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f},
			want: false,
		},
		{
			uuid: UUID{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsMax(tt.uuid); got != tt.want {
				t.Errorf("IsMax() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestV6Generator_Read(t *testing.T) {
	// https://www.ietf.org/archive/id/draft-ietf-uuidrev-rfc4122bis-00.html#name-example-of-a-uuidv6-value
	// -----------------------------------------------
	// field                 bits    value
	// -----------------------------------------------
	// time_high              32     0x1EC9414C
	// time_mid               16     0x232A
	// time_low_and_version   16     0x6B00
	// clk_seq_hi_res          8     0xB3
	// clock_seq_low           8     0xC8
	// node                   48     0x9E6BDECED846
	// -----------------------------------------------
	// total                 128
	// -----------------------------------------------
	// final_hex: 1EC9414C-232A-6B00-B3C8-9E6BDECED846
	fmt.Println()
	g := V6Generator{
		now: func() time.Time {
			b, _ := hex.DecodeString("01EC9414C232AB00")
			ns := binary.BigEndian.Uint64(b) * 100
			return time.Unix(0, int64(ns)+gregEpoch.UnixNano())
		},
		cs:   0x33c8,
		rand: bytes.NewReader([]byte{0x9e, 0x6b, 0xde, 0xce, 0xd8, 0x46}),
	}
	var got UUID
	err := g.Read(&got)
	if (err != nil) != false {
		t.Errorf("Read() error = %v", err)
		return
	}
	want := Must(Parse("1EC9414C-232A-6B00-B3C8-9E6BDECED846"))
	if !Equal(got, want) {
		t.Errorf("Read() got = %v, want %v", got, want)
	}
}

func TestV7Generator_Read(t *testing.T) {
	// https://www.ietf.org/archive/id/draft-ietf-uuidrev-rfc4122bis-00.html#name-example-of-a-uuidv7-value
	// -------------------------------
	// field      bits    value
	// -------------------------------
	// unix_ts_ms   48    0x17F22E279B0
	// ver           4    0x7
	// rand_a       12    0xCC3
	// var           2    b10
	// rand_b       62    b01, 0x8C4DC0C0C07398F
	// -------------------------------
	// total       128
	// -------------------------------
	// final: 017F22E2-79B0-7CC3-98C4-DC0C0C07398F
	fmt.Println()
	g := V7Generator{
		now: func() time.Time {
			return time.UnixMilli(1645557742000)
		},
		rand: bytes.NewReader([]byte{
			0x0c, 0xc3,
			0x18, 0xc4, 0xdc, 0x0c, 0x0c, 0x07, 0x39, 0x8f,
		}),
	}
	var got UUID
	err := g.Read(&got)
	if (err != nil) != false {
		t.Errorf("Read() error = %v", err)
		return
	}
	want := Must(Parse("017F22E2-79B0-7CC3-98C4-DC0C0C07398F"))
	if !Equal(got, want) {
		t.Errorf("Read() got = %v, want %v", got, want)
	}
}

func TestV8Generator_Read(t *testing.T) {
	// https://www.ietf.org/archive/id/draft-ietf-uuidrev-rfc4122bis-00.html#name-example-of-a-uuidv8-value
	// -------------------------------
	// field      bits    value
	// -------------------------------
	// custom_a     48    0x320C3D4DCC00
	// ver           4    0x8
	// custom_b     12    0x75B
	// var           2    b10
	// custom_c     62    b00, 0xEC932D5F69181C0
	// -------------------------------
	// total       128
	// -------------------------------
	// final: 320C3D4D-CC00-875B-8EC9-32D5F69181C0
	fmt.Println()
	g := V8Generator{
		r: bytes.NewReader([]byte{
			0x32, 0x0c, 0x3d, 0x4d, 0xcc, 0x00,
			0x07, 0x5b,
			0x0e, 0xc9, 0x32, 0xd5, 0xf6, 0x91, 0x81, 0xc0,
		}),
	}
	var got UUID
	err := g.Read(&got)
	if (err != nil) != false {
		t.Errorf("Read() error = %v", err)
		return
	}
	want := Must(Parse("320C3D4D-CC00-875B-8EC9-32D5F69181C0"))
	if !Equal(got, want) {
		t.Errorf("Read() got = %v, want %v", got, want)
	}
}
