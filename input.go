package csslexer

import (
	"io"
)

// NOTE: We use rune instead of byte to handle Unicode characters properly.

// EOF is a special rune that represents the end of the input.
const EOF = rune(0)

// nullRune is an empty slice of runes, used to represent no runes.
var nullRune = []rune{}

// Input represents a stream of runes read from a source.
type Input struct {
	runes []rune // The runes in the input stream.
	pos   int    // The current position in the input stream.
	start int    // The start position of the current token being read.
	err   error  // Any error encountered while reading the input.
}

// NewInput creates a new Input instance from the given string.
func NewInput(input string) *Input {
	return NewInputRunes([]rune(input))
}

// NewInputRunes creates a new Input instance from the given slice of runes.
func NewInputRunes(runes []rune) *Input {
	return &Input{
		runes: runes,
		pos:   0,
		start: 0,
		err:   nil,
	}
}

// NewInputBytes creates a new Input instance from the given byte slice.
func NewInputBytes(input []byte) *Input {
	return NewInputRunes([]rune(string(input)))
}

// NewInputReader creates a new Input instance from the given io.Reader.
func NewInputReader(r io.Reader) *Input {
	var b []byte

	if r != nil {
		var err error
		b, err = io.ReadAll(r)
		if err != nil {
			return &Input{
				runes: nullRune,
				pos:   0,
				start: 0,
				err:   err,
			}
		}
	}

	return NewInputBytes(b)
}

// PeekErr checks if there is an error at the current position plus the specified offset.
func (z *Input) PeekErr(pos int) error {
	if z.err != nil {
		return z.err
	}

	if z.pos+pos >= len(z.runes) {
		return io.EOF
	}

	return nil
}

// Err returns the error at the current position.
func (z *Input) Err() error {
	return z.PeekErr(0)
}
