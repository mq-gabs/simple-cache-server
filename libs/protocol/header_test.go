package protocol

import (
	"encoding/binary"
	"errors"
	"testing"
)

func TestDecodeHeader(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		want    *Header
		wantErr error
	}{
		{
			name: "valid header",
			buf: func() []byte {
				buf := make([]byte, HeaderSize)

				binary.BigEndian.PutUint16(buf[0:MagicSize], uint16(SCAS))

				buf[offsetVersion] = byte(Version(0x01))
				buf[offsetCommand] = byte(Command(0x02))
				buf[offsetFlags] = byte(Flag(0x03))

				binary.BigEndian.PutUint32(
					buf[offsetPayloadLength:offsetPayloadLength+PayloadLengthSize],
					uint32(128),
				)

				return buf
			}(),
			want: &Header{
				Version:       Version(0x01),
				Command:       Command(0x02),
				Flags:         Flag(0x03),
				PayloadLength: PayloadLength(128),
			},
		},
		{
			name:    "invalid header size too small",
			buf:     make([]byte, HeaderSize-1),
			wantErr: ErrInvalidHeaderSize,
		},
		{
			name:    "invalid header size too large",
			buf:     make([]byte, HeaderSize+1),
			wantErr: ErrInvalidHeaderSize,
		},
		{
			name: "invalid magic",
			buf: func() []byte {
				buf := make([]byte, HeaderSize)

				binary.BigEndian.PutUint16(buf[0:MagicSize], 0xFFFF)

				return buf
			}(),
			wantErr: ErrInvalidMagic,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeHeader(tt.buf)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}

				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got == nil {
				t.Fatal("expected header, got nil")
			}

			if *got != *tt.want {
				t.Fatalf("expected %+v, got %+v", tt.want, got)
			}
		})
	}
}

func TestDecodeHeader_MaxPayloadLength(t *testing.T) {
	buf := make([]byte, HeaderSize)

	binary.BigEndian.PutUint16(buf[0:MagicSize], uint16(SCAS))

	binary.BigEndian.PutUint32(
		buf[offsetPayloadLength:offsetPayloadLength+PayloadLengthSize],
		0xFFFFFFFF,
	)

	header, err := DecodeHeader(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if header.PayloadLength != PayloadLength(0xFFFFFFFF) {
		t.Fatalf("unexpected payload length: %v", header.PayloadLength)
	}
}
