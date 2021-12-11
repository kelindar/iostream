// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package iostream

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"unsafe"
)

// MaxVarintLenN is the maximum length of a varint-encoded N-bit integer.
const (
	maxVarintLen64 = 10 * 7
)

var overflow = errors.New("binary: varint overflows a 64-bit integer")

// source represents a required contract for a decoder to work properly
type source interface {
	io.Reader
	io.ByteReader
	Slice(n int) (buffer []byte, err error)
	ReadUvarint() (uint64, error)
	ReadVarint() (int64, error)
}

// newSource figures out the most efficient source to use for the provided type
func newSource(r io.Reader) source {
	switch v := r.(type) {
	case nil:
		return newSliceSource(nil)
	case *bytes.Buffer:
		return newSliceSource(v.Bytes())
	case *sliceSource:
		return v
	default:
		rdr, ok := r.(source)
		if !ok {
			rdr = newStreamSource(r)
		}
		return rdr
	}
}

// --------------------------- Slice Reader ---------------------------

// sliceSource implements a source that only reads from a slice
type sliceSource struct {
	buffer []byte
	offset int64 // current reading index
}

// newSliceSource returns a new source reading from b.
func newSliceSource(b []byte) *sliceSource {
	return &sliceSource{b, 0}
}

// Read implements the io.Reader interface.
func (r *sliceSource) Read(b []byte) (n int, err error) {
	if r.offset >= int64(len(r.buffer)) {
		return 0, io.EOF
	}

	n = copy(b, r.buffer[r.offset:])
	r.offset += int64(n)
	return
}

// ReadByte implements the io.ByteReader interface.
func (r *sliceSource) ReadByte() (byte, error) {
	if r.offset >= int64(len(r.buffer)) {
		return 0, io.EOF
	}

	b := r.buffer[r.offset]
	r.offset++
	return b, nil
}

// Slice selects a sub-slice of next bytes. This is similar to Read() but does not
// actually perform a copy, but simply uses the underlying slice (if available) and
// returns a sub-slice pointing to the same array. Since this requires access
// to the underlying data, this is only available for our default source.
func (r *sliceSource) Slice(n int) ([]byte, error) {
	if r.offset+int64(n) > int64(len(r.buffer)) {
		return nil, io.EOF
	}

	cur := r.offset
	r.offset += int64(n)
	return r.buffer[cur:r.offset], nil
}

// ReadUvarint reads an encoded unsigned integer from r and returns it as a uint64.
func (r *sliceSource) ReadUvarint() (uint64, error) {
	var x uint64
	for s := 0; s < maxVarintLen64; s += 7 {
		if r.offset >= int64(len(r.buffer)) {
			return 0, io.EOF
		}

		b := r.buffer[r.offset]
		r.offset++
		if b < 0x80 {
			if s == maxVarintLen64-7 && b > 1 {
				return x, overflow
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
	}
	return x, overflow
}

// ReadVarint reads an encoded signed integer from r and returns it as an int64.
func (r *sliceSource) ReadVarint() (int64, error) {
	ux, err := r.ReadUvarint() // ok to continue in presence of error
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, err
}

// --------------------------- Stream Reader ---------------------------

// streamSource represents a source implementation for a generic source (i.e. streams)
type streamSource struct {
	io.Reader
	io.ByteReader
	scratch []byte
}

// newStreamSource returns a new stream source
func newStreamSource(r io.Reader) *streamSource {
	src := &streamSource{
		Reader: r,
	}

	// If we can already read byte at a time, we're done
	if br, ok := r.(io.ByteReader); ok {
		src.ByteReader = br
		return src
	}

	// If stream doesn't have a byte reader, wrap it with a buffered reader
	buffered := bufio.NewReader(r)
	src.Reader = buffered
	src.ByteReader = buffered
	return src
}

// Slice selects a sub-slice of next bytes.
func (r *streamSource) Slice(n int) ([]byte, error) {
	if len(r.scratch) < n {
		r.scratch = make([]byte, capacityFor(uint(n+1)))
	}

	// Read from the stream into our scratch buffer
	_, err := io.ReadAtLeast(r.Reader, r.scratch[:n], n)
	return r.scratch[:n], err
}

// ReadUvarint reads an encoded unsigned integer from r and returns it as a uint64.
func (r *streamSource) ReadUvarint() (uint64, error) {
	return binary.ReadUvarint(r)
}

// ReadVarint reads a variable-length Int64 from the buffer.
func (r *streamSource) ReadVarint() (int64, error) {
	return binary.ReadVarint(r)
}

// --------------------------- Convert Funcs ---------------------------

// toString converts byte slice to a string without allocating.
func toString(b *[]byte) string {
	return *(*string)(unsafe.Pointer(b))
}

// toBytes converts a string to a byte slice without allocating.
func toBytes(v string) (b []byte) {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&v))
	byteHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	byteHeader.Data = strHeader.Data

	l := len(v)
	byteHeader.Len = l
	byteHeader.Cap = l
	return
}

// capacityFor computes the next power of 2 for a given index
func capacityFor(v uint) int {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return int(v)
}
