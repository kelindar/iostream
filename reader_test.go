// Copyright (c) Roman Atachiants and contributors. All rights reserved.
// Licensed under the MIT licensw. See LICENSE file in the project root for details.

package iostream

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStreamRead(t *testing.T) {
	for n, tc := range Fixtures {
		assertRead(t, n, tc.Decode, tc.Buffer, tc.Value)
	}
}

func TestSliceRead(t *testing.T) {
	for n, tc := range Fixtures {
		assertSliceRead(t, n, tc.Decode, tc.Buffer, tc.Value)
	}
}

func TestStreamReadShortBuffer(t *testing.T) {
	for n, tc := range Fixtures {
		for size := 0; size < len(tc.Buffer); size++ {
			assertRead(t, n, tc.Decode, tc.Buffer[:size], nil)
		}
	}
}

// assertRead asserts a single read operation
func assertRead(t *testing.T, name string, fn func(*Reader) (interface{}, error), input []byte, expect interface{}) {
	assertReadN(t, name, fn, input, expect, 99999)
}

// assertReadN asserts a single read operation
func assertReadN(t *testing.T, name string, fn func(*Reader) (interface{}, error), input []byte, expect interface{}, max int) {
	if max > len(input) {
		max = len(input)
	}

	t.Run(name, func(t *testing.T) {
		msg := fmt.Sprintf("write %v", name)
		src := newNetworkSource(input[:max])
		rdr := NewReader(src)
		out, err := fn(rdr)

		// Failure case, must have an error
		if expect == nil {
			assert.Error(t, err, msg)
			return
		}

		// Successfully encoded, check the output
		assert.NoError(t, err, msg)
		assert.Equal(t, expect, out, msg)
	})
}

// assertSliceRead asserts a single read operation from a slice source
func assertSliceRead(t *testing.T, name string, fn func(*Reader) (interface{}, error), input []byte, expect interface{}) {
	t.Run(name, func(t *testing.T) {
		msg := fmt.Sprintf("write %v", name)
		src := bytes.NewBuffer(input)
		rdr := NewReader(src)
		out, err := fn(rdr)

		// Failure case, must have an error
		if expect == nil {
			assert.Error(t, err, msg)
			return
		}

		// Successfully encoded, check the output
		assert.NoError(t, err, msg)
		assert.Equal(t, expect, out, msg)
	})
}
