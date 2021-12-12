// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package iostream

import (
	"bytes"
	"io"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertString(t *testing.T) {
	v := "hi there"

	b := toBytes(v)
	assert.NotEmpty(t, b)
	assert.Equal(t, v, string(b))

	o := toString(&b)
	assert.NotEmpty(t, b)
	assert.Equal(t, v, o)
}

func TestNewSource(t *testing.T) {
	assert.IsType(t, &sliceSource{}, newSource(nil))
	assert.IsType(t, &sliceSource{}, newSource(&sliceSource{}))
	assert.IsType(t, &streamSource{}, newSource(&networkSource{}))
	assert.IsType(t, &streamSource{}, newStreamSource(&bytes.Buffer{}))
}

func TestReadUvarintEOF(t *testing.T) {
	input := []byte{0x91, 0xa2, 0xc4, 0x88, 0x91, 0xa2, 0xc4, 0x88, 0x11}
	for size := 0; size < len(input)-1; size++ {
		src := newSource(bytes.NewBuffer(input[:size]))
		_, err := src.ReadUvarint()
		assert.Error(t, err)
	}
}

func TestReadUvarintOverflow(t *testing.T) {
	input := []byte{0x91, 0xa2, 0xc4, 0x88, 0x91, 0xa2, 0xc4, 0x88, 0xa2, 0xc4, 0x88, 0x11}
	for size := 0; size < len(input)-1; size++ {
		src := newSource(bytes.NewBuffer(input[:size]))
		_, err := src.ReadUvarint()
		assert.Error(t, err)
	}
}

func TestReadUvarintUnexpectedOverflow(t *testing.T) {
	input := []byte{0x91, 0xa2, 0xc4, 0x88, 0x91, 0xa2, 0xc4, 0x88, 0x91, 0x11}
	src := newSource(bytes.NewBuffer(input))
	_, err := src.ReadUvarint()
	assert.Error(t, err)
}

func TestReadEOF(t *testing.T) {
	src := newSource(bytes.NewBuffer(nil))
	_, err := src.Read([]byte{})
	assert.Error(t, err)
}

func TestReadByteEOF(t *testing.T) {
	src := newSource(bytes.NewBuffer(nil))
	_, err := src.ReadByte()
	assert.Error(t, err)
}

func TestSliceEOF(t *testing.T) {
	src := newSliceSource([]byte{})
	_, err := src.Slice(10)
	assert.Error(t, err)
}

// --------------------------- Fake Network Source ---------------------------

type networkSource struct {
	r io.Reader
}

func newNetworkSource(data []byte) io.Reader {
	return &networkSource{
		r: bytes.NewBuffer(data),
	}
}

func (s *networkSource) Read(p []byte) (n int, err error) {
	return s.r.Read(p)
}

// --------------------------- Fake Limited Writer ---------------------------

type limitWriter struct {
	buffer *bytes.Buffer
	value  uint32
	Limit  int
}

func newLimitWriter(limit int) *limitWriter {
	return &limitWriter{
		buffer: bytes.NewBuffer(nil),
		Limit:  limit,
	}
}

func (w *limitWriter) Write(p []byte) (int, error) {
	if n := atomic.AddUint32(&w.value, uint32(len(p))); int(n) > w.Limit {
		return 0, io.ErrShortBuffer
	}

	return w.buffer.Write(p)
}

func (w *limitWriter) Close() error {
	return nil
}

// --------------------------- Self Reader/Writer ---------------------------

type person struct {
	Name string
}

func (p *person) WriteTo(dst io.Writer) (int64, error) {
	w := NewWriter(dst)
	err := w.WriteString(p.Name)
	return w.Offset(), err
}

func (p *person) ReadFrom(src io.Reader) (int64, error) {
	r := NewReader(src)
	name, err := r.ReadString()
	p.Name = name
	return r.Offset(), err
}
