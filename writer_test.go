// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT licensw. See LICENSE file in the project root for details.

package iostream

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var Fixtures = map[string]struct {
	Encode func(*Writer) error
	Decode func(*Reader) (interface{}, error)
	Buffer []byte
	Value  interface{}
}{
	"uvarint": {
		Encode: func(w *Writer) error { return w.WriteUvarint(0x1111111111111111) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadUvarint() },
		Buffer: []byte{0x91, 0xa2, 0xc4, 0x88, 0x91, 0xa2, 0xc4, 0x88, 0x11},
		Value:  uint64(0x1111111111111111),
	},
	"uint": {
		Encode: func(w *Writer) error { return w.WriteUint(0x1111111111111111) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadUint() },
		Buffer: []byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
		Value:  uint(0x1111111111111111),
	},
	"uint8": {
		Encode: func(w *Writer) error { return w.WriteUint8(0x11) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadUint8() },
		Buffer: []byte{0x11},
		Value:  uint8(0x11),
	},
	"uint16": {
		Encode: func(w *Writer) error { return w.WriteUint16(0x1111) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadUint16() },
		Buffer: []byte{0x11, 0x11},
		Value:  uint16(0x1111),
	},
	"uint32": {
		Encode: func(w *Writer) error { return w.WriteUint32(0x11111111) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadUint32() },
		Buffer: []byte{0x11, 0x11, 0x11, 0x11},
		Value:  uint32(0x11111111),
	},
	"uint64": {
		Encode: func(w *Writer) error { return w.WriteUint64(0x1111111111111111) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadUint64() },
		Buffer: []byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
		Value:  uint64(0x1111111111111111),
	},
	"varint": {
		Encode: func(w *Writer) error { return w.WriteVarint(0x1111111111111111) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadVarint() },
		Buffer: []byte{0xa2, 0xc4, 0x88, 0x91, 0xa2, 0xc4, 0x88, 0x91, 0x22},
		Value:  int64(0x1111111111111111),
	},
	"varint-negative": {
		Encode: func(w *Writer) error { return w.WriteVarint(-0x10) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadVarint() },
		Buffer: []byte{0x1f},
		Value:  int64(-0x10),
	},
	"int": {
		Encode: func(w *Writer) error { return w.WriteInt(0x1111111111111111) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadInt() },
		Buffer: []byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
		Value:  int(0x1111111111111111),
	},
	"int8": {
		Encode: func(w *Writer) error { return w.WriteInt8(0x11) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadInt8() },
		Buffer: []byte{0x11},
		Value:  int8(0x11),
	},
	"int16": {
		Encode: func(w *Writer) error { return w.WriteInt16(0x1111) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadInt16() },
		Buffer: []byte{0x11, 0x11},
		Value:  int16(0x1111),
	},
	"int32": {
		Encode: func(w *Writer) error { return w.WriteInt32(0x11111111) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadInt32() },
		Buffer: []byte{0x11, 0x11, 0x11, 0x11},
		Value:  int32(0x11111111),
	},
	"int64": {
		Encode: func(w *Writer) error { return w.WriteInt64(0x1111111111111111) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadInt64() },
		Buffer: []byte{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11},
		Value:  int64(0x1111111111111111),
	},
	"float32": {
		Encode: func(w *Writer) error { return w.WriteFloat32(0x11) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadFloat32() },
		Buffer: []byte{0x0, 0x0, 0x88, 0x41},
		Value:  float32(0x11),
	},
	"float64": {
		Encode: func(w *Writer) error { return w.WriteFloat64(0x11) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadFloat64() },
		Buffer: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x31, 0x40},
		Value:  float64(0x11),
	},
	"string": {
		Encode: func(w *Writer) error { return w.WriteString("hello") },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadString() },
		Buffer: []byte{0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f},
		Value:  "hello",
	},
	"bytes": {
		Encode: func(w *Writer) error { return w.WriteBytes([]byte("hello")) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadBytes() },
		Buffer: []byte{0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f},
		Value:  []byte("hello"),
	},
	"bool": {
		Encode: func(w *Writer) error { return w.WriteBool(true) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadBool() },
		Buffer: []byte{0x1},
		Value:  true,
	},
	"complex64": {
		Encode: func(w *Writer) error { return w.WriteComplex64(complex64(1)) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadComplex64() },
		Buffer: []byte{0x0, 0x0, 0x80, 0x3f, 0x0, 0x0, 0x0, 0x0},
		Value:  complex64(1),
	},
	"complex128": {
		Encode: func(w *Writer) error { return w.WriteComplex128(complex128(1)) },
		Decode: func(r *Reader) (interface{}, error) { return r.ReadComplex128() },
		Buffer: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf0, 0x3f, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		Value:  complex128(1),
	},
	"time-binary": {
		Encode: func(w *Writer) error { return w.WriteBinary(time.Unix(60, 0).UTC()) },
		Decode: func(r *Reader) (interface{}, error) {
			var out time.Time
			err := r.ReadBinary(&out)
			return out, err
		},
		Buffer: []byte{0xf, 0x1, 0x0, 0x0, 0x0, 0xe, 0x77, 0x91, 0xf7, 0x3c, 0x0, 0x0, 0x0, 0x0, 0xff, 0xff},
		Value:  time.Unix(60, 0).UTC(),
	},
	"time-text": {
		Encode: func(w *Writer) error { return w.WriteText(time.Unix(60, 0).UTC()) },
		Decode: func(r *Reader) (interface{}, error) {
			var out time.Time
			err := r.ReadText(&out)
			return out, err
		},
		Buffer: []byte{0x14, 0x31, 0x39, 0x37, 0x30, 0x2d, 0x30, 0x31, 0x2d, 0x30, 0x31, 0x54, 0x30, 0x30, 0x3a, 0x30, 0x31, 0x3a, 0x30, 0x30, 0x5a},
		Value:  time.Unix(60, 0).UTC(),
	},
}

func TestWrite(t *testing.T) {
	for n, tc := range Fixtures {
		assertWrite(t, n, tc.Encode, tc.Buffer)
	}
}

func TestWriteFailuresString(t *testing.T) {
	assertWriteN(t, "string-err", func(w *Writer) error {
		return w.WriteString("hello")
	}, nil, 0)

	assertWriteN(t, "bytes-err", func(w *Writer) error {
		return w.WriteBytes([]byte("hello"))
	}, nil, 0)
}

// assertWrite asserts a single write operation
func assertWrite(t *testing.T, name string, fn func(*Writer) error, expect []byte) {
	assertWriteN(t, name, fn, expect, 99999)
}

// assertWriteMax asserts a single write operation
func assertWriteN(t *testing.T, name string, fn func(*Writer) error, expect []byte, max int) {
	t.Run(name, func(t *testing.T) {
		msg := fmt.Sprintf("write %v", name)
		dst := newLimitWriter(max)
		wrt := NewWriter(dst)

		wrt.Reset(dst)
		err := fn(wrt)

		// Failure case, must have an error
		if expect == nil {
			assert.Error(t, err, msg)
			return
		}

		// Successfully encoded, check the output
		assert.NoError(t, err, msg)
		assert.Equal(t, expect, dst.buffer.Bytes(), msg)
	})
}
