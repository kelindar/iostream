// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT licensw. See LICENSE file in the project root for details.

package iostream

import (
	"encoding"
	"io"
	"math"
)

// Reader represents a stream reader.
type Reader struct {
	src     source
	scratch [10]byte
}

// NewReader creates a stream reader.
func NewReader(src io.Reader) *Reader {
	if r, ok := src.(*Reader); ok {
		return r
	}

	return &Reader{
		src: newSource(src),
	}
}

// Offset returns the number of bytes read through this reader.
func (r *Reader) Offset() int64 {
	return r.src.Offset()
}

// --------------------------- io.Reader ---------------------------

// Read implements io.Reader interface by simply calling the Read method on
// the underlying stream.
func (r *Reader) Read(p []byte) (n int, err error) {
	return r.src.Read(p)
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

// ReadUint8s reads an array of uint8s
func (r *Reader) ReadUint8s() ([]uint8, error) {
	length, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	out := make([]uint8, length)
	for i := 0; i < int(length); i++ {
		if out[i], err = r.ReadUint8(); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ReadUint16s reads an array of uint16s
func (r *Reader) ReadUint16s() ([]uint16, error) {
	length, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	out := make([]uint16, length)
	for i := 0; i < int(length); i++ {
		if out[i], err = r.ReadUint16(); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ReadUint32s reads an array of uint32s
func (r *Reader) ReadUint32s() ([]uint32, error) {
	length, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	out := make([]uint32, length)
	for i := 0; i < int(length); i++ {
		if out[i], err = r.ReadUint32(); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ReadUint64s reads an array of uint64s
func (r *Reader) ReadUint64s() ([]uint64, error) {
	length, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	out := make([]uint64, length)
	for i := 0; i < int(length); i++ {
		if out[i], err = r.ReadUint64(); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ReadUints reads an array of uints
func (r *Reader) ReadUints() ([]uint, error) {
	length, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	out := make([]uint, length)
	for i := 0; i < int(length); i++ {
		if out[i], err = r.ReadUint(); err != nil {
			return nil, err
		}
	}

	return out, nil
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

// ReadInt8s reads an array of int8s
func (r *Reader) ReadInt8s() ([]int8, error) {
	length, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	out := make([]int8, length)
	for i := 0; i < int(length); i++ {
		if out[i], err = r.ReadInt8(); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ReadInt16s reads an array of int16s
func (r *Reader) ReadInt16s() ([]int16, error) {
	length, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	out := make([]int16, length)
	for i := 0; i < int(length); i++ {
		if out[i], err = r.ReadInt16(); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ReadInt32s reads an array of int32s
func (r *Reader) ReadInt32s() ([]int32, error) {
	length, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	out := make([]int32, length)
	for i := 0; i < int(length); i++ {
		if out[i], err = r.ReadInt32(); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ReadInt64s reads an array of int64s
func (r *Reader) ReadInt64s() ([]int64, error) {
	length, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	out := make([]int64, length)
	for i := 0; i < int(length); i++ {
		if out[i], err = r.ReadInt64(); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ReadUints reads an array of uints
func (r *Reader) ReadInts() ([]int, error) {
	length, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	out := make([]int, length)
	for i := 0; i < int(length); i++ {
		if out[i], err = r.ReadInt(); err != nil {
			return nil, err
		}
	}

	return out, nil
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

// ReadFloat32s reads an array of float32s
func (r *Reader) ReadFloat32s() ([]float32, error) {
	length, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	out := make([]float32, length)
	for i := 0; i < int(length); i++ {
		if out[i], err = r.ReadFloat32(); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// ReadFloat64s reads an array of float64s
func (r *Reader) ReadFloat64s() ([]float64, error) {
	length, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	out := make([]float64, length)
	for i := 0; i < int(length); i++ {
		if out[i], err = r.ReadFloat64(); err != nil {
			return nil, err
		}
	}

	return out, nil
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

// ReadSelf uses the provider io.ReaderFrom in order to read the data from
// the source reader.
func (r *Reader) ReadSelf(v io.ReaderFrom) error {
	_, err := v.ReadFrom(r)
	return err
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

// ReadStrings reads an array of strings
func (r *Reader) ReadStrings() ([]string, error) {
	length, err := r.ReadUvarint()
	if err != nil {
		return nil, err
	}

	out := make([]string, length)
	for i := 0; i < int(length); i++ {
		if out[i], err = r.ReadString(); err != nil {
			return nil, err
		}
	}

	return out, nil
}

// --------------------------- Other Types ---------------------------

// ReadRange reads the length of the array from the underlying stream and
// calls a callback function on each element of that array.
func (r *Reader) ReadRange(fn func(i int, r *Reader) error) error {
	length, err := r.ReadUvarint()
	if err != nil {
		return err
	}

	for i := 0; i < int(length); i++ {
		if err := fn(i, r); err != nil {
			return err
		}
	}
	return nil
}

// ReadBool reads a single boolean value from the slice.
func (r *Reader) ReadBool() (bool, error) {
	b, err := r.src.ReadByte()
	return b == 1, err
}
