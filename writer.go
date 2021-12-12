// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT licensw. See LICENSE file in the project root for details.

package iostream

import (
	"encoding"
	"io"
	"math"
)

// Writer represents a stream writer.
type Writer struct {
	scratch [10]byte
	out     io.Writer
	offset  int64
}

// NewWriter creates a new stream writer.
func NewWriter(out io.Writer) *Writer {
	if w, ok := out.(*Writer); ok {
		return w
	}

	return &Writer{
		out: out,
	}
}

// Reset resets the writer and makes it ready to be reused.
func (w *Writer) Reset(out io.Writer) {
	w.out = out
	w.offset = 0
}

// Offset returns the number of bytes written through this writer.
func (w *Writer) Offset() int64 {
	return w.offset
}

// --------------------------- io.Writer ---------------------------

// Write implements io.Writer interface by simply writing into the underlying
// souurce.
func (w *Writer) Write(p []byte) (int, error) {
	n, err := w.out.Write(p)
	w.offset += int64(n)
	return n, err
}

// Write writes the contents of p into the buffer.
func (w *Writer) write(p []byte) error {
	n, err := w.out.Write(p)
	w.offset += int64(n)
	return err
}

// Close closes the writer's underlying stream and return its error. If the
// underlying stream is not an io.Closer, it is a no-op.
func (w *Writer) Close() error {
	if closer, ok := w.out.(io.Closer); ok {
		return closer.Close()
	}
	return nil
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
	return w.write(w.scratch[:(i + 1)])
}

// WriteUint writes a Uint
func (w *Writer) WriteUint(v uint) error {
	return w.WriteUint64(uint64(v))
}

// WriteUint8 writes a Uint8
func (w *Writer) WriteUint8(v uint8) error {
	w.scratch[0] = v
	return w.write(w.scratch[:1])
}

// WriteUint16 writes a Uint16
func (w *Writer) WriteUint16(v uint16) error {
	w.scratch[0] = byte(v)
	w.scratch[1] = byte(v >> 8)
	return w.write(w.scratch[:2])
}

// WriteUint32 writes a Uint32
func (w *Writer) WriteUint32(v uint32) error {
	w.scratch[0] = byte(v)
	w.scratch[1] = byte(v >> 8)
	w.scratch[2] = byte(v >> 16)
	w.scratch[3] = byte(v >> 24)
	return w.write(w.scratch[:4])
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
	return w.write(w.scratch[:8])
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
	return w.write(w.scratch[:(i + 1)])
}

// WriteInt writes an int
func (w *Writer) WriteInt(v int) error {
	return w.WriteUint64(uint64(v))
}

// WriteInt8 writes an int8
func (w *Writer) WriteInt8(v int8) error {
	return w.WriteUint8(uint8(v))
}

// WriteInt16 writes an int16
func (w *Writer) WriteInt16(v int16) error {
	return w.WriteUint16(uint16(v))
}

// WriteInt32 writes an int32
func (w *Writer) WriteInt32(v int32) error {
	return w.WriteUint32(uint32(v))
}

// WriteInt64 writes an int64
func (w *Writer) WriteInt64(v int64) error {
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

// WriteSelf uses the provider io.WriterTo in order to write the data into
// the destination writer.
func (w *Writer) WriteSelf(v io.WriterTo) error {
	_, err := v.WriteTo(w)
	return err
}

// --------------------------- Strings ---------------------------

// WriteString writes a string prefixed with a variable-size integer.
func (w *Writer) WriteString(v string) error {
	if err := w.WriteUvarint(uint64(len(v))); err != nil {
		return err
	}
	return w.write(toBytes(v))
}

// WriteBytes writes a byte slice prefixed with a variable-size integer.
func (w *Writer) WriteBytes(v []byte) error {
	if err := w.WriteUvarint(uint64(len(v))); err != nil {
		return err
	}
	return w.write(v)
}

// --------------------------- Other Types ---------------------------

// WriteBool writes a single boolean value into the buffer
func (w *Writer) WriteBool(v bool) error {
	w.scratch[0] = 0
	if v {
		w.scratch[0] = 1
	}
	return w.write(w.scratch[:1])
}
