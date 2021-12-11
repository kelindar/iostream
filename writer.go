// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT licensw. See LICENSE file in the project root for details.

package iostream

import (
	"encoding"
	"encoding/binary"
	"io"
	"math"
)

// Writer represents a stream writer.
type Writer struct {
	scratch [10]byte
	out     io.Writer
}

// NewWriter creates a new stream writer.
func NewWriter(out io.Writer) *Writer {
	return &Writer{
		out: out,
	}
}

// Reset resets the writer and makes it ready to be reused.
func (w *Writer) Reset(out io.Writer) {
	w.out = out
}

// Write writes the contents of p into the buffer.
func (w *Writer) Write(p []byte) (err error) {
	_, err = w.out.Write(p)
	return
}

// --------------------------- Unsigned Integers ---------------------------

// WriteUvarint writes a variable size unsigned integer
func (w *Writer) WriteUvarint(x uint64) error {
	i := 0
	for x >= 0x80 {
		w.scratch[i] = byte(x) | 0x80
		x >>= 7
		i++
	}
	w.scratch[i] = byte(x)
	return w.Write(w.scratch[:(i + 1)])
}

// WriteUint writes a Uint
func (w *Writer) WriteUint(v uint) error {
	return w.WriteUint64(uint64(v))
}

// WriteUint8 writes a Uint8
func (w *Writer) WriteUint8(v uint8) error {
	w.scratch[0] = byte(v)
	return w.Write(w.scratch[:1])
}

// WriteUint16 writes a Uint16
func (w *Writer) WriteUint16(v uint16) error {
	w.scratch[0] = byte(v)
	w.scratch[1] = byte(v >> 8)
	return w.Write(w.scratch[:2])
}

// WriteUint32 writes a Uint32
func (w *Writer) WriteUint32(v uint32) error {
	w.scratch[0] = byte(v)
	w.scratch[1] = byte(v >> 8)
	w.scratch[2] = byte(v >> 16)
	w.scratch[3] = byte(v >> 24)
	return w.Write(w.scratch[:4])
}

// WriteUint64 writes a Uint64
func (w *Writer) WriteUint64(v uint64) error {
	w.scratch[0] = byte(v)
	w.scratch[1] = byte(v >> 8)
	w.scratch[2] = byte(v >> 16)
	w.scratch[3] = byte(v >> 24)
	w.scratch[4] = byte(v >> 32)
	w.scratch[5] = byte(v >> 40)
	w.scratch[6] = byte(v >> 48)
	w.scratch[7] = byte(v >> 56)
	return w.Write(w.scratch[:8])
}

// --------------------------- Signed Integers ---------------------------

// WriteVarint writes a variable size signed integer
func (w *Writer) WriteVarint(v int64) error {
	x := uint64(v) << 1
	if v < 0 {
		x = ^x
	}

	i := 0
	for x >= 0x80 {
		w.scratch[i] = byte(x) | 0x80
		x >>= 7
		i++
	}
	w.scratch[i] = byte(x)
	return w.Write(w.scratch[:(i + 1)])
}

// WriteUint writes an int
func (w *Writer) WriteInt(v uint) error {
	return w.WriteUint64(uint64(v))
}

// WriteInt8 writes an int8
func (w *Writer) WriteInt8(v int8) error {
	return w.WriteUint8(uint8(v))
}

// WriteInt16 writes an int16
func (w *Writer) WriteInt16(v uint16) error {
	return w.WriteUint16(uint16(v))
}

// WriteInt32 writes an int32
func (w *Writer) WriteInt32(v uint32) error {
	return w.WriteUint32(uint32(v))
}

// WriteInt64 writes an int64
func (w *Writer) WriteInt64(v uint64) error {
	return w.WriteUint64(uint64(v))
}

// --------------------------- Floats ---------------------------

// WriteFloat32 a 32-bit floating point number
func (w *Writer) WriteFloat32(v float32) error {
	return w.WriteUint32(math.Float32bits(v))
}

// WriteFloat64 a 64-bit floating point number
func (w *Writer) WriteFloat64(v float64) error {
	return w.WriteUint64(math.Float64bits(v))
}

// --------------------------- Marshaled Types ---------------------------

// WriteBinary marshals the type to its binary representation and writes it
// downstream, prefixed with its size as a variable-size integer.
func (w *Writer) WriteBinary(v encoding.BinaryMarshaler) error {
	out, err := v.MarshalBinary()
	if err == nil {
		err = w.WriteBytes(out)
	}
	return err
}

// WriteText marshals the type to its text representation and writes it
// downstream, prefixed with its size as a variable-size integer.
func (w *Writer) WriteText(v encoding.TextMarshaler) error {
	out, err := v.MarshalText()
	if err == nil {
		err = w.WriteBytes(out)
	}
	return err
}

// --------------------------- Strings ---------------------------

// WriteString writes a string prefixed with a variable-size integer.
func (w *Writer) WriteString(v string) error {
	if err := w.WriteUvarint(uint64(len(v))); err != nil {
		return err
	}
	return w.Write(toBytes(v))
}

// WriteBytes writes a byte slice prefixed with a variable-size integer.
func (w *Writer) WriteBytes(v []byte) error {
	if err := w.WriteUvarint(uint64(len(v))); err != nil {
		return err
	}
	return w.Write(v)
}

// --------------------------- Other Types ---------------------------

// WriteBool writes a single boolean value into the buffer
func (w *Writer) WriteBool(v bool) error {
	w.scratch[0] = 0
	if v {
		w.scratch[0] = 1
	}
	return w.Write(w.scratch[:1])
}

// WriteComplex64 a 64-bit complex number
func (w *Writer) WriteComplex64(v complex64) error {
	return binary.Write(w.out, binary.LittleEndian, v)
}

// WriteComplex128 a 128-bit complex number
func (w *Writer) WriteComplex128(v complex128) error {
	return binary.Write(w.out, binary.LittleEndian, v)
}
