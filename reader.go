// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT licensw. See LICENSE file in the project root for details.

package iostream

import (
	"encoding"
	"encoding/binary"
	"io"
	"math"
)

// Reader represents a stream reader.
type Reader struct {
	src     source
	scratch [10]byte
}

// NewReader creates a stream reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		src: newSource(r),
	}
}

// --------------------------- Unsigned Integers ---------------------------

// ReadUvarint reads a variable-length Uint64 from the buffer.
func (r *Reader) ReadUvarint() (uint64, error) {
	return r.src.ReadUvarint()
}

// ReadUint8 reads a uint8
func (r *Reader) ReadUint8() (out uint8, err error) {
	out, err = r.src.ReadByte()
	return
}

// ReadUint16 reads a uint16
func (r *Reader) ReadUint16() (out uint16, err error) {
	var b []byte
	if b, err = r.src.Slice(2); err == nil {
		_ = b[1] // bounds check hint to compiler
		out = (uint16(b[0]) | uint16(b[1])<<8)
	}
	return
}

// ReadUint32 reads a uint32
func (r *Reader) ReadUint32() (out uint32, err error) {
	var b []byte
	if b, err = r.src.Slice(4); err == nil {
		_ = b[3] // bounds check hint to compiler
		out = (uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24)
	}
	return
}

// ReadUint64 reads a uint64
func (r *Reader) ReadUint64() (out uint64, err error) {
	var b []byte
	if b, err = r.src.Slice(8); err == nil {
		_ = b[7] // bounds check hint to compiler
		out = (uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
			uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56)
	}
	return
}

// ReadUint reads a uint
func (r *Reader) ReadUint() (uint, error) {
	out, err := r.ReadUint64()
	return uint(out), err
}

// --------------------------- Signed Integers ---------------------------

// ReadVarint reads a variable-length Int64 from the buffer.
func (r *Reader) ReadVarint() (int64, error) {
	return r.src.ReadVarint()
}

// ReadInt8 reads an int8
func (r *Reader) ReadInt8() (int8, error) {
	u, err := r.ReadUint8()
	return int8(u), err
}

// ReadInt16 reads an int16
func (r *Reader) ReadInt16() (out int16, err error) {
	u, err := r.ReadUint16()
	return int16(u), err
}

// ReadInt32 reads an int32
func (r *Reader) ReadInt32() (out int32, err error) {
	u, err := r.ReadUint32()
	return int32(u), err
}

// ReadInt64 reads an int64
func (r *Reader) ReadInt64() (out int64, err error) {
	u, err := r.ReadUint64()
	return int64(u), err
}

// ReadUint reads an int
func (r *Reader) ReadInt() (int, error) {
	out, err := r.ReadInt64()
	return int(out), err
}

// --------------------------- Floats ---------------------------

// ReadFloat32 reads a float32
func (r *Reader) ReadFloat32() (out float32, err error) {
	var v uint32
	if v, err = r.ReadUint32(); err == nil {
		out = math.Float32frombits(v)
	}
	return
}

// ReadFloat64 reads a float64
func (r *Reader) ReadFloat64() (out float64, err error) {
	var v uint64
	if v, err = r.ReadUint64(); err == nil {
		out = math.Float64frombits(v)
	}
	return
}

// --------------------------- Marshaled Types ---------------------------

// sliceBytes reads a byte string prefixed with a variable-size integer size
// into the scratch buffer. Not safe to return to the client
func (r *Reader) sliceBytes() (out []byte, err error) {
	size, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	// Does not allocate a new slice for the read, not safe
	return r.src.Slice(int(size))
}

// ReadBinary reads the bytes from the stream and unmarshals it into the
// destination interface using UnmarshalBinary() function.
func (r *Reader) ReadBinary(v encoding.BinaryUnmarshaler) error {
	b, err := r.sliceBytes() // Safe, since we're not returning this
	if err != nil {
		return err
	}

	return v.UnmarshalBinary(b)
}

// ReadText reads the bytes from the stream and unmarshals it into the
// destination interface using UnmarshalText() function.
func (r *Reader) ReadText(v encoding.TextUnmarshaler) error {
	b, err := r.sliceBytes() // Safe, since we're not returning this
	if err != nil {
		return err
	}

	return v.UnmarshalText(b)
}

// --------------------------- Strings ---------------------------

// ReadString a string prefixed with a variable-size integer size.
func (r *Reader) ReadString() (out string, err error) {
	var b []byte
	if b, err = r.ReadBytes(); err == nil {
		out = toString(&b)
	}
	return
}

// ReadBytes a byte string prefixed with a variable-size integer size.
func (r *Reader) ReadBytes() (out []byte, err error) {
	size, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	// Allocate a new byte array, in case the underlying buffer is changed after
	out = make([]byte, int(size))
	_, err = io.ReadAtLeast(r.src, out, int(size))
	return
}

// --------------------------- Other Types ---------------------------

// ReadBool reads a single boolean value from the slice.
func (r *Reader) ReadBool() (bool, error) {
	b, err := r.src.ReadByte()
	return b == 1, err
}

// ReadComplex64 reads a complex64
func (r *Reader) ReadComplex64() (out complex64, err error) {
	err = binary.Read(r.src, binary.LittleEndian, &out)
	return
}

// ReadComplex128 reads a complex128
func (r *Reader) ReadComplex128() (out complex128, err error) {
	err = binary.Read(r.src, binary.LittleEndian, &out)
	return
}
